package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlbAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlbAclCreate,
		Read:   resourceAlicloudAlbAclRead,
		Update: resourceAlicloudAlbAclUpdate,
		Delete: resourceAlicloudAlbAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(16 * time.Minute),
			Delete: schema.DefaultTimeout(16 * time.Minute),
			Update: schema.DefaultTimeout(16 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"acl_entries": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"entry": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"acl_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{2,128}$`), "The name must be `2` to `128` characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter."),
			},
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlbAclCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAcl"
	request := make(map[string]interface{})
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("acl_name"); ok {
		request["AclName"] = v
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	request["ClientToken"] = buildClientToken("CreateAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alb_acl", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AclId"]))
	albService := AlbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, albService.AlbAclStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudAlbAclUpdate(d, meta)
}

func resourceAlicloudAlbAclRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	object, err := albService.DescribeAlbAcl(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alb_acl albService.DescribeAlbAcl Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("acl_name", object["AclName"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("status", object["AclStatus"])
	listAclEntriesObject, err := albService.ListAclEntries(d.Id())
	if err != nil {
		return WrapError(err)
	}
	aclEntriesMaps := make([]map[string]interface{}, 0)
	for _, aclEntriesListItem := range listAclEntriesObject {
		aclEntriesArg := make(map[string]interface{}, 0)
		aclEntriesArg["description"] = aclEntriesListItem["Description"]
		aclEntriesArg["entry"] = aclEntriesListItem["Entry"]
		aclEntriesArg["status"] = aclEntriesListItem["Status"]
		aclEntriesMaps = append(aclEntriesMaps, aclEntriesArg)
	}
	d.Set("acl_entries", aclEntriesMaps)

	listTagResourcesObject, err := albService.ListTagResources(d.Id(), "acl")
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}

func resourceAlicloudAlbAclUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	albService := AlbService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := albService.SetResourceTags(d, "acl"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"ResourceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["NewResourceGroupId"] = v
	}
	request["ResourceType"] = "acl"
	if update {
		action := "MoveResourceGroup"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"NotExist.ResourceGroup"}) || NeedRetry(err) {
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
		d.SetPartial("resource_group_id")
	}
	update = false
	updateAclAttributeReq := map[string]interface{}{
		"AclId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("acl_name") {
		update = true
	}
	if v, ok := d.GetOk("acl_name"); ok {
		updateAclAttributeReq["AclName"] = v
	}
	if update {
		if v, ok := d.GetOkExists("dry_run"); ok {
			updateAclAttributeReq["DryRun"] = v
		}
		action := "UpdateAclAttribute"
		conn, err := client.NewAlbClient()
		if err != nil {
			return WrapError(err)
		}
		request["ClientToken"] = buildClientToken("UpdateAclAttribute")
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, updateAclAttributeReq, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateAclAttributeReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbAclStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("acl_name")
	}
	update = false

	if d.HasChange("acl_entries") {

		oraw, nraw := d.GetChange("acl_entries")
		remove := oraw.(*schema.Set).Difference(nraw.(*schema.Set)).List()
		create := nraw.(*schema.Set).Difference(oraw.(*schema.Set)).List()

		if len(remove) > 0 {
			removeList := SplitSlice(remove, 20)
			for _, item := range removeList {
				removeEntriesFromAclReq := map[string]interface{}{
					"AclId": d.Id(),
				}

				aclEntriesMaps := make([]string, 0)
				for _, aclEntries := range item {
					aclEntriesArg := aclEntries.(map[string]interface{})
					aclEntriesMaps = append(aclEntriesMaps, aclEntriesArg["entry"].(string))
				}
				removeEntriesFromAclReq["Entries"] = aclEntriesMaps

				if v, ok := d.GetOkExists("dry_run"); ok {
					removeEntriesFromAclReq["DryRun"] = v
				}
				action := "RemoveEntriesFromAcl"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("RemoveEntriesFromAcl")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, removeEntriesFromAclReq, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"IncorrectStatus.Acl", "OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling"}) || NeedRetry(err) {
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
				stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbAclStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		if len(create) > 0 {
			createList := SplitSlice(create, 20)
			for _, item := range createList {
				addEntriesToAclReq := map[string]interface{}{
					"AclId": d.Id(),
				}
				aclEntriesMaps := make([]map[string]interface{}, 0)
				for _, aclEntries := range item {
					aclEntriesArg := aclEntries.(map[string]interface{})
					aclEntriesMap := map[string]interface{}{}
					aclEntriesMap["Description"] = aclEntriesArg["description"]
					aclEntriesMap["Entry"] = aclEntriesArg["entry"]
					aclEntriesMaps = append(aclEntriesMaps, aclEntriesMap)
				}
				addEntriesToAclReq["AclEntries"] = aclEntriesMaps

				if v, ok := d.GetOkExists("dry_run"); ok {
					addEntriesToAclReq["DryRun"] = v
				}
				action := "AddEntriesToAcl"
				conn, err := client.NewAlbClient()
				if err != nil {
					return WrapError(err)
				}
				request["ClientToken"] = buildClientToken("AddEntriesToAcl")
				runtime := util.RuntimeOptions{}
				runtime.SetAutoretry(true)
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, addEntriesToAclReq, &runtime)
					if err != nil {
						if IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "Throttling"}) || NeedRetry(err) {
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
				stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, albService.AlbAclStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
		}

		d.SetPartial("acl_entries")
	}
	d.Partial(false)
	return resourceAlicloudAlbAclRead(d, meta)
}

func resourceAlicloudAlbAclDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteAcl"
	var response map[string]interface{}
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AclId": d.Id(),
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["ClientToken"] = buildClientToken("DeleteAcl")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"OperationFailed.ResourceGroupStatusCheckFail", "SystemBusy", "ResourceInUse.Acl", "IncorrectStatus.Acl"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"ResourceNotFound.Acl"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
