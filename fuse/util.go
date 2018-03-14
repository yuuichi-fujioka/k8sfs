package fuse

import (
	"log"
	"strings"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	"encoding/json"
	"github.com/ghodss/yaml"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type metaObj struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

func NewMetaObj(typeMeta *metav1.TypeMeta, objectMeta *metav1.ObjectMeta) *metaObj {
	return &metaObj{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
	}
}

func GenYaml(obj runtime.Object) ([]byte, error) {
	jsondata, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	yaml, err := yaml.JSONToYAML(jsondata)
	if err != nil {
		return nil, err
	}
	return yaml, nil
}

func GetFile(d DirEntry, name string) nodefs.File {
	log.Printf("XXX GetFile: %s\n", name)
	if name == "" {
		return d.GetFile()
	}

	names := strings.Split(name, "/")
	for k, child := range d.GetChildDirs() {
		if k == names[0] {
			return GetFile(child, strings.Join(names[1:], "/"))
		}
	}
	for k, child := range d.GetChildFiles() {
		if k == names[0] {
			return child
		}
	}
	return nil
}

func SetAttrTime(obj *metaObj, out *fuse.Attr) {
	ctime := uint64(obj.GetCreationTimestamp().Unix())
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
}

func SetAttrTimeCluster(out *fuse.Attr) {
	ctime := k8s.ClusterCreatedAt
	out.Ctime = ctime
	out.Mtime = ctime
	out.Atime = ctime
}
