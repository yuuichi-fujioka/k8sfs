package fuse

import (
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func TestMain(mountpoint, namespace string, readonly bool) {
	topLevelNamespace = namespace
	if topLevelNamespace == "" {
		go watchAllNs()
	} else {
		go watchNs()
	}

	readOnlyMode = readonly

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
		log.Printf("[Watch] start watchNs\n")
		wi, err := k8s.Clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				log.Printf("[Watch] finish watchNs\n")
				break
			}

			ns, ok := ev.Object.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			nsname := ns.Name

			switch ev.Type {
			case watch.Added:
				log.Printf("[Watch] ns/%s is Added\n", nsname)
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
				log.Printf("[Watch] ns/%s is Modified\n", nsname)
				nsDir.UpdateNamespace(ev.Object)
			case watch.Deleted:
				// Delete
				log.Printf("[Watch] ns/%s is Deleted\n", nsname)
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
