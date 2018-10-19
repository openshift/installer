package machines

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"github.com/pkg/errors"
)

var userDataListTmpl = template.Must(template.New("user-data-list").Parse(`
kind: List
apiVersion: v1
metadata:
  resourceVersion: ""
  selfLink: ""
items:
{{- range $name, $content := . }}
- apiVersion: v1
  kind: Secret
  metadata:
    name: {{$name}}
    namespace: openshift-cluster-api
  type: Opaque
  data:
    userData: {{$content}}
{{- end}}
`))

func userDataList(data map[string][]byte) ([]byte, error) {
	encodedData := map[string]string{}
	for name, content := range data {
		encodedData[name] = base64.StdEncoding.EncodeToString(content)
	}
	buf := &bytes.Buffer{}
	if err := userDataListTmpl.Execute(buf, encodedData); err != nil {
		return nil, errors.Wrap(err, "failed to execute content.UserDataListTmpl")
	}
	return buf.Bytes(), nil
}
