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

func resourceAlicloudVpcNatIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudVpcNatIpCreate,
		Read:   resourceAlicloudVpcNatIpRead,
		Update: resourceAlicloudVpcNatIpUpdate,
		Delete: resourceAlicloudVpcNatIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"dry_run": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nat_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nat_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nat_ip_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"nat_ip_cidr_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_ip_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"nat_ip_name": {
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

func resourceAlicloudVpcNatIpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateNatIp"
	request := make(map[string]interface{})
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["NatGatewayId"] = d.Get("nat_gateway_id")
	if v, ok := d.GetOk("nat_ip"); ok {
		request["NatIp"] = v
	}
	if v, ok := d.GetOk("nat_ip_cidr"); ok {
		request["NatIpCidr"] = v
	}
	if v, ok := d.GetOk("nat_ip_cidr_id"); ok {
		request["NatIpCidrId"] = v
	}
	if v, ok := d.GetOk("nat_ip_description"); ok {
		request["NatIpDescription"] = v
	}
	if v, ok := d.GetOk("nat_ip_name"); ok {
		request["NatIpName"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("CreateNatIp")
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpc_nat_ip", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["NatGatewayId"], ":", response["NatIpId"]))
	vpcService := VpcService{client}
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, vpcService.VpcNatIpStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudVpcNatIpRead(d, meta)
}
func resourceAlicloudVpcNatIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeVpcNatIp(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_vpc_nat_ip vpcService.DescribeVpcNatIp Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("nat_gateway_id", object["NatGatewayId"])
	d.Set("nat_ip", object["NatIp"])
	d.Set("nat_ip_id", object["NatIpId"])
	d.Set("nat_ip_cidr", object["NatIpCidr"])
	d.Set("nat_ip_description", object["NatIpDescription"])
	d.Set("nat_ip_name", object["NatIpName"])
	d.Set("status", object["NatIpStatus"])
	return nil
}
func resourceAlicloudVpcNatIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"RegionId": client.RegionId,
		"NatIpId":  parts[1],
	}
	if d.HasChange("nat_ip_description") {
		update = true
	}
	if v, ok := d.GetOk("nat_ip_description"); ok {
		request["NatIpDescription"] = v
	}
	if d.HasChange("nat_ip_name") {
		update = true
	}
	if v, ok := d.GetOk("nat_ip_name"); ok {
		request["NatIpName"] = v
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("dry_run") || d.IsNewResource() {
		update = true
		if v, ok := d.GetOkExists("dry_run"); ok {
			request["DryRun"] = v
		}
	}
	if update {
		action := "ModifyNatIpAttribute"
		conn, err := client.NewVpcClient()
		if err != nil {
			return WrapError(err)
		}
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			request["ClientToken"] = buildClientToken("ModifyNatIpAttribute")
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
	}
	return resourceAlicloudVpcNatIpRead(d, meta)
}
func resourceAlicloudVpcNatIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	action := "DeleteNatIp"
	var response map[string]interface{}
	conn, err := client.NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"NatIpId": parts[1],
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request["DryRun"] = v
	}
	request["RegionId"] = client.RegionId
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		request["ClientToken"] = buildClientToken("DeleteNatIp")
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
	stateConf := BuildStateConf([]string{}, []string{}, d.Timeout(schema.TimeoutDelete), 5*time.Second, vpcService.VpcNatIpStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
