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

func resourceAlicloudEcsDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsDiskCreate,
		Read:   resourceAlicloudEcsDiskRead,
		Update: resourceAlicloudEcsDiskUpdate,
		Delete: resourceAlicloudEcsDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"advanced_features": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud", "cloud_efficiency", "cloud_essd", "cloud_ssd"}, false),
				Default:      "cloud_efficiency",
			},
			"dedicated_block_storage_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delete_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_with_instance": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'name' has been deprecated from provider version 1.122.0. New field 'disk_name' instead.",
				ConflictsWith: []string{"disk_name"},
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_auto_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"encrypt_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"encrypted": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id"},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"payment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PayAsYouGo", "Subscription"}, false),
			},
			"performance_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("category").(string) != "cloud_essd"
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"encrypted"},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_set_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_set_partition_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
			},
			"zone_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"availability_zone"},
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Deprecated:    "Field 'availability_zone' has been deprecated from provider version 1.122.0. New field 'zone_id' instead",
				ConflictsWith: []string{"zone_id"},
			},
		},
	}
}

func resourceAlicloudEcsDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "CreateDisk"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("advanced_features"); ok {
		request["AdvancedFeatures"] = v
	}

	if v, ok := d.GetOk("category"); ok {
		request["DiskCategory"] = v
	}

	if v, ok := d.GetOk("dedicated_block_storage_cluster_id"); ok {
		request["DedicatedBlockStorageClusterId"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("disk_name"); ok {
		request["DiskName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["DiskName"] = v
	}

	if v, ok := d.GetOk("encrypt_algorithm"); ok {
		request["EncryptAlgorithm"] = v
	}

	if v, ok := d.GetOkExists("encrypted"); ok {
		request["Encrypted"] = v
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}

	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}

	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("size"); ok {
		request["Size"] = v
	}

	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}

	if v, ok := d.GetOk("storage_set_id"); ok {
		request["StorageSetId"] = v
	}

	if v, ok := d.GetOk("storage_set_partition_number"); ok {
		request["StorageSetPartitionNumber"] = v
	}

	if v, ok := d.GetOk("tags"); ok {
		count := 1
		for key, value := range v.(map[string]interface{}) {
			request[fmt.Sprintf("Tag.%d.Key", count)] = key
			request[fmt.Sprintf("Tag.%d.Value", count)] = value
			count++
		}
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	} else if v, ok := d.GetOk("availability_zone"); ok {
		request["ZoneId"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_disk", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["DiskId"]))
	stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, ecsService.EcsDiskStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcsDiskUpdate(d, meta)
}
func resourceAlicloudEcsDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDisk(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_disk ecsService.DescribeEcsDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("category", object["Category"])
	d.Set("delete_auto_snapshot", object["DeleteAutoSnapshot"])
	d.Set("delete_with_instance", object["DeleteWithInstance"])
	d.Set("description", object["Description"])
	d.Set("disk_name", object["DiskName"])
	d.Set("name", object["DiskName"])
	d.Set("enable_auto_snapshot", object["EnableAutoSnapshot"])
	d.Set("encrypted", object["Encrypted"])
	d.Set("instance_id", object["InstanceId"])
	d.Set("kms_key_id", object["KMSKeyId"])
	d.Set("payment_type", convertEcsDiskPaymentTypeResponse(object["DiskChargeType"]))
	d.Set("performance_level", object["PerformanceLevel"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("size", formatInt(object["Size"]))
	d.Set("snapshot_id", object["SourceSnapshotId"])
	d.Set("status", object["Status"])
	if v, ok := object["Tags"].(map[string]interface{}); ok {
		d.Set("tags", tagsToMap(v["Tag"]))
	}
	d.Set("zone_id", object["ZoneId"])
	d.Set("availability_zone", object["ZoneId"])
	return nil
}
func resourceAlicloudEcsDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "disk"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	if !d.IsNewResource() && d.HasChange("size") {
		request := map[string]interface{}{
			"DiskId": d.Id(),
		}
		request["NewSize"] = d.Get("size")
		if _, ok := d.GetOk("type"); ok {
			request["Type"] = d.Get("type")
		}
		action := "ResizeDisk"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("ResizeDisk")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("size")
		d.SetPartial("type")
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		request := map[string]interface{}{
			"ResourceId": d.Id(),
		}
		request["RegionId"] = client.RegionId
		request["ResourceGroupId"] = d.Get("resource_group_id")
		request["ResourceType"] = "disk"
		action := "JoinResourceGroup"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("resource_group_id")
	}
	update := false
	request := map[string]interface{}{
		"DiskId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("category") {
		update = true
		request["DiskCategory"] = d.Get("category")
	}
	if !d.IsNewResource() && d.HasChange("performance_level") {
		update = true
		request["PerformanceLevel"] = d.Get("performance_level")
	}
	if update {
		if _, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = d.Get("dry_run")
		}
		action := "ModifyDiskSpec"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available", "In_use"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDiskStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("category")
		d.SetPartial("dry_run")
		d.SetPartial("performance_level")
	}
	update = false
	modifyDiskChargeTypeReq := map[string]interface{}{
		"DiskIds": convertListToJsonString([]interface{}{d.Id()}),
	}
	if !d.IsNewResource() && d.HasChange("instance_id") {
		update = true
	}
	modifyDiskChargeTypeReq["ClientToken"] = buildClientToken("ModifyDiskChargeType")
	modifyDiskChargeTypeReq["InstanceId"] = d.Get("instance_id")
	modifyDiskChargeTypeReq["RegionId"] = client.RegionId
	modifyDiskChargeTypeReq["AutoPay"] = true
	if d.HasChange("payment_type") && (d.Get("payment_type").(string) == "Subscription" || !d.IsNewResource()) {
		update = true
		modifyDiskChargeTypeReq["DiskChargeType"] = convertEcsDiskPaymentTypeRequest(d.Get("payment_type").(string))
	}
	if update {
		action := "ModifyDiskChargeType"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, modifyDiskChargeTypeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyDiskChargeTypeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("instance_id")
		d.SetPartial("payment_type")
	}
	update = false
	modifyDiskAttributeReq := map[string]interface{}{
		"DiskId": d.Id(),
	}
	if d.HasChange("delete_auto_snapshot") || d.IsNewResource() {
		update = true
		modifyDiskAttributeReq["DeleteAutoSnapshot"] = d.Get("delete_auto_snapshot")
	}
	if d.HasChange("delete_with_instance") || d.IsNewResource() {
		update = true
		modifyDiskAttributeReq["DeleteWithInstance"] = d.Get("delete_with_instance")
	}
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
		modifyDiskAttributeReq["Description"] = d.Get("description")
	}
	if !d.IsNewResource() && d.HasChange("disk_name") {
		update = true
		modifyDiskAttributeReq["DiskName"] = d.Get("disk_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		modifyDiskAttributeReq["DiskName"] = d.Get("name")
	}
	if d.HasChange("enable_auto_snapshot") || d.IsNewResource() {
		update = true
		modifyDiskAttributeReq["EnableAutoSnapshot"] = d.Get("enable_auto_snapshot")
	}
	if update {
		action := "ModifyDiskAttribute"
		conn, err := client.NewEcsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, modifyDiskAttributeReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, modifyDiskAttributeReq)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("delete_auto_snapshot")
		d.SetPartial("delete_with_instance")
		d.SetPartial("description")
		d.SetPartial("name")
		d.SetPartial("disk_name")
		d.SetPartial("enable_auto_snapshot")
	}
	d.Partial(false)
	return resourceAlicloudEcsDiskRead(d, meta)
}
func resourceAlicloudEcsDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDisk"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DiskId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectDiskStatus.Initializing"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDiskId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertEcsDiskPaymentTypeResponse(source interface{}) interface{} {
	switch source {
	case "PostPaid":
		return "PayAsYouGo"
	case "PrePaid":
		return "Subscription"
	}
	return source
}
func convertEcsDiskPaymentTypeRequest(source interface{}) interface{} {
	switch source {
	case "PayAsYouGo":
		return "PostPaid"
	case "Subscription":
		return "PrePaid"
	}
	return source
}
