package elbv2

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

const (
	listenerRulePriorityMin     = 1
	listenerRulePriorityMax     = 50_000
	listenerRulePriorityDefault = 99_999

	listenerActionOrderMin = 1
	listenerActionOrderMax = 50_000
)

func ResourceListenerRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerRuleCreate,
		Read:   resourceListenerRuleRead,
		Update: resourceListenerRuleUpdate,
		Delete: resourceListenerRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     false,
				ValidateFunc: validListenerRulePriority,
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(elbv2.ActionTypeEnum_Values(), true),
						},
						"order": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(listenerActionOrderMin, listenerActionOrderMax),
						},

						"target_group_arn": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumForward),
							ValidateFunc:     verify.ValidARN,
						},

						"forward": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumForward),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_group": {
										Type:     schema.TypeSet,
										MinItems: 2,
										MaxItems: 5,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"arn": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: verify.ValidARN,
												},
												"weight": {
													Type:         schema.TypeInt,
													ValidateFunc: validation.IntBetween(0, 999),
													Default:      1,
													Optional:     true,
												},
											},
										},
									},
									"stickiness": {
										Type:             schema.TypeList,
										Optional:         true,
										DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
										MaxItems:         1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
												"duration": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 604800),
												},
											},
										},
									},
								},
							},
						},

						"redirect": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumRedirect),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "#{host}",
										ValidateFunc: validation.StringLenBetween(1, 128),
									},

									"path": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "/#{path}",
										ValidateFunc: validation.StringLenBetween(1, 128),
									},

									"port": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "#{port}",
									},

									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "#{protocol}",
										ValidateFunc: validation.StringInSlice([]string{
											"#{protocol}",
											"HTTP",
											"HTTPS",
										}, false),
									},

									"query": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "#{query}",
										ValidateFunc: validation.StringLenBetween(0, 128),
									},

									"status_code": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(elbv2.RedirectActionStatusCodeEnum_Values(), false),
									},
								},
							},
						},

						"fixed_response": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumFixedResponse),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"text/plain",
											"text/css",
											"text/html",
											"application/javascript",
											"application/json",
										}, false),
									},

									"message_body": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(0, 1024),
									},

									"status_code": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[245]\d\d$`), ""),
									},
								},
							},
						},

						"authenticate_cognito": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumAuthenticateCognito),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_request_extra_params": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"on_unauthenticated_request": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(elbv2.AuthenticateCognitoActionConditionalBehaviorEnum_Values(), true),
									},
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "openid",
									},
									"session_cookie_name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "AWSELBAuthSessionCookie",
									},
									"session_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  604800,
									},
									"user_pool_arn": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: verify.ValidARN,
									},
									"user_pool_client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"user_pool_domain": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"authenticate_oidc": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfActionTypeNot(elbv2.ActionTypeEnumAuthenticateOidc),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_request_extra_params": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"authorization_endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"client_secret": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
									},
									"issuer": {
										Type:     schema.TypeString,
										Required: true,
									},
									"on_unauthenticated_request": {
										Type:         schema.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.StringInSlice(elbv2.AuthenticateOidcActionConditionalBehaviorEnum_Values(), true),
									},
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "openid",
									},
									"session_cookie_name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "AWSELBAuthSessionCookie",
									},
									"session_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  604800,
									},
									"token_endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
									"user_info_endpoint": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"condition": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_header": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 128),
										},
										Set: schema.HashString,
									},
								},
							},
						},
						"http_header": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_header_name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9!#$%&'*+-.^_`|~]{1,40}$"), ""),
									},
									"values": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 128),
										},
										Required: true,
										Set:      schema.HashString,
									},
								},
							},
						},
						"http_request_method": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_]{1,40}$`), ""),
										},
										Required: true,
										Set:      schema.HashString,
									},
								},
							},
						},
						"path_pattern": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringLenBetween(1, 128),
										},
										Set: schema.HashString,
									},
								},
							},
						},
						"query_string": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"source_ip": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"values": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: verify.ValidCIDRNetworkAddress,
										},
										Required: true,
										Set:      schema.HashString,
									},
								},
							},
						},
					},
				},
			},
			"tags":     tftags.TagsSchema(),
			"tags_all": tftags.TagsSchemaComputed(),
		},
		CustomizeDiff: customdiff.Sequence(
			verify.SetTagsDiff,
		),
	}
}

func suppressIfActionTypeNot(t string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		take := 2
		i := strings.IndexFunc(k, func(r rune) bool {
			if r == '.' {
				take -= 1
				return take == 0
			}
			return false
		})
		at := k[:i+1] + "type"
		return d.Get(at).(string) != t
	}
}

func resourceListenerRuleCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ELBV2Conn
	listenerArn := d.Get("listener_arn").(string)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(tftags.New(d.Get("tags").(map[string]interface{})))

	params := &elbv2.CreateRuleInput{
		ListenerArn: aws.String(listenerArn),
	}
	if len(tags) > 0 {
		params.Tags = Tags(tags.IgnoreAWS())
	}

	var err error

	params.Actions, err = expandLbListenerActions(d.Get("action").([]interface{}))
	if err != nil {
		return fmt.Errorf("creating LB Listener Rule for Listener (%s): %w", listenerArn, err)
	}

	params.Conditions, err = lbListenerRuleConditions(d.Get("condition").(*schema.Set).List())
	if err != nil {
		return fmt.Errorf("creating LB Listener Rule for Listener (%s): %w", listenerArn, err)
	}

	resp, err := retryListenerRuleCreate(conn, d, params, listenerArn)

	// Some partitions may not support tag-on-create
	if params.Tags != nil && verify.ErrorISOUnsupported(conn.PartitionID, err) {
		log.Printf("[WARN] ELBv2 Listener Rule (%s) create failed (%s) with tags. Trying create without tags.", listenerArn, err)
		params.Tags = nil
		resp, err = retryListenerRuleCreate(conn, d, params, listenerArn)
	}

	if err != nil {
		return fmt.Errorf("Error creating LB Listener Rule: %w", err)
	}

	d.SetId(aws.StringValue(resp.Rules[0].RuleArn))

	// Post-create tagging supported in some partitions
	if params.Tags == nil && len(tags) > 0 {
		err := UpdateTags(conn, d.Id(), nil, tags)

		if v, ok := d.GetOk("tags"); (!ok || len(v.(map[string]interface{})) == 0) && verify.ErrorISOUnsupported(conn.PartitionID, err) {
			// if default tags only, log and continue (i.e., should error if explicitly setting tags and they can't be)
			log.Printf("[WARN] error adding tags after create for ELBv2 Listener Rule (%s): %s", d.Id(), err)
			return resourceListenerRuleRead(d, meta)
		}

		if err != nil {
			return fmt.Errorf("creating ELBv2 Listener Rule (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceListenerRuleRead(d, meta)
}

func resourceListenerRuleRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ELBV2Conn
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	var resp *elbv2.DescribeRulesOutput
	var req = &elbv2.DescribeRulesInput{
		RuleArns: []*string{aws.String(d.Id())},
	}

	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.DescribeRules(req)
		if err != nil {
			if d.IsNewResource() && tfawserr.ErrCodeEquals(err, elbv2.ErrCodeRuleNotFoundException) {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	if tfresource.TimedOut(err) {
		resp, err = conn.DescribeRules(req)
	}
	if err != nil {
		if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeRuleNotFoundException) {
			log.Printf("[WARN] DescribeRules - removing %s from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Rules for listener %q: %w", d.Id(), err)
	}

	if len(resp.Rules) != 1 {
		return fmt.Errorf("Error retrieving Rule %q", d.Id())
	}

	rule := resp.Rules[0]

	d.Set("arn", rule.RuleArn)

	// The listener arn isn't in the response but can be derived from the rule arn
	d.Set("listener_arn", ListenerARNFromRuleARN(aws.StringValue(rule.RuleArn)))

	// Rules are evaluated in priority order, from the lowest value to the highest value. The default rule has the lowest priority.
	if aws.StringValue(rule.Priority) == "default" {
		d.Set("priority", listenerRulePriorityDefault)
	} else {
		if priority, err := strconv.Atoi(aws.StringValue(rule.Priority)); err != nil {
			return fmt.Errorf("Cannot convert rule priority %q to int: %w", aws.StringValue(rule.Priority), err)
		} else {
			d.Set("priority", priority)
		}
	}

	sort.Slice(rule.Actions, func(i, j int) bool {
		return aws.Int64Value(rule.Actions[i].Order) < aws.Int64Value(rule.Actions[j].Order)
	})
	actions := make([]interface{}, len(rule.Actions))
	for i, action := range rule.Actions {
		actionMap := make(map[string]interface{})
		actionMap["type"] = aws.StringValue(action.Type)
		actionMap["order"] = aws.Int64Value(action.Order)

		switch actionMap["type"] {
		case elbv2.ActionTypeEnumForward:
			if aws.StringValue(action.TargetGroupArn) != "" {
				actionMap["target_group_arn"] = aws.StringValue(action.TargetGroupArn)
			} else {
				targetGroups := make([]map[string]interface{}, 0, len(action.ForwardConfig.TargetGroups))
				for _, targetGroup := range action.ForwardConfig.TargetGroups {
					targetGroups = append(targetGroups,
						map[string]interface{}{
							"arn":    aws.StringValue(targetGroup.TargetGroupArn),
							"weight": aws.Int64Value(targetGroup.Weight),
						},
					)
				}
				actionMap["forward"] = []map[string]interface{}{
					{
						"target_group": targetGroups,
						"stickiness": []map[string]interface{}{
							{
								"enabled":  aws.BoolValue(action.ForwardConfig.TargetGroupStickinessConfig.Enabled),
								"duration": aws.Int64Value(action.ForwardConfig.TargetGroupStickinessConfig.DurationSeconds),
							},
						},
					},
				}
			}

		case elbv2.ActionTypeEnumRedirect:
			actionMap["redirect"] = []map[string]interface{}{
				{
					"host":        aws.StringValue(action.RedirectConfig.Host),
					"path":        aws.StringValue(action.RedirectConfig.Path),
					"port":        aws.StringValue(action.RedirectConfig.Port),
					"protocol":    aws.StringValue(action.RedirectConfig.Protocol),
					"query":       aws.StringValue(action.RedirectConfig.Query),
					"status_code": aws.StringValue(action.RedirectConfig.StatusCode),
				},
			}

		case elbv2.ActionTypeEnumFixedResponse:
			actionMap["fixed_response"] = []map[string]interface{}{
				{
					"content_type": aws.StringValue(action.FixedResponseConfig.ContentType),
					"message_body": aws.StringValue(action.FixedResponseConfig.MessageBody),
					"status_code":  aws.StringValue(action.FixedResponseConfig.StatusCode),
				},
			}

		case elbv2.ActionTypeEnumAuthenticateCognito:
			authenticationRequestExtraParams := make(map[string]interface{})
			for key, value := range action.AuthenticateCognitoConfig.AuthenticationRequestExtraParams {
				authenticationRequestExtraParams[key] = aws.StringValue(value)
			}

			actionMap["authenticate_cognito"] = []map[string]interface{}{
				{
					"authentication_request_extra_params": authenticationRequestExtraParams,
					"on_unauthenticated_request":          aws.StringValue(action.AuthenticateCognitoConfig.OnUnauthenticatedRequest),
					"scope":                               aws.StringValue(action.AuthenticateCognitoConfig.Scope),
					"session_cookie_name":                 aws.StringValue(action.AuthenticateCognitoConfig.SessionCookieName),
					"session_timeout":                     aws.Int64Value(action.AuthenticateCognitoConfig.SessionTimeout),
					"user_pool_arn":                       aws.StringValue(action.AuthenticateCognitoConfig.UserPoolArn),
					"user_pool_client_id":                 aws.StringValue(action.AuthenticateCognitoConfig.UserPoolClientId),
					"user_pool_domain":                    aws.StringValue(action.AuthenticateCognitoConfig.UserPoolDomain),
				},
			}

		case elbv2.ActionTypeEnumAuthenticateOidc:
			authenticationRequestExtraParams := make(map[string]interface{})
			for key, value := range action.AuthenticateOidcConfig.AuthenticationRequestExtraParams {
				authenticationRequestExtraParams[key] = aws.StringValue(value)
			}

			// The LB API currently provides no way to read the ClientSecret
			// Instead we passthrough the configuration value into the state
			clientSecret := d.Get("action." + strconv.Itoa(i) + ".authenticate_oidc.0.client_secret").(string)

			actionMap["authenticate_oidc"] = []map[string]interface{}{
				{
					"authentication_request_extra_params": authenticationRequestExtraParams,
					"authorization_endpoint":              aws.StringValue(action.AuthenticateOidcConfig.AuthorizationEndpoint),
					"client_id":                           aws.StringValue(action.AuthenticateOidcConfig.ClientId),
					"client_secret":                       clientSecret,
					"issuer":                              aws.StringValue(action.AuthenticateOidcConfig.Issuer),
					"on_unauthenticated_request":          aws.StringValue(action.AuthenticateOidcConfig.OnUnauthenticatedRequest),
					"scope":                               aws.StringValue(action.AuthenticateOidcConfig.Scope),
					"session_cookie_name":                 aws.StringValue(action.AuthenticateOidcConfig.SessionCookieName),
					"session_timeout":                     aws.Int64Value(action.AuthenticateOidcConfig.SessionTimeout),
					"token_endpoint":                      aws.StringValue(action.AuthenticateOidcConfig.TokenEndpoint),
					"user_info_endpoint":                  aws.StringValue(action.AuthenticateOidcConfig.UserInfoEndpoint),
				},
			}
		}

		actions[i] = actionMap
	}
	d.Set("action", actions)

	conditions := make([]interface{}, len(rule.Conditions))
	for i, condition := range rule.Conditions {
		conditionMap := make(map[string]interface{})

		switch aws.StringValue(condition.Field) {
		case "host-header":
			conditionMap["host_header"] = []interface{}{
				map[string]interface{}{
					"values": flex.FlattenStringSet(condition.HostHeaderConfig.Values),
				},
			}

		case "http-header":
			conditionMap["http_header"] = []interface{}{
				map[string]interface{}{
					"http_header_name": aws.StringValue(condition.HttpHeaderConfig.HttpHeaderName),
					"values":           flex.FlattenStringSet(condition.HttpHeaderConfig.Values),
				},
			}

		case "http-request-method":
			conditionMap["http_request_method"] = []interface{}{
				map[string]interface{}{
					"values": flex.FlattenStringSet(condition.HttpRequestMethodConfig.Values),
				},
			}

		case "path-pattern":
			conditionMap["path_pattern"] = []interface{}{
				map[string]interface{}{
					"values": flex.FlattenStringSet(condition.PathPatternConfig.Values),
				},
			}

		case "query-string":
			values := make([]interface{}, len(condition.QueryStringConfig.Values))
			for k, value := range condition.QueryStringConfig.Values {
				values[k] = map[string]interface{}{
					"key":   aws.StringValue(value.Key),
					"value": aws.StringValue(value.Value),
				}
			}
			conditionMap["query_string"] = values

		case "source-ip":
			conditionMap["source_ip"] = []interface{}{
				map[string]interface{}{
					"values": flex.FlattenStringSet(condition.SourceIpConfig.Values),
				},
			}
		}

		conditions[i] = conditionMap
	}
	if err := d.Set("condition", conditions); err != nil {
		return fmt.Errorf("setting condition: %w", err)
	}

	// tags at the end because, if not supported, will skip the rest of Read
	tags, err := ListTags(conn, d.Id())

	if verify.ErrorISOUnsupported(conn.PartitionID, err) {
		log.Printf("[WARN] Unable to list tags for ELBv2 Listener Rule %s: %s", d.Id(), err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("listing tags for (%s): %w", d.Id(), err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("setting tags_all: %w", err)
	}

	return nil
}

func resourceListenerRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ELBV2Conn

	if d.HasChange("priority") {
		params := &elbv2.SetRulePrioritiesInput{
			RulePriorities: []*elbv2.RulePriorityPair{
				{
					RuleArn:  aws.String(d.Id()),
					Priority: aws.Int64(int64(d.Get("priority").(int))),
				},
			},
		}

		_, err := conn.SetRulePriorities(params)
		if err != nil {
			return err
		}
	}

	requestUpdate := false
	params := &elbv2.ModifyRuleInput{
		RuleArn: aws.String(d.Id()),
	}

	if d.HasChange("action") {
		var err error
		params.Actions, err = expandLbListenerActions(d.Get("action").([]interface{}))
		if err != nil {
			return fmt.Errorf("modifying LB Listener Rule (%s) action: %w", d.Id(), err)
		}
		requestUpdate = true
	}

	if d.HasChange("condition") {
		var err error
		params.Conditions, err = lbListenerRuleConditions(d.Get("condition").(*schema.Set).List())
		if err != nil {
			return fmt.Errorf("modifying LB Listener Rule (%s) condition: %w", d.Id(), err)
		}
		requestUpdate = true
	}

	if requestUpdate {
		resp, err := conn.ModifyRule(params)
		if err != nil {
			return fmt.Errorf("Error modifying LB Listener Rule: %w", err)
		}

		if len(resp.Rules) == 0 {
			return errors.New("Error modifying creating LB Listener Rule: no rules returned in response")
		}
	}

	if d.HasChange("tags_all") {
		o, n := d.GetChange("tags_all")

		err := resource.Retry(loadBalancerTagPropagationTimeout, func() *resource.RetryError {
			err := UpdateTags(conn, d.Id(), o, n)

			if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeLoadBalancerNotFoundException) {
				log.Printf("[DEBUG] Retrying tagging of LB Listener Rule (%s) after error: %s", d.Id(), err)
				return resource.RetryableError(err)
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			err = UpdateTags(conn, d.Id(), o, n)
		}

		// ISO partitions may not support tagging, giving error
		if verify.ErrorISOUnsupported(conn.PartitionID, err) {
			log.Printf("[WARN] Unable to update tags for ELBv2 Listener Rule %s: %s", d.Id(), err)
			return resourceListenerRuleRead(d, meta)
		}

		if err != nil {
			return fmt.Errorf("updating LB (%s) tags: %w", d.Id(), err)
		}
	}

	return resourceListenerRuleRead(d, meta)
}

func resourceListenerRuleDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).ELBV2Conn

	_, err := conn.DeleteRule(&elbv2.DeleteRuleInput{
		RuleArn: aws.String(d.Id()),
	})
	if err != nil && !tfawserr.ErrCodeEquals(err, elbv2.ErrCodeRuleNotFoundException) {
		return fmt.Errorf("Error deleting LB Listener Rule: %w", err)
	}
	return nil
}

func retryListenerRuleCreate(conn *elbv2.ELBV2, d *schema.ResourceData, params *elbv2.CreateRuleInput, listenerARN string) (*elbv2.CreateRuleOutput, error) {
	var resp *elbv2.CreateRuleOutput
	if v, ok := d.GetOk("priority"); ok {
		var err error
		params.Priority = aws.Int64(int64(v.(int)))
		resp, err = conn.CreateRule(params)

		if err != nil {
			return nil, err
		}
	} else {
		var priority int64

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			priority, err = highestListenerRulePriority(conn, listenerARN)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			params.Priority = aws.Int64(priority + 1)
			resp, err = conn.CreateRule(params)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, elbv2.ErrCodePriorityInUseException) {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if tfresource.TimedOut(err) {
			priority, err = highestListenerRulePriority(conn, listenerARN)
			if err != nil {
				return nil, fmt.Errorf("getting highest listener rule (%s) priority: %w", listenerARN, err)
			}
			params.Priority = aws.Int64(priority + 1)
			resp, err = conn.CreateRule(params)
		}

		if err != nil {
			return nil, err
		}
	}

	if resp == nil || len(resp.Rules) == 0 {
		return nil, fmt.Errorf("creating LB Listener Rule (%s): no rules returned in response", listenerARN)
	}

	return resp, nil
}

func validListenerRulePriority(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)
	if value < listenerRulePriorityMin || (value > listenerRulePriorityMax && value != listenerRulePriorityDefault) {
		errors = append(errors, fmt.Errorf("%q must be in the range %d-%d for normal rule or %d for the default rule", k, listenerRulePriorityMin, listenerRulePriorityMax, listenerRulePriorityDefault))
	}
	return
}

// from arn:
// arn:aws:elasticloadbalancing:us-east-1:012345678912:listener-rule/app/name/0123456789abcdef/abcdef0123456789/456789abcedf1234
// select submatches:
// (arn:aws:elasticloadbalancing:us-east-1:012345678912:listener)-rule(/app/name/0123456789abcdef/abcdef0123456789)/456789abcedf1234
// concat to become:
// arn:aws:elasticloadbalancing:us-east-1:012345678912:listener/app/name/0123456789abcdef/abcdef0123456789
var lbListenerARNFromRuleARNRegexp = regexp.MustCompile(`^(arn:.+:listener)-rule(/.+)/[^/]+$`)

func ListenerARNFromRuleARN(ruleArn string) string {
	if arnComponents := lbListenerARNFromRuleARNRegexp.FindStringSubmatch(ruleArn); len(arnComponents) > 1 {
		return arnComponents[1] + arnComponents[2]
	}

	return ""
}

func highestListenerRulePriority(conn *elbv2.ELBV2, arn string) (priority int64, err error) {
	var priorities []int
	var nextMarker *string

	for {
		out, aerr := conn.DescribeRules(&elbv2.DescribeRulesInput{
			ListenerArn: aws.String(arn),
			Marker:      nextMarker,
		})
		if aerr != nil {
			return 0, aerr
		}
		for _, rule := range out.Rules {
			if aws.StringValue(rule.Priority) != "default" {
				p, _ := strconv.Atoi(aws.StringValue(rule.Priority))
				priorities = append(priorities, p)
			}
		}
		if out.NextMarker == nil {
			break
		}
		nextMarker = out.NextMarker
	}

	if len(priorities) == 0 {
		return 0, nil
	}

	sort.IntSlice(priorities).Sort()

	return int64(priorities[len(priorities)-1]), nil
}

// lbListenerRuleConditions converts data source generated by Terraform into
// an elbv2.RuleCondition object suitable for submitting to AWS API.
func lbListenerRuleConditions(conditions []interface{}) ([]*elbv2.RuleCondition, error) {
	elbConditions := make([]*elbv2.RuleCondition, len(conditions))
	for i, condition := range conditions {
		elbConditions[i] = &elbv2.RuleCondition{}
		conditionMap := condition.(map[string]interface{})
		var field string
		var attrs int

		if hostHeader, ok := conditionMap["host_header"].([]interface{}); ok && len(hostHeader) > 0 {
			field = "host-header"
			attrs += 1
			values := hostHeader[0].(map[string]interface{})["values"].(*schema.Set)

			elbConditions[i].HostHeaderConfig = &elbv2.HostHeaderConditionConfig{
				Values: flex.ExpandStringSet(values),
			}
		}

		if httpHeader, ok := conditionMap["http_header"].([]interface{}); ok && len(httpHeader) > 0 {
			field = "http-header"
			attrs += 1
			httpHeaderMap := httpHeader[0].(map[string]interface{})
			values := httpHeaderMap["values"].(*schema.Set)

			elbConditions[i].HttpHeaderConfig = &elbv2.HttpHeaderConditionConfig{
				HttpHeaderName: aws.String(httpHeaderMap["http_header_name"].(string)),
				Values:         flex.ExpandStringSet(values),
			}
		}

		if httpRequestMethod, ok := conditionMap["http_request_method"].([]interface{}); ok && len(httpRequestMethod) > 0 {
			field = "http-request-method"
			attrs += 1
			values := httpRequestMethod[0].(map[string]interface{})["values"].(*schema.Set)

			elbConditions[i].HttpRequestMethodConfig = &elbv2.HttpRequestMethodConditionConfig{
				Values: flex.ExpandStringSet(values),
			}
		}

		if pathPattern, ok := conditionMap["path_pattern"].([]interface{}); ok && len(pathPattern) > 0 {
			field = "path-pattern"
			attrs += 1
			values := pathPattern[0].(map[string]interface{})["values"].(*schema.Set)

			elbConditions[i].PathPatternConfig = &elbv2.PathPatternConditionConfig{
				Values: flex.ExpandStringSet(values),
			}
		}

		if queryString, ok := conditionMap["query_string"].(*schema.Set); ok && queryString.Len() > 0 {
			field = "query-string"
			attrs += 1
			values := queryString.List()

			elbConditions[i].QueryStringConfig = &elbv2.QueryStringConditionConfig{
				Values: make([]*elbv2.QueryStringKeyValuePair, len(values)),
			}
			for j, p := range values {
				valuePair := p.(map[string]interface{})
				elbValuePair := &elbv2.QueryStringKeyValuePair{
					Value: aws.String(valuePair["value"].(string)),
				}
				if valuePair["key"].(string) != "" {
					elbValuePair.Key = aws.String(valuePair["key"].(string))
				}
				elbConditions[i].QueryStringConfig.Values[j] = elbValuePair
			}
		}

		if sourceIp, ok := conditionMap["source_ip"].([]interface{}); ok && len(sourceIp) > 0 {
			field = "source-ip"
			attrs += 1
			values := sourceIp[0].(map[string]interface{})["values"].(*schema.Set)

			elbConditions[i].SourceIpConfig = &elbv2.SourceIpConditionConfig{
				Values: flex.ExpandStringSet(values),
			}
		}

		// FIXME Rework this and use ConflictsWith when it finally works with collections:
		// https://github.com/hashicorp/terraform/issues/13016
		// Still need to ensure that one of the condition attributes is set.
		if attrs == 0 {
			return nil, errors.New("One of host_header, http_header, http_request_method, path_pattern, query_string or source_ip must be set in a condition block")
		} else if attrs > 1 {
			return nil, errors.New("Only one of host_header, http_header, http_request_method, path_pattern, query_string or source_ip can be set in a condition block")
		}

		elbConditions[i].Field = aws.String(field)
	}
	return elbConditions, nil
}
