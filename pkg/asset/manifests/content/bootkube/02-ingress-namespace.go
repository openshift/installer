package bootkube

const (
	// IngressNamespace is the constant to represent contents of Ingress_Namespace.yaml file
	IngressNamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
  # This is the namespace used to hold the tectonic ingress controllers
  name: openshift-ingress
  # Give the namespace a label, so we can select for it in networkpolicy
  labels:
    kubernetes.io/ingress.class: tectonic
    name: openshift-ingress
    openshift.io/run-level: "1"
`
)
