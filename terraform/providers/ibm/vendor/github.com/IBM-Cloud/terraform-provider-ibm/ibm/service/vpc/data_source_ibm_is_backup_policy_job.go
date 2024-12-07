// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBackupPolicyJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyJobRead,

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy identifier.",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy job identifier.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this backup policy job will be automatically deleted after it completes. At present, this is always `true`, but may be modifiable in the future.",
			},
			"auto_delete_after": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "If `auto_delete` is `true`, the days after completion that this backup policy job will be deleted. This value may be modifiable in the future.",
			},
			"backup_policy_plan": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The backup policy plan operated this backup policy job (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
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
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this reference is remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this region.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this region.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy job was completed.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy job was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy job.",
			},
			"job_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of backup policy job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the backup policy job on which the unexpected property value was encountered.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"source_volume": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source volume this backup was created from (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this volume.",
						},
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
							Description: "The URL for this volume.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this volume.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this reference is remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this region.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this region.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"source_instance": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source volume this backup was created from (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this volume.",
						},
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
							Description: "The URL for this volume.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this volume.",
						},
					},
				},
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the backup policy job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the backup policy job on which the unexpected property value was encountered.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason:- `internal_error`: Internal error (contact IBM support)- `snapshot_pending`: Cannot delete backup (snapshot) in the `pending` lifecycle state- `snapshot_volume_limit`: The snapshot limit for the source volume has been reached- `source_volume_busy`: The source volume has `busy` set (after multiple retries).",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
			"target_snapshot": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The snapshot operated on by this backup policy job (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this snapshot.",
						},
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
							Description: "The URL for this snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this snapshot.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this snapshot.",
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

func dataSourceIBMIsBackupPolicyJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyJobOptions := &vpcv1.GetBackupPolicyJobOptions{}

	getBackupPolicyJobOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))
	getBackupPolicyJobOptions.SetID(d.Get("identifier").(string))

	backupPolicyJob, response, err := vpcClient.GetBackupPolicyJobWithContext(context, getBackupPolicyJobOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBackupPolicyJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBackupPolicyJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(*backupPolicyJob.ID)
	if err = d.Set("auto_delete", backupPolicyJob.AutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting auto_delete: %s", err))
	}
	if err = d.Set("auto_delete_after", flex.IntValue(backupPolicyJob.AutoDeleteAfter)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting auto_delete_after: %s", err))
	}

	if backupPolicyJob.BackupPolicyPlan != nil {
		err = d.Set("backup_policy_plan", dataSourceBackupPolicyJobFlattenBackupPolicyPlan(*backupPolicyJob.BackupPolicyPlan))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_policy_plan %s", err))
		}
	}
	if err = d.Set("completed_at", flex.DateTimeToString(backupPolicyJob.CompletedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting completed_at: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(backupPolicyJob.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", backupPolicyJob.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("job_type", backupPolicyJob.JobType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting job_type: %s", err))
	}
	if err = d.Set("resource_type", backupPolicyJob.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if backupPolicyJob.Source != nil {
		switch reflect.TypeOf(backupPolicyJob.Source).String() {
		case "*vpcv1.BackupPolicyJobSourceVolumeReference":
			{
				jobSource := backupPolicyJob.Source.(*vpcv1.BackupPolicyJobSourceVolumeReference)
				err = d.Set("source_volume", dataSourceBackupPolicyJobFlattenSourceVolume(*jobSource))
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting source_volume %s", err))
				}
			}
		case "*vpcv1.BackupPolicyJobSourceInstanceReference":
			{
				jobSource := backupPolicyJob.Source.(*vpcv1.BackupPolicyJobSourceInstanceReference)
				err = d.Set("source_instance", dataSourceBackupPolicyJobFlattenSourceInstance(*jobSource))
				if err != nil {
					return diag.FromErr(fmt.Errorf("Error setting source_instance %s", err))
				}
			}
		}

	}
	if err = d.Set("status", backupPolicyJob.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if backupPolicyJob.StatusReasons != nil {
		err = d.Set("status_reasons", dataSourceBackupPolicyJobFlattenStatusReasons(backupPolicyJob.StatusReasons))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_reasons %s", err))
		}
	}

	if backupPolicyJob.TargetSnapshots != nil {
		err = d.Set("target_snapshot", dataSourceBackupPolicyJobFlattenTargetSnapshot(backupPolicyJob.TargetSnapshots))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting target_snapshot %s", err))
		}
	}

	return nil
}

func dataSourceBackupPolicyJobFlattenBackupPolicyPlan(result vpcv1.BackupPolicyPlanReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyJobBackupPolicyPlanToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyJobBackupPolicyPlanToMap(backupPolicyPlanItem vpcv1.BackupPolicyPlanReference) (backupPolicyPlanMap map[string]interface{}) {
	backupPolicyPlanMap = map[string]interface{}{}

	if backupPolicyPlanItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobBackupPolicyPlanDeletedToMap(*backupPolicyPlanItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		backupPolicyPlanMap["deleted"] = deletedList
	}
	if backupPolicyPlanItem.Href != nil {
		backupPolicyPlanMap["href"] = backupPolicyPlanItem.Href
	}
	if backupPolicyPlanItem.ID != nil {
		backupPolicyPlanMap["id"] = backupPolicyPlanItem.ID
	}
	if backupPolicyPlanItem.Name != nil {
		backupPolicyPlanMap["name"] = backupPolicyPlanItem.Name
	}
	if backupPolicyPlanItem.Remote != nil {
		remoteMap, err := resourceIBMIsBackupPolicyPlanRemoteToMap(backupPolicyPlanItem.Remote)
		if err != nil {
			return remoteMap
		}
		backupPolicyPlanMap["remote"] = []map[string]interface{}{remoteMap}
	}
	if backupPolicyPlanItem.ResourceType != nil {
		backupPolicyPlanMap["resource_type"] = backupPolicyPlanItem.ResourceType
	}
	return backupPolicyPlanMap
}

func dataSourceBackupPolicyJobBackupPolicyPlanDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobFlattenSourceVolume(result vpcv1.BackupPolicyJobSourceVolumeReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyJobSourceVolumeToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyJobSourceVolumeToMap(sourceVolumeItem vpcv1.BackupPolicyJobSourceVolumeReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobSourceVolumeDeletedToMap(*sourceVolumeItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		sourceVolumeMap["deleted"] = deletedList
	}
	if sourceVolumeItem.Href != nil {
		sourceVolumeMap["href"] = sourceVolumeItem.Href
	}
	if sourceVolumeItem.ID != nil {
		sourceVolumeMap["id"] = sourceVolumeItem.ID
	}
	if sourceVolumeItem.Name != nil {
		sourceVolumeMap["name"] = sourceVolumeItem.Name
	}
	if sourceVolumeItem.ResourceType != nil {
		sourceVolumeMap["resource_type"] = sourceVolumeItem.ResourceType
	}
	if sourceVolumeItem.Remote != nil {
		remoteMap, err := resourceIBMIsSnapshotConsistencyGroupVolumeRemoteToMap(sourceVolumeItem.Remote)
		if err != nil {
			return remoteMap
		}
		sourceVolumeMap["remote"] = []map[string]interface{}{remoteMap}
	}
	return sourceVolumeMap
}

func resourceIBMIsSnapshotConsistencyGroupVolumeRemoteToMap(model *vpcv1.VolumeRemote) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := resourceIBMIsVolumeConsistencyGroupRegionReferenceToMap(model.Region)
		if err != nil {
			return regionMap, err
		}
		return regionMap, nil
	}
	return regionMap, nil
}

func resourceIBMIsVolumeConsistencyGroupRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	regionMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := resourceIBMIsBackupPolicyPlanRemoteRegionReferenceToMap(model.Region)
		if err != nil {
			return regionMap, err
		}
		return regionMap, nil
	}
	return regionMap, nil
}

func resourceIBMIsBackupPolicyPlanRemoteRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceBackupPolicyJobSourceVolumeDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
func dataSourceBackupPolicyJobFlattenSourceInstance(result vpcv1.BackupPolicyJobSourceInstanceReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyJobSourceInstanceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyJobSourceInstanceToMap(sourceVolumeItem vpcv1.BackupPolicyJobSourceInstanceReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobSourceInstanceDeletedToMap(*sourceVolumeItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		sourceVolumeMap["deleted"] = deletedList
	}
	if sourceVolumeItem.Href != nil {
		sourceVolumeMap["href"] = sourceVolumeItem.Href
	}
	if sourceVolumeItem.ID != nil {
		sourceVolumeMap["id"] = sourceVolumeItem.ID
	}
	if sourceVolumeItem.Name != nil {
		sourceVolumeMap["name"] = sourceVolumeItem.Name
	}
	return sourceVolumeMap
}

func dataSourceBackupPolicyJobSourceInstanceDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobFlattenStatusReasons(result []vpcv1.BackupPolicyJobStatusReason) (statusReasons []map[string]interface{}) {
	for _, statusReasonsItem := range result {
		statusReasons = append(statusReasons, dataSourceBackupPolicyJobStatusReasonsToMap(statusReasonsItem))
	}

	return statusReasons
}

func dataSourceBackupPolicyJobStatusReasonsToMap(statusReasonsItem vpcv1.BackupPolicyJobStatusReason) (statusReasonsMap map[string]interface{}) {
	statusReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		statusReasonsMap["code"] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		statusReasonsMap["message"] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		statusReasonsMap["more_info"] = statusReasonsItem.MoreInfo
	}

	return statusReasonsMap
}

func dataSourceBackupPolicyJobFlattenTargetSnapshot(result []vpcv1.SnapshotReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	for _, snapshotReferenceItem := range result {
		finalMap := dataSourceBackupPolicyJobTargetSnapshotToMap(snapshotReferenceItem)
		finalList = append(finalList, finalMap)
	}

	return finalList
}

func dataSourceBackupPolicyJobTargetSnapshotToMap(targetSnapshotItem vpcv1.SnapshotReference) (targetSnapshotMap map[string]interface{}) {
	targetSnapshotMap = map[string]interface{}{}

	if targetSnapshotItem.CRN != nil {
		targetSnapshotMap["crn"] = targetSnapshotItem.CRN
	}
	if targetSnapshotItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobTargetSnapshotDeletedToMap(*targetSnapshotItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		targetSnapshotMap["deleted"] = deletedList
	}
	if targetSnapshotItem.Href != nil {
		targetSnapshotMap["href"] = targetSnapshotItem.Href
	}
	if targetSnapshotItem.ID != nil {
		targetSnapshotMap["id"] = targetSnapshotItem.ID
	}
	if targetSnapshotItem.Name != nil {
		targetSnapshotMap["name"] = targetSnapshotItem.Name
	}
	if targetSnapshotItem.ResourceType != nil {
		targetSnapshotMap["resource_type"] = targetSnapshotItem.ResourceType
	}

	return targetSnapshotMap
}

func dataSourceBackupPolicyJobTargetSnapshotDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
