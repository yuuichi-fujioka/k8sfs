package k8s

import (
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type SecretAddFunc func(obj *corev1.Secret)
type SecretUpdateFunc func(oldObj, newObj *corev1.Secret)
type SecretDeleteFunc func(obj *corev1.Secret)

func WatchSecrets(namespace string, addFunc SecretAddFunc, updateFunc SecretUpdateFunc, deleteFunc SecretDeleteFunc) chan struct{} {
	secretListWatcher := cache.NewListWatchFromClient(Clientset.CoreV1().RESTClient(), "secrets", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(secretListWatcher, &corev1.Secret{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			secret, ok := obj.(*corev1.Secret)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] secret/%s is Added on %s\n", secret.Name, namespace)
			addFunc(secret)
		},
		UpdateFunc: func(old, new interface{}) {
			oldsecret, ok := old.(*corev1.Secret)
			if !ok {
				panic("!!!!")
			}

			newsecret, ok := new.(*corev1.Secret)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] secret/%s is Modified on %s\n", newsecret.Name, namespace)
			updateFunc(oldsecret, newsecret)
		},

		DeleteFunc: func(obj interface{}) {
			secret, ok := obj.(*corev1.Secret)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] secret/%s is Deleted on %s\n", secret.Name, namespace)
			deleteFunc(secret)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
