package alicloud

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCSKubernetesAddons() *schema.Resource {
	return &schema.Resource{
		Read: dataAlicloudCSKubernetesAddonsRead,

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"required": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataAlicloudCSKubernetesAddonsRead(d *schema.ResourceData, meta interface{}) error {
	filterAddons := make(map[string]*Component)
	var ids, names []string

	clusterId := d.Get("cluster_id").(string)

	addons, err := describeAvailableAddons(d, meta)
	for _, addon := range addons {
		if nameRegex, ok := d.GetOk("name_regex"); ok {
			r, err := regexp.Compile(nameRegex.(string))
			if err != nil {
				return WrapError(err)
			}
			if !r.MatchString(addon.ComponentName) {
				continue
			}
		}
		if ids, ok := d.GetOk("ids"); ok {
			findId := func(id string, ids []string) (ret bool) {
				for _, i := range ids {
					if id == i {
						ret = true
					}
				}
				return
			}
			if !findId(fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, addon.ComponentName), expandStringList(ids.([]interface{}))) {
				continue
			}
		}
		filterAddons[addon.ComponentName] = addon
		ids = append(ids, fmt.Sprintf("%s%s%s", clusterId, COLON_SEPARATED, addon.ComponentName))
		names = append(names, addon.ComponentName)
	}

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "describeAddonsMeta", err)
	}
	result := fetchAddonsMetadata(filterAddons)

	d.Set("cluster_id", clusterId)
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("addons", result)

	d.SetId(tea.ToString(hashcode.String(clusterId)))
	return nil
}

func describeAvailableAddons(d *schema.ResourceData, meta interface{}) (map[string]*Component, error) {
	clusterId := d.Get("cluster_id").(string)

	client, err := meta.(*connectivity.AliyunClient).NewRoaCsClient()
	if err != nil {
		return nil, err
	}
	csClient := CsClient{client}

	availableAddons, err := csClient.DescribeCsKubernetesAllAvailableAddons(clusterId)
	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, ResourceAlicloudCSKubernetesAddon, "DescribeCsKubernetesAllAvailableAddons", err)
	}

	return availableAddons, nil
}

func fetchAddonsMetadata(addonsMap map[string]*Component) []map[string]interface{} {
	result := []map[string]interface{}{}
	for name, addon := range addonsMap {
		state := map[string]interface{}{}
		state["name"] = name
		state["current_version"] = addon.Version
		state["next_version"] = addon.NextVersion
		state["required"] = addon.Required
		result = append(result, state)
	}
	return result
}
