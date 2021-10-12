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

func dataSourceAlicloudCloudFirewallControlPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudFirewallControlPoliciesRead,
		Schema: map[string]*schema.Schema{
			"acl_action": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop", "log"}, false),
			},
			"acl_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"in", "out"}, false),
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh"}, false),
			},
			"proto": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{" TCP", " UDP", "ANY", "ICMP"}, false),
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"acl_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port_group": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_port_group_ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"dest_port_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"destination_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_result_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hit_times": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"order": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"proto": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"release": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudFirewallControlPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeControlPolicy"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("acl_action"); ok {
		request["AclAction"] = v
	}
	if v, ok := d.GetOk("acl_uuid"); ok {
		request["AclUuid"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("destination"); ok {
		request["Destination"] = v
	}
	request["Direction"] = d.Get("direction")
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("proto"); ok {
		request["Proto"] = v
	}
	if v, ok := d.GetOk("source"); ok {
		request["Source"] = v
	}
	if v, ok := d.GetOk("source_ip"); ok {
		request["SourceIp"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["CurrentPage"] = 1
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCloudfwClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-07"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_firewall_control_policies", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Policys", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policys", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["CurrentPage"] = request["CurrentPage"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                      fmt.Sprint(object["AclUuid"], ":", object["Direction"]),
			"acl_action":              object["AclAction"],
			"acl_uuid":                object["AclUuid"],
			"application_id":          object["ApplicationId"],
			"application_name":        object["ApplicationName"],
			"description":             object["Description"],
			"dest_port":               object["DestPort"],
			"dest_port_group":         object["DestPortGroup"],
			"dest_port_group_ports":   object["DestPortGroupPorts"],
			"dest_port_type":          object["DestPortType"],
			"destination":             object["Destination"],
			"destination_group_cidrs": object["DestinationGroupCidrs"],
			"destination_group_type":  object["DestinationGroupType"],
			"destination_type":        object["DestinationType"],
			"direction":               object["Direction"],
			"dns_result":              object["DnsResult"],
			"dns_result_time":         fmt.Sprint(object["DnsResultTime"]),
			"hit_times":               fmt.Sprint(object["HitTimes"]),
			"order":                   formatInt(object["Order"]),
			"proto":                   object["Proto"],
			"release":                 object["Release"],
			"source":                  object["Source"],
			"source_group_cidrs":      object["SourceGroupCidrs"],
			"source_group_type":       object["SourceGroupType"],
			"source_type":             object["SourceType"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
