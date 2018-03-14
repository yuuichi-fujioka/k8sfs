package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type eventsDir struct {
	nodefs.File
	defaultDir
	Namespace string
}

func NewEventsDir(ns string) (string, *eventsDir) {

	return "ev", &eventsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *eventsDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *eventsDir) GetFile() nodefs.File {
	return f
}

func (f *eventsDir) GetDir(name string) DirEntry {
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

func (f *eventsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "ev")
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "ev")
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", "ev")
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "ev", flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *eventsDir) AddEvent(obj runtime.Object) {
	if !f.UpdateEvent(obj) {
		newFile := NewEventFile(obj)
		f.files[newFile.Name] = newFile
	}
}

func (f *eventsDir) UpdateEvent(obj runtime.Object) (updated bool) {

	ev, ok := obj.(*corev1.Event)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := ev.Name
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateEventFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *eventsDir) DeleteEvent(obj runtime.Object) {

	ev, ok := obj.(*corev1.Event)
	if !ok {
		panic("!!!!")
	}
	name := ev.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")
}
