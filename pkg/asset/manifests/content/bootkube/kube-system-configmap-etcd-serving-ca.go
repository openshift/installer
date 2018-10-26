package bootkube

import "text/template"

var (
	// KubeSystemConfigmapEtcdServingCA  is the constant to represent contents of kube-system-configmap-etcd-serving-ca.yaml file
	KubeSystemConfigmapEtcdServingCA = template.Must(template.New("kube-system-configmap-etcd-serving-ca.yaml").Parse(`
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: etcd-serving-ca
  namespace: kube-system
data:
  ca-bundle.crt: |
    {{.EtcdCaCert}}
`))
)
