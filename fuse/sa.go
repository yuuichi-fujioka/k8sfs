package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewServiceAccountFile(sa *corev1.ServiceAccount) *writableFile {
	meta := NewMetaObj(&sa.TypeMeta, &sa.ObjectMeta)
	return NewObjFile(sa, meta, nil)
}

func UpdateServiceAccountFile(f *writableFile, sa *corev1.ServiceAccount) {
	meta := NewMetaObj(&sa.TypeMeta, &sa.ObjectMeta)
	f.Update(sa, meta)
}
