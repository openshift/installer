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

func dataSourceAlicloudAlbHealthCheckTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbHealthCheckTemplatesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"health_check_template_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"health_check_template_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][A-Za-z0-9\\_\\.\\-]{1,127}$`), "The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). The name must start with a letter."),
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
						"health_check_codes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"health_check_connect_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_http_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"healthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unhealthy_threshold": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlbHealthCheckTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListHealthCheckTemplates"
	request := make(map[string]interface{})
	if m, ok := d.GetOk("health_check_template_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("HealthCheckTemplateIds.%d", k+1)] = v.(string)
		}
	}
	if v, ok := d.GetOk("health_check_template_name"); ok {
		request["HealthCheckTemplateNames.1"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var healthCheckTemplateNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		healthCheckTemplateNameRegex = r
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
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_health_check_templates", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.HealthCheckTemplates", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.HealthCheckTemplates", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if healthCheckTemplateNameRegex != nil && !healthCheckTemplateNameRegex.MatchString(fmt.Sprint(item["HealthCheckTemplateName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["HealthCheckTemplateId"])]; !ok {
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
			"health_check_codes":         object["HealthCheckCodes"],
			"health_check_connect_port":  formatInt(object["HealthCheckConnectPort"]),
			"health_check_host":          object["HealthCheckHost"],
			"health_check_http_version":  object["HealthCheckHttpVersion"],
			"health_check_interval":      formatInt(object["HealthCheckInterval"]),
			"health_check_method":        object["HealthCheckMethod"],
			"health_check_path":          object["HealthCheckPath"],
			"health_check_protocol":      object["HealthCheckProtocol"],
			"id":                         fmt.Sprint(object["HealthCheckTemplateId"]),
			"health_check_template_id":   fmt.Sprint(object["HealthCheckTemplateId"]),
			"health_check_template_name": object["HealthCheckTemplateName"],
			"health_check_timeout":       formatInt(object["HealthCheckTimeout"]),
			"healthy_threshold":          formatInt(object["HealthyThreshold"]),
			"unhealthy_threshold":        formatInt(object["UnhealthyThreshold"]),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["HealthCheckTemplateName"])
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
