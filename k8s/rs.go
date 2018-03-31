package k8s

import (
	"log"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type ReplicaSetAddFunc func(obj *v1beta1.ReplicaSet)
type ReplicaSetUpdateFunc func(oldObj, newObj *v1beta1.ReplicaSet)
type ReplicaSetDeleteFunc func(obj *v1beta1.ReplicaSet)

func WatchReplicaSets(namespace string, addFunc ReplicaSetAddFunc, updateFunc ReplicaSetUpdateFunc, deleteFunc ReplicaSetDeleteFunc) chan struct{} {
	rsListWatcher := cache.NewListWatchFromClient(Clientset.ExtensionsV1beta1().RESTClient(), "replicasets", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(rsListWatcher, &v1beta1.ReplicaSet{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			rs, ok := obj.(*v1beta1.ReplicaSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rs/%s is Added on %s\n", rs.Name, namespace)
			addFunc(rs)
		},
		UpdateFunc: func(old, new interface{}) {
			oldrs, ok := old.(*v1beta1.ReplicaSet)
			if !ok {
				panic("!!!!")
			}

			newrs, ok := new.(*v1beta1.ReplicaSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rs/%s is Modified on %s\n", newrs.Name, namespace)
			updateFunc(oldrs, newrs)
		},

		DeleteFunc: func(obj interface{}) {
			rs, ok := obj.(*v1beta1.ReplicaSet)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] rs/%s is Deleted on %s\n", rs.Name, namespace)
			deleteFunc(rs)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
