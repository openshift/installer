package alicloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudAlidnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsRecordCreate,
		Read:   resourceAlicloudAlidnsRecordRead,
		Update: resourceAlicloudAlidnsRecordUpdate,
		Delete: resourceAlicloudAlidnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"line": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"rr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ENABLE", "DISABLE"}, false),
				Default:      "ENABLE",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dnsValueDiffSuppressFunc,
			},
		},
	}
}

func resourceAlicloudAlidnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("line"); ok {
		line := v.(string)
		if line != "default" && d.Get("type").(string) == "FORWORD_URL" {
			return WrapError(Error("The ForwordURLRecord only support default line."))
		}
		request.Line = line
	}
	if v, ok := d.GetOk("priority"); !ok && d.Get("type").(string) == "MX" {
		return WrapError(Error("'priority': required field when 'type' is MX."))
	} else if ok {
		request.Priority = requests.Integer(strconv.Itoa(v.(int)))
	}
	request.RR = d.Get("rr").(string)
	if v, ok := d.GetOk("ttl"); ok {
		request.TTL = requests.NewInteger(v.(int))
	}
	request.Type = d.Get("type").(string)
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	request.Value = d.Get("value").(string)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.AddDomainRecord(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "LastOperationNotFinished"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*alidns.AddDomainRecordResponse)
		d.SetId(fmt.Sprintf("%v", response.RecordId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudAlidnsRecordUpdate(d, meta)
}
func resourceAlicloudAlidnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_record alidnsService.DescribeAlidnsRecord Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", object.DomainName)
	d.Set("line", object.Line)
	d.Set("priority", object.Priority)
	d.Set("rr", object.RR)
	d.Set("status", object.Status)
	d.Set("ttl", object.TTL)
	d.Set("type", object.Type)
	d.Set("value", object.Value)
	return nil
}
func resourceAlicloudAlidnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)

	update := false
	request := alidns.CreateSetDomainRecordStatusRequest()
	request.RecordId = d.Id()
	if d.HasChange("status") {
		update = true
	}
	request.Status = d.Get("status").(string)
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
	}
	request.Lang = d.Get("lang").(string)
	if !d.IsNewResource() && d.HasChange("user_client_ip") {
		update = true
	}
	request.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.SetDomainRecordStatus(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("status")
		d.SetPartial("lang")
		d.SetPartial("user_client_ip")
	}
	update = false
	updateDomainRecordRemarkReq := alidns.CreateUpdateDomainRecordRemarkRequest()
	updateDomainRecordRemarkReq.RecordId = d.Id()
	updateDomainRecordRemarkReq.Lang = d.Get("lang").(string)
	if d.HasChange("remark") {
		update = true
	}
	updateDomainRecordRemarkReq.Remark = d.Get("remark").(string)
	updateDomainRecordRemarkReq.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRecordRemark(updateDomainRecordRemarkReq)
		})
		addDebug(updateDomainRecordRemarkReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateDomainRecordRemarkReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("remark")
	}
	update = false
	updateDomainRecordReq := alidns.CreateUpdateDomainRecordRequest()
	updateDomainRecordReq.RecordId = d.Id()
	if !d.IsNewResource() && d.HasChange("rr") {
		update = true
	}
	updateDomainRecordReq.RR = d.Get("rr").(string)
	if !d.IsNewResource() && d.HasChange("type") {
		update = true
	}
	updateDomainRecordReq.Type = d.Get("type").(string)
	if !d.IsNewResource() && d.HasChange("value") {
		update = true
	}
	updateDomainRecordReq.Value = d.Get("value").(string)
	updateDomainRecordReq.Lang = d.Get("lang").(string)
	if !d.IsNewResource() && d.HasChange("line") {
		update = true
	}
	updateDomainRecordReq.Line = d.Get("line").(string)
	if updateDomainRecordReq.Type == "MX" {
		if !d.IsNewResource() && d.HasChange("priority") {
			update = true
		}
		updateDomainRecordReq.Priority = requests.NewInteger(d.Get("priority").(int))
	}

	if !d.IsNewResource() && d.HasChange("ttl") {
		update = true
	}
	updateDomainRecordReq.TTL = requests.NewInteger(d.Get("ttl").(int))
	updateDomainRecordReq.UserClientIp = d.Get("user_client_ip").(string)
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRecord(updateDomainRecordReq)
		})
		addDebug(updateDomainRecordReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), updateDomainRecordReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("rr")
		d.SetPartial("type")
		d.SetPartial("value")
		d.SetPartial("line")
		d.SetPartial("ttl")
	}
	d.Partial(false)
	return resourceAlicloudAlidnsRecordRead(d, meta)
}
func resourceAlicloudAlidnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainRecordRequest()
	request.RecordId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request.UserClientIp = v.(string)
	}
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DeleteDomainRecord(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError", "RecordForbidden.DNSChange"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
