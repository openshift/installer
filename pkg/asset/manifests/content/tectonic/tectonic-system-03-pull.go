package tectonic

import (
	"text/template"
)

var (
	// PullTectonicSystem  is the variable/constant representing the contents of the respective file
	PullTectonicSystem = template.Must(template.New("tectonic-system-03-pull.json").Parse(`
{
  "apiVersion": "v1",
  "kind": "Secret",
  "type": "kubernetes.io/dockerconfigjson",
  "metadata": {
    "namespace": "tectonic-system",
    "name": "coreos-pull-secret"
  },
  "data": {
    ".dockerconfigjson": "{{.PullSecret}}"
  }
}
`))
)
