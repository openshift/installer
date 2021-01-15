package machines

var etcdSingleMasterData = []byte(`apiVersion: operator.openshift.io/v1
kind: Etcd
metadata:
  name: cluster
spec:
  managementState: Managed
  unsupportedConfigOverrides:
    useUnsupportedUnsafeNonHANonProductionUnstableEtcd: true
`)

var ingressSingleMasterData = []byte(`apiVersion: operator.openshift.io/v1
kind: IngressController
metadata:
  name: default
  namespace: openshift-ingress-operator
spec:
  replicas: 1
`)

var authSingleMasterData = []byte(`apiVersion: operator.openshift.io/v1
kind: Authentication
metadata:
  name: cluster
spec:
  unsupportedConfigOverrides:
    useUnsupportedUnsafeNonHANonProductionUnstableOAuthServer: true
`)
