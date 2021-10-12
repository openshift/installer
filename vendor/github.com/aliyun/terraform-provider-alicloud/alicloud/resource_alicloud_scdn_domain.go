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

func resourceAlicloudScdnDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudScdnDomainCreate,
		Read:   resourceAlicloudScdnDomainRead,
		Update: resourceAlicloudScdnDomainUpdate,
		Delete: resourceAlicloudScdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"biz_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"download", "image", "scdn", "video"}, false),
			},
			"cert_infos": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cert_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"upload", "cas", "free"}, false),
						},
						"ssl_pri": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"ssl_protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
						},
						"ssl_pub": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
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
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
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

func resourceAlicloudScdnDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddScdnDomain"
	request := make(map[string]interface{})
	conn, err := client.NewScdnClient()
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
	sourcesMaps := make([]map[string]interface{}, 0)
	for _, sources := range d.Get("sources").(*schema.Set).List() {
		sourcesArg := sources.(map[string]interface{})
		sourcesMap := map[string]interface{}{}
		sourcesMap["Content"] = sourcesArg["content"]
		sourcesMap["Enabled"] = sourcesArg["enabled"]
		sourcesMap["Port"] = sourcesArg["port"]
		sourcesMap["Priority"] = sourcesArg["priority"]
		sourcesMap["Type"] = sourcesArg["type"]
		sourcesMaps = append(sourcesMaps, sourcesMap)
	}
	if v, err := convertArrayObjectToJsonString(sourcesMaps); err == nil {
		request["Sources"] = v
	} else {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_scdn_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))
	scdnService := ScdnService{client}
	stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudScdnDomainUpdate(d, meta)
}
func resourceAlicloudScdnDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := ScdnService{client}
	object, err := scdnService.DescribeScdnDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_scdn_domain scdnService.DescribeScdnDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("resource_group_id", object["ResourceGroupId"])
	if v, ok := object["Sources"].(map[string]interface{})["Source"].([]interface{}); ok {
		source := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"content":  item["Content"],
				"enabled":  item["Enabled"],
				"port":     item["Port"],
				"priority": item["Priority"],
				"type":     item["Type"],
			}

			source = append(source, temp)
		}
		if err := d.Set("sources", source); err != nil {
			return WrapError(err)
		}
	}
	d.Set("status", object["DomainStatus"])
	describeScdnDomainCertificateInfoObject, err := scdnService.DescribeScdnDomainCertificateInfo(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if certInfosMap, ok := describeScdnDomainCertificateInfoObject["CertInfos"].(map[string]interface{}); ok && certInfosMap != nil {
		if certInfoList, ok := certInfosMap["CertInfo"]; ok && certInfoList != nil {
			certInfosMaps := make([]map[string]interface{}, 0)
			for _, certInfoListItem := range certInfoList.([]interface{}) {
				if certInfoListItemMap, ok := certInfoListItem.(map[string]interface{}); ok {
					certInfoListItemMap["cert_name"] = certInfoListItemMap["CertName"]
					certInfoListItemMap["cert_type"] = certInfoListItemMap["CertType"]
					certInfoListItemMap["ssl_pri"] = certInfoListItemMap["SslPri"]
					certInfoListItemMap["ssl_protocol"] = certInfoListItemMap["SslProtocol"]
					certInfoListItemMap["ssl_pub"] = certInfoListItemMap["SslPub"]
					certInfosMaps = append(certInfosMaps, certInfoListItemMap)
				}
			}
			d.Set("cert_infos", certInfosMaps)
		}
	}

	return nil
}
func resourceAlicloudScdnDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := ScdnService{client}
	var response map[string]interface{}
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if d.HasChange("biz_name") {
		update = true
	}
	if v, ok := d.GetOk("biz_name"); ok {
		request["BizName"] = v
	}
	if update {
		action := "SetScdnDomainBizInfo"
		conn, err := client.NewScdnClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2017-11-15"), StringPointer("AK"), request, nil, &util.RuntimeOptions{})
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
		d.SetPartial("biz_name")
	}
	update = false
	setScdnDomainCertificateReq := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if d.HasChange("ssl_protocol") {
		update = true
	}
	if v, ok := d.GetOk("ssl_protocol"); ok {
		setScdnDomainCertificateReq["SSLProtocol"] = v
	}
	if d.HasChange("cert_name") {
		update = true
		if v, ok := d.GetOk("cert_name"); ok {
			setScdnDomainCertificateReq["CertName"] = v
		}
	}
	if d.HasChange("cert_type") {
		update = true
		if v, ok := d.GetOk("cert_type"); ok {
			setScdnDomainCertificateReq["CertType"] = v
		}
	}
	if d.HasChange("ssl_pri") {
		update = true
		if v, ok := d.GetOk("ssl_pri"); ok {
			setScdnDomainCertificateReq["SSLPri"] = v
		}
	}
	if d.HasChange("ssl_pub") {
		update = true
		if v, ok := d.GetOk("ssl_pub"); ok {
			setScdnDomainCertificateReq["SSLPub"] = v
		}
	}
	if update {
		if v, ok := d.GetOk("force_set"); ok {
			setScdnDomainCertificateReq["ForceSet"] = v
		}
		action := "SetScdnDomainCertificate"
		conn, err := client.NewScdnClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, setScdnDomainCertificateReq, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, setScdnDomainCertificateReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 1*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("ssl_protocol")
		d.SetPartial("cert_name")
		d.SetPartial("cert_type")
		d.SetPartial("ssl_pri")
		d.SetPartial("ssl_pub")
	}
	update = false
	updateScdnDomainReq := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			updateScdnDomainReq["ResourceGroupId"] = v
		}
	}
	if d.HasChange("sources") {
		update = true
		sourcesMaps := make([]map[string]interface{}, 0)
		for _, sources := range d.Get("sources").(*schema.Set).List() {
			sourcesArg := sources.(map[string]interface{})
			sourcesMap := map[string]interface{}{}
			sourcesMap["Content"] = sourcesArg["content"]
			sourcesMap["Enabled"] = sourcesArg["enabled"]
			sourcesMap["Port"] = sourcesArg["port"]
			sourcesMap["Priority"] = sourcesArg["priority"]
			sourcesMap["Type"] = sourcesArg["type"]
			sourcesMaps = append(sourcesMaps, sourcesMap)
		}
		if v, err := convertArrayObjectToJsonString(sourcesMaps); err == nil {
			updateScdnDomainReq["Sources"] = v
		} else {
			return WrapError(err)
		}
	}
	if update {
		action := "UpdateScdnDomain"
		conn, err := client.NewScdnClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, updateScdnDomainReq, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"ServiceBusy"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, updateScdnDomainReq)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{"configure_failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("resource_group_id")
		d.SetPartial("sources")
	}
	if d.HasChange("status") {
		object, err := scdnService.DescribeScdnDomain(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := d.Get("status").(string)
		if object["DomainStatus"].(string) != target {
			if target == "offline" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				action := "StopScdnDomain"
				conn, err := client.NewScdnClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"offline"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			if target == "online" {
				request := map[string]interface{}{
					"DomainName": d.Id(),
				}
				action := "StartScdnDomain"
				conn, err := client.NewScdnClient()
				if err != nil {
					return WrapError(err)
				}
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
				stateConf := BuildStateConf([]string{}, []string{"online"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{}))
				if _, err := stateConf.WaitForState(); err != nil {
					return WrapErrorf(err, IdMsg, d.Id())
				}
			}
			d.SetPartial("status")
		}
	}
	d.Partial(false)
	return resourceAlicloudScdnDomainRead(d, meta)
}
func resourceAlicloudScdnDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	scdnService := ScdnService{client}
	action := "DeleteScdnDomain"
	var response map[string]interface{}
	conn, err := client.NewScdnClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidDomain.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 2*time.Second, scdnService.ScdnDomainStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
