package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudKmsSecretVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsSecretVersionsRead,
		Schema: map[string]*schema.Schema{
			"include_deprecated": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version_stage": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_stages": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudKmsSecretVersionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSecretVersionIds"
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok && len(v.([]interface{})) > 0 {
		for _, i := range v.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}
	request := make(map[string]interface{})

	if v, ok := d.GetOk("include_deprecated"); ok {
		request["IncludeDeprecated"] = v.(string)
	}
	VersionStage, okStage := d.GetOk("version_stage")

	request["SecretName"] = d.Get("secret_name")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var ids []string
	var objects []map[string]interface{}

	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, "alicloud_kms_secret_versions", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.VersionIds.VersionId", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.VersionIds.VersionId", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["VersionId"])]; !ok {
					continue
				}
			}
			if okStage && VersionStage.(string) != "" {
				hasVersionStage := false
				for _, VStage := range item["VersionStages"].(map[string]interface{})["VersionStage"].([]interface{}) {
					if VStage == VersionStage {
						hasVersionStage = true
						break
					}
				}
				if !hasVersionStage {
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

	secretName, err := jsonpath.Get("$.SecretName", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SecretName", response)
	}

	s := make([]map[string]interface{}, len(objects))
	for i, object := range objects {
		mapping := map[string]interface{}{
			"secret_name":    secretName,
			"version_id":     object["VersionId"],
			"version_stages": object["VersionStages"].(map[string]interface{})["VersionStage"].([]interface{}),
		}

		ids = append(ids, fmt.Sprint(object["VersionId"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s[i] = mapping
			continue
		}
		action := "GetSecretValue"
		var response map[string]interface{}
		request := make(map[string]interface{})
		request["RegionId"] = client.RegionId
		request["VersionId"] = object["VersionId"]
		if okStage && VersionStage.(string) != "" {
			request["VersionStage"] = VersionStage
		}
		request["SecretName"] = secretName
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_secret_versions", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		if v, err := jsonpath.Get("$.SecretData", response); err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SecretData", response)
		} else {
			mapping["secret_data"] = v
		}
		if v, err := jsonpath.Get("$.SecretDataType", response); err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SecretDataType", response)
		} else {
			mapping["secret_data_type"] = v
		}
		s[i] = mapping
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("versions", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
