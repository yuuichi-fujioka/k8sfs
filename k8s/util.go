package k8s

import (
	"flag"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Clientset *kubernetes.Clientset
var ClusterCreatedAt uint64

func InitFromArg() {
	Clientset = GenClientSetFromFlags()

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  k8sfs MOUNTPOINT")
	}
	log.Printf("argments: %v\n", flag.Args())

	ns, err := Clientset.CoreV1().Namespaces().Get("kube-system", metav1.GetOptions{})
	if err != nil {
		panic("AAAAAAAAAAAAAAAAA")
	}
	ClusterCreatedAt = uint64(ns.CreationTimestamp.Unix())
}
