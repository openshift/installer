package bootkube

import (
	"text/template"
)

var (
	// Pull is the constant to represent contents of pull.yaml file
	Pull = template.Must(template.New("pull.json").Parse(`
{
  "apiVersion": "v1",
  "kind": "Secret",
  "type": "kubernetes.io/dockerconfigjson",
  "metadata": {
    "namespace": "kube-system",
    "name": "coreos-pull-secret"
  },
  "data": {
    ".dockerconfigjson": "{{.PullSecret}}"
  }
}
`))
)
