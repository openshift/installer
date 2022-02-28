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

func dataSourceAlicloudCmsMetricRuleTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMetricRuleTemplatesRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"metric_rule_template_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_templates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"escalations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"critical": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"comparison_operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"statistics": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"threshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"times": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"info": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"comparison_operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"statistics": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"threshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"times": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"warn": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"comparison_operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"statistics": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"threshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"times": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"metric_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"selector": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"webhook": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_rule_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rest_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_id": {
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

func dataSourceAlicloudCmsMetricRuleTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeMetricRuleTemplateList"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("template_id"); ok {
		request["TemplateId"] = v
	}
	if v, ok := d.GetOk("metric_rule_template_name"); ok {
		request["Name"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var metricRuleTemplateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		metricRuleTemplateNameRegex = r
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
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_metric_rule_templates", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Templates.Template", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Templates.Template", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if metricRuleTemplateNameRegex != nil && !metricRuleTemplateNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["TemplateId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"description":               object["Description"],
			"metric_rule_template_name": object["Name"],
			"rest_version":              fmt.Sprint(object["RestVersion"]),
			"id":                        fmt.Sprint(object["TemplateId"]),
			"template_id":               fmt.Sprint(object["TemplateId"]),
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["TemplateId"])
		cmsService := CmsService{client}
		getResp, err := cmsService.DescribeCmsMetricRuleTemplate(id)
		if err != nil {
			return WrapError(err)
		}

		alertTemplate := make([]map[string]interface{}, 0)
		if alertTemplateList, ok := getResp["AlertTemplates"].(map[string]interface{})["AlertTemplate"].([]interface{}); ok {
			for _, v := range alertTemplateList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"category":    m1["Category"],
						"metric_name": m1["MetricName"],
						"namespace":   m1["Namespace"],
						"rule_name":   m1["RuleName"],
						"selector":    m1["Selector"],
						"webhook":     m1["Webhook"],
					}
					alertTemplate = append(alertTemplate, temp1)
				}
			}
		}
		mapping["alert_templates"] = alertTemplate
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
