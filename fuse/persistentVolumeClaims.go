package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type persistentVolumeClaimsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewPersistentVolumeClaimsDir(ns string) (string, *persistentVolumeClaimsDir) {

	return "pvc", &persistentVolumeClaimsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *persistentVolumeClaimsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
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

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *persistentVolumeClaimsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "pvc")
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "pvc")
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "pvc")
	// TODO
	return fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "pvc", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *persistentVolumeClaimsDir) AddPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) {
	if !f.UpdatePersistentVolumeClaim(pvc) {
		newFile := NewPersistentVolumeClaimFile(pvc)
		f.files[newFile.Name] = newFile
	}
}

func (f *persistentVolumeClaimsDir) UpdatePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) (updated bool) {

	updated = false

	name := pvc.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdatePersistentVolumeClaimFile(file, pvc)
			updated = true
			break
		}
	}
	return
}

func (f *persistentVolumeClaimsDir) DeletePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) {

	name := pvc.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
