package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudAutoProvisioningGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudAutoProvisioningGroupCreate,
		Read:   resourceAlicloudAutoProvisioningGroupRead,
		Update: resourceAlicloudAutoProvisioningGroupUpdate,
		Delete: resourceAlicloudAutoProvisioningGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"total_target_capacity": {
				Type:     schema.TypeString,
				Required: true,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auto_provisioning_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auto_provisioning_group_type": {
				Type:         schema.TypeString,
				Default:      "maintain",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"request", "maintain"}, false),
			},
			"spot_allocation_strategy": {
				Type:         schema.TypeString,
				Default:      "lowest-price",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"lowest-price", "diversified"}, false),
			},
			"spot_instance_interruption_behavior": {
				Type:         schema.TypeString,
				Default:      "stop",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"stop", "terminate"}, false),
			},
			"spot_instance_pools_to_use_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"pay_as_you_go_allocation_strategy": {
				Type:         schema.TypeString,
				Default:      "lowest-price",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"prioritized", "lowest-price"}, false),
			},
			"excess_capacity_termination_policy": {
				Type:         schema.TypeString,
				Default:      "no-termination",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"no-termination", "termination"}, false),
			},
			"valid_from": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"valid_until": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"terminate_instances_with_expiration": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"terminate_instances": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},
			"max_spot_price": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"pay_as_you_go_target_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spot_target_capacity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_target_capacity_type": {
				Type:         schema.TypeString,
				Default:      "Spot",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Spot", "PayAsYouGo"}, false),
			},
			"launch_template_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"launch_template_config": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_price": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weighted_capacity": {
							Type:     schema.TypeString,
							Required: true,
						},
						"priority": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudAutoProvisioningGroupCreate(d *schema.ResourceData, meta interface{}) error {

	request, err := buildAlicloudAutoProvisioningGroupArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.CreateAutoProvisioningGroup(request)
		})
		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ecs.CreateAutoProvisioningGroupResponse)
		d.SetId(response.AutoProvisioningGroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_auto_provisioning_group", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	return resourceAlicloudAutoProvisioningGroupRead(d, meta)
}

func resourceAlicloudAutoProvisioningGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	object, err := ecsService.DescribeAutoProvisioningGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_auto_provisioning_group ecsService.DescribeAutoProvisioningGroup Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("auto_provisioning_group_name", object.AutoProvisioningGroupName)
	d.Set("launch_template_id", object.LaunchTemplateId)
	d.Set("auto_provisioning_group_type", object.AutoProvisioningGroupType)
	d.Set("auto_provisioning_group_id", object.AutoProvisioningGroupId)
	d.Set("create_time", object.CreationTime)
	d.Set("excess_capacity_termination_policy", object.ExcessCapacityTerminationPolicy)
	d.Set("terminate_instances_with_expiration", object.TerminateInstancesWithExpiration)
	d.Set("max_spot_price", object.MaxSpotPrice)
	d.Set("terminate_instances", object.TerminateInstances)
	d.Set("valid_from", object.ValidFrom)
	d.Set("valid_until", object.ValidUntil)
	d.Set("status", object.Status)
	d.Set("state", object.State)
	d.Set("pay_as_you_go_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.PayAsYouGoTargetCapacity))
	d.Set("pay_as_you_go_allocation_strategy", object.PayAsYouGoOptions.AllocationStrategy)
	d.Set("spot_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.SpotTargetCapacity))
	d.Set("spot_allocation_strategy", object.SpotOptions.AllocationStrategy)
	d.Set("spot_instance_interruption_behavior", object.SpotOptions.InstanceInterruptionBehavior)
	d.Set("spot_instance_pools_to_use_count", object.SpotOptions.InstancePoolsToUseCount)
	d.Set("total_target_capacity", fmt.Sprintf("%v", object.TargetCapacitySpecification.TotalTargetCapacity))
	d.Set("default_target_capacity_type", object.TargetCapacitySpecification.DefaultTargetCapacityType)
	d.Set("launch_template_version", object.LaunchTemplateVersion)
	launch_template_config := []map[string]interface{}{}
	para := map[string]interface{}{}
	for _, mappara := range object.LaunchTemplateConfigs.LaunchTemplateConfig {
		para["instance_type"] = mappara.InstanceType
		para["vswitch_id"] = mappara.VSwitchId
		para["weighted_capacity"] = fmt.Sprintf("%v", mappara.WeightedCapacity)
		para["max_price"] = fmt.Sprintf("%v", mappara.MaxPrice)
		para["priority"] = fmt.Sprintf("%v", mappara.Priority)
	}
	launch_template_config = append(launch_template_config, para)
	d.Set("launch_template_config", launch_template_config)

	return nil
}

func resourceAlicloudAutoProvisioningGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateModifyAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.AutoProvisioningGroupId = d.Id()

	if d.HasChange("excess_capacity_termination_policy") {
		if v, ok := d.GetOk("excess_capacity_termination_policy"); ok {
			request.ExcessCapacityTerminationPolicy = v.(string)
		}
	}

	if d.HasChange("default_target_capacity_type") {
		if v, ok := d.GetOk("default_target_capacity_type"); ok {
			request.DefaultTargetCapacityType = v.(string)
		}
	}

	if d.HasChange("terminate_instances_with_expiration") {
		request.TerminateInstancesWithExpiration = requests.NewBoolean(d.Get("terminate_instances_with_expiration").(bool))
	}

	if d.HasChange("max_spot_price") {
		request.MaxSpotPrice = requests.NewFloat(d.Get("max_spot_price").(float64))
	}

	if d.HasChange("total_target_capacity") {
		if v, ok := d.GetOk("total_target_capacity"); ok {
			request.TotalTargetCapacity = v.(string)
		}
	}

	if d.HasChange("pay_as_you_go_target_capacity") {
		if v, ok := d.GetOk("pay_as_you_go_target_capacity"); ok {
			request.PayAsYouGoTargetCapacity = v.(string)
		}
	}

	if d.HasChange("spot_target_capacity") {
		if v, ok := d.GetOk("spot_target_capacity"); ok {
			request.SpotTargetCapacity = v.(string)
		}
	}

	if d.HasChange("auto_provisioning_group_name") {
		if v, ok := d.GetOk("auto_provisioning_group_name"); ok {
			request.AutoProvisioningGroupName = v.(string)
		}
	}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ModifyAutoProvisioningGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return resourceAlicloudAutoProvisioningGroupRead(d, meta)
}

func resourceAlicloudAutoProvisioningGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}

	request := ecs.CreateDeleteAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.AutoProvisioningGroupId = d.Id()
	request.TerminateInstances = requests.NewBoolean(d.Get("terminate_instances").(bool))

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DeleteAutoProvisioningGroup(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidAutoProvisioningGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(ecsService.WaitForAutoProvisioningGroup(d.Id(), Deleted, DefaultTimeoutMedium))
}

func buildAlicloudAutoProvisioningGroupArgs(d *schema.ResourceData, meta interface{}) (*ecs.CreateAutoProvisioningGroupRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateCreateAutoProvisioningGroupRequest()
	request.RegionId = client.RegionId
	request.LaunchTemplateId = d.Get("launch_template_id").(string)
	request.TotalTargetCapacity = d.Get("total_target_capacity").(string)

	if v, ok := d.GetOk("auto_provisioning_group_name"); ok && v.(string) != "" {
		request.AutoProvisioningGroupName = v.(string)
	}
	if v, ok := d.GetOk("auto_provisioning_group_type"); ok && v.(string) != "" {
		request.AutoProvisioningGroupType = v.(string)
	}
	if v, ok := d.GetOk("spot_allocation_strategy"); ok && v.(string) != "" {
		request.SpotAllocationStrategy = v.(string)
	}
	if v, ok := d.GetOk("spot_instance_interruption_behavior"); ok && v.(string) != "" {
		request.SpotInstanceInterruptionBehavior = v.(string)
	}
	if v, ok := d.GetOk("spot_instancePools_to_use_count"); ok && v.(string) != "" {
		request.SpotInstancePoolsToUseCount = requests.NewInteger(v.(int))
	}
	if v, ok := d.GetOk("pay_as_you_go_allocation_strategy"); ok && v.(string) != "" {
		request.PayAsYouGoAllocationStrategy = v.(string)
	}
	if v, ok := d.GetOk("excess_capacity_termination_policy"); ok && v.(string) != "" {
		request.ExcessCapacityTerminationPolicy = v.(string)
	}
	if v, ok := d.GetOk("valid_from"); ok && v.(string) != "" {
		request.ValidFrom = v.(string)
	}
	if v, ok := d.GetOk("valid_until"); ok && v.(string) != "" {
		request.ValidUntil = v.(string)
	}
	if v, ok := d.GetOk("terminate_instances_with_expiration"); ok {
		request.TerminateInstancesWithExpiration = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("terminate_instances"); ok {
		request.TerminateInstances = requests.NewBoolean(v.(bool))
	}
	if v, ok := d.GetOk("max_spot_price"); ok && v.(float64) != 0.0 {
		request.MaxSpotPrice = requests.NewFloat(v.(float64))
	}
	if v, ok := d.GetOk("pay_as_you_go_target_capacity"); ok && v.(string) != "" {
		request.PayAsYouGoTargetCapacity = v.(string)
	}
	if v, ok := d.GetOk("spot_target_capacity"); ok && v.(string) != "" {
		request.SpotTargetCapacity = v.(string)
	}
	if v, ok := d.GetOk("default_target_capacity_type"); ok && v.(string) != "" {
		request.DefaultTargetCapacityType = v.(string)
	}
	if v, ok := d.GetOk("launch_template_version"); ok && v.(string) != "" {
		request.LaunchTemplateVersion = v.(string)
	}
	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}
	configs := d.Get("launch_template_config")
	confs := configs.([]interface{})
	createConfigs := make([]ecs.CreateAutoProvisioningGroupLaunchTemplateConfig, 0, len(confs))
	for _, c := range confs {
		cc := c.(map[string]interface{})
		conf := ecs.CreateAutoProvisioningGroupLaunchTemplateConfig{
			MaxPrice:         cc["max_price"].(string),
			VSwitchId:        cc["vswitch_id"].(string),
			WeightedCapacity: cc["weighted_capacity"].(string),
		}
		if v, ok := cc["instance_type"]; ok && v.(string) != "" {
			conf.InstanceType = cc["instance_type"].(string)
		}
		if v, ok := cc["priority"]; ok && v.(string) != "" {
			conf.Priority = cc["priority"].(string)
		}
		createConfigs = append(createConfigs, conf)
	}
	request.LaunchTemplateConfig = &createConfigs

	return request, nil
}
