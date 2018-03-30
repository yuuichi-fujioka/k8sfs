package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewPersistentVolumeClaimFile(pvc *corev1.PersistentVolumeClaim) *writableFile {
	meta := NewMetaObj(&pvc.TypeMeta, &pvc.ObjectMeta)
	return NewObjFile(pvc, meta, nil)
}

func UpdatePersistentVolumeClaimFile(f *writableFile, pvc *corev1.PersistentVolumeClaim) {
	meta := NewMetaObj(&pvc.TypeMeta, &pvc.ObjectMeta)
	f.Update(pvc, meta)
}
