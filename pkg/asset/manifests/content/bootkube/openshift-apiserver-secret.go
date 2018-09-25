package bootkube

import (
	"text/template"
)

var (
	// OpenshiftApiserverSecret is the constant to represent contents of openshift_apiserversecret.yaml file
	OpenshiftApiserverSecret = template.Must(template.New("openshift-apiserver-secret.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: openshift-apiserver
  namespace: kube-system
type: Opaque
data:
  aggregator-ca.crt: {{.AggregatorCaCert}}
  aggregator-ca.key: {{.AggregatorCaKey}}
  apiserver.key: {{.ApiserverKey}}
  apiserver.crt: {{.ApiserverCert}}
  openshift-apiserver.key: {{.OpenshiftApiserverKey}}
  openshift-apiserver.crt: {{.OpenshiftApiserverCert}}
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
