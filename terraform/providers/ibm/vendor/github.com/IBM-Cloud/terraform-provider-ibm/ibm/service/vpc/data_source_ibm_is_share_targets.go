// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsShareTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceIbmIsShareTargetsRead,
		DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use ibm_is_share_mount_targets instead",

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user-defined name for this share target.",
			},
			"share_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of share targets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this share target.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the share target was created.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share target.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share target.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the mount target.",
						},
						"mount_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"subnet": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnet associated with this file share target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this subnet.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this subnet.",
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"transit_encryption": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC to which this share target is allowing to mount the file share.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this VPC.",
									},
									"resource_type": {
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

func dataSourceIbmIsShareTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listShareTargetsOptions := &vpcbetav1.ListShareMountTargetsOptions{}

	listShareTargetsOptions.SetShareID(d.Get("share").(string))
	if name, ok := d.GetOk("name"); ok {
		listShareTargetsOptions.SetName(name.(string))
	}
	shareTargetCollection, response, err := vpcClient.ListShareMountTargetsWithContext(context, listShareTargetsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListShareTargetsWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(dataSourceIbmIsShareTargetsID(d))

	if shareTargetCollection.MountTargets != nil {
		err = d.Set("share_targets", dataSourceShareTargetCollectionFlattenTargets(shareTargetCollection.MountTargets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
		}
	}

	return nil
}

// dataSourceIbmIsShareTargetsID returns a reasonable ID for the list.
func dataSourceIbmIsShareTargetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareTargetCollectionFlattenTargets(result []vpcbetav1.ShareMountTarget) (targets []map[string]interface{}) {
	for _, targetsItem := range result {
		targets = append(targets, dataSourceShareTargetCollectionTargetsToMap(targetsItem))
	}

	return targets
}

func dataSourceShareTargetCollectionTargetsToMap(targetsItem vpcbetav1.ShareMountTarget) (targetsMap map[string]interface{}) {
	targetsMap = map[string]interface{}{}

	if targetsItem.CreatedAt != nil {
		targetsMap["created_at"] = targetsItem.CreatedAt.String()
	}
	if targetsItem.Href != nil {
		targetsMap["href"] = targetsItem.Href
	}
	if targetsItem.ID != nil {
		targetsMap["id"] = targetsItem.ID
	}
	if targetsItem.LifecycleState != nil {
		targetsMap["lifecycle_state"] = targetsItem.LifecycleState
	}
	if targetsItem.MountPath != nil {
		targetsMap["mount_path"] = targetsItem.MountPath
	}
	if targetsItem.Name != nil {
		targetsMap["name"] = targetsItem.Name
	}
	if targetsItem.TransitEncryption != nil {
		targetsMap["transit_encryption"] = *targetsItem.TransitEncryption
	}
	if targetsItem.ResourceType != nil {
		targetsMap["resource_type"] = targetsItem.ResourceType
	}
	if targetsItem.VPC != nil {
		vpcList := []map[string]interface{}{}
		vpcMap := dataSourceShareMountTargetVpcToMap(*targetsItem.VPC)
		vpcList = append(vpcList, vpcMap)
		targetsMap["vpc"] = vpcList
	}

	return targetsMap
}

func dataSourceShareTargetCollectionTargetsSubnetToMap(subnetItem vpcbetav1.SubnetReference) (subnetMap map[string]interface{}) {
	subnetMap = map[string]interface{}{}

	if subnetItem.CRN != nil {
		subnetMap["crn"] = subnetItem.CRN
	}
	if subnetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetCollectionSubnetDeletedToMap(*subnetItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetMap["deleted"] = deletedList
	}
	if subnetItem.Href != nil {
		subnetMap["href"] = subnetItem.Href
	}
	if subnetItem.ID != nil {
		subnetMap["id"] = subnetItem.ID
	}
	if subnetItem.Name != nil {
		subnetMap["name"] = subnetItem.Name
	}

	return subnetMap
}

func dataSourceShareTargetCollectionSubnetDeletedToMap(deletedItem vpcbetav1.SubnetReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareTargetCollectionTargetsVpcToMap(vpcItem vpcbetav1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetCollectionVpcDeletedToMap(*vpcItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcMap["deleted"] = deletedList
	}
	if vpcItem.Href != nil {
		vpcMap["href"] = vpcItem.Href
	}
	if vpcItem.ID != nil {
		vpcMap["id"] = vpcItem.ID
	}
	if vpcItem.Name != nil {
		vpcMap["name"] = vpcItem.Name
	}

	return vpcMap
}

func dataSourceShareTargetCollectionVpcDeletedToMap(deletedItem vpcbetav1.VPCReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
