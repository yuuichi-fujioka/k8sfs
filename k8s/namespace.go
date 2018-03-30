package k8s

import (
	"fmt"
	"log"

	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func CreateNamespace(name string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	return Clientset.CoreV1().Namespaces().Create(ns)
}

func CreateUpdateNamespaceWithYaml(name string, payload []byte) (*corev1.Namespace, error) {

	ns := &corev1.Namespace{}
	err := yaml.Unmarshal(payload, ns)
	if err != nil {
		return nil, err
	}
	if ns.Name != name {
		return ns, fmt.Errorf("file name and metadata.name are not same")
	}

	// remove unnecessary parameters
	ns.ResourceVersion = ""
	ns.UID = ""

	// Make or Update NS
	_, err = Clientset.CoreV1().Namespaces().Get(name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		return Clientset.CoreV1().Namespaces().Create(ns)
	} else {
		return Clientset.CoreV1().Namespaces().Update(ns)
	}
}

func DeleteNamespace(name string) error {
	return Clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
}

type NamespaceAddFunc func(obj *corev1.Namespace)
type NamespaceUpdateFunc func(oldObj, newObj *corev1.Namespace)
type NamespaceDeleteFunc func(obj *corev1.Namespace)

func WatchNamespaces(addFunc NamespaceAddFunc, updateFunc NamespaceUpdateFunc, deleteFunc NamespaceDeleteFunc) chan struct{} {
	listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
		return Clientset.CoreV1().Namespaces().List(options)
	}
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		options.Watch = true
		return Clientset.CoreV1().Namespaces().Watch(options)
	}

	w := &cache.ListWatch{
		ListFunc:  listFunc,
		WatchFunc: watchFunc,
	}

	_, informer := cache.NewIndexerInformer(w, &corev1.Namespace{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ns/%s is Added\n", ns.Name)
			addFunc(ns)
		},
		UpdateFunc: func(old, new interface{}) {
			oldns, ok := old.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			newns, ok := new.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ns/%s is Modified\n", newns.Name)
			updateFunc(oldns, newns)
		},

		DeleteFunc: func(obj interface{}) {
			ns, ok := obj.(*corev1.Namespace)
			if !ok {
				panic("!!!!")
			}

			log.Printf("[Watch] ns/%s is Deleted\n", ns.Name)
			deleteFunc(ns)
		},
	}, cache.Indexers{})

	closeCh := make(chan struct{})
	go informer.Run(closeCh)
	return closeCh
}
