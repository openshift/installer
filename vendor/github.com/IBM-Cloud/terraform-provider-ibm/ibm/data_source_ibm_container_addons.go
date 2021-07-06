// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func datasourceIBMContainerAddOns() *schema.Resource {
	return &schema.Resource{
		Read: datasourceIBMContainerAddOnsRead,

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster Name or ID",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the resource group.",
			},
			"addons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The List of AddOns",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The addon name such as 'istio'.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The addon version, omit the version if you wish to use the default version.",
						},
						"allowed_upgrade_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The versions that the addon can be upgraded to",
						},
						"deprecated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines if this addon version is deprecated",
						},
						"health_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health state for this addon, a short indication (e.g. critical, pending)",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health status for this addon, provides a description of the state (e.g. error message)",
						},
						"min_kube_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum kubernetes version for this addon.",
						},
						"min_ocp_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum OpenShift version for this addon.",
						},
						"supported_kube_range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The supported kubernetes version range for this addon.",
						},
						"target_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The addon target version.",
						},
						"vlan_spanning_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "VLAN spanning required for multi-zone clusters",
						},
					},
				},
			},
		},
	}
}
func datasourceIBMContainerAddOnsRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	addOnAPI := csClient.AddOns()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	cluster := d.Get("cluster").(string)

	result, err := addOnAPI.GetAddons(cluster, targetEnv)
	if err != nil {
		return err
	}
	d.Set("cluster", cluster)
	addOns, err := flattenAddOnsList(result)
	if err != nil {
		fmt.Printf("Error Flattening Addons list %s", err)
	}
	d.Set("resource_group_id", targetEnv.ResourceGroup)
	d.Set("addons", addOns)
	d.SetId(cluster)
	return nil
}
func flattenAddOnsList(result []v1.AddOn) (addOns []map[string]interface{}, err error) {
	for _, addOn := range result {
		record := map[string]interface{}{}
		record["name"] = addOn.Name
		record["version"] = addOn.Version
		if len(addOn.AllowedUpgradeVersion) > 0 {
			record["allowed_upgrade_versions"] = addOn.AllowedUpgradeVersion
		}
		if &addOn.Deprecated != nil {
			record["deprecated"] = addOn.Deprecated
		}
		if &addOn.HealthState != nil {
			record["health_state"] = addOn.HealthState
		}
		if &addOn.HealthStatus != nil {
			record["health_status"] = addOn.HealthStatus
		}
		if addOn.MinKubeVersion != "" {
			record["min_kube_version"] = addOn.MinKubeVersion
		}
		if addOn.MinOCPVersion != "" {
			record["min_ocp_version"] = addOn.MinOCPVersion
		}
		if addOn.SupportedKubeRange != "" {
			record["supported_kube_range"] = addOn.SupportedKubeRange
		}
		if addOn.TargetVersion != "" {
			record["target_version"] = addOn.TargetVersion
		}
		if &addOn.VlanSpanningRequired != nil {
			record["vlan_spanning_required"] = addOn.VlanSpanningRequired
		}

		addOns = append(addOns, record)
	}

	return addOns, nil
}
