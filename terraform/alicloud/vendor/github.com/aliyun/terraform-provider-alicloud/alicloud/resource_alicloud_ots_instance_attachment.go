package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOtsInstanceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunOtsInstanceAttachmentCreate,
		Read:   resourceAliyunOtsInstanceAttachmentRead,
		Delete: resourceAliyunOtsInstanceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunOtsInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	request := ots.CreateBindInstance2VpcRequest()
	request.RegionId = client.RegionId
	request.InstanceName = d.Get("instance_name").(string)
	request.InstanceVpcName = d.Get("vpc_name").(string)
	request.VirtualSwitchId = d.Get("vswitch_id").(string)

	if vsw, err := vpcService.DescribeVSwitch(d.Get("vswitch_id").(string)); err != nil {
		return WrapError(err)
	} else {
		request.VpcId = vsw.VpcId
	}

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.BindInstance2Vpc(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ots_instance_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId(request.InstanceName)
	return resourceAliyunOtsInstanceAttachmentRead(d, meta)
}

func resourceAliyunOtsInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	// There is a bug that inst does not contain instance name and vswitch ID, so this resource does not support import function.
	//d.Set("instance_name", inst.InstanceName)
	d.Set("vpc_name", object.InstanceVpcName)
	d.Set("vpc_id", object.VpcId)
	return nil
}

func resourceAliyunOtsInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	otsService := OtsService{client}
	object, err := otsService.DescribeOtsInstanceAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}
	request := ots.CreateUnbindInstance2VpcRequest()
	request.RegionId = client.RegionId
	request.InstanceName = d.Id()
	request.InstanceVpcName = object.InstanceVpcName

	raw, err := client.WithOtsClient(func(otsClient *ots.Client) (interface{}, error) {
		return otsClient.UnbindInstance2Vpc(request)
	})
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(otsService.WaitForOtsInstanceVpc(d.Id(), Deleted, DefaultTimeout))
}
