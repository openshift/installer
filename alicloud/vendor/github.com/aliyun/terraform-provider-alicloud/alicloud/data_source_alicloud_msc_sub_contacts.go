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

func dataSourceAlicloudMscSubContacts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMscSubContactsRead,
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
			"contacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"contact_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_account": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_obsolete": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_verified_email": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_verified_mobile": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"last_email_verification_time_stamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_mobile_verification_time_stamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mobile": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"position": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMscSubContactsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListContacts"
	request := make(map[string]interface{})
	request["Locale"] = "en"
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var contactNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		contactNameRegex = r
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
	conn, err := client.NewMscopensubscriptionClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2021-07-13"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_msc_sub_contacts", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		resp, err := jsonpath.Get("$.Contacts", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Contacts", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if contactNameRegex != nil && !contactNameRegex.MatchString(fmt.Sprint(item["ContactName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ContactId"])]; !ok {
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
			"account_uid":                         fmt.Sprint(object["AccountUid"]),
			"id":                                  fmt.Sprint(object["ContactId"]),
			"contact_id":                          fmt.Sprint(object["ContactId"]),
			"contact_name":                        object["ContactName"],
			"email":                               object["Email"],
			"is_account":                          object["IsAccount"],
			"is_obsolete":                         object["IsObsolete"],
			"is_verified_email":                   object["IsVerifiedEmail"],
			"is_verified_mobile":                  object["IsVerifiedMobile"],
			"last_email_verification_time_stamp":  fmt.Sprint(object["LastEmailVerificationTimeStamp"]),
			"last_mobile_verification_time_stamp": fmt.Sprint(object["LastMobileVerificationTimeStamp"]),
			"mobile":                              object["Mobile"],
			"position":                            convertMscSubContactPositionResponse(fmt.Sprint(object["Position"])),
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ContactName"])

		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("contacts", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
