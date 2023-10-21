package alicloud

import (
	"fmt"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudConfigConfigurationRecorders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigConfigurationRecordersRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recorders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_enable_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization_master_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceAlicloudConfigConfigurationRecordersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeConfigurationRecorder"
	request := make(map[string]interface{})
	var response map[string]interface{}
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2019-01-08"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_configuration_recorders", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	s := make([]map[string]interface{}, 0)
	mapping := map[string]interface{}{
		"id":                         fmt.Sprint(formatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])),
		"account_id":                 fmt.Sprint(formatInt(response["ConfigurationRecorder"].(map[string]interface{})["AccountId"])),
		"organization_enable_status": response["ConfigurationRecorder"].(map[string]interface{})["OrganizationEnableStatus"],
		"organization_master_id":     formatInt(response["ConfigurationRecorder"].(map[string]interface{})["OrganizationMasterId"]),
		"resource_types":             response["ConfigurationRecorder"].(map[string]interface{})["ResourceTypes"],
		"status":                     response["ConfigurationRecorder"].(map[string]interface{})["ConfigurationRecorderStatus"],
	}
	s = append(s, mapping)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

	if err := d.Set("recorders", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
