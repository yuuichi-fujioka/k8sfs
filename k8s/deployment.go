package k8s

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeleteDeployment(namespace, name string) error {
	orphanDependents := false
	return Clientset.ExtensionsV1beta1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
		OrphanDependents: &orphanDependents,
	})
}

func DeleteDeploymentLazy(namespace, name string) {
	initDeleteCancelSignal("deploy", namespace, name)
	go func() {
		time.Sleep(1 * time.Second)
		if isFlagedDeleteCancelSignal("deploy", namespace, name) {
			return
		}
		err := DeleteDeployment(namespace, name)
		if err != nil {
			panic(err)
		}
	}()
}
