// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMISInstanceGroupManagerAction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISInstanceGroupManagerActionRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group manager action name",
			},

			"action_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance group manager action ID",
			},

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"instance_group_manager": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance group manager ID of type scheduled",
			},

			"run_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the scheduled action will run.",
			},

			"membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of members the instance group should have at the scheduled time.",
			},

			"max_membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of members in a managed instance group",
			},

			"min_membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum number of members in a managed instance group",
			},

			"target_manager": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this instance group manager of type autoscale.",
			},

			"target_manager_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance group manager name of type autoscale.",
			},

			"cron_spec": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The cron specification for a recurring scheduled action. Actions can be applied a maximum of one time within a 5 min period.",
			},

			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the instance group action- `active`: Action is ready to be run- `completed`: Action was completed successfully- `failed`: Action could not be completed successfully- `incompatible`: Action parameters are not compatible with the group or manager- `omitted`: Action was not applied because this action's manager was disabled.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance group manager action was modified.",
			},
			"action_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of action for the instance group.",
			},

			"last_applied_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the scheduled action was last applied. If empty the action has never been applied.",
			},
			"next_run_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time the scheduled action will next run. If empty the system is currently calculating the next run time.",
			},
			"auto_delete": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_delete_timeout": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance group manager action was modified.",
			},
		},
	}
}

func dataSourceIBMISInstanceGroupManagerActionRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	instanceGroupManagerID := d.Get("instance_group_manager").(string)
	instanceGroupID := d.Get("instance_group").(string)
	actionName := d.Get("name").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.InstanceGroupManagerActionIntf{}

	for {
		listInstanceGroupManagerActionsOptions := vpcv1.ListInstanceGroupManagerActionsOptions{
			InstanceGroupID:        &instanceGroupID,
			InstanceGroupManagerID: &instanceGroupManagerID,
		}

		instanceGroupManagerActionsCollection, response, err := sess.ListInstanceGroupManagerActions(&listInstanceGroupManagerActionsOptions)
		if err != nil {
			return fmt.Errorf("error Getting InstanceGroup Manager Actions %s\n%s", err, response)
		}
		if instanceGroupManagerActionsCollection != nil && *instanceGroupManagerActionsCollection.TotalCount == int64(0) {
			break
		}
		start = GetNext(instanceGroupManagerActionsCollection.Next)
		allrecs = append(allrecs, instanceGroupManagerActionsCollection.Actions...)
		if start == "" {
			break
		}
	}

	for _, data := range allrecs {
		instanceGroupManagerAction := data.(*vpcv1.InstanceGroupManagerAction)
		if actionName == *instanceGroupManagerAction.Name {
			d.SetId(fmt.Sprintf("%s/%s/%s", instanceGroupID, instanceGroupManagerID, *instanceGroupManagerAction.ID))

			if err = d.Set("auto_delete", *instanceGroupManagerAction.AutoDelete); err != nil {
				return fmt.Errorf("error setting auto_delete: %s", err)
			}

			if err = d.Set("auto_delete_timeout", intValue(instanceGroupManagerAction.AutoDeleteTimeout)); err != nil {
				return fmt.Errorf("error setting auto_delete_timeout: %s", err)
			}
			if err = d.Set("created_at", instanceGroupManagerAction.CreatedAt.String()); err != nil {
				return fmt.Errorf("error setting created_at: %s", err)
			}

			if err = d.Set("action_id", *instanceGroupManagerAction.ID); err != nil {
				return fmt.Errorf("error setting instance_group_manager_action : %s", err)
			}

			if err = d.Set("resource_type", *instanceGroupManagerAction.ResourceType); err != nil {
				return fmt.Errorf("error setting resource_type: %s", err)
			}
			if err = d.Set("status", *instanceGroupManagerAction.Status); err != nil {
				return fmt.Errorf("error setting status: %s", err)
			}
			if err = d.Set("updated_at", instanceGroupManagerAction.UpdatedAt.String()); err != nil {
				return fmt.Errorf("error setting updated_at: %s", err)
			}
			if err = d.Set("action_type", *instanceGroupManagerAction.ActionType); err != nil {
				return fmt.Errorf("error setting action_type: %s", err)
			}

			if instanceGroupManagerAction.CronSpec != nil {
				if err = d.Set("cron_spec", *instanceGroupManagerAction.CronSpec); err != nil {
					return fmt.Errorf("error setting cron_spec: %s", err)
				}
			}

			if instanceGroupManagerAction.LastAppliedAt != nil {
				if err = d.Set("last_applied_at", instanceGroupManagerAction.LastAppliedAt.String()); err != nil {
					return fmt.Errorf("error setting last_applied_at: %s", err)
				}
			}
			if instanceGroupManagerAction.NextRunAt != nil {
				if err = d.Set("next_run_at", instanceGroupManagerAction.NextRunAt.String()); err != nil {
					return fmt.Errorf("error setting next_run_at: %s", err)
				}
			}

			instanceGroupManagerScheduledActionGroupGroup := instanceGroupManagerAction.Group
			if instanceGroupManagerScheduledActionGroupGroup != nil && instanceGroupManagerScheduledActionGroupGroup.MembershipCount != nil {
				d.Set("membership_count", intValue(instanceGroupManagerScheduledActionGroupGroup.MembershipCount))
			}
			instanceGroupManagerScheduledActionManagerManagerInt := instanceGroupManagerAction.Manager
			if instanceGroupManagerScheduledActionManagerManagerInt != nil {
				instanceGroupManagerScheduledActionManagerManager := instanceGroupManagerScheduledActionManagerManagerInt.(*vpcv1.InstanceGroupManagerScheduledActionManager)
				if instanceGroupManagerScheduledActionManagerManager != nil && instanceGroupManagerScheduledActionManagerManager.ID != nil {

					if instanceGroupManagerScheduledActionManagerManager.MaxMembershipCount != nil {
						d.Set("max_membership_count", intValue(instanceGroupManagerScheduledActionManagerManager.MaxMembershipCount))
					}
					d.Set("min_membership_count", intValue(instanceGroupManagerScheduledActionManagerManager.MinMembershipCount))
					d.Set("target_manager_name", *instanceGroupManagerScheduledActionManagerManager.Name)
					d.Set("target_manager", *instanceGroupManagerScheduledActionManagerManager.ID)
				}
			}
			return nil
		}
	}
	return fmt.Errorf("instance group manager action %s not found", actionName)
}
