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

func resourceAlicloudDcdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDcdnDomainCreate,
		Read:   resourceAlicloudDcdnDomainRead,
		Update: resourceAlicloudDcdnDomainUpdate,
		Delete: resourceAlicloudDcdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cert_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cert_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"cas", "free", "upload"}, false),
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"force_set": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ssl_pri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ssl_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "off",
			},
			"ssl_pub": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("ssl_protocol").(string) != "on"
				},
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "global", "overseas"}, false),
				Default:      "domestic",
			},
			"security_token": {
				Type:     schema.TypeString,
				Optional: true,
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
							Optional:     true,
							Default:      80,
							ValidateFunc: validation.IntInSlice([]int{443, 80}),
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "20",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
						},
						"weight": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "10",
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"offline", "online"}, false),
				Default:      "online",
			},
			"top_level_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDcdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	var response map[string]interface{}
	action := "AddDcdnDomain"
	request := make(map[string]interface{})
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}

	request["DomainName"] = d.Get("domain_name")
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}

	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}

	if v, ok := d.GetOk("security_token"); ok {
		request["SecurityToken"] = v
	}

	sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
	if err != nil {
		return WrapError(err)
	}
	request["Sources"] = sources
	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}

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
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dcdn_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"check_failed", "configure_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDcdnDomainUpdate(d, meta)
}
func resourceAlicloudDcdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	object, err := dcdnService.DescribeDcdnDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dcdn_domain dcdnService.DescribeDcdnDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("ssl_protocol", convertSSLProtocolResponse(formatSSLProtocolString(object["SSLProtocol"])))
	d.Set("scope", object["Scope"])
	if v, ok := object["Sources"].(map[string]interface{})["Source"].([]interface{}); ok {
		source := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			source = append(source, map[string]interface{}{
				"content":  item["Content"],
				"port":     item["Port"],
				"priority": item["Priority"],
				"type":     item["Type"],
				"weight":   item["Weight"],
			})
		}
		if err := d.Set("sources", source); err != nil {
			return WrapError(err)
		}
	}
	d.Set("status", object["DomainStatus"])

	describeDcdnDomainCertificateInfoObject, err := dcdnService.DescribeDcdnDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("cert_name", describeDcdnDomainCertificateInfoObject["CertName"])
	d.Set("ssl_pub", describeDcdnDomainCertificateInfoObject["SSLPub"])
	return nil
}
func resourceAlicloudDcdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	var response map[string]interface{}
	d.Partial(true)

	if !d.IsNewResource() && d.HasChange("scope") {
		request := map[string]interface{}{
			"DomainName": d.Id(),
		}
		request["Property"] = fmt.Sprintf(`{"coverage":"%s"}`, d.Get("scope").(string))
		action := "ModifyDCdnDomainSchdmByProperty"
		conn, err := client.NewDcdnClient()
		if err != nil {
			return WrapError(err)
		}
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
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("scope")
	}
	update := false
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if d.HasChange("ssl_protocol") {
		update = true
	}
	request["SSLProtocol"] = d.Get("ssl_protocol")
	if d.HasChange("cert_name") {
		update = true
		request["CertName"] = d.Get("cert_name")
	}
	request["Region"] = client.RegionId
	if d.HasChange("ssl_pub") {
		update = true
		request["SSLPub"] = d.Get("ssl_pub")
	}
	if update {
		if _, ok := d.GetOk("cert_type"); ok {
			request["CertType"] = d.Get("cert_type")
		}
		if _, ok := d.GetOk("force_set"); ok {
			request["ForceSet"] = d.Get("force_set")
		}
		if _, ok := d.GetOk("ssl_pri"); ok {
			request["SSLPri"] = d.Get("ssl_pri")
		}
		if _, ok := d.GetOk("security_token"); ok {
			request["SecurityToken"] = d.Get("security_token")
		}
		action := "SetDcdnDomainCertificate"
		conn, err := client.NewDcdnClient()
		if err != nil {
			return WrapError(err)
		}
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
		d.SetPartial("ssl_protocol")
		d.SetPartial("cert_name")
		d.SetPartial("ssl_pub")
	}
	update = false
	updateDcdnDomainReq := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		updateDcdnDomainReq["ResourceGroupId"] = d.Get("resource_group_id")
	}
	if !d.IsNewResource() && d.HasChange("sources") {
		update = true
		sources, err := dcdnService.convertSourcesToString(d.Get("sources").(*schema.Set).List())
		if err != nil {
			return WrapError(err)
		}
		updateDcdnDomainReq["Sources"] = sources
	}
	if update {
		if _, ok := d.GetOk("security_token"); ok {
			updateDcdnDomainReq["SecurityToken"] = d.Get("security_token")
		}
		if _, ok := d.GetOk("top_level_domain"); ok {
			updateDcdnDomainReq["TopLevelDomain"] = d.Get("top_level_domain")
		}
		action := "UpdateDcdnDomain"
		conn, err := client.NewDcdnClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-15"), StringPointer("AK"), nil, updateDcdnDomainReq, &util.RuntimeOptions{})
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
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
		d.SetPartial("sources")
	}
	if d.HasChange("status") {
		object, err := dcdnService.DescribeDcdnDomain(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DomainStatus"].(string) != target {
			if target == "offline" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				if v, ok := d.GetOk("security_token"); ok {
					request["SecurityToken"] = v
				}
				action := "StopDcdnDomain"
				conn, err := client.NewDcdnClient()
				if err != nil {
					return WrapError(err)
				}
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
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "online" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				if v, ok := d.GetOk("security_token"); ok {
					request["SecurityToken"] = v
				}
				action := "StartDcdnDomain"
				conn, err := client.NewDcdnClient()
				if err != nil {
					return WrapError(err)
				}
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
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed", "check_failed"}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudDcdnDomainRead(d, meta)
}
func resourceAlicloudDcdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dcdnService := DcdnService{client}
	action := "DeleteDcdnDomain"
	var response map[string]interface{}
	conn, err := client.NewDcdnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("security_token"); ok {
		request["SecurityToken"] = v
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
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, dcdnService.DcdnDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
func convertSSLProtocolResponse(source string) string {
	switch source {
	case "":
		return "off"
	}
	return source
}
func formatSSLProtocolString(source interface{}) string {
	if source == nil {
		return ""
	} else {
		return source.(string)
	}
}
