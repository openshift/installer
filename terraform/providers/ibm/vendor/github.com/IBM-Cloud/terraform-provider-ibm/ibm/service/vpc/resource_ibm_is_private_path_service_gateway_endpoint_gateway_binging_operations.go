// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	AccessPolicyEnumPermit = "permit"
	AccessPolicyEnumDeny   = "deny"
)

func ResourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsCreate,
		ReadContext:   resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsRead,
		UpdateContext: resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsUpdate,
		DeleteContext: resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private path service gateway identifier.",
			},
			"endpoint_gateway_binding": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private path service gateway identifier.",
			},
			"access_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access polict to set for this endpoint gateway binding.",
			},
		},
	}
}

func resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	ppsgId := d.Get("private_path_service_gateway").(string)
	egwbindingId := d.Get("endpoint_gateway_binding").(string)
	accessPolicy := d.Get("access_policy").(string)
	if accessPolicy == AccessPolicyEnumPermit {
		permitPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.PermitPrivatePathServiceGatewayEndpointGatewayBindingOptions{}

		permitPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetPrivatePathServiceGatewayID(ppsgId)
		permitPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetID(egwbindingId)

		response, err := vpcClient.PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext(context, permitPrivatePathServiceGatewayEndpointGatewayBindingOptions)
		if err != nil {
			log.Printf("[DEBUG] PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response))
		}
	} else {
		denyPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.DenyPrivatePathServiceGatewayEndpointGatewayBindingOptions{}

		denyPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetPrivatePathServiceGatewayID(ppsgId)
		denyPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetID(egwbindingId)

		response, err := vpcClient.DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext(context, denyPrivatePathServiceGatewayEndpointGatewayBindingOptions)
		if err != nil {
			log.Printf("[DEBUG] DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response))
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", ppsgId, egwbindingId))

	return resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	ppsgId := d.Get("private_path_service_gateway").(string)
	egwbindingId := d.Get("endpoint_gateway_binding").(string)
	if d.HasChange("access_policy") {
		_, newAccessPolicy := d.GetChange("access_policy")
		accessPolicy := newAccessPolicy.(string)
		if accessPolicy == AccessPolicyEnumPermit {
			permitPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.PermitPrivatePathServiceGatewayEndpointGatewayBindingOptions{}

			permitPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetPrivatePathServiceGatewayID(ppsgId)
			permitPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetID(egwbindingId)

			response, err := vpcClient.PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext(context, permitPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			if err != nil {
				log.Printf("[DEBUG] PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("PermitPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response))
			}
		} else {
			denyPrivatePathServiceGatewayEndpointGatewayBindingOptions := &vpcv1.DenyPrivatePathServiceGatewayEndpointGatewayBindingOptions{}

			denyPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetPrivatePathServiceGatewayID(ppsgId)
			denyPrivatePathServiceGatewayEndpointGatewayBindingOptions.SetID(egwbindingId)

			response, err := vpcClient.DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext(context, denyPrivatePathServiceGatewayEndpointGatewayBindingOptions)
			if err != nil {
				log.Printf("[DEBUG] DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("DenyPrivatePathServiceGatewayEndpointGatewayBindingWithContext failed %s\n%s", err, response))
			}
		}

	}

	return resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingOperationsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId("")

	return nil
}
