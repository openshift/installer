package bootkube

const (
	// OperatorstatusCrd is the constant to represent contents of Operatorstatus_Crd.yaml file
	OperatorstatusCrd = `
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: operatorstatuses.clusterversion.openshift.io
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: clusterversion.openshift.io
  # list of versions supported by this CustomResourceDefinition
  versions:
    - name: v1
      # Each version can be enabled/disabled by Served flag.
      served: true
      # One and only one version must be marked as the storage version.
      storage: true
  # either Namespaced or Cluster
  scope: Cluster
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: operatorstatuses
    # singular name to be used as an alias on the CLI and for display
    singular: operatorstatus
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: OperatorStatus
`
)
