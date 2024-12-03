// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsSnapshotConsistencyGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsSnapshotConsistencyGroupRead,

		Schema: map[string]*schema.Schema{
			"identifier": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_snapshot_consistency_group", "identifier"),
				Description:  "The snapshot consistency group identifier.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_snapshot_consistency_group", isSnapshotName),
				Description:  "The name for this snapshot consistency group. The name is unique across all snapshot consistency groups in the region.",
			},
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
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of this snapshot consistency group.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this snapshot consistency group.",
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
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Snapshot Consistency Group tags list",
			},
			"access_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_snapshot_consistency_group", "access_tags")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
		},
	}
}

func DataSourceIBMISSnapshotConsistencyGroupValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "identifier",
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isSnapshotName,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString})

	ibmISSnapshotDataSourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_snapshot_consistency_group", Schema: validateSchema}
	return &ibmISSnapshotDataSourceValidator
}

func dataSourceIBMIsSnapshotConsistencyGroupRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("name").(string)
	id := d.Get("identifier").(string)

	if name != "" {
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
		for _, snapshotConsistencyGroup := range allrecs {
			if *snapshotConsistencyGroup.Name == name || *snapshotConsistencyGroup.ID == id {

				d.SetId(fmt.Sprintf("%s", *snapshotConsistencyGroup.ID))

				backupPolicyPlan := []map[string]interface{}{}
				if snapshotConsistencyGroup.BackupPolicyPlan != nil {
					modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceToMap(snapshotConsistencyGroup.BackupPolicyPlan)
					if err != nil {
						return diag.FromErr(err)
					}
					backupPolicyPlan = append(backupPolicyPlan, modelMap)
				}
				if err = d.Set("backup_policy_plan", backupPolicyPlan); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting backup_policy_plan %s", err))
				}

				if err = d.Set("created_at", flex.DateTimeToString(snapshotConsistencyGroup.CreatedAt)); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
				}

				if err = d.Set("crn", snapshotConsistencyGroup.CRN); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
				}

				if err = d.Set("delete_snapshots_on_delete", snapshotConsistencyGroup.DeleteSnapshotsOnDelete); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting delete_snapshots_on_delete: %s", err))
				}

				if err = d.Set("href", snapshotConsistencyGroup.Href); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
				}

				if err = d.Set("lifecycle_state", snapshotConsistencyGroup.LifecycleState); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
				}

				if err = d.Set("name", snapshotConsistencyGroup.Name); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
				}

				resourceGroup := []map[string]interface{}{}
				if snapshotConsistencyGroup.ResourceGroup != nil {
					modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupResourceGroupReferenceToMap(snapshotConsistencyGroup.ResourceGroup)
					if err != nil {
						return diag.FromErr(err)
					}
					resourceGroup = append(resourceGroup, modelMap)
				}
				if err = d.Set("resource_group", resourceGroup); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
				}

				if err = d.Set("resource_type", snapshotConsistencyGroup.ResourceType); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
				}

				snapshots := []map[string]interface{}{}
				if snapshotConsistencyGroup.Snapshots != nil {
					for _, modelItem := range snapshotConsistencyGroup.Snapshots {
						modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(&modelItem)
						if err != nil {
							return diag.FromErr(err)
						}
						snapshots = append(snapshots, modelMap)
					}
				}
				if err = d.Set("snapshots", snapshots); err != nil {
					return diag.FromErr(fmt.Errorf("Error setting snapshots %s", err))
				}
				tags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
				if err != nil {
					log.Printf(
						"Error on get of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
				}
				d.Set("tags", tags)

				accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
				if err != nil {
					log.Printf(
						"Error on get of resource VPC snapshot consistency group (%s) access tags: %s", d.Id(), err)
				}
				d.Set("access_tags", accesstags)
				return nil
			}
		}
		return diag.FromErr(fmt.Errorf("[ERROR] No snapshot consistency group found with name %s", name))
	} else {
		getSnapshotConsistencyGroupOptions := &vpcv1.GetSnapshotConsistencyGroupOptions{}

		getSnapshotConsistencyGroupOptions.SetID(id)

		snapshotConsistencyGroup, response, err := vpcClient.GetSnapshotConsistencyGroupWithContext(context, getSnapshotConsistencyGroupOptions)
		if err != nil {
			log.Printf("[DEBUG] GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetSnapshotConsistencyGroupWithContext failed %s\n%s", err, response))
		}

		d.SetId(fmt.Sprintf("%s", *getSnapshotConsistencyGroupOptions.ID))

		backupPolicyPlan := []map[string]interface{}{}
		if snapshotConsistencyGroup.BackupPolicyPlan != nil {
			modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceToMap(snapshotConsistencyGroup.BackupPolicyPlan)
			if err != nil {
				return diag.FromErr(err)
			}
			backupPolicyPlan = append(backupPolicyPlan, modelMap)
		}
		if err = d.Set("backup_policy_plan", backupPolicyPlan); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_policy_plan %s", err))
		}

		if err = d.Set("created_at", flex.DateTimeToString(snapshotConsistencyGroup.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}

		if err = d.Set("crn", snapshotConsistencyGroup.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}

		if err = d.Set("delete_snapshots_on_delete", snapshotConsistencyGroup.DeleteSnapshotsOnDelete); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting delete_snapshots_on_delete: %s", err))
		}

		if err = d.Set("href", snapshotConsistencyGroup.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}

		if err = d.Set("lifecycle_state", snapshotConsistencyGroup.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
		}

		if err = d.Set("name", snapshotConsistencyGroup.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}

		resourceGroup := []map[string]interface{}{}
		if snapshotConsistencyGroup.ResourceGroup != nil {
			modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupResourceGroupReferenceToMap(snapshotConsistencyGroup.ResourceGroup)
			if err != nil {
				return diag.FromErr(err)
			}
			resourceGroup = append(resourceGroup, modelMap)
		}
		if err = d.Set("resource_group", resourceGroup); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}

		if err = d.Set("resource_type", snapshotConsistencyGroup.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}

		snapshots := []map[string]interface{}{}
		if snapshotConsistencyGroup.Snapshots != nil {
			for _, modelItem := range snapshotConsistencyGroup.Snapshots {
				modelMap, err := dataSourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(&modelItem)
				if err != nil {
					return diag.FromErr(err)
				}
				snapshots = append(snapshots, modelMap)
			}
		}
		if err = d.Set("snapshots", snapshots); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting snapshots %s", err))
		}
		tags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc snapshot consistency group (%s) tags: %s", d.Id(), err)
		}
		d.Set("tags", tags)

		accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *snapshotConsistencyGroup.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error on get of resource VPC snapshot consistency group (%s) access tags: %s", d.Id(), err)
		}
		d.Set("access_tags", accesstags)
		return nil
	}

}

func dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceToMap(model *vpcv1.BackupPolicyPlanReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupBackupPolicyPlanRemoteToMap(model *vpcv1.BackupPolicyPlanRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		regionMap, err := dataSourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		return regionMap, nil
	}
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupSnapshotConsistencyGroupSnapshotsItemToMap(model *vpcv1.SnapshotReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsSnapshotConsistencyGroupSnapshotReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	if model.Remote != nil {
		remoteMap, err := dataSourceIBMIsSnapshotConsistencyGroupSnapshotRemoteToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupSnapshotReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsSnapshotConsistencyGroupSnapshotRemoteToMap(model *vpcv1.SnapshotRemote) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Region != nil {
		log.Println("model.Region")
		log.Println(model.Region)
		regionMap, err := dataSourceIBMIsSnapshotConsistencyGroupRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		return regionMap, nil

	}
	return modelMap, nil
}
