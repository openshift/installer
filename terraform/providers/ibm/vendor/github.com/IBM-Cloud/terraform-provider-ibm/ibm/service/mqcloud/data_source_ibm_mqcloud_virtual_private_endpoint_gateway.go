// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

package mqcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudVirtualPrivateEndpointGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudVirtualPrivateEndpointGatewayRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"virtual_private_endpoint_gateway_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the virtual private endpoint gateway.",
			},
			"trusted_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CRN of the trusted profile to assume for this request.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for the details of the virtual private endpoint gateway.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the virtual private endpoint gateway, created by the user.",
			},
			"target_crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the reserved capacity service instance the user is trying to connect to.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this virtual privage endpoint.",
			},
		},
	}
}

func dataSourceIbmMqcloudVirtualPrivateEndpointGatewayRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVirtualPrivateEndpointGatewayOptions := &mqcloudv1.GetVirtualPrivateEndpointGatewayOptions{}

	getVirtualPrivateEndpointGatewayOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	getVirtualPrivateEndpointGatewayOptions.SetVirtualPrivateEndpointGatewayGuid(d.Get("virtual_private_endpoint_gateway_guid").(string))
	if _, ok := d.GetOk("trusted_profile"); ok {
		getVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(d.Get("trusted_profile").(string))
	}

	virtualPrivateEndpointGatewayDetails, _, err := mqcloudClient.GetVirtualPrivateEndpointGatewayWithContext(context, getVirtualPrivateEndpointGatewayOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVirtualPrivateEndpointGatewayWithContext failed: %s", err.Error()), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getVirtualPrivateEndpointGatewayOptions.ServiceInstanceGuid, *getVirtualPrivateEndpointGatewayOptions.VirtualPrivateEndpointGatewayGuid))

	if err = d.Set("href", virtualPrivateEndpointGatewayDetails.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-href").GetDiag()
	}

	if err = d.Set("name", virtualPrivateEndpointGatewayDetails.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-name").GetDiag()
	}

	if err = d.Set("target_crn", virtualPrivateEndpointGatewayDetails.TargetCrn); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting target_crn: %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-target_crn").GetDiag()
	}

	if err = d.Set("status", virtualPrivateEndpointGatewayDetails.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-status").GetDiag()
	}

	return nil
}
