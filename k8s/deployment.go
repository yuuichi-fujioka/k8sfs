package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteDeployment(namespace, name string) error {
	orphanDependents := false
	return Clientset.ExtensionsV1beta1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
		OrphanDependents: &orphanDependents,
	})
}
