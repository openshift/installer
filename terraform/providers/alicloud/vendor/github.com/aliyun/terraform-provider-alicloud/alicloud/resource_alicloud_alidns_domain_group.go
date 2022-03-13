package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAlidnsDomainGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAlidnsDomainGroupCreate,
		Read:   resourceAlicloudAlidnsDomainGroupRead,
		Update: resourceAlicloudAlidnsDomainGroupUpdate,
		Delete: resourceAlicloudAlidnsDomainGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"group_name"},
			},
			"group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'group_name' has been deprecated from version 1.97.0. Use 'domain_group_name' instead.",
				ConflictsWith: []string{"domain_group_name"},
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudAlidnsDomainGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := alidns.CreateAddDomainGroupRequest()
	if v, ok := d.GetOk("domain_group_name"); ok {
		request.GroupName = v.(string)
	} else if v, ok := d.GetOk("group_name"); ok {
		request.GroupName = v.(string)
	} else {
		return WrapError(Error(`[ERROR] Argument "group_name" or "domain_group_name" must be set one!`))
	}

	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}

	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.AddDomainGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_alidns_domain_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*alidns.AddDomainGroupResponse)
	d.SetId(fmt.Sprintf("%v", response.GroupId))

	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alidnsService := AlidnsService{client}
	object, err := alidnsService.DescribeAlidnsDomainGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alidns_domain_group alidnsService.DescribeAlidnsDomainGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("domain_group_name", object.GroupName)
	d.Set("group_name", object.GroupName)
	return nil
}
func resourceAlicloudAlidnsDomainGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := alidns.CreateUpdateDomainGroupRequest()
	request.GroupId = d.Id()
	if d.HasChange("domain_group_name") {
		update = true
		request.GroupName = d.Get("domain_group_name").(string)
	}
	if d.HasChange("group_name") {
		update = true
		request.GroupName = d.Get("group_name").(string)
	}
	if request.GroupName == "" {
		request.GroupName = d.Get("domain_group_name").(string)
	}
	if d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
			return alidnsClient.UpdateDomainGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudAlidnsDomainGroupRead(d, meta)
}
func resourceAlicloudAlidnsDomainGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := alidns.CreateDeleteDomainGroupRequest()
	request.GroupId = d.Id()
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}
	raw, err := client.WithAlidnsClient(func(alidnsClient *alidns.Client) (interface{}, error) {
		return alidnsClient.DeleteDomainGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
