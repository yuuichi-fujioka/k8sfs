package fuse

import (
	"flag"
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func TestMain() {
	go watchAllNs()
	Serve(flag.Arg(0))
}

func watchAllNs() {
	nsDir := Fs.root.(*namespacesDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
			}

			ns, ok := ev.Object.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			switch ev.Type {
			case watch.Added:
				log.Println("ns/Added")
				nsDir.AddNamespace(ev.Object)
				dir := nsDir.GetDir(ns.Name + "/po")
				poDir := dir.(*podsDir)
				go watchPods(ns.Name, poDir)

			case watch.Modified:
				// Update
				nsDir.UpdateNamespace(ev.Object)
			case watch.Deleted:
				// Delete
				nsDir.DeleteNamespace(ev.Object)
			}
		}
	}
}

func watchPods(ns string, poDir *podsDir) {
	for {
		wi, err := k8s.Clientset.CoreV1().Pods(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
			}

			switch ev.Type {
			case watch.Added:
				log.Println("po/Added")
				poDir.AddPod(ev.Object)

			case watch.Modified:
				// Update
				poDir.UpdatePod(ev.Object)
			case watch.Deleted:
				// Delete
				poDir.DeletePod(ev.Object)
			}
		}
	}
}
