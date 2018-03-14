package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type secretsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewSecretsDir(ns string) (string, *secretsDir) {

	return "secrets", &secretsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *secretsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *secretsDir) GetFile() nodefs.File {
	return f
}

func (f *secretsDir) GetDir(name string) DirEntry {
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

func (f *secretsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "secrets")
	// TODO
	return fuse.ENOSYS
}

func (f *secretsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "secrets")
	// TODO
	return fuse.ENOSYS
}

func (f *secretsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "secrets")
	// TODO
	return fuse.ENOSYS
}

func (f *secretsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "secrets", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *secretsDir) AddSecret(obj runtime.Object) {
	if !f.UpdateSecret(obj) {
		newFile := NewSecretFile(obj)
		f.files[newFile.Name] = newFile
	}
}

func (f *secretsDir) UpdateSecret(obj runtime.Object) (updated bool) {

	secrets, ok := obj.(*corev1.Secret)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := secrets.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateSecretFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *secretsDir) DeleteSecret(obj runtime.Object) {

	secrets, ok := obj.(*corev1.Secret)
	if !ok {
		panic("!!!!")
	}
	name := secrets.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
