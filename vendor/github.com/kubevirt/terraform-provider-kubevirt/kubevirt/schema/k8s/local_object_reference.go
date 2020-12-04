package k8s

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	v1 "k8s.io/api/core/v1"
)

func localObjectReferenceFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Name of the referent.",
			Required:    true,
		},
	}
}

func LocalObjectReferenceSchema(description string) *schema.Schema {
	fields := localObjectReferenceFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: description,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandLocalObjectReferences(localObjectReference []interface{}) *v1.LocalObjectReference {
	if len(localObjectReference) == 0 {
		return nil
	}
	result := &v1.LocalObjectReference{}

	in := localObjectReference[0].(map[string]interface{})

	if v, ok := in["name"].(string); ok {
		result.Name = v
	}

	return result
}

func FlattenLocalObjectReferences(localObjectReference v1.LocalObjectReference) []interface{} {
	att := make(map[string]interface{})

	att["name"] = localObjectReference.Name

	return []interface{}{att}
}
