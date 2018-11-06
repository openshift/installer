package bootkube

import "text/template"

var (
	// KubeSystemConfigmapRootCA  is the constant to represent contents of kube-system-configmap-root-ca.yaml file
	KubeSystemConfigmapRootCA = template.Must(template.New("kube-system-configmap-root-ca.yaml").Parse(`
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: root-ca
  namespace: kube-system
data:
  ca.crt: {{.RootCaCert}}
`))
)
