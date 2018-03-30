package k8s

import (
	"log"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type DeploymentAddFunc func(obj *v1beta1.Deployment)
type DeploymentUpdateFunc func(oldObj, newObj *v1beta1.Deployment)
type DeploymentDeleteFunc func(obj *v1beta1.Deployment)

func WatchDeployments(namespace string, addFunc DeploymentAddFunc, updateFunc DeploymentUpdateFunc, deleteFunc DeploymentDeleteFunc) chan struct{} {
	deployListWatcher := cache.NewListWatchFromClient(Clientset.ExtensionsV1beta1().RESTClient(), "deployments", namespace, fields.Everything())

	_, informer := cache.NewIndexerInformer(deployListWatcher, &v1beta1.Deployment{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			deploy, ok := obj.(*v1beta1.Deployment)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] deploy/%s is Added on %s\n", deploy.Name, namespace)
			addFunc(deploy)
		},
		UpdateFunc: func(old, new interface{}) {
			olddeploy, ok := old.(*v1beta1.Deployment)
			if !ok {
				panic("!!!!")
			}

			newdeploy, ok := new.(*v1beta1.Deployment)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] deploy/%s is Modified on %s\n", newdeploy.Name, namespace)
			updateFunc(olddeploy, newdeploy)
		},

		DeleteFunc: func(obj interface{}) {
			deploy, ok := obj.(*v1beta1.Deployment)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] deploy/%s is Deleted on %s\n", deploy.Name, namespace)
			deleteFunc(deploy)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
