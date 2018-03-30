package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewServiceFile(svc *corev1.Service) *writableFile {
	meta := NewMetaObj(&svc.TypeMeta, &svc.ObjectMeta)
	return NewObjFile(svc, meta, nil)
}

func UpdateServiceFile(f *writableFile, svc *corev1.Service) {
	meta := NewMetaObj(&svc.TypeMeta, &svc.ObjectMeta)
	f.Update(svc, meta)
}
