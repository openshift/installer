package secrets

import (
	"text/template"
)

var (
	// CaCert  is the variable/constant representing the contents of the respective file
	CaCert = template.Must(template.New("ca-cert.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: tectonic-ca-cert-secret
  namespace: tectonic-system
type: Opaque
data:
  ca-cert: {{.IngressCaCert}}
`))
)
