package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewServiceAccountFile(obj runtime.Object) *writableFile {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&sa.TypeMeta, &sa.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateServiceAccountFile(f *writableFile, obj runtime.Object) {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&sa.TypeMeta, &sa.ObjectMeta)
	f.Update(obj, meta)
}
