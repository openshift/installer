// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBackupPolicyJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyJobsRead,

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy identifier.",
			},
			"source_id": {
				Type:        schema.TypeString,
				Description: "Filters the collection to backup policy jobs with a source with the specified identifier",
				Optional:    true,
			},
			"target_snapshots_id": {
				Type:          schema.TypeSet,
				ConflictsWith: []string{"target_snapshots_crn"},
				Description:   "Filters the collection to resources with the target snapshot with the specified identifier",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				Optional:      true,
			},
			"target_snapshots_crn": {
				Type:          schema.TypeSet,
				ConflictsWith: []string{"target_snapshots_id"},
				Description:   "Filters the collection to backup policy jobs with the target snapshot with the specified CRN",
				Elem:          &schema.Schema{Type: schema.TypeString},
				Set:           schema.HashString,
				Optional:      true,
			},
			"backup_policy_plan_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to backup policy jobs with the backup plan with the specified identifier.",
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Filters the collection to backup policy jobs with the specified status",
				Optional:    true,
			},
			"jobs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of backup policy jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy job.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsBackupPolicyJobsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listBackupPolicyJobsOptions := &vpcv1.ListBackupPolicyJobsOptions{}

	listBackupPolicyJobsOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))

	if sourceId, ok := d.GetOk("source_id"); ok {
		listBackupPolicyJobsOptions.SetSourceID(sourceId.(string))
	}

	if targetSnapshotsId, ok := d.GetOk("target_snapshots_id"); ok {
		targetSnapshotsIds := targetSnapshotsId.(*schema.Set)
		if targetSnapshotsIds.Len() != 0 {
			targetSnapshotsIdsArray := make([]string, targetSnapshotsIds.Len())
			for i, key := range targetSnapshotsIds.List() {
				keystr := key.(string)
				targetSnapshotsIdsArray[i] = keystr
			}
			listBackupPolicyJobsOptions.SetTargetSnapshotsID(strings.Join(targetSnapshotsIdsArray, ","))
		}
	}

	if targetSnapshotsCrn, ok := d.GetOk("target_snapshots_crn"); ok {
		targetSnapshotsCrns := targetSnapshotsCrn.(*schema.Set)
		if targetSnapshotsCrns.Len() != 0 {
			targetSnapshotsCrnsArray := make([]string, targetSnapshotsCrns.Len())
			for i, key := range targetSnapshotsCrns.List() {
				keystr := key.(string)
				targetSnapshotsCrnsArray[i] = keystr
			}
			listBackupPolicyJobsOptions.SetTargetSnapshotsCRN(strings.Join(targetSnapshotsCrnsArray, ","))
		}
	}

	if status, ok := d.GetOk("status"); ok {
		listBackupPolicyJobsOptions.SetStatus(status.(string))
	}

	if backupPolicyPlanId, ok := d.GetOk("backup_policy_plan_id"); ok {
		listBackupPolicyJobsOptions.SetBackupPolicyPlanID(backupPolicyPlanId.(string))
	}

	// Support for pagination
	start := ""
	allrecs := []vpcv1.BackupPolicyJob{}

	for {

		if start != "" {
			listBackupPolicyJobsOptions.Start = &start
		}

		backupPolicyJobCollection, response, err := vpcClient.ListBackupPolicyJobsWithContext(context, listBackupPolicyJobsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListBackupPolicyJobsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListBackupPolicyJobsWithContext failed %s\n%s", err, response))
		}

		if backupPolicyJobCollection != nil && *backupPolicyJobCollection.TotalCount == int64(0) {
			break
		}
		start = flex.GetNext(backupPolicyJobCollection.Next)
		allrecs = append(allrecs, backupPolicyJobCollection.Jobs...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsBackupPolicyJobsID(d))

	if allrecs != nil {
		err = d.Set("jobs", dataSourceBackupPolicyJobCollectionFlattenJobs(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting jobs %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsBackupPolicyJobsID returns a reasonable ID for the list.
func dataSourceIBMIsBackupPolicyJobsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceBackupPolicyJobCollectionFlattenJobs(result []vpcv1.BackupPolicyJob) (jobs []map[string]interface{}) {
	for _, jobsItem := range result {
		jobs = append(jobs, dataSourceBackupPolicyJobCollectionJobsToMap(jobsItem))
	}

	return jobs
}

func dataSourceBackupPolicyJobCollectionJobsToMap(jobsItem vpcv1.BackupPolicyJob) (jobsMap map[string]interface{}) {
	// log.Println("Hi I am inside dataSourceBackupPolicyJobCollectionJobsToMap")
	jobsMap = map[string]interface{}{}

	if jobsItem.AutoDelete != nil {
		jobsMap["auto_delete"] = jobsItem.AutoDelete
	}
	if jobsItem.AutoDeleteAfter != nil {
		jobsMap["auto_delete_after"] = jobsItem.AutoDeleteAfter
	}
	if jobsItem.BackupPolicyPlan != nil {
		backupPolicyPlanList := []map[string]interface{}{}
		backupPolicyPlanMap := dataSourceBackupPolicyJobCollectionJobsBackupPolicyPlanToMap(*jobsItem.BackupPolicyPlan)
		backupPolicyPlanList = append(backupPolicyPlanList, backupPolicyPlanMap)
		jobsMap["backup_policy_plan"] = backupPolicyPlanList
	}
	if jobsItem.CompletedAt != nil {
		jobsMap["completed_at"] = jobsItem.CompletedAt.String()
	}
	if jobsItem.CreatedAt != nil {
		jobsMap["created_at"] = jobsItem.CreatedAt.String()
	}
	if jobsItem.Href != nil {
		jobsMap["href"] = jobsItem.Href
	}
	if jobsItem.ID != nil {
		jobsMap["id"] = jobsItem.ID
	}
	if jobsItem.JobType != nil {
		jobsMap["job_type"] = jobsItem.JobType
	}
	if jobsItem.ResourceType != nil {
		jobsMap["resource_type"] = jobsItem.ResourceType
	}
	log.Println("jobsItem.Source")
	log.Println(jobsItem.Source)
	if jobsItem.Source != nil {
		switch reflect.TypeOf(jobsItem.Source).String() {
		case "*vpcv1.BackupPolicyJobSourceVolumeReference":
			{
				jobSource := jobsItem.Source.(*vpcv1.BackupPolicyJobSourceVolumeReference)
				sourceVolumeList := []map[string]interface{}{}
				sourceVolumeMap := dataSourceBackupPolicyJobCollectionJobsSourceVolumeToMap(*jobSource)
				sourceVolumeList = append(sourceVolumeList, sourceVolumeMap)
				jobsMap["source_volume"] = sourceVolumeList
			}
		case "*vpcv1.BackupPolicyJobSourceInstanceReference":
			{
				jobSource := jobsItem.Source.(*vpcv1.BackupPolicyJobSourceInstanceReference)
				sourceVolumeList := []map[string]interface{}{}
				sourceVolumeMap := dataSourceBackupPolicyJobCollectionJobsSourceInstanceToMap(*jobSource)
				sourceVolumeList = append(sourceVolumeList, sourceVolumeMap)
				jobsMap["source_instance"] = sourceVolumeList
			}
		}
	}
	if jobsItem.Status != nil {
		jobsMap["status"] = jobsItem.Status
	}
	if jobsItem.StatusReasons != nil {
		statusReasonsList := []map[string]interface{}{}
		for _, statusReasonsItem := range jobsItem.StatusReasons {
			statusReasonsList = append(statusReasonsList, dataSourceBackupPolicyJobCollectionJobsStatusReasonsToMap(statusReasonsItem))
		}
		jobsMap["status_reasons"] = statusReasonsList
	}
	log.Println("jobsItem.TargetSnapshot")
	log.Println(jobsItem.TargetSnapshots)
	if jobsItem.TargetSnapshots != nil {
		targetSnapshotList := []map[string]interface{}{}
		for _, targetSnapshotsItem := range jobsItem.TargetSnapshots {
			targetSnapshotMap := dataSourceBackupPolicyJobCollectionJobsTargetSnapshotToMap(targetSnapshotsItem)
			targetSnapshotList = append(targetSnapshotList, targetSnapshotMap)
		}
		jobsMap["target_snapshot"] = targetSnapshotList
	}
	// log.Println("jobsItem")
	// log.Println(jobsItem)
	return jobsMap
}

func dataSourceBackupPolicyJobCollectionJobsBackupPolicyPlanToMap(backupPolicyPlanItem vpcv1.BackupPolicyPlanReference) (backupPolicyPlanMap map[string]interface{}) {
	backupPolicyPlanMap = map[string]interface{}{}

	if backupPolicyPlanItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionBackupPolicyPlanDeletedToMap(*backupPolicyPlanItem.Deleted)
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
		remoteMap, err := resourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanRemoteToMap(backupPolicyPlanItem.Remote)
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

func dataSourceBackupPolicyJobCollectionBackupPolicyPlanDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobCollectionJobsSourceVolumeToMap(sourceVolumeItem vpcv1.BackupPolicyJobSourceVolumeReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionSourceVolumeDeletedToMap(*sourceVolumeItem.Deleted)
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
	if sourceVolumeItem.Remote != nil {
		remoteMap, err := resourceIBMIsSnapshotConsistencyGroupVolumeRemoteToMap(sourceVolumeItem.Remote)
		if err != nil {
			return remoteMap
		}
		sourceVolumeMap["remote"] = []map[string]interface{}{remoteMap}
	}
	if sourceVolumeItem.ResourceType != nil {
		sourceVolumeMap["resource_type"] = sourceVolumeItem.ResourceType
	}
	return sourceVolumeMap
}

func dataSourceBackupPolicyJobCollectionJobsSourceInstanceToMap(sourceVolumeItem vpcv1.BackupPolicyJobSourceInstanceReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionSourceInstanceDeletedToMap(*sourceVolumeItem.Deleted)
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

func dataSourceBackupPolicyJobCollectionSourceVolumeDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
func dataSourceBackupPolicyJobCollectionSourceInstanceDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobCollectionJobsStatusReasonsToMap(statusReasonsItem vpcv1.BackupPolicyJobStatusReason) (statusReasonsMap map[string]interface{}) {
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

func dataSourceBackupPolicyJobCollectionJobsTargetSnapshotToMap(targetSnapshotItem vpcv1.SnapshotReference) (targetSnapshotMap map[string]interface{}) {
	targetSnapshotMap = map[string]interface{}{}

	if targetSnapshotItem.CRN != nil {
		targetSnapshotMap["crn"] = targetSnapshotItem.CRN
	}
	if targetSnapshotItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionTargetSnapshotDeletedToMap(*targetSnapshotItem.Deleted)
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

func dataSourceBackupPolicyJobCollectionTargetSnapshotDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobCollectionNextToMap(nextItem vpcv1.BackupPolicyJobCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}
