package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudVpcBgpGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcBgpGroupCreate,
		Read:   resourceAlicloudVpcBgpGroupRead,
		Update: resourceAlicloudVpcBgpGroupUpdate,
		Delete: resourceAlicloudVpcBgpGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"auth_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bgp_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][A-Za-z0-9._-]{1,127}$`), "The name must be `2` to `128` characters in length and can contain digits, periods (.), underscores (_), and hyphens (-)."), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(2, 256), validation.StringDoesNotMatch(regexp.MustCompile(`(^http://.*)|(^https://.*)`), "It cannot begin with \"http://\", \"https://\".")),
			},
			"is_fake_asn": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"local_asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"peer_asn": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcBgpGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBgpGroup"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("auth_key"); ok {
		request["AuthKey"] = v
	}
	if v, ok := d.GetOk("bgp_group_name"); ok {
		request["Name"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOkExists("is_fake_asn"); ok {
		request["IsFakeAsn"] = v
	}
	if v, ok := d.GetOk("local_asn"); ok {
		request["LocalAsn"] = v
	}
	request["PeerAsn"] = d.Get("peer_asn")
	request["RegionId"] = client.RegionId
	request["RouterId"] = d.Get("router_id")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateBgpGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_bgp_group", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BgpGroupId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcBgpGroupRead(d, meta)
}
func resourceAlicloudVpcBgpGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcBgpGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_bgp_group vpcService.DescribeVpcBgpGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("auth_key", object["AuthKey"])
	d.Set("bgp_group_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("local_asn", formatInt(object["LocalAsn"]))
	d.Set("peer_asn", formatInt(object["PeerAsn"]))
	d.Set("router_id", object["RouterId"])
	d.Set("is_fake_asn", object["IsFake"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudVpcBgpGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"BgpGroupId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("auth_key") {
		update = true
		if v, ok := d.GetOk("auth_key"); ok {
			request["AuthKey"] = v
		}
	}
	if d.HasChange("bgp_group_name") {
		update = true
		if v, ok := d.GetOk("bgp_group_name"); ok {
			request["Name"] = v
		}
	}
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("local_asn") {
		update = true
		if v, ok := d.GetOk("local_asn"); ok {
			request["LocalAsn"] = v
		}
	}
	if d.HasChange("peer_asn") {
		update = true
		request["PeerAsn"] = d.Get("peer_asn")
	}
	if update {
		if v, ok := d.GetOkExists("is_fake_asn"); ok {
			request["IsFakeAsn"] = v
		}
		action := "ModifyBgpGroupAttribute"
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyBgpGroupAttribute")
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudVpcBgpGroupRead(d, meta)
}
func resourceAlicloudVpcBgpGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteBgpGroup"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BgpGroupId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteBgpGroup")
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcBgpGroupStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
