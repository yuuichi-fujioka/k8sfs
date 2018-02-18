package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type ReplicationControllerResource struct {
}

func (me *ReplicationControllerResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().ReplicationControllers(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *ReplicationControllerResource) GetName(obj *runtime.Object) string {
	return me.getReplicationController(obj).GetName()
}
func (me *ReplicationControllerResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getReplicationController(obj).GetCreationTimestamp()
}

func (me *ReplicationControllerResource) getReplicationController(obj *runtime.Object) *corev1.ReplicationController {
	replicationController, ok := (*obj).(*corev1.ReplicationController)
	if !ok {
		panic("unexpected object type")
	}
	return replicationController
}
