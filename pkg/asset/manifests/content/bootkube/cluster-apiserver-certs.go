package bootkube

import (
	"text/template"
)

var (
	// ClusterApiserverCerts is the constant to represent contents of cluster_apiservercerts.yaml file
	ClusterApiserverCerts = template.Must(template.New("cluster-apiserver-certs.yaml").Parse(`
apiVersion: v1
kind: Secret
type: kubernetes.io/tls
metadata:
  name: cluster-apiserver-certs
  namespace: openshift-cluster-api
  labels:
    api: clusterapi
    apiserver: "true"
data:
  tls.crt: {{.ClusterapiCaCert}}
  tls.key: {{.ClusterapiCaKey}}
`))
)
