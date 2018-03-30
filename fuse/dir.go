package fuse

import (
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
	Rmdir() (code fuse.Status)

	Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status)
}

type defaultDir struct {
	files    map[string]*writableFile
	dirs     map[string]DirEntry
	tmpFiles map[string]*writableFile
}

func NewDefaultDir() defaultDir {
	return defaultDir{
		files:    map[string]*writableFile{},
		dirs:     map[string]DirEntry{},
		tmpFiles: map[string]*writableFile{},
	}
}

func (f *defaultDir) OpenDir() (c []fuse.DirEntry, code fuse.Status) {
	c = []fuse.DirEntry{}
	for k, _ := range f.dirs {
		c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFDIR})
	}
	for k, _ := range f.GetChildFiles() {
		c = append(c, fuse.DirEntry{Name: k, Mode: fuse.S_IFREG})
	}
	return c, fuse.OK
}

func (f *defaultDir) GetChildDirs() map[string]DirEntry {
	return f.dirs
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
