package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsCreate,
		Read:   resourceAlicloudDnsRead,
		Update: resourceAlicloudDnsUpdate,
		Delete: resourceAlicloudDnsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(5, 67),
				ForceNew:     true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_server": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudDnsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Get("name").(string)
	request.ResourceGroupId = d.Get("resource_group_id").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.AddDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dns", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.AddDomainResponse)
	d.SetId(response.DomainName)
	return resourceAlicloudDnsUpdate(d, meta)
}

func resourceAlicloudDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateChangeDomainGroupRequest()
	request.RegionId = client.RegionId
	request.DomainName = d.Id()
	request.GroupId = d.Get("group_id").(string)

	if d.HasChange("group_id") {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.ChangeDomainGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudDnsRead(d, meta)
}

func resourceAlicloudDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDns(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("group_id", object.GroupId)
	d.Set("domain_id", object.DomainId)
	d.Set("name", object.DomainName)
	d.Set("dns_server", object.DnsServers.DnsServer)
	return nil
}

func resourceAlicloudDnsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDeleteDomainRequest()
	request.DomainName = d.Id()
	request.RegionId = client.RegionId
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "InternalError", "DnsSystemBusyness"}) {
				return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
}
