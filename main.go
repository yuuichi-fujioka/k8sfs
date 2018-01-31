// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Go mirror of libfuse's hello.c

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	// "k8s.io/apimachinery/pkg/api/errors"
	"encoding/json"
	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sFs struct {
	pathfs.FileSystem
	Namespaces []corev1.Namespace
}

func (me *K8sFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr: %s\n", name)
	names := strings.Split(name, "/")
	switch {
	case name == "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	case len(names) == 1:
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
		attr.Ctime = uint64(namespace.GetCreationTimestamp().Unix())
		attr.Mtime = attr.Ctime
		attr.Atime = attr.Ctime

		// TODO caching
		if data, err := me.GetNamespaceYaml(namespacename); err != nil {
			attr.Size = 0
		} else {
			attr.Size = uint64(len(data))
		}

		return attr, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *K8sFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir: %s\n", name)
	if name == "" {
		c = []fuse.DirEntry{}
		for _, namespace := range me.Namespaces {
			c = append(c, fuse.DirEntry{Name: namespace.GetName(), Mode: fuse.S_IFDIR})
			c = append(c, fuse.DirEntry{Name: namespace.GetName() + ".yaml", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *K8sFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("Open: %s\n", name)
	if !strings.HasSuffix(name, ".yaml") {
		return nil, fuse.ENOENT
	}

	namespacename := strings.TrimSuffix(name, ".yaml")

	yaml, err := me.GetNamespaceYaml(namespacename)
	if err != nil {
		return nil, fuse.ENOENT
	}

	return nodefs.NewDataFile([]byte(yaml)), fuse.OK
}

func (me *K8sFs) AddNamespace(ns *corev1.Namespace) {
	me.Namespaces = append(me.Namespaces, *ns)
}

func (me *K8sFs) RemoveNamespace(ns *corev1.Namespace) {
	removedNsName := ns.GetName()
	newlist := me.Namespaces
	for i, namespace := range me.Namespaces {
		if namespace.GetName() == removedNsName {
			newlist = append(me.Namespaces[:i], me.Namespaces[i+1:]...)
			break
		}
	}
	me.Namespaces = newlist
}

func (me *K8sFs) UpdateNamespace(ns *corev1.Namespace) {
	for i, namespace := range me.Namespaces {
		if namespace.GetName() == ns.GetName() {
			me.Namespaces[i] = *ns
			return
		}
	}
}

func (me *K8sFs) GetNamespace(name string) (*corev1.Namespace, error) {
	for _, namespace := range me.Namespaces {
		if namespace.GetName() != name {
			continue
		}
		return &namespace, nil
	}
	return nil, fmt.Errorf("Namespace \"%s\" is not found.", name)
}

func (me *K8sFs) GetNamespaceYaml(name string) ([]byte, error) {
	namespace, err := me.GetNamespace(name)
	if err != nil {
		return nil, err
	}

	jsondata, err := json.Marshal(namespace)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

func main() {

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	log.Printf("argments: %v\n", flag.Args())

	k8sfs := K8sFs{FileSystem: pathfs.NewDefaultFileSystem(), Namespaces: []corev1.Namespace{}}

	go func() {
		wi, err := clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		ch := wi.ResultChan()
		for {
			ev := <-ch
			ns, ok := ev.Object.(*corev1.Namespace)
			if !ok {
				panic("???!?!??!?")
			}

			log.Println(ns)

			switch ev.Type {
			case watch.Added:
				k8sfs.AddNamespace(ns)
			case watch.Modified:
				k8sfs.UpdateNamespace(ns)
			case watch.Deleted:
				k8sfs.RemoveNamespace(ns)
			}
		}
	}()

	nfs := pathfs.NewPathNodeFs(&k8sfs, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
