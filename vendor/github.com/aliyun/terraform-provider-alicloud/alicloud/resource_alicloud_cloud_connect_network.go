package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/smartag"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCloudConnectNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudConnectNetworkCreate,
		Read:   resourceAlicloudCloudConnectNetworkRead,
		Update: resourceAlicloudCloudConnectNetworkUpdate,
		Delete: resourceAlicloudCloudConnectNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
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
			"is_default": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"cidr_block": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateVpnCIDRNetworkAddress,
			},
		},
	}
}

func resourceAlicloudCloudConnectNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := smartag.CreateCreateCloudConnectNetworkRequest()

	request.IsDefault = requests.NewBoolean(d.Get("is_default").(bool))
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		request.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("cidr_block"); ok && v.(string) != "" {
		request.CidrBlock = v.(string)
	}

	raw, err := client.WithSagClient(func(sagClient *smartag.Client) (interface{}, error) {
		return sagClient.CreateCloudConnectNetwork(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_connect_network", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*smartag.CreateCloudConnectNetworkResponse)
	d.SetId(response.CcnId)

	return resourceAlicloudCloudConnectNetworkRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkRead(d *schema.ResourceData, meta interface{}) error {
	ccnService := SagService{meta.(*connectivity.AliyunClient)}
	object, err := ccnService.DescribeCloudConnectNetwork(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("cidr_block", object.CidrBlock)
	d.Set("is_default", object.IsDefault)

	return nil
}

func resourceAlicloudCloudConnectNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	update := false
	request := smartag.CreateModifyCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		update = true
	}
	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		update = true
	}
	if d.HasChange("cidr_block") {
		request.CidrBlock = d.Get("cidr_block").(string)
		update = true
	}
	if update {
		raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
			return ccnClient.ModifyCloudConnectNetwork(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudCloudConnectNetworkRead(d, meta)
}

func resourceAlicloudCloudConnectNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sagService := SagService{client}
	request := smartag.CreateDeleteCloudConnectNetworkRequest()
	request.CcnId = d.Id()

	raw, err := client.WithSagClient(func(ccnClient *smartag.Client) (interface{}, error) {
		return ccnClient.DeleteCloudConnectNetwork(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterCcnInstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(sagService.WaitForCloudConnectNetwork(d.Id(), Deleted, DefaultTimeoutMedium))
}
