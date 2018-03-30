package fuse

import (
	corev1 "k8s.io/api/core/v1"
)

func NewEventFile(ev *corev1.Event) *writableFile {
	meta := NewMetaObj(&ev.TypeMeta, &ev.ObjectMeta)
	return NewObjFile(ev, meta, nil)
}

func UpdateEventFile(f *writableFile, ev *corev1.Event) {
	meta := NewMetaObj(&ev.TypeMeta, &ev.ObjectMeta)
	f.Update(ev, meta)
}
