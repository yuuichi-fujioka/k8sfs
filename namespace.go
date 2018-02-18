package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"

	"encoding/json"
	"github.com/ghodss/yaml"
)

type NamespaceFs struct {
	corev1.Namespace
	PodsFs     NsChildFs
	ServicesFs NsChildFs
}

type NsChildFs interface {
	Watch(nsname string)
	Stop()
	GetAttr(names []string, context *fuse.Context) (*fuse.Attr, fuse.Status)
	OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status)
	Open(names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status)
}

func NewNamespaceFs(ns *corev1.Namespace) NamespaceFs {
	return NamespaceFs{
		Namespace:  *ns,
		PodsFs:     NewSimpleFs(&PodResource{}),
		ServicesFs: NewSimpleFs(&ServiceResource{}),
	}
}

func (me *NamespaceFs) getChildFs(name string) (NsChildFs, error) {
	switch name {
	case "po":
		return me.PodsFs, nil
	case "svc":
		return me.ServicesFs, nil
	}

	return nil, fmt.Errorf("%s is not supported yet", name)
}

func (me *NamespaceFs) WatchAll() {
	me.PodsFs.Watch(me.Name)
	me.ServicesFs.Watch(me.Name)
}

func (me *NamespaceFs) StopAll() {
	me.PodsFs.Stop()
	me.ServicesFs.Stop()
}

func (me *NamespaceFs) GetAttr(name string, names []string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr: %s\n", names)
	switch {
	case len(names) == 0:
		var attr *fuse.Attr
		if strings.HasSuffix(name, ".yaml") {
			attr = &fuse.Attr{Mode: fuse.S_IFREG | 0644}
		} else {
			attr = &fuse.Attr{Mode: fuse.S_IFDIR | 0755}
		}

		attr.Ctime = uint64(me.Namespace.GetCreationTimestamp().Unix())
		attr.Mtime = attr.Ctime
		attr.Atime = attr.Ctime

		// TODO caching
		if data, err := me.GetYaml(); err != nil {
			attr.Size = 0
		} else {
			attr.Size = uint64(len(data))
		}

		return attr, fuse.OK
	case len(names) == 1:
		Ctime := uint64(me.Namespace.GetCreationTimestamp().Unix())
		attr := &fuse.Attr{
			Mode:  fuse.S_IFDIR | 0755,
			Ctime: Ctime,
			Mtime: Ctime,
			Atime: Ctime,
			Size:  0,
		}

		return attr, fuse.OK
	default:
		cfs, err := me.getChildFs(names[0])
		if err != nil {
			return nil, fuse.ENOENT
		}
		attr, status := cfs.GetAttr(names[1:], context)
		return attr, status
	}
	return nil, fuse.ENOENT
}

func (me *NamespaceFs) OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir: %s\n", names)
	if len(names) == 0 {
		c = []fuse.DirEntry{
			fuse.DirEntry{Name: "po", Mode: fuse.S_IFDIR},
			// fuse.DirEntry{Name: "rs", Mode: fuse.S_IFDIR},
			// fuse.DirEntry{Name: "sa", Mode: fuse.S_IFDIR},
			// fuse.DirEntry{Name: "deploy", Mode: fuse.S_IFDIR},
			// fuse.DirEntry{Name: "ds", Mode: fuse.S_IFDIR},
			fuse.DirEntry{Name: "svc", Mode: fuse.S_IFDIR},
			// fuse.DirEntry{Name: "ing", Mode: fuse.S_IFDIR},
		}
		// TODO
		return c, fuse.OK
	}
	if len(names) == 1 {
		cfs, err := me.getChildFs(names[0])
		if err != nil {
			return nil, fuse.ENOENT
		}
		c, status := cfs.OpenDir(names[1:], context)
		return c, status
	}
	// TODO
	return nil, fuse.ENOENT
}

func (me *NamespaceFs) Open(name string, names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	switch {
	case len(names) == 0:

		if !strings.HasSuffix(name, ".yaml") {
			// TODO
			return nil, fuse.ENOENT
		} else {

			yaml, err := me.GetYaml()
			if err != nil {
				return nil, fuse.ENOENT
			}

			return nodefs.NewDataFile([]byte(yaml)), fuse.OK
		}
	default:
		cfs, err := me.getChildFs(names[0])
		if err != nil {
			return nil, fuse.ENOENT
		}
		data, status := cfs.Open(names[1:], flags, context)
		return data, status
	}
}

// TODO Update, Delete

func (me *NamespaceFs) GetYaml() ([]byte, error) {
	jsondata, err := json.Marshal(me.Namespace)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}
