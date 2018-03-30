package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewPodFile(pod *corev1.Pod) *writableFile {
	meta := NewMetaObj(&pod.TypeMeta, &pod.ObjectMeta)
	return NewObjFile(pod, meta, nil)
}

func UpdatePodFile(f *writableFile, pod *corev1.Pod) {
	meta := NewMetaObj(&pod.TypeMeta, &pod.ObjectMeta)
	f.Update(pod, meta)
}
