// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsIkePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsIkePoliciesRead,

		Schema: map[string]*schema.Schema{
			"ike_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of IKE policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this IKE policy.",
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
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this IKE policy.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsIkePoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.IkePolicy{}
	for {
		listIkePoliciesOptions := &vpcv1.ListIkePoliciesOptions{}
		if start != "" {
			listIkePoliciesOptions.Start = &start
		}
		ikePolicyCollection, response, err := vpcClient.ListIkePoliciesWithContext(context, listIkePoliciesOptions)
		if err != nil || ikePolicyCollection == nil {
			log.Printf("[DEBUG] ListIkePoliciesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListIkePoliciesWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(ikePolicyCollection.Next)
		allrecs = append(allrecs, ikePolicyCollection.IkePolicies...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsIkePoliciesID(d))

	err = d.Set("ike_policies", dataSourceIkePolicyCollectionFlattenIkePolicies(allrecs))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ike_policies %s", err))
	}

	return nil
}

// dataSourceIBMIsIkePoliciesID returns a reasonable ID for the list.
func dataSourceIBMIsIkePoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIkePolicyCollectionFlattenIkePolicies(result []vpcv1.IkePolicy) (ikePolicies []map[string]interface{}) {
	for _, ikePoliciesItem := range result {
		ikePolicies = append(ikePolicies, dataSourceIkePolicyCollectionIkePoliciesToMap(ikePoliciesItem))
	}

	return ikePolicies
}

func dataSourceIkePolicyCollectionIkePoliciesToMap(ikePoliciesItem vpcv1.IkePolicy) (ikePoliciesMap map[string]interface{}) {
	ikePoliciesMap = map[string]interface{}{}

	if ikePoliciesItem.AuthenticationAlgorithm != nil {
		ikePoliciesMap["authentication_algorithm"] = ikePoliciesItem.AuthenticationAlgorithm
	}
	if ikePoliciesItem.Connections != nil {
		connectionsList := []map[string]interface{}{}
		for _, connectionsItem := range ikePoliciesItem.Connections {
			connectionsList = append(connectionsList, dataSourceIkePolicyCollectionIkePoliciesConnectionsToMap(connectionsItem))
		}
		ikePoliciesMap["connections"] = connectionsList
	}
	if ikePoliciesItem.CreatedAt != nil {
		ikePoliciesMap["created_at"] = ikePoliciesItem.CreatedAt.String()
	}
	if ikePoliciesItem.DhGroup != nil {
		ikePoliciesMap["dh_group"] = ikePoliciesItem.DhGroup
	}
	if ikePoliciesItem.EncryptionAlgorithm != nil {
		ikePoliciesMap["encryption_algorithm"] = ikePoliciesItem.EncryptionAlgorithm
	}
	if ikePoliciesItem.Href != nil {
		ikePoliciesMap["href"] = ikePoliciesItem.Href
	}
	if ikePoliciesItem.ID != nil {
		ikePoliciesMap["id"] = ikePoliciesItem.ID
	}
	if ikePoliciesItem.IkeVersion != nil {
		ikePoliciesMap["ike_version"] = ikePoliciesItem.IkeVersion
	}
	if ikePoliciesItem.KeyLifetime != nil {
		ikePoliciesMap["key_lifetime"] = ikePoliciesItem.KeyLifetime
	}
	if ikePoliciesItem.Name != nil {
		ikePoliciesMap["name"] = ikePoliciesItem.Name
	}
	if ikePoliciesItem.NegotiationMode != nil {
		ikePoliciesMap["negotiation_mode"] = ikePoliciesItem.NegotiationMode
	}
	if ikePoliciesItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceIkePolicyCollectionIkePoliciesResourceGroupToMap(*ikePoliciesItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		ikePoliciesMap["resource_group"] = resourceGroupList
	}
	if ikePoliciesItem.ResourceType != nil {
		ikePoliciesMap["resource_type"] = ikePoliciesItem.ResourceType
	}

	return ikePoliciesMap
}

func dataSourceIkePolicyCollectionIkePoliciesConnectionsToMap(connectionsItem vpcv1.VPNGatewayConnectionReference) (connectionsMap map[string]interface{}) {
	connectionsMap = map[string]interface{}{}

	if connectionsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceIkePolicyCollectionConnectionsDeletedToMap(*connectionsItem.Deleted)
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

func dataSourceIkePolicyCollectionConnectionsDeletedToMap(deletedItem vpcv1.VPNGatewayConnectionReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceIkePolicyCollectionIkePoliciesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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
