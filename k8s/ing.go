package k8s

import (
	"log"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type IngressAddFunc func(obj *v1beta1.Ingress)
type IngressUpdateFunc func(oldObj, newObj *v1beta1.Ingress)
type IngressDeleteFunc func(obj *v1beta1.Ingress)

func WatchIngresses(namespace string, addFunc IngressAddFunc, updateFunc IngressUpdateFunc, deleteFunc IngressDeleteFunc) chan struct{} {
	ingListWatcher := cache.NewListWatchFromClient(Clientset.ExtensionsV1beta1().RESTClient(), "ingresses", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(ingListWatcher, &v1beta1.Ingress{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ing, ok := obj.(*v1beta1.Ingress)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ing/%s is Added on %s\n", ing.Name, namespace)
			addFunc(ing)
		},
		UpdateFunc: func(old, new interface{}) {
			olding, ok := old.(*v1beta1.Ingress)
			if !ok {
				panic("!!!!")
			}

			newing, ok := new.(*v1beta1.Ingress)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ing/%s is Modified on %s\n", newing.Name, namespace)
			updateFunc(olding, newing)
		},

		DeleteFunc: func(obj interface{}) {
			ing, ok := obj.(*v1beta1.Ingress)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ing/%s is Deleted on %s\n", ing.Name, namespace)
			deleteFunc(ing)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
