package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewConfigMapFile(obj runtime.Object) *writableFile {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&cm.TypeMeta, &cm.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdateConfigMapFile(f *writableFile, obj runtime.Object) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&cm.TypeMeta, &cm.ObjectMeta)
	f.Update(obj, meta)
}
