package tectonic

import (
	"text/template"
)

var (
	// TLSTectonicIngress  is the variable/constant representing the contents of the respective file
	TLSTectonicIngress = template.Must(template.New("tectonic-ingress-02-tls.yaml").Parse(`
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
