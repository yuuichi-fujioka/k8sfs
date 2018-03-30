package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type EventAddFunc func(obj *corev1.Event)
type EventUpdateFunc func(oldObj, newObj *corev1.Event)
type EventDeleteFunc func(obj *corev1.Event)

func WatchEvents(namespace string, addFunc EventAddFunc, updateFunc EventUpdateFunc, deleteFunc EventDeleteFunc) chan struct{} {
	evListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "events", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(evListWatcher, &corev1.Event{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ev, ok := obj.(*corev1.Event)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ev/%s is Added on %s\n", ev.Name, namespace)
			addFunc(ev)
		},
		UpdateFunc: func(old, new interface{}) {
			oldev, ok := old.(*corev1.Event)
			if !ok {
				panic("!!!!")
			}

			newev, ok := new.(*corev1.Event)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ev/%s is Modified on %s\n", newev.Name, namespace)
			updateFunc(oldev, newev)
		},

		DeleteFunc: func(obj interface{}) {
			ev, ok := obj.(*corev1.Event)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ev/%s is Deleted on %s\n", ev.Name, namespace)
			deleteFunc(ev)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
