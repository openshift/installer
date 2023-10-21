package wafregional

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfwaf "github.com/hashicorp/terraform-provider-aws/internal/service/waf"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_wafregional_web_acl", name="Web ACL")
// @Tags(identifierAttribute="arn")
func ResourceWebACL() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceWebACLCreate,
		ReadWithoutTimeout:   resourceWebACLRead,
		UpdateWithoutTimeout: resourceWebACLUpdate,
		DeleteWithoutTimeout: resourceWebACLDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_action": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								waf.WafActionTypeAllow,
								waf.WafActionTypeBlock,
								waf.WafActionTypeCount,
							}, false),
						},
					},
				},
			},
			"logging_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_destination": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
						"redacted_fields": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_to_match": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"data": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"type": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(wafregional.MatchFieldType_Values(), false),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											waf.WafActionTypeAllow,
											waf.WafActionTypeBlock,
											waf.WafActionTypeCount,
										}, false),
									},
								},
							},
						},
						"override_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											waf.WafOverrideActionTypeCount,
											waf.WafOverrideActionTypeNone,
										}, false),
									},
								},
							},
						},
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  waf.WafRuleTypeRegular,
							ValidateFunc: validation.StringInSlice([]string{
								waf.WafRuleTypeRegular,
								waf.WafRuleTypeRateBased,
								waf.WafRuleTypeGroup,
							}, false),
						},
						"rule_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceWebACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	wr := NewRetryer(conn, region)
	out, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
		input := &waf.CreateWebACLInput{
			ChangeToken:   token,
			DefaultAction: tfwaf.ExpandAction(d.Get("default_action").([]interface{})),
			MetricName:    aws.String(d.Get("metric_name").(string)),
			Name:          aws.String(d.Get("name").(string)),
			Tags:          GetTagsIn(ctx),
		}

		return conn.CreateWebACLWithContext(ctx, input)
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating WAF Regional Web ACL (%s): %s", d.Get("name").(string), err)
	}
	resp := out.(*waf.CreateWebACLOutput)
	d.SetId(aws.StringValue(resp.WebACL.WebACLId))

	// The WAF API currently omits this, but use it when it becomes available
	webACLARN := aws.StringValue(resp.WebACL.WebACLArn)
	if webACLARN == "" {
		webACLARN = arn.ARN{
			AccountID: meta.(*conns.AWSClient).AccountID,
			Partition: meta.(*conns.AWSClient).Partition,
			Region:    meta.(*conns.AWSClient).Region,
			Resource:  fmt.Sprintf("webacl/%s", d.Id()),
			Service:   "waf-regional",
		}.String()
	}

	loggingConfiguration := d.Get("logging_configuration").([]interface{})

	if len(loggingConfiguration) == 1 {
		input := &waf.PutLoggingConfigurationInput{
			LoggingConfiguration: expandLoggingConfiguration(loggingConfiguration, webACLARN),
		}

		log.Printf("[DEBUG] Updating WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), input)
		if _, err := conn.PutLoggingConfigurationWithContext(ctx, input); err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), err)
		}
	}

	rules := d.Get("rule").(*schema.Set).List()
	if len(rules) > 0 {
		wr := NewRetryer(conn, region)
		_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			req := &waf.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: tfwaf.ExpandAction(d.Get("default_action").([]interface{})),
				Updates:       diffWebACLRules([]interface{}{}, rules),
				WebACLId:      aws.String(d.Id()),
			}
			return conn.UpdateWebACLWithContext(ctx, req)
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceWebACLRead(ctx, d, meta)...)
}

func resourceWebACLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)

	params := &waf.GetWebACLInput{
		WebACLId: aws.String(d.Id()),
	}

	resp, err := conn.GetWebACLWithContext(ctx, params)
	if err != nil {
		if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, wafregional.ErrCodeWAFNonexistentItemException) {
			log.Printf("[WARN] WAF Regional ACL (%s) not found, removing from state", d.Id())
			d.SetId("")
			return diags
		}

		return sdkdiag.AppendErrorf(diags, "unable to read WAF Regional ACL (%s): %s", d.Id(), err)
	}

	if !d.IsNewResource() && (resp == nil || resp.WebACL == nil) {
		log.Printf("[WARN] WAF Regional ACL (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	// The WAF API currently omits this, but use it when it becomes available
	webACLARN := aws.StringValue(resp.WebACL.WebACLArn)
	if webACLARN == "" {
		webACLARN = arn.ARN{
			AccountID: meta.(*conns.AWSClient).AccountID,
			Partition: meta.(*conns.AWSClient).Partition,
			Region:    meta.(*conns.AWSClient).Region,
			Resource:  fmt.Sprintf("webacl/%s", d.Id()),
			Service:   "waf-regional",
		}.String()
	}
	d.Set("arn", webACLARN)

	if err := d.Set("default_action", tfwaf.FlattenAction(resp.WebACL.DefaultAction)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting default_action: %s", err)
	}
	d.Set("name", resp.WebACL.Name)
	d.Set("metric_name", resp.WebACL.MetricName)
	if err := d.Set("rule", tfwaf.FlattenWebACLRules(resp.WebACL.Rules)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting rule: %s", err)
	}

	getLoggingConfigurationInput := &waf.GetLoggingConfigurationInput{
		ResourceArn: aws.String(d.Get("arn").(string)),
	}
	loggingConfiguration := []interface{}{}

	log.Printf("[DEBUG] Getting WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), getLoggingConfigurationInput)
	getLoggingConfigurationOutput, err := conn.GetLoggingConfigurationWithContext(ctx, getLoggingConfigurationInput)

	if err != nil && !tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
		return sdkdiag.AppendErrorf(diags, "getting WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), err)
	}

	if getLoggingConfigurationOutput != nil {
		loggingConfiguration = flattenLoggingConfiguration(getLoggingConfigurationOutput.LoggingConfiguration)
	}

	if err := d.Set("logging_configuration", loggingConfiguration); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting logging_configuration: %s", err)
	}

	return diags
}

func resourceWebACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	if d.HasChanges("default_action", "rule") {
		o, n := d.GetChange("rule")
		oldR, newR := o.(*schema.Set).List(), n.(*schema.Set).List()

		wr := NewRetryer(conn, region)
		_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			req := &waf.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: tfwaf.ExpandAction(d.Get("default_action").([]interface{})),
				Updates:       diffWebACLRules(oldR, newR),
				WebACLId:      aws.String(d.Id()),
			}
			return conn.UpdateWebACLWithContext(ctx, req)
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("logging_configuration") {
		loggingConfiguration := d.Get("logging_configuration").([]interface{})

		if len(loggingConfiguration) == 1 {
			input := &waf.PutLoggingConfigurationInput{
				LoggingConfiguration: expandLoggingConfiguration(loggingConfiguration, d.Get("arn").(string)),
			}

			log.Printf("[DEBUG] Updating WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), input)
			if _, err := conn.PutLoggingConfigurationWithContext(ctx, input); err != nil {
				return sdkdiag.AppendErrorf(diags, "updating WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), err)
			}
		} else {
			input := &waf.DeleteLoggingConfigurationInput{
				ResourceArn: aws.String(d.Get("arn").(string)),
			}

			log.Printf("[DEBUG] Deleting WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), input)
			if _, err := conn.DeleteLoggingConfigurationWithContext(ctx, input); err != nil {
				return sdkdiag.AppendErrorf(diags, "deleting WAF Regional Web ACL (%s) Logging Configuration: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceWebACLRead(ctx, d, meta)...)
}

func resourceWebACLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).WAFRegionalConn(ctx)
	region := meta.(*conns.AWSClient).Region

	// First, need to delete all rules
	rules := d.Get("rule").(*schema.Set).List()
	if len(rules) > 0 {
		wr := NewRetryer(conn, region)
		_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			req := &waf.UpdateWebACLInput{
				ChangeToken:   token,
				DefaultAction: tfwaf.ExpandAction(d.Get("default_action").([]interface{})),
				Updates:       diffWebACLRules(rules, []interface{}{}),
				WebACLId:      aws.String(d.Id()),
			}
			return conn.UpdateWebACLWithContext(ctx, req)
		})
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Removing WAF Regional ACL Rules: %s", err)
		}
	}

	wr := NewRetryer(conn, region)
	_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
		req := &waf.DeleteWebACLInput{
			ChangeToken: token,
			WebACLId:    aws.String(d.Id()),
		}

		log.Printf("[INFO] Deleting WAF ACL")
		return conn.DeleteWebACLWithContext(ctx, req)
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "Deleting WAF Regional ACL: %s", err)
	}
	return diags
}

func expandLoggingConfiguration(l []interface{}, resourceARN string) *waf.LoggingConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	loggingConfiguration := &waf.LoggingConfiguration{
		LogDestinationConfigs: []*string{
			aws.String(m["log_destination"].(string)),
		},
		RedactedFields: expandRedactedFields(m["redacted_fields"].([]interface{})),
		ResourceArn:    aws.String(resourceARN),
	}

	return loggingConfiguration
}

func expandRedactedFields(l []interface{}) []*waf.FieldToMatch {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})

	if m["field_to_match"] == nil {
		return nil
	}

	redactedFields := make([]*waf.FieldToMatch, 0)

	for _, fieldToMatch := range m["field_to_match"].(*schema.Set).List() {
		if fieldToMatch == nil {
			continue
		}

		redactedFields = append(redactedFields, tfwaf.ExpandFieldToMatch(fieldToMatch.(map[string]interface{})))
	}

	return redactedFields
}

func flattenLoggingConfiguration(loggingConfiguration *waf.LoggingConfiguration) []interface{} {
	if loggingConfiguration == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"log_destination": "",
		"redacted_fields": flattenRedactedFields(loggingConfiguration.RedactedFields),
	}

	if len(loggingConfiguration.LogDestinationConfigs) > 0 {
		m["log_destination"] = aws.StringValue(loggingConfiguration.LogDestinationConfigs[0])
	}

	return []interface{}{m}
}

func flattenRedactedFields(fieldToMatches []*waf.FieldToMatch) []interface{} {
	if len(fieldToMatches) == 0 {
		return []interface{}{}
	}

	fieldToMatchResource := &schema.Resource{
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
	}
	l := make([]interface{}, len(fieldToMatches))

	for i, fieldToMatch := range fieldToMatches {
		l[i] = tfwaf.FlattenFieldToMatch(fieldToMatch)[0]
	}

	m := map[string]interface{}{
		"field_to_match": schema.NewSet(schema.HashResource(fieldToMatchResource), l),
	}

	return []interface{}{m}
}

func diffWebACLRules(oldR, newR []interface{}) []*waf.WebACLUpdate {
	updates := make([]*waf.WebACLUpdate, 0)

	for _, or := range oldR {
		aclRule := or.(map[string]interface{})

		if idx, contains := sliceContainsMap(newR, aclRule); contains {
			newR = append(newR[:idx], newR[idx+1:]...)
			continue
		}
		updates = append(updates, tfwaf.ExpandWebACLUpdate(waf.ChangeActionDelete, aclRule))
	}

	for _, nr := range newR {
		aclRule := nr.(map[string]interface{})
		updates = append(updates, tfwaf.ExpandWebACLUpdate(waf.ChangeActionInsert, aclRule))
	}
	return updates
}
