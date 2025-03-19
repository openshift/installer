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

func DataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsRead,

		Schema: map[string]*schema.Schema{
			"private_path_service_gateway": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The private path service gateway identifier.",
			},
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with the specified account identifier.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with the specified status.",
			},
			"endpoint_gateway_bindings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of endpoint gateway bindings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The account that created the endpoint gateway.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the endpoint gateway binding was created.",
						},
						"expiration_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration date and time for the endpoint gateway binding.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this endpoint gateway binding.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this endpoint gateway binding.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the endpoint gateway binding.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the endpoint gateway binding- `denied`: endpoint gateway binding was denied- `expired`: endpoint gateway binding has expired- `pending`: endpoint gateway binding is awaiting review- `permitted`: endpoint gateway binding was permittedThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the endpoint gateway binding was updated.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listPrivatePathServiceGatewayEndpointGatewayBindingsOptions := &vpcv1.ListPrivatePathServiceGatewayEndpointGatewayBindingsOptions{}

	listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.SetPrivatePathServiceGatewayID(d.Get("private_path_service_gateway").(string))

	if accountIntf, ok := d.GetOk("account"); ok {
		account := accountIntf.(string)
		listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.AccountID = &account
	}

	if statusIntf, ok := d.GetOk("status"); ok {
		status := statusIntf.(string)
		listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.Status = &status
	}

	var pager *vpcv1.PrivatePathServiceGatewayEndpointGatewayBindingsPager
	pager, err = vpcClient.NewPrivatePathServiceGatewayEndpointGatewayBindingsPager(listPrivatePathServiceGatewayEndpointGatewayBindingsOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] PrivatePathServiceGatewayEndpointGatewayBindingsPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("PrivatePathServiceGatewayEndpointGatewayBindingsPager.GetAll() failed %s", err))
	}

	d.SetId(*listPrivatePathServiceGatewayEndpointGatewayBindingsOptions.PrivatePathServiceGatewayID)

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsPrivatePathServiceGatewayEndpointGatewayBindingToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("endpoint_gateway_bindings", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoint_gateway_bindings %s", err))
	}

	return nil
}

func dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsPrivatePathServiceGatewayEndpointGatewayBindingToMap(model *vpcv1.PrivatePathServiceGatewayEndpointGatewayBinding) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.ExpirationAt != nil {
		modelMap["expiration_at"] = model.ExpirationAt.String()
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LifecycleState != nil {
		modelMap["lifecycle_state"] = *model.LifecycleState
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	// if model.UpdatedAt != nil {
	// 	modelMap["updated_at"] = model.UpdatedAt.String()
	// }
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsPrivatePathServiceGatewayEndpointGatewayBindingCollectionFirstToMap(model *vpcv1.PrivatePathServiceGatewayEndpointGatewayBindingCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func dataSourceIBMIsPrivatePathServiceGatewayEndpointGatewayBindingsPrivatePathServiceGatewayEndpointGatewayBindingCollectionNextToMap(model *vpcv1.PrivatePathServiceGatewayEndpointGatewayBindingCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}
