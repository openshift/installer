package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudCenBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCenBandwidthPackageCreate,
		Read:   resourceAlicloudCenBandwidthPackageRead,
		Update: resourceAlicloudCenBandwidthPackageUpdate,
		Delete: resourceAlicloudCenBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"cen_bandwidth_package_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validation.StringLenBetween(2, 128),
				ConflictsWith: []string{"cen_bandwidth_package_name"},
				Deprecated:    "Field 'name' has been deprecated from version 1.98.0. Use 'cen_bandwidth_package_name' and instead.",
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"geographic_region_a_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Asia-Pacific", "Australia", "China", "Europe", "Middle-East", "North-America"}, false),
				ConflictsWith: []string{"geographic_region_ids"},
			},
			"geographic_region_b_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"Asia-Pacific", "Australia", "China", "Europe", "Middle-East", "North-America"}, false),
				ConflictsWith: []string{"geographic_region_ids"},
			},
			"geographic_region_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems:      2,
				MinItems:      1,
				ConflictsWith: []string{"geographic_region_a_id", "geographic_region_b_id"},
				Deprecated:    "Field 'geographic_region_ids' has been deprecated from version 1.98.0. Use 'geographic_region_a_id' and 'geographic_region_b_id' instead.",
			},
			"payment_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"charge_type"},
			},
			"charge_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validation.StringInSlice([]string{"PostPaid", "PrePaid"}, false),
				ConflictsWith: []string{"payment_type"},
				Deprecated:    "Field 'charge_type' has been deprecated from version 1.98.0. Use 'payment_type' and instead.",
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 6, 12}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return PayType(d.Get("charge_type").(string)) == PostPaid || PayType(d.Get("payment_type").(string)) == PostPaid
				},
			},
			"expired_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCenBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}

	request := cbn.CreateCreateCenBandwidthPackageRequest()
	request.AutoPay = requests.NewBoolean(true)

	request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
	if v, ok := d.GetOk("cen_bandwidth_package_name"); ok {
		request.Name = v.(string)
	} else if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}

	if d.Get("geographic_region_a_id").(string) != "" && d.Get("geographic_region_b_id").(string) != "" {
		request.GeographicRegionAId = d.Get("geographic_region_a_id").(string)
		request.GeographicRegionBId = d.Get("geographic_region_b_id").(string)
	} else if len(d.Get("geographic_region_ids").(*schema.Set).List()) > 0 {
		geographicRegionId := d.Get("geographic_region_ids").(*schema.Set).List()
		if len(geographicRegionId) == 1 {
			request.GeographicRegionAId = geographicRegionId[0].(string)
			request.GeographicRegionBId = geographicRegionId[0].(string)
		} else if len(geographicRegionId) == 2 {
			if geographicRegionId[1].(string) == "China" {
				request.GeographicRegionAId = geographicRegionId[1].(string)
				request.GeographicRegionBId = geographicRegionId[0].(string)
			} else {
				request.GeographicRegionAId = geographicRegionId[0].(string)
				request.GeographicRegionBId = geographicRegionId[1].(string)
			}
		}
	} else {
		return WrapError(Error(`[ERROR] Argument "geographic_region_a_id" and "geographic_region_b_id" must be set!`))
	}

	if v, ok := d.GetOk("payment_type"); ok {
		if v.(string) == "PrePaid" {
			request.BandwidthPackageChargeType = convertPaymentTypeRequest(v.(string))
			request.PricingCycle = "Month"
		} else {
			request.BandwidthPackageChargeType = convertPaymentTypeRequest("PostPaid")
		}
	} else if v, ok := d.GetOk("charge_type"); ok {
		if v.(string) == "PrePaid" {
			request.BandwidthPackageChargeType = convertPaymentTypeRequest(v.(string))
			request.PricingCycle = "Month"
		} else {
			request.BandwidthPackageChargeType = convertPaymentTypeRequest("PostPaid")
		}
	}

	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.CreateCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"Operation.Blocking"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		response, _ := raw.(*cbn.CreateCenBandwidthPackageResponse)
		d.SetId(fmt.Sprintf("%v", response.CenBandwidthPackageId))
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cen_bandwidth_package", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{}, []string{"Idle"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, cbnService.CenBandwidthPackageStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}
func resourceAlicloudCenBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cbnService := CbnService{client}
	object, err := cbnService.DescribeCenBandwidthPackage(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cen_bandwidth_package cbnService.DescribeCenBandwidthPackage Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	geographicRegionIds := make([]string, 0)
	geographicRegionIds = append(geographicRegionIds, convertGeographicRegionAIdResponse(object.GeographicRegionAId))
	geographicRegionIds = append(geographicRegionIds, convertGeographicRegionBIdResponse(object.GeographicRegionBId))
	d.Set("geographic_region_ids", geographicRegionIds)
	d.Set("bandwidth", object.Bandwidth)
	d.Set("cen_bandwidth_package_name", object.Name)
	d.Set("name", object.Name)
	d.Set("description", object.Description)
	d.Set("geographic_region_a_id", convertGeographicRegionAIdResponse(object.GeographicRegionAId))
	d.Set("geographic_region_b_id", convertGeographicRegionBIdResponse(object.GeographicRegionBId))
	d.Set("payment_type", convertPaymentTypeResponse(object.BandwidthPackageChargeType))
	d.Set("charge_type", convertPaymentTypeResponse(object.BandwidthPackageChargeType))
	//if convertPaymentTypeResponse(object.BandwidthPackageChargeType) == "PrePaid" {
	//	period, err := computePeriodByUnit(object.CreationTime, object.ExpiredTime, d.Get("period").(int), "Month")
	//	if err != nil {
	//		return WrapError(err)
	//	}
	//	d.Set("period", period)
	//}
	d.Set("expired_time", object.ExpiredTime)
	d.Set("status", object.Status)
	return nil
}
func resourceAlicloudCenBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cenService := CenService{client}
	d.Partial(true)

	if d.HasChange("bandwidth") {
		request := cbn.CreateModifyCenBandwidthPackageSpecRequest()
		request.CenBandwidthPackageId = d.Id()
		request.Bandwidth = requests.NewInteger(d.Get("bandwidth").(int))
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.ModifyCenBandwidthPackageSpec(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidStatus.Resource", "Operation.Blocking"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		if err := cenService.WaitForCenBandwidthPackage(d.Id(), Idle, d.Get("bandwidth").(int), DefaultCenTimeout); err != nil {
			return WrapError(err)
		}
		d.SetPartial("bandwidth")
	}
	update := false
	request := cbn.CreateModifyCenBandwidthPackageAttributeRequest()
	request.CenBandwidthPackageId = d.Id()
	if d.HasChange("cen_bandwidth_package_name") {
		update = true
		request.Name = d.Get("cen_bandwidth_package_name").(string)
	}
	if d.HasChange("name") {
		update = true
		request.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		update = true
		request.Description = d.Get("description").(string)
	}
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.ModifyCenBandwidthPackageAttribute(request)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{"InvalidStatus.Resource", "Operation.Blocking"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(request.GetActionName(), raw)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("cen_bandwidth_package_name")
		d.SetPartial("name")
		d.SetPartial("description")
	}
	d.Partial(false)
	return resourceAlicloudCenBandwidthPackageRead(d, meta)
}
func resourceAlicloudCenBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("payment_type").(string) == "PrePaid" {
		log.Printf("[WARN] Cannot destroy resource Alicloud Resource Cen BandwidthPackage. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	client := meta.(*connectivity.AliyunClient)
	request := cbn.CreateDeleteCenBandwidthPackageRequest()
	request.CenBandwidthPackageId = d.Id()
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.WithCbnClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenBandwidthPackage(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation.CenBandwidthLimitsNotZero", "ParameterBwpInstanceId", "Forbidden.Release"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ParameterBwpInstanceId"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return nil
}
func convertPaymentTypeRequest(source string) string {
	switch source {
	case "PostPaid":
		return "POSTPAY"
	case "PrePaid":
		return "PREPAY"
	}
	return source
}

func convertGeographicRegionAIdResponse(source string) string {
	switch source {
	case "asia-pacific":
		return "Asia-Pacific"
	case "china":
		return "China"
	case "europe":
		return "Europe"
	case "middle-east":
		return "Middle-East"
	case "north-america":
		return "North-America"
	case "australia":
		return "Australia"
	}
	return source
}
func convertGeographicRegionBIdResponse(source string) string {
	switch source {
	case "asia-pacific":
		return "Asia-Pacific"
	case "china":
		return "China"
	case "europe":
		return "Europe"
	case "middle-east":
		return "Middle-East"
	case "north-america":
		return "North-America"
	case "australia":
		return "Australia"
	}
	return source
}
func convertPaymentTypeResponse(source string) string {
	switch source {
	case "POSTPAY":
		return "PostPaid"
	case "PREPAY":
		return "PrePaid"
	}
	return source
}
