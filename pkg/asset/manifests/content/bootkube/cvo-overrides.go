package bootkube

import (
	"text/template"
)

var (
	// CVOOverrides is the constant to represent contents of cvo-override.yaml file
	// This is a gate to prevent CVO from installing these operators which is conflicting
	// with already owned resources by tectonic-operators.
	// This files can be dropped when the overrides list becomes empty.
	CVOOverrides = template.Must(template.New("cvo-override.yaml").Parse(`
apiVersion: clusterversion.openshift.io/v1
kind: CVOConfig
metadata:
  namespace: openshift-cluster-version
  name: cluster-version-operator
upstream: http://localhost:8080/graph
channel: fast
clusterID: {{.CVOClusterID}}
overrides:
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-core-operators
  name: openshift-cluster-kube-apiserver-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-core-operators
  name: openshift-cluster-kube-scheduler-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-core-operators
  name: openshift-cluster-kube-controller-manager-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-core-operators
  name: openshift-cluster-openshift-apiserver-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-core-operators
  name: openshift-cluster-openshift-controller-manager-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-cluster-network-operator
  name: cluster-network-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with tectonic-ingress-controller-operator
  namespace: openshift-cluster-ingress-operator
  name: cluster-ingress-operator
  unmanaged: true
- kind: ServiceAccount                # missing run level 0 on the namespace and has 0000_08
  namespace: openshift-cluster-dns-operator
  name: cluster-dns-operator
  unmanaged: true
- kind: Deployment                    # this conflicts with kube-core-operator
  namespace: openshift-cluster-dns-operator
  name: cluster-dns-operator
  unmanaged: true
`))
)
