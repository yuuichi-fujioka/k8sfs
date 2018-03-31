package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewSecretFile(secret *corev1.Secret) *writableFile {
	meta := NewMetaObj(&secret.TypeMeta, &secret.ObjectMeta)
	return NewObjFile(secret, meta, nil)
}

func UpdateSecretFile(f *writableFile, secret *corev1.Secret) {
	meta := NewMetaObj(&secret.TypeMeta, &secret.ObjectMeta)
	f.Update(secret, meta)
}
