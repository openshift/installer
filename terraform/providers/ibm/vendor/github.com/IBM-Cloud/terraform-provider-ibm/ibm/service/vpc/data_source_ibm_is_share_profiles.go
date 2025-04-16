// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
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
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listShareProfilesOptions := &vpcbetav1.ListShareProfilesOptions{}

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

func dataSourceShareProfileCollectionFlattenProfiles(result []vpcbetav1.ShareProfile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceShareProfileCollectionProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceShareProfileCollectionProfilesToMap(profilesItem vpcbetav1.ShareProfile) (profilesMap map[string]interface{}) {
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

	return profilesMap
}
