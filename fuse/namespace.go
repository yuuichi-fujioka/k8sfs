package fuse

import (
	"github.com/hanwen/go-fuse/fuse"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespacesDir struct {
	defaultDir
}

func NewNamespacesDir() *namespacesDir {
	return &namespacesDir{
		defaultDir: NewDefaultDir(),
	}
}

func (f *namespacesDir) GetAttr(out *fuse.Attr) fuse.Status {
	ctime := uint64(0) // TODO cluster created at
	out.Size = 4096    // block size?
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *namespacesDir) GetName() string {
	return "namespaces"
}

func (f *namespacesDir) AddNamespace(obj *runtime.Object) {
	f.dirs = append(f.dirs, NewNamespaceDir(obj))
	f.files = append(f.files, NewNamespaceFile(obj))
}

type namespaceDir struct {
	objDir
}

func NewNamespaceDir(obj *runtime.Object) *namespaceDir {
	ns, ok := (*obj).(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)

	return &namespaceDir{
		objDir: NewObjDir(meta),
	}
}

func NewNamespaceFile(obj *runtime.Object) *objFile {
	ns, ok := (*obj).(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	return NewObjFile(obj, meta)
}
