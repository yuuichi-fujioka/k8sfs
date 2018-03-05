package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type DirEntry interface {
	GetName() string
	GetAttr(out *fuse.Attr) fuse.Status
	OpenDir() (c []fuse.DirEntry, code fuse.Status)
	GetFile(name string) nodefs.File
	GetDir(name string) DirEntry
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
	log.Printf("XXX OpenDir: %s\n", f.GetName())
	c = []fuse.DirEntry{}
	for _, object := range f.dirs {
		c = append(c, fuse.DirEntry{Name: object.GetName(), Mode: fuse.S_IFDIR})
	}
	for _, object := range f.files {
		c = append(c, fuse.DirEntry{Name: object.Name + "." + object.Ext, Mode: fuse.S_IFREG})
	}
	return c, fuse.OK
}

func (f *defaultDir) GetFile(name string) nodefs.File {
	log.Printf("XXX GetFile: %s\n", name)
	if name == "" {
		return f
	}
	names := strings.Split(name, "/")
	for _, child := range f.dirs {
		if child.GetName() == names[0] {
			return child.GetFile(strings.Join(names[1:], "/"))
		}
	}
	for _, child := range f.files {
		if child.Name+"."+child.Ext == names[0] {
			return child
		}
	}
	return nil
}

func (f *defaultDir) GetName() string {
	return ""
}

func (f *defaultDir) GetDir(name string) DirEntry {
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

func (f *defaultDir) GetAttr(out *fuse.Attr) fuse.Status {
	ctime := uint64(0) // TODO cluster created at
	out.Size = 4096    // block size?
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
	out.Mode = fuse.S_IFDIR | 0755
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
	log.Println("asdf")
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
