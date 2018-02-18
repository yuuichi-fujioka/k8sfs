package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type SecretResource struct {
}

func (me *SecretResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().Secrets(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *SecretResource) GetName(obj *runtime.Object) string {
	return me.getSecret(obj).GetName()
}
func (me *SecretResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getSecret(obj).GetCreationTimestamp()
}

func (me *SecretResource) getSecret(obj *runtime.Object) *corev1.Secret {
	secret, ok := (*obj).(*corev1.Secret)
	if !ok {
		panic("unexpected object type")
	}
	return secret
}
