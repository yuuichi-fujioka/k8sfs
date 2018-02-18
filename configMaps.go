package main

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type ConfigMapResource struct {
}

func (me *ConfigMapResource) MakeWatchInterface(nsname string) (watch.Interface, error) {

	wi, err := clientset.CoreV1().ConfigMaps(nsname).Watch(metav1.ListOptions{})
	return wi, err
}
func (me *ConfigMapResource) GetName(obj *runtime.Object) string {
	return me.getConfigMap(obj).GetName()
}
func (me *ConfigMapResource) GetCreationTimestamp(obj *runtime.Object) metav1.Time {
	return me.getConfigMap(obj).GetCreationTimestamp()
}

func (me *ConfigMapResource) getConfigMap(obj *runtime.Object) *corev1.ConfigMap {
	configMap, ok := (*obj).(*corev1.ConfigMap)
	if !ok {
		panic("unexpected object type")
	}
	return configMap
}
