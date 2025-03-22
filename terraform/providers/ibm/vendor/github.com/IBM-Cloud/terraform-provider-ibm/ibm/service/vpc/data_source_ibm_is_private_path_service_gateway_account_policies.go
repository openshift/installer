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

func DataSourceIBMIsPrivatePathServiceGatewayAccountPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesRead,

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private path service gateway identifier.",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with the specified account identifier.",
			},
			"account_policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of account policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this account policy.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	ppsgId := d.Get("private_path_service_gateway").(string)
	listPrivatePathServiceGatewayAccountPoliciesOptions := &vpcv1.ListPrivatePathServiceGatewayAccountPoliciesOptions{}

	listPrivatePathServiceGatewayAccountPoliciesOptions.SetPrivatePathServiceGatewayID(ppsgId)
	if accountIntf, ok := d.GetOk("account"); ok {
		account := accountIntf.(string)
		listPrivatePathServiceGatewayAccountPoliciesOptions.AccountID = &account
	}
	var pager *vpcv1.PrivatePathServiceGatewayAccountPoliciesPager
	pager, err = vpcClient.NewPrivatePathServiceGatewayAccountPoliciesPager(listPrivatePathServiceGatewayAccountPoliciesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] PrivatePathServiceGatewayAccountPoliciesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("PrivatePathServiceGatewayAccountPoliciesPager.GetAll() failed %s", err))
	}

	d.SetId(ppsgId)

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesPrivatePathServiceGatewayAccountPolicyToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("account_policies", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_policies %s", err))
	}

	return nil
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesPrivatePathServiceGatewayAccountPolicyToMap(model *vpcv1.PrivatePathServiceGatewayAccountPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccessPolicy != nil {
		modelMap["access_policy"] = *model.AccessPolicy
	}
	if model.Account != nil {
		accountMap, err := dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = flex.DateTimeToString(model.CreatedAt)
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	// if model.UpdatedAt != nil {
	// 	modelMap["updated_at"] = flex.DateTimeToString(model.UpdatedAt)
	// }
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesPrivatePathServiceGatewayAccountPolicyCollectionFirstToMap(model *vpcv1.PageLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayAccountPoliciesPrivatePathServiceGatewayAccountPolicyCollectionNextToMap(model *vpcv1.PageLink) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}
