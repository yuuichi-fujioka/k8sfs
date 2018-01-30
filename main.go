// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// A Go mirror of libfuse's hello.c

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type HelloFs struct {
	pathfs.FileSystem
	Namespaces *corev1.NamespaceList
}

func (me *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr: %s\n", name)
	names := strings.Split(name, "/")
	switch {
	case name == "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	case len(names) == 1:
		if strings.HasSuffix(names[0], ".meta") {
			return &fuse.Attr{
				Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
			}, fuse.OK
		} else {
			return &fuse.Attr{
				Mode: fuse.S_IFDIR | 0755,
			}, fuse.OK
		}
	}
	return nil, fuse.ENOENT
}

func (me *HelloFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir: %s\n", name)
	if name == "" {
		c = []fuse.DirEntry{}
		for i := 0; i < len(me.Namespaces.Items); i++ {
			c = append(c, fuse.DirEntry{Name: me.Namespaces.Items[i].GetName(), Mode: fuse.S_IFDIR})
			c = append(c, fuse.DirEntry{Name: me.Namespaces.Items[i].GetName() + ".meta", Mode: fuse.S_IFREG})
		}
		return c, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *HelloFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("Open: %s\n", name)
	if !strings.HasSuffix(name, ".meta") {
		return nil, fuse.ENOENT
	}

	namespacename := strings.TrimSuffix(name, ".meta")
	for i := 0; i < len(me.Namespaces.Items); i++ {
		if me.Namespaces.Items[i].GetName() != namespacename {
			continue
		}
		return nodefs.NewDataFile([]byte(me.Namespaces.Items[i].String())), fuse.OK
	}

	return nil, fuse.ENOENT
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

	opts := metav1.ListOptions{}
	namespacelist, err := clientset.CoreV1().Namespaces().List(opts)
	if err != nil {
		panic(err.Error())
	}

	log.Printf("%d Namespaces are exist\n", len(namespacelist.Items))
	for i := 0; i < len(namespacelist.Items); i++ {
		log.Printf("%v\n", namespacelist.Items[i])
	}

	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}
	log.Printf("argments: %v\n", flag.Args())
	nfs := pathfs.NewPathNodeFs(&HelloFs{FileSystem: pathfs.NewDefaultFileSystem(), Namespaces: namespacelist}, nil)
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
