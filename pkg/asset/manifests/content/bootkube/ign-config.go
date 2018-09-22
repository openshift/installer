package bootkube

import (
	"text/template"
)

var (
	// IgnConfig is the constant to represent contents of ign_config.yaml file
	IgnConfig = template.Must(template.New("ign-config.yaml").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: ignition-worker
  namespace: openshift-cluster-api
type: Opaque
data:
  userData: {{.WorkerIgnConfig}}
`))
)
