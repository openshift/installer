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

func DataSourceIBMISLBListenerPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenerPoliciesRead,

		Schema: map[string]*schema.Schema{
			isLBListenerPolicyLBID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			isLBListenerPolicyListenerID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The listener identifier.",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy action.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the policy on which the unexpected property value was encountered.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this policy was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The listener policy's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy's unique identifier.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this policy.",
						},
						"priority": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Priority of the policy. Lower value indicates higher priority.",
						},
						"provisioning_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this policy.",
						},
						"rules": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The rules for this policy.",
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
										Description: "The rule's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The rule's unique identifier.",
									},
								},
							},
						},
						"target": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "- If `action` is `forward`, the response is a `LoadBalancerPoolReference`- If `action` is `redirect`, the response is a `LoadBalancerListenerPolicyRedirectURL`- If `action` is `https_redirect`, the response is a `LoadBalancerListenerHTTPSRedirect`.",
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
									"http_status_code": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The HTTP status code for this redirect.",
									},
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The redirect target URL.",
									},
									"listener": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
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
													Description: "The listener's canonical URL.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this load balancer listener.",
												},
											},
										},
									},
									"uri": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The redirect relative target URI.",
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

func dataSourceIBMIsLbListenerPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listLoadBalancerListenerPoliciesOptions := &vpcv1.ListLoadBalancerListenerPoliciesOptions{}

	listLoadBalancerListenerPoliciesOptions.SetLoadBalancerID(d.Get(isLBListenerPolicyLBID).(string))
	listLoadBalancerListenerPoliciesOptions.SetListenerID(d.Get(isLBListenerPolicyListenerID).(string))

	loadBalancerListenerPolicyCollection, response, err := vpcClient.ListLoadBalancerListenerPoliciesWithContext(context, listLoadBalancerListenerPoliciesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLoadBalancerListenerPoliciesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLoadBalancerListenerPoliciesWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsLbListenerPoliciesID(d))

	if loadBalancerListenerPolicyCollection.Policies != nil {
		err = d.Set("policies", dataSourceLoadBalancerListenerPolicyCollectionFlattenPolicies(loadBalancerListenerPolicyCollection.Policies))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting policies %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsLbListenerPoliciesID returns a reasonable ID for the list.
func dataSourceIBMIsLbListenerPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLoadBalancerListenerPolicyCollectionFlattenPolicies(result []vpcv1.LoadBalancerListenerPolicy) (policies []map[string]interface{}) {
	for _, policiesItem := range result {
		policies = append(policies, dataSourceLoadBalancerListenerPolicyCollectionPoliciesToMap(policiesItem))
	}

	return policies
}

func dataSourceLoadBalancerListenerPolicyCollectionPoliciesToMap(policiesItem vpcv1.LoadBalancerListenerPolicy) (policiesMap map[string]interface{}) {
	policiesMap = map[string]interface{}{}

	if policiesItem.Action != nil {
		policiesMap["action"] = policiesItem.Action
	}
	if policiesItem.CreatedAt != nil {
		policiesMap["created_at"] = policiesItem.CreatedAt.String()
	}
	if policiesItem.Href != nil {
		policiesMap["href"] = policiesItem.Href
	}
	if policiesItem.ID != nil {
		policiesMap["id"] = policiesItem.ID
	}
	if policiesItem.Name != nil {
		policiesMap["name"] = policiesItem.Name
	}
	if policiesItem.Priority != nil {
		policiesMap["priority"] = policiesItem.Priority
	}
	if policiesItem.ProvisioningStatus != nil {
		policiesMap["provisioning_status"] = policiesItem.ProvisioningStatus
	}
	if policiesItem.Rules != nil {
		rulesList := []map[string]interface{}{}
		for _, rulesItem := range policiesItem.Rules {
			rulesList = append(rulesList, dataSourceLoadBalancerListenerPolicyCollectionPoliciesRulesToMap(rulesItem))
		}
		policiesMap["rules"] = rulesList
	}
	if policiesItem.Target != nil {
		targetList := []map[string]interface{}{}
		target := policiesItem.Target.(*vpcv1.LoadBalancerListenerPolicyTarget)
		targetMap := dataSourceLoadBalancerListenerPolicyCollectionPoliciesTargetToMap(*target)
		targetList = append(targetList, targetMap)
		policiesMap["target"] = targetList
	}

	return policiesMap
}

func dataSourceLoadBalancerListenerPolicyCollectionPoliciesRulesToMap(rulesItem vpcv1.LoadBalancerListenerPolicyRuleReference) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	if rulesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyCollectionRulesDeletedToMap(*rulesItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		rulesMap["deleted"] = deletedList
	}
	if rulesItem.Href != nil {
		rulesMap["href"] = rulesItem.Href
	}
	if rulesItem.ID != nil {
		rulesMap["id"] = rulesItem.ID
	}

	return rulesMap
}

func dataSourceLoadBalancerListenerPolicyCollectionRulesDeletedToMap(deletedItem vpcv1.LoadBalancerListenerPolicyRuleReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerPolicyCollectionPoliciesTargetToMap(targetItem vpcv1.LoadBalancerListenerPolicyTarget) (targetMap map[string]interface{}) {
	targetMap = map[string]interface{}{}

	if targetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyCollectionTargetDeletedToMap(*targetItem.Deleted)
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
	if targetItem.HTTPStatusCode != nil {
		targetMap["http_status_code"] = targetItem.HTTPStatusCode
	}
	if targetItem.URL != nil {
		targetMap["url"] = targetItem.URL
	}
	if targetItem.Listener != nil {
		listenerList := []map[string]interface{}{}
		listenerMap := dataSourceLoadBalancerListenerPolicyCollectionTargetListenerToMap(*targetItem.Listener)
		listenerList = append(listenerList, listenerMap)
		targetMap["listener"] = listenerList
	}
	if targetItem.URI != nil {
		targetMap["uri"] = targetItem.URI
	}

	return targetMap
}

func dataSourceLoadBalancerListenerPolicyCollectionTargetDeletedToMap(deletedItem vpcv1.LoadBalancerPoolReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerPolicyCollectionTargetListenerToMap(listenerItem vpcv1.LoadBalancerListenerReference) (listenerMap map[string]interface{}) {
	listenerMap = map[string]interface{}{}

	if listenerItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyCollectionListenerDeletedToMap(*listenerItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		listenerMap["deleted"] = deletedList
	}
	if listenerItem.Href != nil {
		listenerMap["href"] = listenerItem.Href
	}
	if listenerItem.ID != nil {
		listenerMap["id"] = listenerItem.ID
	}

	return listenerMap
}

func dataSourceLoadBalancerListenerPolicyCollectionListenerDeletedToMap(deletedItem vpcv1.LoadBalancerListenerReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
