package ingress

const (
	// SvcAccount  is the variable/constant representing the contents of the respective file
	SvcAccount = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tectonic-ingress-controller-operator
  namespace: openshift-ingress
`
)
