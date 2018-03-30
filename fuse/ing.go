package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

func NewIngressFile(ing *v1beta1.Ingress) *writableFile {
	meta := NewMetaObj(&ing.TypeMeta, &ing.ObjectMeta)
	return NewObjFile(ing, meta, nil)
}

func UpdateIngressFile(f *writableFile, ing *v1beta1.Ingress) {
	meta := NewMetaObj(&ing.TypeMeta, &ing.ObjectMeta)
	f.Update(ing, meta)
}
