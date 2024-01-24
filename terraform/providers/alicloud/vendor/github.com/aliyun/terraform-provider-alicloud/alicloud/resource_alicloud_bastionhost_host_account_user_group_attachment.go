package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudBastionhostHostAccountUserGroupAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudBastionhostHostAccountUserGroupAttachmentCreate,
		Read:   resourceAlicloudBastionhostHostAccountUserGroupAttachmentRead,
		Update: resourceAlicloudBastionhostHostAccountUserGroupAttachmentUpdate,
		Delete: resourceAlicloudBastionhostHostAccountUserGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"host_account_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudBastionhostHostAccountUserGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(fmt.Sprint(d.Get("instance_id"), ":", d.Get("user_group_id"), ":", d.Get("host_id")))

	return resourceAlicloudBastionhostHostAccountUserGroupAttachmentUpdate(d, meta)
}
func resourceAlicloudBastionhostHostAccountUserGroupAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	yundunBastionhostService := YundunBastionhostService{client}
	object, err := yundunBastionhostService.DescribeBastionhostHostAccountUserGroupAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_bastionhost_host_account_user_attachment yundunBastionhostService.DescribeBastionhostHostAccountUserGroupAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("host_id", parts[2])
	d.Set("instance_id", parts[0])
	d.Set("user_group_id", parts[1])
	hostAccountIdsItems := make([]string, 0)
	for _, item := range object {
		itemMap := item.(map[string]interface{})
		if v, ok := itemMap["IsAuthorized"]; !ok || !v.(bool) {
			continue
		}
		if v, ok := itemMap["HostAccountId"]; ok && v != nil {
			hostAccountIdsItems = append(hostAccountIdsItems, fmt.Sprint(v))
		}
	}
	d.Set("host_account_ids", hostAccountIdsItems)
	return nil
}
func resourceAlicloudBastionhostHostAccountUserGroupAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	if d.HasChange("host_account_ids") {
		parts, err := ParseResourceId(d.Id(), 3)
		if err != nil {
			return WrapError(err)
		}
		action := "AttachHostAccountsToUserGroup"
		request := make(map[string]interface{})
		conn, err := client.NewBastionhostClient()
		if err != nil {
			return WrapError(err)
		}

		oraw, nraw := d.GetChange("host_account_ids")
		request["InstanceId"] = parts[0]
		request["RegionId"] = client.RegionId
		request["UserGroupId"] = parts[1]

		if oraw != nil && len(oraw.(*schema.Set).List()) > 0 {
			action = "DetachHostAccountsFromUserGroup"
			hostRequestMaps := make([]map[string]interface{}, 0)
			hostRequestMap := make(map[string]interface{}, 0)
			hostRequestMap["HostId"] = parts[2]
			hostRequestMap["HostAccountIds"] = oraw.(*schema.Set).List()
			hostRequestMaps = append(hostRequestMaps, hostRequestMap)
			if v, err := convertListMapToJsonString(hostRequestMaps); err != nil {
				return WrapError(err)
			} else {
				request["Hosts"] = v
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				_, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
		if nraw != nil && len(nraw.(*schema.Set).List()) > 0 {
			action = "AttachHostAccountsToUserGroup"
			hostRequestMaps := make([]map[string]interface{}, 0)
			hostRequestMap := make(map[string]interface{}, 0)
			hostRequestMap["HostId"] = parts[2]
			hostRequestMap["HostAccountIds"] = nraw.(*schema.Set).List()
			hostRequestMaps = append(hostRequestMaps, hostRequestMap)
			if v, err := convertListMapToJsonString(hostRequestMaps); err != nil {
				return WrapError(err)
			} else {
				request["Hosts"] = v
			}
			wait := incrementalWait(3*time.Second, 3*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
				_, err := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
				if err != nil {
					if NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
		}
	}
	return resourceAlicloudBastionhostHostAccountUserGroupAttachmentRead(d, meta)
}
func resourceAlicloudBastionhostHostAccountUserGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DetachHostAccountsFromUserGroup"
	var response map[string]interface{}
	conn, err := client.NewBastionhostClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId":  parts[0],
		"UserGroupId": parts[1],
	}
	request["RegionId"] = client.RegionId
	hostRequestMaps := make([]map[string]interface{}, 0)
	hostRequestMap := make(map[string]interface{}, 0)
	hostRequestMap["HostId"] = parts[2]
	hostRequestMap["HostAccountIds"] = d.Get("host_account_ids").(*schema.Set).List()
	hostRequestMaps = append(hostRequestMaps, hostRequestMap)
	if v, err := convertListMapToJsonString(hostRequestMaps); err != nil {
		return WrapError(err)
	} else {
		request["Hosts"] = v
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-12-09"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
