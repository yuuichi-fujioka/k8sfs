package fuse

import (
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type nsWatcher struct {
	Namespace     string
	closeChannels map[string](chan bool)
}

func NewNsWatcher(namespace string) *nsWatcher {
	return &nsWatcher{
		Namespace:     namespace,
		closeChannels: map[string](chan bool){},
	}
}

func (me *nsWatcher) StartAll() {
	go me.watchPods()
	go me.watchServices()
	go me.watchConfigMaps()
	go me.watchDeployments()
	go me.watchEndpoints()
	go me.watchEvents()
	go me.watchIngresses()
	go me.watchPersistentVolumeClaims()
	go me.watchReplicationControllers()
	go me.watchServiceAccounts()
	go me.watchSecrets()
	go me.watchDaemonSets()
	go me.watchReplicaSets()
}

func (me *nsWatcher) StopAll() {
	for _, ch := range me.closeChannels {
		ch <- true
	}
}

func (me *nsWatcher) watchPods() {
	log.Printf("[Watch] start watchPods/%s\n", me.Namespace)

	me.closeChannels["po"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("po")
	poDir := dir.(*podsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Pods(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["po"]:
				log.Printf("[Watch] finish watchPods/%s\n", me.Namespace)
				delete(me.closeChannels, "po")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] po/Added on %s\n", me.Namespace)
					poDir.AddPod(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] po/Modified on %s\n", me.Namespace)
					poDir.UpdatePod(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] po/Deleted on %s\n", me.Namespace)
					poDir.DeletePod(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchServices() {
	log.Printf("[Watch] start watchServices/%s\n", me.Namespace)

	me.closeChannels["svc"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("svc")
	svcDir := dir.(*servicesDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Services(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["svc"]:
				log.Printf("[Watch] finish watchServices/%s\n", me.Namespace)
				delete(me.closeChannels, "svc")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] svc/Added on %s\n", me.Namespace)
					svcDir.AddService(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] svc/Modified on %s\n", me.Namespace)
					svcDir.UpdateService(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] svc/Deleted on %s\n", me.Namespace)
					svcDir.DeleteService(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchConfigMaps() {
	log.Printf("[Watch] start watchConfigMaps/%s\n", me.Namespace)

	me.closeChannels["cm"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("cm")
	cmDir := dir.(*configMapsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ConfigMaps(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["cm"]:
				log.Printf("[Watch] finish watchConfigMaps/%s\n", me.Namespace)
				delete(me.closeChannels, "cm")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] cm/Added on %s\n", me.Namespace)
					cmDir.AddConfigMap(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] cm/Modified on %s\n", me.Namespace)
					cmDir.UpdateConfigMap(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] cm/Deleted on %s\n", me.Namespace)
					cmDir.DeleteConfigMap(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchDeployments() {
	log.Printf("[Watch] start watchDeployments/%s\n", me.Namespace)

	me.closeChannels["deploy"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("deploy")
	deployDir := dir.(*deploymentsDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().Deployments(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["deploy"]:
				log.Printf("[Watch] finish watchDeployments/%s\n", me.Namespace)
				delete(me.closeChannels, "deploy")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] deploy/Added on %s\n", me.Namespace)
					deployDir.AddDeployment(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] deploy/Modified on %s\n", me.Namespace)
					deployDir.UpdateDeployment(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] deploy/Deleted on %s\n", me.Namespace)
					deployDir.DeleteDeployment(ev.Object)
				}
			}
		}
	}

}

func (me *nsWatcher) watchEndpoints() {
	log.Printf("[Watch] start watchEndpoints/%s\n", me.Namespace)

	me.closeChannels["ep"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ep")
	epDir := dir.(*endpointsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Endpoints(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["ep"]:
				log.Printf("[Watch] finish watchEndpoints/%s\n", me.Namespace)
				delete(me.closeChannels, "ep")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] ep/Added on %s\n", me.Namespace)
					epDir.AddEndpoints(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] ep/Modified on %s\n", me.Namespace)
					epDir.UpdateEndpoints(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] ep/Deleted on %s\n", me.Namespace)
					epDir.DeleteEndpoints(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchEvents() {
	log.Printf("[Watch] start watchEvents/%s\n", me.Namespace)

	me.closeChannels["ev"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ev")
	evDir := dir.(*eventsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Events(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["ev"]:
				log.Printf("[Watch] finish watchEvents/%s\n", me.Namespace)
				delete(me.closeChannels, "ev")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] ev/Added on %s\n", me.Namespace)
					evDir.AddEvent(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] ev/Modified on %s\n", me.Namespace)
					evDir.UpdateEvent(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] ev/Deleted on %s\n", me.Namespace)
					evDir.DeleteEvent(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchIngresses() {
	log.Printf("[Watch] start watchIngresses/%s\n", me.Namespace)

	me.closeChannels["ing"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ing")
	ingDir := dir.(*ingresssDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().Ingresses(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["ing"]:
				log.Printf("[Watch] finish watchIngresses/%s\n", me.Namespace)
				delete(me.closeChannels, "ing")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] ing/Added on %s\n", me.Namespace)
					ingDir.AddIngress(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] ing/Modified on %s\n", me.Namespace)
					ingDir.UpdateIngress(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] ing/Deleted on %s\n", me.Namespace)
					ingDir.DeleteIngress(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchPersistentVolumeClaims() {
	log.Printf("[Watch] start watchPersistentVolumeClaims/%s\n", me.Namespace)

	me.closeChannels["pvc"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("pvc")
	pvcDir := dir.(*persistentVolumeClaimsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["pvc"]:
				log.Printf("[Watch] finish watchPersistentVolumeClaims/%s\n", me.Namespace)
				delete(me.closeChannels, "pvc")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] pvc/Added on %s\n", me.Namespace)
					pvcDir.AddPersistentVolumeClaim(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] pvc/Modified on %s\n", me.Namespace)
					pvcDir.UpdatePersistentVolumeClaim(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] pvc/Deleted on %s\n", me.Namespace)
					pvcDir.DeletePersistentVolumeClaim(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchReplicationControllers() {
	log.Printf("[Watch] start watchReplicationControllers/%s\n", me.Namespace)

	me.closeChannels["rc"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("rc")
	rcDir := dir.(*replicationControllersDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ReplicationControllers(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["rc"]:
				log.Printf("[Watch] finish watchReplicationControllers/%s\n", me.Namespace)
				delete(me.closeChannels, "rc")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] rc/Added on %s\n", me.Namespace)
					rcDir.AddReplicationController(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] rc/Modified on %s\n", me.Namespace)
					rcDir.UpdateReplicationController(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] rc/Deleted on %s\n", me.Namespace)
					rcDir.DeleteReplicationController(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchServiceAccounts() {
	log.Printf("[Watch] start watchServiceAccounts/%s\n", me.Namespace)

	me.closeChannels["sa"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("sa")
	saDir := dir.(*serviceAccountsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ServiceAccounts(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["sa"]:
				log.Printf("[Watch] finish watchServiceAccounts/%s\n", me.Namespace)
				delete(me.closeChannels, "sa")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] sa/Added on %s\n", me.Namespace)
					saDir.AddServiceAccount(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] sa/Modified on %s\n", me.Namespace)
					saDir.UpdateServiceAccount(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] sa/Deleted on %s\n", me.Namespace)
					saDir.DeleteServiceAccount(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchSecrets() {
	log.Printf("[Watch] start watchSecrets/%s\n", me.Namespace)

	me.closeChannels["secrets"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("secrets")
	secretsDir := dir.(*secretsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Secrets(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["secrets"]:
				log.Printf("[Watch] finish watchSecrets/%s\n", me.Namespace)
				delete(me.closeChannels, "secrets")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] secrets/Added on %s\n", me.Namespace)
					secretsDir.AddSecret(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] secrets/Modified on %s\n", me.Namespace)
					secretsDir.UpdateSecret(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] secrets/Deleted on %s\n", me.Namespace)
					secretsDir.DeleteSecret(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchDaemonSets() {
	log.Printf("[Watch] start watchDaemonSets/%s\n", me.Namespace)

	me.closeChannels["ds"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ds")
	dsDir := dir.(*daemonSetsDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().DaemonSets(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["ds"]:
				log.Printf("[Watch] finish watchDaemonSets/%s\n", me.Namespace)
				delete(me.closeChannels, "ds")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] ds/Added on %s\n", me.Namespace)
					dsDir.AddDaemonSet(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] ds/Modified on %s\n", me.Namespace)
					dsDir.UpdateDaemonSet(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] ds/Deleted on %s\n", me.Namespace)
					dsDir.DeleteDaemonSet(ev.Object)
				}
			}
		}
	}

}

func (me *nsWatcher) watchReplicaSets() {
	log.Printf("[Watch] start watchReplicaSets/%s\n", me.Namespace)

	me.closeChannels["rs"] = make(chan bool)
	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("rs")
	rsDir := dir.(*replicaSetsDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().ReplicaSets(me.Namespace).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

	loop:
		for {
			select {
			case <-me.closeChannels["rs"]:
				log.Printf("[Watch] finish watchReplicaSets/%s\n", me.Namespace)
				delete(me.closeChannels, "rs")
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Printf("[Watch] rs/Added on %s\n", me.Namespace)
					rsDir.AddReplicaSet(ev.Object)

				case watch.Modified:
					// Update
					log.Printf("[Watch] rs/Modified on %s\n", me.Namespace)
					rsDir.UpdateReplicaSet(ev.Object)
				case watch.Deleted:
					// Delete
					log.Printf("[Watch] rs/Deleted on %s\n", me.Namespace)
					rsDir.DeleteReplicaSet(ev.Object)
				}
			}
		}
	}

}
