package main

import (
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
	}
	return nil, fuse.ENOENT
}

func (me *NamespaceFs) OpenDir(names []string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	// TODO
	return nil, fuse.ENOENT
}

func (me *NamespaceFs) Open(name string, names []string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {

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
}

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
