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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsClusterNetworkSubnetReservedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkSubnetReservedIPRead,

		Schema: map[string]*schema.Schema{
			"cluster_network_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network identifier.",
			},
			"cluster_network_subnet_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network subnet identifier.",
			},
			"cluster_network_subnet_reserved_ip_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The cluster network subnet reserved IP identifier.",
			},
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address.If the address is pending allocation, the value will be `0.0.0.0`.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this cluster network subnet reserved IP member will be automatically deleted when either `target` is deleted, or the cluster network subnet reserved IP is unbound.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network subnet reserved IP was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network subnet reserved IP.",
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
				Description: "The lifecycle state of the cluster network subnet reserved IP.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this cluster network subnet reserved IP. The name is unique across all reserved IPs in a cluster network subnet.",
			},
			"owner": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the cluster network subnet reserved IPThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target this cluster network subnet reserved IP is bound to.If absent, this cluster network subnet reserved IP is provider-owned or unbound.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The URL for this cluster network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this cluster network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this cluster network interface. The name is unique across all interfaces in the cluster network.",
						},
						"resource_type": &schema.Schema{
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

func dataSourceIBMIsClusterNetworkSubnetReservedIPRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkSubnetReservedIPOptions := &vpcv1.GetClusterNetworkSubnetReservedIPOptions{}

	getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(d.Get("cluster_network_subnet_id").(string))
	getClusterNetworkSubnetReservedIPOptions.SetID(d.Get("cluster_network_subnet_reserved_ip_id").(string))

	clusterNetworkSubnetReservedIP, _, err := vpcClient.GetClusterNetworkSubnetReservedIPWithContext(context, getClusterNetworkSubnetReservedIPOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkSubnetReservedIPWithContext failed: %s", err.Error()), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *getClusterNetworkSubnetReservedIPOptions.ClusterNetworkID, *getClusterNetworkSubnetReservedIPOptions.ClusterNetworkSubnetID, *getClusterNetworkSubnetReservedIPOptions.ID))

	if err = d.Set("address", clusterNetworkSubnetReservedIP.Address); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting address: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-address").GetDiag()
	}

	if err = d.Set("auto_delete", clusterNetworkSubnetReservedIP.AutoDelete); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting auto_delete: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-auto_delete").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkSubnetReservedIP.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("href", clusterNetworkSubnetReservedIP.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-href").GetDiag()
	}

	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetworkSubnetReservedIP.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_reasons: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-lifecycle_reasons").GetDiag()
	}

	if err = d.Set("lifecycle_state", clusterNetworkSubnetReservedIP.LifecycleState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-lifecycle_state").GetDiag()
	}

	if err = d.Set("name", clusterNetworkSubnetReservedIP.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-name").GetDiag()
	}

	if err = d.Set("owner", clusterNetworkSubnetReservedIP.Owner); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting owner: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-owner").GetDiag()
	}

	if err = d.Set("resource_type", clusterNetworkSubnetReservedIP.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-resource_type").GetDiag()
	}

	if !core.IsNil(clusterNetworkSubnetReservedIP.Target) {
		target := []map[string]interface{}{}
		targetMap, err := DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(clusterNetworkSubnetReservedIP.Target)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "target-to-map").GetDiag()
		}
		target = append(target, targetMap)
		if err = d.Set("target", target); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting target: %s", err), "(Data) ibm_is_cluster_network_subnet_reserved_ip", "read", "set-target").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(model *vpcv1.ClusterNetworkSubnetReservedIPLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(model vpcv1.ClusterNetworkSubnetReservedIPTargetIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext); ok {
		return DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model.(*vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext))
	} else if _, ok := model.(*vpcv1.ClusterNetworkSubnetReservedIPTarget); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ClusterNetworkSubnetReservedIPTarget)
		if model.Deleted != nil {
			deletedMap, err := DataSourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.ClusterNetworkSubnetReservedIPTargetIntf subtype encountered")
	}
}

func DataSourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model *vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model.Deleted)
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
