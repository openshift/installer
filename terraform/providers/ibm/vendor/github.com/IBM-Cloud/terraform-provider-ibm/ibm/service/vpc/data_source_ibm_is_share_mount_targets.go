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

func DataSourceIBMIsShareTargets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareTargetsRead,

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
			"mount_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of share targets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The access control mode for the share",
						},
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
						"virtual_network_interface": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The virtual network interface for this file share mount target.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this virtual network interface.",
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
										Description: "The URL for this virtual network interface.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this virtual network interface.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this virtual network interface.",
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

func dataSourceIBMIsShareTargetsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		err = d.Set("mount_targets", dataSourceShareMountTargetCollectionFlattenTargets(shareTargetCollection.MountTargets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsShareTargetsID returns a reasonable ID for the list.
func dataSourceIBMIsShareTargetsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareMountTargetCollectionFlattenTargets(result []vpcbetav1.ShareMountTarget) (targets []map[string]interface{}) {
	for _, targetsItem := range result {
		targets = append(targets, dataSourceShareMountTargetCollectionTargetsToMap(targetsItem))
	}

	return targets
}

func dataSourceShareMountTargetCollectionTargetsToMap(targetsItem vpcbetav1.ShareMountTarget) (targetsMap map[string]interface{}) {
	targetsMap = map[string]interface{}{}

	if targetsItem.AccessControlMode != nil {
		targetsMap["access_control_mode"] = *targetsItem.AccessControlMode
	}
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
	if targetsItem.ResourceType != nil {
		targetsMap["resource_type"] = targetsItem.ResourceType
	}
	if targetsItem.TransitEncryption != nil {
		targetsMap["transit_encryption"] = *targetsItem.TransitEncryption
	}

	if targetsItem.VPC != nil {
		targetsMap["vpc"] = dataSourceShareMountTargetFlattenVpc(*targetsItem.VPC)
	}

	if targetsItem.VirtualNetworkInterface != nil {
		targetsMap["virtual_network_interface"] = dataSourceShareMountTargetFlattenVNI(*targetsItem.VirtualNetworkInterface)
	}

	if targetsItem.Subnet != nil {
		targetsMap["subnet"] = dataSourceShareMountTargetFlattenSubnet(*targetsItem.Subnet)
	}
	return targetsMap
}
