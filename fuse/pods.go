package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type podsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewPodsDir(ns string) (string, *podsDir) {

	return "po", &podsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *podsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *podsDir) GetFile() nodefs.File {
	return f
}

func (f *podsDir) GetDir(name string) DirEntry {
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

func (f *podsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "po")
	// TODO
	return fuse.ENOSYS
}

func (f *podsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "po")
	// TODO
	return fuse.ENOSYS
}

func (f *podsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "po")
	// TODO
	return fuse.ENOSYS
}

func (f *podsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "po", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *podsDir) AddPod(pod *corev1.Pod) {
	if !f.UpdatePod(pod) {
		newFile := NewPodFile(pod)
		f.files[newFile.Name] = newFile
	}
}

func (f *podsDir) UpdatePod(pod *corev1.Pod) (updated bool) {

	updated = false

	name := pod.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdatePodFile(file, pod)
			updated = true
			break
		}
	}
	return
}

func (f *podsDir) DeletePod(pod *corev1.Pod) {

	name := pod.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
