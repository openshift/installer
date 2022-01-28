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

func dataSourceAlicloudCloudSsoAccessConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudSsoAccessConfigurationsRead,
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
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_policies": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"add_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"permission_policy_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"permission_policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"permission_policy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"relay_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status_notifications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAlicloudCloudSsoAccessConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAccessConfigurations"
	request := make(map[string]interface{})
	request["DirectoryId"] = d.Get("directory_id")
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var accessConfigurationNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		accessConfigurationNameRegex = r
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
	conn, err := client.NewCloudssoClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2021-05-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_access_configurations", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.AccessConfigurations", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AccessConfigurations", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if accessConfigurationNameRegex != nil && !accessConfigurationNameRegex.MatchString(fmt.Sprint(item["AccessConfigurationName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(request["DirectoryId"], ":", item["AccessConfigurationId"])]; !ok {
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
			"id":                        fmt.Sprint(request["DirectoryId"], ":", object["AccessConfigurationId"]),
			"access_configuration_id":   fmt.Sprint(object["AccessConfigurationId"]),
			"access_configuration_name": object["AccessConfigurationName"],
			"create_time":               object["CreateTime"],
			"description":               object["Description"],
			"directory_id":              request["DirectoryId"],
			"relay_state":               object["RelayState"],
			"session_duration":          formatInt(object["SessionDuration"]),
			"status_notifications":      object["StatusNotifications"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["AccessConfigurationName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(request["DirectoryId"], ":", object["AccessConfigurationId"])
		cloudssoService := CloudssoService{client}
		getResp, err := cloudssoService.ListPermissionPoliciesInAccessConfiguration(id)
		if err != nil {
			return WrapError(err)
		}

		permissionPolicies := make([]map[string]interface{}, 0)
		if permissionPoliciesList, ok := getResp["PermissionPolicies"]; ok && permissionPoliciesList != nil {
			for _, v := range permissionPoliciesList.([]interface{}) {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"add_time":                   m1["AddTime"],
						"permission_policy_document": m1["PermissionPolicyDocument"],
						"permission_policy_name":     m1["PermissionPolicyName"],
						"permission_policy_type":     m1["PermissionPolicyType"],
					}
					permissionPolicies = append(permissionPolicies, temp1)
				}
			}
		}

		mapping["permission_policies"] = permissionPolicies
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("configurations", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
