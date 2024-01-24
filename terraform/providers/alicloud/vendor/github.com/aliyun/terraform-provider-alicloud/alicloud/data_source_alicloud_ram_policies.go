package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRamPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRamPoliciesRead,
		Schema: map[string]*schema.Schema{
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
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
			},
			"group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"role_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 64),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attachment_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"document": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_date": {
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
				Default:  true,
			},
		},
	}
}

func dataSourceAlicloudRamPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPolicies"
	request := make(map[string]interface{})
	request["MaxItems"] = PageSizeLarge
	var objects []map[string]interface{}
	userFilterPoliciesMap := make(map[string]interface{})
	groupFilterPoliciesMap := make(map[string]interface{})
	roleFilterPoliciesMap := make(map[string]interface{})

	dataMap := []map[string]interface{}{}
	userName, userNameOk := d.GetOk("user_name")
	groupName, groupNameOk := d.GetOk("group_name")
	roleName, roleNameOk := d.GetOk("role_name")
	policyType, policyTypeOk := d.GetOk("type")

	var policyNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		policyNameRegex = r
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
	conn, err := client.NewRamClient()
	if err != nil {
		return WrapError(err)
	}
	// policies for user
	if userNameOk {
		userAction := "ListPoliciesForUser"
		userRequest := map[string]interface{}{
			"UserName": userName,
		}
		response, err = conn.DoRequest(StringPointer(userAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, userRequest, &util.RuntimeOptions{})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies")
		}

		userResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		for _, v := range userResp.([]interface{}) {
			userFilterPoliciesMap[v.(map[string]interface{})["PolicyType"].(string)+v.(map[string]interface{})["PolicyName"].(string)] = v
		}
		dataMap = append(dataMap, userFilterPoliciesMap)
	}

	// policies for group
	if groupNameOk {
		groupAction := "ListPoliciesForGroup"
		groupRequest := map[string]interface{}{
			"GroupName": groupName,
		}
		response, err = conn.DoRequest(StringPointer(groupAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, groupRequest, &util.RuntimeOptions{})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies")
		}
		groupResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		for _, v := range groupResp.([]interface{}) {
			groupFilterPoliciesMap[v.(map[string]interface{})["PolicyType"].(string)+v.(map[string]interface{})["PolicyName"].(string)] = v
		}
		dataMap = append(dataMap, groupFilterPoliciesMap)
	}

	// policies for role
	if roleNameOk {
		roleAction := "ListPoliciesForRole"
		roleRequest := map[string]interface{}{
			"RoleName": roleName,
		}
		response, err = conn.DoRequest(StringPointer(roleAction), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, roleRequest, &util.RuntimeOptions{})

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies")
		}
		roleResp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		for _, v := range roleResp.([]interface{}) {
			roleFilterPoliciesMap[v.(map[string]interface{})["PolicyType"].(string)+v.(map[string]interface{})["PolicyName"].(string)] = v
		}
		dataMap = append(dataMap, roleFilterPoliciesMap)
	}

	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ram_policies", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Policies.Policy", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Policies.Policy", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if policyNameRegex != nil {
				if !policyNameRegex.MatchString(fmt.Sprint(item["PolicyName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["PolicyName"])]; !ok {
					continue
				}
			}
			if policyTypeOk && policyType.(string) != item["PolicyType"] {
				continue
			}

			if len(dataMap) > 0 {
				res := false
				for _, v := range dataMap {
					if _, ok := v[item["PolicyType"].(string)+item["PolicyName"].(string)]; ok {
						res = true
						break
					}
				}
				if !res {
					continue
				}
			}

			objects = append(objects, item)
		}
		if marker, ok := response["Marker"].(string); ok && marker != "" {
			request["Marker"] = marker
		} else {
			break
		}
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"attachment_count": formatInt(object["AttachmentCount"]),
			"default_version":  object["DefaultVersion"],
			"description":      object["Description"],
			"id":               fmt.Sprint(object["PolicyName"]),
			"policy_name":      fmt.Sprint(object["PolicyName"]),
			"name":             object["PolicyName"],
			"update_date":      object["UpdateDate"],
			"type":             object["PolicyType"],
			"create_date":      object["CreateDate"],
			"user_name":        object["UserName"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["PolicyName"]))
			names = append(names, object["PolicyName"])
			s = append(s, mapping)
			continue
		}

		id := fmt.Sprint(object["PolicyName"])
		action := "GetPolicy"
		request := map[string]interface{}{
			"PolicyName": id,
			"PolicyType": object["PolicyType"],
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2015-05-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapError(err)
		}
		addDebug(action, response, request)
		v, err := jsonpath.Get("$", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, id, "$", response)
		}
		getResp := v.(map[string]interface{})

		mapping["policy_document"] = getResp["DefaultPolicyVersion"].(map[string]interface{})["PolicyDocument"]
		mapping["document"] = getResp["DefaultPolicyVersion"].(map[string]interface{})["PolicyDocument"]
		mapping["version_id"] = getResp["DefaultPolicyVersion"].(map[string]interface{})["VersionId"]
		ids = append(ids, fmt.Sprint(object["PolicyName"]))
		names = append(names, object["PolicyName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
