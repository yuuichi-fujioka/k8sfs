package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type ServiceAccountAddFunc func(obj *corev1.ServiceAccount)
type ServiceAccountUpdateFunc func(oldObj, newObj *corev1.ServiceAccount)
type ServiceAccountDeleteFunc func(obj *corev1.ServiceAccount)

func WatchServiceAccounts(namespace string, addFunc ServiceAccountAddFunc, updateFunc ServiceAccountUpdateFunc, deleteFunc ServiceAccountDeleteFunc) chan struct{} {
	saListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "serviceaccounts", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(saListWatcher, &corev1.ServiceAccount{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			sa, ok := obj.(*corev1.ServiceAccount)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] sa/%s is Added on %s\n", sa.Name, namespace)
			addFunc(sa)
		},
		UpdateFunc: func(old, new interface{}) {
			oldsa, ok := old.(*corev1.ServiceAccount)
			if !ok {
				panic("!!!!")
			}

			newsa, ok := new.(*corev1.ServiceAccount)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] sa/%s is Modified on %s\n", newsa.Name, namespace)
			updateFunc(oldsa, newsa)
		},

		DeleteFunc: func(obj interface{}) {
			sa, ok := obj.(*corev1.ServiceAccount)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] sa/%s is Deleted on %s\n", sa.Name, namespace)
			deleteFunc(sa)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
