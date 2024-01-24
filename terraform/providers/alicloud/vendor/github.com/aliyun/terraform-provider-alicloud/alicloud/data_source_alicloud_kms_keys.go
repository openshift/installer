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

func dataSourceAlicloudKmsKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKmsKeysRead,
		Schema: map[string]*schema.Schema{
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"filters": {
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled", "PendingDeletion"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"automatic_rotation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_rotation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"material_expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_rotation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rotation_interval": {
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
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func dataSourceAlicloudKmsKeysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListKeys"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		request["Filters"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var descriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		descriptionRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kms_keys", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Keys.Key", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Keys.Key", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["KeyId"])]; !ok {
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
			"arn":    object["KeyArn"],
			"id":     fmt.Sprint(object["KeyId"]),
			"key_id": fmt.Sprint(object["KeyId"]),
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["KeyId"]))
			s = append(s, mapping)
			continue
		}

		kmsService := KmsService{client}
		id := fmt.Sprint(object["KeyId"])
		getResp, err := kmsService.DescribeKmsKey(id)
		if _, ok := getResp["KeyState"]; !ok && err != nil {
			return WrapError(err)
		}
		if descriptionRegex != nil {
			if !descriptionRegex.MatchString(fmt.Sprint(getResp["Description"])) {
				continue
			}
		}
		if statusOk && status != "" && status != getResp["KeyState"].(string) {
			continue
		}
		mapping["automatic_rotation"] = getResp["AutomaticRotation"]
		mapping["creator"] = getResp["Creator"]
		mapping["creation_date"] = getResp["CreationDate"]
		mapping["delete_date"] = getResp["DeleteDate"]
		mapping["description"] = getResp["Description"]
		mapping["key_spec"] = getResp["KeySpec"]
		mapping["key_usage"] = getResp["KeyUsage"]
		mapping["last_rotation_date"] = getResp["LastRotationDate"]
		mapping["material_expire_time"] = getResp["MaterialExpireTime"]
		mapping["next_rotation_date"] = getResp["NextRotationDate"]
		mapping["origin"] = getResp["Origin"]
		mapping["primary_key_version"] = getResp["PrimaryKeyVersion"]
		mapping["protection_level"] = getResp["ProtectionLevel"]
		mapping["rotation_interval"] = getResp["RotationInterval"]
		mapping["status"] = getResp["KeyState"]
		ids = append(ids, fmt.Sprint(object["KeyId"]))
		names = append(names)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("keys", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
