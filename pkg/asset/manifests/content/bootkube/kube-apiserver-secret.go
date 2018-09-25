package bootkube

import (
	"text/template"
)

var (
	// KubeApiserverSecret is the constant to represent contents of kube_apiserversecret.yaml file
	KubeApiserverSecret = template.Must(template.New("kube-apiserver-secret.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: kube-apiserver
  namespace: kube-system
type: Opaque
data:
  aggregator-ca.crt: {{.AggregatorCaCert}}
  aggregator-ca.key: {{.AggregatorCaKey}}
  apiserver.key: {{.ApiserverKey}}
  apiserver.crt: {{.ApiserverCert}}
  apiserver-proxy.key: {{.ApiserverProxyKey}}
  apiserver-proxy.crt: {{.ApiserverProxyCert}}
  service-account.pub: {{.ServiceaccountPub}}
  service-account.key: {{.ServiceaccountKey}}
  root-ca.crt: {{.RootCaCert}}
  kube-ca.crt: {{.KubeCaCert}}
  etcd-client-ca.crt: {{.EtcdCaCert}}
  etcd-client.crt: {{.EtcdClientCert}}
  etcd-client.key: {{.EtcdClientKey}}
  oidc-ca.crt: {{.OidcCaCert}}
  service-serving-ca.crt: {{.ServiceServingCaCert}}
  service-serving-ca.key: {{.ServiceServingCaKey}}
  kubeconfig: {{.OpenshiftLoopbackKubeconfig}}
`))
)
