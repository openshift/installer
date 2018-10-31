package bootkube

import (
	"text/template"
)

var (
	// KubeSystemSecretEtcdClient is the constant to represent contents of kube-system-secret-etcd-client.yaml file
	KubeSystemSecretEtcdClient = template.Must(template.New("kube-system-secret-etcd-client.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: etcd-client
  namespace: kube-system
type: SecretTypeTLS
data:
  tls.crt: {{ .EtcdClientCert }}
  tls.key: {{ .EtcdClientKey }}
`))
)
