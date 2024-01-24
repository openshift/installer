package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudDfsFileSystem() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDfsFileSystemCreate,
		Read:   resourceAlicloudDfsFileSystemRead,
		Update: resourceAlicloudDfsFileSystemUpdate,
		Delete: resourceAlicloudDfsFileSystemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"HDFS"}, false),
			},
			"provisioned_throughput_in_mi_bps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 1024),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("throughput_mode"); ok && v.(string) == "Provisioned" {
						return false
					}
					return true
				},
			},
			"space_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"storage_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PERFORMANCE", "STANDARD"}, false),
			},
			"throughput_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringInSlice([]string{"Provisioned", "Standard"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateFileSystem"
	request := make(map[string]interface{})
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["FileSystemName"] = d.Get("file_system_name")
	if v, ok := d.GetOk("partition_number"); ok {
		request["PartitionNumber"] = v
	}
	request["ProtocolType"] = d.Get("protocol_type")
	if v, ok := d.GetOk("provisioned_throughput_in_mi_bps"); ok {
		request["ProvisionedThroughputInMiBps"] = v
	}
	request["InputRegionId"] = client.RegionId
	request["SpaceCapacity"] = d.Get("space_capacity")
	request["StorageType"] = d.Get("storage_type")
	if v, ok := d.GetOk("throughput_mode"); ok {
		request["ThroughputMode"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_file_system", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))

	return resourceAlicloudDfsFileSystemRead(d, meta)
}
func resourceAlicloudDfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsService := DfsService{client}
	object, err := dfsService.DescribeDfsFileSystem(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_file_system dfsService.DescribeDfsFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", object["Description"])
	d.Set("file_system_name", object["FileSystemName"])
	d.Set("protocol_type", object["ProtocolType"])
	d.Set("provisioned_throughput_in_mi_bps", formatInt(object["ProvisionedThroughputInMiBps"]))
	d.Set("space_capacity", formatInt(object["SpaceCapacity"]))
	d.Set("storage_type", object["StorageType"])
	d.Set("throughput_mode", object["ThroughputMode"])
	d.Set("zone_id", object["ZoneId"])
	return nil
}
func resourceAlicloudDfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
	}
	request["InputRegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("file_system_name") {
		update = true
		request["FileSystemName"] = d.Get("file_system_name")
	}
	if d.HasChange("provisioned_throughput_in_mi_bps") {
		update = true
		if v, ok := d.GetOk("provisioned_throughput_in_mi_bps"); ok {
			request["ProvisionedThroughputInMiBps"] = v
		}
	}
	if d.HasChange("space_capacity") {
		update = true
		request["SpaceCapacity"] = d.Get("space_capacity")
	}
	if d.HasChange("throughput_mode") {
		update = true
		if v, ok := d.GetOk("throughput_mode"); ok {
			request["ThroughputMode"] = v
		}
	}
	if update {
		action := "ModifyFileSystem"
		conn, err := client.NewAlidfsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudDfsFileSystemRead(d, meta)
}
func resourceAlicloudDfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteFileSystem"
	var response map[string]interface{}
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
	}

	request["InputRegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidParameter.FileSystemNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
