package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudImageExport() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudImageExportCreate,
		Read:   resourceAliCloudImageExportRead,
		Delete: resourceAliCloudImageExportDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"oss_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudImageExportCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client: client}

	request := ecs.CreateExportImageRequest()
	request.RegionId = client.RegionId
	request.ImageId = d.Get("image_id").(string)
	request.OSSBucket = d.Get("oss_bucket").(string)
	request.OSSPrefix = d.Get("oss_prefix").(string)
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ExportImage(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_image_export", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	response := raw.(*ecs.ExportImageResponse)
	taskId := response.TaskId
	d.SetId(request.ImageId)
	stateConf := BuildStateConf([]string{"Waiting", "Processing"}, []string{"Finished"}, d.Timeout(schema.TimeoutCreate), 2*time.Minute, ecsService.TaskStateRefreshFunc(taskId, []string{"CreateFailed", "UnAvailable"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAliCloudImageExportRead(d, meta)

}

func resourceAliCloudImageExportRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client: client}

	object, err := ecsService.DescribeImageById(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("image_id", object.ImageId)
	return WrapError(err)
}

func resourceAliCloudImageExportDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	ossService := OssService{client: client}
	var requestInfo *oss.Client
	raw, err := client.WithOssClient(func(ossClient *oss.Client) (interface{}, error) {
		requestInfo = ossClient
		return ossClient.Bucket(d.Get("oss_bucket").(string))
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "oss_bucket", AliyunOssGoSdk)
	}
	addDebug("oss_bucket", raw, requestInfo, map[string]string{"bucketName": d.Get("oss_bucket").(string)})
	bucket, _ := raw.(*oss.Bucket)

	objectName := fmt.Sprintf(d.Get("oss_prefix").(string) + "_" + d.Id() + "_system.raw.tar.gz")
	err = bucket.DeleteObject(objectName)
	if err != nil {
		if IsExpectedErrors(err, []string{"No Content", "Not Found"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Get("oss_prefix").(string), "DeleteObject", AliyunOssGoSdk)
	}
	addDebug("oss_prefix", nil, requestInfo, map[string]string{"oss_prefix": d.Get("oss_prefix").(string)})
	return WrapError(ossService.WaitForOssBucketObject(bucket, d.Get("oss_prefix").(string), Deleted, DefaultTimeoutMedium))
}
