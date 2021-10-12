package alicloud

import (
	"time"

	"github.com/denverdino/aliyungo/common"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRouterInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRouterInterfaceCreate,
		Read:   resourceAlicloudRouterInterfaceRead,
		Update: resourceAlicloudRouterInterfaceUpdate,
		Delete: resourceAlicloudRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"opposite_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"router_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(VRouter), string(VBR)}, false),
				ForceNew:         true,
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(InitiatingSide), string(AcceptingSide)}, false),
				ForceNew: true,
			},
			"specification": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringInSlice(GetAllRouterInterfaceSpec(), false),
				DiffSuppressFunc: routerInterfaceAcceptsideDiffSuppressFunc,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"health_check_source_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"health_check_target_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: routerInterfaceVBRTypeDiffSuppressFunc,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"period": {
				Type:             schema.TypeInt,
				Optional:         true,
				DiffSuppressFunc: PostPaidDiffSuppressFunc,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 24, 36})),
			},
			"access_point_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_access_point_id' has been deprecated from version 1.11.0.",
			},
			"opposite_access_point_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Attribute 'opposite_access_point_id' has been deprecated from version 1.11.0.",
			},
			"opposite_router_type": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_router_type' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_type' instead.",
			},
			"opposite_router_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_router_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_router_id' instead.",
			},
			"opposite_interface_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_interface_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_interface_id' instead.",
			},
			"opposite_interface_owner_id": {
				Type:       schema.TypeString,
				Computed:   true,
				Deprecated: "Attribute 'opposite_interface_owner_id' has been deprecated from version 1.11.0. Use resource alicloud_router_interface_connection's 'opposite_interface_owner_id' instead.",
			},
		},
	}
}

func resourceAlicloudRouterInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}
	request, err := buildAlicloudRouterInterfaceCreateArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
		return vpcClient.CreateRouterInterface(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_router_interface", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*vpc.CreateRouterInterfaceResponse)
	d.SetId(response.RouterInterfaceId)

	if err := vpcService.WaitForRouterInterface(d.Id(), client.RegionId, Idle, 300); err != nil {
		return WrapError(err)
	}

	return resourceAlicloudRouterInterfaceUpdate(d, meta)
}

func resourceAlicloudRouterInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	d.Partial(true)

	request, attributeUpdate, err := buildAlicloudRouterInterfaceModifyAttrArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	request.RegionId = client.RegionId
	if attributeUpdate {
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	if d.HasChange("specification") && !d.IsNewResource() {
		d.SetPartial("specification")
		request := vpc.CreateModifyRouterInterfaceSpecRequest()
		request.RegionId = string(client.Region)
		request.RouterInterfaceId = d.Id()
		request.Spec = d.Get("specification").(string)
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.ModifyRouterInterfaceSpec(request)
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	}

	d.Partial(false)
	return resourceAlicloudRouterInterfaceRead(d, meta)
}

func resourceAlicloudRouterInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	object, err := vpcService.DescribeRouterInterface(d.Id(), client.RegionId)
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
	}

	d.Set("role", object.Role)
	d.Set("specification", object.Spec)
	d.Set("name", object.Name)
	d.Set("router_id", object.RouterId)
	d.Set("router_type", object.RouterType)
	d.Set("description", object.Description)
	d.Set("access_point_id", object.AccessPointId)
	d.Set("opposite_region", object.OppositeRegionId)
	d.Set("opposite_router_type", object.OppositeRouterType)
	d.Set("opposite_router_id", object.OppositeRouterId)
	d.Set("opposite_interface_id", object.OppositeInterfaceId)
	d.Set("opposite_interface_owner_id", object.OppositeInterfaceOwnerId)
	d.Set("health_check_source_ip", object.HealthCheckSourceIp)
	d.Set("health_check_target_ip", object.HealthCheckTargetIp)
	if object.ChargeType == "Prepaid" {
		d.Set("instance_charge_type", PrePaid)
		//period, err := computePeriodByUnit(object.CreationTime, object.EndTime, d.Get("period").(int), "Month")
		//if err != nil {
		//	return WrapError(err)
		//}
		//d.Set("period", period)
	} else {
		d.Set("instance_charge_type", PostPaid)
	}
	return nil

}

func resourceAlicloudRouterInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	if object, err := vpcService.DescribeRouterInterface(d.Id(), client.RegionId); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	} else if object.Status == string(Active) {
		if err := vpcService.DeactivateRouterInterface(d.Id()); err != nil {
			return WrapError(err)
		}
	}

	request := vpc.CreateDeleteRouterInterfaceRequest()
	request.RegionId = string(client.Region)
	request.RouterInterfaceId = d.Id()
	request.ClientToken = buildClientToken(request.GetActionName())
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		args := *request
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DeleteRouterInterface(&args)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectStatus", "DependencyViolation.RouterInterfaceReferedByRouteEntry"}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidInstanceId.NotFound"}) {
			return nil
		}
		WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	return WrapError(vpcService.WaitForRouterInterface(d.Id(), client.RegionId, Deleted, DefaultTimeoutMedium))
}

func buildAlicloudRouterInterfaceCreateArgs(d *schema.ResourceData, meta interface{}) (*vpc.CreateRouterInterfaceRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	oppositeRegion := d.Get("opposite_region").(string)
	if err := ecsService.JudgeRegionValidation("opposite_region", oppositeRegion); err != nil {
		return nil, WrapError(err)
	}

	request := vpc.CreateCreateRouterInterfaceRequest()
	request.RegionId = client.RegionId
	request.RouterType = d.Get("router_type").(string)
	request.RouterId = d.Get("router_id").(string)
	request.Role = d.Get("role").(string)
	request.Spec = d.Get("specification").(string)
	request.InstanceChargeType = d.Get("instance_charge_type").(string)
	if request.InstanceChargeType == string(PrePaid) {
		period := d.Get("period").(int)
		request.Period = requests.NewInteger(period)
		request.PricingCycle = string(Month)
		if period > 9 {
			request.Period = requests.NewInteger(period / 12)
			request.PricingCycle = string(Year)
		}
		request.AutoPay = requests.NewBoolean(true)
	}
	request.OppositeRegionId = oppositeRegion
	// Accepting side router interface spec only be Negative and router type only be VRouter.
	if request.Role == string(AcceptingSide) {
		request.Spec = string(Negative)
		request.RouterType = string(VRouter)
	} else {
		if request.Spec == "" {
			return request, WrapError(Error("'specification': required field is not set when role is %s.", InitiatingSide))
		}
	}

	// Get VBR access point
	if request.RouterType == string(VBR) {
		describeVirtualBorderRoutersRequest := vpc.CreateDescribeVirtualBorderRoutersRequest()
		values := []string{request.RouterId}
		filters := []vpc.DescribeVirtualBorderRoutersFilter{{
			Key:   "VbrId",
			Value: &values,
		}}
		describeVirtualBorderRoutersRequest.Filter = &filters
		raw, err := client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
			return vpcClient.DescribeVirtualBorderRouters(describeVirtualBorderRoutersRequest)
		})
		if err != nil {
			return request, WrapErrorf(err, DefaultErrorMsg, "alicloud_router_interface", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeVirtualBorderRoutersResponse)
		if response.TotalCount > 0 {
			request.AccessPointId = response.VirtualBorderRouterSet.VirtualBorderRouterType[0].AccessPointId
		}
	}
	request.ClientToken = buildClientToken(request.GetActionName())
	return request, nil
}

func buildAlicloudRouterInterfaceModifyAttrArgs(d *schema.ResourceData, meta interface{}) (*vpc.ModifyRouterInterfaceAttributeRequest, bool, error) {

	sourceIp, sourceOk := d.GetOk("health_check_source_ip")
	targetIp, targetOk := d.GetOk("health_check_target_ip")
	if sourceOk && !targetOk || !sourceOk && targetOk {
		return nil, false, WrapError(Error("The 'health_check_source_ip' and 'health_check_target_ip' should be specified or not at one time."))
	}

	request := vpc.CreateModifyRouterInterfaceAttributeRequest()
	request.RouterInterfaceId = d.Id()

	attributeUpdate := false

	if d.HasChange("health_check_source_ip") {
		d.SetPartial("health_check_source_ip")
		request.HealthCheckSourceIp = sourceIp.(string)
		request.HealthCheckTargetIp = targetIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("health_check_target_ip") {
		d.SetPartial("health_check_target_ip")
		request.HealthCheckTargetIp = targetIp.(string)
		request.HealthCheckSourceIp = sourceIp.(string)
		attributeUpdate = true
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		request.Name = d.Get("name").(string)
		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		request.Description = d.Get("description").(string)
		attributeUpdate = true
	}

	return request, attributeUpdate, nil
}
