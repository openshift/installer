package alicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunVpnCustomerGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunVpnCustomerGatewayCreate,
		Read:   resourceAliyunVpnCustomerGatewayRead,
		Update: resourceAliyunVpnCustomerGatewayUpdate,
		Delete: resourceAliyunVpnCustomerGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.SingleIP(),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
		},
	}
}

func resourceAliyunVpnCustomerGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateCreateCustomerGatewayRequest()
	request.RegionId = client.RegionId
	request.IpAddress = d.Get("ip_address").(string)
	if v := d.Get("name").(string); v != "" {
		request.Name = v
	}

	if v := d.Get("description").(string); v != "" {
		request.Description = v
	}
	request.ClientToken = buildClientToken(request.GetActionName())

	wait := incrementalWait(3*time.Second, 5*time.Second)
	var raw interface{}
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.CreateCustomerGateway(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling, "OperationConflict"}) {
				wait()
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_vpn_customer_gateway", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	response, _ := raw.(*vpc.CreateCustomerGatewayResponse)

	d.SetId(response.CustomerGatewayId)

	err = vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Null, 60)
	if err != nil {
		return WrapError(err)
	}
	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}

	object, err := vpnGatewayService.DescribeVpnCustomerGateway(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("ip_address", object.IpAddress)
	d.Set("name", object.Name)
	d.Set("description", object.Description)

	return nil
}

func resourceAliyunVpnCustomerGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := vpc.CreateModifyCustomerGatewayAttributeRequest()
	request.RegionId = client.RegionId
	request.CustomerGatewayId = d.Id()
	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.ModifyCustomerGatewayAttribute(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAliyunVpnCustomerGatewayRead(d, meta)
}

func resourceAliyunVpnCustomerGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpnGatewayService := VpnGatewayService{client}
	request := vpc.CreateDeleteCustomerGatewayRequest()
	request.RegionId = client.RegionId
	request.CustomerGatewayId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteCustomerGateway(&args)
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
		if IsExpectedErrors(err, []string{"InvalidCustomerGatewayInstanceId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpnGatewayService.WaitForVpnCustomerGateway(d.Id(), Deleted, DefaultTimeout))
}
