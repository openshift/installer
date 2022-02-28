package alicloud

import (
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSlbServerCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSlbServerCertificateCreate,
		Read:   resourceAlicloudSlbServerCertificateRead,
		Update: resourceAlicloudSlbServerCertificateUpdate,
		Delete: resourceAlicloudSlbServerCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_certificate": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: slbServerCertificateDiffSuppressFunc,
			},
			"alicloud_certificate_region_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alicloud_certifacte_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'alicloud_certifacte_id' has been deprecated from provider version 1.68.0. Use 'alicloud_certificate_id' replaces it.",
			},
			"alicloud_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"alicloud_certifacte_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Field 'alicloud_certifacte_name' has been deprecated from provider version 1.68.0. Use 'alicloud_certificate_name' replaces it.",
			},
			"alicloud_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudSlbServerCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := slb.CreateUploadServerCertificateRequest()
	request.RegionId = client.RegionId

	if val, ok := d.GetOk("name"); ok && val != "" {
		request.ServerCertificateName = val.(string)
	}

	if val, ok := d.GetOk("server_certificate"); ok && val != "" {
		request.ServerCertificate = val.(string)
	}

	if val, ok := d.GetOk("private_key"); ok && val != "" {
		request.PrivateKey = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_region_id"); ok && val != "" {
		request.AliCloudCertificateRegionId = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_id"); ok && val != "" {
		request.AliCloudCertificateId = val.(string)
	}

	if val, ok := d.GetOk("alicloud_certificate_name"); ok && val != "" {
		request.AliCloudCertificateName = val.(string)
	}

	if val, ok := d.GetOk("resource_group_id"); ok && val != "" {
		request.ResourceGroupId = val.(string)
	}

	// check server_certificate and private_key
	if request.AliCloudCertificateId == "" {
		if val := strings.Trim(request.ServerCertificate, " "); val == "" {
			return WrapError(Error("UploadServerCertificate got an error, as server_certificate should be not null when alicloud_certificate_id is null."))
		}

		if val := strings.Trim(request.PrivateKey, " "); val == "" {
			return WrapError(Error("UploadServerCertificate got an error, as either private_key or private_file  should be not null when alicloud_certificate_id is null."))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.UploadServerCertificate(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.UploadServerCertificateResponse)
	d.SetId(response.ServerCertificateId)

	return resourceAlicloudSlbServerCertificateUpdate(d, meta)
}

func resourceAlicloudSlbServerCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	tags, err := slbService.DescribeTags(d.Id(), nil, TagResourceCertificate)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", slbService.tagsToMap(tags))

	serverCertificate, err := slbService.DescribeSlbServerCertificate(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if err := d.Set("name", serverCertificate.ServerCertificateName); err != nil {
		return WrapError(err)
	}

	if serverCertificate.AliCloudCertificateId != "" {
		if err := d.Set("alicloud_certificate_id", serverCertificate.AliCloudCertificateId); err != nil {
			return WrapError(err)
		}
	}

	if serverCertificate.AliCloudCertificateName != "" {
		if err := d.Set("alicloud_certificate_name", serverCertificate.AliCloudCertificateName); err != nil {
			return WrapError(err)
		}
	}

	if serverCertificate.ResourceGroupId != "" {
		if err := d.Set("resource_group_id", serverCertificate.ResourceGroupId); err != nil {
			return WrapError(err)
		}
	}

	return nil
}

func resourceAlicloudSlbServerCertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	if err := slbService.setInstanceTags(d, TagResourceCertificate); err != nil {
		return WrapError(err)
	}
	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudSlbServerCertificateRead(d, meta)
	}
	if !d.IsNewResource() && d.HasChange("name") {
		request := slb.CreateSetServerCertificateNameRequest()
		request.RegionId = client.RegionId
		request.ServerCertificateId = d.Id()
		request.ServerCertificateName = d.Get("name").(string)
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.SetServerCertificateName(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}
	return resourceAlicloudSlbServerCertificateRead(d, meta)
}

func resourceAlicloudSlbServerCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}

	request := slb.CreateDeleteServerCertificateRequest()
	request.RegionId = client.RegionId
	request.ServerCertificateId = d.Id()
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DeleteServerCertificate(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"CertificateAndPrivateKeyIsRefered"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)

		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ServerCertificateId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return WrapError(slbService.WaitForSlbServerCertificate(d.Id(), Deleted, DefaultTimeoutMedium))

}
