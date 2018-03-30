package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type ServiceAddFunc func(obj *corev1.Service)
type ServiceUpdateFunc func(oldObj, newObj *corev1.Service)
type ServiceDeleteFunc func(obj *corev1.Service)

func WatchServices(namespace string, addFunc ServiceAddFunc, updateFunc ServiceUpdateFunc, deleteFunc ServiceDeleteFunc) chan struct{} {
	svcListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "services", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(svcListWatcher, &corev1.Service{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] svc/%s is Added on %s\n", svc.Name, namespace)
			addFunc(svc)
		},
		UpdateFunc: func(old, new interface{}) {
			oldsvc, ok := old.(*corev1.Service)
			if !ok {
				panic("!!!!")
			}

			newsvc, ok := new.(*corev1.Service)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] svc/%s is Modified on %s\n", newsvc.Name, namespace)
			updateFunc(oldsvc, newsvc)
		},

		DeleteFunc: func(obj interface{}) {
			svc, ok := obj.(*corev1.Service)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] svc/%s is Deleted on %s\n", svc.Name, namespace)
			deleteFunc(svc)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
