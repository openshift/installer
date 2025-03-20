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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func DataSourceIbmMqcloudVirtualPrivateEndpointGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmMqcloudVirtualPrivateEndpointGatewaysRead,

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"trusted_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CRN of the trusted profile to assume for this request.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the virtual private endpoint gateway, created by the user.",
			},
			"virtual_private_endpoint_gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of virtual private endpoint gateways.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URL for the details of the virtual private endpoint gateway.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the virtual private endpoint gateway which was allocated on creation.",
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
				},
			},
		},
	}
}

func dataSourceIbmMqcloudVirtualPrivateEndpointGatewaysRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_virtual_private_endpoint_gateways", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listVirtualPrivateEndpointGatewaysOptions := &mqcloudv1.ListVirtualPrivateEndpointGatewaysOptions{}

	listVirtualPrivateEndpointGatewaysOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	if _, ok := d.GetOk("trusted_profile"); ok {
		listVirtualPrivateEndpointGatewaysOptions.SetTrustedProfile(d.Get("trusted_profile").(string))
	}

	var pager *mqcloudv1.VirtualPrivateEndpointGatewaysPager
	pager, err = mqcloudClient.NewVirtualPrivateEndpointGatewaysPager(listVirtualPrivateEndpointGatewaysOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_virtual_private_endpoint_gateways", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("VirtualPrivateEndpointGatewaysPager.GetAll() failed %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateways", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	var name string

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		d.SetId(name)
	} else {
		d.SetId(dataSourceIbmMqcloudVirtualPrivateEndpointGatewaysID(d))
	}

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIbmMqcloudVirtualPrivateEndpointGatewaysVirtualPrivateEndpointGatewayDetailsToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_mqcloud_virtual_private_endpoint_gateways", "read", "VirtualPrivateEndpointsGateways-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("virtual_private_endpoint_gateways", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting virtual_private_endpoint_gateways %s", err), "(Data) ibm_mqcloud_virtual_private_endpoint_gateways", "read", "virtual_private_endpoint_gateways-set").GetDiag()
	}

	return nil
}

// dataSourceIbmMqcloudVirtualPrivateEndpointGatewaysID returns a reasonable ID for the list.
func dataSourceIbmMqcloudVirtualPrivateEndpointGatewaysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmMqcloudVirtualPrivateEndpointGatewaysVirtualPrivateEndpointGatewayDetailsToMap(model *mqcloudv1.VirtualPrivateEndpointGatewayDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["target_crn"] = *model.TargetCrn
	modelMap["status"] = *model.Status
	return modelMap, nil
}
