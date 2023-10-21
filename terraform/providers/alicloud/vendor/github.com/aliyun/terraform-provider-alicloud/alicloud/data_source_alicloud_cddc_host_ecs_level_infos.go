package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCddcHostEcsLevelInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCddcHostEcsLevelInfosRead,
		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mysql", "mssql", "pgsql", "redis"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3"}, false),
			},
			"image_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"WindowsWithMssqlEntAlwaysonLicense", "WindowsWithMssqlStdLicense", "WindowsWithMssqlEntLicense", "WindowsWithMssqlWebLicense", "AliLinux"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"res_class_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_class_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ecs_class": {
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
		},
	}
}

func dataSourceAlicloudCddcHostEcsLevelInfosRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeHostEcsLevelInfo"
	request := make(map[string]interface{})
	request["DbType"] = d.Get("db_type")
	request["RegionId"] = client.RegionId
	request["ZoneId"] = d.Get("zone_id")
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("image_category"); ok {
		request["ImageCategory"] = v
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewCddcClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-03-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cddc_host_ecs_level_infos", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.HostEcsLevelInfos", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.HostEcsLevelInfos", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		for _, item := range object["Items"].([]interface{}) {
			itemArr := item.(map[string]interface{})
			mapping := map[string]interface{}{
				"description":    itemArr["Description"],
				"res_class_code": itemArr["RdsClassCode"],
				"ecs_class_code": itemArr["EcsClassCode"],
				"ecs_class":      itemArr["EcsClass"],
			}
			s = append(s, mapping)
		}
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("infos", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
