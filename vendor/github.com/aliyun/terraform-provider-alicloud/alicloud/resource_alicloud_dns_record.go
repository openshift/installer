package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsRecordCreate,
		Read:   resourceAlicloudDnsRecordRead,
		Update: resourceAlicloudDnsRecordUpdate,
		Delete: resourceAlicloudDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_record": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRR,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"A", "NS", "MX", "TXT", "CNAME", "SRV", "AAAA", "CAA", "REDIRECT_URL", "FORWORD_URL"}, false),
			},
			"value": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: dnsValueDiffSuppressFunc,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"priority": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: dnsPriorityDiffSuppressFunc,
			},
			"routing": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudDnsRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateAddDomainRecordRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Get("name").(string)
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	request.Value = d.Get("value").(string)
	request.TTL = requests.NewInteger(d.Get("ttl").(int))

	if v, ok := d.GetOk("priority"); !ok && request.Type == "MX" {
		return WrapError(Error("'priority': required field when 'type' is MX."))
	} else if ok {
		request.Priority = requests.Integer(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("routing"); ok {
		routing := v.(string)
		if routing != "default" && request.Type == "FORWORD_URL" {
			return WrapError(Error("The ForwordURLRecord only support default line."))
		}
		request.Line = routing
	}

	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.AddDomainRecord(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InternalError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*alidns.AddDomainRecordResponse)
		d.SetId(response.RecordId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dns_record", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = d.Id()
	request.RR = d.Get("host_record").(string)
	request.Type = d.Get("type").(string)
	if request.Type == "MX" {
		request.Priority = requests.NewInteger(d.Get("priority").(int))
	}
	request.TTL = requests.NewInteger(d.Get("ttl").(int))
	request.Line = d.Get("routing").(string)

	request.Value = d.Get("value").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.UpdateDomainRecord(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAlicloudDnsRecordRead(d, meta)
}

func resourceAlicloudDnsRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsRecord(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("ttl", object.TTL)
	d.Set("priority", object.Priority)
	d.Set("name", object.DomainName)
	d.Set("host_record", object.RR)
	d.Set("type", object.Type)
	d.Set("value", object.Value)
	d.Set("routing", object.Line)
	d.Set("status", object.Status)
	d.Set("locked", object.Locked)

	return nil
}

func resourceAlicloudDnsRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dnsService := &DnsService{client: client}
	request := alidns.CreateDeleteDomainRecordRequest()
	request.RegionId = client.RegionId
	request.RecordId = d.Id()
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainRecord(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"DomainRecordNotBelongToUser"}) {
				return nil
			}
			if IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError"}) {
				return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		_, err = dnsService.DescribeDnsRecord(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
	})
}
