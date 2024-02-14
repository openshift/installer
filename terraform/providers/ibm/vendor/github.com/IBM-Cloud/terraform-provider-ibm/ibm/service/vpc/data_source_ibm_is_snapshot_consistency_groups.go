// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsSnapshotConsistencyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsSnapshotConsistencyGroupsRead,

		Schema: map[string]*schema.Schema{
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `resource_group.id` property matching the specified identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
			},
			"backup_policy_plan": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to backup policy jobs with a `backup_policy_plan.id` property matching the specified identifier.",
			},
			"snapshot_consistency_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of snapshot consistency groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"backup_policy_plan": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, the backup policy plan which created this snapshot consistency group.",
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
										Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
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
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this snapshot consistency group was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of this snapshot consistency group.",
						},
						"delete_snapshots_on_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether deleting the snapshot consistency group will also delete the snapshots in the group.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this snapshot consistency group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this snapshot consistency group.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of this snapshot consistency group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this snapshot consistency group. The name is unique across all snapshot consistency groups in the region.",
						},
						"resource_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource group info",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"service_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The [service tags](https://cloud.ibm.com/apidocs/tagging#types-of-tags)[`is.instance:` prefix](https://cloud.ibm.com/docs/vpc?topic=vpc-snapshots-vpc-faqs) associated with this snapshot consistency group.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"snapshots": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The member snapshots that are data-consistent with respect to captured time. (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of this snapshot.",
									},
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
										Description: "The name for this snapshot. The name is unique across all snapshots in the region.",
									},
									"remote": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource is remote to this region,and identifies the native region.",
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

						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "User Tags for the snapshot consistency group",
						},

						"access_tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Set:         flex.ResourceIBMVPCHash,
							Description: "List of access tags",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsSnapshotConsistencyGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.SnapshotConsistencyGroup{}
	for {
		listSnapshotConsistencyGroupsOptions := &vpcv1.ListSnapshotConsistencyGroupsOptions{}
		if start != "" {
			listSnapshotConsistencyGroupsOptions.Start = &start
		}
		if rgFilterOk, ok := d.GetOk("resource_group"); ok {
			rgFilter := rgFilterOk.(string)
			listSnapshotConsistencyGroupsOptions.ResourceGroupID = &rgFilter
		}
		if nameFilterOk, ok := d.GetOk("name"); ok {
			nameFilter := nameFilterOk.(string)
			listSnapshotConsistencyGroupsOptions.Name = &nameFilter
		}
		if backupPolicyPlanIdFilterOk, ok := d.GetOk("backup_policy_plan"); ok {
			backupPolicyPlanIdFilter := backupPolicyPlanIdFilterOk.(string)
			listSnapshotConsistencyGroupsOptions.BackupPolicyPlanID = &backupPolicyPlanIdFilter
		}

		snapshotConsistencyGroup, response, err := vpcClient.ListSnapshotConsistencyGroups(listSnapshotConsistencyGroupsOptions)
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error fetching snapshots %s\n%s", err, response))
		}
		start = flex.GetNext(snapshotConsistencyGroup.Next)
		allrecs = append(allrecs, snapshotConsistencyGroup.SnapshotConsistencyGroups...)
		if start == "" {
			break
		}
	}

	snapshotConsistencyGroupsInfo := make([]map[string]interface{}, 0)
	for _, snapshotConsistencyGroup := range allrecs {
		l := map[string]interface{}{
			"id":                         *snapshotConsistencyGroup.ID,
			"name":                       *snapshotConsistencyGroup.Name,
			"href":                       *snapshotConsistencyGroup.Href,
			"crn":                        *snapshotConsistencyGroup.CRN,
			"delete_snapshots_on_delete": *snapshotConsistencyGroup.DeleteSnapshotsOnDelete,
			"lifecycle_state":            *snapshotConsistencyGroup.LifecycleState,
			"resource_type":              *snapshotConsistencyGroup.ResourceType,
			"created_at":                 (*snapshotConsistencyGroup.CreatedAt).String(),
		}

		//backup policy plan
		backupPolicyPlanList := []map[string]interface{}{}
		if snapshotConsistencyGroup.BackupPolicyPlan != nil {
			backupPolicyPlan := map[string]interface{}{}
			if snapshotConsistencyGroup.BackupPolicyPlan.Deleted != nil {
				snapshotConsistencyGroupBackupPolicyPlanDeletedMap := map[string]interface{}{}
				snapshotConsistencyGroupBackupPolicyPlanDeletedMap["more_info"] = snapshotConsistencyGroup.BackupPolicyPlan.Deleted.MoreInfo
				backupPolicyPlan["deleted"] = []map[string]interface{}{snapshotConsistencyGroupBackupPolicyPlanDeletedMap}
			}
			backupPolicyPlan["href"] = snapshotConsistencyGroup.BackupPolicyPlan.Href
			backupPolicyPlan["id"] = snapshotConsistencyGroup.BackupPolicyPlan.ID
			backupPolicyPlan["name"] = snapshotConsistencyGroup.BackupPolicyPlan.Name
			backupPolicyPlan["resource_type"] = snapshotConsistencyGroup.BackupPolicyPlan.ResourceType
			backupPolicyPlanList = append(backupPolicyPlanList, backupPolicyPlan)
		}
		l[isSnapshotBackupPolicyPlan] = backupPolicyPlanList

		// source snapshot
		if !core.IsNil(snapshotConsistencyGroup.Snapshots) {
			snapshots := []map[string]interface{}{}
			for _, snapshotsItem := range snapshotConsistencyGroup.Snapshots {
				snapshotsItemMap, err := resourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(&snapshotsItem)
				if err != nil {
					return diag.FromErr(err)
				}
				snapshots = append(snapshots, snapshotsItemMap)
			}
			l["snapshots"] = snapshots
		} else {
			l["snapshots"] = []map[string]interface{}{}
		}

		if snapshotConsistencyGroup.ResourceGroup != nil && snapshotConsistencyGroup.ResourceGroup.ID != nil {
			l["resource_group"] = *snapshotConsistencyGroup.ResourceGroup.ID
		}

		if snapshotConsistencyGroup.ServiceTags != nil {
			l["service_tags"] = snapshotConsistencyGroup.ServiceTags
		}

		tags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
		}
		l["tags"] = tags

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"[ERROR] Error on get of resource snapshot (%s) access tags: %s", d.Id(), err)
		}
		l["access_tags"] = accesstags

		snapshotConsistencyGroupsInfo = append(snapshotConsistencyGroupsInfo, l)
	}

	d.SetId(dataSourceIBMIsSnapshotConsistencyGroupsID(d))
	if err = d.Set("snapshot_consistency_groups", snapshotConsistencyGroupsInfo); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting snapshot_consistency_groups %s", err))
	}
	return nil
}

// dataSourceIBMIsSnapshotConsistencyGroupsID returns a reasonable ID for the list.
func dataSourceIBMIsSnapshotConsistencyGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
