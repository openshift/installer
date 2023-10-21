package ec2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ec2_client_vpn_endpoint", name="Client VPN Endpoint")
// @Tags(identifierAttribute="id")
func ResourceClientVPNEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClientVPNEndpointCreate,
		ReadWithoutTimeout:   resourceClientVPNEndpointRead,
		DeleteWithoutTimeout: resourceClientVPNEndpointDelete,
		UpdateWithoutTimeout: resourceClientVPNEndpointUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authentication_options": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active_directory_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"root_certificate_chain_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"saml_provider_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"self_service_saml_provider_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(ec2.ClientVpnAuthenticationType_Values(), false),
						},
					},
				},
			},
			"client_cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"client_connect_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"lambda_function_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			"client_login_banner_options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"banner_text": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 1400),
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"connection_log_options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloudwatch_log_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloudwatch_log_stream": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				MinItems: 1,
				MaxItems: 5,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"self_service_portal": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ec2.SelfServicePortalDisabled,
				ValidateFunc: validation.StringInSlice(ec2.SelfServicePortal_Values(), false),
			},
			"server_certificate_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: verify.ValidARN,
			},
			"session_timeout_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      24,
				ValidateFunc: validation.IntInSlice([]int{8, 10, 12, 24}),
			},
			"split_tunnel": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"transport_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      ec2.TransportProtocolUdp,
				ValidateFunc: validation.StringInSlice(ec2.TransportProtocol_Values(), false),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpn_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  443,
				ValidateFunc: validation.IntInSlice([]int{
					443,
					1194,
				}),
			},
		},
	}
}

func resourceClientVPNEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.CreateClientVpnEndpointInput{
		ClientCidrBlock:      aws.String(d.Get("client_cidr_block").(string)),
		ServerCertificateArn: aws.String(d.Get("server_certificate_arn").(string)),
		SplitTunnel:          aws.Bool(d.Get("split_tunnel").(bool)),
		TagSpecifications:    getTagSpecificationsIn(ctx, ec2.ResourceTypeClientVpnEndpoint),
		TransportProtocol:    aws.String(d.Get("transport_protocol").(string)),
		VpnPort:              aws.Int64(int64(d.Get("vpn_port").(int))),
	}

	if v, ok := d.GetOk("authentication_options"); ok && v.(*schema.Set).Len() > 0 {
		input.AuthenticationOptions = expandClientVPNAuthenticationRequests(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("client_connect_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ClientConnectOptions = expandClientConnectOptions(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("client_login_banner_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ClientLoginBannerOptions = expandClientLoginBannerOptions(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("connection_log_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ConnectionLogOptions = expandConnectionLogOptions(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("dns_servers"); ok && len(v.([]interface{})) > 0 {
		input.DnsServers = flex.ExpandStringList(v.([]interface{}))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		input.SecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("self_service_portal"); ok {
		input.SelfServicePortal = aws.String(v.(string))
	}

	if v, ok := d.GetOk("session_timeout_hours"); ok {
		input.SessionTimeoutHours = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		input.VpcId = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating EC2 Client VPN Endpoint: %s", input)
	output, err := conn.CreateClientVpnEndpointWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 Client VPN Endpoint: %s", err)
	}

	d.SetId(aws.StringValue(output.ClientVpnEndpointId))

	return append(diags, resourceClientVPNEndpointRead(ctx, d, meta)...)
}

func resourceClientVPNEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	ep, err := FindClientVPNEndpointByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Client VPN Endpoint (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Client VPN Endpoint (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("client-vpn-endpoint/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	if err := d.Set("authentication_options", flattenClientVPNAuthentications(ep.AuthenticationOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting authentication_options: %s", err)
	}
	d.Set("client_cidr_block", ep.ClientCidrBlock)
	if ep.ClientConnectOptions != nil {
		if err := d.Set("client_connect_options", []interface{}{flattenClientConnectResponseOptions(ep.ClientConnectOptions)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting client_connect_options: %s", err)
		}
	} else {
		d.Set("client_connect_options", nil)
	}
	if ep.ClientLoginBannerOptions != nil {
		if err := d.Set("client_login_banner_options", []interface{}{flattenClientLoginBannerResponseOptions(ep.ClientLoginBannerOptions)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting client_login_banner_options: %s", err)
		}
	} else {
		d.Set("client_login_banner_options", nil)
	}
	if ep.ConnectionLogOptions != nil {
		if err := d.Set("connection_log_options", []interface{}{flattenConnectionLogResponseOptions(ep.ConnectionLogOptions)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting connection_log_options: %s", err)
		}
	} else {
		d.Set("connection_log_options", nil)
	}
	d.Set("description", ep.Description)
	d.Set("dns_name", ep.DnsName)
	d.Set("dns_servers", aws.StringValueSlice(ep.DnsServers))
	d.Set("security_group_ids", aws.StringValueSlice(ep.SecurityGroupIds))
	if aws.StringValue(ep.SelfServicePortalUrl) != "" {
		d.Set("self_service_portal", ec2.SelfServicePortalEnabled)
	} else {
		d.Set("self_service_portal", ec2.SelfServicePortalDisabled)
	}
	d.Set("server_certificate_arn", ep.ServerCertificateArn)
	d.Set("session_timeout_hours", ep.SessionTimeoutHours)
	d.Set("split_tunnel", ep.SplitTunnel)
	d.Set("transport_protocol", ep.TransportProtocol)
	d.Set("vpc_id", ep.VpcId)
	d.Set("vpn_port", ep.VpnPort)

	SetTagsOut(ctx, ep.Tags)

	return diags
}

func resourceClientVPNEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		var waitForClientConnectResponseOptionsUpdate bool
		input := &ec2.ModifyClientVpnEndpointInput{
			ClientVpnEndpointId: aws.String(d.Id()),
		}

		if d.HasChange("client_connect_options") {
			waitForClientConnectResponseOptionsUpdate = true

			if v, ok := d.GetOk("client_connect_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.ClientConnectOptions = expandClientConnectOptions(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("client_login_banner_options") {
			if v, ok := d.GetOk("client_login_banner_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.ClientLoginBannerOptions = expandClientLoginBannerOptions(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("connection_log_options") {
			if v, ok := d.GetOk("connection_log_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.ConnectionLogOptions = expandConnectionLogOptions(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChange("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		if d.HasChange("dns_servers") {
			dnsServers := d.Get("dns_servers").([]interface{})
			enabled := len(dnsServers) > 0

			input.DnsServers = &ec2.DnsServersOptionsModifyStructure{
				Enabled: aws.Bool(enabled),
			}
			if enabled {
				input.DnsServers.CustomDnsServers = flex.ExpandStringList(dnsServers)
			}
		}

		if d.HasChange("security_group_ids") {
			input.SecurityGroupIds = flex.ExpandStringSet(d.Get("security_group_ids").(*schema.Set))
			// "InvalidParameterValue: Security Groups cannot be modified without specifying Vpc Id"
			input.VpcId = aws.String(d.Get("vpc_id").(string))
		}

		if d.HasChange("self_service_portal") {
			input.SelfServicePortal = aws.String(d.Get("self_service_portal").(string))
		}

		if d.HasChange("session_timeout_hours") {
			input.SessionTimeoutHours = aws.Int64(int64(d.Get("session_timeout_hours").(int)))
		}

		if d.HasChange("server_certificate_arn") {
			input.ServerCertificateArn = aws.String(d.Get("server_certificate_arn").(string))
		}

		if d.HasChange("split_tunnel") {
			input.SplitTunnel = aws.Bool(d.Get("split_tunnel").(bool))
		}

		if d.HasChange("vpn_port") {
			input.VpnPort = aws.Int64(int64(d.Get("vpn_port").(int)))
		}

		if d.HasChange("vpc_id") {
			input.VpcId = aws.String(d.Get("vpc_id").(string))
		}

		if _, err := conn.ModifyClientVpnEndpointWithContext(ctx, input); err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying EC2 Client VPN Endpoint (%s): %s", d.Id(), err)
		}

		if waitForClientConnectResponseOptionsUpdate {
			if _, err := WaitClientVPNEndpointClientConnectResponseOptionsUpdated(ctx, conn, d.Id()); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for EC2 Client VPN Endpoint (%s) ClientConnectResponseOptions update: %s", d.Id(), err)
			}
		}
	}

	return append(diags, resourceClientVPNEndpointRead(ctx, d, meta)...)
}

func resourceClientVPNEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[DEBUG] Deleting EC2 Client VPN Endpoint: %s", d.Id())
	_, err := conn.DeleteClientVpnEndpointWithContext(ctx, &ec2.DeleteClientVpnEndpointInput{
		ClientVpnEndpointId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidClientVPNEndpointIdNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 Client VPN Endpoint (%s): %s", d.Id(), err)
	}

	if _, err := WaitClientVPNEndpointDeleted(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 Client VPN Endpoint (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func expandClientVPNAuthenticationRequest(tfMap map[string]interface{}) *ec2.ClientVpnAuthenticationRequest {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.ClientVpnAuthenticationRequest{}

	var authnType string
	if v, ok := tfMap["type"].(string); ok && v != "" {
		authnType = v
		apiObject.Type = aws.String(v)
	}

	switch authnType {
	case ec2.ClientVpnAuthenticationTypeCertificateAuthentication:
		if v, ok := tfMap["root_certificate_chain_arn"].(string); ok && v != "" {
			apiObject.MutualAuthentication = &ec2.CertificateAuthenticationRequest{
				ClientRootCertificateChainArn: aws.String(v),
			}
		}

	case ec2.ClientVpnAuthenticationTypeDirectoryServiceAuthentication:
		if v, ok := tfMap["active_directory_id"].(string); ok && v != "" {
			apiObject.ActiveDirectory = &ec2.DirectoryServiceAuthenticationRequest{
				DirectoryId: aws.String(v),
			}
		}

	case ec2.ClientVpnAuthenticationTypeFederatedAuthentication:
		if v, ok := tfMap["saml_provider_arn"].(string); ok && v != "" {
			apiObject.FederatedAuthentication = &ec2.FederatedAuthenticationRequest{
				SAMLProviderArn: aws.String(v),
			}

			if v, ok := tfMap["self_service_saml_provider_arn"].(string); ok && v != "" {
				apiObject.FederatedAuthentication.SelfServiceSAMLProviderArn = aws.String(v)
			}
		}
	}

	return apiObject
}

func expandClientVPNAuthenticationRequests(tfList []interface{}) []*ec2.ClientVpnAuthenticationRequest {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []*ec2.ClientVpnAuthenticationRequest

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandClientVPNAuthenticationRequest(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenClientVPNAuthentication(apiObject *ec2.ClientVpnAuthentication) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Type; v != nil {
		tfMap["type"] = aws.StringValue(v)
	}

	if apiObject.MutualAuthentication != nil {
		if v := apiObject.MutualAuthentication.ClientRootCertificateChain; v != nil {
			tfMap["root_certificate_chain_arn"] = aws.StringValue(v)
		}
	} else if apiObject.ActiveDirectory != nil {
		if v := apiObject.ActiveDirectory.DirectoryId; v != nil {
			tfMap["active_directory_id"] = aws.StringValue(v)
		}
	} else if apiObject.FederatedAuthentication != nil {
		if v := apiObject.FederatedAuthentication.SamlProviderArn; v != nil {
			tfMap["saml_provider_arn"] = aws.StringValue(v)
		}

		if v := apiObject.FederatedAuthentication.SelfServiceSamlProviderArn; v != nil {
			tfMap["self_service_saml_provider_arn"] = aws.StringValue(v)
		}
	}

	return tfMap
}

func flattenClientVPNAuthentications(apiObjects []*ec2.ClientVpnAuthentication) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenClientVPNAuthentication(apiObject))
	}

	return tfList
}

func expandClientConnectOptions(tfMap map[string]interface{}) *ec2.ClientConnectOptions {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.ClientConnectOptions{}

	var enabled bool
	if v, ok := tfMap["enabled"].(bool); ok {
		enabled = v
	}

	if enabled {
		if v, ok := tfMap["lambda_function_arn"].(string); ok && v != "" {
			apiObject.LambdaFunctionArn = aws.String(v)
		}
	}

	apiObject.Enabled = aws.Bool(enabled)

	return apiObject
}

func flattenClientConnectResponseOptions(apiObject *ec2.ClientConnectResponseOptions) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Enabled; v != nil {
		tfMap["enabled"] = aws.BoolValue(v)
	}

	if v := apiObject.LambdaFunctionArn; v != nil {
		tfMap["lambda_function_arn"] = aws.StringValue(v)
	}

	return tfMap
}

func expandClientLoginBannerOptions(tfMap map[string]interface{}) *ec2.ClientLoginBannerOptions {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.ClientLoginBannerOptions{}

	var enabled bool
	if v, ok := tfMap["enabled"].(bool); ok {
		enabled = v
	}

	if enabled {
		if v, ok := tfMap["banner_text"].(string); ok && v != "" {
			apiObject.BannerText = aws.String(v)
		}
	}

	apiObject.Enabled = aws.Bool(enabled)

	return apiObject
}

func flattenClientLoginBannerResponseOptions(apiObject *ec2.ClientLoginBannerResponseOptions) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BannerText; v != nil {
		tfMap["banner_text"] = aws.StringValue(v)
	}

	if v := apiObject.Enabled; v != nil {
		tfMap["enabled"] = aws.BoolValue(v)
	}

	return tfMap
}

func expandConnectionLogOptions(tfMap map[string]interface{}) *ec2.ConnectionLogOptions {
	if tfMap == nil {
		return nil
	}

	apiObject := &ec2.ConnectionLogOptions{}

	var enabled bool
	if v, ok := tfMap["enabled"].(bool); ok {
		enabled = v
	}

	if enabled {
		if v, ok := tfMap["cloudwatch_log_group"].(string); ok && v != "" {
			apiObject.CloudwatchLogGroup = aws.String(v)
		}

		if v, ok := tfMap["cloudwatch_log_stream"].(string); ok && v != "" {
			apiObject.CloudwatchLogStream = aws.String(v)
		}
	}

	apiObject.Enabled = aws.Bool(enabled)

	return apiObject
}

func flattenConnectionLogResponseOptions(apiObject *ec2.ConnectionLogResponseOptions) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.CloudwatchLogGroup; v != nil {
		tfMap["cloudwatch_log_group"] = aws.StringValue(v)
	}

	if v := apiObject.CloudwatchLogStream; v != nil {
		tfMap["cloudwatch_log_stream"] = aws.StringValue(v)
	}

	if v := apiObject.Enabled; v != nil {
		tfMap["enabled"] = aws.BoolValue(v)
	}

	return tfMap
}
