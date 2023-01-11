package events

import (
	"fmt"
	"log"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eventbridge"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func ResourceConnection() *schema.Resource {
	connectionHttpParameters := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"body": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"is_value_secret": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"header": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"is_value_secret": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"query_string": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"is_value_secret": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}

	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`^[\.\-_A-Za-z0-9]+`), ""),
				),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 512),
			},
			"authorization_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(eventbridge.ConnectionAuthorizationType_Values(), true),
			},
			"auth_parameters": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"auth_parameters.0.api_key",
								"auth_parameters.0.basic",
								"auth_parameters.0.oauth",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 512),
										),
									},
									"value": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 512),
										),
									},
								},
							},
						},
						"basic": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"auth_parameters.0.api_key",
								"auth_parameters.0.basic",
								"auth_parameters.0.oauth",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"username": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 512),
										),
									},
									"password": {
										Type:      schema.TypeString,
										Required:  true,
										Sensitive: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 512),
										),
									},
								},
							},
						},
						"oauth": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							ExactlyOneOf: []string{
								"auth_parameters.0.api_key",
								"auth_parameters.0.basic",
								"auth_parameters.0.oauth",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization_endpoint": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.All(
											validation.StringLenBetween(1, 2048),
										),
									},
									"http_method": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(eventbridge.ConnectionOAuthHttpMethod_Values(), true),
									},
									"oauth_http_parameters": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem:     connectionHttpParameters,
									},
									"client_parameters": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.All(
														validation.StringLenBetween(1, 512),
													),
												},
												"client_secret": {
													Type:      schema.TypeString,
													Required:  true,
													Sensitive: true,
													ValidateFunc: validation.All(
														validation.StringLenBetween(1, 512),
													),
												},
											},
										},
									},
								},
							},
						},
						"invocation_http_parameters": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     connectionHttpParameters,
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EventsConn

	name := d.Get("name").(string)
	input := &eventbridge.CreateConnectionInput{
		AuthorizationType: aws.String(d.Get("authorization_type").(string)),
		AuthParameters:    expandCreateConnectionAuthRequestParameters(d.Get("auth_parameters").([]interface{})),
		Name:              aws.String(name),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating EventBridge connection: %s", input)

	_, err := conn.CreateConnection(input)

	if err != nil {
		return fmt.Errorf("error creating EventBridge connection (%s): %w", name, err)
	}

	d.SetId(name)

	_, err = waitConnectionCreated(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error waiting for EventBridge connection (%s) to create: %w", d.Id(), err)
	}

	return resourceConnectionRead(d, meta)
}

func resourceConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EventsConn

	output, err := FindConnectionByName(conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EventBridge connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading EventBridge connection (%s): %w", d.Id(), err)
	}

	d.Set("arn", output.ConnectionArn)
	d.Set("authorization_type", output.AuthorizationType)
	d.Set("description", output.Description)
	d.Set("name", output.Name)
	d.Set("secret_arn", output.SecretArn)

	if output.AuthParameters != nil {
		authParameters := flattenConnectionAuthParameters(output.AuthParameters, d)
		if err := d.Set("auth_parameters", authParameters); err != nil {
			return fmt.Errorf("error setting auth_parameters error: %w", err)
		}
	}

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EventsConn

	input := &eventbridge.UpdateConnectionInput{
		Name: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("authorization_type"); ok {
		input.AuthorizationType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("auth_parameters"); ok {
		input.AuthParameters = expandUpdateConnectionAuthRequestParameters(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Updating EventBridge connection: %s", input)
	_, err := conn.UpdateConnection(input)

	if err != nil {
		return fmt.Errorf("error updating EventBridge connection (%s): %w", d.Id(), err)
	}

	_, err = waitConnectionUpdated(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error waiting for EventBridge connection (%s) to update: %w", d.Id(), err)
	}

	return resourceConnectionRead(d, meta)
}

func resourceConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EventsConn

	log.Printf("[INFO] Deleting EventBridge connection (%s)", d.Id())
	_, err := conn.DeleteConnection(&eventbridge.DeleteConnectionInput{
		Name: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, eventbridge.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting EventBridge connection (%s): %w", d.Id(), err)
	}

	_, err = waitConnectionDeleted(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error waiting for EventBridge connection (%s) to delete: %w", d.Id(), err)
	}

	return nil
}

func expandCreateConnectionAuthRequestParameters(config []interface{}) *eventbridge.CreateConnectionAuthRequestParameters {
	authParameters := &eventbridge.CreateConnectionAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["api_key"]; ok {
			authParameters.ApiKeyAuthParameters = expandCreateConnectionAPIKeyAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["basic"]; ok {
			authParameters.BasicAuthParameters = expandCreateConnectionBasicAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["oauth"]; ok {
			authParameters.OAuthParameters = expandCreateConnectionOAuthAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["invocation_http_parameters"]; ok {
			authParameters.InvocationHttpParameters = expandConnectionHTTPParameters(val.([]interface{}))
		}
	}

	return authParameters
}

func expandCreateConnectionAPIKeyAuthRequestParameters(config []interface{}) *eventbridge.CreateConnectionApiKeyAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	apiKeyAuthParameters := &eventbridge.CreateConnectionApiKeyAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["key"].(string); ok && val != "" {
			apiKeyAuthParameters.ApiKeyName = aws.String(val)
		}
		if val, ok := param["value"].(string); ok && val != "" {
			apiKeyAuthParameters.ApiKeyValue = aws.String(val)
		}
	}
	return apiKeyAuthParameters
}

func expandCreateConnectionBasicAuthRequestParameters(config []interface{}) *eventbridge.CreateConnectionBasicAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	basicAuthParameters := &eventbridge.CreateConnectionBasicAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["username"].(string); ok && val != "" {
			basicAuthParameters.Username = aws.String(val)
		}
		if val, ok := param["password"].(string); ok && val != "" {
			basicAuthParameters.Password = aws.String(val)
		}
	}
	return basicAuthParameters
}

func expandCreateConnectionOAuthAuthRequestParameters(config []interface{}) *eventbridge.CreateConnectionOAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	oAuthParameters := &eventbridge.CreateConnectionOAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["authorization_endpoint"].(string); ok && val != "" {
			oAuthParameters.AuthorizationEndpoint = aws.String(val)
		}
		if val, ok := param["http_method"].(string); ok && val != "" {
			oAuthParameters.HttpMethod = aws.String(val)
		}
		if val, ok := param["oauth_http_parameters"]; ok {
			oAuthParameters.OAuthHttpParameters = expandConnectionHTTPParameters(val.([]interface{}))
		}
		if val, ok := param["client_parameters"]; ok {
			oAuthParameters.ClientParameters = expandCreateConnectionOAuthClientRequestParameters(val.([]interface{}))
		}
	}
	return oAuthParameters
}

func expandCreateConnectionOAuthClientRequestParameters(config []interface{}) *eventbridge.CreateConnectionOAuthClientRequestParameters {
	oAuthClientRequestParameters := &eventbridge.CreateConnectionOAuthClientRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["client_id"].(string); ok && val != "" {
			oAuthClientRequestParameters.ClientID = aws.String(val)
		}
		if val, ok := param["client_secret"].(string); ok && val != "" {
			oAuthClientRequestParameters.ClientSecret = aws.String(val)
		}
	}
	return oAuthClientRequestParameters
}

func expandConnectionHTTPParameters(config []interface{}) *eventbridge.ConnectionHttpParameters {
	if len(config) == 0 {
		return nil
	}
	httpParameters := &eventbridge.ConnectionHttpParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["body"]; ok {
			httpParameters.BodyParameters = expandConnectionHTTPParametersBody(val.([]interface{}))
		}
		if val, ok := param["header"]; ok {
			httpParameters.HeaderParameters = expandConnectionHTTPParametersHeader(val.([]interface{}))
		}
		if val, ok := param["query_string"]; ok {
			httpParameters.QueryStringParameters = expandConnectionHTTPParametersQueryString(val.([]interface{}))
		}
	}
	return httpParameters
}

func expandConnectionHTTPParametersBody(config []interface{}) []*eventbridge.ConnectionBodyParameter {
	if len(config) == 0 {
		return nil
	}
	var parameters []*eventbridge.ConnectionBodyParameter
	for _, c := range config {
		parameter := eventbridge.ConnectionBodyParameter{}

		input := c.(map[string]interface{})
		if val, ok := input["key"].(string); ok && val != "" {
			parameter.Key = aws.String(val)
		}
		if val, ok := input["value"].(string); ok && val != "" {
			parameter.Value = aws.String(val)
		}
		if val, ok := input["is_value_secret"].(bool); ok {
			parameter.IsValueSecret = aws.Bool(val)
		}
		parameters = append(parameters, &parameter)
	}
	return parameters
}

func expandConnectionHTTPParametersHeader(config []interface{}) []*eventbridge.ConnectionHeaderParameter {
	if len(config) == 0 {
		return nil
	}
	var parameters []*eventbridge.ConnectionHeaderParameter
	for _, c := range config {
		parameter := eventbridge.ConnectionHeaderParameter{}

		input := c.(map[string]interface{})
		if val, ok := input["key"].(string); ok && val != "" {
			parameter.Key = aws.String(val)
		}
		if val, ok := input["value"].(string); ok && val != "" {
			parameter.Value = aws.String(val)
		}
		if val, ok := input["is_value_secret"].(bool); ok {
			parameter.IsValueSecret = aws.Bool(val)
		}
		parameters = append(parameters, &parameter)
	}
	return parameters
}

func expandConnectionHTTPParametersQueryString(config []interface{}) []*eventbridge.ConnectionQueryStringParameter {
	if len(config) == 0 {
		return nil
	}
	var parameters []*eventbridge.ConnectionQueryStringParameter
	for _, c := range config {
		parameter := eventbridge.ConnectionQueryStringParameter{}

		input := c.(map[string]interface{})
		if val, ok := input["key"].(string); ok && val != "" {
			parameter.Key = aws.String(val)
		}
		if val, ok := input["value"].(string); ok && val != "" {
			parameter.Value = aws.String(val)
		}
		if val, ok := input["is_value_secret"].(bool); ok {
			parameter.IsValueSecret = aws.Bool(val)
		}
		parameters = append(parameters, &parameter)
	}
	return parameters
}

func flattenConnectionAuthParameters(
	authParameters *eventbridge.ConnectionAuthResponseParameters,
	resourceData *schema.ResourceData,
) []map[string]interface{} {
	config := make(map[string]interface{})

	if authParameters.ApiKeyAuthParameters != nil {
		config["api_key"] = flattenConnectionAPIKeyAuthParameters(authParameters.ApiKeyAuthParameters, resourceData)
	}

	if authParameters.BasicAuthParameters != nil {
		config["basic"] = flattenConnectionBasicAuthParameters(authParameters.BasicAuthParameters, resourceData)
	}

	if authParameters.OAuthParameters != nil {
		config["oauth"] = flattenConnectionOAuthParameters(authParameters.OAuthParameters, resourceData)
	}

	if authParameters.InvocationHttpParameters != nil {
		config["invocation_http_parameters"] = flattenConnectionHTTPParameters(authParameters.InvocationHttpParameters, resourceData, "auth_parameters.0.invocation_http_parameters")
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenConnectionAPIKeyAuthParameters(apiKeyAuthParameters *eventbridge.ConnectionApiKeyAuthResponseParameters, resourceData *schema.ResourceData) []map[string]interface{} {
	if apiKeyAuthParameters == nil {
		return nil
	}

	config := make(map[string]interface{})
	if apiKeyAuthParameters.ApiKeyName != nil {
		config["key"] = aws.StringValue(apiKeyAuthParameters.ApiKeyName)
	}

	if v, ok := resourceData.GetOk("auth_parameters.0.api_key.0.value"); ok {
		config["value"] = v.(string)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenConnectionBasicAuthParameters(basicAuthParameters *eventbridge.ConnectionBasicAuthResponseParameters, resourceData *schema.ResourceData) []map[string]interface{} {
	if basicAuthParameters == nil {
		return nil
	}

	config := make(map[string]interface{})
	if basicAuthParameters.Username != nil {
		config["username"] = aws.StringValue(basicAuthParameters.Username)
	}

	if v, ok := resourceData.GetOk("auth_parameters.0.basic.0.password"); ok {
		config["password"] = v.(string)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenConnectionOAuthParameters(oAuthParameters *eventbridge.ConnectionOAuthResponseParameters, resourceData *schema.ResourceData) []map[string]interface{} {
	if oAuthParameters == nil {
		return nil
	}

	config := make(map[string]interface{})
	if oAuthParameters.AuthorizationEndpoint != nil {
		config["authorization_endpoint"] = aws.StringValue(oAuthParameters.AuthorizationEndpoint)
	}
	if oAuthParameters.HttpMethod != nil {
		config["http_method"] = aws.StringValue(oAuthParameters.HttpMethod)
	}
	config["oauth_http_parameters"] = flattenConnectionHTTPParameters(oAuthParameters.OAuthHttpParameters, resourceData, "auth_parameters.0.oauth.0.oauth_http_parameters")
	config["client_parameters"] = flattenConnectionOAuthClientResponseParameters(oAuthParameters.ClientParameters, resourceData)

	result := []map[string]interface{}{config}
	return result
}

func flattenConnectionOAuthClientResponseParameters(oAuthClientRequestParameters *eventbridge.ConnectionOAuthClientResponseParameters, resourceData *schema.ResourceData) []map[string]interface{} {
	if oAuthClientRequestParameters == nil {
		return nil
	}

	config := make(map[string]interface{})
	if oAuthClientRequestParameters.ClientID != nil {
		config["client_id"] = aws.StringValue(oAuthClientRequestParameters.ClientID)
	}

	if v, ok := resourceData.GetOk("auth_parameters.0.oauth.0.client_parameters.0.client_secret"); ok {
		config["client_secret"] = v.(string)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenConnectionHTTPParameters(
	httpParameters *eventbridge.ConnectionHttpParameters,
	resourceData *schema.ResourceData,
	path string,
) []map[string]interface{} {
	if httpParameters == nil {
		return nil
	}

	var bodyParameters []map[string]interface{}
	for i, param := range httpParameters.BodyParameters {
		config := make(map[string]interface{})
		config["is_value_secret"] = aws.BoolValue(param.IsValueSecret)
		config["key"] = aws.StringValue(param.Key)

		if param.Value != nil {
			config["value"] = aws.StringValue(param.Value)
		} else if v, ok := resourceData.GetOk(fmt.Sprintf("%s.0.body.%d.value", path, i)); ok {
			config["value"] = v.(string)
		}
		bodyParameters = append(bodyParameters, config)
	}

	var headerParameters []map[string]interface{}
	for i, param := range httpParameters.HeaderParameters {
		config := make(map[string]interface{})
		config["is_value_secret"] = aws.BoolValue(param.IsValueSecret)
		config["key"] = aws.StringValue(param.Key)

		if param.Value != nil {
			config["value"] = aws.StringValue(param.Value)
		} else if v, ok := resourceData.GetOk(fmt.Sprintf("%s.0.header.%d.value", path, i)); ok {
			config["value"] = v.(string)
		}
		headerParameters = append(headerParameters, config)
	}

	var queryStringParameters []map[string]interface{}
	for i, param := range httpParameters.QueryStringParameters {
		config := make(map[string]interface{})
		config["is_value_secret"] = aws.BoolValue(param.IsValueSecret)
		config["key"] = aws.StringValue(param.Key)

		if param.Value != nil {
			config["value"] = aws.StringValue(param.Value)
		} else if v, ok := resourceData.GetOk(fmt.Sprintf("%s.0.query_string.%d.value", path, i)); ok {
			config["value"] = v.(string)
		}
		queryStringParameters = append(queryStringParameters, config)
	}

	parameters := make(map[string]interface{})
	parameters["body"] = bodyParameters
	parameters["header"] = headerParameters
	parameters["query_string"] = queryStringParameters

	result := []map[string]interface{}{parameters}
	return result
}

func expandUpdateConnectionAuthRequestParameters(config []interface{}) *eventbridge.UpdateConnectionAuthRequestParameters {
	authParameters := &eventbridge.UpdateConnectionAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["api_key"]; ok {
			authParameters.ApiKeyAuthParameters = expandUpdateConnectionAPIKeyAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["basic"]; ok {
			authParameters.BasicAuthParameters = expandUpdateConnectionBasicAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["oauth"]; ok {
			authParameters.OAuthParameters = expandUpdateConnectionOAuthAuthRequestParameters(val.([]interface{}))
		}
		if val, ok := param["invocation_http_parameters"]; ok {
			authParameters.InvocationHttpParameters = expandConnectionHTTPParameters(val.([]interface{}))
		}
	}

	return authParameters
}

func expandUpdateConnectionAPIKeyAuthRequestParameters(config []interface{}) *eventbridge.UpdateConnectionApiKeyAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	apiKeyAuthParameters := &eventbridge.UpdateConnectionApiKeyAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["key"].(string); ok && val != "" {
			apiKeyAuthParameters.ApiKeyName = aws.String(val)
		}
		if val, ok := param["value"].(string); ok && val != "" {
			apiKeyAuthParameters.ApiKeyValue = aws.String(val)
		}
	}
	return apiKeyAuthParameters
}

func expandUpdateConnectionBasicAuthRequestParameters(config []interface{}) *eventbridge.UpdateConnectionBasicAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	basicAuthParameters := &eventbridge.UpdateConnectionBasicAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["username"].(string); ok && val != "" {
			basicAuthParameters.Username = aws.String(val)
		}
		if val, ok := param["password"].(string); ok && val != "" {
			basicAuthParameters.Password = aws.String(val)
		}
	}
	return basicAuthParameters
}

func expandUpdateConnectionOAuthAuthRequestParameters(config []interface{}) *eventbridge.UpdateConnectionOAuthRequestParameters {
	if len(config) == 0 {
		return nil
	}
	oAuthParameters := &eventbridge.UpdateConnectionOAuthRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["authorization_endpoint"].(string); ok && val != "" {
			oAuthParameters.AuthorizationEndpoint = aws.String(val)
		}
		if val, ok := param["http_method"].(string); ok && val != "" {
			oAuthParameters.HttpMethod = aws.String(val)
		}
		if val, ok := param["oauth_http_parameters"]; ok {
			oAuthParameters.OAuthHttpParameters = expandConnectionHTTPParameters(val.([]interface{}))
		}
		if val, ok := param["client_parameters"]; ok {
			oAuthParameters.ClientParameters = expandUpdateConnectionOAuthClientRequestParameters(val.([]interface{}))
		}
	}
	return oAuthParameters
}

func expandUpdateConnectionOAuthClientRequestParameters(config []interface{}) *eventbridge.UpdateConnectionOAuthClientRequestParameters {
	oAuthClientRequestParameters := &eventbridge.UpdateConnectionOAuthClientRequestParameters{}
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["client_id"].(string); ok && val != "" {
			oAuthClientRequestParameters.ClientID = aws.String(val)
		}
		if val, ok := param["client_secret"].(string); ok && val != "" {
			oAuthClientRequestParameters.ClientSecret = aws.String(val)
		}
	}
	return oAuthClientRequestParameters
}
