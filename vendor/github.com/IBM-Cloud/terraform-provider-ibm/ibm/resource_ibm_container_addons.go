// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/internal/hashcode"
)

func resourceIBMContainerAddOns() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMContainerAddOnsCreate,
		Read:     resourceIBMContainerAddOnsRead,
		Update:   resourceIBMContainerAddOnsUpdate,
		Delete:   resourceIBMContainerAddOnsDelete,
		Exists:   resourceIBMContainerAddOnsExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster Name or ID",
			},
			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the resource group.",
				ForceNew:    true,
				Computed:    true,
			},
			"addons": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceIBMContainerAddonsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The addon name such as 'istio'.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    false,
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
func resourceIBMContainerAddOnsCreate(d *schema.ResourceData, meta interface{}) error {
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
	existingAddons, err := addOnAPI.GetAddons(cluster, targetEnv)
	if err != nil {
		fmt.Println("[ WARN ] Error getting Addons.")
	}

	payload, err := expandAddOns(d, meta, cluster, targetEnv, existingAddons)
	if err != nil {
		return fmt.Errorf("Error in getting addons from expandAddOns %s", err)
	}
	payload.Enable = true
	_, err = addOnAPI.ConfigureAddons(cluster, &payload, targetEnv)
	if err != nil {
		return err
	}
	_, err = waitForContainerAddOns(d, meta, cluster, schema.TimeoutCreate)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Enabling Addon (%s) : %s", d.Id(), err)
	}
	d.SetId(cluster)

	return resourceIBMContainerAddOnsRead(d, meta)
}
func expandAddOns(d *schema.ResourceData, meta interface{}, cluster string, targetEnv v1.ClusterTargetHeader, existingAddons []v1.AddOn) (addOns v1.ConfigureAddOns, err error) {
	addOnSet := d.Get("addons").(*schema.Set).List()
	if existingAddons == nil || len(existingAddons) < 1 {
		for _, aoSet := range addOnSet {
			ao, _ := aoSet.(map[string]interface{})
			addOn := v1.AddOn{
				Name: ao["name"].(string),
			}
			if ao["version"] != nil {
				addOn.Version = ao["version"].(string)
			}
			addOns.AddonsList = append(addOns.AddonsList, addOn)
		}
	}
	if existingAddons != nil && len(existingAddons) > 0 {
		for _, aoSet := range addOnSet {
			ao, _ := aoSet.(map[string]interface{})
			exist := false
			for _, existAddon := range existingAddons {
				if existAddon.Name == ao["name"].(string) {
					exist = true
					if existAddon.Version != ao["version"].(string) {
						err := updateAddOnVersion(d, meta, ao, cluster, targetEnv)
						if err != nil {
							return addOns, err
						}
					}
				}
			}
			if !exist {
				addOn := v1.AddOn{
					Name: ao["name"].(string),
				}
				if ao["version"] != nil {
					addOn.Version = ao["version"].(string)
				}
				addOns.AddonsList = append(addOns.AddonsList, addOn)
			}
		}
	}

	return addOns, nil
}
func updateAddOnVersion(d *schema.ResourceData, meta interface{}, u map[string]interface{}, cluster string, targetEnv v1.ClusterTargetHeader) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	addOnAPI := csClient.AddOns()

	update := v1.AddOn{
		Name: u["name"].(string),
	}
	if u["version"].(string) != "" {
		update.Version = u["version"].(string)
	}
	updateList := v1.ConfigureAddOns{}
	updateList.AddonsList = append(updateList.AddonsList, update)
	updateList.Update = true
	_, err = addOnAPI.ConfigureAddons(cluster, &updateList, targetEnv)
	if err != nil {
		return err
	}
	if !d.IsNewResource() {
		_, err = waitForContainerAddOns(d, meta, cluster, schema.TimeoutUpdate)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for Updating Addon (%s) : %s", d.Id(), err)
		}
	}

	return nil
}
func resourceIBMContainerAddOnsRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	addOnAPI := csClient.AddOns()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	cluster := d.Id()

	result, err := addOnAPI.GetAddons(cluster, targetEnv)
	if err != nil {
		return err
	}
	d.Set("cluster", cluster)
	addOns, err := flattenAddOns(result)
	if err != nil {
		fmt.Printf("Error Flattening Addons list %s", err)
	}
	d.Set("resource_group_id", targetEnv.ResourceGroup)
	d.Set("addons", addOns)
	return nil
}
func flattenAddOns(result []v1.AddOn) (resp *schema.Set, err error) {
	addOns := []interface{}{}
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

	return schema.NewSet(resourceIBMContainerAddonsHash, addOns), nil
}
func resourceIBMContainerAddOnsUpdate(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	addOnAPI := csClient.AddOns()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	cluster := d.Id()

	if d.HasChange("addons") && !d.IsNewResource() {
		oldList, newList := d.GetChange("addons")
		if oldList == nil {
			oldList = new(schema.Set)
		}
		if newList == nil {
			newList = new(schema.Set)
		}
		os := oldList.(*schema.Set)
		ns := newList.(*schema.Set)
		for _, nA := range ns.List() {
			newPack := nA.(map[string]interface{})
			for _, oA := range os.List() {
				oldPack := oA.(map[string]interface{})
				if (strings.Compare(newPack["name"].(string), oldPack["name"].(string)) == 0) && (strings.Compare(newPack["version"].(string), oldPack["version"].(string)) != 0) && (newPack["version"].(string) != "") {
					err := updateAddOnVersion(d, meta, newPack, cluster, targetEnv)
					if err != nil {
						return err
					}
					ns.Remove(nA)
					os.Remove(oA)
				}
			}
		}
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(add) > 0 {
			addOnParams := v1.ConfigureAddOns{}
			for _, addon := range add {
				newAddon := addon.(map[string]interface{})
				addOnParam := v1.AddOn{
					Name: newAddon["name"].(string),
				}
				if newAddon["version"] != nil {
					addOnParam.Version = newAddon["version"].(string)
				}
				addOnParams.AddonsList = append(addOnParams.AddonsList, addOnParam)

			}
			addOnParams.Enable = true
			_, err = addOnAPI.ConfigureAddons(cluster, &addOnParams, targetEnv)
			if err != nil {
				return err
			}
			_, err = waitForContainerAddOns(d, meta, cluster, schema.TimeoutCreate)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for Enabling Addon (%s) : %s", d.Id(), err)
			}
		}
		if len(remove) > 0 {
			addOnParams := v1.ConfigureAddOns{}
			for _, addOn := range remove {
				oldAddOn := addOn.(map[string]interface{})
				addOnParam := v1.AddOn{
					Name: oldAddOn["name"].(string),
				}
				if oldAddOn["version"] != nil {
					addOnParam.Version = oldAddOn["version"].(string)
				}
				addOnParams.AddonsList = append(addOnParams.AddonsList, addOnParam)
			}
			addOnParams.Enable = false
			_, err = addOnAPI.ConfigureAddons(cluster, &addOnParams, targetEnv)
			if err != nil {
				return err
			}
		}
	}

	return resourceIBMContainerAddOnsRead(d, meta)
}
func resourceIBMContainerAddOnsDelete(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return err
	}
	addOnAPI := csClient.AddOns()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return err
	}
	cluster := d.Id()
	payload, err := expandAddOns(d, meta, cluster, targetEnv, nil)
	if err != nil {
		return fmt.Errorf("Error in getting addons from expandAddOns %s", err)
	}

	payload.Enable = false
	_, err = addOnAPI.ConfigureAddons(cluster, &payload, targetEnv)
	if err != nil {
		return err
	}

	return nil
}
func waitForContainerAddOns(d *schema.ResourceData, meta interface{}, cluster, timeout string) (interface{}, error) {
	addOnClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return false, err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending", "updating", ""},
		Target:  []string{"normal", "warning", "critical", "available"},
		Refresh: func() (interface{}, string, error) {
			targetEnv, err := getClusterTargetHeader(d, meta)
			if err != nil {
				return nil, "", err
			}
			addOns, err := addOnClient.AddOns().GetAddons(cluster, targetEnv)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The resource addons %s does not exist anymore: %v", d.Id(), err)
				}
				return nil, "", err
			}
			for _, addOn := range addOns {
				if addOn.HealthState == "pending" || addOn.HealthState == "updating" || addOn.HealthState == "" {
					return addOns, addOn.HealthState, nil
				}
			}
			return addOns, "available", nil
		},
		Timeout:    d.Timeout(timeout),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
func resourceIBMContainerAddOnsExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	csClient, err := meta.(ClientSession).ContainerAPI()
	if err != nil {
		return false, err
	}
	addOnAPI := csClient.AddOns()

	targetEnv, err := getClusterTargetHeader(d, meta)
	if err != nil {
		return false, err
	}
	cluster := d.Id()

	_, err = addOnAPI.GetAddons(cluster, targetEnv)
	if err != nil {
		if apiErr, ok := err.(bmxerror.RequestFailure); ok {
			if apiErr.StatusCode() == 404 {
				return false, nil
			}
		}
		return false, fmt.Errorf("Error communicating with the API: %s", err)
	}

	return true, nil
}

func resourceIBMContainerAddonsHash(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", a["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["version"].(string)))

	return hashcode.String(buf.String())
}
