package fuse

import (
	"flag"
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func TestMain() {
	go func() {
		nsDir := Fs.root.(*namespacesDir)
		wi, err := k8s.Clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
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
				nsDir.AddNamespace(&ev.Object)
			case watch.Modified:
				// Update
				nsDir.UpdateNamespace(&ev.Object)
			case watch.Deleted:
				// Delete
				nsDir.DeleteNamespace(&ev.Object)
			}
		}
	}()
	Serve(flag.Arg(0))
}
