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
				go watchPods(ns.Name)
				go watchServices(ns.Name)
				go watchConfigMaps(ns.Name)
				go watchDeployments(ns.Name)
				go watchEndpoints(ns.Name)
				go watchEvents(ns.Name)
				go watchIngresses(ns.Name)
				go watchPersistentVolumeClaims(ns.Name)
				go watchReplicationControllers(ns.Name)
				go watchServiceAccounts(ns.Name)
				go watchSecrets(ns.Name)

			case watch.Modified:
				// Update
				nsDir.UpdateNamespace(ev.Object)
			case watch.Deleted:
				// Delete
				nsDir.DeleteNamespace(ev.Object)
				// TODO Stop, watcher
			}
		}
	}
}
