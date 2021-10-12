package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudResourceManagerPolicyAttachments() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudResourceManagerPolicyAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"language": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Custom", "System"}, false),
			},
			"principal_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"principal_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IMSGroup", "IMSUser", "ServiceRole"}, false),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attach_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_type": {
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
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudResourceManagerPolicyAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListPolicyAttachments"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("language"); ok {
		request["Language"] = v
	}
	if v, ok := d.GetOk("policy_name"); ok {
		request["PolicyName"] = v
	}
	if v, ok := d.GetOk("policy_type"); ok {
		request["PolicyType"] = v
	}
	if v, ok := d.GetOk("principal_name"); ok {
		request["PrincipalName"] = v
	}
	if v, ok := d.GetOk("principal_type"); ok {
		request["PrincipalType"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewResourcemanagerClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-31"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_resource_manager_policy_attachments", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.PolicyAttachments.PolicyAttachment", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.PolicyAttachments.PolicyAttachment", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["PolicyName"], ":", object["PolicyType"], ":", object["PrincipalName"], ":", object["PrincipalType"], ":", object["ResourceGroupId"]),
			"attach_date":       object["AttachDate"],
			"description":       object["Description"],
			"policy_name":       object["PolicyName"],
			"policy_type":       object["PolicyType"],
			"principal_name":    object["PrincipalName"],
			"principal_type":    object["PrincipalType"],
			"resource_group_id": object["ResourceGroupId"],
		}
		ids = append(ids, fmt.Sprint(object["PolicyName"], ":", object["PolicyType"], ":", object["PrincipalName"], ":", object["PrincipalType"], ":", object["ResourceGroupId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("attachments", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
