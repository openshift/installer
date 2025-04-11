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

func DataSourceIBMIsClusterNetworkSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsClusterNetworkSubnetRead,

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
			"available_ipv4_address_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of IPv4 addresses in this cluster network subnet that are not in use, and have not been reserved by the user or the provider.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the cluster network subnet was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this cluster network subnet.",
			},
			"ip_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for this cluster network subnet.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"ipv4_cidr_block": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IPv4 range of this cluster network subnet, expressed in CIDR format.",
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
				Description: "The lifecycle state of the cluster network subnet.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this cluster network subnet. The name is unique across all cluster network subnets in the cluster network.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"total_ipv4_address_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of IPv4 addresses in this cluster network subnet.Note: This is calculated as 2<sup>(32 - prefix length)</sup>. For example, the prefix length `/24` gives:<br> 2<sup>(32 - 24)</sup> = 2<sup>8</sup> = 256 addresses.",
			},
		},
	}
}

func dataSourceIBMIsClusterNetworkSubnetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnet", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClusterNetworkSubnetOptions := &vpcv1.GetClusterNetworkSubnetOptions{}

	getClusterNetworkSubnetOptions.SetClusterNetworkID(d.Get("cluster_network_id").(string))
	getClusterNetworkSubnetOptions.SetID(d.Get("cluster_network_subnet_id").(string))

	clusterNetworkSubnet, _, err := vpcClient.GetClusterNetworkSubnetWithContext(context, getClusterNetworkSubnetOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClusterNetworkSubnetWithContext failed: %s", err.Error()), "(Data) ibm_is_cluster_network_subnet", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getClusterNetworkSubnetOptions.ClusterNetworkID, *getClusterNetworkSubnetOptions.ID))

	if err = d.Set("available_ipv4_address_count", flex.IntValue(clusterNetworkSubnet.AvailableIpv4AddressCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting available_ipv4_address_count: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-available_ipv4_address_count").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(clusterNetworkSubnet.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("href", clusterNetworkSubnet.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-href").GetDiag()
	}

	if err = d.Set("ip_version", clusterNetworkSubnet.IPVersion); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ip_version: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-ip_version").GetDiag()
	}

	if err = d.Set("ipv4_cidr_block", clusterNetworkSubnet.Ipv4CIDRBlock); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ipv4_cidr_block: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-ipv4_cidr_block").GetDiag()
	}

	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range clusterNetworkSubnet.LifecycleReasons {
		lifecycleReasonsItemMap, err := DataSourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(&lifecycleReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_cluster_network_subnet", "read", "lifecycle_reasons-to-map").GetDiag()
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_reasons: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-lifecycle_reasons").GetDiag()
	}

	if err = d.Set("lifecycle_state", clusterNetworkSubnet.LifecycleState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-lifecycle_state").GetDiag()
	}

	if err = d.Set("name", clusterNetworkSubnet.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-name").GetDiag()
	}

	if err = d.Set("resource_type", clusterNetworkSubnet.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-resource_type").GetDiag()
	}

	if err = d.Set("total_ipv4_address_count", flex.IntValue(clusterNetworkSubnet.TotalIpv4AddressCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_ipv4_address_count: %s", err), "(Data) ibm_is_cluster_network_subnet", "read", "set-total_ipv4_address_count").GetDiag()
	}

	return nil
}

func DataSourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(model *vpcv1.ClusterNetworkSubnetLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
