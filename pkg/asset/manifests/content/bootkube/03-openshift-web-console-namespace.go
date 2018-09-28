package bootkube

const (
	// OpenshiftWebConsoleNamespace is the constant to represent contents of Openshift_WebConsoleNamespace.yaml file
	OpenshiftWebConsoleNamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
  # This is the namespace used to hold the openshift console.
  # They require openshift console run in this namespace.
  name: openshift-web-console
  labels:
    name: openshift-web-console
`
)
