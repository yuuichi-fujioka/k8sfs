package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type configMapsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewConfigMapsDir(ns string) (string, *configMapsDir) {

	return "cm", &configMapsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *configMapsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *configMapsDir) GetFile() nodefs.File {
	return f
}

func (f *configMapsDir) GetDir(name string) DirEntry {
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

func (f *configMapsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "cm")
	// TODO
	return fuse.ENOSYS
}

func (f *configMapsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "cm")
	// TODO
	return fuse.ENOSYS
}

func (f *configMapsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "cm")
	// TODO
	return fuse.ENOSYS
}

func (f *configMapsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "cm", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *configMapsDir) AddConfigMap(cm *corev1.ConfigMap) {
	if !f.UpdateConfigMap(cm) {
		newFile := NewConfigMapFile(cm)
		f.files[newFile.Name] = newFile
	}
}

func (f *configMapsDir) UpdateConfigMap(cm *corev1.ConfigMap) (updated bool) {

	updated = false

	name := cm.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateConfigMapFile(file, cm)
			updated = true
			break
		}
	}
	return
}

func (f *configMapsDir) DeleteConfigMap(cm *corev1.ConfigMap) {

	name := cm.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
