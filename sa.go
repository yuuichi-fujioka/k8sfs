package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type ServiceAccountResource struct {
}

func (me *ServiceAccountResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().ServiceAccounts(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *ServiceAccountResource) GetName(obj *runtime.Object) string {
	return me.getServiceAccount(obj).GetName()
}
func (me *ServiceAccountResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getServiceAccount(obj).GetCreationTimestamp()
}

func (me *ServiceAccountResource) getServiceAccount(obj *runtime.Object) *corev1.ServiceAccount {
	serviceAccount, ok := (*obj).(*corev1.ServiceAccount)
	if !ok {
		panic("unexpected object type")
	}
	return serviceAccount
}
