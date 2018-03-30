package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type ConfigMapAddFunc func(obj *corev1.ConfigMap)
type ConfigMapUpdateFunc func(oldObj, newObj *corev1.ConfigMap)
type ConfigMapDeleteFunc func(obj *corev1.ConfigMap)

func WatchConfigMaps(namespace string, addFunc ConfigMapAddFunc, updateFunc ConfigMapUpdateFunc, deleteFunc ConfigMapDeleteFunc) chan struct{} {
	cmListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "configmaps", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(cmListWatcher, &corev1.ConfigMap{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm, ok := obj.(*corev1.ConfigMap)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] cm/%s is Added on %s\n", cm.Name, namespace)
			addFunc(cm)
		},
		UpdateFunc: func(old, new interface{}) {
			oldcm, ok := old.(*corev1.ConfigMap)
			if !ok {
				panic("!!!!")
			}

			newcm, ok := new.(*corev1.ConfigMap)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] cm/%s is Modified on %s\n", newcm.Name, namespace)
			updateFunc(oldcm, newcm)
		},

		DeleteFunc: func(obj interface{}) {
			cm, ok := obj.(*corev1.ConfigMap)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] cm/%s is Deleted on %s\n", cm.Name, namespace)
			deleteFunc(cm)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
