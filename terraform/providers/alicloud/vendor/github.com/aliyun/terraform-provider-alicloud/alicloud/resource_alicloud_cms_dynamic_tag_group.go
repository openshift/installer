package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsDynamicTagGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsDynamicTagGroupCreate,
		Read:   resourceAlicloudCmsDynamicTagGroupRead,
		Delete: resourceAlicloudCmsDynamicTagGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"contact_group_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"tag_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_express": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"tag_value_match_function": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"match_express_filter_relation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"template_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCmsDynamicTagGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateDynamicTagGroup"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["ContactGroupList"] = d.Get("contact_group_list")

	if v, ok := d.GetOk("match_express"); ok {
		for matchExpressPtr, matchExpress := range v.(*schema.Set).List() {
			matchExpressArg := matchExpress.(map[string]interface{})
			request["MatchExpress."+fmt.Sprint(matchExpressPtr+1)+".TagValue"] = matchExpressArg["tag_value"]
			request["MatchExpress."+fmt.Sprint(matchExpressPtr+1)+".TagValueMatchFunction"] = matchExpressArg["tag_value_match_function"]
		}
	}
	if v, ok := d.GetOk("match_express_filter_relation"); ok {
		request["MatchExpressFilterRelation"] = v
	}
	request["TagRegionId"] = client.RegionId
	request["TagKey"] = d.Get("tag_key")
	if v, ok := d.GetOk("template_id_list"); ok {
		request["TemplateIdList"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_dynamic_tag_group", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAlicloudCmsDynamicTagGroupRead(d, meta)
}

func resourceAlicloudCmsDynamicTagGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsDynamicTagGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_dynamic_tag_group cmsService.DescribeCmsDynamicTagGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("status", object["Status"])
	d.Set("match_express_filter_relation", object["MatchExpressFilterRelation"])

	d.Set("template_id_list", object["TemplateIdList"])

	if matchExpressMap, ok := object["MatchExpress"]; ok && matchExpressMap != nil {
		resourceData := make([]map[string]interface{}, 0)
		for _, matchExpressListItem := range matchExpressMap.(map[string]interface{}) {
			for _, val := range matchExpressListItem.([]interface{}) {
				matchExpressObject := make(map[string]interface{}, 0)
				matchExpressObject["tag_value"] = val.(map[string]interface{})["TagValue"]
				matchExpressObject["tag_value_match_function"] = val.(map[string]interface{})["TagValueMatchFunction"]
				resourceData = append(resourceData, matchExpressObject)
			}
		}
		d.Set("match_express", resourceData)
	}

	d.Set("tag_key", object["TagKey"])
	return nil
}
func resourceAlicloudCmsDynamicTagGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteDynamicTagGroup"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"DynamicTagRuleId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	return nil
}
