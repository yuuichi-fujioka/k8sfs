package fuse

import (
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func TestMain(mountpoint, namespace string) {
	topLevelNamespace = namespace
	if topLevelNamespace == "" {
		go watchAllNs()
	} else {
		go watchNs()
	}

	Serve(mountpoint)
}

func watchNs() {
	nsw := NewNsWatcher(topLevelNamespace)
	nsw.StartAll()
}

func watchAllNs() {
	watchers := map[string]*nsWatcher{}
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

			nsname := ns.Name

			switch ev.Type {
			case watch.Added:
				log.Println("ns/Added")
				nsDir.AddNamespace(ev.Object)

				nsw, ok := watchers[nsname]
				if !ok {
					nsw = NewNsWatcher(nsname)
					watchers[nsname] = nsw
				} else {
					nsw.StopAll()
				}
				nsw.StartAll()
			case watch.Modified:
				// Update
				nsDir.UpdateNamespace(ev.Object)
			case watch.Deleted:
				// Delete
				nsDir.DeleteNamespace(ev.Object)
				// TODO Stop, watcher
				nsw, ok := watchers[nsname]
				if ok {
					nsw.StopAll()
				}
			}
		}
	}
}
