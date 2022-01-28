package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSddpInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSddpInstancesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"odps_set": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_bucket_set": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rds_set": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSddpInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeUserStatus"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewSddpClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-03"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sddp_instances", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.UserStatus", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.UserStatus", response)
	}
	userstatus := resp.(map[string]interface{})
	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"authed":         userstatus["Authed"],
		"id":             fmt.Sprint(userstatus["InstanceId"]),
		"instance_id":    fmt.Sprint(userstatus["InstanceId"]),
		"instance_num":   fmt.Sprint(userstatus["InstanceNum"]),
		"odps_set":       userstatus["OdpsSet"],
		"oss_bucket_set": userstatus["OssBucketSet"],
		"oss_size":       fmt.Sprint(userstatus["OssSize"]),
		"payment_type":   convertSddpInstancePaymentTypeToStandard(userstatus["ChargeType"]),
		"rds_set":        userstatus["RdsSet"],
		"status":         fmt.Sprint(userstatus["InstanceStatus"]),
	}
	s = append(s, mapping)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
