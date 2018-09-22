package bootkube

import (
	"text/template"
)

var (
	// MachineConfigServerTLSSecret is the constant to represent contents of machine_configservertlssecret.yaml file
	MachineConfigServerTLSSecret = template.Must(template.New("machine-config-server-tls-secret.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: machine-config-server-tls
  namespace: openshift-machine-config-operator
type: Opaque
data:
  tls.crt: {{.McsTLSCert}}
  tls.key: {{.McsTLSKey}}
`))
)
