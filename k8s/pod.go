package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type PodAddFunc func(obj *corev1.Pod)
type PodUpdateFunc func(oldObj, newObj *corev1.Pod)
type PodDeleteFunc func(obj *corev1.Pod)

func WatchPods(namespace string, addFunc PodAddFunc, updateFunc PodUpdateFunc, deleteFunc PodDeleteFunc) chan struct{} {
	podListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "pods", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(podListWatcher, &corev1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pod/%s is Added on %s\n", pod.Name, namespace)
			addFunc(pod)
		},
		UpdateFunc: func(old, new interface{}) {
			oldpod, ok := old.(*corev1.Pod)
			if !ok {
				panic("!!!!")
			}

			newpod, ok := new.(*corev1.Pod)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pod/%s is Modified on %s\n", newpod.Name, namespace)
			updateFunc(oldpod, newpod)
		},

		DeleteFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] pod/%s is Deleted on %s\n", pod.Name, namespace)
			deleteFunc(pod)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
