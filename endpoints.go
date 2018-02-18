package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type EndpointResource struct {
}

func (me *EndpointResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().Endpoints(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *EndpointResource) GetName(obj *runtime.Object) string {
	return me.getEndpoint(obj).GetName()
}
func (me *EndpointResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getEndpoint(obj).GetCreationTimestamp()
}

func (me *EndpointResource) getEndpoint(obj *runtime.Object) *corev1.Endpoints {
	endpoint, ok := (*obj).(*corev1.Endpoints)
	if !ok {
		panic("unexpected object type")
	}
	return endpoint
}
