// Copyright IBM Corp. 2024 All Rights Reserved.
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

func DataSourceIBMIsClusterNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkRead,

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this cluster network.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network.",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the cluster network.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this cluster network. The name must not be used by another cluster network in the region.",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"subnet_prefixes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IP address ranges available for subnets for this cluster network.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allocation_policy": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The allocation policy for this subnet prefix:- `auto`: Subnets created by total count in this cluster network can use this prefix.",
						},
						"cidr": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CIDR block for this prefix.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this cluster network resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone this cluster network resides in.",
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
	}
}

func dataSourceIBMIsClusterNetworkRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

	getClusterNetworkOptions.SetID(d.Get("cluster_network_id").(string))

	clusterNetwork, _, err := vpcClient.GetClusterNetworkWithContext(context, getClusterNetworkOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkWithContext failed: %s", err.Error()), "(Data) ibm_is_cluster_network", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getClusterNetworkOptions.ID)

	if err = d.Set("created_at", flex.DateTimeToString(clusterNetwork.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_cluster_network", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("crn", clusterNetwork.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_cluster_network", "read", "set-crn").GetDiag()
	}

	if err = d.Set("href", clusterNetwork.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_cluster_network", "read", "set-href").GetDiag()
	}

	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetwork.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_reasons: %s", err), "(Data) ibm_is_cluster_network", "read", "set-lifecycle_reasons").GetDiag()
	}

	if err = d.Set("lifecycle_state", clusterNetwork.LifecycleState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_cluster_network", "read", "set-lifecycle_state").GetDiag()
	}

	if err = d.Set("name", clusterNetwork.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_cluster_network", "read", "set-name").GetDiag()
	}

	profile := []map[string]interface{}{}
	profileMap, err := DataSourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(clusterNetwork.Profile)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "profile-to-map").GetDiag()
	}
	profile = append(profile, profileMap)
	if err = d.Set("profile", profile); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting profile: %s", err), "(Data) ibm_is_cluster_network", "read", "set-profile").GetDiag()
	}

	resourceGroup := []map[string]interface{}{}
	resourceGroupMap, err := DataSourceIBMIsClusterNetworkResourceGroupReferenceToMap(clusterNetwork.ResourceGroup)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "resource_group-to-map").GetDiag()
	}
	resourceGroup = append(resourceGroup, resourceGroupMap)
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group: %s", err), "(Data) ibm_is_cluster_network", "read", "set-resource_group").GetDiag()
	}

	if err = d.Set("resource_type", clusterNetwork.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_cluster_network", "read", "set-resource_type").GetDiag()
	}

	subnetPrefixes := []map[string]interface{}{}
	for _, subnetPrefixesItem := range clusterNetwork.SubnetPrefixes {
		subnetPrefixesItemMap, err := DataSourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(&subnetPrefixesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "subnet_prefixes-to-map").GetDiag()
		}
		subnetPrefixes = append(subnetPrefixes, subnetPrefixesItemMap)
	}
	if err = d.Set("subnet_prefixes", subnetPrefixes); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting subnet_prefixes: %s", err), "(Data) ibm_is_cluster_network", "read", "set-subnet_prefixes").GetDiag()
	}

	vpc := []map[string]interface{}{}
	vpcMap, err := DataSourceIBMIsClusterNetworkVPCReferenceToMap(clusterNetwork.VPC)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "vpc-to-map").GetDiag()
	}
	vpc = append(vpc, vpcMap)
	if err = d.Set("vpc", vpc); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting vpc: %s", err), "(Data) ibm_is_cluster_network", "read", "set-vpc").GetDiag()
	}

	zone := []map[string]interface{}{}
	zoneMap, err := DataSourceIBMIsClusterNetworkZoneReferenceToMap(clusterNetwork.Zone)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network", "read", "zone-to-map").GetDiag()
	}
	zone = append(zone, zoneMap)
	if err = d.Set("zone", zone); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting zone: %s", err), "(Data) ibm_is_cluster_network", "read", "set-zone").GetDiag()
	}

	return nil
}

func DataSourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(model *vpcv1.ClusterNetworkLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(model *vpcv1.ClusterNetworkProfileReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(model *vpcv1.ClusterNetworkSubnetPrefix) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["allocation_policy"] = *model.AllocationPolicy
	modelMap["cidr"] = *model.CIDR
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworkDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
