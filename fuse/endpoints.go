package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type endpointsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewEndpointsDir(ns string) (string, *endpointsDir) {

	return "ep", &endpointsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *endpointsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
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

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *endpointsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "ep")
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "ep")
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Rmdir(name string) (code fuse.Status) {
	log.Printf("Rmdir: %s", "ep")
	// TODO
	return fuse.ENOSYS
}

func (f *endpointsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "ep", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *endpointsDir) AddEndpoints(ep *corev1.Endpoints) {
	if !f.UpdateEndpoints(ep) {
		newFile := NewEndpointsFile(ep)
		f.files[newFile.Name] = newFile
	}
}

func (f *endpointsDir) UpdateEndpoints(ep *corev1.Endpoints) (updated bool) {

	updated = false

	name := ep.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateEndpointsFile(file, ep)
			updated = true
			break
		}
	}
	return
}

func (f *endpointsDir) DeleteEndpoints(ep *corev1.Endpoints) {

	name := ep.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
