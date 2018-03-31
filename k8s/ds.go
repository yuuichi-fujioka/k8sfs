package k8s

import (
	"log"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type DaemonSetAddFunc func(obj *v1beta1.DaemonSet)
type DaemonSetUpdateFunc func(oldObj, newObj *v1beta1.DaemonSet)
type DaemonSetDeleteFunc func(obj *v1beta1.DaemonSet)

func WatchDaemonSets(namespace string, addFunc DaemonSetAddFunc, updateFunc DaemonSetUpdateFunc, deleteFunc DaemonSetDeleteFunc) chan struct{} {
	dsListWatcher := cache.NewListWatchFromClient(Clientset.ExtensionsV1beta1().RESTClient(), "daemonsets", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(dsListWatcher, &v1beta1.DaemonSet{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ds, ok := obj.(*v1beta1.DaemonSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ds/%s is Added on %s\n", ds.Name, namespace)
			addFunc(ds)
		},
		UpdateFunc: func(old, new interface{}) {
			oldds, ok := old.(*v1beta1.DaemonSet)
			if !ok {
				panic("!!!!")
			}

			newds, ok := new.(*v1beta1.DaemonSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ds/%s is Modified on %s\n", newds.Name, namespace)
			updateFunc(oldds, newds)
		},

		DeleteFunc: func(obj interface{}) {
			ds, ok := obj.(*v1beta1.DaemonSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ds/%s is Deleted on %s\n", ds.Name, namespace)
			deleteFunc(ds)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
