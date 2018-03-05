package fuse

import (
	"log"

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
			case watch.Deleted:
				// Delete
			}
		}
	}()
	return &K8sFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		root:       nsDir,
	}
}

func (me *K8sFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	log.Printf("GetAttr: %s\n", name)
	f := me.root.GetFile(name)
	if f == nil {
		return nil, fuse.ENOENT
	}
	attr := &fuse.Attr{}
	status := f.GetAttr(attr)
	return attr, status
}

func (me *K8sFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	log.Printf("OpenDir: %s\n", name)
	c, code = me.root.GetDir(name).OpenDir()
	return
}

func (me *K8sFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	log.Printf("Open: %s\n", name)
	f := me.root.GetFile(name)
	return f, fuse.OK
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
