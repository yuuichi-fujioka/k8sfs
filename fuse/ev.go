package fuse

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func NewEventFile(obj runtime.Object) *writableFile {
	ev, ok := obj.(*corev1.Event)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ev.TypeMeta, &ev.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdateEventFile(f *writableFile, obj runtime.Object) {
	ev, ok := obj.(*corev1.Event)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ev.TypeMeta, &ev.ObjectMeta)
	f.Update(obj, meta)
}
