package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudGaAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudGaAclCreate,
		Read:   resourceAlicloudGaAclRead,
		Update: resourceAlicloudGaAclUpdate,
		Delete: resourceAlicloudGaAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entry": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"entry_description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9._/-]{1,256}$`), "The description of the IP entry. The description must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_)."),
						},
					},
				},
			},
			"acl_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{2,128}$`), "The name must be `2` to `128` characters in length, and can contain letters, digits, periods (.), hyphens (-) and underscores (_). It must start with a letter."),
			},
			"address_ip_version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudGaAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAcl"
	request := make(map[string]interface{})
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	if m, ok := d.GetOk("acl_entries"); ok {
		for k, aclEntries := range m.(*schema.Set).List() {
			aclEntriesArg := aclEntries.(map[string]interface{})
			request[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
			request[fmt.Sprintf("AclEntries.%d.EntryDescription", k+1)] = aclEntriesArg["entry_description"].(string)
		}
	}
	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}
	request["AddressIPVersion"] = d.Get("address_ip_version")
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("CreateAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ga_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))
	gaService := GaService{client}
	stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudGaAclRead(d, meta)
}
func resourceAlicloudGaAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	object, err := gaService.DescribeGaAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ga_acl gaService.DescribeGaAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["AclEntries"].([]interface{}); ok {
		aclEntries := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"entry":             item["Entry"],
				"entry_description": item["EntryDescription"],
			}
			aclEntries = append(aclEntries, temp)
		}
		if err := d.Set("acl_entries", aclEntries); err != nil {
			return WrapError(err)
		}
	}

	d.Set("acl_name", object["AclName"])
	d.Set("address_ip_version", object["AddressIPVersion"])
	d.Set("status", object["AclStatus"])
	return nil
}
func resourceAlicloudGaAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"AclId": d.Id(),
	}
	if d.HasChange("acl_name") {
		update = true
	}
	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}
	request["RegionId"] = client.RegionId
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
		action := "UpdateAclAttribute"
		conn, err := client.NewGaplusClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAclAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 20*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
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
		d.SetPartial("acl_name")
	}

	if d.HasChange("acl_entries") {
		oraw, nraw := d.GetChange("acl_entries")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeEntriesFromAclReq := map[string]interface{}{
				"AclId": d.Id(),
			}
			for k, aclEntries := range remove {
				aclEntriesArg := aclEntries.(map[string]interface{})
				removeEntriesFromAclReq[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
			}
			removeEntriesFromAclReq["RegionId"] = client.RegionId

			if v, ok := d.GetOkExists("dry_run"); ok {
				removeEntriesFromAclReq["DryRun"] = v
			}
			action := "RemoveEntriesFromAcl"
			conn, err := client.NewGaplusClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("RemoveEntriesFromAcl")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 20*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, removeEntriesFromAclReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, removeEntriesFromAclReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		if len(create) > 0 {
			addEntriesToAclReq := map[string]interface{}{
				"AclId": d.Id(),
			}
			for k, aclEntries := range create {
				aclEntriesArg := aclEntries.(map[string]interface{})
				addEntriesToAclReq[fmt.Sprintf("AclEntries.%d.Entry", k+1)] = aclEntriesArg["entry"].(string)
				addEntriesToAclReq[fmt.Sprintf("AclEntries.%d.EntryDescription", k+1)] = aclEntriesArg["entry_description"].(string)
			}
			addEntriesToAclReq["RegionId"] = client.RegionId
			if v, ok := d.GetOkExists("dry_run"); ok {
				addEntriesToAclReq["DryRun"] = v
			}
			action := "AddEntriesToAcl"
			conn, err := client.NewGaplusClient()
			if err != nil {
				return WrapError(err)
			}
			request["ClientToken"] = buildClientToken("AddEntriesToAcl")
			runtime := util.RuntimeOptions{}
			runtime.SetAutoretry(true)
			wait := incrementalWait(3*time.Second, 20*time.Second)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, addEntriesToAclReq, &runtime)
				if err != nil {
					if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
						wait()
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			addDebug(action, response, addEntriesToAclReq)
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{}, []string{"active"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
			if _, err := stateConf.WaitForState(); err != nil {
				return WrapErrorf(err, IdMsg, d.Id())
			}
		}

		d.SetPartial("acl_entries")
	}

	d.Partial(false)
	return resourceAlicloudGaAclRead(d, meta)
}
func resourceAlicloudGaAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	gaService := GaService{client}
	action := "DeleteAcl"
	var response map[string]interface{}
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId": d.Id(),
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	request["ClientToken"] = buildClientToken("DeleteAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 20*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"StateError.Acl"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"NotExist.Acl"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, gaService.GaAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
