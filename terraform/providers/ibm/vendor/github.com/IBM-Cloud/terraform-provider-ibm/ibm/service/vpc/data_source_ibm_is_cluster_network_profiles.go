// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsClusterNetworkProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkProfilesRead,

		Schema: map[string]*schema.Schema{
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of cluster network profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this cluster network profile belongs to.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this cluster network profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this cluster network profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"supported_instance_profiles": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance profiles that support this cluster network profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this virtual server instance profile.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this virtual server instance profile.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"zones": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Zones in this region that support this cluster network profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsClusterNetworkProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_profiles", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listClusterNetworkProfilesOptions := &vpcv1.ListClusterNetworkProfilesOptions{}

	var pager *vpcv1.ClusterNetworkProfilesPager
	pager, err = vpcClient.NewClusterNetworkProfilesPager(listClusterNetworkProfilesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_profiles", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ClusterNetworkProfilesPager.GetAll() failed %s", err), "(Data) ibm_is_cluster_network_profiles", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsClusterNetworkProfilesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsClusterNetworkProfilesClusterNetworkProfileToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_profiles", "read", "ClusterNetworks-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("profiles", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting profiles %s", err), "(Data) ibm_is_cluster_network_profiles", "read", "profiles-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsClusterNetworkProfilesID returns a reasonable ID for the list.
func dataSourceIBMIsClusterNetworkProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsClusterNetworkProfilesClusterNetworkProfileToMap(model *vpcv1.ClusterNetworkProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["family"] = *model.Family
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	supportedInstanceProfiles := []map[string]interface{}{}
	for _, supportedInstanceProfilesItem := range model.SupportedInstanceProfiles {
		supportedInstanceProfilesItemMap, err := DataSourceIBMIsClusterNetworkProfilesInstanceProfileReferenceToMap(&supportedInstanceProfilesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		supportedInstanceProfiles = append(supportedInstanceProfiles, supportedInstanceProfilesItemMap)
	}
	modelMap["supported_instance_profiles"] = supportedInstanceProfiles
	zones := []map[string]interface{}{}
	for _, zonesItem := range model.Zones {
		zonesItemMap, err := DataSourceIBMIsClusterNetworkProfilesZoneReferenceToMap(&zonesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		zones = append(zones, zonesItemMap)
	}
	modelMap["zones"] = zones
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkProfilesInstanceProfileReferenceToMap(model *vpcv1.InstanceProfileReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkProfilesZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
