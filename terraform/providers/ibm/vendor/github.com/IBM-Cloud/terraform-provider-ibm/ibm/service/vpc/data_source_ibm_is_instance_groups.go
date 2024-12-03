// Copyright IBM Corp. 2022 All Rights Reserved.
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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISInstanceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMIsInstanceGroupsRead,

		Schema: map[string]*schema.Schema{
			"instance_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of instance groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Required if specifying a load balancer pool only. Used by the instance group when scaling up instances to supply the port for the load balancer pool member.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the instance group was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this instance group.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this instance group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this instance group.",
						},
						"instance_template": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The template used to create new instances for this group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this instance template.",
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
										Description: "The URL for this instance template.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance template.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this instance template.",
									},
								},
							},
						},
						"load_balancer_pool": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The load balancer pool managed by this group. Instances createdby this group will have a new load balancer pool member in thatpool created.",
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
										Description: "The pool's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this load balancer pool.",
									},
								},
							},
						},
						"managers": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The managers for the instance group.",
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
										Description: "The URL for this instance group manager.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance group manager.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this instance group manager.",
									},
								},
							},
						},
						"membership_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of instances in the instance group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this instance group.",
						},
						"resource_group": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
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
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the instance group- `deleting`: Group is being deleted- `healthy`: Group has `membership_count` instances- `scaling`: Instances in the group are being created or deleted to reach             `membership_count`- `unhealthy`: Group is unable to reach `membership_count` instances.",
						},
						"subnets": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnets to use when creating new instances.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
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
										Description: "The URL for this subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the instance group was updated.",
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC the instance group resides in.",
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
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						isInstanceGroupAccessTags: {
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

func DataSourceIBMIsInstanceGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.InstanceGroup{}
	listInstanceGroupsOptions := &vpcv1.ListInstanceGroupsOptions{}
	for {

		if start != "" {
			listInstanceGroupsOptions.Start = &start
		}

		instanceGroupCollection, response, err := vpcClient.ListInstanceGroupsWithContext(context, listInstanceGroupsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListInstanceGroupsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListInstanceGroupsWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(instanceGroupCollection.Next)
		allrecs = append(allrecs, instanceGroupCollection.InstanceGroups...)
		if start == "" {
			break
		}
	}

	d.SetId(DataSourceIBMIsInstanceGroupsID(d))

	instanceGroups := []map[string]interface{}{}

	for _, instanceGroupItem := range allrecs {
		instanceGroup, err := DataSourceIBMIsInstanceGroupsInstanceGroupToMap(&instanceGroupItem, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		instanceGroups = append(instanceGroups, instanceGroup)
	}

	if err = d.Set("instance_groups", instanceGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_groups %s", err))
	}

	return nil
}

// DataSourceIBMIsInstanceGroupsID returns a reasonable ID for the list.
func DataSourceIBMIsInstanceGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsInstanceGroupsInstanceGroupCollectionFirstToMap(model *vpcv1.InstanceGroupCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsInstanceGroupToMap(model *vpcv1.InstanceGroup, meta interface{}) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApplicationPort != nil {
		modelMap["application_port"] = *model.ApplicationPort
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.InstanceTemplate != nil {
		instanceTemplateMap, err := DataSourceIBMIsInstanceGroupsInstanceTemplateReferenceToMap(model.InstanceTemplate)
		if err != nil {
			return modelMap, err
		}
		modelMap["instance_template"] = []map[string]interface{}{instanceTemplateMap}
	}
	if model.LoadBalancerPool != nil {
		loadBalancerPoolMap, err := DataSourceIBMIsInstanceGroupsLoadBalancerPoolReferenceToMap(model.LoadBalancerPool)
		if err != nil {
			return modelMap, err
		}
		modelMap["load_balancer_pool"] = []map[string]interface{}{loadBalancerPoolMap}
	}
	if model.Managers != nil {
		managers := []map[string]interface{}{}
		for _, managersItem := range model.Managers {
			managersItemMap, err := DataSourceIBMIsInstanceGroupsInstanceGroupManagerReferenceToMap(&managersItem)
			if err != nil {
				return modelMap, err
			}
			managers = append(managers, managersItemMap)
		}
		modelMap["managers"] = managers
	}
	if model.MembershipCount != nil {
		modelMap["membership_count"] = *model.MembershipCount
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroup != nil {
		resourceGroupMap, err := DataSourceIBMIsInstanceGroupsResourceGroupReferenceToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Subnets != nil {
		subnets := []map[string]interface{}{}
		for _, subnetsItem := range model.Subnets {
			subnetsItemMap, err := DataSourceIBMIsInstanceGroupsSubnetReferenceToMap(&subnetsItem)
			if err != nil {
				return modelMap, err
			}
			subnets = append(subnets, subnetsItemMap)
		}
		modelMap["subnets"] = subnets
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VPC != nil {
		vpcMap, err := DataSourceIBMIsInstanceGroupsVPCReferenceToMap(model.VPC)
		if err != nil {
			return modelMap, err
		}
		modelMap["vpc"] = []map[string]interface{}{vpcMap}
	}
	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *model.CRN, "", isInstanceGroupAccessTagType)
	if err != nil {
		log.Printf(
			"[ERROR] Error on get of resource instance group (%s) access tags: %s", *model.ID, err)
	}
	modelMap[isInstanceGroupAccessTags] = accesstags
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsInstanceTemplateReferenceToMap(model *vpcv1.InstanceTemplateReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceGroupsInstanceTemplateReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsInstanceTemplateReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsLoadBalancerPoolReferenceToMap(model *vpcv1.LoadBalancerPoolReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceGroupsLoadBalancerPoolReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsLoadBalancerPoolReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsInstanceGroupManagerReferenceToMap(model *vpcv1.InstanceGroupManagerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceGroupsInstanceGroupManagerReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsInstanceGroupManagerReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceGroupsSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsSubnetReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CRN != nil {
		modelMap["crn"] = *model.CRN
	}
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsInstanceGroupsVPCReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Href != nil {
		modelMap["href"] = *model.Href
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	return modelMap, nil
}

func DataSourceIBMIsInstanceGroupsVPCReferenceDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
