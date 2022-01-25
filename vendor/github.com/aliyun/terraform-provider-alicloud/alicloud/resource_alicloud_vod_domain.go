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

func resourceAlicloudVodDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVodDomainCreate,
		Read:   resourceAlicloudVodDomainRead,
		Update: resourceAlicloudVodDomainUpdate,
		Delete: resourceAlicloudVodDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_content": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"source_priority": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice([]string{"20", "30"}, false),
						},
						"source_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ipaddr", "domain", "oss"}, false),
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gmt_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weight": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_pub": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"domestic", "overseas", "global"}, false),
				Default:      "domestic",
			},
			"top_level_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudVodDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "AddVodDomain"
	request := make(map[string]interface{})
	conn, err := client.NewVodClient()
	if err != nil {
		return WrapError(err)
	}
	request["DomainName"] = d.Get("domain_name")
	if v, ok := d.GetOk("check_url"); ok {
		request["CheckUrl"] = v
	}

	if v, ok := d.GetOk("scope"); ok {
		request["Scope"] = v
	}
	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}
	sourcesMaps := make([]map[string]interface{}, 0)
	for _, sources := range d.Get("sources").(*schema.Set).List() {
		sourcesArg := sources.(map[string]interface{})
		sourcesMap := map[string]interface{}{}
		sourcesMap["content"] = sourcesArg["source_content"]
		sourcesMap["priority"] = sourcesArg["source_priority"]
		sourcesMap["type"] = sourcesArg["source_type"]
		sourcesMap["port"] = sourcesArg["source_port"]
		sourcesMaps = append(sourcesMaps, sourcesMap)
	}
	if v, err := convertArrayObjectToJsonString(sourcesMaps); err == nil {
		request["Sources"] = v
	} else {
		return WrapError(err)
	}
	if v, ok := d.GetOk("top_level_domain"); ok {
		request["TopLevelDomain"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-03-21"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vod_domain", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["DomainName"]))

	return resourceAlicloudVodDomainUpdate(d, meta)
}
func resourceAlicloudVodDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
	object, err := vodService.DescribeVodDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vod_domain vodService.DescribeVodDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("scope", object["Scope"])
	d.Set("status", object["DomainStatus"])
	d.Set("gmt_modified", object["GmtModified"])
	d.Set("gmt_created", object["GmtCreated"])
	d.Set("weight", object["Weight"])
	d.Set("ssl_pub", object["SSLPub"])
	d.Set("ssl_protocol", object["SSLProtocol"])
	d.Set("cert_name", object["CertName"])
	d.Set("cname", object["Cname"])
	if v, ok := object["Sources"].(map[string]interface{})["Source"].([]interface{}); ok {
		source := make([]map[string]interface{}, 0)
		for _, val := range v {
			item := val.(map[string]interface{})
			temp := map[string]interface{}{
				"source_content":  item["Content"],
				"source_port":     item["Port"],
				"source_priority": item["Priority"],
				"source_type":     item["Type"],
			}

			source = append(source, temp)
		}
		if err := d.Set("sources", source); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
func resourceAlicloudVodDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}
	if d.HasChange("sources") {
		update = true
		sourcesMaps := make([]map[string]interface{}, 0)
		for _, sources := range d.Get("sources").(*schema.Set).List() {
			sourcesArg := sources.(map[string]interface{})
			sourcesMap := map[string]interface{}{}
			sourcesMap["content"] = sourcesArg["source_content"]
			sourcesMap["priority"] = sourcesArg["source_priority"]
			sourcesMap["type"] = sourcesArg["source_type"]
			sourcesMap["port"] = sourcesArg["source_port"]
			sourcesMaps = append(sourcesMaps, sourcesMap)
		}
		if v, err := convertArrayObjectToJsonString(sourcesMaps); err == nil {
			request["Sources"] = v
		} else {
			return WrapError(err)
		}
	}
	if d.HasChange("top_level_domain") {
		update = true
		if v, ok := d.GetOk("top_level_domain"); ok {
			request["TopLevelDomain"] = v
		}
	}
	if update {
		action := "UpdateVodDomain"
		conn, err := client.NewVodClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-03-21"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"3", "2001"}) || NeedRetry(err) {
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

	if err := vodService.SetResourceTags(d, "DOMAIN"); err != nil {
		return WrapError(err)
	}
	stateConf := BuildStateConf([]string{"configuring"}, []string{"online"}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVodDomainRead(d, meta)
}
func resourceAlicloudVodDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vodService := VodService{client}
	action := "DeleteVodDomain"
	var response map[string]interface{}
	conn, err := client.NewVodClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DomainName": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-03-21"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"2001"}) || NeedRetry(err) {
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

	stateConf := BuildStateConf([]string{"deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vodService.VodStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return nil
}
