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

func resourceAlicloudEcdNetworkPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdNetworkPackageCreate,
		Read:   resourceAlicloudEcdNetworkPackageRead,
		Update: resourceAlicloudEcdNetworkPackageUpdate,
		Delete: resourceAlicloudEcdNetworkPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"office_site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"internet_charge_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcdNetworkPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNetworkPackage"
	request := make(map[string]interface{})
	ecdService := EcdService{client}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}

	request["Bandwidth"] = d.Get("bandwidth")

	request["OfficeSiteId"] = d.Get("office_site_id")

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_network_package", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["NetworkPackageId"]))

	stateConf := BuildStateConf([]string{"Creating"}, []string{"InUse"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcdNetworkPackageRead(d, meta)
}
func resourceAlicloudEcdNetworkPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdNetworkPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_network_package ecdService.DescribeEcdNetworkPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["Bandwidth"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bandwidth", formatInt(v))
	}
	d.Set("internet_charge_type", object["InternetChargeType"])
	d.Set("office_site_id", object["OfficeSiteId"])
	d.Set("status", object["NetworkPackageStatus"])
	return nil
}
func resourceAlicloudEcdNetworkPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	request := map[string]interface{}{
		"RegionId":         client.RegionId,
		"NetworkPackageId": d.Id(),
	}
	if d.HasChange("bandwidth") {
		request["Bandwidth"] = d.Get("bandwidth")
	}
	action := "ModifyNetworkPackage"
	ecdService := EcdService{client}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IncorrectNetworkPackageStatus.ModificationNotSupport"}) {
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
	stateConf := BuildStateConf([]string{"Creating"}, []string{"InUse"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudEcdNetworkPackageRead(d, meta)
}
func resourceAlicloudEcdNetworkPackageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteNetworkPackages"
	var response map[string]interface{}
	ecdService := EcdService{client}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"NetworkPackageId": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"IncorrectNetworkPackageStatus.DeletionNotSupport"}) {
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
	stateConf := BuildStateConf([]string{"Releasing"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, ecdService.EcdNetworkPackageRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
