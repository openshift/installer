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

func DataSourceIBMIsIpsecPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsIpsecPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "ipsec_policy"},
				Description:  "The IPsec policy name.",
			},
			"ipsec_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "ipsec_policy"},
				Description:  "The IPsec policy identifier.",
			},
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
			"key_lifetime": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The key lifetime in seconds.",
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
	}
}

func dataSourceIBMIsIpsecPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	identifier := d.Get("ipsec_policy").(string)
	var IPSecPolicy *vpcv1.IPsecPolicy
	if name != "" {
		start := ""
		allrecs := []vpcv1.IPsecPolicy{}
		for {
			listIPSecPoliciesyOptions := &vpcv1.ListIpsecPoliciesOptions{}
			if start != "" {
				listIPSecPoliciesyOptions.Start = &start
			}
			ipSecPolicy, response, err := vpcClient.ListIpsecPolicies(listIPSecPoliciesyOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("Error Fetching IPSec Policies %s\n%s", err, response))
			}
			start = flex.GetNext(ipSecPolicy.Next)
			allrecs = append(allrecs, ipSecPolicy.IpsecPolicies...)
			if start == "" {
				break
			}
		}
		ipsec_policy_found := false
		for _, ipSecPolicyItem := range allrecs {
			if *ipSecPolicyItem.Name == name {
				IPSecPolicy = &ipSecPolicyItem
				ipsec_policy_found = true
				break
			}
		}

		if !ipsec_policy_found {
			log.Printf("[DEBUG] No ipsec policy found with given name %s", name)
			return diag.FromErr(fmt.Errorf("No ipsec policy found with given name %s", name))
		}

	} else {
		getIPSecPolicyOptions := &vpcv1.GetIpsecPolicyOptions{}

		getIPSecPolicyOptions.SetID(identifier)

		ipsecPolicy1, response, err := vpcClient.GetIpsecPolicyWithContext(context, getIPSecPolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] GetIpsecPolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetIpsecPolicyWithContext failed %s\n%s", err, response))
		}
		IPSecPolicy = ipsecPolicy1
	}

	d.SetId(*IPSecPolicy.ID)
	if err = d.Set("authentication_algorithm", IPSecPolicy.AuthenticationAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting authentication_algorithm: %s", err))
	}

	if IPSecPolicy.Connections != nil {
		err = d.Set("connections", dataSourceIPsecPolicyFlattenConnections(IPSecPolicy.Connections))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting connections %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(IPSecPolicy.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("encapsulation_mode", IPSecPolicy.EncapsulationMode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encapsulation_mode: %s", err))
	}
	if err = d.Set("encryption_algorithm", IPSecPolicy.EncryptionAlgorithm); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption_algorithm: %s", err))
	}
	if err = d.Set("href", IPSecPolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("key_lifetime", flex.IntValue(IPSecPolicy.KeyLifetime)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_lifetime: %s", err))
	}
	if err = d.Set("name", IPSecPolicy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("pfs", IPSecPolicy.Pfs); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pfs: %s", err))
	}

	if IPSecPolicy.ResourceGroup != nil {
		err = d.Set("resource_group", dataSourceIPsecPolicyFlattenResourceGroup(*IPSecPolicy.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", IPSecPolicy.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("transform_protocol", IPSecPolicy.TransformProtocol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting transform_protocol: %s", err))
	}

	return nil
}

func dataSourceIPsecPolicyFlattenConnections(result []vpcv1.VPNGatewayConnectionReference) (connections []map[string]interface{}) {
	for _, connectionsItem := range result {
		connections = append(connections, dataSourceIPsecPolicyConnectionsToMap(connectionsItem))
	}

	return connections
}

func dataSourceIPsecPolicyConnectionsToMap(connectionsItem vpcv1.VPNGatewayConnectionReference) (connectionsMap map[string]interface{}) {
	connectionsMap = map[string]interface{}{}

	if connectionsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceIPsecPolicyConnectionsDeletedToMap(*connectionsItem.Deleted)
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

func dataSourceIPsecPolicyConnectionsDeletedToMap(deletedItem vpcv1.VPNGatewayConnectionReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceIPsecPolicyFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceIPsecPolicyResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceIPsecPolicyResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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
