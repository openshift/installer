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

func dataSourceAlicloudRosTemplateScratches() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRosTemplateScratchesRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"GENERATE_COMPLETE", "GENERATE_FAILED", "GENERATE_IN_PROGRESS"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"template_scratch_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ResourceImport", "ArchitectureReplication"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scratches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"logical_id_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preference_parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"parameter_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"parameter_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"source_tag": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_tags": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"resource_type_filter": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"source_resource_group": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type_filter": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"source_resources": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"stacks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"stack_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_scratch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_scratch_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudRosTemplateScratchesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListTemplateScratches"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = convertListToJsonString([]interface{}{v})
	}
	if v, ok := d.GetOk("template_scratch_type"); ok {
		request["TemplateScratchType"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
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
	var response map[string]interface{}
	conn, err := client.NewRosClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-09-10"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ros_template_scratches", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.TemplateScratches", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.TemplateScratches", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TemplateScratchId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":           object["CreateTime"],
			"description":           object["Description"],
			"logical_id_strategy":   object["LogicalIdStrategy"],
			"status":                object["Status"],
			"id":                    fmt.Sprint(object["TemplateScratchId"]),
			"template_scratch_id":   fmt.Sprint(object["TemplateScratchId"]),
			"template_scratch_type": object["TemplateScratchType"],
		}

		preferenceParameters := make([]map[string]interface{}, 0)
		if preferenceParametersList, ok := object["PreferenceParameters"].([]interface{}); ok {
			for _, v := range preferenceParametersList {
				if m1, ok := v.(map[string]interface{}); ok {
					if v, ok := m1["ParameterValue"]; ok && fmt.Sprint(v) != "" {
						temp1 := map[string]interface{}{
							"parameter_key":   m1["ParameterKey"],
							"parameter_value": v,
						}
						preferenceParameters = append(preferenceParameters, temp1)
					}
				}
			}
		}
		mapping["preference_parameters"] = preferenceParameters

		sourceResourceGroupSli := make([]map[string]interface{}, 0)
		if v, ok := object["SourceResourceGroup"]; ok {
			if sourceResourceGroup, ok := v.(map[string]interface{}); ok && len(sourceResourceGroup) > 0 {
				sourceResourceGroupMap := make(map[string]interface{})
				sourceResourceGroupMap["resource_group_id"] = sourceResourceGroup["ResourceGroupId"]
				resourceTypeFilter := make([]interface{}, 0)
				if v, ok := sourceResourceGroup["ResourceTypeFilter"]; ok {
					if vv, ok := v.([]interface{}); ok && len(vv) > 0 {
						resourceTypeFilter = append(resourceTypeFilter, vv...)
					}
				}
				sourceResourceGroupMap["resource_type_filter"] = resourceTypeFilter
				sourceResourceGroupSli = append(sourceResourceGroupSli, sourceResourceGroupMap)
			}
		}
		mapping["source_resource_group"] = sourceResourceGroupSli

		sourceResources := make([]map[string]interface{}, 0)
		if sourceResourcesList, ok := object["SourceResources"].([]interface{}); ok {
			for _, v := range sourceResourcesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"resource_id":   m1["ResourceId"],
						"resource_type": m1["ResourceType"],
					}
					sourceResources = append(sourceResources, temp1)
				}
			}
		}
		mapping["source_resources"] = sourceResources

		sourceTagSli := make([]map[string]interface{}, 0)
		if v, ok := object["SourceTag"]; ok {
			if sourceTag, ok := v.(map[string]interface{}); ok && len(sourceTag) > 0 {
				sourceTagMap := make(map[string]interface{})
				sourceTagMap["resource_tags"] = sourceTag["ResourceTags"]
				resourceTypeFilter := make([]interface{}, 0)
				if v, ok := sourceTag["ResourceTypeFilter"]; ok {
					if vv, ok := v.([]interface{}); ok && len(vv) > 0 {
						resourceTypeFilter = append(resourceTypeFilter, vv...)
					}
				}
				sourceTagMap["resource_type_filter"] = resourceTypeFilter
				sourceTagSli = append(sourceTagSli, sourceTagMap)
			}
		}
		mapping["source_tag"] = sourceTagSli

		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["TemplateScratchId"])
		rosService := RosService{client}
		getResp, err := rosService.DescribeRosTemplateScratch(id)
		if err != nil {
			return WrapError(err)
		}

		stacks := make([]map[string]interface{}, 0)
		if stacksList, ok := getResp["Stacks"].([]interface{}); ok {
			for _, v := range stacksList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"stack_id": m1["StackId"],
					}
					stacks = append(stacks, temp1)
				}
			}
		}
		mapping["stacks"] = stacks
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("scratches", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
