#!/bin/bash

if ! which kubectl > /dev/null 2>/dev/null; then
    echo kubectl is not installed
    exit 1
fi

if [ $(mount | grep /mnt/k8s | grep pathfs | wc -l) -ne 1 ]; then
    echo k8sfs is not mounted at /mnt/k8s
    exit 2
fi

if [ -d /mnt/k8s/crudtest ]; then
    echo ns/crudtest is exist
    exit 3
fi

trap "kubectl delete ns crudtest" EXIT

echo TEST 1: Create a namespace with mkdir
mkdir /mnt/k8s/crudtest

sleep 5 # TODO fix way to handling that creting a namespace is finished.

if [ ! -f /mnt/k8s/crudtest.yaml ]; then
    echo ns/crudtest should be created if /mnt/k8s/crudtest/ is created.
    exit 4
fi

echo TEST 1: OK

echo TEST 2: Delete a namespace with rm xxx.yaml
rm -f /mnt/k8s/crudtest.yaml

sleep 30 # TODO fix way to handling that deleting a namespace is finished.

if [ ! -d /mnt/k8s/crudtest ]; then
    echo ns/crudtest should be deleted if /mnt/k8sf/crudtest.yaml is deleted.
    exit 5
fi

echo TEST 2: OK

echo TEST 3: Read some resources.

## setup
cat << __EOF__ | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: crudtest
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: crudtest
  namespace: crudtest
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: crudtest
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: crudtest
    namespace: crudtest
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: testcm
  namespace: crudtest
data:
  connection.count.max: "3"
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: crudtest
  name: crudtest
  namespace: crudtest
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: crudtest
    spec:
      containers:
      - image: nginx
        imagePullPolicy: Never
        name: crudtest
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: crudtest
  name: crudtest
  namespace: crudtest
spec:
  ports:
  - name: "80"
    port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: crudtest
  type: ClusterIP
---
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: crudtest-ingress
  namespace: crudtest
spec:
  rules:
  - host: example.com
    http:
      paths:
      - backend:
          serviceName: crudtest
          servicePort: 80
__EOF__

kubectl --namespace crudtest rollout status deployment/crudtest

tree /mnt/k8s/crudtest/

# /mnt/k8s/crudtest/
# ├── cm
# │   └── testcm.yaml
# ├── deploy
# │   └── crudtest.yaml
# ├── ds
# ├── ep
# │   └── crudtest.yaml
# ├── ev
# │   ├── crudtest.151ca6e82e85e166.yaml
# │   ├── crudtest-5c85c5f77.151ca6e8388ae213.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e83ae933f7.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e85320ff12.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e8ead31a0b.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e8fc95797e.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e90f9fcf9b.yaml
# │   ├── crudtest-5c85c5f77-58f45.151ca6e96258da08.yaml
# │   ├── crudtest-ingress.151ca6e856211ab5.yaml
# │   ├── crudtest-ingress.151ca6e85ab7d921.yaml
# │   ├── crudtest-ingress.151ca6e85ae95b64.yaml
# │   ├── crudtest-ingress.151ca6ec6d63c1ca.yaml
# │   ├── crudtest-ingress.151ca6ec6d9b7d8f.yaml
# │   └── crudtest-ingress.151ca6ec6dab0b4b.yaml
# ├── ing
# │   └── crudtest-ingress.yaml
# ├── po
# │   └── crudtest-5c85c5f77-58f45.yaml
# ├── pvc
# ├── rc
# ├── rs
# │   └── crudtest-5c85c5f77.yaml
# ├── sa
# │   ├── crudtest.yaml
# │   └── default.yaml
# ├── secrets
# │   ├── crudtest-token-9nckd.yaml
# │   └── default-token-8kvcg.yaml
# └── svc
#     └── crudtest.yaml
# 
# 13 directories, 25 files

# TODO Check

# TEST 4: TODO create a tmp file(e.g. .aaa.yaml.swp)
# TEST 5: TODO create a namespace yaml file
# TEST 6: TODO delete deployment with rm yaml
