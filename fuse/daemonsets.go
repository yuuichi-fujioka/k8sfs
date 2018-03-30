package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type daemonSetsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewDaemonSetsDir(ns string) (string, *daemonSetsDir) {

	return "ds", &daemonSetsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *daemonSetsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *daemonSetsDir) GetFile() nodefs.File {
	return f
}

func (f *daemonSetsDir) GetDir(name string) DirEntry {
	if name == "" {
		return f
	}

	names := strings.Split(name, "/")

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *daemonSetsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "ds")
	// TODO
	return fuse.ENOSYS
}

func (f *daemonSetsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "ds")
	// TODO
	return fuse.ENOSYS
}

func (f *daemonSetsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "ds")
	// TODO
	return fuse.ENOSYS
}

func (f *daemonSetsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "ds", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *daemonSetsDir) AddDaemonSet(obj runtime.Object) {
	if !f.UpdateDaemonSet(obj) {
		newFile := NewDaemonSetFile(obj)
		f.files[newFile.Name] = newFile
	}
}

func (f *daemonSetsDir) UpdateDaemonSet(obj runtime.Object) (updated bool) {

	ds, ok := obj.(*v1beta1.DaemonSet)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := ds.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateDaemonSetFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *daemonSetsDir) DeleteDaemonSet(obj runtime.Object) {

	ds, ok := obj.(*v1beta1.DaemonSet)
	if !ok {
		panic("!!!!")
	}
	name := ds.Name
	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
