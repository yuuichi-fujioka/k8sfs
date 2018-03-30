package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewEndpointsFile(ep *corev1.Endpoints) *writableFile {
	meta := NewMetaObj(&ep.TypeMeta, &ep.ObjectMeta)
	return NewObjFile(ep, meta, nil)
}

func UpdateEndpointsFile(f *writableFile, ep *corev1.Endpoints) {
	meta := NewMetaObj(&ep.TypeMeta, &ep.ObjectMeta)
	f.Update(ep, meta)
}
