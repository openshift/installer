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

func resourceAlicloudDbfsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbfsInstanceCreate,
		Read:   resourceAlicloudDbfsInstanceRead,
		Update: resourceAlicloudDbfsInstanceUpdate,
		Delete: resourceAlicloudDbfsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(15 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"attach_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attach_point": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard"}, false),
			},
			"delete_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ecs_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ecs_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enable_raid": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"performance_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PL0", "PL1", "PL2", "PL3"}, false),
			},
			"raid_stripe_unit_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
			"used_scene": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDbfsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDbfs"
	request := make(map[string]interface{})
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	request["Category"] = d.Get("category")
	if v, ok := d.GetOkExists("delete_snapshot"); ok {
		request["DeleteSnapshot"] = v
	}
	if v, ok := d.GetOkExists("enable_raid"); ok {
		request["EnableRaid"] = v
	}
	if v, ok := d.GetOkExists("encryption"); ok {
		request["Encryption"] = v
	}
	request["FsName"] = d.Get("instance_name")
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KMSKeyId"] = v
	}
	if v, ok := d.GetOk("performance_level"); ok {
		request["PerformanceLevel"] = v
	}
	if v, ok := d.GetOk("raid_stripe_unit_number"); ok {
		request["RaidStripeUnitNumber"] = v
	}
	request["RegionId"] = client.RegionId
	request["SizeG"] = d.Get("size")
	if v, ok := d.GetOk("snapshot_id"); ok {
		request["SnapshotId"] = v
	}
	if v, ok := d.GetOk("used_scene"); ok {
		request["UsedScene"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	request["ClientToken"] = buildClientToken("CreateDbfs")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dbfs_instance", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["FsId"]))
	dbfsService := DbfsService{client}
	stateConf := BuildStateConf([]string{}, []string{"unattached"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDbfsInstanceUpdate(d, meta)
}
func resourceAlicloudDbfsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	object, err := dbfsService.DescribeDbfsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dbfs_instance dbfsService.DescribeDbfsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("category", object["Category"])
	if ecsListList, ok := object["EcsList"]; ok && ecsListList != nil {
		ecsListMaps := make([]map[string]interface{}, 0)
		for _, ecsListListItem := range ecsListList.([]interface{}) {
			if ecsListListItemMap, ok := ecsListListItem.(map[string]interface{}); ok {
				ecsListListItemMap["ecs_id"] = ecsListListItemMap["EcsId"]
				ecsListMaps = append(ecsListMaps, ecsListListItemMap)
			}
			d.Set("ecs_list", ecsListMaps)
		}
	}

	d.Set("enable_raid", object["EnableRaid"])
	d.Set("encryption", object["Encryption"])
	d.Set("kms_key_id", object["KMSKeyId"])
	d.Set("performance_level", object["PerformanceLevel"])
	if v, ok := object["SizeG"]; ok && fmt.Sprint(v) != "0" {
		d.Set("size", formatInt(v))
	}
	d.Set("status", object["Status"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("zone_id", object["ZoneId"])
	d.Set("instance_name", object["FsName"])
	return nil
}
func resourceAlicloudDbfsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"FsId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
	}
	request["FsName"] = d.Get("instance_name")
	if update {
		action := "RenameDbfs"
		conn, err := client.NewDbfsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("instance_name")
	}
	update = false
	resizeDbfsReq := map[string]interface{}{
		"FsId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("size") {
		update = true
	}
	resizeDbfsReq["NewSizeG"] = d.Get("size")
	if update {
		action := "ResizeDbfs"
		conn, err := client.NewDbfsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, resizeDbfsReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, resizeDbfsReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		dbfsService := DbfsService{client}
		stateConf := BuildStateConf([]string{}, []string{"attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("size")
	}

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		remove := oraw.(map[string]interface{})
		create := nraw.(map[string]interface{})

		if len(remove) > 0 {

			deleteTagsBatchReq := map[string]interface{}{
				"DbfsList": "[\"" + d.Id() + "\"]",
				"RegionId": client.RegionId,
			}

			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range remove {
				tagsMap := map[string]interface{}{}
				tagsMap["TagKey"] = key
				tagsMap["TagValue"] = value
				tagsMaps = append(tagsMaps, tagsMap)
			}
			deleteTagsBatchReq["Tags"], _ = convertListMapToJsonString(tagsMaps)

			action := "DeleteTagsBatch"
			conn, err := client.NewDbfsClient()
			if err != nil {
				return WrapError(err)
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, deleteTagsBatchReq, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, deleteTagsBatchReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		if len(create) > 0 {

			addTagsBatchReq := map[string]interface{}{
				"DbfsList": "[\"" + d.Id() + "\"]",
				"RegionId": client.RegionId,
			}

			tagsMaps := make([]map[string]interface{}, 0)
			for key, value := range create {
				tagsMap := map[string]interface{}{}
				tagsMap["TagKey"] = key
				tagsMap["TagValue"] = value
				tagsMaps = append(tagsMaps, tagsMap)
			}
			addTagsBatchReq["Tags"], _ = convertListMapToJsonString(tagsMaps)

			action := "AddTagsBatch"
			conn, err := client.NewDbfsClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("AddTagsBatch")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, addTagsBatchReq, &runtime)
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addTagsBatchReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}

		d.SetPartial("tags")
	}

	if d.HasChange("ecs_list") {
		oldEcsList, newEcsList := d.GetChange("ecs_list")
		oldEcsListSet := oldEcsList.(*schema.Set)
		newEcsListSet := newEcsList.(*schema.Set)
		removed := oldEcsListSet.Difference(newEcsListSet)
		added := newEcsListSet.Difference(oldEcsListSet)

		if removed.Len() > 0 {
			detachdbfsrequest := map[string]interface{}{
				"FsId": d.Id(),
			}
			detachdbfsrequest["RegionId"] = client.RegionId
			detachdbfsrequest["ECSInstanceId"] = d.Get("ecs_instance_id")
			for _, ecsArg := range removed.List() {
				ecsMap := ecsArg.(map[string]interface{})
				detachdbfsrequest["ECSInstanceId"] = ecsMap["ecs_id"]

				action := "DetachDbfs"
				conn, err := client.NewDbfsClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, detachdbfsrequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, detachdbfsrequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}

				dbfsService := DbfsService{client}
				stateConf := BuildStateConf([]string{}, []string{"unattached", "attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		if added.Len() > 0 {
			attachdbfsrequest := map[string]interface{}{
				"FsId": d.Id(),
			}
			attachdbfsrequest["RegionId"] = client.RegionId

			action := "AttachDbfs"
			for _, ecsArg := range added.List() {
				ecsMap := ecsArg.(map[string]interface{})
				attachdbfsrequest["ECSInstanceId"] = ecsMap["ecs_id"]

				conn, err := client.NewDbfsClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, attachdbfsrequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, attachdbfsrequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
				dbfsService := DbfsService{client}
				stateConf := BuildStateConf([]string{}, []string{"attached"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		d.SetPartial("ecs_list")
	}
	d.Partial(false)
	return resourceAlicloudDbfsInstanceRead(d, meta)
}
func resourceAlicloudDbfsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbfsService := DbfsService{client}
	action := "DeleteDbfs"
	var response map[string]interface{}
	conn, err := client.NewDbfsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FsId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"EntityNotExist.DBFS"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, dbfsService.DbfsInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertDbfsInstancePaymentTypeResponse(source string) string {
	switch source {
	case "postpaid":
		return "PayAsYouGo"
	}
	return source
}
