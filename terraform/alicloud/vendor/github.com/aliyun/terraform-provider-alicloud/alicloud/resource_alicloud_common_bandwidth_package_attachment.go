package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliyunCommonBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunCommonBandwidthPackageAttachmentCreate,
		Read:   resourceAliyunCommonBandwidthPackageAttachmentRead,
		Delete: resourceAliyunCommonBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliyunCommonBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request := vpc.CreateAddCommonBandwidthPackageIpRequest()
	request.RegionId = client.RegionId
	request.BandwidthPackageId = Trim(d.Get("bandwidth_package_id").(string))
	request.IpInstanceId = Trim(d.Get("instance_id").(string))
	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.AddCommonBandwidthPackageIp(request)
	})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.AddCommonBandwidthPackageIp(request)
		})
		//Waiting for unassociate the common bandwidth package
		if err != nil && !IsExpectedErrors(err, []string{"IpInstanceId.AlreadyInBandwidthPackage"}) {
			if IsExpectedErrors(err, []string{"TaskConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_common_bandwidth_package_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	//check the common bandwidth package attachment
	d.SetId(request.BandwidthPackageId + COLON_SEPARATED + request.IpInstanceId)
	if err := vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Available, 5*DefaultTimeout); err != nil {
		return WrapError(err)
	}
	return resourceAliyunCommonBandwidthPackageAttachmentRead(d, meta)
}

func resourceAliyunCommonBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]
	_, err = vpcService.DescribeCommonBandwidthPackageAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("bandwidth_package_id", bandwidthPackageId)
	d.Set("instance_id", ipInstanceId)
	return nil
}

func resourceAliyunCommonBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	bandwidthPackageId, ipInstanceId := parts[0], parts[1]

	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.RemoveCommonBandwidthPackageIp(request)
		})
		//Waiting for unassociate the common bandwidth package
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "OperationConflict"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"OperationUnsupported.IpNotInCbwp"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForCommonBandwidthPackageAttachment(d.Id(), Deleted, DefaultTimeoutMedium))
}
