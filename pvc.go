package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type PersistentVolumeClaimResource struct {
}

func (me *PersistentVolumeClaimResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().PersistentVolumeClaims(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *PersistentVolumeClaimResource) GetName(obj *runtime.Object) string {
	return me.getPersistentVolumeClaim(obj).GetName()
}
func (me *PersistentVolumeClaimResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getPersistentVolumeClaim(obj).GetCreationTimestamp()
}

func (me *PersistentVolumeClaimResource) getPersistentVolumeClaim(obj *runtime.Object) *corev1.PersistentVolumeClaim {
	persistentVolumeClaim, ok := (*obj).(*corev1.PersistentVolumeClaim)
	if !ok {
		panic("unexpected object type")
	}
	return persistentVolumeClaim
}
