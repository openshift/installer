package waf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func SizeConstraintSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"arn": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"size_constraints": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"comparison_operator": {
						Type:     schema.TypeString,
						Required: true,
					},
					"field_to_match": {
						Type:     schema.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"data": {
									Type:     schema.TypeString,
									Optional: true,
								},
								"type": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
					"size": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"text_transformation": {
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func DiffSizeConstraints(oldS, newS []interface{}) []*waf.SizeConstraintSetUpdate {
	updates := make([]*waf.SizeConstraintSetUpdate, 0)

	for _, os := range oldS {
		constraint := os.(map[string]interface{})

		if idx, contains := sliceContainsMap(newS, constraint); contains {
			newS = append(newS[:idx], newS[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.SizeConstraintSetUpdate{
			Action: aws.String(waf.ChangeActionDelete),
			SizeConstraint: &waf.SizeConstraint{
				FieldToMatch:       ExpandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: aws.String(constraint["comparison_operator"].(string)),
				Size:               aws.Int64(int64(constraint["size"].(int))),
				TextTransformation: aws.String(constraint["text_transformation"].(string)),
			},
		})
	}

	for _, ns := range newS {
		constraint := ns.(map[string]interface{})

		updates = append(updates, &waf.SizeConstraintSetUpdate{
			Action: aws.String(waf.ChangeActionInsert),
			SizeConstraint: &waf.SizeConstraint{
				FieldToMatch:       ExpandFieldToMatch(constraint["field_to_match"].([]interface{})[0].(map[string]interface{})),
				ComparisonOperator: aws.String(constraint["comparison_operator"].(string)),
				Size:               aws.Int64(int64(constraint["size"].(int))),
				TextTransformation: aws.String(constraint["text_transformation"].(string)),
			},
		})
	}
	return updates
}

func FlattenSizeConstraints(sc []*waf.SizeConstraint) []interface{} {
	out := make([]interface{}, len(sc))
	for i, c := range sc {
		m := make(map[string]interface{})
		m["comparison_operator"] = aws.StringValue(c.ComparisonOperator)
		if c.FieldToMatch != nil {
			m["field_to_match"] = FlattenFieldToMatch(c.FieldToMatch)
		}
		m["size"] = aws.Int64Value(c.Size)
		m["text_transformation"] = aws.StringValue(c.TextTransformation)
		out[i] = m
	}
	return out
}

func FlattenGeoMatchConstraint(ts []*waf.GeoMatchConstraint) []interface{} {
	out := make([]interface{}, len(ts))
	for i, t := range ts {
		m := make(map[string]interface{})
		m["type"] = aws.StringValue(t.Type)
		m["value"] = aws.StringValue(t.Value)
		out[i] = m
	}
	return out
}

func DiffGeoMatchSetConstraints(oldT, newT []interface{}) []*waf.GeoMatchSetUpdate {
	updates := make([]*waf.GeoMatchSetUpdate, 0)

	for _, od := range oldT {
		constraint := od.(map[string]interface{})

		if idx, contains := sliceContainsMap(newT, constraint); contains {
			newT = append(newT[:idx], newT[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.GeoMatchSetUpdate{
			Action: aws.String(waf.ChangeActionDelete),
			GeoMatchConstraint: &waf.GeoMatchConstraint{
				Type:  aws.String(constraint["type"].(string)),
				Value: aws.String(constraint["value"].(string)),
			},
		})
	}

	for _, nd := range newT {
		constraint := nd.(map[string]interface{})

		updates = append(updates, &waf.GeoMatchSetUpdate{
			Action: aws.String(waf.ChangeActionInsert),
			GeoMatchConstraint: &waf.GeoMatchConstraint{
				Type:  aws.String(constraint["type"].(string)),
				Value: aws.String(constraint["value"].(string)),
			},
		})
	}
	return updates
}

func DiffRegexPatternSetPatternStrings(oldPatterns, newPatterns []interface{}) []*waf.RegexPatternSetUpdate {
	updates := make([]*waf.RegexPatternSetUpdate, 0)

	for _, op := range oldPatterns {
		if idx, contains := verify.SliceContainsString(newPatterns, op.(string)); contains {
			newPatterns = append(newPatterns[:idx], newPatterns[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.RegexPatternSetUpdate{
			Action:             aws.String(waf.ChangeActionDelete),
			RegexPatternString: aws.String(op.(string)),
		})
	}

	for _, np := range newPatterns {
		updates = append(updates, &waf.RegexPatternSetUpdate{
			Action:             aws.String(waf.ChangeActionInsert),
			RegexPatternString: aws.String(np.(string)),
		})
	}
	return updates
}

func DiffRulePredicates(oldP, newP []interface{}) []*waf.RuleUpdate {
	updates := make([]*waf.RuleUpdate, 0)

	for _, op := range oldP {
		predicate := op.(map[string]interface{})

		if idx, contains := sliceContainsMap(newP, predicate); contains {
			newP = append(newP[:idx], newP[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.RuleUpdate{
			Action: aws.String(waf.ChangeActionDelete),
			Predicate: &waf.Predicate{
				Negated: aws.Bool(predicate["negated"].(bool)),
				Type:    aws.String(predicate["type"].(string)),
				DataId:  aws.String(predicate["data_id"].(string)),
			},
		})
	}

	for _, np := range newP {
		predicate := np.(map[string]interface{})

		updates = append(updates, &waf.RuleUpdate{
			Action: aws.String(waf.ChangeActionInsert),
			Predicate: &waf.Predicate{
				Negated: aws.Bool(predicate["negated"].(bool)),
				Type:    aws.String(predicate["type"].(string)),
				DataId:  aws.String(predicate["data_id"].(string)),
			},
		})
	}
	return updates
}

func DiffRuleGroupActivatedRules(oldRules, newRules []interface{}) []*waf.RuleGroupUpdate {
	updates := make([]*waf.RuleGroupUpdate, 0)

	for _, op := range oldRules {
		rule := op.(map[string]interface{})

		if idx, contains := sliceContainsMap(newRules, rule); contains {
			newRules = append(newRules[:idx], newRules[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.RuleGroupUpdate{
			Action:        aws.String(waf.ChangeActionDelete),
			ActivatedRule: ExpandActivatedRule(rule),
		})
	}

	for _, np := range newRules {
		rule := np.(map[string]interface{})

		updates = append(updates, &waf.RuleGroupUpdate{
			Action:        aws.String(waf.ChangeActionInsert),
			ActivatedRule: ExpandActivatedRule(rule),
		})
	}
	return updates
}

func FlattenActivatedRules(activatedRules []*waf.ActivatedRule) []interface{} {
	out := make([]interface{}, len(activatedRules))
	for i, ar := range activatedRules {
		rule := map[string]interface{}{
			"priority": aws.Int64Value(ar.Priority),
			"rule_id":  aws.StringValue(ar.RuleId),
			"type":     aws.StringValue(ar.Type),
		}
		if ar.Action != nil {
			rule["action"] = []interface{}{
				map[string]interface{}{
					"type": aws.StringValue(ar.Action.Type),
				},
			}
		}
		out[i] = rule
	}
	return out
}

func ExpandActivatedRule(rule map[string]interface{}) *waf.ActivatedRule {
	r := &waf.ActivatedRule{
		Priority: aws.Int64(int64(rule["priority"].(int))),
		RuleId:   aws.String(rule["rule_id"].(string)),
		Type:     aws.String(rule["type"].(string)),
	}

	if a, ok := rule["action"].([]interface{}); ok && len(a) > 0 {
		m := a[0].(map[string]interface{})
		r.Action = &waf.WafAction{
			Type: aws.String(m["type"].(string)),
		}
	}
	return r
}

func FlattenRegexMatchTuples(tuples []*waf.RegexMatchTuple) []interface{} {
	out := make([]interface{}, len(tuples))
	for i, t := range tuples {
		m := make(map[string]interface{})

		if t.FieldToMatch != nil {
			m["field_to_match"] = FlattenFieldToMatch(t.FieldToMatch)
		}
		m["regex_pattern_set_id"] = aws.StringValue(t.RegexPatternSetId)
		m["text_transformation"] = aws.StringValue(t.TextTransformation)

		out[i] = m
	}
	return out
}

func ExpandRegexMatchTuple(tuple map[string]interface{}) *waf.RegexMatchTuple {
	ftm := tuple["field_to_match"].([]interface{})
	return &waf.RegexMatchTuple{
		FieldToMatch:       ExpandFieldToMatch(ftm[0].(map[string]interface{})),
		RegexPatternSetId:  aws.String(tuple["regex_pattern_set_id"].(string)),
		TextTransformation: aws.String(tuple["text_transformation"].(string)),
	}
}

func DiffRegexMatchSetTuples(oldT, newT []interface{}) []*waf.RegexMatchSetUpdate {
	updates := make([]*waf.RegexMatchSetUpdate, 0)

	for _, ot := range oldT {
		tuple := ot.(map[string]interface{})

		if idx, contains := sliceContainsMap(newT, tuple); contains {
			newT = append(newT[:idx], newT[idx+1:]...)
			continue
		}

		updates = append(updates, &waf.RegexMatchSetUpdate{
			Action:          aws.String(waf.ChangeActionDelete),
			RegexMatchTuple: ExpandRegexMatchTuple(tuple),
		})
	}

	for _, nt := range newT {
		tuple := nt.(map[string]interface{})

		updates = append(updates, &waf.RegexMatchSetUpdate{
			Action:          aws.String(waf.ChangeActionInsert),
			RegexMatchTuple: ExpandRegexMatchTuple(tuple),
		})
	}
	return updates
}

func RegexMatchSetTupleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["field_to_match"]; ok {
		ftms := v.([]interface{})
		ftm := ftms[0].(map[string]interface{})

		if v, ok := ftm["data"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(v.(string))))
		}
		buf.WriteString(fmt.Sprintf("%s-", ftm["type"].(string)))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["regex_pattern_set_id"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["text_transformation"].(string)))

	return create.StringHashcode(buf.String())
}
