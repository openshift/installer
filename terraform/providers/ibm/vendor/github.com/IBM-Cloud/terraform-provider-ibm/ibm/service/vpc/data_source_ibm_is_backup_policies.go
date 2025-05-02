// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBackupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPoliciesRead,

		Schema: map[string]*schema.Schema{
			"resource_group": {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources in the resource group with the specified identifier",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with the exact specified name",
				Optional:    true,
			},
			"tag": {
				Type:        schema.TypeString,
				Description: "Filters the collection to resources with the exact tag value",
				Optional:    true,
			},

			"backup_policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of backup policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the backup policy.",
						},
						"last_job_completed_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the most recent job for this backup policy completed.",
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this backup policy.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsBackupPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	matchBackupPolicies := []vpcv1.BackupPolicy{}

	var resourceGroup string
	if v, ok := d.GetOk("resource_group"); ok {
		resourceGroup = v.(string)
	}

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	var tag string
	if v, ok := d.GetOk("tag"); ok {
		tag = v.(string)
	}

	for {
		listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{}
		if start != "" {
			listBackupPoliciesOptions.Start = &start
		}
		if resourceGroup != "" {
			listBackupPoliciesOptions.SetResourceGroupID(resourceGroup)
		}
		if name != "" {
			listBackupPoliciesOptions.SetName(name)
		}
		if tag != "" {
			listBackupPoliciesOptions.SetTag(tag)
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
		matchBackupPolicies = append(matchBackupPolicies, backupPolicyCollection.BackupPolicies...)
		if start == "" {
			break
		}

	}

	d.SetId(dataSourceIBMIsBackupPoliciesID(d))

	if matchBackupPolicies != nil {
		err = d.Set("backup_policies", dataSourceBackupPolicyCollectionFlattenBackupPolicies(matchBackupPolicies))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting backup_policies %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsBackupPoliciesID returns a reasonable ID for the list.
func dataSourceIBMIsBackupPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceBackupPolicyCollectionFlattenBackupPolicies(result []vpcv1.BackupPolicy) (backupPolicies []map[string]interface{}) {
	for _, backupPoliciesItem := range result {
		backupPolicies = append(backupPolicies, dataSourceBackupPolicyCollectionBackupPoliciesToMap(backupPoliciesItem))
	}

	return backupPolicies
}

func dataSourceBackupPolicyCollectionBackupPoliciesToMap(backupPoliciesItem vpcv1.BackupPolicy) (backupPoliciesMap map[string]interface{}) {
	backupPoliciesMap = map[string]interface{}{}

	if backupPoliciesItem.CreatedAt != nil {
		backupPoliciesMap["created_at"] = backupPoliciesItem.CreatedAt.String()
	}
	if backupPoliciesItem.CRN != nil {
		backupPoliciesMap["crn"] = backupPoliciesItem.CRN
	}
	if backupPoliciesItem.Href != nil {
		backupPoliciesMap["href"] = backupPoliciesItem.Href
	}
	if backupPoliciesItem.ID != nil {
		backupPoliciesMap["id"] = backupPoliciesItem.ID
	}
	if backupPoliciesItem.LifecycleState != nil {
		backupPoliciesMap["lifecycle_state"] = backupPoliciesItem.LifecycleState
	}
	if backupPoliciesItem.LastJobCompletedAt != nil {
		backupPoliciesMap["last_job_completed_at"] = flex.DateTimeToString(backupPoliciesItem.LastJobCompletedAt)
	}
	if backupPoliciesItem.MatchResourceTypes != nil {
		backupPoliciesMap["match_resource_types"] = backupPoliciesItem.MatchResourceTypes
	}
	if backupPoliciesItem.MatchUserTags != nil {
		backupPoliciesMap["match_user_tags"] = backupPoliciesItem.MatchUserTags
	}
	if backupPoliciesItem.Name != nil {
		backupPoliciesMap["name"] = backupPoliciesItem.Name
	}
	if backupPoliciesItem.Plans != nil {
		plansList := []map[string]interface{}{}
		for _, plansItem := range backupPoliciesItem.Plans {
			plansList = append(plansList, dataSourceBackupPolicyCollectionBackupPoliciesPlansToMap(plansItem))
		}
		backupPoliciesMap["plans"] = plansList
	}
	if backupPoliciesItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceBackupPolicyCollectionBackupPoliciesResourceGroupToMap(*backupPoliciesItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		backupPoliciesMap["resource_group"] = resourceGroupList
	}
	if backupPoliciesItem.ResourceType != nil {
		backupPoliciesMap["resource_type"] = backupPoliciesItem.ResourceType
	}

	return backupPoliciesMap
}

func dataSourceBackupPolicyCollectionBackupPoliciesPlansToMap(plansItem vpcv1.BackupPolicyPlanReference) (plansMap map[string]interface{}) {
	plansMap = map[string]interface{}{}

	if plansItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyCollectionPlansDeletedToMap(*plansItem.Deleted)
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

func dataSourceBackupPolicyCollectionPlansDeletedToMap(deletedItem vpcv1.BackupPolicyPlanReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyCollectionBackupPoliciesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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
