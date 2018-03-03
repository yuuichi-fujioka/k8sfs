package fuse

import (
	"log"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type DirEntry interface {
	GetName() string
	GetAttr(out *fuse.Attr) fuse.Status
	OpenDir() (c []fuse.DirEntry, code fuse.Status)
}

type defaultDir struct {
	nodefs.File

	files []*objFile
	dirs  []DirEntry
}

func NewDefaultDir() defaultDir {
	return defaultDir{
		File:  nodefs.NewDefaultFile(),
		files: make([]*objFile, 0),
		dirs:  make([]DirEntry, 0),
	}
}

func (f *defaultDir) OpenDir() (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("XXX OpenDir: \n")
	c = []fuse.DirEntry{}
	for _, object := range f.dirs {
		c = append(c, fuse.DirEntry{Name: object.GetName(), Mode: fuse.S_IFDIR})
	}
	for _, object := range f.files {
		c = append(c, fuse.DirEntry{Name: object.Name + "." + object.Ext, Mode: fuse.S_IFREG})
	}
	return c, fuse.OK
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
	ctime := uint64(f.GetCreationTimestamp().Unix())
	out.Size = 4096 // block size?
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
	out.Mode = fuse.S_IFDIR | 0755
	return fuse.OK
}

func (f *objDir) GetName() string {
	return f.Name
}
