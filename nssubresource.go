package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"encoding/json"
	"github.com/ghodss/yaml"

	"k8s.io/apimachinery/pkg/runtime"
)

type SubResource interface {
	MakeWatchInterface(nsname string) (watch.Interface, error)
	GetName(obj *runtime.Object) string
	GetCreationTimestamp(obj *runtime.Object) metav1.Time
}

type SimpleFs struct {
	SubResource
	watch.Interface
	objects []*runtime.Object
}

func NewSimpleFs(subresource SubResource) *SimpleFs {
	return &SimpleFs{
		SubResource: subresource,
		objects:     []*runtime.Object{},
	}
}

func (me *SimpleFs) Watch(nsname string) {
	// watch XXXs C_UD
	wi, err := me.MakeWatchInterface(nsname)
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

			switch ev.Type {
			case watch.Added:
				log.Printf("%s is Addes to %s", ev.Object, nsname)
				me.add(&ev.Object)
			case watch.Modified:
				log.Printf("%s@%s is Modified ", ev.Object, nsname)
				me.update(&ev.Object)
			case watch.Deleted:
				log.Printf("%s@%s is Killed", ev.Object, nsname)
				me.remove(&ev.Object)
			}
		}
		log.Printf("watching is finished %v\n", wi)
	}()
}

func (me *SimpleFs) Stop() {
	me.Interface.Stop()
}

func (me *SimpleFs) add(obj *runtime.Object) {
	me.objects = append(me.objects, obj)
}

func (me *SimpleFs) remove(obj *runtime.Object) {
	newlist := me.objects
	for i, object := range me.objects {
		if me.GetName(obj) == me.GetName(object) {
			newlist = append(me.objects[:i], me.objects[i+1:]...)
			break
		}
	}
	me.objects = newlist
}

func (me *SimpleFs) update(obj *runtime.Object) {
	for i, object := range me.objects {
		if me.GetName(obj) == me.GetName(object) {
			me.objects[i] = obj
			return
		}
	}
}

func (me *SimpleFs) get(name string) (*runtime.Object, error) {
	for _, object := range me.objects {
		if me.GetName(object) != name {
			continue
		}
		return object, nil
	}
	return nil, fmt.Errorf(" \"%s\" is not found.", name)
}

func (me *SimpleFs) GetAttr(names []string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("XXX GetAttr: %s\n", names)
	switch {
	case len(names) == 1:

		objectname := strings.TrimSuffix(names[0], ".yaml")

		object, err := me.get(objectname)
		if err != nil {
			log.Printf("%s is not found\n", objectname)
			return nil, fuse.ENOENT
		}

		ctime := uint64(me.GetCreationTimestamp(object).Unix())
		attr := &fuse.Attr{
			Mode:  fuse.S_IFREG | 0644,
			Ctime: ctime,
			Mtime: ctime,
			Atime: ctime,
		}

		// TODO caching
		if data, err := me.GetYaml(objectname); err != nil {
			attr.Size = 0
		} else {
			attr.Size = uint64(len(data))
		}

		return attr, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *SimpleFs) OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("XXX OpenDir: %s\n", names)
	if len(names) == 0 {
		c = []fuse.DirEntry{}
		for _, object := range me.objects {
			c = append(c, fuse.DirEntry{Name: me.GetName(object) + ".yaml", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *SimpleFs) Open(names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	name := names[0]

	if !strings.HasSuffix(name, ".yaml") {
		return nil, fuse.ENOENT
	} else {
		objectname := strings.TrimSuffix(name, ".yaml")

		yaml, err := me.GetYaml(objectname)
		if err != nil {
			return nil, fuse.ENOENT
		}

		return nodefs.NewDataFile([]byte(yaml)), fuse.OK
	}
}

// TODO Update, Delete

func (me *SimpleFs) GetYaml(name string) ([]byte, error) {
	object, err := me.get(name)
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}
