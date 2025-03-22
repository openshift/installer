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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsShareSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareSnapshotRead,

		Schema: map[string]*schema.Schema{
			"share": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"share_snapshot": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The share snapshot identifier.",
			},
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
	}
}

func dataSourceIBMIsShareSnapshotRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getShareSnapshotOptions := &vpcv1.GetShareSnapshotOptions{}

	getShareSnapshotOptions.SetShareID(d.Get("share").(string))
	getShareSnapshotOptions.SetID(d.Get("share_snapshot").(string))

	shareSnapshot, _, err := vpcClient.GetShareSnapshotWithContext(context, getShareSnapshotOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareSnapshotWithContext failed: %s", err.Error()), "(Data) ibm_is_share_snapshot", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getShareSnapshotOptions.ShareID, *getShareSnapshotOptions.ID))

	backupPolicyPlan := []map[string]interface{}{}
	if shareSnapshot.BackupPolicyPlan != nil {
		modelMap, err := DataSourceIBMIsShareSnapshotBackupPolicyPlanReferenceToMap(shareSnapshot.BackupPolicyPlan)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshot", "read", "backup_policy_plan-to-map").GetDiag()
		}
		backupPolicyPlan = append(backupPolicyPlan, modelMap)
	}
	if err = d.Set("backup_policy_plan", backupPolicyPlan); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting backup_policy_plan: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-backup_policy_plan").GetDiag()
	}

	if err = d.Set("captured_at", flex.DateTimeToString(shareSnapshot.CapturedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting captured_at: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-captured_at").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(shareSnapshot.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("crn", shareSnapshot.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-crn").GetDiag()
	}

	if err = d.Set("fingerprint", shareSnapshot.Fingerprint); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting fingerprint: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-fingerprint").GetDiag()
	}

	if err = d.Set("href", shareSnapshot.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-href").GetDiag()
	}

	if err = d.Set("lifecycle_state", shareSnapshot.LifecycleState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting lifecycle_state: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-lifecycle_state").GetDiag()
	}

	if err = d.Set("minimum_size", flex.IntValue(shareSnapshot.MinimumSize)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting minimum_size: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-minimum_size").GetDiag()
	}

	if err = d.Set("name", shareSnapshot.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-name").GetDiag()
	}

	resourceGroup := []map[string]interface{}{}
	if shareSnapshot.ResourceGroup != nil {
		modelMap, err := DataSourceIBMIsShareSnapshotResourceGroupReferenceToMap(shareSnapshot.ResourceGroup)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshot", "read", "resource_group-to-map").GetDiag()
		}
		resourceGroup = append(resourceGroup, modelMap)
	}
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-resource_group").GetDiag()
	}

	if err = d.Set("resource_type", shareSnapshot.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-resource_type").GetDiag()
	}

	if err = d.Set("status", shareSnapshot.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-status").GetDiag()
	}
	if shareSnapshot.UserTags != nil {
		if err = d.Set("tags", shareSnapshot.UserTags); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-tags").GetDiag()
		}
	}
	statusReasons := []map[string]interface{}{}
	if shareSnapshot.StatusReasons != nil {
		for _, modelItem := range shareSnapshot.StatusReasons {
			modelMap, err := DataSourceIBMIsShareSnapshotShareSnapshotStatusReasonToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshot", "read", "status_reasons-to-map").GetDiag()
			}
			statusReasons = append(statusReasons, modelMap)
		}
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status_reasons: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-status_reasons").GetDiag()
	}

	zone := []map[string]interface{}{}
	if shareSnapshot.Zone != nil {
		modelMap, err := DataSourceIBMIsShareSnapshotZoneReferenceToMap(shareSnapshot.Zone)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_share_snapshot", "read", "zone-to-map").GetDiag()
		}
		zone = append(zone, modelMap)
	}
	if err = d.Set("zone", zone); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting zone: %s", err), "(Data) ibm_is_share_snapshot", "read", "set-zone").GetDiag()
	}

	return nil
}

func DataSourceIBMIsShareSnapshotBackupPolicyPlanReferenceToMap(model *vpcv1.BackupPolicyPlanReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsShareSnapshotDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsShareSnapshotBackupPolicyPlanRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsShareSnapshotRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotShareSnapshotStatusReasonToMap(model *vpcv1.ShareSnapshotStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsShareSnapshotZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
