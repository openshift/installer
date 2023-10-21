package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEcdDesktops() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdDesktopsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"desktop_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"end_user_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Deleted", "Expired", "Pending", "Running", "Starting", "Stopped", "Stopping"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"desktops": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desktop_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_user_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdDesktopsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDesktops"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("desktop_name"); ok {
		request["DesktopName"] = v
	}
	if v, ok := d.GetOk("office_site_id"); ok {
		request["OfficeSiteId"] = v
	}

	if v, ok := d.GetOk("policy_group_id"); ok {
		request["PolicyGroupId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["DesktopStatus"] = v
	}
	if m, ok := d.GetOk("end_user_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("EndUserId.%d", k+1)] = v.(string)
		}
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var desktopNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		desktopNameRegex = r
	}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_desktops", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Desktops", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Desktops", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if desktopNameRegex != nil && !desktopNameRegex.MatchString(fmt.Sprint(item["DesktopName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DesktopId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"cpu":                  formatInt(object["Cpu"]),
			"create_time":          object["CreationTime"],
			"id":                   fmt.Sprint(object["DesktopId"]),
			"desktop_id":           fmt.Sprint(object["DesktopId"]),
			"desktop_name":         object["DesktopName"],
			"desktop_type":         object["DesktopType"],
			"directory_id":         object["DirectoryId"],
			"expired_time":         object["ExpiredTime"],
			"end_user_ids":         object["EndUserIds"],
			"image_id":             object["ImageId"],
			"memory":               fmt.Sprint(object["Memory"]),
			"network_interface_id": object["NetworkInterfaceId"],
			"payment_type":         convertEcdDesktopPaymentTypeResponse(object["ChargeType"]),
			"policy_group_id":      object["PolicyGroupId"],
			"status":               object["DesktopStatus"],
			"system_disk_size":     formatInt(object["SystemDiskSize"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DesktopName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("desktops", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
