package fuse

import (
	"log"
	"strings"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	"k8s.io/apimachinery/pkg/runtime"
)

func NewObjFile(obj runtime.Object, meta *metaObj, handler WFReleaseHandler) *writableFile {
	yaml, err := GenYaml(obj)
	if err != nil {
		panic("!!!")
	}

	f := &writableFile{
		Name:    meta.Name + ".yaml",
		File:    nodefs.NewDefaultFile(),
		data:    yaml,
		ctime:   uint64(meta.GetCreationTimestamp().Unix()),
		handler: handler,
	}
	return f
}

func (f *writableFile) Update(obj runtime.Object, meta *metaObj) {
	yaml, err := GenYaml(obj)
	if err != nil {
		panic("!!!")
	}
	f.data = yaml
	f.ctime = uint64(meta.GetCreationTimestamp().Unix())
}

type WFReleaseHandler interface {
	HandleRelease(*writableFile)
}

type writableFile struct {
	Name string
	nodefs.File
	data    []byte
	ctime   uint64
	handler WFReleaseHandler
}

func NewWFile(name string, handler WFReleaseHandler) *writableFile {
	f := &writableFile{
		Name:    name,
		File:    nodefs.NewDefaultFile(),
		data:    make([]byte, 0),
		ctime:   uint64(time.Now().Unix()),
		handler: handler,
	}
	return f
}

func (f *writableFile) GetFile() nodefs.File {
	return f
}

func (f *writableFile) GetAttr(out *fuse.Attr) fuse.Status {
	out.Size = uint64(len(f.data))
	out.Ctime = f.ctime
	out.Mtime = f.ctime
	out.Atime = f.ctime
	out.Mode = fuse.S_IFREG | 0644
	return fuse.OK
}

func (f *writableFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	log.Printf("Write: %s %d\n", f.Name, off)
	end := off + int64(len(data))
	if int64(len(f.data)) < end {
		f.data = append(f.data, make([]byte, end-int64(len(f.data)))...)
	}
	copy(f.data[off:], data)
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
	if strings.HasSuffix(f.Name, ".yaml") {
		log.Printf("Relase: %s %s\n", f.Name, f.data)
	} else {
		log.Printf("Relase: %s\n", f.Name)
	}
	if f.handler != nil {
		f.handler.HandleRelease(f)
	}
}

func (f *writableFile) Truncate(size uint64) fuse.Status {

	if uint64(len(f.data)) < size {
		f.data = append(f.data, make([]byte, size-uint64(len(f.data)))...)
	}
	for i := uint64(0); i < size; i++ {
		f.data[i] = 0
	}
	return fuse.OK
}

func (f *writableFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	log.Printf("Allocate: %s\n", f.Name)
	return fuse.ENOSYS
}
func (f *writableFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	log.Printf("Utimens: %s\n", f.Name)
	return fuse.ENOSYS
}
func (f *writableFile) Fsync(flags int) (code fuse.Status) {
	log.Printf("Fsync: %s\n", f.Name)
	return fuse.ENOSYS
}
func (f *writableFile) Flush() fuse.Status {
	log.Printf("Flush: %s\n", f.Name)
	return fuse.ENOSYS
}
func (f *writableFile) Flock(flags int) fuse.Status {
	log.Printf("Flock: %s\n", f.Name)
	return fuse.ENOSYS
}
