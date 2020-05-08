package explain

import (
	"fmt"
	"io"
	"sort"
	"strings"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	fieldIndent     = 4
	fieldDescIndent = 6
)

type printer struct {
	Writer io.Writer
}

func (p printer) PrintKindAndVersion() {
	io.WriteString(p.Writer, fmt.Sprintf("KIND:     %s\n", "InstallConfig"))
	io.WriteString(p.Writer, fmt.Sprintf("VERSION:  %s\n\n", "v1"))
}

func (p printer) PrintResource(schema *apiextv1.JSONSchemaProps) {
	resource := schema.Type
	if schema.Items != nil && schema.Items.Schema != nil {
		resource = fmt.Sprintf("[]%s", schema.Items.Schema.Type)
	}
	io.WriteString(p.Writer, fmt.Sprintf("RESOURCE: <%s>\n", resource))

	desc := schema.Description
	if len(desc) == 0 {
		desc = "<empty>"
	}
	write(2, p.Writer, desc)
	io.WriteString(p.Writer, "\n")

}

func (p printer) PrintFields(schema *apiextv1.JSONSchemaProps) {
	required := sets.NewString(schema.Required...)
	properties := map[string]apiextv1.JSONSchemaProps{}
	if schema.Items != nil && schema.Items.Schema != nil && len(schema.Items.Schema.Properties) > 0 {
		properties = schema.Items.Schema.Properties
		required.Insert(schema.Items.Schema.Required...)
	}
	if len(schema.Properties) > 0 {
		properties = schema.Properties
	}
	if len(properties) == 0 {
		return
	}

	var keys []string
	for pname := range properties {
		keys = append(keys, pname)
	}
	sort.Strings(keys)

	io.WriteString(p.Writer, "FIELDS:\n")
	for _, pname := range keys {
		pschema := properties[pname]
		p.printField(pname, required.Has(pname), &pschema)
	}
}

func (p printer) printField(name string, required bool, schema *apiextv1.JSONSchemaProps) {
	ftype := schema.Type
	if schema.Items != nil && schema.Items.Schema != nil {
		ftype = fmt.Sprintf("[]%s", schema.Items.Schema.Type)
	}
	title := fmt.Sprintf("%s <%s>", name, ftype)
	if required {
		title = fmt.Sprintf("%s -required-", title)
	}
	write(fieldIndent, p.Writer, title)

	if schema.Default != nil {
		write(fieldDescIndent, p.Writer, fmt.Sprintf("Default: %s", defaultString(*schema.Default)))
	}
	if len(schema.Format) > 0 {
		write(fieldDescIndent, p.Writer, fmt.Sprintf("Format: %s", schema.Format))
	}
	if len(schema.Enum) > 0 {
		write(fieldDescIndent, p.Writer, fmt.Sprintf("Valid Values: %s", strings.Join(validValues(schema.Enum), ",")))
	}

	fdesc := schema.Description
	if fdesc == "" {
		fdesc = "<empty>"
	}
	write(6, p.Writer, fdesc)
	if schema.Items != nil && schema.Items.Schema != nil && len(schema.Items.Schema.Description) > 0 {
		write(6, p.Writer, schema.Items.Schema.Description)
	}
	io.WriteString(p.Writer, "\n")
}

func write(indentLevel int, w io.Writer, s string) {
	if strings.TrimSpace(s) == "" {
		io.WriteString(w, "\n")
	}
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent = indent + " "
	}
	io.WriteString(w, indent+s+"\n")
}

func defaultString(obj apiextv1.JSON) string { return string(obj.Raw) }

func validValues(objs []apiextv1.JSON) []string {
	ret := make([]string, len(objs))
	for idx, obj := range objs {
		ret[idx] = string(obj.Raw)
	}
	return ret
}
