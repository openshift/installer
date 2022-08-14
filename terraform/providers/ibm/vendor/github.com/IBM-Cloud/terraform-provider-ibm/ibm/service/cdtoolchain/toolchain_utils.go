package cdtoolchain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetParametersForCreate(d *schema.ResourceData, resource *schema.Resource, remapFields map[string]string) map[string]interface{} {
	params := make(map[string]interface{})

	if _, ok := d.GetOk("parameters"); ok {
		srcParams := d.Get("parameters.0").(map[string]interface{})
		parametersSchema := resource.Schema["parameters"].Elem.(*schema.Resource).Schema
		for key, element := range parametersSchema {
			if !element.Computed && srcParams[key] != nil {
				params[getTargetField(key, remapFields)] = srcParams[key]
			}
		}
	}

	if _, ok := d.GetOk("initialization"); ok {
		srcInit := d.Get("initialization.0").(map[string]interface{})
		initSchema := resource.Schema["initialization"].Elem.(*schema.Resource).Schema
		for key, _ := range initSchema {
			if srcInit[key] != nil {
				params[getTargetField(key, remapFields)] = srcInit[key]
			}
		}
	}

	return params
}

func GetParametersForUpdate(d *schema.ResourceData, resource *schema.Resource, remapFields map[string]string) map[string]interface{} {
	params := make(map[string]interface{})
	srcParams := d.Get("parameters.0").(map[string]interface{})
	parametersSchema := resource.Schema["parameters"].Elem.(*schema.Resource).Schema
	for key, element := range parametersSchema {
		if !element.Computed && srcParams[key] != nil && d.HasChange("parameters.0."+key) {
			params[getTargetField(key, remapFields)] = srcParams[key]
		}
	}
	return params
}

func GetParametersFromRead(readParams map[string]interface{}, resource *schema.Resource, remapFields map[string]string) map[string]interface{} {
	params := make(map[string]interface{})
	parametersSchema := resource.Schema["parameters"].Elem.(*schema.Resource).Schema
	for key := range parametersSchema {
		readKey := getTargetField(key, remapFields)
		if readParams[readKey] != nil {
			params[key] = readParams[readKey]
		}
	}
	return params
}

func getTargetField(field string, remapFields map[string]string) string {
	if remapFields != nil {
		if val, ok := remapFields[field]; ok {
			return val
		}
	}
	return field
}
