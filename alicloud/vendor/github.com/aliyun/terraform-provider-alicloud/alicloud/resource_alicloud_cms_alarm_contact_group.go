package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsAlarmContactGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmContactGroupCreate,
		Read:   resourceAlicloudCmsAlarmContactGroupRead,
		Update: resourceAlicloudCmsAlarmContactGroupUpdate,
		Delete: resourceAlicloudCmsAlarmContactGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alarm_contact_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contacts": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"describe": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_subscribed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsAlarmContactGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreatePutContactGroupRequest()
	request.ContactGroupName = d.Get("alarm_contact_group_name").(string)
	if v, ok := d.GetOk("contacts"); ok {
		contactNames := expandStringList(v.(*schema.Set).List())
		request.ContactNames = &contactNames
	}

	if v, ok := d.GetOk("describe"); ok {
		request.Describe = v.(string)
	}

	if v, ok := d.GetOkExists("enable_subscribed"); ok {
		request.EnableSubscribed = requests.NewBoolean(v.(bool))
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutContactGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alarm_contact_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.PutContactGroupResponse)

	if response.Code != "200" {
		return WrapError(Error("PutContactGroup failed for " + response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactGroupName))

	return resourceAlicloudCmsAlarmContactGroupRead(d, meta)
}
func resourceAlicloudCmsAlarmContactGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContactGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_alarm_contact_group cmsService.DescribeCmsAlarmContactGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alarm_contact_group_name", d.Id())
	d.Set("contacts", object.Contacts.Contact)
	d.Set("describe", object.Describe)
	d.Set("enable_subscribed", object.EnableSubscribed)
	return nil
}
func resourceAlicloudCmsAlarmContactGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := cms.CreatePutContactGroupRequest()
	request.ContactGroupName = d.Id()
	if d.HasChange("contacts") {
		update = true
		contactNames := expandStringList(d.Get("contacts").(*schema.Set).List())
		request.ContactNames = &contactNames

	}
	if d.HasChange("describe") {
		update = true
		request.Describe = d.Get("describe").(string)
	}
	if d.HasChange("enable_subscribed") {
		update = true
		request.EnableSubscribed = requests.NewBoolean(d.Get("enable_subscribed").(bool))
	}
	if update {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.PutContactGroup(request)
		})
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cms.PutContactGroupResponse)

		if response.Code != "200" {
			return WrapError(Error("PutContactGroup failed for " + response.Message))
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCmsAlarmContactGroupRead(d, meta)
}
func resourceAlicloudCmsAlarmContactGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cms.CreateDeleteContactGroupRequest()
	request.ContactGroupName = d.Id()
	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContactGroup(request)
	})
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.DeleteContactGroupResponse)

	if response.Code != "200" {
		return WrapError(Error("DeleteContactGroup failed for " + response.Message))
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
