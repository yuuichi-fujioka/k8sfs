package k8s

import (
	"fmt"

	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
