package main

import (
	"github.com/yuuichi-fujioka/k8sfs/fuse"
	"github.com/yuuichi-fujioka/k8sfs/k8s"
)

func main() {
	k8s.InitFromArg()
	fuse.TestMain()
}
