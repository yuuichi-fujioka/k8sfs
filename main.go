package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	// "k8s.io/apimachinery/pkg/api/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

	k8sfs := K8sFs{FileSystem: pathfs.NewDefaultFileSystem(), Namespaces: []NamespaceFs{}}

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
