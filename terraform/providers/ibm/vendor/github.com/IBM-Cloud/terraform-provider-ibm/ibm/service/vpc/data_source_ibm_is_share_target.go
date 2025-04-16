// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-beta-go-sdk/vpcbetav1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIbmIsShareTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext:        dataSourceIbmIsShareTargetRead,
		DeprecationMessage: "This resource is deprecated and will be removed in a future release. Please use ibm_is_share_mount_target instead",

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The file share identifier.",
			},
			"share_target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The share target identifier.",
			},
			"share_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"share", "share_name"},
				Description:  "The file share name.",
			},
			"share_target_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"share_target", "share_target_name"},
				Description:  "The share target name.",
			},
			"transit_encryption": {
				Type:     schema.TypeString,
				Computed: true,
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user-defined name for this share target.",
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
	}
}

func dataSourceIbmIsShareTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1BetaAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	share_id := d.Get("share").(string)
	share_name := d.Get("share_name").(string)
	share_target := d.Get("share_target").(string)
	share_target_name := d.Get("share_target_name").(string)
	var shareTarget *vpcbetav1.ShareMountTarget
	if share_name != "" {
		listSharesOptions := &vpcbetav1.ListSharesOptions{}
		listSharesOptions.Name = &share_name
		shareCollection, response, err := vpcClient.ListSharesWithContext(context, listSharesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListSharesWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		for _, sharesItem := range shareCollection.Shares {
			if *sharesItem.Name == share_name {
				share_id = *sharesItem.ID
				break
			}
		}
	}
	if share_target_name != "" {
		listShareTargetsOptions := &vpcbetav1.ListShareMountTargetsOptions{}

		listShareTargetsOptions.SetShareID(share_id)
		listShareTargetsOptions.SetName(share_target_name)

		shareTargetCollection, response, err := vpcClient.ListShareMountTargetsWithContext(context, listShareTargetsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListShareTargetsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		for _, targetsItem := range shareTargetCollection.MountTargets {
			if *targetsItem.Name == share_target_name {
				shareTarget = &targetsItem
				break
			}
		}
	} else {
		getShareTargetOptions := &vpcbetav1.GetShareMountTargetOptions{}
		getShareTargetOptions.SetShareID(share_id)
		getShareTargetOptions.SetID(share_target)
		shareTarget1, response, err := vpcClient.GetShareMountTargetWithContext(context, getShareTargetOptions)
		if err != nil {
			log.Printf("[DEBUG] GetShareTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		shareTarget = shareTarget1
	}

	d.SetId(fmt.Sprintf("%s/%s", share_id, *shareTarget.ID))
	if err = d.Set("created_at", shareTarget.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", shareTarget.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", shareTarget.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("mount_path", shareTarget.MountPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_path: %s", err))
	}
	if err = d.Set("name", shareTarget.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if shareTarget.TransitEncryption != nil {
		if err = d.Set("transit_encryption", *shareTarget.TransitEncryption); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting transit_encryption: %s", err))
		}
	}

	if err = d.Set("resource_type", shareTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if shareTarget.Subnet != nil {
		err = d.Set("subnet", dataSourceShareMountTargetFlattenSubnet(*shareTarget.Subnet))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subnet %s", err))
		}
	}

	if shareTarget.VPC != nil {
		err = d.Set("vpc", dataSourceShareMountTargetFlattenVpc(*shareTarget.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	return nil
}
