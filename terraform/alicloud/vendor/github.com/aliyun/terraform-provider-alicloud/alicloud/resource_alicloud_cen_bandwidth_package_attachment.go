package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCenBandwidthPackageAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageAttachmentCreate,
		Read:   resourceAlicloudCenBandwidthPackageAttachmentRead,
		Delete: resourceAlicloudCenBandwidthPackageAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"bandwidth_package_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateAssociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.AssociateCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.BwpInstanceStatus", "InvalidOperation.BwpBusinessStatus", "InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_bandwidth_package_attachment", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	d.SetId(cenBwpId)
	if err := cenService.WaitForCenBandwidthPackageAttachment(d.Id(), InUse, DefaultCenTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudCenBandwidthPackageAttachmentRead(d, meta)
}

func resourceAlicloudCenBandwidthPackageAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}

	object, err := cenService.DescribeCenBandwidthPackageAttachment(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("instance_id", object.CenIds.CenId[0])
	d.Set("bandwidth_package_id", object.CenBandwidthPackageId)

	return nil
}

func resourceAlicloudCenBandwidthPackageAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	cenId := d.Get("instance_id").(string)
	cenBwpId := d.Get("bandwidth_package_id").(string)

	request := cbn.CreateUnassociateCenBandwidthPackageRequest()
	request.CenId = cenId
	request.CenBandwidthPackageId = cenBwpId

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.UnassociateCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.BwpInstanceStatus", "InvalidOperation.BwpBusinessStatus", "InvalidOperation.CenInstanceStatus", "Operation.Blocking"}) {
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

	return WrapError(cenService.WaitForCenBandwidthPackageAttachment(cenBwpId, Deleted, DefaultCenTimeout))
}
