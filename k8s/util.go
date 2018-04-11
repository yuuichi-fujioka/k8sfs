package k8s

import (
	"strings"
	"sync"

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

var deleteCancelSig map[string]bool
var mutexSigMap sync.RWMutex

func init() {
	deleteCancelSig = map[string]bool{}
	mutexSigMap = sync.RWMutex{}
}

func initDeleteCancelSignal(resource string, names ...string) {
	label := strings.Join(append([]string{resource}, names...), "/")

	mutexSigMap.Lock()
	defer mutexSigMap.Unlock()

	deleteCancelSig[label] = false
}

func flagDeleteCancelSignal(resource string, names ...string) {
	label := strings.Join(append([]string{resource}, names...), "/")

	mutexSigMap.Lock()
	defer mutexSigMap.Unlock()

	if _, ok := deleteCancelSig[label]; ok {
		deleteCancelSig[label] = true
	}
}

func isFlagedDeleteCancelSignal(resource string, names ...string) bool {
	label := strings.Join(append([]string{resource}, names...), "/")

	mutexSigMap.Lock()
	defer mutexSigMap.Unlock()

	v, ok := deleteCancelSig[label]
	if !ok {
		return false
	}
	delete(deleteCancelSig, label)
	return v
}
