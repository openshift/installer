package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudVpcTrafficMirrorFilterIngressRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudVpcTrafficMirrorFilterIngressRulesRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Created", "Creating", "Deleting", "Modifying"}, false),
			},
			"traffic_mirror_filter_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"traffic_mirror_filter_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudVpcTrafficMirrorFilterIngressRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTrafficMirrorFilters"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["TrafficMirrorFilterIds"] = []string{fmt.Sprint(d.Get("traffic_mirror_filter_id"))}
	request["MaxResults"] = PageSizeSmall
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_vpc_traffic_mirror_filter_ingress_rules", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.TrafficMirrorFilters", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrafficMirrorFilters", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		EgressRules := v.(map[string]interface{})["IngressRules"]
		if EgressRulesMap, ok := EgressRules.([]interface{}); ok {
			for _, v := range EgressRulesMap {
				item := v.(map[string]interface{})
				if len(idsMap) > 0 {
					if _, ok := idsMap[fmt.Sprint(item["TrafficMirrorFilterId"], ":", item["TrafficMirrorFilterRuleId"])]; !ok {
						continue
					}
				}
				if statusOk && status.(string) != "" && status.(string) != item["TrafficMirrorFilterRuleStatus"].(string) {
					continue
				}
				objects = append(objects, item)
			}
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"destination_cidr_block":        object["DestinationCidrBlock"],
			"destination_port_range":        object["DestinationPortRange"],
			"priority":                      formatInt(object["Priority"]),
			"rule_action":                   object["Action"],
			"protocol":                      object["Protocol"],
			"source_cidr_block":             object["SourceCidrBlock"],
			"source_port_range":             object["SourcePortRange"],
			"status":                        object["TrafficMirrorFilterRuleStatus"],
			"traffic_mirror_filter_id":      object["TrafficMirrorFilterId"],
			"id":                            fmt.Sprint(object["TrafficMirrorFilterId"], ":", object["TrafficMirrorFilterRuleId"]),
			"traffic_mirror_filter_rule_id": fmt.Sprint(object["TrafficMirrorFilterRuleId"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
