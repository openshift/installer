package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunRouteTableAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunRouteTableAttachmentCreate,
		Read:   resourceAliyunRouteTableAttachmentRead,
		Delete: resourceAliyunRouteTableAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunRouteTableAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := vpc.CreateAssociateRouteTableRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = Trim(d.Get("route_table_id").(string))
	request.VSwitchId = Trim(d.Get("vswitch_id").(string))
	request.ClientToken = buildClientToken(request.GetActionName())
	if err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AssociateRouteTable(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_route_table_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetId(request.RouteTableId + COLON_SEPARATED + request.VSwitchId)
	err := vpcService.WaitForRouteTableAttachment(d.Id(), Available, DefaultTimeout)
	if err != nil {
		return WrapError(err)
	}
	if err := vpcService.WaitForVSwitch(request.VSwitchId, Available, DefaultTimeoutMedium); err != nil {
		return WrapError(err)
	}
	return resourceAliyunRouteTableAttachmentRead(d, meta)
}

func resourceAliyunRouteTableAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	object, err := vpcService.DescribeRouteTableAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("route_table_id", object["RouteTableId"])
	d.Set("vswitch_id", parts[1])
	return nil
}

func resourceAliyunRouteTableAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}

	request := vpc.CreateUnassociateRouteTableRequest()
	request.RegionId = client.RegionId
	request.RouteTableId = parts[0]
	request.VSwitchId = parts[1]
	request.ClientToken = buildClientToken(request.GetActionName())
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.UnassociateRouteTable(&args)
		})
		//Waiting for unassociate the route table
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
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
	return WrapError(vpcService.WaitForRouteTableAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
