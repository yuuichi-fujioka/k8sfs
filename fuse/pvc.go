package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewPersistentVolumeClaimFile(obj runtime.Object) *writableFile {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&pvc.TypeMeta, &pvc.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdatePersistentVolumeClaimFile(f *writableFile, obj runtime.Object) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&pvc.TypeMeta, &pvc.ObjectMeta)
	f.Update(obj, meta)
}
