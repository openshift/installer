package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlidnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsDomainCreate,
		Read:   resourceAlicloudAlidnsDomainRead,
		Update: resourceAlicloudAlidnsDomainUpdate,
		Delete: resourceAlicloudAlidnsDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"puny_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudAlidnsDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainRequest()
	request.DomainName = d.Get("domain_name").(string)
	if v, ok := d.GetOk("group_id"); ok {
		request.GroupId = v.(string)
	}

	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}

	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomain(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_domain", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainResponse)
	d.SetId(fmt.Sprintf("%v", response.DomainName))

	return resourceAlicloudAlidnsDomainUpdate(d, meta)
}
func resourceAlicloudAlidnsDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsDomain(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_domain alidnsService.DescribeAlidnsDomain Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_name", d.Id())
	d.Set("dns_servers", object.DnsServers.DnsServer)
	d.Set("domain_id", object.DomainId)
	d.Set("group_id", object.GroupId)
	d.Set("group_name", object.GroupName)
	d.Set("puny_code", object.PunyCode)
	d.Set("remark", object.Remark)
	d.Set("resource_group_id", object.ResourceGroupId)

	listTagResourcesObject, err := alidnsService.ListTagResources(d.Id())
	if err != nil {
		return WrapError(err)
	}

	tags := make(map[string]string)
	for _, t := range listTagResourcesObject.TagResources {
		tags[t.TagKey] = t.TagValue
	}
	d.Set("tags", tags)
	return nil
}
func resourceAlicloudAlidnsDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := alidnsService.SetResourceTags(d, "DOMAIN"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := alidns.CreateUpdateDomainRemarkRequest()
	request.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if d.HasChange("remark") {
		update = true
		request.Remark = d.Get("remark").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainRemark(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("lang")
		d.SetPartial("remark")
	}
	update = false
	changeDomainGroupReq := alidns.CreateChangeDomainGroupRequest()
	changeDomainGroupReq.DomainName = d.Id()
	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
		changeDomainGroupReq.GroupId = d.Get("group_id").(string)
	}
	if !d.IsNewResource() && d.HasChange("lang") {
		update = true
		changeDomainGroupReq.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.ChangeDomainGroup(changeDomainGroupReq)
		})
		addDebug(changeDomainGroupReq.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), changeDomainGroupReq.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("group_id")
		d.SetPartial("lang")
	}
	d.Partial(false)
	return resourceAlicloudAlidnsDomainRead(d, meta)
}
func resourceAlicloudAlidnsDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainRequest()
	request.DomainName = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.DeleteDomain(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"RecordForbidden.DNSChange", "DnsSystemBusyness", "InternalError"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidDomainName.NoExist"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
