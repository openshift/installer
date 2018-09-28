package bootkube

const (
	// TectonicNamespace is the constant to represent contents of Tectonic_Namespace.yaml file
	TectonicNamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
  name: tectonic-system  # Create the namespace first.
  labels:  # network policy can only select by labels
    name: tectonic-system
    openshift.io/run-level: "1"
`
)
