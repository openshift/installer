package ingress

import (
	"text/template"
)

var (
	// ClusterConfig  is the variable/constant representing the contents of the respective file
	ClusterConfig = template.Must(template.New("cluster-config.yaml").Parse(`
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
