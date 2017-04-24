package server

import (
	"bytes"
	"text/template"

	"github.com/coreos/tectonic-installer/installer/binassets"
)

const (
	bcryptCost = 12
)

// mustTemplateAsset parses a named binasset as a Template and panics if the
// asset is not found or parses with a non-nil error.
func mustTemplateAsset(binasset string) *template.Template {
	return template.Must(template.New("").Parse(string(binassets.MustAsset(binasset))))
}

// renderTemplate executes the give template with data.
func renderTemplate(tmpl *template.Template, data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}
