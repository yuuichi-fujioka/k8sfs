package main

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type DeploymentResource struct {
}

func (me *DeploymentResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.ExtensionsV1beta1().Deployments(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *DeploymentResource) GetName(obj *runtime.Object) string {
	return me.getDeployment(obj).GetName()
}
func (me *DeploymentResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getDeployment(obj).GetCreationTimestamp()
}

func (me *DeploymentResource) getDeployment(obj *runtime.Object) *v1beta1.Deployment {
	service, ok := (*obj).(*v1beta1.Deployment)
	if !ok {
		panic("unexpected object type")
	}
	return service
}
