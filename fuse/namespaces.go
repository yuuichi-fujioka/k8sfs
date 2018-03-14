package fuse

import (
	"log"
	"strings"
	"time"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespacesDir struct {
	nodefs.File
	defaultDir
}

func NewNamespacesDir() *namespacesDir {
	return &namespacesDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
	}
}

func (f *namespacesDir) GetName() string {
	return "namespaces"
}

func (f *namespacesDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *namespacesDir) GetFile() nodefs.File {
	return f
}

func (f *namespacesDir) GetDir(name string) DirEntry {
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

func (f *namespacesDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.GetName())
	code = f.RemoveTmpFile(name)
	if code != fuse.ENOENT {
		return
	}

	nsname := strings.TrimSuffix(name, ".yaml")
	err := k8s.DeleteNamespace(nsname)
	if err != nil {
		return fuse.EIO
	}
	return fuse.OK
}

func (f *namespacesDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.GetName())

	_, err := k8s.CreateNamespace(name)
	if err != nil {
		return fuse.EIO
	}

	for {
		if f.GetDir(name) != nil {
			return fuse.OK
		}
		// TODO Event handling should be smart.
		time.Sleep(10 * time.Millisecond)
	}
}

func (f *namespacesDir) Rmdir() (code fuse.Status) {
	log.Printf("Rmdir: %s", f.GetName())
	// TODO
	return fuse.ENOSYS
}

func (f *namespacesDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.GetName(), flags, mode)
	// TODO
	return f.AddTmpFile(name), fuse.OK
}

func (f *namespacesDir) AddNamespace(obj runtime.Object) {
	if !f.UpdateNamespace(obj) {
		f.dirs = append(f.dirs, NewNamespaceDir(obj))
		f.files = append(f.files, NewNamespaceFile(obj))
	}
}

func (f *namespacesDir) UpdateNamespace(obj runtime.Object) (updated bool) {

	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	updated = false

	name := ns.Name
	for _, dir := range f.dirs {
		if dir.GetName() == name {
			nsDir, ok := (dir).(*namespaceDir)
			if !ok {
				panic("!!!")
			}
			nsDir.Update(obj)
			updated = true
			break
		}
	}
	for _, file := range f.files {
		if file.Name == name+".yaml" {
			UpdateNamespaceFile(file, obj)
			updated = true
			break
		}
	}
	return
}

func (f *namespacesDir) DeleteNamespace(obj runtime.Object) {

	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}
	name := ns.Name

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
		if file.Name == name+".yaml" {
			newlist2 = append(f.files[:i], f.files[i+1:]...)
			break
		}
	}
	f.files = newlist2
}