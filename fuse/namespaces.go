package fuse

import (
	"log"
	"strings"
	"time"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
)

type namespacesDir struct {
	nodefs.File
	defaultDir
}

func NewNamespacesDir() (string, *namespacesDir) {
	return "namespaces", &namespacesDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
	}
}

func (f *namespacesDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
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

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *namespacesDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "namespaces")
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
	log.Printf("Mkdir: %s at %s", name, "namespaces")

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
	log.Printf("Rmdir: %s", "namespaces")
	// TODO
	return fuse.ENOSYS
}

func (f *namespacesDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "namespaces", flags, mode)
	// TODO
	return f.AddTmpFile(name, f), fuse.OK
}

func (f *namespacesDir) AddNamespace(obj *corev1.Namespace) {
	if !f.UpdateNamespace(obj) {
		name, newDir := NewNamespaceDir(obj)
		f.dirs[name] = newDir
		newFile := NewNamespaceFile(obj, f)
		f.files[newFile.Name] = newFile
	}
}

func (f *namespacesDir) UpdateNamespace(obj *corev1.Namespace) (updated bool) {

	updated = false

	name := obj.Name
	for k, dir := range f.dirs {
		if k == name {
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

func (f *namespacesDir) DeleteNamespace(obj *corev1.Namespace) {
	name := obj.Name

	delete(f.dirs, name)
	delete(f.files, name+".yaml")

}

func (f *namespacesDir) HandleRelease(wf *writableFile) {
	log.Printf("Namespace/HandleRelease: %s", wf.Name)
	if strings.HasPrefix(wf.Name, ".") {
		// This is the Hidden File.
		// TODO Delete this file
		return
	}
	if !strings.HasSuffix(wf.Name, ".yaml") {
		// This is the Hidden File.
		// TODO Delete this file
		return
	}
	if !wf.changed {
		// Namespace is Not Changed
		return
	}

	nsname := strings.TrimSuffix(wf.Name, ".yaml")
	ns, err := k8s.CreateUpdateNamespaceWithYaml(nsname, wf.data)
	if err != nil {
		log.Printf("Creating/Updating Namespace is failed because %v. data: %s\n", err, wf.data)
		return
	}
	log.Printf("Ns will be created/updated: %v\n", ns)
}
