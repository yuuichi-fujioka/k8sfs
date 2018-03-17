package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewEndpointsFile(obj runtime.Object) *writableFile {
	ep, ok := obj.(*corev1.Endpoints)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ep.TypeMeta, &ep.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateEndpointsFile(f *writableFile, obj runtime.Object) {
	ep, ok := obj.(*corev1.Endpoints)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ep.TypeMeta, &ep.ObjectMeta)
	f.Update(obj, meta)
}
