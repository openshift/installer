package alicloud

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudReservedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudReservedInstanceCreate,
		Read:   resourceAliCloudReservedInstanceRead,
		Update: resourceAliCloudReservedInstanceUpdate,
		Delete: resourceAliCloudReservedInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Region",
				ValidateFunc: validation.StringInSlice([]string{"Region", "Zone"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Windows", "Linux"}, false),
			},
			"period_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Year",
				ValidateFunc: validation.StringInSlice([]string{"Year"}, false),
			},
			"period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1, 3}),
			},
			"offering_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"No Upfront", "Partial Upfront", "All Upfront"}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}
func resourceAliCloudReservedInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreatePurchaseReservedInstancesOfferingRequest()
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = v.(string)
	}
	request.RegionId = client.RegionId
	if scope, ok := d.GetOk("scope"); ok {
		request.Scope = scope.(string)
		if v, ok := d.GetOk("zone_id"); ok {
			if scope == "Zone" && v == "" {
				return WrapError(Error("Required when Scope is Zone."))
			}
			request.ZoneId = v.(string)
		}
	}
	if v, ok := d.GetOk("instance_amount"); ok {
		request.InstanceAmount = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("platform"); ok {
		request.Platform = v.(string)
	}
	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = v.(string)
	}
	if v, ok := d.GetOk("period"); ok {
		request.Period = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("offering_type"); ok {
		request.OfferingType = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		request.ReservedInstanceName = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		request.Description = v.(string)
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request.ResourceGroupId = v.(string)
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.PurchaseReservedInstancesOffering(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_reserved_instance", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response := raw.(*ecs.PurchaseReservedInstancesOfferingResponse)
	if len(response.ReservedInstanceIdSets.ReservedInstanceId) != 1 {
		return WrapError(Error("API returned wrong number of collections"))
	}
	d.SetId(response.ReservedInstanceIdSets.ReservedInstanceId[0])

	if err := ecsService.WaitForReservedInstance(d.Id(), Active, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateModifyReservedInstanceAttributeRequest()

	request.ReservedInstanceId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("name") || d.HasChange("description") {
		if v, ok := d.GetOk("name"); ok {
			request.ReservedInstanceName = v.(string)
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = v.(string)
		}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifyReservedInstanceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	return resourceAliCloudReservedInstanceRead(d, meta)
}
func resourceAliCloudReservedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	reservedInstances, err := ecsService.DescribeReservedInstance(d.Id())
	if err != nil {
		return WrapError(err)
	}
	d.Set("instance_type", reservedInstances.InstanceType)
	d.Set("scope", reservedInstances.Scope)
	d.Set("zone_id", reservedInstances.ZoneId)
	d.Set("instance_amount", reservedInstances.InstanceAmount)
	d.Set("platform", strings.Title(reservedInstances.Platform))
	d.Set("offering_type", reservedInstances.OfferingType)
	d.Set("name", reservedInstances.ReservedInstanceName)
	d.Set("description", reservedInstances.Description)
	d.Set("resource_group_id", reservedInstances.ReservedInstanceId)

	return WrapError(err)
}

func resourceAliCloudReservedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	// PurchaseReservedInstancesOffering can not be release.
	return nil
}
