package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type ReplicationControllerAddFunc func(obj *corev1.ReplicationController)
type ReplicationControllerUpdateFunc func(oldObj, newObj *corev1.ReplicationController)
type ReplicationControllerDeleteFunc func(obj *corev1.ReplicationController)

func WatchReplicationControllers(namespace string, addFunc ReplicationControllerAddFunc, updateFunc ReplicationControllerUpdateFunc, deleteFunc ReplicationControllerDeleteFunc) chan struct{} {
	rcListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "replicationcontrollers", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(rcListWatcher, &corev1.ReplicationController{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			rc, ok := obj.(*corev1.ReplicationController)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rc/%s is Added on %s\n", rc.Name, namespace)
			addFunc(rc)
		},
		UpdateFunc: func(old, new interface{}) {
			oldrc, ok := old.(*corev1.ReplicationController)
			if !ok {
				panic("!!!!")
			}

			newrc, ok := new.(*corev1.ReplicationController)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rc/%s is Modified on %s\n", newrc.Name, namespace)
			updateFunc(oldrc, newrc)
		},

		DeleteFunc: func(obj interface{}) {
			rc, ok := obj.(*corev1.ReplicationController)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rc/%s is Deleted on %s\n", rc.Name, namespace)
			deleteFunc(rc)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
