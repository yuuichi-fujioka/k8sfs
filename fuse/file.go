package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"k8s.io/apimachinery/pkg/runtime"
)

type objFile struct {
	nodefs.File
	metaObj
	yaml []byte
	Ext  string
}

func NewObjFile(obj *runtime.Object, meta *metaObj) *objFile {
	yaml, err := GenYaml(obj)
	if err != nil {
		panic("!!!")
	}

	return &objFile{
		File:    nodefs.NewDefaultFile(),
		yaml:    yaml,
		metaObj: *meta,
		Ext:     "yaml",
	}
}

func (f *objFile) Read(buf []byte, off int64) (res fuse.ReadResult, code fuse.Status) {
	end := int(off) + int(len(buf))
	if end > len(f.yaml) {
		end = len(f.yaml)
	}

	return fuse.ReadResultData(f.yaml[off:end]), fuse.OK
}

func (f *objFile) GetAttr(out *fuse.Attr) fuse.Status {
	ctime := uint64(f.GetCreationTimestamp().Unix())
	out.Size = uint64(len(f.yaml))
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
	out.Mode = fuse.S_IFREG | 0644
	return fuse.OK
}

func (f *objFile) Update(obj *runtime.Object, meta *metaObj) {
	yaml, err := GenYaml(obj)
	if err != nil {
		panic("!!!")
	}
	f.yaml = yaml
	f.metaObj = *meta
}

type writableFile struct {
	Name string
	nodefs.File
	data []byte
}

func NewWFile(name string) *writableFile {
	f := &writableFile{
		Name: name,
		File: nodefs.NewDefaultFile(),
		data: make([]byte, 0),
	}
	return f
}

func (f *writableFile) GetAttr(out *fuse.Attr) fuse.Status {
	ctime := uint64(0)
	out.Size = uint64(len(f.data))
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
	out.Mode = fuse.S_IFREG | 0644
	return fuse.OK
}

func (f *writableFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	log.Printf("Write: %s %d\n", f.Name, off)
	f.data = append(f.data[:off], data...)
	return uint32(len(data)), fuse.OK
}

func (f *writableFile) Read(buf []byte, off int64) (res fuse.ReadResult, code fuse.Status) {
	log.Printf("Read: %s %d\n", f.Name, off)
	end := int(off) + int(len(buf))
	if end > len(f.data) {
		end = len(f.data)
	}

	return fuse.ReadResultData(f.data[off:end]), fuse.OK
}

func (f *writableFile) InnerFile() nodefs.File {
	return f.File
}

func (f *writableFile) Release() {
	if !strings.HasPrefix(f.Name, ".") {
		log.Printf("Relase: %s %s\n", f.Name, f.data)
	} else {
		log.Printf("Relase: %s\n", f.Name)
	}
}
