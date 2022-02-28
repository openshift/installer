package alicloud

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAttachmentCreate,
		Read:   resourceAliyunEssAttachmentRead,
		Update: resourceAliyunEssAttachmentUpdate,
		Delete: resourceAliyunEssAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instance_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"force": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunEssAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("scaling_group_id").(string))
	return resourceAliyunEssAttachmentUpdate(d, meta)
}

func resourceAliyunEssAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	d.Partial(true)

	if d.HasChange("instance_ids") {
		object, err := essService.DescribeEssScalingGroup(d.Id())
		if err != nil {
			return WrapError(err)
		}
		if object.LifecycleState == string(Inactive) {
			return WrapError(Error("Scaling group current status is %s, please active it before attaching or removing ECS instances.", object.LifecycleState))
		} else {
			if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
				return WrapError(err)
			}
		}
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := convertArrayInterfaceToArrayString(ns.Difference(os).List())

		if len(add) > 0 {
			request := ess.CreateAttachInstancesRequest()
			request.RegionId = client.RegionId
			request.ScalingGroupId = d.Id()
			request.InstanceId = &add

			err := resource.Retry(5*time.Minute, func() *resource.RetryError {
				raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
					return essClient.AttachInstances(request)
				})
				if err != nil {
					if IsExpectedErrors(err, []string{"IncorrectCapacity.MaxSize"}) {
						instances, err := essService.DescribeEssAttachment(d.Id(), make([]string, 0))
						if !NotFoundError(err) {
							return resource.NonRetryableError(err)
						}
						var autoAdded, attached []string
						if len(instances) > 0 {
							for _, inst := range instances {
								if inst.CreationType == "Attached" {
									attached = append(attached, inst.InstanceId)
								} else {
									autoAdded = append(autoAdded, inst.InstanceId)
								}
							}
						}
						if len(add) > object.MaxSize {
							return resource.NonRetryableError(WrapError(Error("To attach %d instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size.", len(add), object.MaxSize)))
						}

						if len(autoAdded) > 0 {
							if d.Get("force").(bool) {
								if err := essService.EssRemoveInstances(d.Id(), autoAdded); err != nil {
									return resource.NonRetryableError(WrapError(err))
								}
								time.Sleep(5)
								return resource.RetryableError(WrapError(err))
							} else {
								return resource.NonRetryableError(WrapError(Error("To attach the instances, the total capacity will be greater than the scaling group max size %d."+
									"Please enlarge scaling group max size or set 'force' to true to remove autocreated instances: %#v.", object.MaxSize, autoAdded)))
							}
						}

						if len(attached) > 0 {
							return resource.NonRetryableError(WrapError(Error("To attach the instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size or remove already attached instances: %#v.", object.MaxSize, attached)))
						}
					}
					if IsExpectedErrors(err, []string{"ScalingActivityInProgress"}) {
						time.Sleep(5)
						return resource.RetryableError(WrapError(err))
					}
					return resource.NonRetryableError(WrapError(err))
				}
				addDebug(request.GetActionName(), raw, request.RpcRequest, request)
				return nil
			})
			if err != nil {
				return WrapError(err)
			}

			err = resource.Retry(3*time.Minute, func() *resource.RetryError {

				instances, err := essService.DescribeEssAttachment(d.Id(), add)
				if err != nil {
					return resource.NonRetryableError(WrapError(err))
				}
				if len(instances) < 0 {
					return resource.RetryableError(WrapError(Error("There are no ECS instances have been attached.")))
				}

				for _, inst := range instances {
					if inst.LifecycleState != string(InService) {
						return resource.RetryableError(WrapError(Error("There are still ECS instances are not %s.", string(InService))))
					}
				}
				return nil
			})
			if err != nil {
				return WrapError(err)
			}
		}
		if len(remove) > 0 {
			if err := essService.EssRemoveInstances(d.Id(), convertArrayInterfaceToArrayString(remove)); err != nil {
				return WrapError(err)
			}
		}

		d.SetPartial("instance_ids")
	}

	d.Partial(false)

	return resourceAliyunEssAttachmentRead(d, meta)
}

func resourceAliyunEssAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	object, err := essService.DescribeEssAttachment(d.Id(), make([]string, 0))

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	var instanceIds []string
	for _, inst := range object {
		instanceIds = append(instanceIds, inst.InstanceId)
	}

	d.Set("scaling_group_id", object[0].ScalingGroupId)
	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceAliyunEssAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}

	removed := convertArrayInterfaceToArrayString(d.Get("instance_ids").(*schema.Set).List())

	if len(removed) < 1 {
		return nil
	}
	object, err := essService.DescribeEssScalingGroup(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if err := essService.WaitForEssScalingGroup(object.ScalingGroupId, Active, DefaultTimeout); err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	if err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := ess.CreateRemoveInstancesRequest()
		request.RegionId = client.RegionId
		request.ScalingGroupId = d.Id()

		if len(removed) > 0 {
			request.InstanceId = &removed
		} else {
			return nil
		}
		raw, err := essService.client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.RemoveInstances(request)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{"IncorrectCapacity.MinSize"}) {
				instances, err := essService.DescribeEssAttachment(d.Id(), removed)
				if len(instances) > 0 {
					if object.MinSize == 0 {
						return resource.RetryableError(WrapError(err))
					}
					return resource.NonRetryableError(WrapError(Error("To remove %d instances, the total capacity will be lesser than the scaling group min size %d. "+
						"Please shorten scaling group min size and try again.", len(removed), object.MinSize)))
				}
			}
			if IsExpectedErrors(err, []string{"ScalingActivityInProgress", "IncorrectScalingGroupStatus"}) {
				time.Sleep(5)
				return resource.RetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
			}
			if IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			return resource.NonRetryableError(WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR))
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		time.Sleep(3 * time.Second)
		instances, err := essService.DescribeEssAttachment(d.Id(), removed)
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(WrapError(err))
		}
		if len(instances) > 0 {
			removed = make([]string, 0)
			for _, inst := range instances {
				removed = append(removed, inst.InstanceId)
			}
			return resource.RetryableError(WrapError(Error("There are still ECS instances in the scaling group.")))
		}

		return nil
	}); err != nil {
		return WrapError(err)
	}

	return WrapError(essService.WaitForEssAttachment(d.Id(), Deleted, DefaultTimeout))
}

func convertArrayInterfaceToArrayString(elm []interface{}) (arr []string) {
	if len(elm) < 1 {
		return
	}
	for _, e := range elm {
		arr = append(arr, e.(string))
	}
	return
}
