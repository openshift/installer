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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsPrivatePathServiceGatewayAccountPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsPrivatePathServiceGatewayAccountPolicyCreate,
		ReadContext:   resourceIBMIsPrivatePathServiceGatewayAccountPolicyRead,
		UpdateContext: resourceIBMIsPrivatePathServiceGatewayAccountPolicyUpdate,
		DeleteContext: resourceIBMIsPrivatePathServiceGatewayAccountPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private path service gateway identifier.",
			},
			"access_policy": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_private_path_service_gateway_account_policy", "access_policy"),
				Description:  "The access policy for the account:- permit: access will be permitted- deny:  access will be denied- review: access will be manually reviewed.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The account for this access policy.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the account policy was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this account policy.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the account policy was updated.",
			},
			"account_policy": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this account policy.",
			},
		},
	}
}

func ResourceIBMIsPrivatePathServiceGatewayAccountPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "access_policy",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "deny, permit, review",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_private_path_service_gateway_account_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsPrivatePathServiceGatewayAccountPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createPrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.CreatePrivatePathServiceGatewayAccountPolicyOptions{}

	createPrivatePathServiceGatewayAccountPolicyOptions.SetPrivatePathServiceGatewayID(d.Get("private_path_service_gateway").(string))
	createPrivatePathServiceGatewayAccountPolicyOptions.SetAccessPolicy(d.Get("access_policy").(string))
	accountId := d.Get("account").(string)
	account := &vpcv1.AccountIdentity{
		ID: &accountId,
	}
	createPrivatePathServiceGatewayAccountPolicyOptions.SetAccount(account)

	privatePathServiceGatewayAccountPolicy, response, err := vpcClient.CreatePrivatePathServiceGatewayAccountPolicyWithContext(context, createPrivatePathServiceGatewayAccountPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] CreatePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreatePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createPrivatePathServiceGatewayAccountPolicyOptions.PrivatePathServiceGatewayID, *privatePathServiceGatewayAccountPolicy.ID))

	return resourceIBMIsPrivatePathServiceGatewayAccountPolicyRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayAccountPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.GetPrivatePathServiceGatewayAccountPolicyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getPrivatePathServiceGatewayAccountPolicyOptions.SetPrivatePathServiceGatewayID(parts[0])
	getPrivatePathServiceGatewayAccountPolicyOptions.SetID(parts[1])

	privatePathServiceGatewayAccountPolicy, response, err := vpcClient.GetPrivatePathServiceGatewayAccountPolicyWithContext(context, getPrivatePathServiceGatewayAccountPolicyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetPrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("access_policy", privatePathServiceGatewayAccountPolicy.AccessPolicy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting access_policy: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(privatePathServiceGatewayAccountPolicy.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", privatePathServiceGatewayAccountPolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("resource_type", privatePathServiceGatewayAccountPolicy.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	// if err = d.Set("updated_at", flex.DateTimeToString(privatePathServiceGatewayAccountPolicy.UpdatedAt)); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	// }
	if err = d.Set("account_policy", privatePathServiceGatewayAccountPolicy.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_path_service_gateway_account_policy_id: %s", err))
	}

	return nil
}

func resourceIBMIsPrivatePathServiceGatewayAccountPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updatePrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.UpdatePrivatePathServiceGatewayAccountPolicyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updatePrivatePathServiceGatewayAccountPolicyOptions.SetPrivatePathServiceGatewayID(parts[0])
	updatePrivatePathServiceGatewayAccountPolicyOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.PrivatePathServiceGatewayAccountPolicyPatch{}

	if d.HasChange("access_policy") {
		newAccessPolicy := d.Get("access_policy").(string)
		patchVals.AccessPolicy = &newAccessPolicy
		hasChange = true
	}

	if hasChange {
		updatePrivatePathServiceGatewayAccountPolicyOptions.PrivatePathServiceGatewayAccountPolicyPatch, _ = patchVals.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] Error calling AsPatch for PrivatePathServiceGatewayAccountPolicyPatch %s", err)
			return diag.FromErr(err)
		}
		_, response, err := vpcClient.UpdatePrivatePathServiceGatewayAccountPolicyWithContext(context, updatePrivatePathServiceGatewayAccountPolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdatePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdatePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsPrivatePathServiceGatewayAccountPolicyRead(context, d, meta)
}

func resourceIBMIsPrivatePathServiceGatewayAccountPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deletePrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.DeletePrivatePathServiceGatewayAccountPolicyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deletePrivatePathServiceGatewayAccountPolicyOptions.SetPrivatePathServiceGatewayID(parts[0])
	deletePrivatePathServiceGatewayAccountPolicyOptions.SetID(parts[1])

	response, err := vpcClient.DeletePrivatePathServiceGatewayAccountPolicyWithContext(context, deletePrivatePathServiceGatewayAccountPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeletePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeletePrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
