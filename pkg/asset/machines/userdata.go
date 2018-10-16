package machines

import (
	"bytes"
	"encoding/base64"
	"text/template"

	"github.com/pkg/errors"
)

var userDataTmpl = template.Must(template.New("user-data").Parse(`
apiVersion: v1
kind: Secret
metadata:
  name: {{.Name}}
  namespace: openshift-cluster-api
type: Opaque
data:
  userData: {{.UserDataContent}}
`))

func userData(secretName string, content []byte) ([]byte, error) {
	templateData := struct {
		Name            string
		UserDataContent string
	}{
		Name:            secretName,
		UserDataContent: base64.StdEncoding.EncodeToString(content),
	}
	buf := &bytes.Buffer{}
	if err := userDataTmpl.Execute(buf, templateData); err != nil {
		return nil, errors.Wrap(err, "failed to execute content.UserDataTmpl")
	}
	return buf.Bytes(), nil
}
