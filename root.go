package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
)

type K8sFs struct {
	pathfs.FileSystem
	Namespaces []NamespaceFs
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
	me.Namespaces = append(me.Namespaces, NamespaceFs{Namespace: *ns})
}

func (me *K8sFs) RemoveNamespace(ns *corev1.Namespace) {
	removedNsName := ns.GetName()
	newlist := me.Namespaces
	for i, namespace := range me.Namespaces {
		if namespace.Namespace.GetName() == removedNsName {
			newlist = append(me.Namespaces[:i], me.Namespaces[i+1:]...)
			break
		}
	}
	me.Namespaces = newlist
}

func (me *K8sFs) UpdateNamespace(ns *corev1.Namespace) {
	for i, namespace := range me.Namespaces {
		if namespace.GetName() == ns.GetName() {
			me.Namespaces[i] = NamespaceFs{Namespace: *ns}
			return
		}
	}
}

func (me *K8sFs) GetNamespace(name string) (*NamespaceFs, error) {
	for _, namespace := range me.Namespaces {
		if namespace.GetName() != name {
			continue
		}
		return &namespace, nil
	}
	return nil, fmt.Errorf("Namespace \"%s\" is not found.", name)
}
