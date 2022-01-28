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

func resourceAlicloudEcsPrefixList() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcsPrefixListCreate,
		Read:   resourceAlicloudEcsPrefixListRead,
		Update: resourceAlicloudEcsPrefixListUpdate,
		Delete: resourceAlicloudEcsPrefixListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"address_family": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"max_entries": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 200),
			},
			"entry": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 200,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"prefix_list_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudEcsPrefixListCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreatePrefixList"
	request := make(map[string]interface{})
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}

	request["AddressFamily"] = d.Get("address_family")
	request["MaxEntries"] = d.Get("max_entries")
	request["PrefixListName"] = d.Get("prefix_list_name")
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreatePrefixList")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("entry"); ok {
		for entryPtr, entry := range v.(*schema.Set).List() {
			entryArg := entry.(map[string]interface{})
			request["Entry."+fmt.Sprint(entryPtr+1)+".Cidr"] = entryArg["cidr"]
			request["Entry."+fmt.Sprint(entryPtr+1)+".Description"] = entryArg["description"]
		}
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecs_prefix_list", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["PrefixListId"]))

	return resourceAlicloudEcsPrefixListRead(d, meta)
}
func resourceAlicloudEcsPrefixListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsPrefixList(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecs_prefix_list ecsService.DescribeEcsPrefixList Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("address_family", object["AddressFamily"])
	d.Set("description", object["Description"])
	if entryMap, ok := object["Entries"].(map[string]interface{}); ok && entryMap != nil {
		if entryList, ok := entryMap["Entry"]; ok && entryList != nil {
			entryMaps := make([]map[string]interface{}, 0)
			for _, entryListItem := range entryList.([]interface{}) {
				if entryListItemMap, ok := entryListItem.(map[string]interface{}); ok {
					res := make(map[string]interface{}, 0)
					res["cidr"] = entryListItemMap["Cidr"]
					res["description"] = entryListItemMap["Description"]
					entryMaps = append(entryMaps, res)
				}
			}
			d.Set("entry", entryMaps)
		}
	}

	if v, ok := object["MaxEntries"]; ok && fmt.Sprint(v) != "0" {
		d.Set("max_entries", formatInt(v))
	}
	d.Set("prefix_list_name", object["PrefixListName"])
	return nil
}
func resourceAlicloudEcsPrefixListUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false

	request := map[string]interface{}{
		"PrefixListId": d.Id(),
	}
	request["RegionId"] = client.RegionId

	if d.HasChange("prefix_list_name") {
		update = true
		request["PrefixListName"] = d.Get("prefix_list_name")
	}

	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}

	if d.HasChange("entry") {
		update = true
		oraw, nraw := d.GetChange("entry")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		if len(remove) != 0 {
			for entryPtr, entry := range remove {
				entryArg := entry.(map[string]interface{})
				request["RemoveEntry."+fmt.Sprint(entryPtr+1)+".Cidr"] = entryArg["cidr"]
			}
		}
		added := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()
		if len(added) != 0 {
			for entryPtr, entry := range added {
				entryArg := entry.(map[string]interface{})
				request["AddEntry."+fmt.Sprint(entryPtr+1)+".Cidr"] = entryArg["cidr"]
				request["AddEntry."+fmt.Sprint(entryPtr+1)+".Description"] = entryArg["description"]
			}
		}
	}
	if update {
		action := "ModifyPrefixList"
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
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudEcsPrefixListRead(d, meta)
}
func resourceAlicloudEcsPrefixListDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeletePrefixList"
	var response map[string]interface{}
	conn, err := client.NewEcsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"PrefixListId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
