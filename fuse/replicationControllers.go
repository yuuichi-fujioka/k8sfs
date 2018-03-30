package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type replicationControllersDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewReplicationControllersDir(ns string) (string, *replicationControllersDir) {

	return "rc", &replicationControllersDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *replicationControllersDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *replicationControllersDir) GetFile() nodefs.File {
	return f
}

func (f *replicationControllersDir) GetDir(name string) DirEntry {
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

func (f *replicationControllersDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "rc")
	// TODO
	return fuse.ENOSYS
}

func (f *replicationControllersDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "rc")
	// TODO
	return fuse.ENOSYS
}

func (f *replicationControllersDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "rc")
	// TODO
	return fuse.ENOSYS
}

func (f *replicationControllersDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "rc", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *replicationControllersDir) AddReplicationController(rc *corev1.ReplicationController) {
	if !f.UpdateReplicationController(rc) {
		newFile := NewReplicationControllerFile(rc)
		f.files[newFile.Name] = newFile
	}
}

func (f *replicationControllersDir) UpdateReplicationController(rc *corev1.ReplicationController) (updated bool) {

	updated = false

	name := rc.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateReplicationControllerFile(file, rc)
			updated = true
			break
		}
	}
	return
}

func (f *replicationControllersDir) DeleteReplicationController(rc *corev1.ReplicationController) {

	name := rc.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
