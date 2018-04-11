package k8s

import (
	"fmt"
	"time"

	"github.com/ghodss/yaml"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateUpdateDeploymentWithYaml(namespace, name string, payload []byte) (*v1beta1.Deployment, error) {
	flagDeleteCancelSignal("deploy", name)

	deploy := &v1beta1.Deployment{}
	err := yaml.Unmarshal(payload, deploy)
	if err != nil {
		return nil, err
	}
	if deploy.Name != name {
		return deploy, fmt.Errorf("file name and metadata.name are not same")
	}

	// remove unnecessary parameters
	deploy.ResourceVersion = ""
	deploy.UID = ""

	// Make or Update NS
	_, err = Clientset.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		return Clientset.ExtensionsV1beta1().Deployments(namespace).Create(deploy)
	} else {
		return Clientset.ExtensionsV1beta1().Deployments(namespace).Update(deploy)
	}
}

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
