package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudOtsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOtsInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"write_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"read_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": tagsSchemaComputed(),
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOtsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}

	allInstanceNames, err := otsService.ListOtsInstance(PageSizeLarge, 1)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instances", "ListOtsInstance", AlibabaCloudSdkGoERROR)
	}

	idsMap := make(map[string]bool)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, x := range v.([]interface{}) {
			if x == nil {
				continue
			}
			idsMap[x.(string)] = true
		}
	}

	var nameReg *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		nameReg = regexp.MustCompile(v.(string))
	}

	var filteredInstanceNames []string
	for _, instanceName := range allInstanceNames {
		// name_regex mismatch
		if nameReg != nil && !nameReg.MatchString(instanceName) {
			continue
		}
		// ids mismatch
		if len(idsMap) != 0 {
			if _, ok := idsMap[instanceName]; !ok {
				continue
			}
		}
		filteredInstanceNames = append(filteredInstanceNames, instanceName)
	}

	// get full instance info via GetInstance
	var allInstances []ots.InstanceInfo
	for _, instanceName := range filteredInstanceNames {
		instanceInfo, err := otsService.DescribeOtsInstance(instanceName)
		if err != nil {
			return WrapError(err)
		}
		allInstances = append(allInstances, instanceInfo)
	}

	// filter by tag.
	var filteredInstances []ots.InstanceInfo
	if v, ok := d.GetOk("tags"); ok {
		if vmap, ok := v.(map[string]interface{}); ok && len(vmap) > 0 {
			for _, instance := range allInstances {
				if tagsMapEqual(vmap, otsTagsToMap(instance.TagInfos.TagInfo)) {
					filteredInstances = append(filteredInstances, instance)
				}
			}
		} else {
			filteredInstances = allInstances[:]
		}
	} else {
		filteredInstances = allInstances[:]
	}
	return otsInstancesDecriptionAttributes(d, filteredInstances, meta)
}

func otsInstancesDecriptionAttributes(d *schema.ResourceData, instances []ots.InstanceInfo, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, instance := range instances {
		mapping := map[string]interface{}{
			"id":             instance.InstanceName,
			"name":           instance.InstanceName,
			"status":         string(convertOtsInstanceStatusConvert(instance.Status)),
			"write_capacity": instance.WriteCapacity,
			"read_capacity":  instance.ReadCapacity,
			"cluster_type":   instance.ClusterType,
			"create_time":    instance.CreateTime,
			"user_id":        instance.UserId,
			"network":        instance.Network,
			"description":    instance.Description,
			"entity_quota":   instance.Quota.EntityQuota,
			"tags":           otsTagsToMap(instance.TagInfos.TagInfo),
		}
		names = append(names, instance.InstanceName)
		ids = append(ids, instance.InstanceName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
