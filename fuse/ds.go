package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewDaemonSetFile(obj runtime.Object) *writableFile {
	ds, ok := obj.(*v1beta1.DaemonSet)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ds.TypeMeta, &ds.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateDaemonSetFile(f *writableFile, obj runtime.Object) {
	ds, ok := obj.(*v1beta1.DaemonSet)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ds.TypeMeta, &ds.ObjectMeta)
	f.Update(obj, meta)
}
