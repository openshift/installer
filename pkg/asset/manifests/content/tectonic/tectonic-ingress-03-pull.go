package tectonic

import (
	"text/template"
)

var (
	// PullTectonicIngress  is the variable/constant representing the contents of the respective file
	PullTectonicIngress = template.Must(template.New("tectonic-ingress-03-pull.json").Parse(`
{
  "apiVersion": "v1",
  "kind": "Secret",
  "type": "kubernetes.io/dockerconfigjson",
  "metadata": {
    "namespace": "openshift-ingress",
    "name": "coreos-pull-secret"
  },
  "data": {
    ".dockerconfigjson": "{{.PullSecret}}"
  }
}
`))
)
