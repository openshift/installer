package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudRdsBackups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRdsBackupsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"backup_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Automated", "Manual"}, false),
			},
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Success", "Failed"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_initiator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_intranet_download_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_size": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consistent_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"copy_only_backup": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"encryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_avail": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"meta_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"backup_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"store_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRdsBackupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeBackups"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	if v, ok := d.GetOk("backup_mode"); ok {
		request["BackupMode"] = v
	}
	request["DBInstanceId"] = d.Get("db_instance_id")
	if v, ok := d.GetOk("end_time"); ok {
		request["EndTime"] = v
	}
	if v, ok := d.GetOk("start_time"); ok {
		request["StartTime"] = v
	}
	if v, ok := d.GetOk("backup_status"); ok {
		request["BackupStatus"] = v
	}
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
	var response map[string]interface{}
	conn, err := client.NewRdsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-08-15"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_rds_backups", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Items.Backup", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Items.Backup", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DBInstanceId"], ":", item["BackupId"])]; !ok {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"backup_download_url":          object["BackupDownloadURL"],
			"backup_end_time":              object["BackupEndTime"],
			"id":                           fmt.Sprint(object["DBInstanceId"], ":", object["BackupId"]),
			"backup_id":                    fmt.Sprint(object["BackupId"]),
			"backup_initiator":             object["BackupInitiator"],
			"backup_intranet_download_url": object["BackupIntranetDownloadURL"],
			"backup_method":                object["BackupMethod"],
			"backup_mode":                  object["BackupMode"],
			"backup_size":                  fmt.Sprint(object["BackupSize"]),
			"backup_start_time":            object["BackupStartTime"],
			"backup_type":                  object["BackupType"],
			"consistent_time":              fmt.Sprint(object["ConsistentTime"]),
			"copy_only_backup":             object["CopyOnlyBackup"],
			"db_instance_id":               object["DBInstanceId"],
			"encryption":                   object["Encryption"],
			"host_instance_id":             object["HostInstanceID"],
			"is_avail":                     formatInt(object["IsAvail"]),
			"meta_status":                  object["MetaStatus"],
			"backup_status":                object["BackupStatus"],
			"storage_class":                object["StorageClass"],
			"store_status":                 object["StoreStatus"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("backups", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
