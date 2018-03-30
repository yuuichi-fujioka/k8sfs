package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
)

func NewDeploymentFile(deploy *v1beta1.Deployment) *writableFile {
	meta := NewMetaObj(&deploy.TypeMeta, &deploy.ObjectMeta)
	return NewObjFile(deploy, meta, nil)
}

func UpdateDeploymentFile(f *writableFile, deploy *v1beta1.Deployment) {
	meta := NewMetaObj(&deploy.TypeMeta, &deploy.ObjectMeta)
	f.Update(deploy, meta)
}
