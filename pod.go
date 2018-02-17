package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"encoding/json"
	"github.com/ghodss/yaml"
)

type PodsFs struct {
	watch.Interface
	pods []*corev1.Pod
}

func NewPodsFs() *PodsFs {
	return &PodsFs{
		pods: []*corev1.Pod{},
	}
}

func (me *PodsFs) Watch(nsname string) {
	// watch Pods C_UD
	wi, err := clientset.CoreV1().Pods(nsname).Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	me.Interface = wi

	go func() {
		ch := me.Interface.ResultChan()
		log.Printf("start watching %v\n", wi)
		for {
			ev, ok := <-ch
			if !ok {
				break
			}
			po, ok := ev.Object.(*corev1.Pod)
			if !ok {
				break
			}

			switch ev.Type {
			case watch.Added:
				log.Printf("%s is Addes to %s", po.GetName(), po.GetNamespace())
				me.addPod(po)
			case watch.Modified:
				log.Printf("%s@%s is Modified ", po.GetName(), po.GetNamespace())
				me.updatePod(po)
			case watch.Deleted:
				log.Printf("%s@%s is Killed", po.GetName(), po.GetNamespace())
				me.removePod(po)
			}
		}
		log.Printf("watching is finished %v\n", wi)
	}()
}

func (me *PodsFs) Stop() {
	me.Interface.Stop()
}

func (me *PodsFs) addPod(po *corev1.Pod) {
	me.pods = append(me.pods, po)
}

func (me *PodsFs) removePod(po *corev1.Pod) {
	newlist := me.pods
	for i, pod := range me.pods {
		if pod.GetName() == po.GetName() {
			newlist = append(me.pods[:i], me.pods[i+1:]...)
			break
		}
	}
	me.pods = newlist
}

func (me *PodsFs) updatePod(po *corev1.Pod) {
	for i, pod := range me.pods {
		if pod.GetName() == po.GetName() {
			me.pods[i] = po
			return
		}
	}
}

func (me *PodsFs) getPod(name string) (*corev1.Pod, error) {
	for _, pod := range me.pods {
		if pod.GetName() != name {
			continue
		}
		return pod, nil
	}
	return nil, fmt.Errorf("Pod \"%s\" is not found.", name)
}

func (me *PodsFs) GetAttr(names []string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("Pod GetAttr: %s\n", names)
	switch {
	case len(names) == 1:

		podname := strings.TrimSuffix(names[0], ".yaml")

		pod, err := me.getPod(podname)
		if err != nil {
			return nil, fuse.ENOENT
		}

		ctime := uint64(pod.GetCreationTimestamp().Unix())
		attr := &fuse.Attr{
			Mode:  fuse.S_IFREG | 0644,
			Ctime: ctime,
			Mtime: ctime,
			Atime: ctime,
		}

		// TODO caching
		if data, err := me.GetYaml(podname); err != nil {
			attr.Size = 0
		} else {
			attr.Size = uint64(len(data))
		}

		return attr, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *PodsFs) OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("Pod OpenDir: %s\n", names)
	if len(names) == 0 {
		c = []fuse.DirEntry{}
		for _, pod := range me.pods {
			c = append(c, fuse.DirEntry{Name: pod.GetName() + ".yaml", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *PodsFs) Open(names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	name := names[0]

	if !strings.HasSuffix(name, ".yaml") {
		return nil, fuse.ENOENT
	} else {
		podname := strings.TrimSuffix(name, ".yaml")

		yaml, err := me.GetYaml(podname)
		if err != nil {
			return nil, fuse.ENOENT
		}

		return nodefs.NewDataFile([]byte(yaml)), fuse.OK
	}
}

// TODO Update, Delete

func (me *PodsFs) GetYaml(name string) ([]byte, error) {
	pod, err := me.getPod(name)
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(pod)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}
