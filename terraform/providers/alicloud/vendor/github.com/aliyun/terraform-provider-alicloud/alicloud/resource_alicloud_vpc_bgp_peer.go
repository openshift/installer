package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudVpcBgpPeer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcBgpPeerCreate,
		Read:   resourceAlicloudVpcBgpPeerRead,
		Update: resourceAlicloudVpcBgpPeerUpdate,
		Delete: resourceAlicloudVpcBgpPeerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bfd_multi_hop": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("enable_bfd"); ok && fmt.Sprint(v) == "true" {
						return false
					}
					return true
				},
			},
			"bgp_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"enable_bfd": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"peer_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudVpcBgpPeerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateBgpPeer"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("bfd_multi_hop"); ok {
		request["BfdMultiHop"] = v
	}
	request["BgpGroupId"] = d.Get("bgp_group_id")
	if v, ok := d.GetOkExists("enable_bfd"); ok {
		request["EnableBfd"] = v
	}
	if v, ok := d.GetOk("ip_version"); ok {
		request["IpVersion"] = v
	}
	if v, ok := d.GetOk("peer_ip_address"); ok {
		request["PeerIpAddress"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateBgpPeer")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_bgp_peer", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["BgpPeerId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcBgpPeerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcBgpPeerRead(d, meta)
}
func resourceAlicloudVpcBgpPeerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeVpcBgpPeer(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_bgp_peer vpcService.DescribeVpcBgpPeer Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, ok := object["BfdMultiHop"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bfd_multi_hop", formatInt(v))
	}
	d.Set("bgp_group_id", object["BgpGroupId"])
	d.Set("enable_bfd", object["EnableBfd"])
	d.Set("ip_version", object["IpVersion"])
	d.Set("peer_ip_address", object["PeerIpAddress"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudVpcBgpPeerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"BgpPeerId": d.Id(),
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("bfd_multi_hop") {
		update = true
		if v, ok := d.GetOk("bfd_multi_hop"); ok {
			request["BfdMultiHop"] = v
		}
	}
	if d.HasChange("enable_bfd") {
		update = true
		if v, ok := d.GetOkExists("enable_bfd"); ok {
			request["EnableBfd"] = v
		}
	}
	if d.HasChange("peer_ip_address") {
		update = true
		if v, ok := d.GetOk("peer_ip_address"); ok {
			request["PeerIpAddress"] = v
		}
	}
	if update {
		action := "ModifyBgpPeerAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyBgpPeerAttribute")
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
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, vpcService.VpcBgpPeerStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}
	return resourceAlicloudVpcBgpPeerRead(d, meta)
}
func resourceAlicloudVpcBgpPeerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteBgpPeer"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"BgpPeerId": d.Id(),
	}

	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteBgpPeer")
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcBgpPeerStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
