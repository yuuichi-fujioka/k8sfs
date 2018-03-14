package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewReplicaSetFile(obj runtime.Object) *writableFile {
	rs, ok := obj.(*v1beta1.ReplicaSet)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&rs.TypeMeta, &rs.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdateReplicaSetFile(f *writableFile, obj runtime.Object) {
	rs, ok := obj.(*v1beta1.ReplicaSet)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&rs.TypeMeta, &rs.ObjectMeta)
	f.Update(obj, meta)
}
