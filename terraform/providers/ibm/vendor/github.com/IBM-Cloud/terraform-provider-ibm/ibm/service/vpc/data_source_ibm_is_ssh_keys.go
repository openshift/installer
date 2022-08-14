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

const (
	isKeys                 = "keys"
	isKeyCreatedAt         = "created_at"
	isKeyCRN               = "crn"
	isKeysHref             = "href"
	isKeyId                = "id"
	isKeyResourceGroupHref = "href"
	isKeyResourceGroupId   = "id"
	isKeyResourceGroupName = "name"
	isKeysLimit            = "limit"
)

func DataSourceIBMIsSshKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsSshKeysRead,

		Schema: map[string]*schema.Schema{
			isKeys: &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isKeyCreatedAt: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the key was created.",
						},
						isKeyCRN: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this key.",
						},
						isKeyFingerprint: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fingerprint for this key.  The value is returned base64-encoded and prefixed with the hash algorithm (always `SHA256`).",
						},
						isKeysHref: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this key.",
						},
						isKeyId: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this key.",
						},
						isKeyLength: &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The length of this key (in bits).",
						},
						isKeyName: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this key. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						isKeyPublicKey: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public SSH key, consisting of two space-separated fields: the algorithm name, and the base64-encoded key.",
						},
						isKeyResourceGroup: &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this key.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									isKeyResourceGroupHref: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									isKeyResourceGroupId: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									isKeyResourceGroupName: &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						isKeyType: &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The crypto-system used by this key.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsSshKeysRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.Key{}
	listKeysOptions := &vpcv1.ListKeysOptions{}

	for {
		if start != "" {
			listKeysOptions.Start = &start
		}

		keyCollection, response, err := vpcClient.ListKeysWithContext(context, listKeysOptions)
		if err != nil || keyCollection == nil {
			log.Printf("[DEBUG] ListKeysWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListKeysWithContext failed %s\n%s", err, response))
		}

		start = flex.GetNext(keyCollection.Next)
		allrecs = append(allrecs, keyCollection.Keys...)

		if start == "" {
			break
		}

	}

	d.SetId(dataSourceIBMIsSshKeysID(d))

	err = d.Set(isKeys, dataSourceKeyCollectionFlattenKeys(allrecs))
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error setting keys %s", err))
	}

	return nil
}

// dataSourceIBMIsSshKeysID returns a reasonable ID for the list.
func dataSourceIBMIsSshKeysID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceKeyCollectionFlattenKeys(result []vpcv1.Key) (keys []map[string]interface{}) {
	for _, keysItem := range result {
		keys = append(keys, dataSourceKeyCollectionKeysToMap(keysItem))
	}

	return keys
}

func dataSourceKeyCollectionKeysToMap(keysItem vpcv1.Key) (keysMap map[string]interface{}) {
	keysMap = map[string]interface{}{}

	if keysItem.CreatedAt != nil {
		keysMap[isKeyCreatedAt] = keysItem.CreatedAt.String()
	}
	if keysItem.CRN != nil {
		keysMap[isKeyCRN] = keysItem.CRN
	}
	if keysItem.Fingerprint != nil {
		keysMap[isKeyFingerprint] = keysItem.Fingerprint
	}
	if keysItem.Href != nil {
		keysMap[isKeysHref] = keysItem.Href
	}
	if keysItem.ID != nil {
		keysMap[isKeyId] = keysItem.ID
	}
	if keysItem.Length != nil {
		keysMap[isKeyLength] = keysItem.Length
	}
	if keysItem.Name != nil {
		keysMap[isKeyName] = keysItem.Name
	}
	if keysItem.PublicKey != nil {
		keysMap[isKeyPublicKey] = keysItem.PublicKey
	}
	if keysItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceKeyCollectionKeysResourceGroupToMap(*keysItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		keysMap[isKeyResourceGroup] = resourceGroupList
	}
	if keysItem.Type != nil {
		keysMap[isKeyType] = keysItem.Type
	}

	return keysMap
}

func dataSourceKeyCollectionKeysResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap[isKeyResourceGroupHref] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap[isKeyResourceGroupId] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap[isKeyResourceGroupName] = resourceGroupItem.Name
	}

	return resourceGroupMap
}
