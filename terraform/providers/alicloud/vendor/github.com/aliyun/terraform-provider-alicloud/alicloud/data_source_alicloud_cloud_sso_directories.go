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

func dataSourceAlicloudCloudSsoDirectories() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudSsoDirectoriesRead,
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"directory_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mfa_authentication_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"saml_identity_provider_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"create_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"encoded_metadata_document": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entity_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"login_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"sso_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"scim_synchronization_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tasks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_configuration_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"access_configuration_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"failure_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"principal_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"principal_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"principal_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"task_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"task_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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

func dataSourceAlicloudCloudSsoDirectoriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDirectories"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var directoryNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		directoryNameRegex = r
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_directories", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Directories", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Directories", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if directoryNameRegex != nil && !directoryNameRegex.MatchString(fmt.Sprint(item["DirectoryName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["DirectoryId"])]; !ok {
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
			"create_time":    object["CreateTime"],
			"id":             fmt.Sprint(object["DirectoryId"]),
			"directory_id":   fmt.Sprint(object["DirectoryId"]),
			"directory_name": object["DirectoryName"],
			"region":         object["Region"],
		}

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["DirectoryName"])

		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}

		cloudssoService := CloudssoService{client}
		getDirectoryObject, err := cloudssoService.DescribeCloudSsoDirectory(fmt.Sprint(object["DirectoryId"]))
		if err != nil {
			return WrapError(err)
		}
		mapping["mfa_authentication_status"] = getDirectoryObject["MFAAuthenticationStatus"]
		mapping["scim_synchronization_status"] = getDirectoryObject["SCIMSynchronizationStatus"]
		if SAMLIdentityProviderConfiguration, ok := getDirectoryObject["SAMLIdentityProviderConfiguration"]; ok && len(SAMLIdentityProviderConfiguration.(map[string]interface{})) > 0 {
			SAMLIdentityProviderConfigurationSli := make([]map[string]interface{}, 0)
			SAMLIdentityProviderConfigurationMap := make(map[string]interface{})
			SAMLIdentityProviderConfigurationMap["sso_status"] = SAMLIdentityProviderConfiguration.(map[string]interface{})["SSOStatus"]
			SAMLIdentityProviderConfigurationMap["entity_id"] = SAMLIdentityProviderConfiguration.(map[string]interface{})["EntityId"]
			SAMLIdentityProviderConfigurationMap["create_time"] = SAMLIdentityProviderConfiguration.(map[string]interface{})["CreateTime"]
			SAMLIdentityProviderConfigurationMap["login_url"] = SAMLIdentityProviderConfiguration.(map[string]interface{})["LoginUrl"]
			if v, ok := SAMLIdentityProviderConfiguration.(map[string]interface{})["EncodedMetadataDocument"]; ok {
				SAMLIdentityProviderConfigurationMap["encoded_metadata_document"] = v
			}
			SAMLIdentityProviderConfigurationSli = append(SAMLIdentityProviderConfigurationSli, SAMLIdentityProviderConfigurationMap)
			mapping["saml_identity_provider_configuration"] = SAMLIdentityProviderConfigurationSli
		}

		getResp, err := cloudssoService.GetDirectoryTasks(fmt.Sprint(object["DirectoryId"]))
		if err != nil {
			return WrapError(err)
		}

		if getDirectoryTasks, ok := getResp["Tasks"]; ok {
			tasks := make([]map[string]interface{}, 0)
			for _, v := range getDirectoryTasks.([]interface{}) {
				if t, ok := v.(map[string]interface{}); ok {
					temp := map[string]interface{}{
						"access_configuration_id":   t["AccessConfigurationId"],
						"access_configuration_name": t["AccessConfigurationName"],
						"end_time":                  t["EndTime"],
						"failure_reason":            t["FailureReason"],
						"principal_id":              t["PrincipalId"],
						"principal_name":            t["PrincipalName"],
						"principal_type":            t["PrincipalType"],
						"start_time":                t["StartTime"],
						"status":                    t["Status"],
						"target_id":                 t["TargetId"],
						"target_name":               t["TargetName"],
						"target_path":               t["TargetPath"],
						"target_type":               t["TargetType"],
						"task_id":                   t["TaskId"],
						"task_type":                 t["TaskType"],
					}
					tasks = append(tasks, temp)
				}
			}
			mapping["tasks"] = tasks
		}

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("directories", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
