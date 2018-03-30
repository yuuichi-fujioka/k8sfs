package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type PersistentVolumeClaimAddFunc func(obj *corev1.PersistentVolumeClaim)
type PersistentVolumeClaimUpdateFunc func(oldObj, newObj *corev1.PersistentVolumeClaim)
type PersistentVolumeClaimDeleteFunc func(obj *corev1.PersistentVolumeClaim)

func WatchPersistentVolumeClaims(namespace string, addFunc PersistentVolumeClaimAddFunc, updateFunc PersistentVolumeClaimUpdateFunc, deleteFunc PersistentVolumeClaimDeleteFunc) chan struct{} {
	pvcListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "persistentvolumeclaims", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(pvcListWatcher, &corev1.PersistentVolumeClaim{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pvc, ok := obj.(*corev1.PersistentVolumeClaim)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pvc/%s is Added on %s\n", pvc.Name, namespace)
			addFunc(pvc)
		},
		UpdateFunc: func(old, new interface{}) {
			oldpvc, ok := old.(*corev1.PersistentVolumeClaim)
			if !ok {
				panic("!!!!")
			}

			newpvc, ok := new.(*corev1.PersistentVolumeClaim)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pvc/%s is Modified on %s\n", newpvc.Name, namespace)
			updateFunc(oldpvc, newpvc)
		},

		DeleteFunc: func(obj interface{}) {
			pvc, ok := obj.(*corev1.PersistentVolumeClaim)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pvc/%s is Deleted on %s\n", pvc.Name, namespace)
			deleteFunc(pvc)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
