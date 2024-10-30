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

func DataSourceIBMIsIpsecPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsIpsecPoliciesRead,

		Schema: map[string]*schema.Schema{
			"ipsec_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of IPsec policies.",
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
							Description: "The VPN gateway connections that use this IPsec policy.",
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
							Description: "The date and time that this IPsec policy was created.",
						},
						"encapsulation_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The encapsulation mode used. Only `tunnel` is supported.",
						},
						"encryption_algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The encryption algorithm.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPsec policy's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this IPsec policy.",
						},
						"key_lifetime": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The key lifetime in seconds.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this IPsec policy.",
						},
						"pfs": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Perfect Forward Secrecy.",
						},
						"resource_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this IPsec policy.",
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
						"transform_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The transform protocol used. Only `esp` is supported.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsIpsecPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.IPsecPolicy{}
	for {
		listIpsecPoliciesOptions := &vpcv1.ListIpsecPoliciesOptions{}
		if start != "" {
			listIpsecPoliciesOptions.Start = &start
		}
		iPsecPolicyCollection, response, err := vpcClient.ListIpsecPoliciesWithContext(context, listIpsecPoliciesOptions)
		if err != nil || iPsecPolicyCollection == nil {
			log.Printf("[DEBUG] ListIpsecPoliciesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListIpsecPoliciesWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(iPsecPolicyCollection.Next)
		allrecs = append(allrecs, iPsecPolicyCollection.IpsecPolicies...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsIpsecPoliciesID(d))

	err = d.Set("ipsec_policies", dataSourceIPsecPolicyCollectionFlattenIpsecPolicies(allrecs))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ipsec_policies %s", err))
	}

	return nil
}

// dataSourceIBMIsIpsecPoliciesID returns a reasonable ID for the list.
func dataSourceIBMIsIpsecPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIPsecPolicyCollectionFlattenIpsecPolicies(result []vpcv1.IPsecPolicy) (ipsecPolicies []map[string]interface{}) {
	for _, ipsecPoliciesItem := range result {
		ipsecPolicies = append(ipsecPolicies, dataSourceIPsecPolicyCollectionIpsecPoliciesToMap(ipsecPoliciesItem))
	}

	return ipsecPolicies
}

func dataSourceIPsecPolicyCollectionIpsecPoliciesToMap(ipsecPoliciesItem vpcv1.IPsecPolicy) (ipsecPoliciesMap map[string]interface{}) {
	ipsecPoliciesMap = map[string]interface{}{}

	if ipsecPoliciesItem.AuthenticationAlgorithm != nil {
		ipsecPoliciesMap["authentication_algorithm"] = ipsecPoliciesItem.AuthenticationAlgorithm
	}
	if ipsecPoliciesItem.Connections != nil {
		connectionsList := []map[string]interface{}{}
		for _, connectionsItem := range ipsecPoliciesItem.Connections {
			connectionsList = append(connectionsList, dataSourceIPsecPolicyCollectionIpsecPoliciesConnectionsToMap(connectionsItem))
		}
		ipsecPoliciesMap["connections"] = connectionsList
	}
	if ipsecPoliciesItem.CreatedAt != nil {
		ipsecPoliciesMap["created_at"] = ipsecPoliciesItem.CreatedAt.String()
	}
	if ipsecPoliciesItem.EncapsulationMode != nil {
		ipsecPoliciesMap["encapsulation_mode"] = ipsecPoliciesItem.EncapsulationMode
	}
	if ipsecPoliciesItem.EncryptionAlgorithm != nil {
		ipsecPoliciesMap["encryption_algorithm"] = ipsecPoliciesItem.EncryptionAlgorithm
	}
	if ipsecPoliciesItem.Href != nil {
		ipsecPoliciesMap["href"] = ipsecPoliciesItem.Href
	}
	if ipsecPoliciesItem.ID != nil {
		ipsecPoliciesMap["id"] = ipsecPoliciesItem.ID
	}
	if ipsecPoliciesItem.KeyLifetime != nil {
		ipsecPoliciesMap["key_lifetime"] = ipsecPoliciesItem.KeyLifetime
	}
	if ipsecPoliciesItem.Name != nil {
		ipsecPoliciesMap["name"] = ipsecPoliciesItem.Name
	}
	if ipsecPoliciesItem.Pfs != nil {
		ipsecPoliciesMap["pfs"] = ipsecPoliciesItem.Pfs
	}
	if ipsecPoliciesItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceIPsecPolicyCollectionIpsecPoliciesResourceGroupToMap(*ipsecPoliciesItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		ipsecPoliciesMap["resource_group"] = resourceGroupList
	}
	if ipsecPoliciesItem.ResourceType != nil {
		ipsecPoliciesMap["resource_type"] = ipsecPoliciesItem.ResourceType
	}
	if ipsecPoliciesItem.TransformProtocol != nil {
		ipsecPoliciesMap["transform_protocol"] = ipsecPoliciesItem.TransformProtocol
	}

	return ipsecPoliciesMap
}

func dataSourceIPsecPolicyCollectionIpsecPoliciesConnectionsToMap(connectionsItem vpcv1.VPNGatewayConnectionReference) (connectionsMap map[string]interface{}) {
	connectionsMap = map[string]interface{}{}

	if connectionsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceIPsecPolicyCollectionConnectionsDeletedToMap(*connectionsItem.Deleted)
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

func dataSourceIPsecPolicyCollectionConnectionsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceIPsecPolicyCollectionIpsecPoliciesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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
