package main

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type IngressResource struct {
}

func (me *IngressResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.ExtensionsV1beta1().Ingresses(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *IngressResource) GetName(obj *runtime.Object) string {
	return me.getIngress(obj).GetName()
}
func (me *IngressResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getIngress(obj).GetCreationTimestamp()
}

func (me *IngressResource) getIngress(obj *runtime.Object) *v1beta1.Ingress {
	ingress, ok := (*obj).(*v1beta1.Ingress)
	if !ok {
		panic("unexpected object type")
	}
	return ingress
}
