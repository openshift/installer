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

func dataSourceAlicloudEcsSnapshots() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsSnapshotsRead,
		Schema: map[string]*schema.Schema{
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"flash", "standard"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"snapshot_link_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"auto", "user", "all"}, false),
				Default:      "all",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Data"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"accomplished", "failed", "progressing", "all"}, false),
			},
			"tags": tagsSchema(),
			"usage": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"image", "disk", "image_disk", "none"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"snapshots": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instant_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"instant_access_retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remain_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_sn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"snapshot_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"usage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEcsSnapshotsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeSnapshots"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("snapshot_link_id"); ok {
		request["SnapshotLinkId"] = v
	}
	if v, ok := d.GetOk("snapshot_name"); ok {
		request["SnapshotName"] = v
	}
	if v, ok := d.GetOk("snapshot_type"); ok {
		request["SnapshotType"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["SnapshotType"] = v
	}
	if v, ok := d.GetOk("source_disk_type"); ok {
		request["SourceDiskType"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}
		request["Tag"] = tags
	}
	if v, ok := d.GetOk("usage"); ok {
		request["Usage"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var snapshotNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		snapshotNameRegex = r
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
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_snapshots", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Snapshots.Snapshot", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Snapshots.Snapshot", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if snapshotNameRegex != nil {
				if !snapshotNameRegex.MatchString(fmt.Sprint(item["SnapshotName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["SnapshotId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
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
			"category":                      object["Category"],
			"description":                   object["Description"],
			"disk_id":                       object["SourceDiskId"],
			"encrypted":                     object["Encrypted"],
			"instant_access":                object["InstantAccess"],
			"instant_access_retention_days": formatInt(object["InstantAccessRetentionDays"]),
			"product_code":                  object["ProductCode"],
			"progress":                      object["Progress"],
			"remain_time":                   formatInt(object["RemainTime"]),
			"resource_group_id":             object["ResourceGroupId"],
			"retention_days":                formatInt(object["RetentionDays"]),
			"id":                            fmt.Sprint(object["SnapshotId"]),
			"snapshot_id":                   fmt.Sprint(object["SnapshotId"]),
			"snapshot_name":                 object["SnapshotName"],
			"snapshot_sn":                   object["SnapshotSN"],
			"snapshot_type":                 object["SnapshotType"],
			"type":                          object["SnapshotType"],
			"source_disk_size":              object["SourceDiskSize"],
			"source_disk_type":              object["SourceDiskType"],
			"source_storage_type":           object["SourceStorageType"],
			"status":                        object["Status"],
			"usage":                         object["Usage"],
			"source_disk_id":                object["SourceDiskId"],
			"creation_time":                 object["CreationTime"],
			"name":                          object["SnapshotName"],
		}

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.Tags.Tag", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(object["SnapshotId"]))
		names = append(names, object["SnapshotName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("snapshots", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
