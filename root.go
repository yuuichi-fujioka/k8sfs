package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sFs struct {
	pathfs.FileSystem
	Namespaces []*NamespaceFs
}

func (me *K8sFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Printf("Mkdir: %s\n", name)
	names := strings.Split(name, "/")
	switch {
	case len(names) == 1:
		// Create a Namespace
		namespacename := names[0]
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespacename,
			},
		}
		_, err := clientset.CoreV1().Namespaces().Create(ns)
		if err != nil {
			log.Printf("error has occured. %v\n", err)
			return fuse.EIO
		}
		// TODO Event handling should be smart.
		for {
			_, err := me.GetNamespace(namespacename)
			if err == nil {
				return fuse.OK
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	return fuse.ENOSYS
}

func (me *K8sFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	log.Printf("Unlink: %s\n", name)
	names := strings.Split(name, "/")
	switch {
	case len(names) == 1:
		// Dlete a Namespace
		namespacename := strings.TrimSuffix(names[0], ".yaml")
		err := clientset.CoreV1().Namespaces().Delete(namespacename, &metav1.DeleteOptions{})
		if err != nil {
			log.Printf("error has occured. %v\n", err)
			return fuse.EIO
		}

		return fuse.OK
	}
	return fuse.ENOSYS
}

func (me *K8sFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr: %s\n", name)
	names := strings.Split(name, "/")
	switch {
	case name == "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	case len(names) >= 1:
		var attr *fuse.Attr
		var namespacename string
		if strings.HasSuffix(names[0], ".yaml") {
			attr = &fuse.Attr{Mode: fuse.S_IFREG | 0644}
			namespacename = strings.TrimSuffix(names[0], ".yaml")
		} else {
			attr = &fuse.Attr{Mode: fuse.S_IFDIR | 0755}
			namespacename = names[0]
		}

		namespace, err := me.GetNamespace(namespacename)
		if err != nil {
			return nil, fuse.ENOENT
		}
		attr, status := namespace.GetAttr(names[0], names[1:], context)
		return attr, status
	}
	return nil, fuse.ENOENT
}

func (me *K8sFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir: %s\n", name)
	if name == "" {
		c = []fuse.DirEntry{}
		for _, namespace := range me.Namespaces {
			c = append(c, fuse.DirEntry{Name: namespace.Namespace.GetName(), Mode: fuse.S_IFDIR})
			c = append(c, fuse.DirEntry{Name: namespace.Namespace.GetName() + ".yaml", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}

	names := strings.Split(name, "/")
	namespace, err := me.GetNamespace(names[0])
	if err != nil {
		return nil, fuse.ENOENT
	}
	c, code = namespace.OpenDir(names[1:], context)
	return
}

func (me *K8sFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("Open: %s\n", name)
	names := strings.Split(name, "/")
	if len(names) == 0 {
		return nil, fuse.ENOENT
	}
	var namespacename string
	if strings.HasSuffix(names[0], ".yaml") {
		namespacename = strings.TrimSuffix(names[0], ".yaml")
	} else {
		namespacename = names[0]
	}

	namespace, err := me.GetNamespace(namespacename)
	if err != nil {
		return nil, fuse.ENOENT
	}
	file, code = namespace.Open(names[0], names[1:], flags, context)
	return
}

func (me *K8sFs) AddNamespace(ns *corev1.Namespace) {
	nsfs := NewNamespaceFs(ns)
	nsfs.WatchAll()

	me.Namespaces = append(me.Namespaces, nsfs)
}

func (me *K8sFs) RemoveNamespace(ns *corev1.Namespace) {
	removedNsName := ns.GetName()
	newlist := me.Namespaces
	for i, namespace := range me.Namespaces {
		if namespace.Namespace.GetName() == removedNsName {
			namespace.StopAll()
			newlist = append(me.Namespaces[:i], me.Namespaces[i+1:]...)
			break
		}
	}
	me.Namespaces = newlist
}

func (me *K8sFs) UpdateNamespace(ns *corev1.Namespace) {
	for _, namespace := range me.Namespaces {
		if namespace.GetName() == ns.GetName() {
			namespace.Namespace = ns
			return
		}
	}
}

func (me *K8sFs) GetNamespace(name string) (*NamespaceFs, error) {
	for _, namespace := range me.Namespaces {
		if namespace.GetName() != name {
			continue
		}
		return namespace, nil
	}
	return nil, fmt.Errorf("Namespace \"%s\" is not found.", name)
}
