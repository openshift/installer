// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBackupPolicyPlans() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyPlansRead,

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique user-defined name for this backup policy plan.",
			},
			"plans": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of backup policy plans.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the plan is active.",
						},
						"attach_user_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The user tags to attach to backups (snapshots) created by this plan.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"copy_user_tags": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether to copy the source's user tags to the created backups (snapshots).",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the backup policy plan was created.",
						},
						"cron_spec": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron specification for the backup schedule.",
						},
						"deletion_trigger": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delete_after": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum number of days to keep each backup after creation.",
									},
									"delete_over_count": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum number of recent backups to keep. If absent, there is no maximum.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this backup policy plan.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy plan.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this backup policy plan.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this backup policy plan.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsBackupPolicyPlansRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listBackupPolicyPlansOptions := &vpcv1.ListBackupPolicyPlansOptions{}

	listBackupPolicyPlansOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))

	backupPolicyPlanCollection, response, err := vpcClient.ListBackupPolicyPlansWithContext(context, listBackupPolicyPlansOptions)
	if err != nil {
		log.Printf("[DEBUG] ListBackupPolicyPlansWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] ListBackupPolicyPlansWithContext failed %s\n%s", err, response))
	}

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	var matchPlans []vpcv1.BackupPolicyPlan
	var name string
	var suppliedFilter bool

	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
		suppliedFilter = true
		for _, data := range backupPolicyPlanCollection.Plans {
			if *data.Name == name {
				matchPlans = append(matchPlans, data)
			}
		}
		backupPolicyPlanCollection.Plans = matchPlans
	}
	if suppliedFilter {
		if len(backupPolicyPlanCollection.Plans) == 0 {
			return diag.FromErr(fmt.Errorf("[ERROR] no plans found with name %s", name))
		}
		d.SetId(name)
	} else {
		d.SetId(dataSourceIBMIsBackupPolicyPlansID(d))
	}

	if backupPolicyPlanCollection.Plans != nil {
		err = d.Set("plans", dataSourceBackupPolicyPlanCollectionFlattenPlans(backupPolicyPlanCollection.Plans))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting plans %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsBackupPolicyPlansID returns a reasonable ID for the list.
func dataSourceIBMIsBackupPolicyPlansID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceBackupPolicyPlanCollectionFlattenPlans(result []vpcv1.BackupPolicyPlan) (plans []map[string]interface{}) {
	for _, plansItem := range result {
		plans = append(plans, dataSourceBackupPolicyPlanCollectionPlansToMap(plansItem))
	}

	return plans
}

func dataSourceBackupPolicyPlanCollectionPlansToMap(plansItem vpcv1.BackupPolicyPlan) (plansMap map[string]interface{}) {
	plansMap = map[string]interface{}{}

	if plansItem.Active != nil {
		plansMap["active"] = plansItem.Active
	}
	if plansItem.AttachUserTags != nil {
		plansMap["attach_user_tags"] = plansItem.AttachUserTags
	}
	if plansItem.CopyUserTags != nil {
		plansMap["copy_user_tags"] = plansItem.CopyUserTags
	}
	if plansItem.CreatedAt != nil {
		plansMap["created_at"] = plansItem.CreatedAt.String()
	}
	if plansItem.CronSpec != nil {
		plansMap["cron_spec"] = plansItem.CronSpec
	}
	if plansItem.DeletionTrigger != nil {
		deletionTriggerList := []map[string]interface{}{}
		deletionTriggerMap := dataSourceBackupPolicyPlanCollectionPlansDeletionTriggerToMap(*plansItem.DeletionTrigger)
		deletionTriggerList = append(deletionTriggerList, deletionTriggerMap)
		plansMap["deletion_trigger"] = deletionTriggerList
	}
	if plansItem.Href != nil {
		plansMap["href"] = plansItem.Href
	}
	if plansItem.ID != nil {
		plansMap["id"] = plansItem.ID
	}
	if plansItem.LifecycleState != nil {
		plansMap["lifecycle_state"] = plansItem.LifecycleState
	}
	if plansItem.Name != nil {
		plansMap["name"] = plansItem.Name
	}
	if plansItem.ResourceType != nil {
		plansMap["resource_type"] = plansItem.ResourceType
	}

	return plansMap
}

func dataSourceBackupPolicyPlanCollectionClonePolicyZonesToMap(zonesItem vpcv1.ZoneReference) (zonesMap map[string]interface{}) {
	zonesMap = map[string]interface{}{}

	if zonesItem.Href != nil {
		zonesMap["href"] = zonesItem.Href
	}
	if zonesItem.Name != nil {
		zonesMap["name"] = zonesItem.Name
	}

	return zonesMap
}

func dataSourceBackupPolicyPlanCollectionPlansDeletionTriggerToMap(deletionTriggerItem vpcv1.BackupPolicyPlanDeletionTrigger) (deletionTriggerMap map[string]interface{}) {
	deletionTriggerMap = map[string]interface{}{}

	if deletionTriggerItem.DeleteAfter != nil {
		deletionTriggerMap["delete_after"] = deletionTriggerItem.DeleteAfter
	}
	if deletionTriggerItem.DeleteOverCount != nil {
		deletionTriggerMap["delete_over_count"] = deletionTriggerItem.DeleteOverCount
	}

	return deletionTriggerMap
}
