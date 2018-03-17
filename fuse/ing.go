package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewIngressFile(obj runtime.Object) *writableFile {
	ing, ok := obj.(*v1beta1.Ingress)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ing.TypeMeta, &ing.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateIngressFile(f *writableFile, obj runtime.Object) {
	ing, ok := obj.(*v1beta1.Ingress)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ing.TypeMeta, &ing.ObjectMeta)
	f.Update(obj, meta)
}
