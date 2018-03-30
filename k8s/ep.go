package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type EndpointsAddFunc func(obj *corev1.Endpoints)
type EndpointsUpdateFunc func(oldObj, newObj *corev1.Endpoints)
type EndpointsDeleteFunc func(obj *corev1.Endpoints)

func WatchEndpoints(namespace string, addFunc EndpointsAddFunc, updateFunc EndpointsUpdateFunc, deleteFunc EndpointsDeleteFunc) chan struct{} {
	epListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "endpoints", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(epListWatcher, &corev1.Endpoints{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ep/%s is Added on %s\n", ep.Name, namespace)
			addFunc(ep)
		},
		UpdateFunc: func(old, new interface{}) {
			oldep, ok := old.(*corev1.Endpoints)
			if !ok {
				panic("!!!!")
			}

			newep, ok := new.(*corev1.Endpoints)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ep/%s is Modified on %s\n", newep.Name, namespace)
			updateFunc(oldep, newep)
		},

		DeleteFunc: func(obj interface{}) {
			ep, ok := obj.(*corev1.Endpoints)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ep/%s is Deleted on %s\n", ep.Name, namespace)
			deleteFunc(ep)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
