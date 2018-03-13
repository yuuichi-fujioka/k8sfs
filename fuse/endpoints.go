package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type endpointsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewEndpointsDir(ns string) *endpointsDir {

	return &endpointsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *endpointsDir) GetName() string {
	return "ep"
}

func (f *endpointsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *endpointsDir) GetFile() nodefs.File {
	return f
}

func (f *endpointsDir) GetDir(name string) DirEntry {
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

func (f *endpointsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *endpointsDir) AddEndpoints(obj runtime.Object) {
	if !f.UpdateEndpoints(obj) {
		f.files = append(f.files, NewEndpointsFile(obj))
	}
}

func (f *endpointsDir) UpdateEndpoints(obj runtime.Object) (updated bool) {

	ep, ok := obj.(*corev1.Endpoints)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := ep.Name
	for _, file := range f.files {
		if file.Name == name {
			UpdateEndpointsFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *endpointsDir) DeleteEndpoints(obj runtime.Object) {

	ep, ok := obj.(*corev1.Endpoints)
	if !ok {
		panic("!!!!")
	}
	name := ep.Name

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
		if file.Name == name {
			newlist2 = append(f.files[:i], f.files[i+1:]...)
			break
		}
	}
	f.files = newlist2
}
