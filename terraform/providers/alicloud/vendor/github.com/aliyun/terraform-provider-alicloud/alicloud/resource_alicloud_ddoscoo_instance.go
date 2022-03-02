package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudDdoscooInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDdoscooInstanceCreate,
		Read:   resourceAlicloudDdoscooInstanceRead,
		Update: resourceAlicloudDdoscooInstanceUpdate,
		Delete: resourceAlicloudDdoscooInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"base_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_bandwidth": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain_count": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:     true,
				Default:      1,
				ForceNew:     true,
			},
			"product_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ddoscoo",
				ValidateFunc: validation.StringInSlice([]string{"ddoscoo", "ddoscoo_intl"}, false),
			},
		},
	}
}

func resourceAlicloudDdoscooInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	var response map[string]interface{}
	action := "CreateInstance"
	request := make(map[string]interface{})
	conn, err := client.NewBssopenapiClient()
	if err != nil {
		return WrapError(err)
	}
	request["Period"] = requests.NewInteger(d.Get("period").(int))
	request["ProductCode"] = "ddos"
	request["ProductType"] = d.Get("product_type").(string)
	request["SubscriptionType"] = "Subscription"
	request["Parameter"] = []map[string]string{
		{
			"Code":  "ServicePartner",
			"Value": "coop-line-001",
		},
		{
			"Code":  "Bandwidth",
			"Value": d.Get("bandwidth").(string),
		},
		{
			"Code":  "BaseBandwidth",
			"Value": d.Get("base_bandwidth").(string),
		},
		{
			"Code":  "DomainCount",
			"Value": d.Get("domain_count").(string),
		},
		{
			"Code":  "PortCount",
			"Value": d.Get("port_count").(string),
		},
		{
			"Code":  "ServiceBandwidth",
			"Value": d.Get("service_bandwidth").(string),
		},
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-12-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			if err != nil {
				if IsExpectedErrors(err, []string{"NotApplicable"}) {
					conn.Endpoint = String(connectivity.BssOpenAPIEndpointInternational)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ddoscoo_instance", action, AlibabaCloudSdkGoERROR)
	}
	if response["Code"].(string) != "Success" {
		return WrapErrorf(fmt.Errorf("%v", response), DefaultErrorMsg, "alicloud_ddoscoo_instance", action, AlibabaCloudSdkGoERROR)
	}
	response = response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(response["InstanceId"]))
	stateConf := BuildStateConf([]string{"Pending"}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ddoscooService.DdosStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudDdoscooInstanceUpdate(d, client)
}

func resourceAlicloudDdoscooInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	insInfo, err := ddoscooService.DescribeDdoscooInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	specInfo, err := ddoscooService.DescribeDdoscooInstanceSpec(d)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}

		return WrapError(err)
	}

	d.Set("name", insInfo.Remark)
	d.Set("bandwidth", strconv.Itoa(specInfo.ElasticBandwidth))
	d.Set("base_bandwidth", strconv.Itoa(specInfo.BaseBandwidth))
	d.Set("domain_count", strconv.Itoa(specInfo.DomainLimit))
	d.Set("port_count", strconv.Itoa(specInfo.PortLimit))
	d.Set("service_bandwidth", strconv.Itoa(specInfo.BandwidthMbps))

	return nil
}

func resourceAlicloudDdoscooInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ddoscooService := DdoscooService{client}

	d.Partial(true)

	if d.HasChange("name") {
		if err := ddoscooService.UpdateDdoscooInstanceName(d.Id(), d.Get("name").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("name")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDdoscooInstanceRead(d, meta)
	}

	if d.HasChange("bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("bandwidth", "Bandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("bandwidth")
	}

	if d.HasChange("base_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("base_bandwidth", "BaseBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("base_bandwidth")
	}

	if d.HasChange("domain_count") {
		if err := ddoscooService.UpdateInstanceSpec("domain_count", "DomainCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("domain_count")
	}

	if d.HasChange("port_count") {
		if err := ddoscooService.UpdateInstanceSpec("port_count", "PortCount", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("port_count")
	}

	if d.HasChange("service_bandwidth") {
		if err := ddoscooService.UpdateInstanceSpec("service_bandwidth", "ServiceBandwidth", d, meta); err != nil {
			return WrapError(err)
		}

		d.SetPartial("service_bandwidth")
	}

	d.Partial(false)
	return resourceAlicloudDdoscooInstanceRead(d, meta)
}

type JsonObject struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func resourceAlicloudDdoscooInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ddoscoo.CreateReleaseInstanceRequest()
	request.RegionId = "cn-hangzhou"
	request.InstanceId = d.Id()

	raw, err := client.WithDdoscooClient(func(ddoscooClient *ddoscoo.Client) (interface{}, error) {
		return ddoscooClient.ReleaseInstance(request)
	})

	var DeleteResult JsonObject
	if err != nil {
		_ = json.Unmarshal([]byte(raw.(*ddoscoo.ReleaseInstanceResponse).GetHttpContentString()), &DeleteResult)

		if DeleteResult.Code == "InstanceNotExpire" {
			log.Printf("[INFO]  instance cannot be deleted and must wait it to be expired and release it automatically")
			return WrapError(Error("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically."))
		}
		if IsExpectedErrors(err, []string{"InstanceNotFound"}) {
			return nil
		}

		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
