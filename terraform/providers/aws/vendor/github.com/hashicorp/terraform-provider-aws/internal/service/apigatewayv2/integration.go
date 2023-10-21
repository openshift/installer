package apigatewayv2

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_apigatewayv2_integration")
func ResourceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIntegrationCreate,
		ReadWithoutTimeout:   resourceIntegrationRead,
		UpdateWithoutTimeout: resourceIntegrationUpdate,
		DeleteWithoutTimeout: resourceIntegrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIntegrationImport,
		},

		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connection_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 1024),
			},
			"connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      apigatewayv2.ConnectionTypeInternet,
				ValidateFunc: validation.StringInSlice(apigatewayv2.ConnectionType_Values(), false),
			},
			"content_handling_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(apigatewayv2.ContentHandlingStrategy_Values(), false),
			},
			"credentials_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: verify.ValidARN,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"integration_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validHTTPMethod(),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Default HTTP method for Lambda integration is POST.
					if v := d.Get("integration_type").(string); (v == apigatewayv2.IntegrationTypeAws || v == apigatewayv2.IntegrationTypeAwsProxy) && old == "POST" && new == "" {
						return true
					}

					return false
				},
			},
			"integration_response_selection_expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"integration_subtype": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"integration_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(apigatewayv2.IntegrationType_Values(), false),
			},
			"integration_uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"passthrough_behavior": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      apigatewayv2.PassthroughBehaviorWhenNoMatch,
				ValidateFunc: validation.StringInSlice(apigatewayv2.PassthroughBehavior_Values(), false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Not set for HTTP APIs.
					if old == "" && new == apigatewayv2.PassthroughBehaviorWhenNoMatch {
						return true
					}
					return false
				},
			},
			"payload_format_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1.0",
				ValidateFunc: validation.StringInSlice([]string{
					"1.0",
					"2.0",
				}, false),
			},
			"request_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				// Length between [1-512].
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"request_templates": {
				Type:     schema.TypeMap,
				Optional: true,
				// Length between [0-32768].
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"response_parameters": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 0,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mappings": {
							Type:     schema.TypeMap,
							Required: true,
							// Length between [1-512].
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						"status_code": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"template_selection_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeout_milliseconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tls_config": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 0,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_name_to_verify": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	req := &apigatewayv2.CreateIntegrationInput{
		ApiId:           aws.String(d.Get("api_id").(string)),
		IntegrationType: aws.String(d.Get("integration_type").(string)),
	}
	if v, ok := d.GetOk("connection_id"); ok {
		req.ConnectionId = aws.String(v.(string))
	}
	if v, ok := d.GetOk("connection_type"); ok {
		req.ConnectionType = aws.String(v.(string))
	}
	if v, ok := d.GetOk("content_handling_strategy"); ok {
		req.ContentHandlingStrategy = aws.String(v.(string))
	}
	if v, ok := d.GetOk("credentials_arn"); ok {
		req.CredentialsArn = aws.String(v.(string))
	}
	if v, ok := d.GetOk("description"); ok {
		req.Description = aws.String(v.(string))
	}
	if v, ok := d.GetOk("integration_method"); ok {
		req.IntegrationMethod = aws.String(v.(string))
	}
	if v, ok := d.GetOk("integration_subtype"); ok {
		req.IntegrationSubtype = aws.String(v.(string))
	}
	if v, ok := d.GetOk("integration_uri"); ok {
		req.IntegrationUri = aws.String(v.(string))
	}
	if v, ok := d.GetOk("passthrough_behavior"); ok {
		req.PassthroughBehavior = aws.String(v.(string))
	}
	if v, ok := d.GetOk("payload_format_version"); ok {
		req.PayloadFormatVersion = aws.String(v.(string))
	}
	if v, ok := d.GetOk("request_parameters"); ok {
		req.RequestParameters = flex.ExpandStringMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("request_templates"); ok {
		req.RequestTemplates = flex.ExpandStringMap(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("response_parameters"); ok && v.(*schema.Set).Len() > 0 {
		req.ResponseParameters = expandIntegrationResponseParameters(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("template_selection_expression"); ok {
		req.TemplateSelectionExpression = aws.String(v.(string))
	}
	if v, ok := d.GetOk("timeout_milliseconds"); ok {
		req.TimeoutInMillis = aws.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("tls_config"); ok {
		req.TlsConfig = expandTLSConfig(v.([]interface{}))
	}

	log.Printf("[DEBUG] Creating API Gateway v2 integration: %s", req)
	resp, err := conn.CreateIntegrationWithContext(ctx, req)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating API Gateway v2 integration: %s", err)
	}

	d.SetId(aws.StringValue(resp.IntegrationId))

	return append(diags, resourceIntegrationRead(ctx, d, meta)...)
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	resp, err := conn.GetIntegrationWithContext(ctx, &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(d.Get("api_id").(string)),
		IntegrationId: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) && !d.IsNewResource() {
		log.Printf("[WARN] API Gateway v2 integration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway v2 integration: %s", err)
	}

	d.Set("connection_id", resp.ConnectionId)
	d.Set("connection_type", resp.ConnectionType)
	d.Set("content_handling_strategy", resp.ContentHandlingStrategy)
	d.Set("credentials_arn", resp.CredentialsArn)
	d.Set("description", resp.Description)
	d.Set("integration_method", resp.IntegrationMethod)
	d.Set("integration_response_selection_expression", resp.IntegrationResponseSelectionExpression)
	d.Set("integration_subtype", resp.IntegrationSubtype)
	d.Set("integration_type", resp.IntegrationType)
	d.Set("integration_uri", resp.IntegrationUri)
	d.Set("passthrough_behavior", resp.PassthroughBehavior)
	d.Set("payload_format_version", resp.PayloadFormatVersion)
	err = d.Set("request_parameters", flex.PointersMapToStringList(resp.RequestParameters))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting request_parameters: %s", err)
	}
	err = d.Set("request_templates", flex.PointersMapToStringList(resp.RequestTemplates))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting request_templates: %s", err)
	}
	err = d.Set("response_parameters", flattenIntegrationResponseParameters(resp.ResponseParameters))
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "setting response_parameters: %s", err)
	}
	d.Set("template_selection_expression", resp.TemplateSelectionExpression)
	d.Set("timeout_milliseconds", resp.TimeoutInMillis)
	if err := d.Set("tls_config", flattenTLSConfig(resp.TlsConfig)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tls_config: %s", err)
	}

	return diags
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	req := &apigatewayv2.UpdateIntegrationInput{
		ApiId:         aws.String(d.Get("api_id").(string)),
		IntegrationId: aws.String(d.Id()),
		// Always specify the integration type.
		IntegrationType: aws.String(d.Get("integration_type").(string)),
	}
	if d.HasChange("connection_id") {
		req.ConnectionId = aws.String(d.Get("connection_id").(string))
	}
	if d.HasChange("connection_type") {
		req.ConnectionType = aws.String(d.Get("connection_type").(string))
	}
	if d.HasChange("content_handling_strategy") {
		req.ContentHandlingStrategy = aws.String(d.Get("content_handling_strategy").(string))
	}
	if d.HasChange("credentials_arn") {
		req.CredentialsArn = aws.String(d.Get("credentials_arn").(string))
	}
	if d.HasChange("description") {
		req.Description = aws.String(d.Get("description").(string))
	}
	if d.HasChange("integration_method") {
		req.IntegrationMethod = aws.String(d.Get("integration_method").(string))
	}
	// Always specify any integration subtype.
	if v, ok := d.GetOk("integration_subtype"); ok {
		req.IntegrationSubtype = aws.String(v.(string))
	}
	if d.HasChange("integration_uri") {
		req.IntegrationUri = aws.String(d.Get("integration_uri").(string))
	}
	if d.HasChange("passthrough_behavior") {
		req.PassthroughBehavior = aws.String(d.Get("passthrough_behavior").(string))
	}
	if d.HasChange("payload_format_version") {
		req.PayloadFormatVersion = aws.String(d.Get("payload_format_version").(string))
	}
	if d.HasChange("request_parameters") {
		o, n := d.GetChange("request_parameters")
		add, del, nop := verify.DiffStringMaps(o.(map[string]interface{}), n.(map[string]interface{}))

		// Parameters are removed by setting the associated value to "".
		for k := range del {
			del[k] = aws.String("")
		}
		variables := del
		for k, v := range add {
			variables[k] = v
		}
		// Also specify any request parameters that are unchanged as for AWS service integrations some parameters are always required:
		// https://docs.aws.amazon.com/apigateway/latest/developerguide/http-api-develop-integrations-aws-services-reference.html
		for k, v := range nop {
			variables[k] = v
		}

		req.RequestParameters = variables
	}
	if d.HasChange("request_templates") {
		req.RequestTemplates = flex.ExpandStringMap(d.Get("request_templates").(map[string]interface{}))
	}
	if d.HasChange("response_parameters") {
		o, n := d.GetChange("response_parameters")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		del := os.Difference(ns).List()

		req.ResponseParameters = expandIntegrationResponseParameters(ns.List())

		// Parameters are removed by setting the associated value to {}.
		for _, tfMapRaw := range del {
			tfMap, ok := tfMapRaw.(map[string]interface{})

			if !ok {
				continue
			}

			if v, ok := tfMap["status_code"].(string); ok && v != "" {
				if req.ResponseParameters == nil {
					req.ResponseParameters = map[string]map[string]*string{}
				}
				req.ResponseParameters[v] = map[string]*string{}
			}
		}
	}
	if d.HasChange("template_selection_expression") {
		req.TemplateSelectionExpression = aws.String(d.Get("template_selection_expression").(string))
	}
	if d.HasChange("timeout_milliseconds") {
		req.TimeoutInMillis = aws.Int64(int64(d.Get("timeout_milliseconds").(int)))
	}
	if d.HasChange("tls_config") {
		req.TlsConfig = expandTLSConfig(d.Get("tls_config").([]interface{}))
	}

	log.Printf("[DEBUG] Updating API Gateway v2 integration: %s", req)
	_, err := conn.UpdateIntegrationWithContext(ctx, req)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating API Gateway v2 integration: %s", err)
	}

	return append(diags, resourceIntegrationRead(ctx, d, meta)...)
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	log.Printf("[DEBUG] Deleting API Gateway v2 integration (%s)", d.Id())
	_, err := conn.DeleteIntegrationWithContext(ctx, &apigatewayv2.DeleteIntegrationInput{
		ApiId:         aws.String(d.Get("api_id").(string)),
		IntegrationId: aws.String(d.Id()),
	})
	if tfawserr.ErrCodeEquals(err, apigatewayv2.ErrCodeNotFoundException) {
		return diags
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting API Gateway v2 integration: %s", err)
	}

	return diags
}

func resourceIntegrationImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of import ID (%s), use: 'api-id/integration-id'", d.Id())
	}

	apiId := parts[0]
	integrationId := parts[1]

	conn := meta.(*conns.AWSClient).APIGatewayV2Conn(ctx)

	resp, err := conn.GetIntegrationWithContext(ctx, &apigatewayv2.GetIntegrationInput{
		ApiId:         aws.String(apiId),
		IntegrationId: aws.String(integrationId),
	})
	if err != nil {
		return nil, err
	}

	if aws.BoolValue(resp.ApiGatewayManaged) {
		return nil, fmt.Errorf("API Gateway v2 integration (%s) was created via quick create", integrationId)
	}

	d.SetId(integrationId)
	d.Set("api_id", apiId)

	return []*schema.ResourceData{d}, nil
}

func expandTLSConfig(vConfig []interface{}) *apigatewayv2.TlsConfigInput {
	config := &apigatewayv2.TlsConfigInput{}

	if len(vConfig) == 0 || vConfig[0] == nil {
		return config
	}
	mConfig := vConfig[0].(map[string]interface{})

	if vServerNameToVerify, ok := mConfig["server_name_to_verify"].(string); ok && vServerNameToVerify != "" {
		config.ServerNameToVerify = aws.String(vServerNameToVerify)
	}

	return config
}

func flattenTLSConfig(config *apigatewayv2.TlsConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"server_name_to_verify": aws.StringValue(config.ServerNameToVerify),
	}}
}

func expandIntegrationResponseParameters(tfList []interface{}) map[string]map[string]*string {
	if len(tfList) == 0 {
		return nil
	}

	responseParameters := map[string]map[string]*string{}

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		if vStatusCode, ok := tfMap["status_code"].(string); ok && vStatusCode != "" {
			if v, ok := tfMap["mappings"].(map[string]interface{}); ok && len(v) > 0 {
				responseParameters[vStatusCode] = flex.ExpandStringMap(v)
			}
		}
	}

	return responseParameters
}

func flattenIntegrationResponseParameters(responseParameters map[string]map[string]*string) []interface{} {
	if len(responseParameters) == 0 {
		return nil
	}

	var tfList []interface{}

	for statusCode, mappings := range responseParameters {
		if len(mappings) == 0 {
			continue
		}

		tfMap := map[string]interface{}{}

		tfMap["status_code"] = statusCode
		tfMap["mappings"] = aws.StringValueMap(mappings)

		tfList = append(tfList, tfMap)
	}

	return tfList
}
