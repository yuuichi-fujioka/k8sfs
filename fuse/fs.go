package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"flag"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type K8sFs struct {
	pathfs.FileSystem
	root DirEntry
}

func NewK8sFs() *K8sFs {
	nsDir := NewNamespacesDir()
	// TODO: begin testcode
	go func() {
		wi, err := clientset.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				return
			}

			switch ev.Type {
			case watch.Added:
				log.Println("Added")
				nsDir.AddNamespace(&ev.Object)
			case watch.Modified:
				// Update
				nsDir.UpdateNamespace(&ev.Object)
			case watch.Deleted:
				// Delete
				nsDir.DeleteNamespace(&ev.Object)
			}
		}
	}()
	// TODO: finish testcode
	return &K8sFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		root:       nsDir,
	}
}

func (me *K8sFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("fs/GetAttr: %s\n", name)
	f := GetFile(me.root, name)
	if f == nil {
		return nil, fuse.ENOENT
	}
	attr := &fuse.Attr{}
	status := f.GetAttr(attr)
	return attr, status
}

func (me *K8sFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("fs/OpenDir: %s\n", name)
	c, code = me.root.GetDir(name).OpenDir()
	return
}

func (me *K8sFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("fs/Open: %s\n", name)
	f := GetFile(me.root, name)
	return f, fuse.OK
}

func (me *K8sFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	log.Printf("fs/Mkdir: %s\n", name)

	names := strings.Split(name, "/")
	parentName := strings.Join(names[:len(names)-1], "/")
	return me.root.GetDir(parentName).Mkdir(names[len(names)-1], mode)
}

func (me *K8sFs) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	log.Printf("fs/Unlink: %s\n", name)
	names := strings.Split(name, "/")
	parentName := strings.Join(names[:len(names)-1], "/")
	return me.root.GetDir(parentName).Unlink(names[len(names)-1])
}

func (me *K8sFs) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	log.Printf("fs/Rmdir: %s\n", name)
	return me.root.GetDir(name).Rmdir()
}

func (me *K8sFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("fs/Create: %s %o %o\n", name, flags, mode)
	return nil, fuse.ENOSYS
}

func Serve(mountPoint string) {
	k8sfs := NewK8sFs()
	nfs := pathfs.NewPathNodeFs(k8sfs, nil)
	server, _, err := nodefs.MountRoot(mountPoint, nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}

// TODO: begin testcode
var clientset *kubernetes.Clientset

func TestMain() {
	clientset = k8s.GenClientSetFromFlags()

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  k8sfs MOUNTPOINT")
	}
	log.Printf("argments: %v\n", flag.Args())

	Serve(flag.Arg(0))
}

// TODO: finish testcode
