package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type EventResource struct {
}

func (me *EventResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().Events(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *EventResource) GetName(obj *runtime.Object) string {
	return me.getEvent(obj).GetName()
}
func (me *EventResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getEvent(obj).GetCreationTimestamp()
}

func (me *EventResource) getEvent(obj *runtime.Object) *corev1.Event {
	event, ok := (*obj).(*corev1.Event)
	if !ok {
		panic("unexpected object type")
	}
	return event
}
