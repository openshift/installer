// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package kubernetes

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMContainerAddOns() *schema.Resource {
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
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		log.Println("[ WARN ] Error getting Addons.")
	}
	payload, err := expandAddOns(d, meta, cluster, targetEnv, existingAddons)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in getting addons from expandAddOns during Create: %s", err)
	}
	payload.Enable = true
	_, err = addOnAPI.ConfigureAddons(cluster, &payload, targetEnv)
	if err != nil {
		return err
	}
	_, err = waitForContainerAddOns(d, meta, cluster, schema.TimeoutCreate)
	if err != nil {
		return fmt.Errorf("[ERROR] Error waiting for Addon to reach normal during create (%s) : %s", d.Id(), err)
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
	if len(existingAddons) > 0 {
		csClient, err := meta.(conns.ClientSession).ContainerAPI()
		if err != nil {
			return addOns, err
		}
		addOnAPI := csClient.AddOns()
		for _, aoSet := range addOnSet {
			ao, _ := aoSet.(map[string]interface{})
			exist := false
			for _, existAddon := range existingAddons {
				if existAddon.Name == ao["name"].(string) {
					exist = true
					if existAddon.Version != ao["version"].(string) {
						if flex.StringContains(existAddon.AllowedUpgradeVersion, ao["version"].(string)) {
							// This block upgrates addon version if addon has `allowed_upgrade_versions`
							err := updateAddOnVersion(d, meta, ao, cluster, targetEnv)
							if err != nil {
								return addOns, err
							}
						} else if (ao["version"].(string) == existAddon.TargetVersion) && (!flex.StringContains(existAddon.AllowedUpgradeVersion, ao["version"].(string))) {
							// This block reinstalls addons that dont have upgradation capability
							//Uninstall AddOn with old version
							rmParams := v1.ConfigureAddOns{}
							rmParam := v1.AddOn{
								Name: existAddon.Name,
							}
							if existAddon.Version != "" {
								rmParam.Version = existAddon.Version
							}
							rmParams.AddonsList = append(rmParams.AddonsList, rmParam)
							rmParams.Enable = false
							_, err = addOnAPI.ConfigureAddons(cluster, &rmParams, targetEnv)
							if err != nil {
								return addOns, fmt.Errorf("[ERROR] Error uninstalling addon %s  on %s during create : %s", d.Id(), existAddon.Name, err)
							}
							//Install AddOn with new version
							addParams := v1.ConfigureAddOns{}
							addParam := v1.AddOn{
								Name: ao["name"].(string),
							}
							if ao["version"] != nil {
								rmParam.Version = ao["version"].(string)
							}
							addParams.AddonsList = append(addParams.AddonsList, addParam)
							addParams.Enable = true
							_, err = addOnAPI.ConfigureAddons(cluster, &addParams, targetEnv)
							if err != nil {
								return addOns, fmt.Errorf("[ERROR] Error installing addon %s  on %s during update : %s", d.Id(), ao["name"], err)
							}

						} else {
							return addOns, fmt.Errorf("[ERROR] The given addon is not provided with upgradable or updatable version")
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
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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

	return nil
}
func resourceIBMContainerAddOnsRead(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		record["deprecated"] = addOn.Deprecated
		record["health_state"] = addOn.HealthState
		record["health_status"] = addOn.HealthStatus

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

		record["vlan_spanning_required"] = addOn.VlanSpanningRequired

		addOns = append(addOns, record)
	}

	return schema.NewSet(resourceIBMContainerAddonsHash, addOns), nil
}
func resourceIBMContainerAddOnsUpdate(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
				if (strings.Compare(newPack["name"].(string), oldPack["name"].(string)) == 0) && (strings.Compare(newPack["version"].(string), oldPack["version"].(string)) != 0) {
					if flex.StringContains(flex.ExpandStringList(oldPack["allowed_upgrade_versions"].([]interface{})), newPack["version"].(string)) {
						// This block upgrates addon version if addon has `allowed_upgrade_versions`
						err := updateAddOnVersion(d, meta, newPack, cluster, targetEnv)
						if err != nil {
							return err
						}
						ns.Remove(nA)
						os.Remove(oA)
					} else if (newPack["version"].(string) == oldPack["target_version"].(string)) && (!flex.StringContains(flex.ExpandStringList(oldPack["allowed_upgrade_versions"].([]interface{})), newPack["version"].(string))) {
						// This block reinstalls addons that dont have upgradation capability
						//Uninstall AddOn with old version
						rmParams := v1.ConfigureAddOns{}
						rmParam := v1.AddOn{
							Name: oldPack["name"].(string),
						}
						if oldPack["version"] != nil {
							rmParam.Version = oldPack["version"].(string)
						}
						rmParams.AddonsList = append(rmParams.AddonsList, rmParam)
						rmParams.Enable = false
						_, err = addOnAPI.ConfigureAddons(cluster, &rmParams, targetEnv)
						if err != nil {
							return fmt.Errorf("[ERROR] Error uninstalling addon %s  on %s during update : %s", d.Id(), oldPack["name"], err)
						}
						//Install AddOn with new version
						addParams := v1.ConfigureAddOns{}
						addParam := v1.AddOn{
							Name: newPack["name"].(string),
						}
						if newPack["version"] != nil {
							rmParam.Version = newPack["version"].(string)
						}
						addParams.AddonsList = append(addParams.AddonsList, addParam)
						addParams.Enable = true
						_, err = addOnAPI.ConfigureAddons(cluster, &addParams, targetEnv)
						if err != nil {
							return fmt.Errorf("[ERROR] Error installing addon %s  on %s during update : %s", d.Id(), newPack["name"], err)
						}
						ns.Remove(nA)
						os.Remove(oA)

					} else {
						return fmt.Errorf("[ERROR] The given addon is not provided with upgradable or updatable version")
					}
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
		_, err = waitForContainerAddOns(d, meta, cluster, schema.TimeoutUpdate)
		if err != nil {
			return fmt.Errorf("[ERROR] Error waiting for Addon to reach normal during update (%s) : %s", d.Id(), err)
		}
	}

	return resourceIBMContainerAddOnsRead(d, meta)
}

func resourceIBMContainerAddOnsDelete(d *schema.ResourceData, meta interface{}) error {
	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		return fmt.Errorf("[ERROR] Error in getting addons from expandAddOns during Destroy %s", err)
	}

	payload.Enable = false
	_, err = addOnAPI.ConfigureAddons(cluster, &payload, targetEnv)
	if err != nil {
		return err
	}

	return nil
}
func waitForContainerAddOns(d *schema.ResourceData, meta interface{}, cluster, timeout string) (interface{}, error) {
	addOnClient, err := meta.(conns.ClientSession).ContainerAPI()
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
					return nil, "", fmt.Errorf("[ERROR] The resource addons %s does not exist anymore: %v", d.Id(), err)
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

	csClient, err := meta.(conns.ClientSession).ContainerAPI()
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
		return false, fmt.Errorf("[ERROR] Error getting container addons: %s", err)
	}

	return true, nil
}

func resourceIBMContainerAddonsHash(v interface{}) int {
	var buf bytes.Buffer
	a := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", a["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", a["version"].(string)))

	return conns.String(buf.String())
}
