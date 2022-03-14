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

func dataSourceAlicloudEcsDisks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEcsDisksRead,
		Schema: map[string]*schema.Schema{
			"additional_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"auto_snapshot_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_essd", "cloud_ssd", "ephemeral_ssd"}, false),
			},
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
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
			"disk_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "data", "all"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"system", "data", "all"}, false),
				Deprecated:   "Field 'type' has been deprecated from provider version 1.122.0. New field 'disk_type' instead.",
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"enable_automated_snapshot_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"operation_locks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lock_reason": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"portable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Attaching", "Available", "Creating", "Detaching", "In_use", "Migrating", "ReIniting", "Transferring"}, false),
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'availability_zone' has been deprecated from provider version 1.122.0. New field 'zone_id' instead",
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_snapshot_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_auto_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detached_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enable_auto_snapshot": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_automated_snapshot_policy": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"encrypted": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expired_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iops": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"iops_read": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"iops_write": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mount_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attached_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"operation_locks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"lock_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"performance_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"portable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"snapshot_id": {
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
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_type": {
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
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudEcsDisksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDisks"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("additional_attributes"); ok {
		request["AdditionalAttributes"] = convertListToJsonString(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("auto_snapshot_policy_id"); ok {
		request["AutoSnapshotPolicyId"] = v
	}
	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}
	if v, ok := d.GetOkExists("delete_auto_snapshot"); ok {
		request["DeleteAutoSnapshot"] = v
	}
	if v, ok := d.GetOkExists("delete_with_instance"); ok {
		request["DeleteWithInstance"] = v
	}
	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	}
	if v, ok := d.GetOk("disk_type"); ok {
		request["DiskType"] = v
	} else if v, ok := d.GetOk("type"); ok {
		request["DiskType"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOkExists("enable_auto_snapshot"); ok {
		request["EnableAutoSnapshot"] = v
	}
	if v, ok := d.GetOkExists("enable_automated_snapshot_policy"); ok {
		request["EnableAutomatedSnapshotPolicy"] = v
	}
	if v, ok := d.GetOkExists("enable_shared"); ok {
		request["EnableShared"] = v
	}
	if v, ok := d.GetOkExists("encrypted"); ok {
		if v == string(OnFlag) {
			request["Encrypted"] = true
		} else {
			request["Encrypted"] = false
		}
	}
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("operation_locks"); ok {
		operationLocksMaps := make([]map[string]interface{}, 0)
		for _, operationLocks := range v.([]interface{}) {
			operationLocksMap := make(map[string]interface{})
			operationLocksArg := operationLocks.(map[string]interface{})
			operationLocksMap["LockReason"] = operationLocksArg["lock_reason"]
			operationLocksMaps = append(operationLocksMaps, operationLocksMap)
		}
		request["OperationLocks"] = operationLocksMaps

	}
	if v, ok := d.GetOk("payment_type"); ok {
		request["DiskChargeType"] = convertEcsDiskPaymentTypeRequest(v.(string))
	}
	if v, ok := d.GetOkExists("portable"); ok {
		request["Portable"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
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
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	if v, ok := d.GetOk("page_number"); ok && v.(int) > 0 {
		request["PageNumber"] = v.(int)
	} else {
		request["PageNumber"] = 1
	}
	if v, ok := d.GetOk("page_size"); ok && v.(int) > 0 {
		request["PageSize"] = v.(int)
	} else {
		request["PageSize"] = PageSizeLarge
	}
	var objects []interface{}
	var diskNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		diskNameRegex = r
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
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ecs_disks", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Disks.Disk", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Disks.Disk", response)
		}
		result, _ := resp.([]interface{})
		if isPagingRequest(d) {
			objects = result
			break
		}
		for _, v := range result {
			item := v.(map[string]interface{})
			if diskNameRegex != nil {
				if !diskNameRegex.MatchString(fmt.Sprint(item["DiskName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DiskId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, v := range objects {
		object := v.(map[string]interface{})
		mapping := map[string]interface{}{
			"attached_time":                    object["AttachedTime"],
			"auto_snapshot_policy_id":          object["AutoSnapshotPolicyId"],
			"category":                         object["Category"],
			"delete_auto_snapshot":             object["DeleteAutoSnapshot"],
			"delete_with_instance":             object["DeleteWithInstance"],
			"description":                      object["Description"],
			"detached_time":                    object["DetachedTime"],
			"device":                           object["Device"],
			"id":                               fmt.Sprint(object["DiskId"]),
			"disk_id":                          fmt.Sprint(object["DiskId"]),
			"disk_name":                        object["DiskName"],
			"enable_auto_snapshot":             object["EnableAutoSnapshot"],
			"enable_automated_snapshot_policy": object["EnableAutomatedSnapshotPolicy"],
			"encrypted":                        string(OnFlag),
			"expired_time":                     object["ExpiredTime"],
			"iops":                             formatInt(object["IOPS"]),
			"iops_read":                        formatInt(object["IOPSRead"]),
			"iops_write":                       formatInt(object["IOPSWrite"]),
			"image_id":                         object["ImageId"],
			"instance_id":                      object["InstanceId"],
			"kms_key_id":                       object["KMSKeyId"],
			"mount_instance_num":               formatInt(object["MountInstanceNum"]),
			"payment_type":                     convertEcsDiskPaymentTypeResponse(object["DiskChargeType"].(string)),
			"performance_level":                object["PerformanceLevel"],
			"portable":                         object["Portable"],
			"product_code":                     object["ProductCode"],
			"resource_group_id":                object["ResourceGroupId"],
			"size":                             formatInt(object["Size"]),
			"snapshot_id":                      object["SourceSnapshotId"],
			"status":                           object["Status"],
			"zone_id":                          object["ZoneId"],
			"availability_zone":                object["ZoneId"],
			"disk_type":                        object["Type"],
			"type":                             object["Type"],
			"region_id":                        object["RegionId"],
			"creation_time":                    object["CreationTime"],
			"name":                             object["DiskName"],
		}
		if !object["Encrypted"].(bool) {
			mapping["encrypted"] = string(OffFlag)
		}

		mountInstance := make([]map[string]interface{}, 0)
		if mountInstances, ok := object["MountInstances"]; ok {
			if mountInstanceList, ok := mountInstances.(map[string]interface{})["MountInstance"].([]interface{}); ok {
				for _, v := range mountInstanceList {
					if m1, ok := v.(map[string]interface{}); ok {
						temp1 := map[string]interface{}{
							"attached_time": m1["AttachedTime"],
							"device":        m1["Device"],
							"instance_id":   m1["InstanceId"],
						}
						mountInstance = append(mountInstance, temp1)
					}
				}
			}
		}
		mapping["mount_instances"] = mountInstance

		operationLock := make([]map[string]interface{}, 0)
		if operationLockList, ok := object["OperationLocks"].(map[string]interface{})["OperationLock"].([]interface{}); ok {
			for _, v := range operationLockList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"lock_reason": m1["LockReason"],
					}
					operationLock = append(operationLock, temp1)
				}
			}
		}
		mapping["operation_locks"] = operationLock

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
		ids = append(ids, fmt.Sprint(object["DiskId"]))
		names = append(names, object["DiskName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("disks", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("total_count", formatInt(response["TotalCount"])); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
