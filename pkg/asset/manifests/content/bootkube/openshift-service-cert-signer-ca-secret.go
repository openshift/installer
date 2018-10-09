package bootkube

import (
	"text/template"
)

var (
	// OpenshiftServiceCertSignerSecret is the constant to represent the contents of openshift-service-signer-secret.yaml
	OpenshiftServiceCertSignerSecret = template.Must(template.New("openshift-service-signer-secret.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: service-serving-cert-signer-signing-key
  namespace: openshift-service-cert-signer
type: kubernetes.io/tls
data:
  tls.crt: {{.ServiceServingCaCert}}
  tls.key: {{.ServiceServingCaKey}}
`))
)
