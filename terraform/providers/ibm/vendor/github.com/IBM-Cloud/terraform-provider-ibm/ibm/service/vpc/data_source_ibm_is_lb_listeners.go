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

func DataSourceIBMISLBListeners() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenersRead,

		Schema: map[string]*schema.Schema{
			isLBListenerLBID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The load balancer identifier.",
			},
			"listeners": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of listeners.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accept_proxy_protocol": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, this listener will accept and forward PROXY protocol information. Supported by load balancers in the `application` family (otherwise always `false`). Additional restrictions:- If this listener has `https_redirect` specified, its `accept_proxy_protocol` value must  match the `accept_proxy_protocol` value of the `https_redirect` listener.- If this listener is the target of another listener's `https_redirect`, its  `accept_proxy_protocol` value must match that listener's `accept_proxy_protocol` value.",
						},
						"certificate_instance": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The certificate instance used for SSL termination. It is applicable only to `https`protocol.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this certificate instance.",
									},
								},
							},
						},
						"connection_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The connection limit of the listener.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this listener was created.",
						},
						"default_pool": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The default pool associated with the listener.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The listener's canonical URL.",
						},
						"https_redirect": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If specified, the target listener that requests are redirected to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_status_code": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The HTTP status code for this redirect.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this load balancer listener.",
						},
						"policies": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The policies for this listener.",
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
										Description: "The listener policy's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The policy's unique identifier.",
									},
								},
							},
						},
						"port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The listener port number, or the inclusive lower bound of the port range. Each listener in the load balancer must have a unique `port` and `protocol` combination.",
						},
						"port_max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The inclusive upper bound of the range of ports used by this listener.Only load balancers in the `network` family support more than one port per listener.",
						},
						"port_min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The inclusive lower bound of the range of ports used by this listener.Only load balancers in the `network` family support more than one port per listener.",
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The listener protocol. Load balancers in the `network` family support `tcp`. Load balancers in the `application` family support `tcp`, `http`, and `https`. Each listener in the load balancer must have a unique `port` and `protocol` combination.",
						},
						"provisioning_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this listener.",
						},
						isLBListenerIdleConnectionTimeout: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "idle connection timeout of listener",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsLbListenersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listLoadBalancerListenersOptions := &vpcv1.ListLoadBalancerListenersOptions{}

	listLoadBalancerListenersOptions.SetLoadBalancerID(d.Get(isLBListenerLBID).(string))

	loadBalancerListenerCollection, response, err := vpcClient.ListLoadBalancerListenersWithContext(context, listLoadBalancerListenersOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLoadBalancerListenersWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLoadBalancerListenersWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsLbListenersID(d))

	if loadBalancerListenerCollection.Listeners != nil {
		err = d.Set("listeners", dataSourceLoadBalancerListenerCollectionFlattenListeners(loadBalancerListenerCollection.Listeners))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting listeners %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsLbListenersID returns a reasonable ID for the list.
func dataSourceIBMIsLbListenersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLoadBalancerListenerCollectionFlattenListeners(result []vpcv1.LoadBalancerListener) (listeners []map[string]interface{}) {
	for _, listenersItem := range result {
		listeners = append(listeners, dataSourceLoadBalancerListenerCollectionListenersToMap(listenersItem))
	}

	return listeners
}

func dataSourceLoadBalancerListenerCollectionListenersToMap(listenersItem vpcv1.LoadBalancerListener) (listenersMap map[string]interface{}) {
	listenersMap = map[string]interface{}{}

	if listenersItem.AcceptProxyProtocol != nil {
		listenersMap["accept_proxy_protocol"] = listenersItem.AcceptProxyProtocol
	}
	if listenersItem.CertificateInstance != nil {
		certificateInstanceList := []map[string]interface{}{}
		certificateInstanceMap := dataSourceLoadBalancerListenerCollectionListenersCertificateInstanceToMap(*listenersItem.CertificateInstance)
		certificateInstanceList = append(certificateInstanceList, certificateInstanceMap)
		listenersMap["certificate_instance"] = certificateInstanceList
	}
	if listenersItem.ConnectionLimit != nil {
		listenersMap["connection_limit"] = listenersItem.ConnectionLimit
	}
	if listenersItem.IdleConnectionTimeout != nil {
		listenersMap[isLBListenerIdleConnectionTimeout] = listenersItem.IdleConnectionTimeout
	}
	if listenersItem.CreatedAt != nil {
		listenersMap["created_at"] = listenersItem.CreatedAt.String()
	}
	if listenersItem.DefaultPool != nil {
		defaultPoolList := []map[string]interface{}{}
		defaultPoolMap := dataSourceLoadBalancerListenerCollectionListenersDefaultPoolToMap(*listenersItem.DefaultPool)
		defaultPoolList = append(defaultPoolList, defaultPoolMap)
		listenersMap["default_pool"] = defaultPoolList
	}
	if listenersItem.Href != nil {
		listenersMap["href"] = listenersItem.Href
	}
	if listenersItem.HTTPSRedirect != nil {
		httpsRedirectList := []map[string]interface{}{}
		httpsRedirectMap := dataSourceLoadBalancerListenerCollectionListenersHTTPSRedirectToMap(*listenersItem.HTTPSRedirect)
		httpsRedirectList = append(httpsRedirectList, httpsRedirectMap)
		listenersMap["https_redirect"] = httpsRedirectList
	}
	if listenersItem.ID != nil {
		listenersMap["id"] = listenersItem.ID
	}
	if listenersItem.Policies != nil {
		policiesList := []map[string]interface{}{}
		for _, policiesItem := range listenersItem.Policies {
			policiesList = append(policiesList, dataSourceLoadBalancerListenerCollectionListenersPoliciesToMap(policiesItem))
		}
		listenersMap["policies"] = policiesList
	}
	if listenersItem.Port != nil {
		listenersMap["port"] = listenersItem.Port
	}
	if listenersItem.PortMax != nil {
		listenersMap["port_max"] = listenersItem.PortMax
	}
	if listenersItem.PortMin != nil {
		listenersMap["port_min"] = listenersItem.PortMin
	}
	if listenersItem.Protocol != nil {
		listenersMap["protocol"] = listenersItem.Protocol
	}
	if listenersItem.ProvisioningStatus != nil {
		listenersMap["provisioning_status"] = listenersItem.ProvisioningStatus
	}

	return listenersMap
}

func dataSourceLoadBalancerListenerCollectionListenersCertificateInstanceToMap(certificateInstanceItem vpcv1.CertificateInstanceReference) (certificateInstanceMap map[string]interface{}) {
	certificateInstanceMap = map[string]interface{}{}

	if certificateInstanceItem.CRN != nil {
		certificateInstanceMap["crn"] = certificateInstanceItem.CRN
	}

	return certificateInstanceMap
}

func dataSourceLoadBalancerListenerCollectionListenersDefaultPoolToMap(defaultPoolItem vpcv1.LoadBalancerPoolReference) (defaultPoolMap map[string]interface{}) {
	defaultPoolMap = map[string]interface{}{}

	if defaultPoolItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerCollectionDefaultPoolDeletedToMap(*defaultPoolItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		defaultPoolMap["deleted"] = deletedList
	}
	if defaultPoolItem.Href != nil {
		defaultPoolMap["href"] = defaultPoolItem.Href
	}
	if defaultPoolItem.ID != nil {
		defaultPoolMap["id"] = defaultPoolItem.ID
	}
	if defaultPoolItem.Name != nil {
		defaultPoolMap["name"] = defaultPoolItem.Name
	}

	return defaultPoolMap
}

func dataSourceLoadBalancerListenerCollectionDefaultPoolDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerCollectionListenersHTTPSRedirectToMap(httpsRedirectItem vpcv1.LoadBalancerListenerHTTPSRedirect) (httpsRedirectMap map[string]interface{}) {
	httpsRedirectMap = map[string]interface{}{}

	if httpsRedirectItem.HTTPStatusCode != nil {
		httpsRedirectMap["http_status_code"] = httpsRedirectItem.HTTPStatusCode
	}
	if httpsRedirectItem.Listener != nil {
		listenerList := []map[string]interface{}{}
		listenerMap := dataSourceLoadBalancerListenerCollectionHTTPSRedirectListenerToMap(*httpsRedirectItem.Listener)
		listenerList = append(listenerList, listenerMap)
		httpsRedirectMap["listener"] = listenerList
	}
	if httpsRedirectItem.URI != nil {
		httpsRedirectMap["uri"] = httpsRedirectItem.URI
	}

	return httpsRedirectMap
}

func dataSourceLoadBalancerListenerCollectionHTTPSRedirectListenerToMap(listenerItem vpcv1.LoadBalancerListenerReference) (listenerMap map[string]interface{}) {
	listenerMap = map[string]interface{}{}

	if listenerItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerCollectionListenerDeletedToMap(*listenerItem.Deleted)
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

func dataSourceLoadBalancerListenerCollectionListenerDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerCollectionListenersPoliciesToMap(policiesItem vpcv1.LoadBalancerListenerPolicyReference) (policiesMap map[string]interface{}) {
	policiesMap = map[string]interface{}{}

	if policiesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerCollectionPoliciesDeletedToMap(*policiesItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		policiesMap["deleted"] = deletedList
	}
	if policiesItem.Href != nil {
		policiesMap["href"] = policiesItem.Href
	}
	if policiesItem.ID != nil {
		policiesMap["id"] = policiesItem.ID
	}

	return policiesMap
}

func dataSourceLoadBalancerListenerCollectionPoliciesDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
