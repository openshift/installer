// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsShareSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"share": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The file share identifier, or `-` to wildcard all accessible file shares.",
			},
			"backup_policy_plan": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to backup policy jobs with a `backup_policy_plan.id` property matching the specified identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
			},
			"snapshots": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of share snapshots.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_policy_plan": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, the backup policy plan which created this share snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
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
										Description: "The name for this backup policy plan. The name is unique across all plans in the backup policy.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
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
						"captured_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time the data capture for this share snapshot was completed.If absent, this snapshot's data has not yet been captured.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the share snapshot was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this share snapshot.",
						},
						"fingerprint": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fingerprint for this snapshot.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share snapshot.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this share snapshot.",
						},
						"minimum_size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum size of a share created from this snapshot. When a snapshot is created, this will be set to the size of the `source_share`.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this share snapshot. The name is unique across all snapshots for the file share.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this file share.",
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
										Description: "The name for this resource group.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the share snapshot:- `available`: The share snapshot is available for use.- `failed`: The share snapshot is irrecoverably unusable.- `pending`: The share snapshot is being provisioned and is not yet usable.- `unusable`: The share snapshot is not currently usable (see `status_reasons`)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current status (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A reason code for the status:- `encryption_key_deleted`: File share snapshot is unusable  because its `encryption_key` was deleted- `internal_error`: Internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
						"user_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The [user tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags) associated with this share snapshot.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"zone": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The zone this share snapshot resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this zone.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The globally unique name for this zone.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The [user tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags) associated with this share snapshot.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsShareSnapshotsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshots", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listShareSnapshotsOptions := &vpcv1.ListShareSnapshotsOptions{}

	if shareIntf, ok := d.GetOk("share"); ok {
		listShareSnapshotsOptions.SetShareID(shareIntf.(string))
	} else {
		listShareSnapshotsOptions.SetShareID("-")
	}

	if _, ok := d.GetOk("backup_policy_plan"); ok {
		listShareSnapshotsOptions.SetBackupPolicyPlanID(d.Get("backup_policy_plan").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		listShareSnapshotsOptions.SetName(d.Get("name").(string))
	}

	var pager *vpcv1.ShareSnapshotsPager
	pager, err = vpcClient.NewShareSnapshotsPager(listShareSnapshotsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshots", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ShareSnapshotsPager.GetAll() failed %s", err), "(Data) ibm_is_share_snapshots", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsShareSnapshotsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsShareSnapshotsShareSnapshotToMap(&modelItem)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshots", "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("snapshots", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting snapshots %s", err), "(Data) ibm_is_share_snapshots", "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIBMIsShareSnapshotsID returns a reasonable ID for the list.
func dataSourceIBMIsShareSnapshotsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsShareSnapshotsShareSnapshotToMap(model *vpcv1.ShareSnapshot) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BackupPolicyPlan != nil {
		backupPolicyPlanMap, err := DataSourceIBMIsShareSnapshotsBackupPolicyPlanReferenceToMap(model.BackupPolicyPlan)
		if err != nil {
			return modelMap, err
		}
		modelMap["backup_policy_plan"] = []map[string]interface{}{backupPolicyPlanMap}
	}
	if model.CapturedAt != nil {
		modelMap["captured_at"] = model.CapturedAt.String()
	}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = *model.CRN
	modelMap["fingerprint"] = *model.Fingerprint
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["lifecycle_state"] = *model.LifecycleState
	modelMap["minimum_size"] = flex.IntValue(model.MinimumSize)
	modelMap["name"] = *model.Name
	if model.ResourceGroup != nil {
		resourceGroupMap, err := DataSourceIBMIsShareSnapshotsResourceGroupReferenceToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	modelMap["status"] = *model.Status
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range model.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsShareSnapshotsShareSnapshotStatusReasonToMap(&statusReasonsItem)
		if err != nil {
			return modelMap, err
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	modelMap["status_reasons"] = statusReasons
	if model.UserTags != nil {
		modelMap["user_tags"] = model.UserTags
	}
	zoneMap, err := DataSourceIBMIsShareSnapshotsZoneReferenceToMap(model.Zone)
	if err != nil {
		return modelMap, err
	}
	modelMap["zone"] = []map[string]interface{}{zoneMap}

	if model.UserTags != nil {
		modelMap["tags"] = model.UserTags
	}

	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsBackupPolicyPlanReferenceToMap(model *vpcv1.BackupPolicyPlanReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsShareSnapshotsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsShareSnapshotsBackupPolicyPlanRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsShareSnapshotsRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsShareSnapshotStatusReasonToMap(model *vpcv1.ShareSnapshotStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotsZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
