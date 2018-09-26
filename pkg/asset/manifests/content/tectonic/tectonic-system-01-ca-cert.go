package tectonic

import (
	"text/template"
)

var (
	// CaCertTectonicSystem  is the variable/constant representing the contents of the respective file
	CaCertTectonicSystem = template.Must(template.New("tectonic-system-01-ca-cert.yaml").Parse(`
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
