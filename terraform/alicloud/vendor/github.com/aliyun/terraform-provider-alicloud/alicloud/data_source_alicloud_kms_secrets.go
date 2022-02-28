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

func dataSourceAlicloudKmsSecrets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsSecretsRead,
		Schema: map[string]*schema.Schema{
			"fetch_tags": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"filters": {
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
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"planned_delete_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudKmsSecretsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListSecrets"
	request := make(map[string]interface{})
	if v, ok := d.GetOkExists("fetch_tags"); ok {
		request["FetchTags"] = v
	}
	if v, ok := d.GetOk("filters"); ok {
		request["Filters"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var secretNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		secretNameRegex = r
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
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	var response map[string]interface{}
	conn, err := client.NewKmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-01-20"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_secrets", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.SecretList.Secret", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.SecretList.Secret", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if secretNameRegex != nil {
				if !secretNameRegex.MatchString(fmt.Sprint(item["SecretName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SecretName"])]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if len(item["Tags"].(map[string]interface{})["Tag"].([]interface{})) != len(tagsMap) {
					continue
				}
				match := true
				for _, tag := range item["Tags"].(map[string]interface{})["Tag"].([]interface{}) {
					if v, ok := tagsMap[tag.(map[string]interface{})["TagKey"].(string)]; !ok || v.(string) != tag.(map[string]interface{})["TagValue"].(string) {
						match = false
						break
					}
				}
				if !match {
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
			"planned_delete_time": object["PlannedDeleteTime"],
			"id":                  fmt.Sprint(object["SecretName"]),
			"secret_name":         fmt.Sprint(object["SecretName"]),
			"secret_type":         object["SecretType"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["SecretName"]))
			names = append(names, object["SecretName"])
			s = append(s, mapping)
			continue
		}

		kmsService := KmsService{client}
		id := fmt.Sprint(object["SecretName"])
		getResp, err := kmsService.DescribeKmsSecret(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["arn"] = getResp["Arn"]
		mapping["description"] = getResp["Description"]
		mapping["encryption_key_id"] = getResp["EncryptionKeyId"]
		getResp1, err := kmsService.GetSecretValue(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["secret_data"] = getResp1["SecretData"]
		mapping["secret_data_type"] = getResp1["SecretDataType"]
		mapping["version_id"] = getResp1["VersionId"]
		mapping["version_stages"] = getResp1["VersionStages"].(map[string]interface{})["VersionStage"]

		ids = append(ids, fmt.Sprint(object["SecretName"]))
		names = append(names, object["SecretName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("secrets", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
