package bootkube

import (
	"text/template"
)

var (
	// KubeCloudConfig is the constant to represent contents of kube_cloudconfig.yaml file
	KubeCloudConfig = template.Must(template.New("kube-cloud-config.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: kube-cloud-cfg
  namespace: kube-system
type: Opaque
data:
  config: ""
`))
)
