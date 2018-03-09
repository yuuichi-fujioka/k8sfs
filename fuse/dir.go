package fuse

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type DirEntry interface {
	GetName() string
	GetAttr(out *fuse.Attr) fuse.Status
	OpenDir() (c []fuse.DirEntry, code fuse.Status)
	GetFile() nodefs.File
	GetDir(name string) DirEntry
	GetChildDirs() []DirEntry
	GetChildFiles() []namedFile

	Unlink(name string) (code fuse.Status)
	Mkdir(name string, mode uint32) fuse.Status
	Rmdir() (code fuse.Status)

	Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status)
}

type namedFile interface {
	GetName() string
	GetFile() nodefs.File
}

type defaultDir struct {
	files    []*objFile
	dirs     []DirEntry
	tmpFiles []*writableFile
}

func NewDefaultDir() defaultDir {
	return defaultDir{
		files: make([]*objFile, 0),
		dirs:  make([]DirEntry, 0),
	}
}

func (f *defaultDir) OpenDir() (c []fuse.DirEntry, code fuse.Status) {
	c = []fuse.DirEntry{}
	for _, object := range f.dirs {
		c = append(c, fuse.DirEntry{Name: object.GetName(), Mode: fuse.S_IFDIR})
	}
	for _, object := range f.files {
		c = append(c, fuse.DirEntry{Name: object.Name + "." + object.Ext, Mode: fuse.S_IFREG})
	}
	for _, object := range f.tmpFiles {
		c = append(c, fuse.DirEntry{Name: object.Name, Mode: fuse.S_IFREG})
	}
	return c, fuse.OK
}

func (f *defaultDir) GetChildDirs() []DirEntry {
	return f.dirs
}

func (f *defaultDir) GetChildFiles() []namedFile {
	files := make([]namedFile, 0)
	for _, n := range f.files {
		files = append(files, n)
	}
	for _, n := range f.tmpFiles {
		files = append(files, n)
	}
	return files
}

func (f *defaultDir) RemoveTmpFile(name string) fuse.Status {
	for i, tmp := range f.tmpFiles {
		if tmp.GetName() == name {
			f.tmpFiles = append(f.tmpFiles[:i], f.tmpFiles[i+1:]...)
			return fuse.OK
		}
	}
	return fuse.ENOENT
}

func (f *defaultDir) AddTmpFile(name string) nodefs.File {
	tmp := NewWFile(name)
	f.tmpFiles = append(f.tmpFiles, tmp)
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
	out.Mode = fuse.S_IFDIR | 0755
	SetAttrTime(&f.metaObj, out)
	return fuse.OK
}

func (f *objDir) GetName() string {
	return f.Name
}
