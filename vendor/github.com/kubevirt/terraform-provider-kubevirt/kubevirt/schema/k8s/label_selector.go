package k8s

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/kubevirt/terraform-provider-kubevirt/kubevirt/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func labelSelectorFields(updatable bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"match_expressions": {
			Type:        schema.TypeList,
			Description: "A list of label selector requirements. The requirements are ANDed.",
			Optional:    true,
			ForceNew:    !updatable,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:        schema.TypeString,
						Description: "The label key that the selector applies to.",
						Optional:    true,
						ForceNew:    !updatable,
					},
					"operator": {
						Type:        schema.TypeString,
						Description: "A key's relationship to a set of values. Valid operators ard `In`, `NotIn`, `Exists` and `DoesNotExist`.",
						Optional:    true,
						ForceNew:    !updatable,
					},
					"values": {
						Type:        schema.TypeSet,
						Description: "An array of string values. If the operator is `In` or `NotIn`, the values array must be non-empty. If the operator is `Exists` or `DoesNotExist`, the values array must be empty. This array is replaced during a strategic merge patch.",
						Optional:    true,
						ForceNew:    !updatable,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Set:         schema.HashString,
					},
				},
			},
		},
		"match_labels": {
			Type:        schema.TypeMap,
			Description: "A map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of `match_expressions`, whose key field is \"key\", the operator is \"In\", and the values array contains only \"value\". The requirements are ANDed.",
			Optional:    true,
			ForceNew:    !updatable,
		},
	}
}

// Flatteners

func flattenLabelSelector(in *metav1.LabelSelector) []interface{} {
	att := make(map[string]interface{})
	if len(in.MatchLabels) > 0 {
		att["match_labels"] = utils.FlattenStringMap(in.MatchLabels)
	}
	if len(in.MatchExpressions) > 0 {
		att["match_expressions"] = flattenLabelSelectorRequirement(in.MatchExpressions)
	}
	return []interface{}{att}
}

func flattenLabelSelectorRequirement(in []metav1.LabelSelectorRequirement) []interface{} {
	att := make([]interface{}, len(in), len(in))
	for i, n := range in {
		m := make(map[string]interface{})
		m["key"] = n.Key
		m["operator"] = n.Operator
		m["values"] = utils.NewStringSet(schema.HashString, n.Values)
		att[i] = m
	}
	return att
}

// Expanders

func expandLabelSelector(l []interface{}) *metav1.LabelSelector {
	if len(l) == 0 || l[0] == nil {
		return &metav1.LabelSelector{}
	}
	in := l[0].(map[string]interface{})
	obj := &metav1.LabelSelector{}
	if v, ok := in["match_labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.MatchLabels = utils.ExpandStringMap(v)
	}
	if v, ok := in["match_expressions"].([]interface{}); ok && len(v) > 0 {
		obj.MatchExpressions = expandLabelSelectorRequirement(v)
	}
	return obj
}

func expandLabelSelectorRequirement(l []interface{}) []metav1.LabelSelectorRequirement {
	if len(l) == 0 || l[0] == nil {
		return []metav1.LabelSelectorRequirement{}
	}
	obj := make([]metav1.LabelSelectorRequirement, len(l), len(l))
	for i, n := range l {
		in := n.(map[string]interface{})
		obj[i] = metav1.LabelSelectorRequirement{
			Key:      in["key"].(string),
			Operator: metav1.LabelSelectorOperator(in["operator"].(string)),
			Values:   utils.SliceOfString(in["values"].(*schema.Set).List()),
		}
	}
	return obj
}
