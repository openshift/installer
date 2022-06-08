// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISLBListenerPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenerPolicyRead,

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
			isLBListenerPolicyID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy identifier.",
			},
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
	}
}

func dataSourceIBMIsLbListenerPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getLoadBalancerListenerPolicyOptions := &vpcv1.GetLoadBalancerListenerPolicyOptions{}

	getLoadBalancerListenerPolicyOptions.SetLoadBalancerID(d.Get(isLBListenerPolicyLBID).(string))
	getLoadBalancerListenerPolicyOptions.SetListenerID(d.Get(isLBListenerPolicyListenerID).(string))
	getLoadBalancerListenerPolicyOptions.SetID(d.Get(isLBListenerPolicyID).(string))

	loadBalancerListenerPolicy, response, err := vpcClient.GetLoadBalancerListenerPolicyWithContext(context, getLoadBalancerListenerPolicyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLoadBalancerListenerPolicyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLoadBalancerListenerPolicyWithContext failed %s\n%s", err, response))
	}
	d.SetId(*loadBalancerListenerPolicy.ID)
	if err = d.Set("action", loadBalancerListenerPolicy.Action); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting action: %s", err))
	}

	if err = d.Set("created_at", loadBalancerListenerPolicy.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", loadBalancerListenerPolicy.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("name", loadBalancerListenerPolicy.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("priority", flex.IntValue(loadBalancerListenerPolicy.Priority)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting priority: %s", err))
	}
	if err = d.Set("provisioning_status", loadBalancerListenerPolicy.ProvisioningStatus); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting provisioning_status: %s", err))
	}

	if loadBalancerListenerPolicy.Rules != nil {
		err = d.Set("rules", dataSourceLoadBalancerListenerPolicyFlattenRules(loadBalancerListenerPolicy.Rules))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rules %s", err))
		}
	}

	if loadBalancerListenerPolicy.Target != nil {
		target := loadBalancerListenerPolicy.Target.(*vpcv1.LoadBalancerListenerPolicyTarget)
		err = d.Set("target", dataSourceLoadBalancerListenerPolicyFlattenTarget(*target))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting target %s", err))
		}
	}

	return nil
}

func dataSourceLoadBalancerListenerPolicyFlattenRules(result []vpcv1.LoadBalancerListenerPolicyRuleReference) (rules []map[string]interface{}) {
	for _, rulesItem := range result {
		rules = append(rules, dataSourceLoadBalancerListenerPolicyRulesToMap(rulesItem))
	}

	return rules
}

func dataSourceLoadBalancerListenerPolicyRulesToMap(rulesItem vpcv1.LoadBalancerListenerPolicyRuleReference) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	if rulesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyRulesDeletedToMap(*rulesItem.Deleted)
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

func dataSourceLoadBalancerListenerPolicyRulesDeletedToMap(deletedItem vpcv1.LoadBalancerListenerPolicyRuleReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerPolicyFlattenTarget(result vpcv1.LoadBalancerListenerPolicyTarget) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerListenerPolicyTargetToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerListenerPolicyTargetToMap(targetItem vpcv1.LoadBalancerListenerPolicyTarget) (targetMap map[string]interface{}) {
	targetMap = map[string]interface{}{}

	if targetItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyTargetDeletedToMap(*targetItem.Deleted)
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
		listenerMap := dataSourceLoadBalancerListenerPolicyTargetListenerToMap(*targetItem.Listener)
		listenerList = append(listenerList, listenerMap)
		targetMap["listener"] = listenerList
	}
	if targetItem.URI != nil {
		targetMap["uri"] = targetItem.URI
	}

	return targetMap
}

func dataSourceLoadBalancerListenerPolicyTargetDeletedToMap(deletedItem vpcv1.LoadBalancerPoolReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerPolicyTargetListenerToMap(listenerItem vpcv1.LoadBalancerListenerReference) (listenerMap map[string]interface{}) {
	listenerMap = map[string]interface{}{}

	if listenerItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPolicyListenerDeletedToMap(*listenerItem.Deleted)
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

func dataSourceLoadBalancerListenerPolicyListenerDeletedToMap(deletedItem vpcv1.LoadBalancerListenerReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
