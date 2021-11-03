package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudActiontrailTrails() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudActiontrailTrailsRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"include_organization_trail": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"include_shadow_trails": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disable", "Enable", "Fresh"}, false),
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
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
			"trails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_rw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_organization_trail": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_project_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"actiontrails": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_rw": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_organization_trail": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_key_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_project_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sls_write_role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trail_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Deprecated: "Field 'actiontrails' has been deprecated from version 1.95.0. Use 'trails' instead.",
			},
		},
	}
}

func dataSourceAlicloudActiontrailTrailsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeTrails"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("include_organization_trail"); ok {
		request["IncludeOrganizationTrail"] = v
	}
	if v, ok := d.GetOkExists("include_shadow_trails"); ok {
		request["IncludeShadowTrails"] = v
	}
	var objects []map[string]interface{}
	var trailNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		trailNameRegex = r
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
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_actiontrail_trails", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.TrailList", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TrailList", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if trailNameRegex != nil {
			if !trailNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Name"])]; !ok {
				continue
			}
		}
		if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"event_rw":              object["EventRW"],
			"is_organization_trail": object["IsOrganizationTrail"],
			"oss_bucket_name":       object["OssBucketName"],
			"oss_key_prefix":        object["OssKeyPrefix"],
			"oss_write_role_arn":    object["OssWriteRoleArn"],
			"sls_project_arn":       object["SlsProjectArn"],
			"sls_write_role_arn":    object["SlsWriteRoleArn"],
			"status":                object["Status"],
			"id":                    fmt.Sprint(object["Name"]),
			"trail_name":            fmt.Sprint(object["Name"]),
			"trail_region":          object["TrailRegion"],
		}
		ids = append(ids, fmt.Sprint(object["Name"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("actiontrails", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("trails", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
