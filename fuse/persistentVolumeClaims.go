package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type persistentVolumeClaimsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewPersistentVolumeClaimsDir(ns string) *persistentVolumeClaimsDir {

	return &persistentVolumeClaimsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *persistentVolumeClaimsDir) GetName() string {
	return "pvc"
}

func (f *persistentVolumeClaimsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *persistentVolumeClaimsDir) GetFile() nodefs.File {
	return f
}

func (f *persistentVolumeClaimsDir) GetDir(name string) DirEntry {
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

func (f *persistentVolumeClaimsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) AddPersistentVolumeClaim(obj runtime.Object) {
	if !f.UpdatePersistentVolumeClaim(obj) {
		f.files = append(f.files, NewPersistentVolumeClaimFile(obj))
	}
}

func (f *persistentVolumeClaimsDir) UpdatePersistentVolumeClaim(obj runtime.Object) (updated bool) {

	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := pvc.Name
	for _, file := range f.files {
		if file.Name == name {
			UpdatePersistentVolumeClaimFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *persistentVolumeClaimsDir) DeletePersistentVolumeClaim(obj runtime.Object) {

	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		panic("!!!!")
	}
	name := pvc.Name

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
