// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsFlowLog() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsFlowLogRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"identifier", "name"},
				Description:  "The unique user-defined name for this flow log collector.",
			},

			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"identifier", "name"},
				Description:  "The flow log collector identifier.",
			},
			"active": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this collector is active.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to `true`, this flow log collector will be automatically deleted when the target is deleted.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the flow log collector was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this flow log collector.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this flow log collector.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the flow log collector.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this flow log collector.",
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
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			"storage_bucket": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Cloud Object Storage bucket where the collected flows are logged.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name of this COS bucket.",
						},
					},
				},
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target this collector is collecting flow logs for. If the target is an instance,subnet, or VPC, flow logs will not be collected for any network interfaces within thetarget that are themselves the target of a more specific flow log collector.",
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
							Description: "The URL for this network interface.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this network interface.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this network interface.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual server instance.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this flow log collector is associated with.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
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
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this VPC.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsFlowLogRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	identifier := d.Get("identifier").(string)
	var flowLogCollector *vpcv1.FlowLogCollector

	if name != "" {
		start := ""
		allrecs := []vpcv1.FlowLogCollector{}
		for {
			listOptions := &vpcv1.ListFlowLogCollectorsOptions{}
			if start != "" {
				listOptions.Start = &start
			}
			flowlogCollectors, response, err := sess.ListFlowLogCollectors(listOptions)
			if err != nil {
				return diag.FromErr(fmt.Errorf("Error Fetching Flow Logs for VPC %s\n%s", err, response))
			}
			start = GetNext(flowlogCollectors.Next)
			allrecs = append(allrecs, flowlogCollectors.FlowLogCollectors...)
			if start == "" {
				break
			}
		}
		for _, flowlogCollector := range allrecs {
			if *flowlogCollector.Name == name {
				flowLogCollector = &flowlogCollector
				break
			}
		}
		if flowLogCollector == nil {
			return diag.FromErr(fmt.Errorf("No flow log collector found with name (%s)", name))
		}
	} else if identifier != "" {
		getFlowLogCollectorOptions := &vpcv1.GetFlowLogCollectorOptions{}

		getFlowLogCollectorOptions.SetID(d.Get("identifier").(string))

		flowlogCollector, response, err := sess.GetFlowLogCollectorWithContext(context, getFlowLogCollectorOptions)
		if err != nil {
			log.Printf("[DEBUG] GetFlowLogCollectorWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("GetFlowLogCollectorWithContext failed %s\n%s", err, response))
		}
		flowLogCollector = flowlogCollector
	}

	d.SetId(*flowLogCollector.ID)
	if err = d.Set("active", flowLogCollector.Active); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting active: %s", err))
	}
	if err = d.Set("auto_delete", flowLogCollector.AutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting auto_delete: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(flowLogCollector.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", flowLogCollector.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", flowLogCollector.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", flowLogCollector.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", flowLogCollector.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("identifier", *flowLogCollector.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting identifier: %s", err))
	}

	if flowLogCollector.ResourceGroup != nil {
		err = d.Set("resource_group", dataSourceFlowLogCollectorFlattenResourceGroup(*flowLogCollector.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}
	}

	if flowLogCollector.StorageBucket != nil {
		err = d.Set("storage_bucket", dataSourceFlowLogCollectorFlattenStorageBucket(*flowLogCollector.StorageBucket))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting storage_bucket %s", err))
		}
	}

	if flowLogCollector.Target != nil {
		targetIntf := flowLogCollector.Target
		target := targetIntf.(*vpcv1.FlowLogCollectorTarget)
		err = d.Set("target", dataSourceFlowLogCollectorFlattenTarget(*target))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting target %s", err))
		}
	}

	if flowLogCollector.VPC != nil {
		err = d.Set("vpc", dataSourceFlowLogCollectorFlattenVPC(*flowLogCollector.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
		}
	}

	return nil
}

func dataSourceFlowLogCollectorFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFlowLogCollectorResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFlowLogCollectorResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceFlowLogCollectorFlattenStorageBucket(result vpcv1.CloudObjectStorageBucketReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFlowLogCollectorStorageBucketToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFlowLogCollectorStorageBucketToMap(storageBucketItem vpcv1.CloudObjectStorageBucketReference) (storageBucketMap map[string]interface{}) {
	storageBucketMap = map[string]interface{}{}

	if storageBucketItem.Name != nil {
		storageBucketMap["name"] = storageBucketItem.Name
	}

	return storageBucketMap
}

func dataSourceFlowLogCollectorFlattenTarget(result vpcv1.FlowLogCollectorTarget) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFlowLogCollectorTargetToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFlowLogCollectorTargetToMap(targetItem vpcv1.FlowLogCollectorTarget) (targetMap map[string]interface{}) {
	targetMap = map[string]interface{}{}

	if targetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceFlowLogCollectorTargetDeletedToMap(*targetItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		targetMap["deleted"] = deletedList
	}
	if targetItem.Href != nil {
		targetMap["href"] = targetItem.Href
	}
	if targetItem.ID != nil {
		targetMap["id"] = targetItem.ID
	}
	if targetItem.Name != nil {
		targetMap["name"] = targetItem.Name
	}
	if targetItem.ResourceType != nil {
		targetMap["resource_type"] = targetItem.ResourceType
	}
	if targetItem.CRN != nil {
		targetMap["crn"] = targetItem.CRN
	}

	return targetMap
}

func dataSourceFlowLogCollectorTargetDeletedToMap(deletedItem vpcv1.NetworkInterfaceReferenceTargetContextDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceFlowLogCollectorFlattenVPC(result vpcv1.VPCReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceFlowLogCollectorVPCToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceFlowLogCollectorVPCToMap(vpcItem vpcv1.VPCReference) (vpcMap map[string]interface{}) {
	vpcMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceFlowLogCollectorVPCDeletedToMap(*vpcItem.Deleted)
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

func dataSourceFlowLogCollectorVPCDeletedToMap(deletedItem vpcv1.VPCReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
