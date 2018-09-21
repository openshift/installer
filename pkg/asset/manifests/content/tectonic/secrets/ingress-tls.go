package secrets

import (
	"text/template"
)

var (
	// IngressTLS  is the variable/constant representing the contents of the respective file
	IngressTLS = template.Must(template.New("ingress-tls.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: tectonic-ingress-tls
  namespace: openshift-ingress
type: Opaque
data:
  tls.crt: {{.IngressTLSCert}}
  tls.key: {{.IngressTLSKey}}
  bundle.crt: {{.IngressTLSBundle}}
`))
)
