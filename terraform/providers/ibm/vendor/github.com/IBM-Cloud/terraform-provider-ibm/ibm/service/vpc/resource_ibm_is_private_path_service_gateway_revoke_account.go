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

func ResourceIBMIsPrivatePathServiceGatewayRevokeAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsPrivatePathServiceGatewayRevokeAccountCreate,
		ReadContext:   resourceIBMIsPrivatePathServiceGatewayRevokeAccountRead,
		UpdateContext: resourceIBMIsPrivatePathServiceGatewayRevokeAccountUpdate,
		DeleteContext: resourceIBMIsPrivatePathServiceGatewayRevokeAccountDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private path service gateway identifier.",
			},
			"account": {
				Type:     schema.TypeString,
				Required: true,
				//ForceNew:    true,
				Description: "The account for this access policy.",
			},
		},
	}
}

func resourceIBMIsPrivatePathServiceGatewayRevokeAccountCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	revokePrivatePathServiceGatewayOptions := &vpcv1.RevokeAccountForPrivatePathServiceGatewayOptions{}

	revokePrivatePathServiceGatewayOptions.SetPrivatePathServiceGatewayID(d.Get("private_path_service_gateway").(string))

	accountId := d.Get("account").(string)
	account := &vpcv1.AccountIdentity{
		ID: &accountId,
	}
	revokePrivatePathServiceGatewayOptions.SetAccount(account)

	response, err := vpcClient.RevokeAccountForPrivatePathServiceGatewayWithContext(context, revokePrivatePathServiceGatewayOptions)
	if err != nil {
		log.Printf("[DEBUG] RevokeAccountForPrivatePathServiceGatewayWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("RevokeAccountForPrivatePathServiceGatewayWithContext failed %s\n%s", err, response))
	}

	d.SetId(*revokePrivatePathServiceGatewayOptions.PrivatePathServiceGatewayID)

	return resourceIBMIsPrivatePathServiceGatewayRevokeAccountRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayRevokeAccountRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIBMIsPrivatePathServiceGatewayRevokeAccountUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return resourceIBMIsPrivatePathServiceGatewayRevokeAccountRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayRevokeAccountDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId("")

	return nil
}
