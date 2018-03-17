package fuse

import (
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewDeploymentFile(obj runtime.Object) *writableFile {
	dep, ok := obj.(*v1beta1.Deployment)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&dep.TypeMeta, &dep.ObjectMeta)
	return NewObjFile(obj, meta, nil)
}

func UpdateDeploymentFile(f *writableFile, obj runtime.Object) {
	dep, ok := obj.(*v1beta1.Deployment)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&dep.TypeMeta, &dep.ObjectMeta)
	f.Update(obj, meta)
}
