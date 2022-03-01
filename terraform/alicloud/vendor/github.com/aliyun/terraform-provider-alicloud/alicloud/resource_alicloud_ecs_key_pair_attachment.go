package alicloud

import (
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcsKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsKeyPairAttachmentCreate,
		Read:   resourceAlicloudEcsKeyPairAttachmentRead,
		Delete: resourceAlicloudEcsKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"key_pair_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"key_name"},
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				Deprecated:    "Field 'key_name' has been deprecated from provider version 1.121.0. New field 'key_pair_name' instead.",
				ConflictsWith: []string{"key_pair_name"},
			},
			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudEcsKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AttachKeyPair"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request["InstanceIds"] = convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())
	if v, ok := d.GetOk("key_pair_name"); ok {
		request["KeyPairName"] = v
	} else if v, ok := d.GetOk("key_name"); ok {
		request["KeyPairName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "key_name" or "key_pair_name" must be set one!`))
	}

	request["RegionId"] = client.RegionId
	ecsService := EcsService{client}
	force := d.Get("force").(bool)
	idsMap := make(map[string]string)
	var newIds []string
	if force {
		ids, err := ecsService.QueryInstancesWithKeyPair("", d.Get("key_pair_name").(string))
		if err != nil {
			return WrapError(err)
		}
		for _, id := range ids {
			idsMap[id] = id
		}
		for _, id := range d.Get("instance_ids").(*schema.Set).List() {
			if _, ok := idsMap[id.(string)]; !ok {
				newIds = append(newIds, id.(string))
			}
		}
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_key_pair_attachment", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["KeyPairName"], ":", request["InstanceIds"]))

	if force {
		requestReboot := make(map[string]interface{})
		requestReboot["RegionId"] = client.RegionId
		requestReboot["ForceStop"] = requests.NewBoolean(true)
		action = "RebootInstance"
		for _, id := range newIds {
			requestReboot["InstanceId"] = id

			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2014-05-26"), StringPointer("AK"), requestReboot, nil, &util.RuntimeOptions{})
			if err != nil {
				return WrapError(err)
			}
			addDebug(action, response, request)
		}
		for _, id := range newIds {
			if err := ecsService.WaitForEcsInstance(id, Running, DefaultLongTimeout); err != nil {
				return WrapError(err)
			}
		}
	}

	return resourceAlicloudEcsKeyPairAttachmentRead(d, meta)
}
func resourceAlicloudEcsKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsKeyPairAttachment(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("key_name", object.KeyPairName)
	d.Set("instance_ids", d.Get("instance_ids"))
	return nil
}
func resourceAlicloudEcsKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DetachKeyPair"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	separatorIndex := strings.LastIndexByte(d.Id(), ':')
	KeyName := d.Id()[:separatorIndex]

	request := map[string]interface{}{
		"KeyPairName": KeyName,
	}

	request["RegionId"] = client.RegionId
	InstanceIds := d.Id()[separatorIndex+1:]
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["InstanceIds"] = InstanceIds
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)

		ecsService := EcsService{client}
		instance_ids, err := ecsService.QueryInstancesWithKeyPair(InstanceIds, KeyName)
		if err != nil {
			return resource.NonRetryableError(WrapError(err))
		}
		if len(instance_ids) > 0 {
			var ids []interface{}
			for _, id := range instance_ids {
				ids = append(ids, id)
			}
			InstanceIds = convertListToJsonString(ids)
			return resource.RetryableError(WrapError(fmt.Errorf("detach Key Pair timeout and the instances including %s has not yet been detached. ", InstanceIds)))
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
