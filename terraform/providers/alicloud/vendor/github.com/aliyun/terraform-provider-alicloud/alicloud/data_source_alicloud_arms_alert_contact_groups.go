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

func dataSourceAlicloudArmsAlertContactGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudArmsAlertContactGroupsRead,
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
			"alert_contact_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"contact_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"contact_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_contact_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alert_contact_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudArmsAlertContactGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "SearchAlertContactGroup"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("alert_contact_group_name"); ok {
		request["ContactGroupName"] = v
	}
	if v, ok := d.GetOk("contact_id"); ok {
		request["ContactId"] = v
	}
	if v, ok := d.GetOk("contact_name"); ok {
		request["ContactName"] = v
	}
	request["IsDetail"] = true
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var alertContactGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		alertContactGroupNameRegex = r
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
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_alert_contact_groups", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ContactGroups", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if alertContactGroupNameRegex != nil && !alertContactGroupNameRegex.MatchString(fmt.Sprint(item["ContactGroupName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["ContactGroupId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                       fmt.Sprint(object["ContactGroupId"]),
			"alert_contact_group_id":   fmt.Sprint(object["ContactGroupId"]),
			"alert_contact_group_name": object["ContactGroupName"],
			"create_time":              fmt.Sprint(object["CreateTime"]),
		}
		contactIdsItems := make([]string, 0)
		if contacts, ok := object["Contacts"]; ok && contacts != nil {
			for _, contactsItem := range contacts.([]interface{}) {
				if contactId, ok := contactsItem.(map[string]interface{})["ContactId"]; ok && contactId != nil {
					contactIdsItems = append(contactIdsItems, fmt.Sprint(contactId))
				}
			}
		}
		mapping["contact_ids"] = contactIdsItems
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ContactGroupName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
