package assets

import (
	"bytes"
	"text/template"
)

// MustTemplateAsset parses a named binasset as a Template and panics if the
// asset is not found or parses with a non-nil error.
func MustTemplateAsset(binasset string) *template.Template {
	return template.Must(template.New("").Parse(string(MustAsset(binasset))))
}

// RenderTemplate executes the give template with data.
func RenderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}
