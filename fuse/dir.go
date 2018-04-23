package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type DirEntry interface {
	GetAttr(out *fuse.Attr) fuse.Status
	OpenDir() (c []fuse.DirEntry, code fuse.Status)
	GetFile() nodefs.File
	GetDir(name string) DirEntry
	GetChildDirs() map[string]DirEntry
	GetChildFiles() map[string]nodefs.File

	Unlink(name string) (code fuse.Status)
	Mkdir(name string, mode uint32) fuse.Status
	Rmdir(name string) (code fuse.Status)

	Rename(oldName, newName string) (code fuse.Status)
	RenameOnSameDir(oldName, newName string) (code fuse.Status)
	Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status)
}

type defaultDir struct {
	files    map[string]*writableFile
	dirs     map[string]DirEntry
	tmpFiles map[string]*writableFile
	tmpDirs  map[string]DirEntry
}

func NewDefaultDir() defaultDir {
	return defaultDir{
		files:    map[string]*writableFile{},
		dirs:     map[string]DirEntry{},
		tmpFiles: map[string]*writableFile{},
		tmpDirs:  map[string]DirEntry{},
	}
}

func (f *defaultDir) OpenDir() (c []fuse.DirEntry, code fuse.Status) {
	c = []fuse.DirEntry{}
	for k, _ := range f.dirs {
		c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFDIR})
	}
	for k, _ := range f.tmpDirs {
		c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFDIR})
	}
	for k, _ := range f.GetChildFiles() {
		c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFREG})
	}
	return c, fuse.OK
}

func (f *defaultDir) GetChildDirs() map[string]DirEntry {
	dirs := map[string]DirEntry{}
	for k, v := range f.tmpDirs {
		dirs[k] = v
	}
	for k, v := range f.dirs {
		dirs[k] = v
	}
	return dirs
}

func (f *defaultDir) GetChildFiles() map[string]nodefs.File {
	files := map[string]nodefs.File{}
	for k, n := range f.tmpFiles {
		files[k] = n
	}
	for k, n := range f.files {
		files[k] = n
	}
	return files
}

func (f *defaultDir) RemoveTmpFile(name string) fuse.Status {
	if _, ok := f.tmpFiles[name]; ok {

		delete(f.tmpFiles, name)
		return fuse.OK
	}
	return fuse.ENOENT
}

func (f *defaultDir) AddTmpFile(name string, handler WFReleaseHandler) nodefs.File {
	tmp := NewWFile(name, handler)
	f.tmpFiles[name] = tmp
	return tmp
}

func (f *defaultDir) RemoveTmpDir(name string) fuse.Status {
	if _, ok := f.tmpDirs[name]; ok {

		delete(f.tmpDirs, name)
		return fuse.OK
	}
	return fuse.ENOENT
}

func (f *defaultDir) AddTmpDir(name string) DirEntry {
	tmp := NewTmpDir()
	f.tmpDirs[name] = tmp
	return tmp
}

func (f *defaultDir) Rename(oldName, newName string) (code fuse.Status) {
	return fuse.ENOSYS
}

func (f *defaultDir) RenameOnSameDir(oldName, newName string) (code fuse.Status) {
	return fuse.ENOSYS
}

type tmpDir struct {
	nodefs.File
	defaultDir
}

func NewTmpDir() *tmpDir {
	return &tmpDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
	}
}

func (f *tmpDir) GetAttr(out *fuse.Attr) fuse.Status {
	SetAttrTimeCluster(out)
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	return fuse.OK
}

func (f *tmpDir) GetFile() nodefs.File {
	return f
}

func (f *tmpDir) GetDir(name string) DirEntry {
	if name == "" {
		return f
	}

	names := strings.Split(name, "/")

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	for k, child := range f.tmpDirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *tmpDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, "tmpdir")
	code = f.RemoveTmpFile(name)
	return code
}

func (f *tmpDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, "tmpdir")
	f.AddTmpDir(name)
	return fuse.OK
}

func (f *tmpDir) Rmdir(name string) (code fuse.Status) {
	log.Printf("Rmdir: %s", "tmpdir")
	code = f.RemoveTmpDir(name)
	return code
}

func (f *tmpDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, "tmpdir", flags, mode)
	return f.AddTmpFile(name, nil), fuse.OK
}

func (f *tmpDir) RenameOnSameDir(oldName, newName string) (code fuse.Status) {
	log.Printf("RenameOnSameDir: %s to %s at %s", oldName, newName, "tmpdir")

	if _, ok := f.tmpFiles[oldName]; ok {
		f.tmpFiles[newName] = f.tmpFiles[oldName]
		delete(f.tmpFiles, oldName)
		f.tmpFiles[newName].Name = newName
	}

	if _, ok := f.tmpDirs[oldName]; ok {
		f.tmpDirs[newName] = f.tmpDirs[oldName]
		delete(f.tmpDirs, oldName)
	}
	return fuse.OK
}

type objDir struct {
	defaultDir
	metaObj
}

func NewObjDir(obj *metaObj) objDir {
	return objDir{
		defaultDir: NewDefaultDir(),
		metaObj:    *obj,
	}
}

func (f *objDir) GetAttr(out *fuse.Attr) fuse.Status {
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	SetAttrTime(&f.metaObj, out)
	return fuse.OK
}
