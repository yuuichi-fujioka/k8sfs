package k8s

import (
	corev1 "k8s.io/api/core/v1"
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

func DeleteNamespace(name string) error {
	return Clientset.CoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
}
