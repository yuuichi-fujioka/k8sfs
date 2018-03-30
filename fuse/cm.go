package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewConfigMapFile(cm *corev1.ConfigMap) *writableFile {
	meta := NewMetaObj(&cm.TypeMeta, &cm.ObjectMeta)
	return NewObjFile(cm, meta, nil)
}

func UpdateConfigMapFile(f *writableFile, cm *corev1.ConfigMap) {
	meta := NewMetaObj(&cm.TypeMeta, &cm.ObjectMeta)
	f.Update(cm, meta)
}
