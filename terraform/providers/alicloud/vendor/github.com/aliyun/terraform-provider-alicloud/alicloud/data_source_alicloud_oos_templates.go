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

func dataSourceAlicloudOosTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudOosTemplatesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_date": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_date_after": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"has_trigger": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"share_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Private", "Public"}, false),
			},
			"sort_field": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TotalExecutionCount",
				ValidateFunc: validation.StringInSlice([]string{"CreatedDate", "Popularity", "TemplateName", "TotalExecutionCount"}, false),
			},
			"sort_order": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Descending",
				ValidateFunc: validation.StringInSlice([]string{"Ascending", "Descending"}, false),
			},
			"tags": tagsSchema(),
			"template_format": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON", "YAML"}, false),
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
			"template_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"has_trigger": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"template_format": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudOosTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTemplates"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOk("created_by"); ok {
		request["CreatedBy"] = v
	}
	if v, ok := d.GetOk("created_date"); ok {
		request["CreatedDateBefore"] = v
	}
	if v, ok := d.GetOk("created_date_after"); ok {
		request["CreatedDateAfter"] = v
	}
	if v, ok := d.GetOkExists("has_trigger"); ok {
		request["HasTrigger"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("share_type"); ok {
		request["ShareType"] = v
	}
	if v, ok := d.GetOk("sort_field"); ok {
		request["SortField"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	if v, ok := d.GetOk("template_format"); ok {
		request["TemplateFormat"] = v
	}
	if v, ok := d.GetOk("template_type"); ok {
		request["TemplateType"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var templateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		templateNameRegex = r
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
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_oos_templates", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Templates", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Templates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if templateNameRegex != nil {
				if !templateNameRegex.MatchString(fmt.Sprint(item["TemplateName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TemplateName"])]; !ok {
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
			"category":         object["Category"],
			"created_by":       object["CreatedBy"],
			"created_date":     object["CreatedDate"],
			"description":      object["Description"],
			"has_trigger":      object["HasTrigger"],
			"share_type":       object["ShareType"],
			"tags":             object["Tags"],
			"template_format":  object["TemplateFormat"],
			"template_id":      object["TemplateId"],
			"id":               fmt.Sprint(object["TemplateName"]),
			"template_name":    fmt.Sprint(object["TemplateName"]),
			"template_type":    object["TemplateType"],
			"template_version": object["TemplateVersion"],
			"updated_by":       object["UpdatedBy"],
			"updated_date":     object["UpdatedDate"],
		}
		ids = append(ids, fmt.Sprint(object["TemplateName"]))
		names = append(names, object["TemplateName"])
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
