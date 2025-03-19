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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsPrivatePathServiceGatewayAccountPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPrivatePathServiceGatewayAccountPolicyRead,

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private path service gateway identifier.",
			},
			"account_policy": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The account policy identifier.",
			},
			"access_policy": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access policy for the account:- permit: access will be permitted- deny:  access will be denied- review: access will be manually reviewedThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The account for this access policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
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
		},
	}
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPrivatePathServiceGatewayAccountPolicyOptions := &vpcv1.GetPrivatePathServiceGatewayAccountPolicyOptions{}

	getPrivatePathServiceGatewayAccountPolicyOptions.SetPrivatePathServiceGatewayID(d.Get("private_path_service_gateway").(string))
	getPrivatePathServiceGatewayAccountPolicyOptions.SetID(d.Get("account_policy").(string))

	privatePathServiceGatewayAccountPolicy, response, err := vpcClient.GetPrivatePathServiceGatewayAccountPolicyWithContext(context, getPrivatePathServiceGatewayAccountPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetPrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPrivatePathServiceGatewayAccountPolicyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getPrivatePathServiceGatewayAccountPolicyOptions.PrivatePathServiceGatewayID, *getPrivatePathServiceGatewayAccountPolicyOptions.ID))

	if err = d.Set("access_policy", privatePathServiceGatewayAccountPolicy.AccessPolicy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting access_policy: %s", err))
	}

	account := []map[string]interface{}{}
	if privatePathServiceGatewayAccountPolicy.Account != nil {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayAccountPolicyAccountReferenceToMap(privatePathServiceGatewayAccountPolicy.Account)
		if err != nil {
			return diag.FromErr(err)
		}
		account = append(account, modelMap)
	}
	if err = d.Set("account", account); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account %s", err))
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

	return nil
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPolicyAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}
