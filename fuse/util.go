package fuse

import (
	"encoding/json"
	"github.com/ghodss/yaml"

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

func GenYaml(obj *runtime.Object) ([]byte, error) {
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
