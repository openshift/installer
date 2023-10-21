package elbv2

import (
	"context"
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
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_alb_listener", name="Listener")
// @SDKResource("aws_lb_listener", name="Listener")
// @Tags(identifierAttribute="id")
func ResourceListener() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceListenerCreate,
		ReadWithoutTimeout:   resourceListenerRead,
		UpdateWithoutTimeout: resourceListenerUpdate,
		DeleteWithoutTimeout: resourceListenerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.Sequence(
			verify.SetTagsDiff,
		),

		Schema: map[string]*schema.Schema{
			"alpn_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP1Only",
					"HTTP2Only",
					"HTTP2Optional",
					"HTTP2Preferred",
					"None",
				}, true),
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			"default_action": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authenticate_cognito": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumAuthenticateCognito),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authentication_request_extra_params": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"on_unauthenticated_request": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice(
											elbv2.AuthenticateCognitoActionConditionalBehaviorEnum_Values(),
											true,
										),
									},
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"session_cookie_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"session_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
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
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumAuthenticateOidc),
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
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ValidateFunc: validation.StringInSlice(
											elbv2.AuthenticateOidcActionConditionalBehaviorEnum_Values(),
											true,
										),
									},
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"session_cookie_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"session_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
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
						"fixed_response": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumFixedResponse),
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
										Type:     schema.TypeString,
										Optional: true,
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
						"forward": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumForward),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_group": {
										Type:     schema.TypeSet,
										MinItems: 1,
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
												"duration": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 604800),
												},
												"enabled": {
													Type:     schema.TypeBool,
													Optional: true,
													Default:  false,
												},
											},
										},
									},
								},
							},
						},
						"order": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 50000),
						},
						"redirect": {
							Type:             schema.TypeList,
							Optional:         true,
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumRedirect),
							MaxItems:         1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "#{host}",
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "/#{path}",
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
										Type:     schema.TypeString,
										Optional: true,
										Default:  "#{query}",
									},
									"status_code": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											elbv2.RedirectActionStatusCodeEnum_Values(),
											false,
										),
									},
								},
							},
						},
						"target_group_arn": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressIfDefaultActionTypeNot(elbv2.ActionTypeEnumForward),
							ValidateFunc:     verify.ValidARN,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								elbv2.ActionTypeEnum_Values(),
								true,
							),
						},
					},
				},
			},
			"load_balancer_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IsPortNumber,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					return strings.ToUpper(v.(string))
				},
				ValidateFunc: validation.StringInSlice(elbv2.ProtocolEnum_Values(), true),
			},
			"ssl_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func suppressIfDefaultActionTypeNot(t string) schema.SchemaDiffSuppressFunc {
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

func resourceListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBV2Conn(ctx)

	lbARN := d.Get("load_balancer_arn").(string)
	input := &elbv2.CreateListenerInput{
		LoadBalancerArn: aws.String(lbARN),
		Tags:            GetTagsIn(ctx),
	}

	if alpnPolicy, ok := d.GetOk("alpn_policy"); ok {
		input.AlpnPolicy = make([]*string, 1)
		input.AlpnPolicy[0] = aws.String(alpnPolicy.(string))
	}

	if certificateArn, ok := d.GetOk("certificate_arn"); ok {
		input.Certificates = make([]*elbv2.Certificate, 1)
		input.Certificates[0] = &elbv2.Certificate{
			CertificateArn: aws.String(certificateArn.(string)),
		}
	}

	if v, ok := d.GetOk("default_action"); ok && len(v.([]interface{})) > 0 {
		var err error
		input.DefaultActions, err = expandLbListenerActions(v.([]interface{}))
		if err != nil {
			return sdkdiag.AppendFromErr(diags, err)
		}
	}

	if v, ok := d.GetOk("port"); ok {
		input.Port = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("protocol"); ok {
		input.Protocol = aws.String(v.(string))
	} else if strings.Contains(lbARN, "loadbalancer/app/") {
		// Keep previous default of HTTP for Application Load Balancers
		input.Protocol = aws.String(elbv2.ProtocolEnumHttp)
	}

	if sslPolicy, ok := d.GetOk("ssl_policy"); ok {
		input.SslPolicy = aws.String(sslPolicy.(string))
	}

	output, err := retryListenerCreate(ctx, conn, input)

	// Some partitions (e.g. ISO) may not support tag-on-create.
	if input.Tags != nil && errs.IsUnsupportedOperationInPartitionError(conn.PartitionID, err) {
		input.Tags = nil

		output, err = retryListenerCreate(ctx, conn, input)
	}

	// Tags are not supported on creation with some load balancer types (i.e. Gateway)
	// Retry creation without tags
	if input.Tags != nil && tfawserr.ErrMessageContains(err, ErrValidationError, TagsOnCreationErrMessage) {
		input.Tags = nil

		output, err = retryListenerCreate(ctx, conn, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating ELBv2 Listener (%s): %s", lbARN, err)
	}

	d.SetId(aws.StringValue(output.Listeners[0].ListenerArn))

	// For partitions not supporting tag-on-create, attempt tag after create.
	if tags := GetTagsIn(ctx); input.Tags == nil && len(tags) > 0 {
		err := createTags(ctx, conn, d.Id(), tags)

		// If default tags only, continue. Otherwise, error.
		if v, ok := d.GetOk(names.AttrTags); (!ok || len(v.(map[string]interface{})) == 0) && errs.IsUnsupportedOperationInPartitionError(conn.PartitionID, err) {
			return append(diags, resourceListenerRead(ctx, d, meta)...)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "setting ELBv2 Listener (%s) tags: %s", d.Id(), err)
		}
	}

	return append(diags, resourceListenerRead(ctx, d, meta)...)
}

func resourceListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	const (
		loadBalancerListenerReadTimeout = 2 * time.Minute
	)
	conn := meta.(*conns.AWSClient).ELBV2Conn(ctx)

	var listener *elbv2.Listener

	err := retry.RetryContext(ctx, loadBalancerListenerReadTimeout, func() *retry.RetryError {
		var err error
		listener, err = FindListenerByARN(ctx, conn, d.Id())

		if d.IsNewResource() && tfawserr.ErrCodeEquals(err, elbv2.ErrCodeListenerNotFoundException) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		listener, err = FindListenerByARN(ctx, conn, d.Id())
	}

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, elbv2.ErrCodeListenerNotFoundException) {
		log.Printf("[WARN] ELBv2 Listener (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "describing ELBv2 Listener (%s): %s", d.Id(), err)
	}

	if listener == nil {
		if d.IsNewResource() {
			return sdkdiag.AppendErrorf(diags, "describing ELBv2 Listener (%s): empty response", d.Id())
		}
		log.Printf("[WARN] ELBv2 Listener (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	d.Set("arn", listener.ListenerArn)
	d.Set("load_balancer_arn", listener.LoadBalancerArn)
	d.Set("port", listener.Port)
	d.Set("protocol", listener.Protocol)
	d.Set("ssl_policy", listener.SslPolicy)

	if listener.Certificates != nil && len(listener.Certificates) == 1 && listener.Certificates[0] != nil {
		d.Set("certificate_arn", listener.Certificates[0].CertificateArn)
	}

	if listener.AlpnPolicy != nil && len(listener.AlpnPolicy) == 1 && listener.AlpnPolicy[0] != nil {
		d.Set("alpn_policy", listener.AlpnPolicy[0])
	}

	sort.Slice(listener.DefaultActions, func(i, j int) bool {
		return aws.Int64Value(listener.DefaultActions[i].Order) < aws.Int64Value(listener.DefaultActions[j].Order)
	})

	if err := d.Set("default_action", flattenLbListenerActions(d, listener.DefaultActions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting default_action for ELBv2 listener (%s): %s", d.Id(), err)
	}

	return diags
}

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	const (
		loadBalancerListenerUpdateTimeout = 5 * time.Minute
	)
	conn := meta.(*conns.AWSClient).ELBV2Conn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &elbv2.ModifyListenerInput{
			ListenerArn: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("port"); ok {
			input.Port = aws.Int64(int64(v.(int)))
		}

		if v, ok := d.GetOk("protocol"); ok {
			input.Protocol = aws.String(v.(string))
		}

		if v, ok := d.GetOk("ssl_policy"); ok {
			input.SslPolicy = aws.String(v.(string))
		}

		if v, ok := d.GetOk("certificate_arn"); ok {
			input.Certificates = make([]*elbv2.Certificate, 1)
			input.Certificates[0] = &elbv2.Certificate{
				CertificateArn: aws.String(v.(string)),
			}
		}

		if v, ok := d.GetOk("alpn_policy"); ok {
			input.AlpnPolicy = aws.StringSlice([]string{v.(string)})
		}

		if d.HasChange("default_action") {
			var err error
			input.DefaultActions, err = expandLbListenerActions(d.Get("default_action").([]interface{}))
			if err != nil {
				return sdkdiag.AppendFromErr(diags, err)
			}
		}

		err := retry.RetryContext(ctx, loadBalancerListenerUpdateTimeout, func() *retry.RetryError {
			_, err := conn.ModifyListenerWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeCertificateNotFoundException) {
				return retry.RetryableError(err)
			}

			if err != nil {
				return retry.NonRetryableError(err)
			}

			return nil
		})

		if tfresource.TimedOut(err) {
			_, err = conn.ModifyListenerWithContext(ctx, input)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying ELBv2 Listener (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceListenerRead(ctx, d, meta)...)
}

func resourceListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBV2Conn(ctx)

	_, err := conn.DeleteListenerWithContext(ctx, &elbv2.DeleteListenerInput{
		ListenerArn: aws.String(d.Id()),
	})
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Listener (%s): %s", d.Id(), err)
	}

	return diags
}

func retryListenerCreate(ctx context.Context, conn *elbv2.ELBV2, params *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
	const (
		loadBalancerListenerCreateTimeout = 5 * time.Minute
	)
	var output *elbv2.CreateListenerOutput

	err := retry.RetryContext(ctx, loadBalancerListenerCreateTimeout, func() *retry.RetryError {
		var err error

		output, err = conn.CreateListenerWithContext(ctx, params)

		if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeCertificateNotFoundException) {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateListenerWithContext(ctx, params)
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.Listeners) == 0 {
		return nil, fmt.Errorf("creating ELBv2 Listener: no listeners returned in response")
	}

	return output, nil
}

func expandLbListenerActions(l []interface{}) ([]*elbv2.Action, error) {
	if len(l) == 0 {
		return nil, nil
	}

	var actions []*elbv2.Action
	var err error

	for i, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		action := &elbv2.Action{
			Order: aws.Int64(int64(i + 1)),
			Type:  aws.String(tfMap["type"].(string)),
		}

		if order, ok := tfMap["order"].(int); ok && order != 0 {
			action.Order = aws.Int64(int64(order))
		}

		switch tfMap["type"].(string) {
		case elbv2.ActionTypeEnumForward:
			if v, ok := tfMap["target_group_arn"].(string); ok && v != "" {
				action.TargetGroupArn = aws.String(v)
			} else if v, ok := tfMap["forward"].([]interface{}); ok {
				action.ForwardConfig = expandLbListenerActionForwardConfig(v)
			} else {
				err = errors.New("for actions of type 'forward', you must specify a 'forward' block or 'target_group_arn'")
			}

		case elbv2.ActionTypeEnumRedirect:
			if v, ok := tfMap["redirect"].([]interface{}); ok {
				action.RedirectConfig = expandLbListenerRedirectActionConfig(v)
			} else {
				err = errors.New("for actions of type 'redirect', you must specify a 'redirect' block")
			}

		case elbv2.ActionTypeEnumFixedResponse:
			if v, ok := tfMap["fixed_response"].([]interface{}); ok {
				action.FixedResponseConfig = expandLbListenerFixedResponseConfig(v)
			} else {
				err = errors.New("for actions of type 'fixed-response', you must specify a 'fixed_response' block")
			}

		case elbv2.ActionTypeEnumAuthenticateCognito:
			if v, ok := tfMap["authenticate_cognito"].([]interface{}); ok {
				action.AuthenticateCognitoConfig = expandLbListenerAuthenticateCognitoConfig(v)
			} else {
				err = errors.New("for actions of type 'authenticate-cognito', you must specify a 'authenticate_cognito' block")
			}

		case elbv2.ActionTypeEnumAuthenticateOidc:
			if v, ok := tfMap["authenticate_oidc"].([]interface{}); ok {
				action.AuthenticateOidcConfig = expandAuthenticateOIDCConfig(v)
			} else {
				err = errors.New("for actions of type 'authenticate-oidc', you must specify a 'authenticate_oidc' block")
			}
		}

		actions = append(actions, action)
	}

	return actions, err
}

func expandLbListenerAuthenticateCognitoConfig(l []interface{}) *elbv2.AuthenticateCognitoActionConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	config := &elbv2.AuthenticateCognitoActionConfig{
		AuthenticationRequestExtraParams: flex.ExpandStringMap(tfMap["authentication_request_extra_params"].(map[string]interface{})),
		UserPoolArn:                      aws.String(tfMap["user_pool_arn"].(string)),
		UserPoolClientId:                 aws.String(tfMap["user_pool_client_id"].(string)),
		UserPoolDomain:                   aws.String(tfMap["user_pool_domain"].(string)),
	}

	if v, ok := tfMap["on_unauthenticated_request"].(string); ok && v != "" {
		config.OnUnauthenticatedRequest = aws.String(v)
	}

	if v, ok := tfMap["scope"].(string); ok && v != "" {
		config.Scope = aws.String(v)
	}

	if v, ok := tfMap["session_cookie_name"].(string); ok && v != "" {
		config.SessionCookieName = aws.String(v)
	}

	if v, ok := tfMap["session_timeout"].(int); ok && v != 0 {
		config.SessionTimeout = aws.Int64(int64(v))
	}

	return config
}

func expandAuthenticateOIDCConfig(l []interface{}) *elbv2.AuthenticateOidcActionConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	config := &elbv2.AuthenticateOidcActionConfig{
		AuthenticationRequestExtraParams: flex.ExpandStringMap(tfMap["authentication_request_extra_params"].(map[string]interface{})),
		AuthorizationEndpoint:            aws.String(tfMap["authorization_endpoint"].(string)),
		ClientId:                         aws.String(tfMap["client_id"].(string)),
		ClientSecret:                     aws.String(tfMap["client_secret"].(string)),
		Issuer:                           aws.String(tfMap["issuer"].(string)),
		TokenEndpoint:                    aws.String(tfMap["token_endpoint"].(string)),
		UserInfoEndpoint:                 aws.String(tfMap["user_info_endpoint"].(string)),
	}

	if v, ok := tfMap["on_unauthenticated_request"].(string); ok && v != "" {
		config.OnUnauthenticatedRequest = aws.String(v)
	}

	if v, ok := tfMap["scope"].(string); ok && v != "" {
		config.Scope = aws.String(v)
	}

	if v, ok := tfMap["session_cookie_name"].(string); ok && v != "" {
		config.SessionCookieName = aws.String(v)
	}

	if v, ok := tfMap["session_timeout"].(int); ok && v != 0 {
		config.SessionTimeout = aws.Int64(int64(v))
	}

	return config
}

func expandLbListenerFixedResponseConfig(l []interface{}) *elbv2.FixedResponseActionConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	return &elbv2.FixedResponseActionConfig{
		ContentType: aws.String(tfMap["content_type"].(string)),
		MessageBody: aws.String(tfMap["message_body"].(string)),
		StatusCode:  aws.String(tfMap["status_code"].(string)),
	}
}

func expandLbListenerRedirectActionConfig(l []interface{}) *elbv2.RedirectActionConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})

	if !ok {
		return nil
	}

	return &elbv2.RedirectActionConfig{
		Host:       aws.String(tfMap["host"].(string)),
		Path:       aws.String(tfMap["path"].(string)),
		Port:       aws.String(tfMap["port"].(string)),
		Protocol:   aws.String(tfMap["protocol"].(string)),
		Query:      aws.String(tfMap["query"].(string)),
		StatusCode: aws.String(tfMap["status_code"].(string)),
	}
}

func expandLbListenerActionForwardConfig(l []interface{}) *elbv2.ForwardActionConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}

	config := &elbv2.ForwardActionConfig{}

	if v, ok := tfMap["target_group"].(*schema.Set); ok && v.Len() > 0 {
		config.TargetGroups = expandLbListenerActionForwardConfigTargetGroups(v.List())
	}

	if v, ok := tfMap["stickiness"].([]interface{}); ok && len(v) > 0 {
		config.TargetGroupStickinessConfig = expandLbListenerActionForwardConfigTargetGroupStickinessConfig(v)
	}

	return config
}

func expandLbListenerActionForwardConfigTargetGroups(l []interface{}) []*elbv2.TargetGroupTuple {
	if len(l) == 0 {
		return nil
	}

	var groups []*elbv2.TargetGroupTuple

	for _, tfMapRaw := range l {
		tfMap, ok := tfMapRaw.(map[string]interface{})
		if !ok {
			continue
		}

		group := &elbv2.TargetGroupTuple{
			TargetGroupArn: aws.String(tfMap["arn"].(string)),
			Weight:         aws.Int64(int64(tfMap["weight"].(int))),
		}

		groups = append(groups, group)
	}

	return groups
}

func expandLbListenerActionForwardConfigTargetGroupStickinessConfig(l []interface{}) *elbv2.TargetGroupStickinessConfig {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	tfMap, ok := l[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return &elbv2.TargetGroupStickinessConfig{
		Enabled:         aws.Bool(tfMap["enabled"].(bool)),
		DurationSeconds: aws.Int64(int64(tfMap["duration"].(int))),
	}
}

func flattenLbListenerActions(d *schema.ResourceData, Actions []*elbv2.Action) []interface{} {
	if len(Actions) == 0 {
		return []interface{}{}
	}

	var vActions []interface{}

	for i, action := range Actions {
		m := map[string]interface{}{
			"type":  aws.StringValue(action.Type),
			"order": aws.Int64Value(action.Order),
		}

		switch aws.StringValue(action.Type) {
		case elbv2.ActionTypeEnumForward:
			if aws.StringValue(action.TargetGroupArn) != "" {
				m["target_group_arn"] = aws.StringValue(action.TargetGroupArn)
			} else {
				m["forward"] = flattenLbListenerActionForwardConfig(action.ForwardConfig)
			}

		case elbv2.ActionTypeEnumRedirect:
			m["redirect"] = flattenLbListenerActionRedirectConfig(action.RedirectConfig)

		case elbv2.ActionTypeEnumFixedResponse:
			m["fixed_response"] = flattenLbListenerActionFixedResponseConfig(action.FixedResponseConfig)

		case elbv2.ActionTypeEnumAuthenticateCognito:
			m["authenticate_cognito"] = flattenLbListenerActionAuthenticateCognitoConfig(action.AuthenticateCognitoConfig)

		case elbv2.ActionTypeEnumAuthenticateOidc:
			// The LB API currently provides no way to read the ClientSecret
			// Instead we passthrough the configuration value into the state
			var clientSecret string
			if v, ok := d.GetOk("default_action." + strconv.Itoa(i) + ".authenticate_oidc.0.client_secret"); ok {
				clientSecret = v.(string)
			}

			m["authenticate_oidc"] = flattenAuthenticateOIDCActionConfig(action.AuthenticateOidcConfig, clientSecret)
		}

		vActions = append(vActions, m)
	}

	return vActions
}

func flattenAuthenticateOIDCActionConfig(config *elbv2.AuthenticateOidcActionConfig, clientSecret string) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"authentication_request_extra_params": aws.StringValueMap(config.AuthenticationRequestExtraParams),
		"authorization_endpoint":              aws.StringValue(config.AuthorizationEndpoint),
		"client_id":                           aws.StringValue(config.ClientId),
		"client_secret":                       clientSecret,
		"issuer":                              aws.StringValue(config.Issuer),
		"on_unauthenticated_request":          aws.StringValue(config.OnUnauthenticatedRequest),
		"scope":                               aws.StringValue(config.Scope),
		"session_cookie_name":                 aws.StringValue(config.SessionCookieName),
		"session_timeout":                     aws.Int64Value(config.SessionTimeout),
		"token_endpoint":                      aws.StringValue(config.TokenEndpoint),
		"user_info_endpoint":                  aws.StringValue(config.UserInfoEndpoint),
	}

	return []interface{}{m}
}

func flattenLbListenerActionAuthenticateCognitoConfig(config *elbv2.AuthenticateCognitoActionConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"authentication_request_extra_params": aws.StringValueMap(config.AuthenticationRequestExtraParams),
		"on_unauthenticated_request":          aws.StringValue(config.OnUnauthenticatedRequest),
		"scope":                               aws.StringValue(config.Scope),
		"session_cookie_name":                 aws.StringValue(config.SessionCookieName),
		"session_timeout":                     aws.Int64Value(config.SessionTimeout),
		"user_pool_arn":                       aws.StringValue(config.UserPoolArn),
		"user_pool_client_id":                 aws.StringValue(config.UserPoolClientId),
		"user_pool_domain":                    aws.StringValue(config.UserPoolDomain),
	}

	return []interface{}{m}
}

func flattenLbListenerActionFixedResponseConfig(config *elbv2.FixedResponseActionConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"content_type": aws.StringValue(config.ContentType),
		"message_body": aws.StringValue(config.MessageBody),
		"status_code":  aws.StringValue(config.StatusCode),
	}

	return []interface{}{m}
}

func flattenLbListenerActionForwardConfig(config *elbv2.ForwardActionConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"target_group": flattenLbListenerActionForwardConfigTargetGroups(config.TargetGroups),
		"stickiness":   flattenLbListenerActionForwardConfigTargetGroupStickinessConfig(config.TargetGroupStickinessConfig),
	}

	return []interface{}{m}
}

func flattenLbListenerActionForwardConfigTargetGroups(groups []*elbv2.TargetGroupTuple) []interface{} {
	if len(groups) == 0 {
		return []interface{}{}
	}

	var vGroups []interface{}

	for _, group := range groups {
		if group == nil {
			continue
		}

		m := map[string]interface{}{
			"arn":    aws.StringValue(group.TargetGroupArn),
			"weight": aws.Int64Value(group.Weight),
		}

		vGroups = append(vGroups, m)
	}

	return vGroups
}

func flattenLbListenerActionForwardConfigTargetGroupStickinessConfig(config *elbv2.TargetGroupStickinessConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"enabled":  aws.BoolValue(config.Enabled),
		"duration": aws.Int64Value(config.DurationSeconds),
	}

	return []interface{}{m}
}

func flattenLbListenerActionRedirectConfig(config *elbv2.RedirectActionConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"host":        aws.StringValue(config.Host),
		"path":        aws.StringValue(config.Path),
		"port":        aws.StringValue(config.Port),
		"protocol":    aws.StringValue(config.Protocol),
		"query":       aws.StringValue(config.Query),
		"status_code": aws.StringValue(config.StatusCode),
	}

	return []interface{}{m}
}
