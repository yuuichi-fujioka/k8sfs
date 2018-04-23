package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	v1beta1 "k8s.io/api/extensions/v1beta1"
)

type ingresssDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewIngressesDir(ns string) (string, *ingresssDir) {

	return "ing", &ingresssDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *ingresssDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *ingresssDir) GetFile() nodefs.File {
	return f
}

func (f *ingresssDir) GetDir(name string) DirEntry {
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

func (f *ingresssDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "ing")
	// TODO
	return fuse.ENOSYS
}

func (f *ingresssDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "ing")
	// TODO
	return fuse.ENOSYS
}

func (f *ingresssDir) Rmdir(name string) (code fuse.Status) {
	log.Printf("Rmdir: %s", "ing")
	// TODO
	return fuse.ENOSYS
}

func (f *ingresssDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "ing", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *ingresssDir) AddIngress(ing *v1beta1.Ingress) {
	if !f.UpdateIngress(ing) {
		newFile := NewIngressFile(ing)
		f.files[newFile.Name] = newFile
	}
}

func (f *ingresssDir) UpdateIngress(ing *v1beta1.Ingress) (updated bool) {

	updated = false

	name := ing.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateIngressFile(file, ing)
			updated = true
			break
		}
	}
	return
}

func (f *ingresssDir) DeleteIngress(ing *v1beta1.Ingress) {

	name := ing.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
