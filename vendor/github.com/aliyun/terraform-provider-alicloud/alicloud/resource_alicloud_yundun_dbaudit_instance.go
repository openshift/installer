package alicloud

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/bssopenapi"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	RELEASE_HANG_MINS = 3
)

func resourceAlicloudDbauditInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDbauditInstanceCreate,
		Read:   resourceAlicloudDbauditInstanceRead,
		Update: resourceAlicloudDbauditInstanceUpdate,
		Delete: resourceAlicloudDbauditInstanceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"plan_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntInSlice([]int{1, 3, 6, 12, 24, 36}),
				Required:     true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),

			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlicloudDbauditInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := buildDbauditCreateRequest(d, meta)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_yundun_dbaudit_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	instanceId := response.Data.InstanceId
	if !response.Success {
		return WrapError(Error(response.Message))
	}
	d.SetId(instanceId)

	dbauditService := DbauditService{client}

	// check RAM policy
	dbauditService.ProcessRolePolicy()
	// wait for order complete
	stateConf := BuildStateConf([]string{}, []string{"PENDING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	// start instance
	if err := dbauditService.StartDbauditInstance(instanceId, d.Get("vswitch_id").(string)); err != nil {
		return WrapError(err)
	}
	// wait for pending
	stateConf = BuildStateConf([]string{"PENDING", "CREATING"}, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"UPGRADING", "UPGRADE_FAILED", "CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return resourceAlicloudDbauditInstanceUpdate(d, meta)
}

func resourceAlicloudDbauditInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbauditService := DbauditService{client}
	instance, err := dbauditService.DescribeYundunDbauditInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("description", instance.Description)
	//period, err := computePeriodByUnit(instance.StartTime/1000, instance.ExpireTime/1000, d.Get("period").(int), "Month")
	//if err != nil {
	//	return WrapError(err)
	//}
	//d.Set("period", period)
	d.Set("plan_code", instance.LicenseCode)
	d.Set("region_id", client.RegionId)
	d.Set("vswitch_id", instance.VswitchId)

	tags, err := dbauditService.DescribeTags(d.Id(), nil, TagResourceInstance)
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", dbauditService.tagsToMap(tags))
	return nil
}

func resourceAlicloudDbauditInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dbauditService := DbauditService{client}

	d.Partial(true)

	if d.HasChange("tags") {
		if err := dbauditService.setInstanceTags(d, TagResourceInstance); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}

	if d.HasChange("description") {
		if err := dbauditService.UpdateDbauditInstanceDescription(d.Id(), d.Get("description").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("description")
	}

	if d.HasChange("resource_group_id") {
		if err := dbauditService.UpdateResourceGroup(d.Id(), d.Get("resource_group_id").(string)); err != nil {
			return WrapError(err)
		}
		d.SetPartial("resource_group_id")
	}

	if d.IsNewResource() {
		d.Partial(false)
		return resourceAlicloudDbauditInstanceRead(d, meta)
	}

	if d.HasChange("plan_code") {
		if err := dbauditService.UpdateInstanceSpec("plan_code", "PlanCode", d, meta); err != nil {
			return WrapError(err)
		}
		stateConf := BuildStateConf([]string{"UPGRADING"}, []string{"PENDING", "RUNNING"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second, dbauditService.DbauditInstanceRefreshFunc(d.Id(), []string{"CREATING", "UPGRADE_FAILED", "CREATE_FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
		d.SetPartial("plan_code")
	}

	d.Partial(false)
	// wait for order complete
	return resourceAlicloudDbauditInstanceRead(d, meta)
}

func resourceAlicloudDbauditInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudDbauditInstance. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}

func buildDbauditCreateRequest(d *schema.ResourceData, meta interface{}) *bssopenapi.CreateInstanceRequest {
	request := bssopenapi.CreateCreateInstanceRequest()
	request.ProductCode = "dbaudit"
	request.SubscriptionType = "Subscription"
	request.Period = requests.NewInteger(d.Get("period").(int))
	client := meta.(*connectivity.AliyunClient)

	request.Parameter = &[]bssopenapi.CreateInstanceParameter{
		// force to buy vpc version
		{
			Code:  "NetworkType",
			Value: "vpc",
		},
		{
			Code:  "SeriesCode",
			Value: "alpha",
		},
		{
			Code:  "PlanCode",
			Value: d.Get("plan_code").(string),
		},
		{
			Code:  "RegionId",
			Value: client.RegionId,
		},
	}
	return request
}
