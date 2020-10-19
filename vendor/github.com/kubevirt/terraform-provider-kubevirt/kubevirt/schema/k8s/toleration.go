package k8s

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	v1 "k8s.io/api/core/v1"
)

func tolerationFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"effect": {
			Type:         schema.TypeString,
			Description:  "Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"NoSchedule", "PreferNoSchedule", "NoExecute"}, false),
		},
		"key": {
			Type:        schema.TypeString,
			Description: "Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.",
			Optional:    true,
		},
		"operator": {
			Type:         schema.TypeString,
			Description:  "Operator represents a key's relationship to the value. Valid operators are Exists and Equal. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.",
			Default:      "Equal",
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"Exists", "Equal"}, false),
		},
		"toleration_seconds": {
			// Use TypeString to allow an "unspecified" value,
			Type:         schema.TypeString,
			Description:  "TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.",
			Optional:     true,
			ValidateFunc: utils.ValidateTypeStringNullableInt,
		},
		"value": {
			Type:        schema.TypeString,
			Description: "Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.",
			Optional:    true,
		},
	}
}

func TolerationSchema() *schema.Schema {
	fields := tolerationFields()

	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "If specified, the pod's toleration. Optional: Defaults to empty",
		Elem: &schema.Resource{
			Schema: fields,
		},
	}

}

func ExpandTolerations(tolerations []interface{}) ([]v1.Toleration, error) {
	if len(tolerations) == 0 {
		return []v1.Toleration{}, nil
	}
	ts := make([]v1.Toleration, len(tolerations))
	for i, t := range tolerations {
		m := t.(map[string]interface{})
		ts[i] = v1.Toleration{}

		if value, ok := m["effect"].(string); ok {
			ts[i].Effect = v1.TaintEffect(value)
		}
		if value, ok := m["key"].(string); ok {
			ts[i].Key = value
		}
		if value, ok := m["operator"].(string); ok {
			ts[i].Operator = v1.TolerationOperator(value)
		}
		if value, ok := m["toleration_seconds"].(string); ok && value != "" {
			seconds, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid toleration_seconds must be int or \"\", got \"%s\"", value)
			}
			ts[i].TolerationSeconds = utils.PtrToInt64(seconds)
		}
		if value, ok := m["value"]; ok {
			ts[i].Value = value.(string)
		}
	}
	return ts, nil
}

func FlattenTolerations(tolerations []v1.Toleration) []interface{} {
	att := []interface{}{}
	for _, v := range tolerations {
		// The API Server may automatically add several Tolerations to pods, strip these to avoid TF diff.
		if strings.Contains(v.Key, "node.kubernetes.io/") {
			log.Printf("[INFO] ignoring toleration with key: %s", v.Key)
			continue
		}
		obj := map[string]interface{}{}

		if v.Effect != "" {
			obj["effect"] = string(v.Effect)
		}
		if v.Key != "" {
			obj["key"] = v.Key
		}
		if v.Operator != "" {
			obj["operator"] = string(v.Operator)
		}
		if v.TolerationSeconds != nil {
			obj["toleration_seconds"] = strconv.FormatInt(*v.TolerationSeconds, 10)
		}
		if v.Value != "" {
			obj["value"] = v.Value
		}
		att = append(att, obj)
	}
	return att
}
