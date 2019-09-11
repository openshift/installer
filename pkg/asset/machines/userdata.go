package machines

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"github.com/pkg/errors"
)

var userDataTmpl = template.Must(template.New("user-data").Parse(`apiVersion: v1
kind: Secret
metadata:
  name: {{.name}}
  namespace: openshift-machine-api
type: Opaque
data:
  disableTemplating: "dHJ1ZQo="
  userData: {{.content}}
`))

func userDataSecret(name string, content []byte) ([]byte, error) {
	encodedData := map[string]string{
		"name":    name,
		"content": base64.StdEncoding.EncodeToString(content),
	}
	buf := &bytes.Buffer{}
	if err := userDataTmpl.Execute(buf, encodedData); err != nil {
		return nil, errors.Wrap(err, "failed to execute user-data template")
	}
	return buf.Bytes(), nil
}
