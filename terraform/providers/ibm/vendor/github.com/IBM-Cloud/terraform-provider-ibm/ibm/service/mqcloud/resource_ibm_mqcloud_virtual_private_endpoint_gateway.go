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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/mqcloud-go-sdk/mqcloudv1"
)

func ResourceIbmMqcloudVirtualPrivateEndpointGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmMqcloudVirtualPrivateEndpointGatewayCreate,
		ReadContext:   resourceIbmMqcloudVirtualPrivateEndpointGatewayRead,
		DeleteContext: resourceIbmMqcloudVirtualPrivateEndpointGatewayDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"service_instance_guid": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_virtual_private_endpoint_gateway", "service_instance_guid"),
				Description:  "The GUID that uniquely identifies the MQaaS service instance.",
			},
			"trusted_profile": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_virtual_private_endpoint_gateway", "trusted_profile"),
				Description:  "The CRN of the trusted profile to assume for this request.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_virtual_private_endpoint_gateway", "name"),
				Description:  "The name of the virtual private endpoint gateway, created by the user.",
			},
			"target_crn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_mqcloud_virtual_private_endpoint_gateway", "target_crn"),
				Description:  "The CRN of the reserved capacity service instance the user is trying to connect to.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL for the details of the virtual private endpoint gateway.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this virtual privage endpoint.",
			},
			"virtual_private_endpoint_gateway_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the virtual private endpoint gateway which was allocated on creation.",
			},
		},
	}
}

func ResourceIbmMqcloudVirtualPrivateEndpointGatewayValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "service_instance_guid",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "trusted_profile",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\.\/]*$|^crn:\[\.\.\.\]$`,
			MinValueLength:             9,
			MaxValueLength:             512,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z]|[a-z][-a-z0-9]*[a-z0-9]$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "target_crn",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^crn:v[0-9]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]+:[a-z0-9-]*:([a-z]\/[a-z0-9-]+)?:[a-z0-9-]*:[a-z0-9-]*:[a-zA-Z0-9-_\.\/]*$|^crn:\[\.\.\.\]$`,
			MinValueLength:             9,
			MaxValueLength:             512,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_mqcloud_virtual_private_endpoint_gateway", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmMqcloudVirtualPrivateEndpointGatewayCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Create Virtual Private Endpoint Gateway failed: %s", err.Error()), "ibm_mqcloud_virtual_private_endpoint_gateway", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createVirtualPrivateEndpointGatewayOptions := &mqcloudv1.CreateVirtualPrivateEndpointGatewayOptions{}

	createVirtualPrivateEndpointGatewayOptions.SetServiceInstanceGuid(d.Get("service_instance_guid").(string))
	createVirtualPrivateEndpointGatewayOptions.SetName(d.Get("name").(string))
	createVirtualPrivateEndpointGatewayOptions.SetTargetCrn(d.Get("target_crn").(string))
	if _, ok := d.GetOk("trusted_profile"); ok {
		createVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(d.Get("trusted_profile").(string))
	}

	virtualPrivateEndpointGatewayDetails, _, err := mqcloudClient.CreateVirtualPrivateEndpointGatewayWithContext(context, createVirtualPrivateEndpointGatewayOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateVirtualPrivateEndpointGatewayWithContext failed: %s", err.Error()), "ibm_mqcloud_virtual_private_endpoint_gateway", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createVirtualPrivateEndpointGatewayOptions.ServiceInstanceGuid, *virtualPrivateEndpointGatewayDetails.ID))

	return resourceIbmMqcloudVirtualPrivateEndpointGatewayRead(context, d, meta)
}

func resourceIbmMqcloudVirtualPrivateEndpointGatewayRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVirtualPrivateEndpointGatewayOptions := &mqcloudv1.GetVirtualPrivateEndpointGatewayOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "sep-id-parts").GetDiag()
	}

	getVirtualPrivateEndpointGatewayOptions.SetServiceInstanceGuid(parts[0])
	getVirtualPrivateEndpointGatewayOptions.SetVirtualPrivateEndpointGatewayGuid(parts[1])
	if _, ok := d.GetOk("trusted_profile"); ok {
		getVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(d.Get("trusted_profile").(string))
	}

	virtualPrivateEndpointGatewayDetails, response, err := mqcloudClient.GetVirtualPrivateEndpointGatewayWithContext(context, getVirtualPrivateEndpointGatewayOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVirtualPrivateEndpointGatewayWithContext failed: %s", err.Error()), "ibm_mqcloud_virtual_private_endpoint_gateway", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", virtualPrivateEndpointGatewayDetails.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-name").GetDiag()
	}
	if err = d.Set("target_crn", virtualPrivateEndpointGatewayDetails.TargetCrn); err != nil {
		err = fmt.Errorf("Error setting target_crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-target_crn").GetDiag()
	}
	if err = d.Set("href", virtualPrivateEndpointGatewayDetails.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-href").GetDiag()
	}
	if err = d.Set("status", virtualPrivateEndpointGatewayDetails.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-status").GetDiag()
	}
	if err = d.Set("virtual_private_endpoint_gateway_guid", virtualPrivateEndpointGatewayDetails.ID); err != nil {
		err = fmt.Errorf("Error setting virtual_private_endpoint_gateway_guid: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "read", "set-virtual_private_endpoint_gateway_guid").GetDiag()
	}

	return nil
}

func resourceIbmMqcloudVirtualPrivateEndpointGatewayDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mqcloudClient, err := meta.(conns.ClientSession).MqcloudV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	err = checkSIPlan(d, meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Delete Virtual Private Endpoint Gateway failed: %s", err.Error()), "ibm_mqcloud_virtual_private_endpoint_gateway", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteVirtualPrivateEndpointGatewayOptions := &mqcloudv1.DeleteVirtualPrivateEndpointGatewayOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_mqcloud_virtual_private_endpoint_gateway", "delete", "sep-id-parts").GetDiag()
	}

	deleteVirtualPrivateEndpointGatewayOptions.SetServiceInstanceGuid(parts[0])
	deleteVirtualPrivateEndpointGatewayOptions.SetVirtualPrivateEndpointGatewayGuid(parts[1])
	if _, ok := d.GetOk("trusted_profile"); ok {
		deleteVirtualPrivateEndpointGatewayOptions.SetTrustedProfile(d.Get("trusted_profile").(string))
	}

	_, err = mqcloudClient.DeleteVirtualPrivateEndpointGatewayWithContext(context, deleteVirtualPrivateEndpointGatewayOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteVirtualPrivateEndpointGatewayWithContext failed: %s", err.Error()), "ibm_mqcloud_virtual_private_endpoint_gateway", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
