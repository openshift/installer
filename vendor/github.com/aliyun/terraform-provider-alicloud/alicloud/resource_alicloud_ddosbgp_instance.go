package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddosbgp"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDdosbgpInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdosbgpInstanceCreate,
		Read:   resourceAlicloudDdosbgpInstanceRead,
		Update: resourceAlicloudDdosbgpInstanceUpdate,
		Delete: resourceAlicloudDdosbgpInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				Default:      string(Enterprise),
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enterprise", "Professional"}, false),
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"base_bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  20,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"ip_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, false),
			},
			"ip_count": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      12,
				ForceNew:     true,
			},
		},
	}
}

func resourceAlicloudDdosbgpInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := buildDdosbgpCreateRequest(client.RegionId, d, meta)
	var response *bssopenapi.CreateInstanceResponse
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithBssopenapiClient(func(bssopenapiClient *bssopenapi.Client) (interface{}, error) {
			return bssopenapiClient.CreateInstance(request)
		})

		if err != nil {
			if IsExpectedErrors(err, []string{"NotApplicable"}) {
				request.RegionId = string(connectivity.APSouthEast1)
				request.Domain = connectivity.BssOpenAPIEndpointInternational
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)

		response = raw.(*bssopenapi.CreateInstanceResponse)

		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddosbgp_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	instanceId := response.Data.InstanceId
	if !response.Success {
		return WrapError(Error(response.Message))
	}
	d.SetId(instanceId)

	return resourceAlicloudDdosbgpInstanceUpdate(d, meta)
}

func resourceAlicloudDdosbgpInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosbgpService := DdosbgpService{client}
	insInfo, err := ddosbgpService.DescribeDdosbgpInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	specInfo, err := ddosbgpService.DescribeDdosbgpInstanceSpec(d.Id(), client.RegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	ddosbgpInstanceType := string(Enterprise)
	if insInfo.InstanceType == "0" {
		ddosbgpInstanceType = string(Professional)
	}

	d.Set("name", insInfo.Remark)
	d.Set("region", specInfo.Region)
	d.Set("bandwidth", specInfo.PackConfig.PackAdvThre)
	d.Set("base_bandwidth", specInfo.PackConfig.PackBasicThre)
	d.Set("ip_type", insInfo.IpType)
	d.Set("ip_count", specInfo.PackConfig.IpSpec)
	d.Set("type", ddosbgpInstanceType)

	return nil
}

func resourceAlicloudDdosbgpInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	if d.HasChange("name") {
		request := ddosbgp.CreateModifyRemarkRequest()
		request.InstanceId = d.Id()
		request.RegionId = client.RegionId
		request.ResourceRegionId = client.RegionId

		request.Remark = d.Get("name").(string)

		raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
			return ddosbgpClient.ModifyRemark(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAlicloudDdosbgpInstanceRead(d, meta)
}

func resourceAlicloudDdosbgpInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddosbgpService := DdosbgpService{client}

	request := ddosbgp.CreateReleaseInstanceRequest()
	request.InstanceId = d.Id()
	request.RegionId = client.RegionId

	raw, err := client.WithDdosbgpClient(func(ddosbgpClient *ddosbgp.Client) (interface{}, error) {
		return ddosbgpClient.ReleaseInstance(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(ddosbgpService.WaitForDdosbgpInstance(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildDdosbgpCreateRequest(region string, d *schema.ResourceData, meta interface{}) *bssopenapi.CreateInstanceRequest {
	request := bssopenapi.CreateCreateInstanceRequest()
	request.ProductCode = "ddos"
	request.ProductType = "ddosbgp"
	request.SubscriptionType = "Subscription"
	request.Period = requests.NewInteger(d.Get("period").(int))

	ddosbgpInstanceType := "1"
	if d.Get("type").(string) == string(Professional) {
		ddosbgpInstanceType = "0"
	}

	baseBandWidth := d.Get("base_bandwidth").(int)
	bandWidth := d.Get("bandwidth").(int)
	ipCount := d.Get("ip_count").(int)

	ddosbgpInstanceIpType := "v4"
	if d.Get("ip_type").(string) == string(IPv6) {
		ddosbgpInstanceIpType = "v6"
	}

	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		{
			Code:  "Type",
			Value: ddosbgpInstanceType,
		},
		{
			Code:  "Region",
			Value: region,
		},
		{
			Code:  "IpType",
			Value: ddosbgpInstanceIpType,
		},
		{
			Code:  "BaseBandwidth",
			Value: strconv.Itoa(baseBandWidth),
		},
		{
			Code:  "Bandwidth",
			Value: strconv.Itoa(bandWidth),
		},
		{
			Code:  "IpCount",
			Value: strconv.Itoa(ipCount),
		},
	}

	return request
}
