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
	GetChildFiles() []*objFile
}

type defaultDir struct {
	files []*objFile
	dirs  []DirEntry
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
	return c, fuse.OK
}

func (f *defaultDir) GetChildDirs() []DirEntry {
	return f.dirs
}

func (f *defaultDir) GetChildFiles() []*objFile {
	return f.files
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
