package alicloud

import (
	"fmt"
	"math"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"reflect"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingGroupCreate,
		Read:   resourceAliyunEssScalingGroupRead,
		Update: resourceAliyunEssScalingGroupUpdate,
		Delete: resourceAliyunEssScalingGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"desired_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"scaling_group_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cooldown": {
				Type:         schema.TypeInt,
				Default:      300,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"vswitch_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Field 'vswitch_id' has been deprecated from provider version 1.7.1, and new field 'vswitch_ids' can replace it.",
			},
			"vswitch_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
			},
			"removal_policies": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
				MaxItems: 2,
				MinItems: 1,
			},
			"db_instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MinItems: 0,
			},
			"multi_az_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PRIORITY",
				ValidateFunc: validation.StringInSlice([]string{"PRIORITY", "BALANCE", "COST_OPTIMIZED"}, false),
				ForceNew:     true,
			},
			"on_demand_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			"on_demand_percentage_above_base_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"spot_instance_pools": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 10),
			},
			"spot_instance_remedy": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"group_deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"launch_template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssScalingGroupCreate(d *schema.ResourceData, meta interface{}) error {

	request, err := buildAlicloudEssScalingGroupArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.CreateScalingGroup(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{Throttling, "IncorrectLoadBalancerHealthCheck", "IncorrectLoadBalancerStatus"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.CreateScalingGroupResponse)
		d.SetId(response.ScalingGroupId)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalinggroup", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	if err := essService.WaitForEssScalingGroup(d.Id(), Inactive, DefaultTimeout); err != nil {
		return WrapError(err)
	}

	return resourceAliyunEssScalingGroupUpdate(d, meta)
}

func resourceAliyunEssScalingGroupRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("min_size", object.MinSize)
	d.Set("max_size", object.MaxSize)
	d.Set("desired_capacity", object.DesiredCapacity)
	d.Set("scaling_group_name", object.ScalingGroupName)
	d.Set("default_cooldown", object.DefaultCooldown)
	d.Set("multi_az_policy", object.MultiAZPolicy)
	d.Set("on_demand_base_capacity", object.OnDemandBaseCapacity)
	d.Set("on_demand_percentage_above_base_capacity", object.OnDemandPercentageAboveBaseCapacity)
	d.Set("spot_instance_pools", object.SpotInstancePools)
	d.Set("spot_instance_remedy", object.SpotInstanceRemedy)
	d.Set("group_deletion_protection", object.GroupDeletionProtection)
	var polices []string
	if len(object.RemovalPolicies.RemovalPolicy) > 0 {
		for _, v := range object.RemovalPolicies.RemovalPolicy {
			polices = append(polices, v)
		}
	}
	d.Set("removal_policies", polices)
	var dbIds []string
	if len(object.DBInstanceIds.DBInstanceId) > 0 {
		for _, v := range object.DBInstanceIds.DBInstanceId {
			dbIds = append(dbIds, v)
		}
	}
	d.Set("db_instance_ids", dbIds)

	var slbIds []string
	if len(object.LoadBalancerIds.LoadBalancerId) > 0 {
		for _, v := range object.LoadBalancerIds.LoadBalancerId {
			slbIds = append(slbIds, v)
		}
	}
	d.Set("loadbalancer_ids", slbIds)

	var vswitchIds []string
	if len(object.VSwitchIds.VSwitchId) > 0 {
		for _, v := range object.VSwitchIds.VSwitchId {
			vswitchIds = append(vswitchIds, v)
		}
	}
	d.Set("vswitch_ids", vswitchIds)
	d.Set("launch_template_id", object.LaunchTemplateId)

	return nil
}

func resourceAliyunEssScalingGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateModifyScalingGroupRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Id()

	d.Partial(true)
	if d.HasChange("scaling_group_name") {
		request.ScalingGroupName = d.Get("scaling_group_name").(string)
	}

	if d.HasChange("min_size") {
		request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	}

	if d.HasChange("max_size") {
		request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	}
	if d.HasChange("desired_capacity") {
		request.DesiredCapacity = requests.NewInteger(d.Get("desired_capacity").(int))
	}
	if d.HasChange("default_cooldown") {
		request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))
	}

	if d.HasChange("vswitch_ids") {
		vSwitchIds := expandStringList(d.Get("vswitch_ids").(*schema.Set).List())
		request.VSwitchIds = &vSwitchIds
	}

	if d.HasChange("removal_policies") {
		policyies := expandStringList(d.Get("removal_policies").([]interface{}))
		s := reflect.ValueOf(request).Elem()
		for i, p := range policyies {
			s.FieldByName(fmt.Sprintf("RemovalPolicy%d", i+1)).Set(reflect.ValueOf(p))
		}
	}

	if d.HasChange("on_demand_base_capacity") {
		request.OnDemandBaseCapacity = requests.NewInteger(d.Get("on_demand_base_capacity").(int))
	}

	if d.HasChange("on_demand_percentage_above_base_capacity") {
		request.OnDemandPercentageAboveBaseCapacity = requests.NewInteger(d.Get("on_demand_percentage_above_base_capacity").(int))
	}

	if d.HasChange("spot_instance_pools") {
		request.SpotInstancePools = requests.NewInteger(d.Get("spot_instance_pools").(int))
	}

	if d.HasChange("spot_instance_remedy") {
		request.SpotInstanceRemedy = requests.NewBoolean(d.Get("spot_instance_remedy").(bool))
	}

	if d.HasChange("group_deletion_protection") {
		request.GroupDeletionProtection = requests.NewBoolean(d.Get("group_deletion_protection").(bool))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingGroup(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	d.SetPartial("scaling_group_name")
	d.SetPartial("min_size")
	d.SetPartial("max_size")
	d.SetPartial("desired_capacity")
	d.SetPartial("default_cooldown")
	d.SetPartial("vswitch_ids")
	d.SetPartial("removal_policies")
	d.SetPartial("on_demand_base_capacity")
	d.SetPartial("on_demand_percentage_above_base_capacity")
	d.SetPartial("spot_instance_pools")
	d.SetPartial("spot_instance_remedy")
	d.SetPartial("group_deletion_protection")
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if d.HasChange("loadbalancer_ids") {
		oldLoadbalancers, newLoadbalancers := d.GetChange("loadbalancer_ids")
		err = attachOrDetachLoadbalancers(d, client, oldLoadbalancers.(*schema.Set), newLoadbalancers.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("loadbalancer_ids")
	}

	if d.HasChange("db_instance_ids") {
		oldDbInstanceIds, newDbInstanceIds := d.GetChange("db_instance_ids")
		err = attachOrDetachDbInstances(d, client, oldDbInstanceIds.(*schema.Set), newDbInstanceIds.(*schema.Set))
		if err != nil {
			return WrapError(err)
		}
		d.SetPartial("db_instance_ids")
	}
	d.Partial(false)
	return resourceAliyunEssScalingGroupRead(d, meta)
}

func resourceAliyunEssScalingGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	request := ess.CreateDeleteScalingGroupRequest()
	request.RegionId = client.RegionId
	request.ScalingGroupId = d.Id()
	request.ForceDelete = requests.NewBoolean(true)

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingGroup(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(essService.WaitForEssScalingGroup(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudEssScalingGroupArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupRequest, error) {
	client := meta.(*connectivity.AliyunClient)
	slbService := SlbService{client}
	request := ess.CreateCreateScalingGroupRequest()
	request.RegionId = client.RegionId
	request.MinSize = requests.NewInteger(d.Get("min_size").(int))
	request.MaxSize = requests.NewInteger(d.Get("max_size").(int))
	request.DefaultCooldown = requests.NewInteger(d.Get("default_cooldown").(int))

	if v, ok := d.GetOk("scaling_group_name"); ok && v.(string) != "" {
		request.ScalingGroupName = v.(string)
	}

	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("vswitch_ids"); ok {
		ids := expandStringList(v.(*schema.Set).List())
		request.VSwitchIds = &ids
	}

	if dbs, ok := d.GetOk("db_instance_ids"); ok {
		request.DBInstanceIds = convertListToJsonString(dbs.(*schema.Set).List())
	}

	if lbs, ok := d.GetOk("loadbalancer_ids"); ok {
		for _, lb := range lbs.(*schema.Set).List() {
			if err := slbService.WaitForSlb(lb.(string), Active, DefaultTimeout); err != nil {
				return nil, WrapError(err)
			}
		}
		request.LoadBalancerIds = convertListToJsonString(lbs.(*schema.Set).List())
	}

	if v, ok := d.GetOk("multi_az_policy"); ok && v.(string) != "" {
		request.MultiAZPolicy = v.(string)
	}

	if v, ok := d.GetOk("on_demand_base_capacity"); ok {
		request.OnDemandBaseCapacity = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("on_demand_percentage_above_base_capacity"); ok {
		request.OnDemandPercentageAboveBaseCapacity = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("spot_instance_pools"); ok {
		request.SpotInstancePools = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("spot_instance_remedy"); ok {
		request.SpotInstanceRemedy = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("group_deletion_protection"); ok {
		request.GroupDeletionProtection = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("launch_template_id"); ok {
		request.LaunchTemplateId = v.(string)
	}

	return request, nil
}

func attachOrDetachLoadbalancers(d *schema.ResourceData, client *connectivity.AliyunClient, oldLoadbalancerSet *schema.Set, newLoadbalancerSet *schema.Set) error {
	detachLoadbalancerSet := oldLoadbalancerSet.Difference(newLoadbalancerSet)
	attachLoadbalancerSet := newLoadbalancerSet.Difference(oldLoadbalancerSet)
	// attach
	if attachLoadbalancerSet.Len() > 0 {
		var subLists = partition(attachLoadbalancerSet, int(AttachDetachLoadbalancersBatchsize))
		for _, subList := range subLists {
			attachLoadbalancersRequest := ess.CreateAttachLoadBalancersRequest()
			attachLoadbalancersRequest.RegionId = client.RegionId
			attachLoadbalancersRequest.ScalingGroupId = d.Id()
			attachLoadbalancersRequest.ForceAttach = requests.NewBoolean(true)
			attachLoadbalancersRequest.LoadBalancer = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.AttachLoadBalancers(attachLoadbalancersRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), attachLoadbalancersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(attachLoadbalancersRequest.GetActionName(), raw, attachLoadbalancersRequest.RpcRequest, attachLoadbalancersRequest)
		}
	}
	// detach
	if detachLoadbalancerSet.Len() > 0 {
		var subLists = partition(detachLoadbalancerSet, int(AttachDetachLoadbalancersBatchsize))
		for _, subList := range subLists {
			detachLoadbalancersRequest := ess.CreateDetachLoadBalancersRequest()
			detachLoadbalancersRequest.RegionId = client.RegionId
			detachLoadbalancersRequest.ScalingGroupId = d.Id()
			detachLoadbalancersRequest.ForceDetach = requests.NewBoolean(false)
			detachLoadbalancersRequest.LoadBalancer = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DetachLoadBalancers(detachLoadbalancersRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), detachLoadbalancersRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(detachLoadbalancersRequest.GetActionName(), raw, detachLoadbalancersRequest.RpcRequest, detachLoadbalancersRequest)
		}
	}
	return nil
}

func attachOrDetachDbInstances(d *schema.ResourceData, client *connectivity.AliyunClient, oldDbInstanceIdSet *schema.Set, newDbInstanceIdSet *schema.Set) error {
	detachDbInstanceSet := oldDbInstanceIdSet.Difference(newDbInstanceIdSet)
	attachDbInstanceSet := newDbInstanceIdSet.Difference(oldDbInstanceIdSet)
	// attach
	if attachDbInstanceSet.Len() > 0 {
		var subLists = partition(attachDbInstanceSet, int(AttachDetachDbinstancesBatchsize))
		for _, subList := range subLists {
			attachDbInstancesRequest := ess.CreateAttachDBInstancesRequest()
			attachDbInstancesRequest.RegionId = client.RegionId
			attachDbInstancesRequest.ScalingGroupId = d.Id()
			attachDbInstancesRequest.ForceAttach = requests.NewBoolean(true)
			attachDbInstancesRequest.DBInstance = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.AttachDBInstances(attachDbInstancesRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), attachDbInstancesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(attachDbInstancesRequest.GetActionName(), raw, attachDbInstancesRequest.RpcRequest, attachDbInstancesRequest)
		}
	}
	// detach
	if detachDbInstanceSet.Len() > 0 {
		var subLists = partition(detachDbInstanceSet, int(AttachDetachDbinstancesBatchsize))
		for _, subList := range subLists {
			detachDbInstancesRequest := ess.CreateDetachDBInstancesRequest()
			detachDbInstancesRequest.RegionId = client.RegionId
			detachDbInstancesRequest.ScalingGroupId = d.Id()
			detachDbInstancesRequest.ForceDetach = requests.NewBoolean(true)
			detachDbInstancesRequest.DBInstance = &subList
			raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
				return essClient.DetachDBInstances(detachDbInstancesRequest)
			})
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), detachDbInstancesRequest.GetActionName(), AlibabaCloudSdkGoERROR)
			}
			addDebug(detachDbInstancesRequest.GetActionName(), raw, detachDbInstancesRequest.RpcRequest, detachDbInstancesRequest)
		}
	}
	return nil
}

func partition(instanceIds *schema.Set, batchSize int) [][]string {
	var res [][]string
	size := instanceIds.Len()
	batchCount := int(math.Ceil(float64(size) / float64(batchSize)))
	idList := expandStringList(instanceIds.List())
	for i := 1; i <= batchCount; i++ {
		fromIndex := batchSize * (i - 1)
		toIndex := int(math.Min(float64(batchSize*i), float64(size)))
		subList := idList[fromIndex:toIndex]
		res = append(res, subList)
	}
	return res
}
