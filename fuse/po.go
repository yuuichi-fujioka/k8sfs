package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewPodFile(obj runtime.Object) *writableFile {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&pod.TypeMeta, &pod.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdatePodFile(f *writableFile, obj runtime.Object) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&pod.TypeMeta, &pod.ObjectMeta)
	f.Update(obj, meta)
}
