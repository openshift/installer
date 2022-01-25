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

func dataSourceAlicloudImpAppTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudImpAppTemplatesRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"attached", "unattached"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_template_creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"config_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"integration_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scene": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sdk_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard_room_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudImpAppTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAppTemplates"
	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var appTemplateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		appTemplateNameRegex = r
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
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewImpClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-06-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_imp_app_templates", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Result.AppTemplateInfoList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Result.AppTemplateInfoList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if appTemplateNameRegex != nil && !appTemplateNameRegex.MatchString(fmt.Sprint(item["AppTemplateName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["AppTemplateId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"app_template_creator": object["AppTemplateCreator"],
			"id":                   fmt.Sprint(object["AppTemplateId"]),
			"app_template_id":      fmt.Sprint(object["AppTemplateId"]),
			"app_template_name":    object["AppTemplateName"],
			"component_list":       object["ComponentList"],
			"create_time":          object["CreateTime"],
			"integration_mode":     object["IntegrationMode"],
			"scene":                object["Scene"],
			"sdk_info":             object["SdkInfo"],
			"standard_room_info":   object["StandardRoomInfo"],
			"status":               object["Status"],
		}
		if configList, ok := object["ConfigList"]; ok && configList != nil {
			configListMaps := make([]map[string]interface{}, 0)
			for _, configListItem := range configList.([]interface{}) {
				if configListItemMap, ok := configListItem.(map[string]interface{}); ok {
					configListMap := make(map[string]interface{}, 0)
					configListMap["key"] = configListItemMap["Key"]
					configListMap["value"] = configListItemMap["Value"]
					configListMaps = append(configListMaps, configListMap)
				}
			}
			mapping["config_list"] = configListMaps
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AppTemplateName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("templates", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
