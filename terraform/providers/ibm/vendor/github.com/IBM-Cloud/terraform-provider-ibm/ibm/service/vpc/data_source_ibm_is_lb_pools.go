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

func DataSourceIBMISLBPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbPoolsRead,

		Schema: map[string]*schema.Schema{
			"lb": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			"pools": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of pools.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancing algorithm.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this pool was created.",
						},
						"health_monitor": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The health monitor of this pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"delay": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The health check interval in seconds. Interval must be greater than timeout value.",
									},
									"max_retries": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The health check max retries.",
									},
									"port": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The health check port number. If specified, this overrides the ports specified in the server member resources.",
									},
									"timeout": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The health check timeout in seconds.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The protocol type of this load balancer pool health monitor.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the health monitor on which the unexpected property value was encountered.",
									},
									"url_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The health check URL path. Applicable only if the health monitor `type` is `http` or`https`. This value must be in the format of an [origin-form request target](https://tools.ietf.org/html/rfc7230#section-5.3.1).",
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
						"instance_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The instance group that is managing this pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this instance group.",
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
										Description: "The URL for this instance group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this instance group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this instance group.",
									},
								},
							},
						},
						"members": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The backend server members of the pool.",
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
										Description: "The member's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool member.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this load balancer pool.",
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol used for this load balancer pool.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the pool on which the unexpected property value was encountered.",
						},
						"provisioning_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this pool.",
						},
						"proxy_protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The PROXY protocol setting for this pool:- `v1`: Enabled with version 1 (human-readable header format)- `v2`: Enabled with version 2 (binary header format)- `disabled`: DisabledSupported by load balancers in the `application` family (otherwise always `disabled`).",
						},
						"session_persistence": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The session persistence of this pool.The enumerated values for this property are expected to expand in the future. Whenprocessing this property, check for and log unknown values. Optionally haltprocessing and surface the error, or bypass the pool on which the unexpectedproperty value was encountered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cookie_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The session persistence cookie name. Applicable only for type `app_cookie`. Names starting with `IBM` are not allowed.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The session persistence type. The `http_cookie` and `app_cookie` types are applicable only to the `http` and `https` protocols.",
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

func dataSourceIBMIsLbPoolsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listLoadBalancerPoolsOptions := &vpcv1.ListLoadBalancerPoolsOptions{}

	listLoadBalancerPoolsOptions.SetLoadBalancerID(d.Get("lb").(string))

	loadBalancerPoolCollection, response, err := sess.ListLoadBalancerPoolsWithContext(context, listLoadBalancerPoolsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLoadBalancerPoolsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLoadBalancerPoolsWithContext failed %s\n%s", err, response))
	}
	if err = d.Set("lb", d.Get("lb").(string)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lb: %s", err))
	}

	d.SetId(dataSourceIBMIsLbPoolsID(d))

	if loadBalancerPoolCollection.Pools != nil {
		err = d.Set("pools", dataSourceLoadBalancerPoolCollectionFlattenPools(loadBalancerPoolCollection.Pools))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting pools %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsLbPoolsID returns a reasonable ID for the list.
func dataSourceIBMIsLbPoolsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLoadBalancerPoolCollectionFlattenPools(result []vpcv1.LoadBalancerPool) (pools []map[string]interface{}) {
	for _, poolsItem := range result {
		pools = append(pools, dataSourceLoadBalancerPoolCollectionPoolsToMap(poolsItem))
	}

	return pools
}

func dataSourceLoadBalancerPoolCollectionPoolsToMap(poolsItem vpcv1.LoadBalancerPool) (poolsMap map[string]interface{}) {
	poolsMap = map[string]interface{}{}

	if poolsItem.Algorithm != nil {
		poolsMap["algorithm"] = poolsItem.Algorithm
	}
	if poolsItem.CreatedAt != nil {
		poolsMap["created_at"] = poolsItem.CreatedAt.String()
	}
	if poolsItem.HealthMonitor != nil {
		healthMonitorList := []map[string]interface{}{}
		healthMonitorMap := dataSourceLoadBalancerPoolCollectionPoolsHealthMonitorToMap(*poolsItem.HealthMonitor)
		healthMonitorList = append(healthMonitorList, healthMonitorMap)
		poolsMap["health_monitor"] = healthMonitorList
	}
	if poolsItem.Href != nil {
		poolsMap["href"] = poolsItem.Href
	}
	if poolsItem.ID != nil {
		poolsMap["id"] = poolsItem.ID
	}
	if poolsItem.InstanceGroup != nil {
		instanceGroupList := []map[string]interface{}{}
		instanceGroupMap := dataSourceLoadBalancerPoolCollectionPoolsInstanceGroupToMap(*poolsItem.InstanceGroup)
		instanceGroupList = append(instanceGroupList, instanceGroupMap)
		poolsMap["instance_group"] = instanceGroupList
	}
	if poolsItem.Members != nil {
		membersList := []map[string]interface{}{}
		for _, membersItem := range poolsItem.Members {
			membersList = append(membersList, dataSourceLoadBalancerPoolCollectionPoolsMembersToMap(membersItem))
		}
		poolsMap["members"] = membersList
	}
	if poolsItem.Name != nil {
		poolsMap["name"] = poolsItem.Name
	}
	if poolsItem.Protocol != nil {
		poolsMap["protocol"] = poolsItem.Protocol
	}
	if poolsItem.ProvisioningStatus != nil {
		poolsMap["provisioning_status"] = poolsItem.ProvisioningStatus
	}
	if poolsItem.ProxyProtocol != nil {
		poolsMap["proxy_protocol"] = poolsItem.ProxyProtocol
	}
	if poolsItem.SessionPersistence != nil {
		sessionPersistenceList := []map[string]interface{}{}
		sessionPersistenceMap := dataSourceLoadBalancerPoolCollectionPoolsSessionPersistenceToMap(*poolsItem.SessionPersistence)
		sessionPersistenceList = append(sessionPersistenceList, sessionPersistenceMap)
		poolsMap["session_persistence"] = sessionPersistenceList
	}

	return poolsMap
}

func dataSourceLoadBalancerPoolCollectionPoolsHealthMonitorToMap(healthMonitorItem vpcv1.LoadBalancerPoolHealthMonitor) (healthMonitorMap map[string]interface{}) {
	healthMonitorMap = map[string]interface{}{}

	if healthMonitorItem.Delay != nil {
		healthMonitorMap["delay"] = healthMonitorItem.Delay
	}
	if healthMonitorItem.MaxRetries != nil {
		healthMonitorMap["max_retries"] = healthMonitorItem.MaxRetries
	}
	if healthMonitorItem.Port != nil {
		healthMonitorMap["port"] = healthMonitorItem.Port
	}
	if healthMonitorItem.Timeout != nil {
		healthMonitorMap["timeout"] = healthMonitorItem.Timeout
	}
	if healthMonitorItem.Type != nil {
		healthMonitorMap["type"] = healthMonitorItem.Type
	}
	if healthMonitorItem.URLPath != nil {
		healthMonitorMap["url_path"] = healthMonitorItem.URLPath
	}

	return healthMonitorMap
}

func dataSourceLoadBalancerPoolCollectionPoolsInstanceGroupToMap(instanceGroupItem vpcv1.InstanceGroupReference) (instanceGroupMap map[string]interface{}) {
	instanceGroupMap = map[string]interface{}{}

	if instanceGroupItem.CRN != nil {
		instanceGroupMap["crn"] = instanceGroupItem.CRN
	}
	if instanceGroupItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerPoolCollectionInstanceGroupDeletedToMap(*instanceGroupItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		instanceGroupMap["deleted"] = deletedList
	}
	if instanceGroupItem.Href != nil {
		instanceGroupMap["href"] = instanceGroupItem.Href
	}
	if instanceGroupItem.ID != nil {
		instanceGroupMap["id"] = instanceGroupItem.ID
	}
	if instanceGroupItem.Name != nil {
		instanceGroupMap["name"] = instanceGroupItem.Name
	}

	return instanceGroupMap
}

func dataSourceLoadBalancerPoolCollectionInstanceGroupDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerPoolCollectionPoolsMembersToMap(membersItem vpcv1.LoadBalancerPoolMemberReference) (membersMap map[string]interface{}) {
	membersMap = map[string]interface{}{}

	if membersItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerPoolCollectionMembersDeletedToMap(*membersItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		membersMap["deleted"] = deletedList
	}
	if membersItem.Href != nil {
		membersMap["href"] = membersItem.Href
	}
	if membersItem.ID != nil {
		membersMap["id"] = membersItem.ID
	}

	return membersMap
}

func dataSourceLoadBalancerPoolCollectionMembersDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerPoolCollectionPoolsSessionPersistenceToMap(sessionPersistenceItem vpcv1.LoadBalancerPoolSessionPersistence) (sessionPersistenceMap map[string]interface{}) {
	sessionPersistenceMap = map[string]interface{}{}

	if sessionPersistenceItem.CookieName != nil {
		sessionPersistenceMap["cookie_name"] = sessionPersistenceItem.CookieName
	}
	if sessionPersistenceItem.Type != nil {
		sessionPersistenceMap["type"] = sessionPersistenceItem.Type
	}

	return sessionPersistenceMap
}
