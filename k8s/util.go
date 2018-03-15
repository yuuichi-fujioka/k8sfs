package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Clientset *kubernetes.Clientset
var ClusterCreatedAt uint64

func InitFromArg() {
	Clientset = GenClientSetFromFlags()
	initClusterCreateAt()
}

func Init(kubeconfig string) {
	Clientset = GenClientSet(kubeconfig)
	initClusterCreateAt()
}

func initClusterCreateAt() {

	ns, err := Clientset.CoreV1().Namespaces().Get("kube-system", metav1.GetOptions{})
	if err != nil {
		panic("AAAAAAAAAAAAAAAAA")
	}
	ClusterCreatedAt = uint64(ns.CreationTimestamp.Unix())
}
