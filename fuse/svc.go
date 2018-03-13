package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewServiceFile(obj runtime.Object) *writableFile {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&svc.TypeMeta, &svc.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdateServiceFile(f *writableFile, obj runtime.Object) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&svc.TypeMeta, &svc.ObjectMeta)
	f.Update(obj, meta)
}
