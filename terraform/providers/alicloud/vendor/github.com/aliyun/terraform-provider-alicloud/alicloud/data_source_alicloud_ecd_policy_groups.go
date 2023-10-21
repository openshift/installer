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

func dataSourceAlicloudEcdPolicyGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcdPolicyGroupsRead,
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
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AVAILABLE", "CREATING"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorize_access_policy_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"authorize_security_policy_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cidr_ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_protocol": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port_range": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"clipboard": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eds_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"html_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"html_file_transfer": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_drive": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"usb_redirect": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"visual_quality": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"watermark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"watermark_transparency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"watermark_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcdPolicyGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribePolicyGroups"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var policyGroupNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		policyGroupNameRegex = r
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
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecd_policy_groups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DescribePolicyGroups", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DescribePolicyGroups", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if policyGroupNameRegex != nil && !policyGroupNameRegex.MatchString(fmt.Sprint(item["Name"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PolicyGroupId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["PolicyStatus"].(string) {
				continue
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
			"clipboard":              object["Clipboard"],
			"domain_list":            object["DomainList"],
			"eds_count":              formatInt(object["EdsCount"]),
			"html_access":            object["Html5Access"],
			"html_file_transfer":     object["Html5FileTransfer"],
			"local_drive":            object["LocalDrive"],
			"id":                     fmt.Sprint(object["PolicyGroupId"]),
			"policy_group_id":        fmt.Sprint(object["PolicyGroupId"]),
			"policy_group_name":      object["Name"],
			"policy_group_type":      object["PolicyGroupType"],
			"status":                 object["PolicyStatus"],
			"usb_redirect":           object["UsbRedirect"],
			"visual_quality":         object["VisualQuality"],
			"watermark":              object["Watermark"],
			"watermark_transparency": object["WatermarkTransparency"],
			"watermark_type":         object["WatermarkType"],
		}

		authorizeAccessPolicyRules := make([]map[string]interface{}, 0)
		if authorizeAccessPolicyRulesList, ok := object["AuthorizeAccessPolicyRules"].([]interface{}); ok {
			for _, v := range authorizeAccessPolicyRulesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"cidr_ip":     m1["CidrIp"],
						"description": m1["Description"],
					}
					authorizeAccessPolicyRules = append(authorizeAccessPolicyRules, temp1)
				}
			}
		}
		mapping["authorize_access_policy_rules"] = authorizeAccessPolicyRules
		if v, ok := object["AuthorizeSecurityPolicyRules"].([]interface{}); ok {
			authorizeSecurityPolicyRules := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"cidr_ip":     item["CidrIp"],
					"description": item["Description"],
					"ip_protocol": item["IpProtocol"],
					"policy":      item["Policy"],
					"port_range":  item["PortRange"],
					"priority":    item["Priority"],
					"type":        item["Type"],
				}

				authorizeSecurityPolicyRules = append(authorizeSecurityPolicyRules, temp)
			}
			mapping["authorize_security_policy_rules"] = authorizeSecurityPolicyRules
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("groups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
