package fuse

import (
	"log"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type namespaceDir struct {
	nodefs.File
	defaultDir

	metaObj
}

func NewNamespaceDir(obj runtime.Object) (string, *namespaceDir) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)

	d := &namespaceDir{
		File:       nodefs.NewDefaultFile(),
		defaultDir: NewDefaultDir(),
		metaObj:    *meta,
	}
	name, configMapsDir := NewConfigMapsDir(ns.Name)
	d.dirs[name] = configMapsDir
	name, deploymentsDir := NewDeploymentsDir(ns.Name)
	d.dirs[name] = deploymentsDir
	name, endpointsDir := NewEndpointsDir(ns.Name)
	d.dirs[name] = endpointsDir
	name, eventsDir := NewEventsDir(ns.Name)
	d.dirs[name] = eventsDir
	name, ingressesDir := NewIngressesDir(ns.Name)
	d.dirs[name] = ingressesDir
	name, persistentVolumeClaimsDir := NewPersistentVolumeClaimsDir(ns.Name)
	d.dirs[name] = persistentVolumeClaimsDir
	name, podsDir := NewPodsDir(ns.Name)
	d.dirs[name] = podsDir
	name, replicationControllersDir := NewReplicationControllersDir(ns.Name)
	d.dirs[name] = replicationControllersDir
	name, secretsDir := NewSecretsDir(ns.Name)
	d.dirs[name] = secretsDir
	name, serviceAccountsDir := NewServiceAccountsDir(ns.Name)
	d.dirs[name] = serviceAccountsDir
	name, servicesDir := NewServicesDir(ns.Name)
	d.dirs[name] = servicesDir
	name, daemonSetsDir := NewDaemonSetsDir(ns.Name)
	d.dirs[name] = daemonSetsDir
	name, replicaSetsDir := NewReplicaSetsDir(ns.Name)
	d.dirs[name] = replicaSetsDir
	return d.Name, d
}

func (f *namespaceDir) GetAttr(out *fuse.Attr) fuse.Status {
	out.Size = 4096 // block size?
	if readOnlyMode {
		out.Mode = fuse.S_IFDIR | 0555
	} else {
		out.Mode = fuse.S_IFDIR | 0755
	}
	SetAttrTime(&f.metaObj, out)
	return fuse.OK
}

func (f *namespaceDir) GetFile() nodefs.File {
	return f
}

func (f *namespaceDir) GetDir(name string) DirEntry {
	if name == "" {
		return f
	}

	names := strings.Split(name, "/")

	for k, child := range f.dirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	for k, child := range f.tmpDirs {
		if k == names[0] {
			return child.GetDir(strings.Join(names[1:], "/"))
		}
	}

	return nil
}

func (f *namespaceDir) Unlink(name string) (code fuse.Status) {
	log.Printf("Unlink: %s at %s", name, f.Name)
	// TODO
	return fuse.ENOSYS
}

func (f *namespaceDir) Mkdir(name string, mode uint32) fuse.Status {
	log.Printf("Mkdir: %s at %s", name, f.Name)

	if strings.HasPrefix(name, ".") {
		f.AddTmpDir(name)
		return fuse.OK
	} else {
		return fuse.ENOSYS
	}
}

func (f *namespaceDir) Rmdir(name string) (code fuse.Status) {
	log.Printf("Rmdir: %s", f.Name)
	// TODO
	code = f.RemoveTmpDir(name)
	return
}

func (f *namespaceDir) Create(name string, flags uint32, mode uint32) (file nodefs.File, code fuse.Status) {
	log.Printf("Create: %s on %s with 0x%x 0x%x", name, f.Name, flags, mode)
	// TODO
	return nil, fuse.ENOSYS
}

func (f *namespaceDir) Update(obj runtime.Object) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	f.metaObj = *meta
}

func NewNamespaceFile(obj runtime.Object, handler WFReleaseHandler) *writableFile {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	return NewObjFile(obj, meta, handler)
}

func UpdateNamespaceFile(f *writableFile, obj runtime.Object) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		panic("!!!!")
	}

	meta := NewMetaObj(&ns.TypeMeta, &ns.ObjectMeta)
	f.Update(obj, meta)
}
