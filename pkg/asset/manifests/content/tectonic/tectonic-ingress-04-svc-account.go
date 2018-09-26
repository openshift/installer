package tectonic

const (
	// SvcAccountTectonicIngress  is the variable/constant representing the contents of the respective file
	SvcAccountTectonicIngress = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tectonic-ingress-controller-operator
  namespace: openshift-ingress
`
)
