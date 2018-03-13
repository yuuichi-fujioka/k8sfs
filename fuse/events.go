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

func NewEventsDir(ns string) *eventsDir {

	return &eventsDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		Namespace:  ns,
	}
}

func (f *eventsDir) GetName() string {
	return "ev"
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

	for _, child := range f.dirs {
		if child.GetName() == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *eventsDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *eventsDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *eventsDir) AddEvent(obj runtime.Object) {
	if !f.UpdateEvent(obj) {
		f.files = append(f.files, NewEventFile(obj))
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
		if file.Name == name {
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
