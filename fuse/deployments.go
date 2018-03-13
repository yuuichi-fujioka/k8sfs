package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type deploymentsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewDeploymentsDir(ns string) *deploymentsDir {

	return &deploymentsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *deploymentsDir) GetName() string {
	return "deploy"
}

func (f *deploymentsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *deploymentsDir) GetFile() nodefs.File {
	return f
}

func (f *deploymentsDir) GetDir(name string) DirEntry {
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

func (f *deploymentsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *deploymentsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *deploymentsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *deploymentsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *deploymentsDir) AddDeployment(obj runtime.Object) {
	if !f.UpdateDeployment(obj) {
		f.files = append(f.files, NewDeploymentFile(obj))
	}
}

func (f *deploymentsDir) UpdateDeployment(obj runtime.Object) (updated bool) {

	deploy, ok := obj.(*v1beta1.Deployment)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := deploy.Name
	for _, file := range f.files {
		if file.Name == name {
			UpdateDeploymentFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *deploymentsDir) DeleteDeployment(obj runtime.Object) {

	deploy, ok := obj.(*v1beta1.Deployment)
	if !ok {
		panic("!!!!")
	}
	name := deploy.Name

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
