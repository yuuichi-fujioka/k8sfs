package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewReplicationControllerFile(obj runtime.Object) *writableFile {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&rc.TypeMeta, &rc.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateReplicationControllerFile(f *writableFile, obj runtime.Object) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&rc.TypeMeta, &rc.ObjectMeta)
	f.Update(obj, meta)
}
