// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsShareTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsShareTargetRead,

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The file share identifier.",
			},
			"mount_target": {
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
			"mount_target_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"mount_target", "mount_target_name"},
				Description:  "The share target name.",
			},
			"transit_encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_control_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The access control mode for the share",
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
			"primary_ip": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The primary IP address of the virtual network interface for the share mount target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address..",
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
							Description: "The URL for this reserved IP.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this reserved IP.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
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
	}
}

func dataSourceIBMIsShareTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	share_id := d.Get("share").(string)
	share_name := d.Get("share_name").(string)
	share_target := d.Get("mount_target").(string)
	share_target_name := d.Get("mount_target_name").(string)
	var shareTarget *vpcv1.ShareMountTarget
	if share_name != "" {
		listSharesOptions := &vpcv1.ListSharesOptions{}
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
		listShareTargetsOptions := &vpcv1.ListShareMountTargetsOptions{}

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
		getShareTargetOptions := &vpcv1.GetShareMountTargetOptions{}
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
	if shareTarget.AccessControlMode != nil {
		d.Set("access_control_mode", *shareTarget.AccessControlMode)
	}
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
	if err = d.Set("resource_type", shareTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if shareTarget.TransitEncryption != nil {
		if err = d.Set("transit_encryption", *shareTarget.TransitEncryption); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting transit_encryption: %s", err))
		}
	}

	if shareTarget.PrimaryIP != nil {
		err = d.Set("primary_ip", dataSourceShareMountTargetFlattenPrimaryIP(*shareTarget.PrimaryIP))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	if shareTarget.VPC != nil {
		err = d.Set("vpc", dataSourceShareMountTargetFlattenVpc(*shareTarget.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	if shareTarget.VirtualNetworkInterface != nil {
		err = d.Set("virtual_network_interface", dataSourceShareMountTargetFlattenVNI(*shareTarget.VirtualNetworkInterface))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	if shareTarget.Subnet != nil {
		err = d.Set("subnet", dataSourceShareMountTargetFlattenSubnet(*shareTarget.Subnet))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subnet %s", err))
		}
	}

	return nil
}

func dataSourceShareMountTargetFlattenVpc(result vpcv1.VPCReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareMountTargetVpcToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareMountTargetVpcToMap(vpcItem vpcv1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareMountTargetVpcDeletedToMap(*vpcItem.Deleted)
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

	if vpcItem.ResourceType != nil {
		vpcMap["resource_type"] = vpcItem.ResourceType
	}

	return vpcMap
}

func dataSourceShareMountTargetVpcDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareMountTargetFlattenPrimaryIP(result vpcv1.ReservedIPReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareTargetPrimaryIPToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareTargetPrimaryIPToMap(primaryIPItem vpcv1.ReservedIPReference) (primaryIPMap map[string]interface{}) {
	primaryIPMap = map[string]interface{}{}

	if primaryIPItem.Address != nil {
		primaryIPMap["address"] = primaryIPItem.Address
	}
	if primaryIPItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetPrimaryIPDeletedToMap(*primaryIPItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		primaryIPMap["deleted"] = deletedList
	}
	if primaryIPItem.Href != nil {
		primaryIPMap["href"] = primaryIPItem.Href
	}
	if primaryIPItem.ID != nil {
		primaryIPMap["id"] = primaryIPItem.ID
	}
	if primaryIPItem.Name != nil {
		primaryIPMap["name"] = primaryIPItem.Name
	}

	if primaryIPItem.ResourceType != nil {
		primaryIPMap["resource_type"] = primaryIPItem.ResourceType
	}

	return primaryIPMap
}

func dataSourceShareTargetPrimaryIPDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareMountTargetFlattenSubnet(result vpcv1.SubnetReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareTargetSubnetToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareTargetSubnetToMap(subnetItem vpcv1.SubnetReference) (subnetMap map[string]interface{}) {
	subnetMap = map[string]interface{}{}

	if subnetItem.CRN != nil {
		subnetMap["crn"] = subnetItem.CRN
	}
	if subnetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetSubnetDeletedToMap(*subnetItem.Deleted)
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

	if subnetItem.ResourceType != nil {
		subnetMap["resource_type"] = subnetItem.ResourceType
	}

	return subnetMap
}

func dataSourceShareTargetSubnetDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceShareMountTargetFlattenVNI(result vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareTargetVNIToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareTargetVNIToMap(VNIItem vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext) (subnetMap map[string]interface{}) {
	subnetMap = map[string]interface{}{}

	if VNIItem.CRN != nil {
		subnetMap["crn"] = VNIItem.CRN
	}
	if VNIItem.Href != nil {
		subnetMap["href"] = VNIItem.Href
	}
	if VNIItem.ID != nil {
		subnetMap["id"] = VNIItem.ID
	}
	if VNIItem.Name != nil {
		subnetMap["name"] = VNIItem.Name
	}

	if VNIItem.ResourceType != nil {
		subnetMap["resource_type"] = VNIItem.ResourceType
	}

	return subnetMap
}
