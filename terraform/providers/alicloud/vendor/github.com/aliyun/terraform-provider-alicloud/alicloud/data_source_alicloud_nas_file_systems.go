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

func dataSourceAlicloudFileSystems() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFileSystemsRead,

		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Capacity", "Performance", "standard", "advance"}, false),
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NFS", "SMB"}, false),
			},
			"description_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRegexp,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"descriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"systems": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metered_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"encrypt_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"file_system_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudFileSystemsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeFileSystems"
	request := make(map[string]interface{})
	request["RegionId"] = client.Region
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	var objects []map[string]interface{}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var filesystemDescriptionRegex *regexp.Regexp
	if v, ok := d.GetOk("description_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		filesystemDescriptionRegex = r
	}
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_file_systems", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)
		resp, err := jsonpath.Get("$.FileSystems.FileSystem", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.FileSystems.FileSystem", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if filesystemDescriptionRegex != nil {
				if !filesystemDescriptionRegex.MatchString(fmt.Sprint(item["Description"])) {
					continue
				}
			}
			if v, ok := d.GetOk("storage_type"); ok && v.(string) != "" && item["StorageType"].(string) != v.(string) {
				continue
			}
			if v, ok := d.GetOk("protocol_type"); ok && v.(string) != "" && item["ProtocolType"].(string) != v.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["FileSystemId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	descriptions := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":               fmt.Sprint(object["FileSystemId"]),
			"region_id":        object["RegionId"],
			"create_time":      object["CreateTime"],
			"description":      object["Description"],
			"protocol_type":    object["ProtocolType"],
			"storage_type":     object["StorageType"],
			"metered_size":     formatInt(object["MeteredSize"]),
			"encrypt_type":     object["EncryptType"],
			"file_system_type": object["FileSystemType"],
			"capacity":         object["Capacity"],
			"kms_key_id":       object["KMSKeyId"],
			"zone_id":          object["ZoneId"],
		}
		ids = append(ids, fmt.Sprint(object["FileSystemId"]))
		descriptions = append(descriptions, object["Description"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("descriptions", descriptions); err != nil {
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
