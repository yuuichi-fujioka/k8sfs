package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

func NewDaemonSetFile(ds *v1beta1.DaemonSet) *writableFile {
	meta := NewMetaObj(&ds.TypeMeta, &ds.ObjectMeta)
	return NewObjFile(ds, meta, nil)
}

func UpdateDaemonSetFile(f *writableFile, ds *v1beta1.DaemonSet) {
	meta := NewMetaObj(&ds.TypeMeta, &ds.ObjectMeta)
	f.Update(ds, meta)
}
