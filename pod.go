package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type PodResource struct {
}

func (me *PodResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().Pods(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *PodResource) GetName(obj *runtime.Object) string {
	return me.getPod(obj).GetName()
}
func (me *PodResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getPod(obj).GetCreationTimestamp()
}

func (me *PodResource) getPod(obj *runtime.Object) *corev1.Pod {
	pod, ok := (*obj).(*corev1.Pod)
	if !ok {
		panic("unexpected object type")
	}
	return pod
}
