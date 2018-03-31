package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

func NewReplicaSetFile(rs *v1beta1.ReplicaSet) *writableFile {
	meta := NewMetaObj(&rs.TypeMeta, &rs.ObjectMeta)
	return NewObjFile(rs, meta, nil)
}

func UpdateReplicaSetFile(f *writableFile, rs *v1beta1.ReplicaSet) {
	meta := NewMetaObj(&rs.TypeMeta, &rs.ObjectMeta)
	f.Update(rs, meta)
}
