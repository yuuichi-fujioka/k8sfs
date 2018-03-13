package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type serviceAccountsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewServiceAccountsDir(ns string) *serviceAccountsDir {

	return &serviceAccountsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *serviceAccountsDir) GetName() string {
	return "sa"
}

func (f *serviceAccountsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *serviceAccountsDir) GetFile() nodefs.File {
	return f
}

func (f *serviceAccountsDir) GetDir(name string) DirEntry {
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

func (f *serviceAccountsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *serviceAccountsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *serviceAccountsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *serviceAccountsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *serviceAccountsDir) AddServiceAccount(obj runtime.Object) {
	if !f.UpdateServiceAccount(obj) {
		f.files = append(f.files, NewServiceAccountFile(obj))
	}
}

func (f *serviceAccountsDir) UpdateServiceAccount(obj runtime.Object) (updated bool) {

	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := sa.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateServiceAccountFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *serviceAccountsDir) DeleteServiceAccount(obj runtime.Object) {

	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		panic("!!!!")
	}
	name := sa.Name

	newlist := f.dirs
	for i, dir := range f.dirs {
		if dir.GetName() == name {
			newlist = append(f.dirs[:i], f.dirs[i+1:]...)
			break
		}
	}
	f.dirs = newlist

	newlist2 := f.files
	for i, file := range f.files {
		if file.Name == name+".yaml" {
			newlist2 = append(f.files[:i], f.files[i+1:]...)
			break
		}
	}
	f.files = newlist2
}
