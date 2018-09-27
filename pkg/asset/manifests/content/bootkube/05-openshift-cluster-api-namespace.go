package bootkube

const (
	// OpenshiftClusterAPINamespace is the constant to represent contents of Openshift_ClusterApiNamespace.yaml file
	OpenshiftClusterAPINamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
  # This is the namespace used to hold cluster-api components.
  name: openshift-cluster-api
  labels:
    name: openshift-cluster-api
    openshift.io/run-level: "1"
`
)
