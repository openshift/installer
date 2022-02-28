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

func dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumesRead,
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
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 8),
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cache_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chap_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"chap_in_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chunk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gateway_block_volume_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"index_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lun_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"oss_bucket_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oss_bucket_ssl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oss_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_download": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_upload": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_state": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudCloudStorageGatewayGatewayBlockVolumesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeGatewayBlockVolumes"
	request := make(map[string]interface{})
	request["GatewayId"] = d.Get("gateway_id")
	var objects []map[string]interface{}
	var gatewayBlockVolumeNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		gatewayBlockVolumeNameRegex = r
	}

	status, statusOk := d.GetOk("status")
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
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_block_volumes", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	resp, err := jsonpath.Get("$.BlockVolumes.BlockVolume", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.BlockVolumes.BlockVolume", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if gatewayBlockVolumeNameRegex != nil && !gatewayBlockVolumeNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(request["GatewayId"], ":", item["IndexId"])]; !ok {
				continue
			}
		}
		if statusOk && status.(int) != item["VolumeState"] {
			continue
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"address":                   object["Address"],
			"cache_mode":                object["CacheMode"],
			"chap_enabled":              object["ChapEnabled"],
			"chap_in_user":              fmt.Sprint(object["ChapInUser"]),
			"chunk_size":                formatInt(object["ChunkSize"]),
			"disk_id":                   object["DiskId"],
			"disk_type":                 object["DiskType"],
			"enabled":                   object["Enabled"],
			"gateway_block_volume_name": object["Name"],
			"gateway_id":                request["GatewayId"],
			"id":                        fmt.Sprint(request["GatewayId"], ":", object["IndexId"]),
			"index_id":                  fmt.Sprint(object["IndexId"]),
			"local_path":                object["LocalPath"],
			"lun_id":                    formatInt(object["LunId"]),
			"oss_bucket_name":           object["OssBucketName"],
			"oss_bucket_ssl":            object["OssBucketSsl"],
			"oss_endpoint":              object["OssEndpoint"],
			"port":                      formatInt(object["Port"]),
			"protocol":                  object["Protocol"],
			"size":                      formatInt(object["Size"]),
			"state":                     object["State"],
			"status":                    formatInt(object["Status"]),
			"target":                    object["Target"],
			"total_download":            formatInt(object["TotalDownload"]),
			"total_upload":              formatInt(object["TotalUpload"]),
			"volume_state":              formatInt(object["VolumeState"]),
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

	if err := d.Set("volumes", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
