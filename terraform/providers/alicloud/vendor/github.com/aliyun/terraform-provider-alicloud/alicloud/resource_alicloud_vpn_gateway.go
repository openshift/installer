package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnGatewayCreate,
		Read:   resourceAliyunVpnGatewayRead,
		Update: resourceAliyunVpnGatewayUpdate,
		Delete: resourceAliyunVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
				Computed:     true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},

			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateFunc:     validation.Any(validation.IntBetween(1, 9), validation.IntInSlice([]int{12, 24, 36})),
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
			},

			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{5, 10, 20, 50, 100, 200, 500, 1000}),
			},

			"enable_ipsec": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enable_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ssl_connections": {
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				DiffSuppressFunc: vpnSslConnectionsDiffSuppressFunc,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"business_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateVpnGatewayRequest()
	request.RegionId = client.RegionId

	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = d.Get("name").(string)
	}

	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.VSwitchId = d.Get("vswitch_id").(string)
	}

	request.VpcId = d.Get("vpc_id").(string)

	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		if v.(string) == string(PostPaid) {
			request.InstanceChargeType = string("POSTPAY")
		} else {
			request.InstanceChargeType = string("PREPAY")
		}
	}

	if v, ok := d.GetOk("period"); ok && v.(int) != 0 && request.InstanceChargeType == string("PREPAY") {
		request.Period = requests.NewInteger(v.(int))
	}

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))

	if v, ok := d.GetOkExists("enable_ipsec"); ok {
		request.EnableIpsec = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("enable_ssl"); ok {
		request.EnableSsl = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ssl_connections"); ok && v.(int) != 0 {
		request.SslConnections = requests.NewInteger(v.(int))
	}

	request.AutoPay = requests.NewBoolean(true)

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateVpnGateway(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_gateway", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateVpnGatewayResponse)
	d.SetId(response.VpnGatewayId)

	time.Sleep(10 * time.Second)
	if err := vpnGatewayService.WaitForVpnGateway(d.Id(), Active, 2*DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpnGatewayUpdate(d, meta)
}

func resourceAliyunVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeVpnGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("vpc_id", object.VpcId)
	d.Set("internet_ip", object.InternetIp)
	d.Set("status", object.Status)
	d.Set("vswitch_id", object.VSwitchId)
	d.Set("enable_ipsec", "enable" == strings.ToLower(object.IpsecVpn))
	d.Set("enable_ssl", "enable" == strings.ToLower(object.SslVpn))
	d.Set("ssl_connections", object.SslMaxConnections)
	d.Set("business_status", object.BusinessStatus)

	spec := strings.Split(object.Spec, "M")[0]
	bandwidth, err := strconv.Atoi(spec)

	if err == nil {
		d.Set("bandwidth", bandwidth)
	} else {
		return WrapError(err)
	}

	if string("PostpayByFlow") == object.ChargeType {
		d.Set("instance_charge_type", string(PostPaid))
	} else {
		d.Set("instance_charge_type", string(PrePaid))
	}

	return nil
}

func resourceAliyunVpnGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateModifyVpnGatewayAttributeRequest()
	request.RegionId = client.RegionId
	request.VpnGatewayId = d.Id()
	update := false
	d.Partial(true)
	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}

	if update {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyVpnGatewayAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		d.SetPartial("name")
		d.SetPartial("description")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAliyunVpnGatewayRead(d, meta)
	}

	if d.HasChange("bandwidth") {
		return fmt.Errorf("Now Cann't Support modify vpn gateway bandwidth, try to modify on the web console")
	}

	if d.HasChange("enable_ipsec") || d.HasChange("enable_ssl") {
		return fmt.Errorf("Now Cann't Support modify ipsec/ssl switch, try to modify on the web console")
	}

	d.Partial(false)
	return resourceAliyunVpnGatewayRead(d, meta)
}

func resourceAliyunVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateDeleteVpnGatewayRequest()
	request.RegionId = client.RegionId
	request.VpnGatewayId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnGateway(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			/*Vpn known issue: while the vpn is configuring, it will return unknown error*/
			if IsExpectedErrors(err, []string{"UnknownError"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnGatewayInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(vpnGatewayService.WaitForVpnGateway(d.Id(), Deleted, DefaultTimeoutMedium))
}
