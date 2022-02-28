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

func dataSourceAlicloudCloudSsoUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudSsoUsersRead,
		Schema: map[string]*schema.Schema{
			"directory_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provision_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Manual", "Synchronized"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled"}, false),
			},
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
			"users": {
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
						"directory_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mfa_devices": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"effective_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"provision_type": {
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
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
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

func dataSourceAlicloudCloudSsoUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListUsers"
	request := make(map[string]interface{})
	request["DirectoryId"] = d.Get("directory_id")
	if v, ok := d.GetOk("provision_type"); ok {
		request["ProvisionType"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var userNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		userNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_sso_users", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Users", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Users", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if userNameRegex != nil && !userNameRegex.MatchString(fmt.Sprint(item["UserName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprintf("%s:%s", request["DirectoryId"], item["UserId"])]; !ok {
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
			"create_time":    object["CreateTime"],
			"description":    object["Description"],
			"directory_id":   request["DirectoryId"],
			"display_name":   object["DisplayName"],
			"email":          object["Email"],
			"first_name":     object["FirstName"],
			"last_name":      object["LastName"],
			"provision_type": object["ProvisionType"],
			"status":         object["Status"],
			"id":             fmt.Sprint(request["DirectoryId"], ":", object["UserId"]),
			"user_id":        fmt.Sprint(object["UserId"]),
			"user_name":      object["UserName"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["UserName"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(request["DirectoryId"], ":", object["UserId"])
		cloudssoService := CloudssoService{client}
		getResp, err := cloudssoService.ListMFADevicesForUser(id)
		if err != nil {
			return WrapError(err)
		}

		mFADevices := make([]map[string]interface{}, 0)
		if mFADevicesList, ok := getResp["MFADevices"]; ok {
			for _, v := range mFADevicesList.([]interface{}) {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"device_id":      m1["DeviceId"],
						"device_name":    m1["DeviceName"],
						"device_type":    m1["DeviceType"],
						"effective_time": m1["EffectiveTime"],
					}
					mFADevices = append(mFADevices, temp1)
				}
			}
		}
		mapping["mfa_devices"] = mFADevices
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("users", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
