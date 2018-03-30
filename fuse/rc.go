package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewReplicationControllerFile(rc *corev1.ReplicationController) *writableFile {
	meta := NewMetaObj(&rc.TypeMeta, &rc.ObjectMeta)
	return NewObjFile(rc, meta, nil)
}

func UpdateReplicationControllerFile(f *writableFile, rc *corev1.ReplicationController) {
	meta := NewMetaObj(&rc.TypeMeta, &rc.ObjectMeta)
	f.Update(rc, meta)
}
