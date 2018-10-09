package bootkube

const (
	// OpenshiftServiceCertSignerNamespace is the constant to represent the contents of 09-openshift-service-signer-namespace.yaml
	OpenshiftServiceCertSignerNamespace = `
---
apiVersion: v1
kind: Namespace
metadata:
  # This is the namespace used to hold the service-serving-cert-signer.
  name: openshift-service-cert-signer
  labels:
    openshift.io/run-level: "1"
`
)
