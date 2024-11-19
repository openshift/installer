// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISLBPoolMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbPoolMembersRead,

		Schema: map[string]*schema.Schema{
			"lb": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			"pool": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pool identifier.",
			},
			"members": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of members.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this member was created.",
						},
						"health": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health of the server member in the pool.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The member's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this load balancer pool member.",
						},
						"port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number of the application running in the server member.",
						},
						"provisioning_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this member.",
						},
						"target": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The pool member target. Load balancers in the `network` family support virtual serverinstances. Load balancers in the `application` family support IP addresses.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this virtual server instance.",
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
										Description: "The URL for this virtual server instance.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this virtual server instance.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this virtual server instance (and default system hostname).",
									},
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
								},
							},
						},
						"weight": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Weight of the server member. Applicable only if the pool algorithm is`weighted_round_robin`.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsLbPoolMembersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listLoadBalancerPoolMembersOptions := &vpcv1.ListLoadBalancerPoolMembersOptions{}

	listLoadBalancerPoolMembersOptions.SetLoadBalancerID(d.Get("lb").(string))
	listLoadBalancerPoolMembersOptions.SetPoolID(d.Get("pool").(string))

	loadBalancerPoolMemberCollection, response, err := sess.ListLoadBalancerPoolMembersWithContext(context, listLoadBalancerPoolMembersOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLoadBalancerPoolMembersWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLoadBalancerPoolMembersWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsLbPoolMembersID(d))

	if loadBalancerPoolMemberCollection.Members != nil {
		err = d.Set("members", dataSourceLoadBalancerPoolMemberCollectionFlattenMembers(loadBalancerPoolMemberCollection.Members))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting members %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsLbPoolMembersID returns a reasonable ID for the list.
func dataSourceIBMIsLbPoolMembersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLoadBalancerPoolMemberCollectionFlattenMembers(result []vpcv1.LoadBalancerPoolMember) (members []map[string]interface{}) {
	for _, membersItem := range result {
		members = append(members, dataSourceLoadBalancerPoolMemberCollectionMembersToMap(membersItem))
	}

	return members
}

func dataSourceLoadBalancerPoolMemberCollectionMembersToMap(membersItem vpcv1.LoadBalancerPoolMember) (membersMap map[string]interface{}) {
	membersMap = map[string]interface{}{}

	if membersItem.CreatedAt != nil {
		membersMap["created_at"] = membersItem.CreatedAt.String()
	}
	if membersItem.Health != nil {
		membersMap["health"] = membersItem.Health
	}
	if membersItem.Href != nil {
		membersMap["href"] = membersItem.Href
	}
	if membersItem.ID != nil {
		membersMap["id"] = membersItem.ID
	}
	if membersItem.Port != nil {
		membersMap["port"] = membersItem.Port
	}
	if membersItem.ProvisioningStatus != nil {
		membersMap["provisioning_status"] = membersItem.ProvisioningStatus
	}
	if membersItem.Target != nil {
		targetList := []map[string]interface{}{}
		target := membersItem.Target.(*vpcv1.LoadBalancerPoolMemberTarget)
		targetMap := dataSourceLoadBalancerPoolMemberCollectionMembersTargetToMap(*target)
		targetList = append(targetList, targetMap)
		membersMap["target"] = targetList
	}
	if membersItem.Weight != nil {
		membersMap["weight"] = membersItem.Weight
	}

	return membersMap
}

func dataSourceLoadBalancerPoolMemberCollectionMembersTargetToMap(targetItem vpcv1.LoadBalancerPoolMemberTarget) (targetMap map[string]interface{}) {
	targetMap = map[string]interface{}{}

	if targetItem.CRN != nil {
		targetMap["crn"] = targetItem.CRN
	}
	if targetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerPoolMemberCollectionTargetDeletedToMap(*targetItem.Deleted)
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
	if targetItem.Address != nil {
		targetMap["address"] = targetItem.Address
	}

	return targetMap
}

func dataSourceLoadBalancerPoolMemberCollectionTargetDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
