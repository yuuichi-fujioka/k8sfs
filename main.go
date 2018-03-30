package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/yuuichi-fujioka/k8sfs/fuse"
	"github.com/yuuichi-fujioka/k8sfs/k8s"
)

func main() {
	startHandlingSignal()

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	namespace := flag.String("namespace", "", "top level namespace. if this is blank, all namespaces will be mount")
	readOnly := flag.Bool("readonly", false, "read only mode.")

	flag.Parse()

	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n  k8sfs MOUNTPOINT")
	}
	log.Printf("argments: %v\n", flag.Args())

	k8s.Init(*kubeconfig)
	fuse.TestMain(flag.Arg(0), *namespace, *readOnly)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
