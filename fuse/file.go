package fuse

import (
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
