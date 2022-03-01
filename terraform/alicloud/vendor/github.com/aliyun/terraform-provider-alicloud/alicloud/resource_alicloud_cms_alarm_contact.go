package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCmsAlarmContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsAlarmContactCreate,
		Read:   resourceAlicloudCmsAlarmContactRead,
		Update: resourceAlicloudCmsAlarmContactUpdate,
		Delete: resourceAlicloudCmsAlarmContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alarm_contact_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"channels_aliim": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_ding_web_hook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_mail": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"channels_sms": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"describe": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lang": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"en", "zh-cn"}, false),
			},
		},
	}
}

func resourceAlicloudCmsAlarmContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cms.CreatePutContactRequest()
	request.ContactName = d.Get("alarm_contact_name").(string)
	if v, ok := d.GetOk("channels_aliim"); ok {
		request.ChannelsAliIM = v.(string)
	}

	if v, ok := d.GetOk("channels_ding_web_hook"); ok {
		request.ChannelsDingWebHook = v.(string)
	}

	if v, ok := d.GetOk("channels_mail"); ok {
		request.ChannelsMail = v.(string)
	}

	if v, ok := d.GetOk("channels_sms"); ok {
		request.ChannelsSMS = v.(string)
	}

	request.Describe = d.Get("describe").(string)
	if v, ok := d.GetOk("lang"); ok {
		request.Lang = v.(string)
	}

	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.PutContact(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alarm_contact", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.PutContactResponse)

	if response.Code != "200" {
		return WrapError(Error("PutContact failed for " + response.Message))
	}
	d.SetId(fmt.Sprintf("%v", request.ContactName))

	return resourceAlicloudCmsAlarmContactRead(d, meta)
}
func resourceAlicloudCmsAlarmContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsAlarmContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_alarm_contact cmsService.DescribeCmsAlarmContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("alarm_contact_name", d.Id())
	d.Set("channels_aliim", object.Channels.AliIM)
	d.Set("channels_ding_web_hook", object.Channels.DingWebHook)
	d.Set("channels_mail", object.Channels.Mail)
	d.Set("channels_sms", object.Channels.SMS)
	d.Set("describe", object.Desc)
	d.Set("lang", object.Lang)
	return nil
}
func resourceAlicloudCmsAlarmContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	update := false
	request := cms.CreatePutContactRequest()
	request.ContactName = d.Id()
	if d.HasChange("channels_aliim") {
		update = true
	}
	request.ChannelsAliIM = d.Get("channels_aliim").(string)
	if d.HasChange("channels_ding_web_hook") {
		update = true
	}
	request.ChannelsDingWebHook = d.Get("channels_ding_web_hook").(string)
	if d.HasChange("channels_mail") {
		update = true
	}
	request.ChannelsMail = d.Get("channels_mail").(string)
	if d.HasChange("channels_sms") {
		update = true
	}
	request.ChannelsSMS = d.Get("channels_sms").(string)
	if d.HasChange("describe") {
		update = true
	}
	request.Describe = d.Get("describe").(string)
	if d.HasChange("lang") {
		update = true
		request.Lang = d.Get("lang").(string)
	}
	if update {
		raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
			return cmsClient.PutContact(request)
		})
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cms.PutContactResponse)

		if response.Code != "200" {
			return WrapError(Error("PutContact failed for " + response.Message))
		}
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudCmsAlarmContactRead(d, meta)
}
func resourceAlicloudCmsAlarmContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cms.CreateDeleteContactRequest()
	request.ContactName = d.Id()
	raw, err := client.WithCmsClient(func(cmsClient *cms.Client) (interface{}, error) {
		return cmsClient.DeleteContact(request)
	})
	addDebug(request.GetActionName(), raw)
	response, _ := raw.(*cms.DeleteContactResponse)

	if response.Code != "200" {
		return WrapError(Error("DeleteContact failed for " + response.Message))
	}
	if err != nil {
		if IsExpectedErrors(err, []string{"400", "403", "404", "ContactNotExists"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
