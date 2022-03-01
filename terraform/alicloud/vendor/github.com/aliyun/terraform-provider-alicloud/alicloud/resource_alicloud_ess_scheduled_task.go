package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEssScheduledTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScheduledTaskCreate,
		Read:   resourceAliyunEssScheduledTaskRead,
		Update: resourceAliyunEssScheduledTaskUpdate,
		Delete: resourceAliyunEssScheduledTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scheduled_action": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"scaling_group_id"},
			},
			"launch_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"scaling_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"scheduled_action"},
			},
			"launch_expiration_time": {
				Type:         schema.TypeInt,
				Default:      600,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 21600),
			},
			"min_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"desired_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			//RecurrenceType, RecurrenceValue and RecurrenceEndTime must be specified.
			"recurrence_type": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Daily", "Weekly", "Monthly", "Cron"}, false),
			},
			"recurrence_value": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"recurrence_end_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"task_enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssScheduledTaskCreate(d *schema.ResourceData, meta interface{}) error {

	request := buildAlicloudEssScheduledTaskArgs(d)
	client := meta.(*connectivity.AliyunClient)
	request.RegionId = client.RegionId
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScheduledTask(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scheduled_task", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.CreateScheduledTaskResponse)
	d.SetId(response.ScheduledTaskId)

	return resourceAliyunEssScheduledTaskRead(d, meta)
}

func resourceAliyunEssScheduledTaskRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	object, err := essService.DescribeEssScheduledTask(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scheduled_action", object.ScheduledAction)
	d.Set("launch_time", object.LaunchTime)
	d.Set("scheduled_task_name", object.ScheduledTaskName)
	d.Set("description", object.Description)
	d.Set("launch_expiration_time", object.LaunchExpirationTime)
	d.Set("recurrence_type", object.RecurrenceType)
	d.Set("recurrence_value", object.RecurrenceValue)
	d.Set("recurrence_end_time", object.RecurrenceEndTime)
	d.Set("task_enabled", object.TaskEnabled)
	d.Set("min_value", object.MinValue)
	d.Set("max_value", object.MaxValue)
	d.Set("desired_capacity", object.DesiredCapacity)
	d.Set("scaling_group_id", object.ScalingGroupId)

	return nil
}

func resourceAliyunEssScheduledTaskUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	request := ess.CreateModifyScheduledTaskRequest()
	request.RegionId = client.RegionId
	request.ScheduledTaskId = d.Id()
	request.LaunchExpirationTime = requests.NewInteger(d.Get("launch_expiration_time").(int))

	if d.HasChange("scheduled_task_name") {
		request.ScheduledTaskName = d.Get("scheduled_task_name").(string)
	}

	if d.HasChange("scaling_group_id") {
		request.ScalingGroupId = d.Get("scaling_group_id").(string)
	}

	if d.HasChange("min_value") {
		request.MinValue = requests.NewInteger(d.Get("min_value").(int))
	}

	if d.HasChange("max_value") {
		request.MaxValue = requests.NewInteger(d.Get("max_value").(int))
	}

	if d.HasChange("desired_capacity") {
		request.DesiredCapacity = requests.NewInteger(d.Get("desired_capacity").(int))
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
	}

	if d.HasChange("scheduled_action") {
		request.ScheduledAction = d.Get("scheduled_action").(string)
	}

	if d.HasChange("launch_time") {
		request.LaunchTime = d.Get("launch_time").(string)
	}

	if d.HasChange("recurrence_type") || d.HasChange("recurrence_value") || d.HasChange("recurrence_end_time") {
		request.RecurrenceType = d.Get("recurrence_type").(string)
		request.RecurrenceValue = d.Get("recurrence_value").(string)
		request.RecurrenceEndTime = d.Get("recurrence_end_time").(string)
	}

	if d.HasChange("task_enabled") {
		request.TaskEnabled = requests.NewBoolean(d.Get("task_enabled").(bool))
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScheduledTask(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAliyunEssScheduledTaskRead(d, meta)
}

func resourceAliyunEssScheduledTaskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	request := ess.CreateDeleteScheduledTaskRequest()
	request.ScheduledTaskId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScheduledTask(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScheduledTaskId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return WrapError(essService.WaitForEssScheduledTask(d.Id(), Deleted, DefaultTimeout))
}

func buildAlicloudEssScheduledTaskArgs(d *schema.ResourceData) *ess.CreateScheduledTaskRequest {
	request := ess.CreateCreateScheduledTaskRequest()
	if v, ok := d.GetOk("scheduled_action"); ok && v.(string) != "" {
		request.ScheduledAction = v.(string)
	}

	if v, ok := d.GetOk("launch_time"); ok && v.(string) != "" {
		request.LaunchTime = v.(string)
	}

	if v, ok := d.GetOk("task_enabled"); ok {
		request.TaskEnabled = requests.NewBoolean(v.(bool))
	}

	if v, ok := d.GetOk("scheduled_task_name"); ok && v.(string) != "" {
		request.ScheduledTaskName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok && v.(string) != "" {
		request.Description = v.(string)
	}

	if v, ok := d.GetOk("recurrence_type"); ok && v.(string) != "" {
		request.RecurrenceType = v.(string)
	}

	if v, ok := d.GetOk("recurrence_value"); ok && v.(string) != "" {
		request.RecurrenceValue = v.(string)
	}

	if v, ok := d.GetOk("recurrence_end_time"); ok && v.(string) != "" {
		request.RecurrenceEndTime = v.(string)
	}

	if v, ok := d.GetOk("scaling_group_id"); ok && v.(string) != "" {
		request.ScalingGroupId = v.(string)
	}

	if v, ok := d.GetOk("launch_expiration_time"); ok && v.(int) != 0 {
		request.LaunchExpirationTime = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("min_value"); ok && v.(int) != 0 {
		request.MinValue = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("max_value"); ok && v.(int) != 0 {
		request.MaxValue = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("desired_capacity"); ok && v.(int) != 0 {
		request.DesiredCapacity = requests.NewInteger(v.(int))
	}

	return request
}
