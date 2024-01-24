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

func resourceAlicloudDdoscooDomainResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooDomainResourceCreate,
		Read:   resourceAlicloudDdoscooDomainResourceRead,
		Update: resourceAlicloudDdoscooDomainResourceUpdate,
		Delete: resourceAlicloudDdoscooDomainResourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"https_ext": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"proxy_types": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_ports": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							ForceNew: true,
						},
						"proxy_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"http", "https", "websocket", "websockets"}, false),
						},
					},
				},
				ForceNew: true,
			},
			"real_servers": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rs_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
		},
	}
}

func resourceAlicloudDdoscooDomainResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDomainResource"
	request := make(map[string]interface{})
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	request["Domain"] = d.Get("domain")
	if v, ok := d.GetOk("https_ext"); ok {
		request["HttpsExt"] = v
	}

	request["InstanceIds"] = d.Get("instance_ids")
	proxyTypesMaps := make([]map[string]interface{}, 0)
	for _, proxyTypes := range d.Get("proxy_types").(*schema.Set).List() {
		proxyTypesMap := make(map[string]interface{})
		proxyTypesArg := proxyTypes.(map[string]interface{})
		proxyTypesMap["ProxyPorts"] = proxyTypesArg["proxy_ports"]
		proxyTypesMap["ProxyType"] = proxyTypesArg["proxy_type"]
		proxyTypesMaps = append(proxyTypesMaps, proxyTypesMap)
	}
	request["ProxyTypes"] = proxyTypesMaps

	request["RealServers"] = d.Get("real_servers")
	request["RsType"] = d.Get("rs_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_domain_resource", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Domain"]))

	return resourceAlicloudDdoscooDomainResourceRead(d, meta)
}
func resourceAlicloudDdoscooDomainResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}
	object, err := ddoscooService.DescribeDdoscooDomainResource(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ddoscoo_domain_resource ddoscooService.DescribeDdoscooDomainResource Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain", d.Id())
	d.Set("https_ext", object["HttpsExt"])
	d.Set("instance_ids", object["InstanceIds"])

	proxyTypes := make([]map[string]interface{}, 0)
	if proxyTypesList, ok := object["ProxyTypes"].([]interface{}); ok {
		for _, v := range proxyTypesList {
			if m1, ok := v.(map[string]interface{}); ok {
				temp1 := map[string]interface{}{
					"proxy_ports": m1["ProxyPorts"],
					"proxy_type":  m1["ProxyType"],
				}
				proxyTypes = append(proxyTypes, temp1)

			}
		}
	}
	if err := d.Set("proxy_types", proxyTypes); err != nil {
		return WrapError(err)
	}
	d.Set("real_servers", object["RealServers"])
	d.Set("rs_type", formatInt(object["RsType"]))
	return nil
}
func resourceAlicloudDdoscooDomainResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Domain": d.Id(),
	}
	ProxyTypes := make([]map[string]interface{}, len(d.Get("proxy_types").(*schema.Set).List()))
	for i, ProxyTypesValue := range d.Get("proxy_types").(*schema.Set).List() {
		ProxyTypesMap := ProxyTypesValue.(map[string]interface{})
		ProxyTypes[i] = map[string]interface{}{
			"ProxyPorts": ProxyTypesMap["proxy_ports"],
			"ProxyType":  ProxyTypesMap["proxy_type"],
		}
	}
	request["ProxyTypes"] = ProxyTypes

	if d.HasChange("real_servers") {
		update = true
	}
	request["RealServers"] = d.Get("real_servers")
	request["RsType"] = d.Get("rs_type")
	if d.HasChange("https_ext") {
		update = true
		request["HttpsExt"] = d.Get("https_ext")
	}
	if d.HasChange("instance_ids") {
		update = true
	}
	request["InstanceIds"] = d.Get("instance_ids")
	if update {
		action := "ModifyDomainResource"
		conn, err := client.NewDdoscooClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	}
	return resourceAlicloudDdoscooDomainResourceRead(d, meta)
}
func resourceAlicloudDdoscooDomainResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDomainResource"
	var response map[string]interface{}
	conn, err := client.NewDdoscooClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Domain": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}
