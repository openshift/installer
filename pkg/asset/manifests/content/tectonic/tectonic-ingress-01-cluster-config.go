package tectonic

import (
	"text/template"
)

var (
	// ClusterConfigTectonicIngress  is the variable/constant representing the contents of the respective file
	ClusterConfigTectonicIngress = template.Must(template.New("tectonic-ingress-01-cluster-config.yaml").Parse(`
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-config-v1
  namespace: openshift-ingress
data:
  ingress-config: |
    apiVersion: v1
    kind: TectonicIngressOperatorConfig
    type: {{.IngressKind}}
    statsPassword: {{.IngressStatusPassword}}
    statsUsername: admin
`))
)
