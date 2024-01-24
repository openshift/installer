package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunSslVpnServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSslVpnServerCreate,
		Read:   resourceAliyunSslVpnServerRead,
		Update: resourceAliyunSslVpnServerUpdate,
		Delete: resourceAliyunSslVpnServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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

			"client_ip_pool": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateCIDRNetworkAddress,
			},

			"local_subnet": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateVpnCIDRNetworkAddress,
			},

			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      VPN_UDP_PROTO,
				ValidateFunc: validation.StringInSlice([]string{VPN_UDP_PROTO, VPN_TCP_PROTO}, false),
			},

			"cipher": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SSL_VPN_ENC_AES_128,
				ValidateFunc: validation.StringInSlice([]string{SSL_VPN_ENC_AES_128, SSL_VPN_ENC_AES_192, SSL_VPN_ENC_AES_256, SSL_VPN_ENC_NONE}, false),
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1194,
				ValidateFunc: validateSslVpnPortValue([]int{22, 2222, 22222, 9000, 9001, 9002, 7505, 80, 443, 53, 68, 123, 4510, 4560, 500, 4500}),
			},
			"compress": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"max_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"internet_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunSslVpnServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateSslVpnServerRequest()
	request.RegionId = string(client.Region)
	request.VpnGatewayId = d.Get("vpn_gateway_id").(string)
	request.ClientIpPool = d.Get("client_ip_pool").(string)
	request.LocalSubnet = d.Get("local_subnet").(string)
	request.Name = d.Get("name").(string)
	request.Proto = d.Get("protocol").(string)
	request.Cipher = d.Get("cipher").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))
	request.Compress = requests.NewBoolean(d.Get("compress").(bool))
	request.ClientToken = buildClientToken(request.GetActionName())

	var response *vpc.CreateSslVpnServerResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateSslVpnServer(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"VpnGateway.Configuring"}) {
				time.Sleep(10 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ = raw.(*vpc.CreateSslVpnServerResponse)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ssl_vpn_server", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(response.SslVpnServerId)
	err = vpnGatewayService.WaitForSslVpnServer(d.Id(), Null, DefaultTimeout)

	if err != nil {
		return WrapError(err)
	}
	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeSslVpnServer(d.Id())

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("vpn_gateway_id", object.VpnGatewayId)
	d.Set("name", object.Name)
	d.Set("local_subnet", object.LocalSubnet)
	d.Set("client_ip_pool", object.ClientIpPool)
	d.Set("cipher", object.Cipher)
	d.Set("protocol", object.Proto)
	d.Set("port", object.Port)
	d.Set("compress", object.Compress)
	d.Set("connections", object.Connections)
	d.Set("max_connections", object.MaxConnections)
	d.Set("internet_ip", object.InternetIp)

	return nil
}

func resourceAliyunSslVpnServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := vpc.CreateModifySslVpnServerRequest()
	request.RegionId = client.RegionId
	request.SslVpnServerId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
	}

	request.ClientIpPool = d.Get("client_ip_pool").(string)
	request.LocalSubnet = d.Get("local_subnet").(string)

	if d.HasChange("protocol") {
		request.Proto = d.Get("protocol").(string)
	}

	if d.HasChange("cipher") {
		request.Cipher = d.Get("cipher").(string)
	}

	if d.HasChange("port") {
		request.Port = requests.NewInteger(d.Get("port").(int))
	}

	if d.HasChange("compress") {
		request.Compress = requests.NewBoolean(d.Get("compress").(bool))
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifySslVpnServer(request)
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAliyunSslVpnServerRead(d, meta)
}

func resourceAliyunSslVpnServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteSslVpnServerRequest()
	request.RegionId = client.RegionId
	request.SslVpnServerId = d.Id()

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteSslVpnServer(request)
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
		if IsExpectedErrors(err, []string{"InvalidSslVpnServerId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpnGatewayService.WaitForSslVpnServer(d.Id(), Deleted, DefaultTimeout))
}
