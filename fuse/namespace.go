package fuse

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespaceDir struct {
	nodefs.File
	defaultDir

	metaObj
}

func NewNamespaceDir(obj *runtime.Object) *namespaceDir {
	ns, ok := (*obj).(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)

	return &namespaceDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		metaObj:    *meta,
	}
}

func (f *namespaceDir) GetName() string {
	return f.Name
}

func (f *namespaceDir) GetAttr(out *fuse.Attr) fuse.Status {
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	SetAttrTime(&f.metaObj, out)
	return fuse.OK
}

func (f *namespaceDir) GetFile() nodefs.File {
	return f
}

func (f *namespaceDir) GetDir(name string) DirEntry {
	if name == "" {
		return f
	}

	return nil
}

func NewNamespaceFile(obj *runtime.Object) *objFile {
	ns, ok := (*obj).(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	return NewObjFile(obj, meta)
}
