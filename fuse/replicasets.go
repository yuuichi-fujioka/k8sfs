package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type replicaSetsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewReplicaSetsDir(ns string) (string, *replicaSetsDir) {

	return "rs", &replicaSetsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *replicaSetsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *replicaSetsDir) GetFile() nodefs.File {
	return f
}

func (f *replicaSetsDir) GetDir(name string) DirEntry {
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

func (f *replicaSetsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "rs")
	// TODO
	return fuse.ENOSYS
}

func (f *replicaSetsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "rs")
	// TODO
	return fuse.ENOSYS
}

func (f *replicaSetsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "rs")
	// TODO
	return fuse.ENOSYS
}

func (f *replicaSetsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "rs", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *replicaSetsDir) AddReplicaSet(obj runtime.Object) {
	if !f.UpdateReplicaSet(obj) {
		newFile := NewReplicaSetFile(obj)
		f.files[newFile.Name] = newFile
	}
}

func (f *replicaSetsDir) UpdateReplicaSet(obj runtime.Object) (updated bool) {

	rs, ok := obj.(*v1beta1.ReplicaSet)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := rs.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateReplicaSetFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *replicaSetsDir) DeleteReplicaSet(obj runtime.Object) {

	rs, ok := obj.(*v1beta1.ReplicaSet)
	if !ok {
		panic("!!!!")
	}
	name := rs.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
