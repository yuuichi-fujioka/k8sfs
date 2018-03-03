package main

import (
	"log"

	"github.com/yuuichi-fujioka/k8sfs/fuse"
	"github.com/yuuichi-fujioka/k8sfs/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func main() {
	testMain()
}

func testMain() {
	nsdir := fuse.NewNamespacesDir()
	clientset := k8s.GenClientSetFromFlags()

	wi, err := clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ch := wi.ResultChan()

	for {
		ev, ok := <-ch
		if !ok {
			return
		}

		switch ev.Type {
		case watch.Added:
			log.Println("Added")
			nsdir.AddNamespace(&ev.Object)
		}
	}
}
