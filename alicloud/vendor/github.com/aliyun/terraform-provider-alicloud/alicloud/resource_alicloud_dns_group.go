package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDnsGroupCreate,
		Read:   resourceAlicloudDnsGroupRead,
		Update: resourceAlicloudDnsGroupUpdate,
		Delete: resourceAlicloudDnsGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceAlicloudDnsGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateAddDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupName = d.Get("name").(string)

	raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
		return dnsClient.AddDomainGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dns_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*alidns.AddDomainGroupResponse)
	d.SetId(response.GroupId)
	return resourceAlicloudDnsGroupRead(d, meta)
}

func resourceAlicloudDnsGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateUpdateDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	if d.HasChange("name") {
		request.GroupName = d.Get("name").(string)
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.UpdateDomainGroup(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudDnsGroupRead(d, meta)
}

func resourceAlicloudDnsGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dnsService := &DnsService{client: client}
	object, err := dnsService.DescribeDnsGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.GroupName)
	return nil
}

func resourceAlicloudDnsGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateDeleteDomainGroupRequest()
	request.RegionId = client.RegionId
	request.GroupId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DeleteDomainGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Fobidden.NotEmptyGroup"}) {
				return resource.RetryableError(WrapErrorf(err, DefaultTimeoutMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
}
