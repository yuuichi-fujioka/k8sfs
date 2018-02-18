package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type ServiceResource struct {
}

func (me *ServiceResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().Services(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *ServiceResource) GetName(obj *runtime.Object) string {
	return me.getService(obj).GetName()
}
func (me *ServiceResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getService(obj).GetCreationTimestamp()
}

func (me *ServiceResource) getService(obj *runtime.Object) *corev1.Service {
	service, ok := (*obj).(*corev1.Service)
	if !ok {
		panic("unexpected object type")
	}
	return service
}
