package ec2

import (
	"context"
	"fmt"
	"log"
	"time"

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

// @SDKResource("aws_vpc_endpoint_service", name="VPC Endpoint Service")
// @Tags(identifierAttribute="id")
func ResourceVPCEndpointService() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVPCEndpointServiceCreate,
		ReadWithoutTimeout:   resourceVPCEndpointServiceRead,
		UpdateWithoutTimeout: resourceVPCEndpointServiceUpdate,
		DeleteWithoutTimeout: resourceVPCEndpointServiceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"acceptance_required": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"allowed_principals": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"base_endpoint_dns_names": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"gateway_load_balancer_arns": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"manages_vpc_endpoints": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"network_load_balancer_arns": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"private_dns_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"private_dns_name_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_ip_address_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(ec2.ServiceConnectivityType_Values(), false),
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceVPCEndpointServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	input := &ec2.CreateVpcEndpointServiceConfigurationInput{
		AcceptanceRequired: aws.Bool(d.Get("acceptance_required").(bool)),
		TagSpecifications:  getTagSpecificationsIn(ctx, ec2.ResourceTypeVpcEndpointService),
	}

	if v, ok := d.GetOk("gateway_load_balancer_arns"); ok && v.(*schema.Set).Len() > 0 {
		input.GatewayLoadBalancerArns = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("network_load_balancer_arns"); ok && v.(*schema.Set).Len() > 0 {
		input.NetworkLoadBalancerArns = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("private_dns_name"); ok {
		input.PrivateDnsName = aws.String(v.(string))
	}

	if v, ok := d.GetOk("supported_ip_address_types"); ok && v.(*schema.Set).Len() > 0 {
		input.SupportedIpAddressTypes = flex.ExpandStringSet(v.(*schema.Set))
	}

	log.Printf("[DEBUG] Creating EC2 VPC Endpoint Service: %s", input)
	output, err := conn.CreateVpcEndpointServiceConfigurationWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating EC2 VPC Endpoint Service: %s", err)
	}

	d.SetId(aws.StringValue(output.ServiceConfiguration.ServiceId))

	if _, err := WaitVPCEndpointServiceAvailable(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPC Endpoint Service (%s) create: %s", d.Id(), err)
	}

	if v, ok := d.GetOk("allowed_principals"); ok && v.(*schema.Set).Len() > 0 {
		input := &ec2.ModifyVpcEndpointServicePermissionsInput{
			AddAllowedPrincipals: flex.ExpandStringSet(v.(*schema.Set)),
			ServiceId:            aws.String(d.Id()),
		}

		if _, err := conn.ModifyVpcEndpointServicePermissionsWithContext(ctx, input); err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying EC2 VPC Endpoint Service (%s) permissions: %s", d.Id(), err)
		}
	}

	return append(diags, resourceVPCEndpointServiceRead(ctx, d, meta)...)
}

func resourceVPCEndpointServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	svcCfg, err := FindVPCEndpointServiceConfigurationByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 VPC Endpoint Service %s not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC Endpoint Service (%s): %s", d.Id(), err)
	}

	d.Set("acceptance_required", svcCfg.AcceptanceRequired)
	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Service:   ec2.ServiceName,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("vpc-endpoint-service/%s", d.Id()),
	}.String()
	d.Set("arn", arn)
	d.Set("availability_zones", aws.StringValueSlice(svcCfg.AvailabilityZones))
	d.Set("base_endpoint_dns_names", aws.StringValueSlice(svcCfg.BaseEndpointDnsNames))
	d.Set("gateway_load_balancer_arns", aws.StringValueSlice(svcCfg.GatewayLoadBalancerArns))
	d.Set("manages_vpc_endpoints", svcCfg.ManagesVpcEndpoints)
	d.Set("network_load_balancer_arns", aws.StringValueSlice(svcCfg.NetworkLoadBalancerArns))
	d.Set("private_dns_name", svcCfg.PrivateDnsName)
	// The EC2 API can return a XML structure with no elements.
	if tfMap := flattenPrivateDNSNameConfiguration(svcCfg.PrivateDnsNameConfiguration); len(tfMap) > 0 {
		if err := d.Set("private_dns_name_configuration", []interface{}{tfMap}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting private_dns_name_configuration: %s", err)
		}
	} else {
		d.Set("private_dns_name_configuration", nil)
	}
	d.Set("service_name", svcCfg.ServiceName)
	if len(svcCfg.ServiceType) > 0 {
		d.Set("service_type", svcCfg.ServiceType[0].ServiceType)
	} else {
		d.Set("service_type", nil)
	}
	d.Set("state", svcCfg.ServiceState)
	d.Set("supported_ip_address_types", aws.StringValueSlice(svcCfg.SupportedIpAddressTypes))

	SetTagsOut(ctx, svcCfg.Tags)

	allowedPrincipals, err := FindVPCEndpointServicePermissionsByServiceID(ctx, conn, d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 VPC Endpoint Service (%s) permissions: %s", d.Id(), err)
	}

	d.Set("allowed_principals", flattenAllowedPrincipals(allowedPrincipals))

	return diags
}

func resourceVPCEndpointServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	if d.HasChanges("acceptance_required", "gateway_load_balancer_arns", "network_load_balancer_arns", "private_dns_name", "supported_ip_address_types") {
		input := &ec2.ModifyVpcEndpointServiceConfigurationInput{
			ServiceId: aws.String(d.Id()),
		}

		if d.HasChange("acceptance_required") {
			input.AcceptanceRequired = aws.Bool(d.Get("acceptance_required").(bool))
		}

		input.AddGatewayLoadBalancerArns, input.RemoveGatewayLoadBalancerArns = flattenAddAndRemoveStringLists(d, "gateway_load_balancer_arns")
		input.AddNetworkLoadBalancerArns, input.RemoveNetworkLoadBalancerArns = flattenAddAndRemoveStringLists(d, "network_load_balancer_arns")

		if d.HasChange("private_dns_name") {
			input.PrivateDnsName = aws.String(d.Get("private_dns_name").(string))
		}

		input.AddSupportedIpAddressTypes, input.RemoveSupportedIpAddressTypes = flattenAddAndRemoveStringLists(d, "supported_ip_address_types")

		log.Printf("[DEBUG] Updating EC2 VPC Endpoint Service: %s", input)
		_, err := conn.ModifyVpcEndpointServiceConfigurationWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating EC2 VPC Endpoint Service (%s): %s", d.Id(), err)
		}

		if _, err := WaitVPCEndpointServiceAvailable(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPC Endpoint Service (%s) update: %s", d.Id(), err)
		}
	}

	if d.HasChange("allowed_principals") {
		input := &ec2.ModifyVpcEndpointServicePermissionsInput{
			ServiceId: aws.String(d.Id()),
		}

		input.AddAllowedPrincipals, input.RemoveAllowedPrincipals = flattenAddAndRemoveStringLists(d, "allowed_principals")

		if _, err := conn.ModifyVpcEndpointServicePermissionsWithContext(ctx, input); err != nil {
			return sdkdiag.AppendErrorf(diags, "modifying EC2 VPC Endpoint Service (%s) permissions: %s", d.Id(), err)
		}
	}

	return append(diags, resourceVPCEndpointServiceRead(ctx, d, meta)...)
}

func resourceVPCEndpointServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Conn(ctx)

	log.Printf("[INFO] Deleting EC2 VPC Endpoint Service: %s", d.Id())
	output, err := conn.DeleteVpcEndpointServiceConfigurationsWithContext(ctx, &ec2.DeleteVpcEndpointServiceConfigurationsInput{
		ServiceIds: aws.StringSlice([]string{d.Id()}),
	})

	if err == nil && output != nil {
		err = UnsuccessfulItemsError(output.Unsuccessful)
	}

	if tfawserr.ErrCodeEquals(err, errCodeInvalidVPCEndpointServiceIdNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting EC2 VPC Endpoint Service (%s): %s", d.Id(), err)
	}

	if _, err := WaitVPCEndpointServiceDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for EC2 VPC Endpoint Service (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func flattenAllowedPrincipal(apiObject *ec2.AllowedPrincipal) *string {
	if apiObject == nil {
		return nil
	}

	return apiObject.Principal
}

func flattenAllowedPrincipals(apiObjects []*ec2.AllowedPrincipal) []*string {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []*string

	for _, apiObject := range apiObjects {
		if apiObject == nil {
			continue
		}

		tfList = append(tfList, flattenAllowedPrincipal(apiObject))
	}

	return tfList
}

func flattenPrivateDNSNameConfiguration(apiObject *ec2.PrivateDnsNameConfiguration) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Name; v != nil {
		tfMap["name"] = aws.StringValue(v)
	}

	if v := apiObject.State; v != nil {
		tfMap["state"] = aws.StringValue(v)
	}

	if v := apiObject.Type; v != nil {
		tfMap["type"] = aws.StringValue(v)
	}

	if v := apiObject.Value; v != nil {
		tfMap["value"] = aws.StringValue(v)
	}

	return tfMap
}
