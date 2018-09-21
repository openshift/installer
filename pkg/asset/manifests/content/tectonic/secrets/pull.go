package secrets

import (
	"text/template"
)

var (
	// Pull  is the variable/constant representing the contents of the respective file
	Pull = template.Must(template.New("pull.json").Parse(`
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
