package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cloudapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunApigatewayVpc() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunApigatewayVpcAccessCreate,
		Read:   resourceAliyunApigatewayVpcAccessRead,
		Delete: resourceAliyunApigatewayVpcAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunApigatewayVpcAccessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := cloudapi.CreateSetVpcAccessRequest()
	request.RegionId = client.RegionId
	request.Name = d.Get("name").(string)
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))
	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.SetVpcAccess(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_api_gateway_vpc", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	d.SetId(fmt.Sprintf("%s%s%s%s%s%s%s", request.Name, COLON_SEPARATED, request.VpcId, COLON_SEPARATED, request.InstanceId, COLON_SEPARATED, request.Port))
	return resourceAliyunApigatewayVpcAccessRead(d, meta)
}

func resourceAliyunApigatewayVpcAccessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudApiService := CloudApiService{client}

	vpc, err := cloudApiService.DescribeApiGatewayVpcAccess(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", vpc.Name)
	d.Set("vpc_id", vpc.VpcId)
	d.Set("instance_id", vpc.InstanceId)
	d.Set("port", vpc.Port)

	return nil
}

func resourceAliyunApigatewayVpcAccessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := cloudapi.CreateRemoveVpcAccessRequest()
	request.RegionId = client.RegionId
	request.VpcId = d.Get("vpc_id").(string)
	request.InstanceId = d.Get("instance_id").(string)
	request.Port = requests.NewInteger(d.Get("port").(int))

	raw, err := client.WithCloudApiClient(func(cloudApiClient *cloudapi.Client) (interface{}, error) {
		return cloudApiClient.RemoveVpcAccess(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return nil

}
