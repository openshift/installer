package alicloud

import (
	"fmt"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbAclCreate,
		Read:   resourceAlicloudSlbAclRead,
		Update: resourceAlicloudSlbAclUpdate,
		Delete: resourceAlicloudSlbAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "ipv4",
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"entry_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				MaxItems: 300,
				MinItems: 0,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudSlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessControlList"
	request := make(map[string]interface{})
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request["RegionId"] = client.RegionId
	if v := d.Get("resource_group_id").(string); v != "" {
		request["ResourceGroupId"] = v
	}
	request["AclName"] = strings.TrimSpace(d.Get("name").(string))
	if v, ok := d.GetOk("address_ip_version"); ok {
		request["AddressIPVersion"] = v
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_slb_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))
	return resourceAlicloudSlbAclUpdate(d, meta)
}

func resourceAlicloudSlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	tags, err := slbService.DescribeTags(d.Id(), nil, TagResourceAcl)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", slbService.tagsToMap(tags))

	object, err := slbService.DescribeSlbAcl(d.Id())
	if err != nil {
		if IsExpectedErrors(err, []string{"AclNotExist"}) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object["AclName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("ip_version", object["AddressIPVersion"])

	if aclEntrys, ok := object["AclEntrys"]; ok {
		if v, ok := aclEntrys.(map[string]interface{})["AclEntry"].([]interface{}); ok {
			aclEntry := make([]map[string]interface{}, 0)
			for _, val := range v {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"comment": item["AclEntryComment"],
					"entry":   item["AclEntryIP"],
				}

				aclEntry = append(aclEntry, temp)
			}
			if err := d.Set("entry_list", aclEntry); err != nil {
				return WrapError(err)
			}
		}
	}
	return nil
}

func resourceAlicloudSlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := slbService.setInstanceTags(d, TagResourceAcl); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if !d.IsNewResource() && d.HasChange("name") {
		request := map[string]interface{}{
			"AclId": d.Id(),
		}
		if v, ok := d.GetOk("name"); ok {
			request["AclName"] = v
		}
		request["RegionId"] = client.RegionId
		action := "SetAccessControlListAttribute"
		conn, err := client.NewSlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		d.SetPartial("name")
	}

	if d.HasChange("entry_list") {
		o, n := d.GetChange("entry_list")
		oe := o.(*schema.Set)
		ne := n.(*schema.Set)
		remove := oe.Difference(ne).List()
		add := ne.Difference(oe).List()

		if len(remove) > 0 {
			removeList := SplitSlice(remove, 50)
			for _, item := range removeList {
				removedRequest := map[string]interface{}{
					"AclId":    d.Id(),
					"RegionId": client.RegionId,
				}
				aclEntries, err := slbService.convertAclEntriesToString(item)
				if err != nil {
					return WrapError(err)
				}
				removedRequest["AclEntrys"] = aclEntries
				action := "RemoveAccessControlListEntry"
				conn, err := client.NewSlbClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, removedRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, removedRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}

		if len(add) > 0 {
			addList := SplitSlice(add, 50)
			for _, item := range addList {
				addedRequest := map[string]interface{}{
					"AclId":    d.Id(),
					"RegionId": client.RegionId,
				}
				aclEntries, err := slbService.convertAclEntriesToString(item)
				if err != nil {
					return WrapError(err)
				}
				addedRequest["AclEntrys"] = aclEntries
				action := "AddAccessControlListEntry"
				conn, err := client.NewSlbClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, addedRequest, &util.RuntimeOptions{})
					if err != nil {
						if NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, addedRequest)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
				}
			}
		}
		d.SetPartial("entry_list")
	}
	d.Partial(false)

	return resourceAlicloudSlbAclRead(d, meta)
}

func resourceAlicloudSlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAccessControlList"
	var response map[string]interface{}
	conn, err := client.NewSlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId":    d.Id(),
		"RegionId": client.RegionId,
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"AclNotExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
