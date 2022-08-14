// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBackupPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyRead,

		Schema: map[string]*schema.Schema{

			"identifier": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The backup policy identifier.",
			},

			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The unique user-defined name for this backup policy.",
			},

			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this backup policy.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy.",
			},
			"last_job_completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the most recent job for this backup policy completed.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the backup policy.",
			},
			"match_resource_types": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"match_user_tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"plans": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The plans for the backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this backup policy plan.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this backup policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func dataSourceIBMIsBackupPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	var backupPolicy *vpcv1.BackupPolicy

	if v, ok := d.GetOk("identifier"); ok {

		id := v.(string)
		getBackupPolicyOptions := &vpcv1.GetBackupPolicyOptions{}
		getBackupPolicyOptions.SetID(id)
		backupPolicyInfo, response, err := sess.GetBackupPolicyWithContext(context, getBackupPolicyOptions)
		if err != nil {
			log.Printf("[DEBUG] GetBackupPolicyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] GetBackupPolicyWithContext failed %s\n%s", err, response))
		}
		backupPolicy = backupPolicyInfo

	} else if v, ok := d.GetOk("name"); ok {

		name := v.(string)
		start := ""
		allrecs := []vpcv1.BackupPolicy{}
		for {
			listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{}
			if start != "" {
				listBackupPoliciesOptions.Start = &start
			}
			backupPolicyCollection, response, err := sess.ListBackupPoliciesWithContext(context, listBackupPoliciesOptions)
			if err != nil {
				log.Printf("[DEBUG] ListBackupPoliciesWithContext failed %s\n%s", err, response)
				return diag.FromErr(fmt.Errorf("[ERROR] ListBackupPoliciesWithContext failed %s\n%s", err, response))
			}
			if backupPolicyCollection != nil && *backupPolicyCollection.TotalCount == int64(0) {
				break
			}
			start = flex.GetNext(backupPolicyCollection.Next)
			allrecs = append(allrecs, backupPolicyCollection.BackupPolicies...)
			if start == "" {
				break
			}
		}
		for _, backupPolicyInfo := range allrecs {
			if *backupPolicyInfo.Name == name {
				backupPolicy = &backupPolicyInfo
				break
			}
		}
		if backupPolicy == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] No backup policy found with name (%s)", name))
		}
	}

	d.SetId(*backupPolicy.ID)

	if err = d.Set("created_at", backupPolicy.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", backupPolicy.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("href", backupPolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if backupPolicy.LastJobCompletedAt != nil {
		if err = d.Set("last_job_completed_at", backupPolicy.LastJobCompletedAt.String()); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting last_job_completed_at: %s", err))
		}
	}
	if err = d.Set("lifecycle_state", backupPolicy.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", backupPolicy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}

	if backupPolicy.Plans != nil {
		err = d.Set("plans", dataSourceBackupPolicyFlattenPlans(backupPolicy.Plans))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting plans %s", err))
		}
	}

	matchResourceType := make([]string, 0)
	if backupPolicy.MatchResourceTypes != nil {
		for _, matchResourceTyp := range backupPolicy.MatchResourceTypes {
			matchResourceType = append(matchResourceType, matchResourceTyp)
		}
	}
	d.Set("match_resource_types", matchResourceType)

	matchUserTags := make([]string, 0)
	if backupPolicy.MatchUserTags != nil {
		for _, matchUserTag := range backupPolicy.MatchUserTags {
			matchUserTags = append(matchUserTags, matchUserTag)
		}
	}
	d.Set("match_user_tags", matchUserTags)

	if backupPolicy.ResourceGroup != nil {
		err = d.Set("resource_group", dataSourceBackupPolicyFlattenResourceGroup(*backupPolicy.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", backupPolicy.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}

	return nil
}

func dataSourceBackupPolicyFlattenPlans(result []vpcv1.BackupPolicyPlanReference) (plans []map[string]interface{}) {
	for _, plansItem := range result {
		plans = append(plans, dataSourceBackupPolicyPlansToMap(plansItem))
	}

	return plans
}

func dataSourceBackupPolicyPlansToMap(plansItem vpcv1.BackupPolicyPlanReference) (plansMap map[string]interface{}) {
	plansMap = map[string]interface{}{}

	if plansItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyPlansDeletedToMap(*plansItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		plansMap["deleted"] = deletedList
	}
	if plansItem.Href != nil {
		plansMap["href"] = plansItem.Href
	}
	if plansItem.ID != nil {
		plansMap["id"] = plansItem.ID
	}
	if plansItem.Name != nil {
		plansMap["name"] = plansItem.Name
	}
	if plansItem.ResourceType != nil {
		plansMap["resource_type"] = plansItem.ResourceType
	}

	return plansMap
}

func dataSourceBackupPolicyPlansDeletedToMap(deletedItem vpcv1.BackupPolicyPlanReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}
