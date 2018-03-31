package fuse

import (
	"log"
	"sync"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	corev1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type nsWatcher struct {
	Namespace     string
	closeChannels map[string](chan bool)
	lock          *sync.RWMutex
}

func NewNsWatcher(namespace string) *nsWatcher {
	return &nsWatcher{
		Namespace:     namespace,
		closeChannels: map[string](chan bool){},
		lock:          &sync.RWMutex{},
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

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("po")
	poDir := dir.(*podsDir)

	closeCh := k8s.WatchPods(
		me.Namespace,
		func(pod *corev1.Pod) {
			poDir.AddPod(pod)
		},
		func(oldpod, newpod *corev1.Pod) {
			poDir.UpdatePod(newpod)
		},
		func(pod *corev1.Pod) {
			poDir.DeletePod(pod)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["po"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "po")
	<-me.closeChannels["po"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchPods/%s\n", me.Namespace)
}

func (me *nsWatcher) watchServices() {
	log.Printf("[Watch] start watchServices/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("svc")
	svcDir := dir.(*servicesDir)

	closeCh := k8s.WatchServices(
		me.Namespace,
		func(service *corev1.Service) {
			svcDir.AddService(service)
		},
		func(oldservice, newservice *corev1.Service) {
			svcDir.UpdateService(newservice)
		},
		func(service *corev1.Service) {
			svcDir.DeleteService(service)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["svc"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "svc")
	<-me.closeChannels["svc"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchServices/%s\n", me.Namespace)

}

func (me *nsWatcher) watchConfigMaps() {
	log.Printf("[Watch] start watchConfigMaps/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("cm")
	cmDir := dir.(*configMapsDir)

	closeCh := k8s.WatchConfigMaps(
		me.Namespace,
		func(cm *corev1.ConfigMap) {
			cmDir.AddConfigMap(cm)
		},
		func(oldcm, newcm *corev1.ConfigMap) {
			cmDir.UpdateConfigMap(newcm)
		},
		func(cm *corev1.ConfigMap) {
			cmDir.DeleteConfigMap(cm)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["cm"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "cm")
	<-me.closeChannels["cm"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchConfigMaps/%s\n", me.Namespace)
}

func (me *nsWatcher) watchDeployments() {
	log.Printf("[Watch] start watchDeployments/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("deploy")
	deployDir := dir.(*deploymentsDir)

	closeCh := k8s.WatchDeployments(
		me.Namespace,
		func(deploy *v1beta1.Deployment) {
			deployDir.AddDeployment(deploy)
		},
		func(olddeploy, newdeploy *v1beta1.Deployment) {
			deployDir.UpdateDeployment(newdeploy)
		},
		func(deploy *v1beta1.Deployment) {
			deployDir.DeleteDeployment(deploy)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["deploy"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "deploy")
	<-me.closeChannels["deploy"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchDeployments/%s\n", me.Namespace)
}

func (me *nsWatcher) watchEndpoints() {
	log.Printf("[Watch] start watchEndpoints/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ep")
	epDir := dir.(*endpointsDir)

	closeCh := k8s.WatchEndpoints(
		me.Namespace,
		func(ep *corev1.Endpoints) {
			epDir.AddEndpoints(ep)
		},
		func(oldep, newep *corev1.Endpoints) {
			epDir.UpdateEndpoints(newep)
		},
		func(ep *corev1.Endpoints) {
			epDir.DeleteEndpoints(ep)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["ep"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "ep")
	<-me.closeChannels["ep"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchEndpoints/%s\n", me.Namespace)
}

func (me *nsWatcher) watchEvents() {
	log.Printf("[Watch] start watchEvents/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ev")
	evDir := dir.(*eventsDir)

	closeCh := k8s.WatchEvents(
		me.Namespace,
		func(ev *corev1.Event) {
			evDir.AddEvent(ev)
		},
		func(oldev, newev *corev1.Event) {
			evDir.UpdateEvent(newev)
		},
		func(ev *corev1.Event) {
			evDir.DeleteEvent(ev)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["ev"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "ev")
	<-me.closeChannels["ev"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchEvents/%s\n", me.Namespace)
}

func (me *nsWatcher) watchIngresses() {
	log.Printf("[Watch] start watchIngresses/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("ing")
	ingDir := dir.(*ingresssDir)

	closeCh := k8s.WatchIngresses(
		me.Namespace,
		func(ing *v1beta1.Ingress) {
			ingDir.AddIngress(ing)
		},
		func(olding, newing *v1beta1.Ingress) {
			ingDir.UpdateIngress(newing)
		},
		func(ing *v1beta1.Ingress) {
			ingDir.DeleteIngress(ing)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["ing"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "ing")
	<-me.closeChannels["ing"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchIngresses/%s\n", me.Namespace)
}

func (me *nsWatcher) watchPersistentVolumeClaims() {
	log.Printf("[Watch] start watchPersistentVolumeClaims/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("pvc")
	pvcDir := dir.(*persistentVolumeClaimsDir)

	closeCh := k8s.WatchPersistentVolumeClaims(
		me.Namespace,
		func(pvc *corev1.PersistentVolumeClaim) {
			pvcDir.AddPersistentVolumeClaim(pvc)
		},
		func(oldpvc, newpvc *corev1.PersistentVolumeClaim) {
			pvcDir.UpdatePersistentVolumeClaim(newpvc)
		},
		func(pvc *corev1.PersistentVolumeClaim) {
			pvcDir.DeletePersistentVolumeClaim(pvc)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["pvc"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "pvc")
	<-me.closeChannels["pvc"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchPersistentVolumeClaims/%s\n", me.Namespace)
}

func (me *nsWatcher) watchReplicationControllers() {
	log.Printf("[Watch] start watchReplicationControllers/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("rc")
	rcDir := dir.(*replicationControllersDir)

	closeCh := k8s.WatchReplicationControllers(
		me.Namespace,
		func(rc *corev1.ReplicationController) {
			rcDir.AddReplicationController(rc)
		},
		func(oldrc, newrc *corev1.ReplicationController) {
			rcDir.UpdateReplicationController(newrc)
		},
		func(rc *corev1.ReplicationController) {
			rcDir.DeleteReplicationController(rc)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["rc"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "rc")
	<-me.closeChannels["rc"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchReplicationControllers/%s\n", me.Namespace)
}

func (me *nsWatcher) watchServiceAccounts() {
	log.Printf("[Watch] start watchServiceAccounts/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("sa")
	saDir := dir.(*serviceAccountsDir)

	closeCh := k8s.WatchServiceAccounts(
		me.Namespace,
		func(sa *corev1.ServiceAccount) {
			saDir.AddServiceAccount(sa)
		},
		func(oldsa, newsa *corev1.ServiceAccount) {
			saDir.UpdateServiceAccount(newsa)
		},
		func(sa *corev1.ServiceAccount) {
			saDir.DeleteServiceAccount(sa)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["sa"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "sa")
	<-me.closeChannels["sa"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchServiceAccounts/%s\n", me.Namespace)
}

func (me *nsWatcher) watchSecrets() {
	log.Printf("[Watch] start watchSecrets/%s\n", me.Namespace)

	nsDir := GetNamespaceDir(me.Namespace)
	dir := nsDir.GetDir("secrets")
	secretDir := dir.(*secretsDir)

	closeCh := k8s.WatchSecrets(
		me.Namespace,
		func(secret *corev1.Secret) {
			secretDir.AddSecret(secret)
		},
		func(oldsecret, newsecret *corev1.Secret) {
			secretDir.UpdateSecret(newsecret)
		},
		func(secret *corev1.Secret) {
			secretDir.DeleteSecret(secret)
		},
	)
	defer func() { closeCh <- struct{}{} }()

	me.lock.Lock()
	me.closeChannels["secret"] = make(chan bool)
	me.lock.Unlock()
	defer delete(me.closeChannels, "secret")
	<-me.closeChannels["secret"] // wait until stopAll is called.

	log.Printf("[Watch] finish watchSecrets/%s\n", me.Namespace)
}

func (me *nsWatcher) watchDaemonSets() {
	log.Printf("[Watch] start watchDaemonSets/%s\n", me.Namespace)

	me.lock.Lock()
	me.closeChannels["ds"] = make(chan bool)
	me.lock.Unlock()
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

	me.lock.Lock()
	me.closeChannels["rs"] = make(chan bool)
	me.lock.Unlock()
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
