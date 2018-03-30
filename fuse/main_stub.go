package fuse

import (
	"github.com/yuuichi-fujioka/k8sfs/k8s"

	corev1 "k8s.io/api/core/v1"
)

func TestMain(mountpoint, namespace string, readonly bool) {
	topLevelNamespace = namespace
	if topLevelNamespace == "" {
		go watchAllNs()
	} else {
		go watchNs()
	}

	readOnlyMode = readonly

	Serve(mountpoint)
}

func watchNs() {
	nsw := NewNsWatcher(topLevelNamespace)
	nsw.StartAll()
}

func watchAllNs() {
	watchers := map[string]*nsWatcher{}
	nsDir := Fs.root.(*namespacesDir)
	k8s.WatchNamespaces(
		func(ns *corev1.Namespace) {
			nsDir.AddNamespace(ns)

			nsw, ok := watchers[ns.Name]
			if !ok {
				nsw = NewNsWatcher(ns.Name)
				watchers[ns.Name] = nsw
			} else {
				nsw.StopAll()
			}
			nsw.StartAll()
		},
		func(old, new *corev1.Namespace) {
			nsDir.UpdateNamespace(new)
		},
		func(ns *corev1.Namespace) {
			nsw, ok := watchers[ns.Name]
			if ok {
				nsw.StopAll()
			}
			nsDir.DeleteNamespace(ns)
		},
	)
}
