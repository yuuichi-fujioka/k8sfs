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
		Namespace: namespace,
		closeChannels: map[string](chan bool){
			"cm":      make(chan bool),
			"deploy":  make(chan bool),
			"ep":      make(chan bool),
			"ev":      make(chan bool),
			"ing":     make(chan bool),
			"po":      make(chan bool),
			"pvc":     make(chan bool),
			"rc":      make(chan bool),
			"sa":      make(chan bool),
			"secrets": make(chan bool),
			"svc":     make(chan bool),
		},
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
}

func (me *nsWatcher) StopAll() {
	me.closeChannels["cm"] <- true
	me.closeChannels["deploy"] <- true
	me.closeChannels["ep"] <- true
	me.closeChannels["ev"] <- true
	me.closeChannels["ing"] <- true
	me.closeChannels["po"] <- true
	me.closeChannels["pvc"] <- true
	me.closeChannels["rc"] <- true
	me.closeChannels["sa"] <- true
	me.closeChannels["secrets"] <- true
	me.closeChannels["svc"] <- true
}

func (me *nsWatcher) watchPods() {
	log.Printf("start watchPods/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/po")
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
				log.Printf("finish watchPods/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("po/Added")
					poDir.AddPod(ev.Object)

				case watch.Modified:
					// Update
					poDir.UpdatePod(ev.Object)
				case watch.Deleted:
					// Delete
					poDir.DeletePod(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchServices() {
	log.Printf("start watchServices/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/svc")
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
				log.Printf("finish watchServices/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("svc/Added")
					svcDir.AddService(ev.Object)

				case watch.Modified:
					// Update
					svcDir.UpdateService(ev.Object)
				case watch.Deleted:
					// Delete
					svcDir.DeleteService(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchConfigMaps() {
	log.Printf("start watchConfigMaps/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/cm")
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
				log.Printf("finish watchConfigMaps/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("cm/Added")
					cmDir.AddConfigMap(ev.Object)

				case watch.Modified:
					// Update
					cmDir.UpdateConfigMap(ev.Object)
				case watch.Deleted:
					// Delete
					cmDir.DeleteConfigMap(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchDeployments() {
	log.Printf("start watchDeployments/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/deploy")
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
				log.Printf("finish watchDeployments/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("deploy/Added")
					deployDir.AddDeployment(ev.Object)

				case watch.Modified:
					// Update
					deployDir.UpdateDeployment(ev.Object)
				case watch.Deleted:
					// Delete
					deployDir.DeleteDeployment(ev.Object)
				}
			}
		}
	}

}

func (me *nsWatcher) watchEndpoints() {
	log.Printf("start watchEndpoints/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/ep")
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
				log.Printf("finish watchEndpoints/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("ep/Added")
					epDir.AddEndpoints(ev.Object)

				case watch.Modified:
					// Update
					epDir.UpdateEndpoints(ev.Object)
				case watch.Deleted:
					// Delete
					epDir.DeleteEndpoints(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchEvents() {
	log.Printf("start watchEvents/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/ev")
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
				log.Printf("finish watchEvents/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("ev/Added")
					evDir.AddEvent(ev.Object)

				case watch.Modified:
					// Update
					evDir.UpdateEvent(ev.Object)
				case watch.Deleted:
					// Delete
					evDir.DeleteEvent(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchIngresses() {
	log.Printf("start watchIngresses/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/ing")
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
				log.Printf("finish watchIngresses/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("ing/Added")
					ingDir.AddIngress(ev.Object)

				case watch.Modified:
					// Update
					ingDir.UpdateIngress(ev.Object)
				case watch.Deleted:
					// Delete
					ingDir.DeleteIngress(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchPersistentVolumeClaims() {
	log.Printf("start watchPersistentVolumeClaims/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/pvc")
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
				log.Printf("finish watchPersistentVolumeClaims/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("pvc/Added")
					pvcDir.AddPersistentVolumeClaim(ev.Object)

				case watch.Modified:
					// Update
					pvcDir.UpdatePersistentVolumeClaim(ev.Object)
				case watch.Deleted:
					// Delete
					pvcDir.DeletePersistentVolumeClaim(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchReplicationControllers() {
	log.Printf("start watchReplicationControllers/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/rc")
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
				log.Printf("finish watchReplicationControllers/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("rc/Added")
					rcDir.AddReplicationController(ev.Object)

				case watch.Modified:
					// Update
					rcDir.UpdateReplicationController(ev.Object)
				case watch.Deleted:
					// Delete
					rcDir.DeleteReplicationController(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchServiceAccounts() {
	log.Printf("start watchServiceAccounts/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/sa")
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
				log.Printf("finish watchServiceAccounts/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("sa/Added")
					saDir.AddServiceAccount(ev.Object)

				case watch.Modified:
					// Update
					saDir.UpdateServiceAccount(ev.Object)
				case watch.Deleted:
					// Delete
					saDir.DeleteServiceAccount(ev.Object)
				}
			}
		}
	}
}

func (me *nsWatcher) watchSecrets() {
	log.Printf("start watchSecrets/%s", me.Namespace)

	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(me.Namespace + "/secrets")
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
				log.Printf("finish watchSecrets/%s", me.Namespace)
				return
			case ev, ok := <-ch:
				if !ok {
					break loop
				}

				switch ev.Type {
				case watch.Added:
					log.Println("secrets/Added")
					secretsDir.AddSecret(ev.Object)

				case watch.Modified:
					// Update
					secretsDir.UpdateSecret(ev.Object)
				case watch.Deleted:
					// Delete
					secretsDir.DeleteSecret(ev.Object)
				}
			}
		}
	}
}
