// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsIkePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsIkePolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "ike_policy"},
				Description:  "The IKE policy name.",
			},
			"ike_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "ike_policy"},
				Description:  "The IKE policy identifier.",
			},
			"authentication_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The authentication algorithm.",
			},
			"connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPN gateway connections that use this IKE policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN connection's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN gateway connection.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this VPN connection.",
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
				Description: "The date and time that this IKE policy was created.",
			},
			"dh_group": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Diffie-Hellman group.",
			},
			"encryption_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encryption algorithm.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IKE policy's canonical URL.",
			},
			"ike_version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The IKE protocol version.",
			},
			"key_lifetime": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The key lifetime in seconds.",
			},
			"negotiation_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IKE negotiation mode. Only `main` is supported.",
			},
			"resource_group": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this IKE policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIBMIsIkePolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	identifier := d.Get("ike_policy").(string)
	var ikePolicy *vpcv1.IkePolicy
	if name != "" {
		start := ""
		allrecs := []vpcv1.IkePolicy{}
		for {
			listIkePoliciesyOptions := &vpcv1.ListIkePoliciesOptions{}
			if start != "" {
				listIkePoliciesyOptions.Start = &start
			}
			ikePolicies, response, err := vpcClient.ListIkePolicies(listIkePoliciesyOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("Error Fetching IKE Policies %s\n%s", err, response))
			}
			start = flex.GetNext(ikePolicies.Next)
			allrecs = append(allrecs, ikePolicies.IkePolicies...)
			if start == "" {
				break
			}
		}
		ike_policy_found := false
		for _, ikePolicyItem := range allrecs {
			if *ikePolicyItem.Name == name {
				ikePolicy = &ikePolicyItem
				ike_policy_found = true
				break
			}
		}
		if !ike_policy_found {
			log.Printf("[DEBUG] No ike policy found with given name %s", name)
			return diag.FromErr(fmt.Errorf("No ike policy found with given name %s", name))
		}

	} else {
		getIkePolicyOptions := &vpcv1.GetIkePolicyOptions{}

		getIkePolicyOptions.SetID(identifier)

		ikePolicy1, response, err := vpcClient.GetIkePolicyWithContext(context, getIkePolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] GetIkePolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetIkePolicyWithContext failed %s\n%s", err, response))
		}
		ikePolicy = ikePolicy1
	}

	d.SetId(*ikePolicy.ID)
	if err = d.Set("authentication_algorithm", ikePolicy.AuthenticationAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting authentication_algorithm: %s", err))
	}

	if ikePolicy.Connections != nil {
		err = d.Set("connections", dataSourceIkePolicyFlattenConnections(ikePolicy.Connections))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting connections %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(ikePolicy.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("dh_group", flex.IntValue(ikePolicy.DhGroup)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dh_group: %s", err))
	}
	if err = d.Set("encryption_algorithm", ikePolicy.EncryptionAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption_algorithm: %s", err))
	}
	if err = d.Set("href", ikePolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("ike_version", flex.IntValue(ikePolicy.IkeVersion)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ike_version: %s", err))
	}
	if err = d.Set("key_lifetime", flex.IntValue(ikePolicy.KeyLifetime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_lifetime: %s", err))
	}
	if err = d.Set("name", ikePolicy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("negotiation_mode", ikePolicy.NegotiationMode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting negotiation_mode: %s", err))
	}

	if ikePolicy.ResourceGroup != nil {
		err = d.Set("resource_group", dataSourceIkePolicyFlattenResourceGroup(*ikePolicy.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", ikePolicy.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func dataSourceIkePolicyFlattenConnections(result []vpcv1.VPNGatewayConnectionReference) (connections []map[string]interface{}) {
	for _, connectionsItem := range result {
		connections = append(connections, dataSourceIkePolicyConnectionsToMap(connectionsItem))
	}

	return connections
}

func dataSourceIkePolicyConnectionsToMap(connectionsItem vpcv1.VPNGatewayConnectionReference) (connectionsMap map[string]interface{}) {
	connectionsMap = map[string]interface{}{}

	if connectionsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceIkePolicyConnectionsDeletedToMap(*connectionsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		connectionsMap["deleted"] = deletedList
	}
	if connectionsItem.Href != nil {
		connectionsMap["href"] = connectionsItem.Href
	}
	if connectionsItem.ID != nil {
		connectionsMap["id"] = connectionsItem.ID
	}
	if connectionsItem.Name != nil {
		connectionsMap["name"] = connectionsItem.Name
	}
	if connectionsItem.ResourceType != nil {
		connectionsMap["resource_type"] = connectionsItem.ResourceType
	}

	return connectionsMap
}

func dataSourceIkePolicyConnectionsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceIkePolicyFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceIkePolicyResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceIkePolicyResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}
