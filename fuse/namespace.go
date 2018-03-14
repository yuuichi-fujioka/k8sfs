package fuse

import (
	"log"
	"strings"

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

func NewNamespaceDir(obj runtime.Object) *namespaceDir {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)

	d := &namespaceDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		metaObj:    *meta,
	}
	d.dirs = append(d.dirs, NewConfigMapsDir(ns.Name))
	d.dirs = append(d.dirs, NewDeploymentsDir(ns.Name))
	d.dirs = append(d.dirs, NewEndpointsDir(ns.Name))
	d.dirs = append(d.dirs, NewEventsDir(ns.Name))
	d.dirs = append(d.dirs, NewIngressesDir(ns.Name))
	d.dirs = append(d.dirs, NewPersistentVolumeClaimsDir(ns.Name))
	d.dirs = append(d.dirs, NewPodsDir(ns.Name))
	d.dirs = append(d.dirs, NewReplicationControllersDir(ns.Name))
	d.dirs = append(d.dirs, NewSecretsDir(ns.Name))
	d.dirs = append(d.dirs, NewServiceAccountsDir(ns.Name))
	d.dirs = append(d.dirs, NewServicesDir(ns.Name))
	d.dirs = append(d.dirs, NewDaemonSetsDir(ns.Name))
	return d
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

	names := strings.Split(name, "/")

	for _, child := range f.dirs {
		if child.GetName() == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *namespaceDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *namespaceDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *namespaceDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *namespaceDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *namespaceDir) Update(obj runtime.Object) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	f.metaObj = *meta
}

func NewNamespaceFile(obj runtime.Object) *writableFile {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	return NewObjFile(obj, meta)
}

func UpdateNamespaceFile(f *writableFile, obj runtime.Object) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	f.Update(obj, meta)
}
