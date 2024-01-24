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

func dataSourceAlicloudDfsFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDfsFileSystemsRead,
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
			"systems": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"file_system_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_point_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"number_of_directories": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"number_of_files": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"provisioned_throughput_in_mi_bps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"space_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_package_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"throughput_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"used_space_size": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudDfsFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListFileSystems"
	request := make(map[string]interface{})
	request["InputRegionId"] = client.RegionId
	var objects []map[string]interface{}
	var fileSystemNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		fileSystemNameRegex = r
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
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dfs_file_systems", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.FileSystems", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.FileSystems", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if fileSystemNameRegex != nil && !fileSystemNameRegex.MatchString(fmt.Sprint(item["FileSystemName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["FileSystemId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":                      object["CreateTime"],
			"description":                      object["Description"],
			"id":                               fmt.Sprint(object["FileSystemId"]),
			"file_system_id":                   fmt.Sprint(object["FileSystemId"]),
			"file_system_name":                 object["FileSystemName"],
			"mount_point_count":                object["MountPointCount"],
			"number_of_directories":            object["NumberOfDirectories"],
			"number_of_files":                  object["NumberOfFiles"],
			"protocol_type":                    object["ProtocolType"],
			"provisioned_throughput_in_mi_bps": object["ProvisionedThroughputInMiBps"],
			"space_capacity":                   object["SpaceCapacity"],
			"storage_package_id":               object["StoragePackageId"],
			"storage_type":                     object["StorageType"],
			"throughput_mode":                  object["ThroughputMode"],
			"used_space_size":                  object["UsedSpaceSize"],
			"zone_id":                          object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["FileSystemName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("systems", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
