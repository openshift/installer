package alicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnConnectionCreate,
		Read:   resourceAliyunVpnConnectionRead,
		Update: resourceAliyunVpnConnectionUpdate,
		Delete: resourceAliyunVpnConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"customer_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpn_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},

			"local_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"remote_subnet": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateCIDRNetworkAddress,
				},
				MinItems: 1,
				MaxItems: 10,
			},

			"effect_immediately": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ike_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"psk": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_VERSION_1,
							ValidateFunc: validation.StringInSlice([]string{IKE_VERSION_1, IKE_VERSION_2}, false),
						},
						"ike_mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      IKE_MODE_MAIN,
							ValidateFunc: validation.StringInSlice([]string{IKE_MODE_MAIN, IKE_MODE_AGGRESSIVE}, false),
						},
						"ike_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ike_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ike_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24}, false),
						},
						"ike_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      86400,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
						"ike_local_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
						"ike_remote_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 100),
						},
					},
				},
			},

			"ipsec_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipsec_enc_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_ENC_AES,
							ValidateFunc: validation.StringInSlice([]string{VPN_ENC_AES, VPN_ENC_AES_3DES, VPN_ENC_AES_192, VPN_ENC_AES_256, VPN_ENC_AES_DES}, false),
						},
						"ipsec_auth_alg": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_AUTH_SHA,
							ValidateFunc: validation.StringInSlice([]string{VPN_AUTH_SHA, VPN_AUTH_MD5, VPN_AUTH_SHA256, VPN_AUTH_SHA386, VPN_AUTH_SHA512}, false),
						},
						"ipsec_pfs": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      VPN_PFS_G2,
							ValidateFunc: validation.StringInSlice([]string{VPN_PFS_G1, VPN_PFS_G2, VPN_PFS_G5, VPN_PFS_G14, VPN_PFS_G24, VPN_PFS_DISABLED}, false),
						},
						"ipsec_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 86400),
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunVpnConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request, err := buildAliyunVpnConnectionArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	var response *vpc.CreateVpnConnectionResponse
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateVpnConnection(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateVpnConnectionResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_connection", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.VpnConnectionId)

	if err := vpnGatewayService.WaitForVpnConnection(d.Id(), Null, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	response, err := vpnGatewayService.DescribeVpnConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("customer_gateway_id", response.CustomerGatewayId)
	d.Set("vpn_gateway_id", response.VpnGatewayId)
	d.Set("name", response.Name)

	localSubnet := strings.Split(response.LocalSubnet, ",")
	d.Set("local_subnet", localSubnet)

	remoteSubnet := strings.Split(response.RemoteSubnet, ",")
	d.Set("remote_subnet", remoteSubnet)

	d.Set("effect_immediately", response.EffectImmediately)
	d.Set("status", response.Status)

	if err := d.Set("ike_config", vpnGatewayService.ParseIkeConfig(response.IkeConfig)); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ipsec_config", vpnGatewayService.ParseIpsecConfig(response.IpsecConfig)); err != nil {
		return WrapError(err)
	}

	return nil
}

func resourceAliyunVpnConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateModifyVpnConnectionAttributeRequest()
	request.RegionId = client.RegionId
	request.ClientToken = buildClientToken(request.GetActionName())
	request.VpnConnectionId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
	}

	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	/* If not set effect_immediately value, VPN connection will automatically set the value to false*/
	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if d.HasChange("ike_config") {
		ike_config, err := vpnGatewayService.AssembleIkeConfig(d.Get("ike_config").([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.IkeConfig = ike_config
	}

	if d.HasChange("ipsec_config") {
		ipsec_config, err := vpnGatewayService.AssembleIpsecConfig(d.Get("ipsec_config").([]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request.IpsecConfig = ipsec_config
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyVpnConnectionAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAliyunVpnConnectionRead(d, meta)
}

func resourceAliyunVpnConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteVpnConnectionRequest()
	request.RegionId = client.RegionId
	request.ClientToken = buildClientToken(request.GetActionName())
	request.VpnConnectionId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteVpnConnection(&args)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidVpnConnectionInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpnGatewayService.WaitForVpnConnection(d.Id(), Deleted, DefaultTimeout))
}

func buildAliyunVpnConnectionArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateVpnConnectionRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	request := vpc.CreateCreateVpnConnectionRequest()
	request.RegionId = client.RegionId
	request.CustomerGatewayId = d.Get("customer_gateway_id").(string)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.LocalSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("local_subnet").(*schema.Set).List())
	request.RemoteSubnet = vpnGatewayService.AssembleNetworkSubnetToString(d.Get("remote_subnet").(*schema.Set).List())

	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v, ok := d.GetOk("effect_immediately"); ok {
		request.EffectImmediately = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("ike_config"); ok {
		ikeConfig, err := vpnGatewayService.AssembleIkeConfig(v.([]interface{}))
		if err != nil {
			return nil, WrapError(err)
		}
		request.IkeConfig = ikeConfig
	}

	if v, ok := d.GetOk("ipsec_config"); ok {
		ipsecConfig, err := vpnGatewayService.AssembleIpsecConfig(v.([]interface{}))
		if err != nil {
			return nil, fmt.Errorf("wrong ipsec_config: %#v", err)
		}
		request.IpsecConfig = ipsecConfig
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	return request, nil
}
