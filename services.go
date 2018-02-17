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

type ServicesFs struct {
	watch.Interface
	services []*corev1.Service
}

func NewServicesFs() *ServicesFs {
	return &ServicesFs{
		services: []*corev1.Service{},
	}
}

func (me *ServicesFs) Watch(nsname string) {
	// watch Services C_UD
	wi, err := clientset.CoreV1().Services(nsname).Watch(metav1.ListOptions{})
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
			po, ok := ev.Object.(*corev1.Service)
			if !ok {
				break
			}

			switch ev.Type {
			case watch.Added:
				log.Printf("%s is Addes to %s", po.GetName(), po.GetNamespace())
				me.addService(po)
			case watch.Modified:
				log.Printf("%s@%s is Modified ", po.GetName(), po.GetNamespace())
				me.updateService(po)
			case watch.Deleted:
				log.Printf("%s@%s is Killed", po.GetName(), po.GetNamespace())
				me.removeService(po)
			}
		}
		log.Printf("watching is finished %v\n", wi)
	}()
}

func (me *ServicesFs) Stop() {
	me.Interface.Stop()
}

func (me *ServicesFs) addService(po *corev1.Service) {
	me.services = append(me.services, po)
}

func (me *ServicesFs) removeService(po *corev1.Service) {
	newlist := me.services
	for i, service := range me.services {
		if service.GetName() == po.GetName() {
			newlist = append(me.services[:i], me.services[i+1:]...)
			break
		}
	}
	me.services = newlist
}

func (me *ServicesFs) updateService(po *corev1.Service) {
	for i, service := range me.services {
		if service.GetName() == po.GetName() {
			me.services[i] = po
			return
		}
	}
}

func (me *ServicesFs) getService(name string) (*corev1.Service, error) {
	for _, service := range me.services {
		if service.GetName() != name {
			continue
		}
		return service, nil
	}
	return nil, fmt.Errorf("Service \"%s\" is not found.", name)
}

func (me *ServicesFs) GetAttr(names []string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("Service GetAttr: %s\n", names)
	switch {
	case len(names) == 1:

		servicename := strings.TrimSuffix(names[0], ".yaml")

		service, err := me.getService(servicename)
		if err != nil {
			return nil, fuse.ENOENT
		}

		ctime := uint64(service.GetCreationTimestamp().Unix())
		attr := &fuse.Attr{
			Mode:  fuse.S_IFREG | 0644,
			Ctime: ctime,
			Mtime: ctime,
			Atime: ctime,
		}

		// TODO caching
		if data, err := me.GetYaml(servicename); err != nil {
			attr.Size = 0
		} else {
			attr.Size = uint64(len(data))
		}

		return attr, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *ServicesFs) OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("Service OpenDir: %s %s\n", names, me.services)
	if len(names) == 0 {
		c = []fuse.DirEntry{}
		for _, service := range me.services {
			c = append(c, fuse.DirEntry{Name: service.GetName() + ".yaml", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *ServicesFs) Open(names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	name := names[0]

	if !strings.HasSuffix(name, ".yaml") {
		return nil, fuse.ENOENT
	} else {
		servicename := strings.TrimSuffix(name, ".yaml")

		yaml, err := me.GetYaml(servicename)
		if err != nil {
			return nil, fuse.ENOENT
		}

		return nodefs.NewDataFile([]byte(yaml)), fuse.OK
	}
}

// TODO Update, Delete

func (me *ServicesFs) GetYaml(name string) ([]byte, error) {
	service, err := me.getService(name)
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}
