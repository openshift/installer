// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsShareProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareProfilesRead,

		Schema: map[string]*schema.Schema{
			"profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of share profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted capacity range (in gigabytes) for a share with this profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default capacity.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The max capacity.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The min capacity.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field",
									},
									"values": {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "The permitted values for this profile field.",
									},
								},
							},
						},
						"iops": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted IOPS range for a share with this profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The default iops.",
									},
									"max": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The max iops.",
									},
									"min": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The min iops.",
									},
									"step": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The increment step value for this profile field.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type for this profile field.",
									},
									"value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The value for this profile field",
									},
									"values": {
										Type:        schema.TypeSet,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Description: "The permitted values for this profile field.",
									},
								},
							},
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this share profile belongs to.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this share profile.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsShareProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listShareProfilesOptions := &vpcv1.ListShareProfilesOptions{}

	shareProfileCollection, response, err := vpcClient.ListShareProfilesWithContext(context, listShareProfilesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListShareProfilesWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(dataSourceIbmIsShareProfilesID(d))

	if shareProfileCollection.Profiles != nil {
		err = d.Set("profiles", dataSourceShareProfileCollectionFlattenProfiles(shareProfileCollection.Profiles))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profiles %s", err))
		}
	}
	if err = d.Set("total_count", shareProfileCollection.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIbmIsShareProfilesID returns a reasonable ID for the list.
func dataSourceIbmIsShareProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareProfileCollectionFlattenProfiles(result []vpcv1.ShareProfile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceShareProfileCollectionProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceShareProfileCollectionProfilesToMap(profilesItem vpcv1.ShareProfile) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.Family != nil {
		profilesMap["family"] = profilesItem.Family
	}
	if profilesItem.Href != nil {
		profilesMap["href"] = profilesItem.Href
	}
	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.ResourceType != nil {
		profilesMap["resource_type"] = profilesItem.ResourceType
	}
	if profilesItem.Capacity != nil {
		capacityList := []map[string]interface{}{}
		capacity := profilesItem.Capacity.(*vpcv1.ShareProfileCapacity)
		capacityMap := dataSourceShareProfileCapacityToMap(*capacity)
		capacityList = append(capacityList, capacityMap)
		profilesMap["capacity"] = capacityList
	}
	if profilesItem.Iops != nil {
		iopsList := []map[string]interface{}{}
		iops := profilesItem.Iops.(*vpcv1.ShareProfileIops)
		iopsMap := dataSourceShareProfileIopsToMap(*iops)
		iopsList = append(iopsList, iopsMap)
		profilesMap["iops"] = iopsList
	}
	return profilesMap
}
