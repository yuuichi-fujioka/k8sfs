package fuse

import (
	"log"

	"github.com/yuuichi-fujioka/k8sfs/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func watchPods(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/po")
	poDir := dir.(*podsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Pods(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchServices(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/svc")
	svcDir := dir.(*servicesDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Services(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchConfigMaps(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/cm")
	cmDir := dir.(*configMapsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ConfigMaps(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchDeployments(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/deploy")
	deployDir := dir.(*deploymentsDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().Deployments(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchEndpoints(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/ep")
	epDir := dir.(*endpointsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Endpoints(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchEvents(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/ev")
	evDir := dir.(*eventsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Events(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchIngresses(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/ing")
	ingDir := dir.(*ingresssDir)
	for {
		wi, err := k8s.Clientset.ExtensionsV1beta1().Ingresses(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchPersistentVolumeClaims(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/pvc")
	pvcDir := dir.(*persistentVolumeClaimsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().PersistentVolumeClaims(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchReplicationControllers(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/rc")
	rcDir := dir.(*replicationControllersDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ReplicationControllers(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchServiceAccounts(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/sa")
	saDir := dir.(*serviceAccountsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().ServiceAccounts(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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

func watchSecrets(ns string) {
	nsDir := Fs.root.(*namespacesDir)
	dir := nsDir.GetDir(ns + "/secrets")
	secretsDir := dir.(*secretsDir)
	for {
		wi, err := k8s.Clientset.CoreV1().Secrets(ns).Watch(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		ch := wi.ResultChan()

		for {
			ev, ok := <-ch
			if !ok {
				break
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
