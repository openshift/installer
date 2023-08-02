package apigateway

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_api_gateway_integration")
func ResourceIntegration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceIntegrationCreate,
		ReadWithoutTimeout:   resourceIntegrationRead,
		UpdateWithoutTimeout: resourceIntegrationUpdate,
		DeleteWithoutTimeout: resourceIntegrationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "/")
				if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected REST-API-ID/RESOURCE-ID/HTTP-METHOD", d.Id())
				}
				restApiID := idParts[0]
				resourceID := idParts[1]
				httpMethod := idParts[2]
				d.Set("http_method", httpMethod)
				d.Set("resource_id", resourceID)
				d.Set("rest_api_id", restApiID)
				d.SetId(fmt.Sprintf("agi-%s-%s-%s", restApiID, resourceID, httpMethod))
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"cache_key_parameters": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"cache_namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connection_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      apigateway.ConnectionTypeInternet,
				ValidateFunc: validation.StringInSlice(apigateway.ConnectionType_Values(), false),
			},
			"content_handling": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validIntegrationContentHandling(),
			},
			"credentials": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"http_method": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validHTTPMethod(),
			},
			"integration_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validHTTPMethod(),
			},
			"passthrough_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"WHEN_NO_MATCH",
					"WHEN_NO_TEMPLATES",
					"NEVER",
				}, false),
			},
			"request_parameters": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"request_templates": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"rest_api_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout_milliseconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(50, 29000),
				Default:      29000,
			},
			"tls_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"insecure_skip_verification": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(apigateway.IntegrationType_Values(), false),
			},
			"uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIntegrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	input := &apigateway.PutIntegrationInput{
		HttpMethod: aws.String(d.Get("http_method").(string)),
		ResourceId: aws.String(d.Get("resource_id").(string)),
		RestApiId:  aws.String(d.Get("rest_api_id").(string)),
		Type:       aws.String(d.Get("type").(string)),
	}

	if v, ok := d.GetOk("cache_key_parameters"); ok && v.(*schema.Set).Len() > 0 {
		input.CacheKeyParameters = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("cache_namespace"); ok {
		input.CacheNamespace = aws.String(v.(string))
	} else if input.CacheKeyParameters != nil {
		input.CacheNamespace = aws.String(d.Get("resource_id").(string))
	}

	if v, ok := d.GetOk("connection_id"); ok {
		input.ConnectionId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("connection_type"); ok {
		input.ConnectionType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_handling"); ok {
		input.ContentHandling = aws.String(v.(string))
	}

	if v, ok := d.GetOk("credentials"); ok {
		input.Credentials = aws.String(v.(string))
	}

	if v, ok := d.GetOk("integration_http_method"); ok {
		input.IntegrationHttpMethod = aws.String(v.(string))
	}

	if v, ok := d.GetOk("passthrough_behavior"); ok {
		input.PassthroughBehavior = aws.String(v.(string))
	}

	if v, ok := d.GetOk("request_parameters"); ok && len(v.(map[string]interface{})) > 0 {
		input.RequestParameters = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("request_templates"); ok && len(v.(map[string]interface{})) > 0 {
		input.RequestTemplates = flex.ExpandStringMap(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("timeout_milliseconds"); ok {
		input.TimeoutInMillis = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("tls_config"); ok && len(v.([]interface{})) > 0 {
		input.TlsConfig = expandTLSConfig(v.([]interface{}))
	}

	if v, ok := d.GetOk("uri"); ok {
		input.Uri = aws.String(v.(string))
	}

	_, err := conn.PutIntegrationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating API Gateway Integration: %s", err)
	}

	d.SetId(fmt.Sprintf("agi-%s-%s-%s", d.Get("rest_api_id").(string), d.Get("resource_id").(string), d.Get("http_method").(string)))

	return append(diags, resourceIntegrationRead(ctx, d, meta)...)
}

func resourceIntegrationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	integration, err := FindIntegrationByThreePartKey(ctx, conn, d.Get("http_method").(string), d.Get("resource_id").(string), d.Get("rest_api_id").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] API Gateway Integration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading API Gateway Integration (%s): %s", d.Id(), err)
	}

	d.Set("cache_key_parameters", aws.StringValueSlice(integration.CacheKeyParameters))
	d.Set("cache_namespace", integration.CacheNamespace)
	d.Set("connection_id", integration.ConnectionId)
	d.Set("connection_type", apigateway.ConnectionTypeInternet)
	if integration.ConnectionType != nil { // nosemgrep:ci.helper-schema-ResourceData-Set-extraneous-nil-check
		d.Set("connection_type", integration.ConnectionType)
	}
	d.Set("content_handling", integration.ContentHandling)
	d.Set("credentials", integration.Credentials)
	d.Set("integration_http_method", integration.HttpMethod)
	d.Set("passthrough_behavior", integration.PassthroughBehavior)
	d.Set("request_parameters", aws.StringValueMap(integration.RequestParameters))
	// We need to explicitly convert key = nil values into key = "", which aws.StringValueMap() removes
	requestTemplates := make(map[string]string)
	for k, v := range integration.RequestTemplates {
		requestTemplates[k] = aws.StringValue(v)
	}
	d.Set("request_templates", requestTemplates)
	d.Set("timeout_milliseconds", integration.TimeoutInMillis)
	d.Set("type", integration.Type)
	d.Set("uri", integration.Uri)

	if err := d.Set("tls_config", flattenTLSConfig(integration.TlsConfig)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tls_config: %s", err)
	}

	return diags
}

func resourceIntegrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	operations := make([]*apigateway.PatchOperation, 0)

	// https://docs.aws.amazon.com/apigateway/api-reference/link-relation/integration-update/#remarks
	// According to the above documentation, only a few parts are addable / removable.
	if d.HasChange("request_templates") {
		o, n := d.GetChange("request_templates")
		prefix := "requestTemplates"

		os := o.(map[string]interface{})
		ns := n.(map[string]interface{})

		// Handle Removal
		for k := range os {
			if _, ok := ns[k]; !ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
				})
			}
		}

		for k, v := range ns {
			// Handle replaces
			if _, ok := os[k]; ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
					Value: aws.String(v.(string)),
				})
			}

			// Handle additions
			if _, ok := os[k]; !ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:    aws.String(apigateway.OpAdd),
					Path:  aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
					Value: aws.String(v.(string)),
				})
			}
		}
	}

	if d.HasChange("request_parameters") {
		o, n := d.GetChange("request_parameters")
		prefix := "requestParameters"

		os := o.(map[string]interface{})
		ns := n.(map[string]interface{})

		// Handle Removal
		for k := range os {
			if _, ok := ns[k]; !ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:   aws.String(apigateway.OpRemove),
					Path: aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
				})
			}
		}

		for k, v := range ns {
			// Handle replaces
			if _, ok := os[k]; ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:    aws.String(apigateway.OpReplace),
					Path:  aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
					Value: aws.String(v.(string)),
				})
			}

			// Handle additions
			if _, ok := os[k]; !ok {
				operations = append(operations, &apigateway.PatchOperation{
					Op:    aws.String(apigateway.OpAdd),
					Path:  aws.String(fmt.Sprintf("/%s/%s", prefix, strings.Replace(k, "/", "~1", -1))),
					Value: aws.String(v.(string)),
				})
			}
		}
	}

	if d.HasChange("cache_key_parameters") {
		o, n := d.GetChange("cache_key_parameters")

		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removalList := os.Difference(ns)
		for _, v := range removalList.List() {
			operations = append(operations, &apigateway.PatchOperation{
				Op:    aws.String(apigateway.OpRemove),
				Path:  aws.String(fmt.Sprintf("/cacheKeyParameters/%s", v.(string))),
				Value: aws.String(""),
			})
		}

		additionList := ns.Difference(os)
		for _, v := range additionList.List() {
			operations = append(operations, &apigateway.PatchOperation{
				Op:    aws.String(apigateway.OpAdd),
				Path:  aws.String(fmt.Sprintf("/cacheKeyParameters/%s", v.(string))),
				Value: aws.String(""),
			})
		}
	}

	if d.HasChange("cache_namespace") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/cacheNamespace"),
			Value: aws.String(d.Get("cache_namespace").(string)),
		})
	}

	// The documentation https://docs.aws.amazon.com/apigateway/api-reference/link-relation/integration-update/ says
	// that uri changes are only supported for non-mock types. Because the uri value is not used in mock
	// resources, it means that the uri can always be updated
	if d.HasChange("uri") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/uri"),
			Value: aws.String(d.Get("uri").(string)),
		})
	}

	if d.HasChange("content_handling") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/contentHandling"),
			Value: aws.String(d.Get("content_handling").(string)),
		})
	}

	if d.HasChange("connection_type") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/connectionType"),
			Value: aws.String(d.Get("connection_type").(string)),
		})
	}

	if d.HasChange("connection_id") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/connectionId"),
			Value: aws.String(d.Get("connection_id").(string)),
		})
	}

	if d.HasChange("timeout_milliseconds") {
		operations = append(operations, &apigateway.PatchOperation{
			Op:    aws.String(apigateway.OpReplace),
			Path:  aws.String("/timeoutInMillis"),
			Value: aws.String(strconv.Itoa(d.Get("timeout_milliseconds").(int))),
		})
	}

	if d.HasChange("tls_config") {
		if v, ok := d.GetOk("tls_config"); ok && len(v.([]interface{})) > 0 {
			m := v.([]interface{})[0].(map[string]interface{})

			operations = append(operations, &apigateway.PatchOperation{
				Op:    aws.String(apigateway.OpReplace),
				Path:  aws.String("/tlsConfig/insecureSkipVerification"),
				Value: aws.String(strconv.FormatBool(m["insecure_skip_verification"].(bool))),
			})
		}
	}

	input := &apigateway.UpdateIntegrationInput{
		HttpMethod:      aws.String(d.Get("http_method").(string)),
		PatchOperations: operations,
		ResourceId:      aws.String(d.Get("resource_id").(string)),
		RestApiId:       aws.String(d.Get("rest_api_id").(string)),
	}

	_, err := conn.UpdateIntegrationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "updating API Gateway Integration (%s): %s", d.Id(), err)
	}

	return append(diags, resourceIntegrationRead(ctx, d, meta)...)
}

func resourceIntegrationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).APIGatewayConn(ctx)

	log.Printf("[DEBUG] Deleting API Gateway Integration: %s", d.Id())
	_, err := conn.DeleteIntegrationWithContext(ctx, &apigateway.DeleteIntegrationInput{
		HttpMethod: aws.String(d.Get("http_method").(string)),
		ResourceId: aws.String(d.Get("resource_id").(string)),
		RestApiId:  aws.String(d.Get("rest_api_id").(string)),
	})

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting API Gateway Integration (%s): %s", d.Id(), err)
	}

	return diags
}

func FindIntegrationByThreePartKey(ctx context.Context, conn *apigateway.APIGateway, httpMethod, resourceID, apiID string) (*apigateway.Integration, error) {
	input := &apigateway.GetIntegrationInput{
		HttpMethod: aws.String(httpMethod),
		ResourceId: aws.String(resourceID),
		RestApiId:  aws.String(apiID),
	}

	output, err := conn.GetIntegrationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, apigateway.ErrCodeNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func expandTLSConfig(vConfig []interface{}) *apigateway.TlsConfig {
	config := &apigateway.TlsConfig{}

	if len(vConfig) == 0 || vConfig[0] == nil {
		return config
	}
	mConfig := vConfig[0].(map[string]interface{})

	if insecureSkipVerification, ok := mConfig["insecure_skip_verification"].(bool); ok {
		config.InsecureSkipVerification = aws.Bool(insecureSkipVerification)
	}
	return config
}

func flattenTLSConfig(config *apigateway.TlsConfig) []interface{} {
	if config == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"insecure_skip_verification": aws.BoolValue(config.InsecureSkipVerification),
	}}
}
