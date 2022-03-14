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

func resourceAlicloudDcdnIpaDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDcdnIpaDomainCreate,
		Read:   resourceAlicloudDcdnIpaDomainRead,
		Update: resourceAlicloudDcdnIpaDomainUpdate,
		Delete: resourceAlicloudDcdnIpaDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "global", "overseas"}, false),
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 65535),
						},
						"priority": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"20", "30"}, false),
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipaddr", "domain"}, false),
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
			},
		},
	}
}

func resourceAlicloudDcdnIpaDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	var response map[string]interface{}
	action := "AddDcdnIpaDomain"
	request := make(map[string]interface{})
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	request["DomainName"] = d.Get("domain_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request["Sources"] = sources
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_ipa_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dcdnService.DcdnIpaDomainStateRefreshFunc(d.Id(), []string{"check_failed", "configure_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDcdnIpaDomainRead(d, meta)
}
func resourceAlicloudDcdnIpaDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	object, err := dcdnService.DescribeDcdnIpaDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dcdn_ipa_domain dcdnService.DescribeDcdnIpaDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("scope", object["Scope"])
	if v, ok := object["Sources"]; ok {
		sources := v.(map[string]interface{})
		if sources, ok := sources["Source"]; ok {
			sourceMap := sources.([]interface{})
			sourceMaps := make([]map[string]interface{}, 0)
			for _, val := range sourceMap {
				item := val.(map[string]interface{})
				temp := map[string]interface{}{
					"content":  item["Content"],
					"port":     formatInt(item["Port"]),
					"priority": item["Priority"],
					"type":     item["Type"],
					"weight":   formatInt(item["Weight"]),
				}
				sourceMaps = append(sourceMaps, temp)
			}
			d.Set("sources", sourceMaps)
		}
	}
	d.Set("status", object["DomainStatus"])
	return nil
}
func resourceAlicloudDcdnIpaDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}
	d.Partial(true)
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("sources") {
		update = true
		sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
		request["Sources"] = sources
	}
	if update {
		action := "UpdateDcdnIpaDomain"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, dcdnService.DcdnIpaDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
		d.SetPartial("sources")
	}

	if d.HasChange("status") {
		object, err := dcdnService.DescribeDcdnIpaDomain(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DomainStatus"].(string) != target {
			if target == "offline" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				action := "StopDcdnIpaDomain"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnIpaDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "online" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				action := "StartDcdnIpaDomain"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnIpaDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudDcdnIpaDomainRead(d, meta)
}
func resourceAlicloudDcdnIpaDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	action := "DeleteDcdnIpaDomain"
	var response map[string]interface{}
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, dcdnService.DcdnIpaDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
