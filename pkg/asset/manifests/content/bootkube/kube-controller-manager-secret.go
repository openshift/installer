package bootkube

import (
	"text/template"
)

var (
	// KubeControllerManagerSecret is the constant to represent contents of kube_controllermanagersecret.yaml file
	KubeControllerManagerSecret = template.Must(template.New("kube-controller-manager-secret.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: kube-controller-manager
  namespace: kube-system
type: Opaque
data:
  service-account.key: {{.ServiceaccountKey}}
  root-ca.crt: {{.RootCaCert}}
  kube-ca.crt: {{.KubeCaCert}}
  kube-ca.key: {{.KubeCaKey}}
`))
)
