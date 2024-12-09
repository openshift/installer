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

func DataSourceIBMISLBListener() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenerRead,

		Schema: map[string]*schema.Schema{
			isLBListenerLBID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			isLBListenerID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The listener identifier.",
			},
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
	}
}

func dataSourceIBMIsLbListenerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getLoadBalancerListenerOptions := &vpcv1.GetLoadBalancerListenerOptions{}

	getLoadBalancerListenerOptions.SetLoadBalancerID(d.Get(isLBListenerLBID).(string))
	getLoadBalancerListenerOptions.SetID(d.Get(isLBListenerID).(string))

	loadBalancerListener, response, err := vpcClient.GetLoadBalancerListenerWithContext(context, getLoadBalancerListenerOptions)
	if err != nil {
		log.Printf("[DEBUG] GetLoadBalancerListenerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetLoadBalancerListenerWithContext failed %s\n%s", err, response))
	}

	d.SetId(*loadBalancerListener.ID)
	if err = d.Set("accept_proxy_protocol", loadBalancerListener.AcceptProxyProtocol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting accept_proxy_protocol: %s", err))
	}

	if loadBalancerListener.CertificateInstance != nil {
		err = d.Set("certificate_instance", dataSourceLoadBalancerListenerFlattenCertificateInstance(*loadBalancerListener.CertificateInstance))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting certificate_instance %s", err))
		}
	}
	if loadBalancerListener.ConnectionLimit != nil {
		if err = d.Set("connection_limit", flex.IntValue(loadBalancerListener.ConnectionLimit)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting connection_limit: %s", err))
		}
	}
	if loadBalancerListener.IdleConnectionTimeout != nil {
		if err = d.Set(isLBListenerIdleConnectionTimeout, flex.IntValue(loadBalancerListener.IdleConnectionTimeout)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting idle_connection_timeout: %s", err))
		}
	}
	if err = d.Set("created_at", loadBalancerListener.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if loadBalancerListener.DefaultPool != nil {
		err = d.Set("default_pool", dataSourceLoadBalancerListenerFlattenDefaultPool(*loadBalancerListener.DefaultPool))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting default_pool %s", err))
		}
	}
	if err = d.Set("href", loadBalancerListener.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if loadBalancerListener.HTTPSRedirect != nil {
		err = d.Set("https_redirect", dataSourceLoadBalancerListenerFlattenHTTPSRedirect(*loadBalancerListener.HTTPSRedirect))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting https_redirect %s", err))
		}
	}

	if loadBalancerListener.Policies != nil {
		err = d.Set("policies", dataSourceLoadBalancerListenerFlattenPolicies(loadBalancerListener.Policies))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting policies %s", err))
		}
	}
	if err = d.Set("port", flex.IntValue(loadBalancerListener.Port)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting port: %s", err))
	}
	if err = d.Set("port_max", flex.IntValue(loadBalancerListener.PortMax)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting port_max: %s", err))
	}
	if err = d.Set("port_min", flex.IntValue(loadBalancerListener.PortMin)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting port_min: %s", err))
	}
	if err = d.Set("protocol", loadBalancerListener.Protocol); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting protocol: %s", err))
	}
	if err = d.Set("provisioning_status", loadBalancerListener.ProvisioningStatus); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting provisioning_status: %s", err))
	}

	return nil
}

func dataSourceLoadBalancerListenerFlattenCertificateInstance(result vpcv1.CertificateInstanceReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerListenerCertificateInstanceToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerListenerCertificateInstanceToMap(certificateInstanceItem vpcv1.CertificateInstanceReference) (certificateInstanceMap map[string]interface{}) {
	certificateInstanceMap = map[string]interface{}{}

	if certificateInstanceItem.CRN != nil {
		certificateInstanceMap["crn"] = certificateInstanceItem.CRN
	}

	return certificateInstanceMap
}

func dataSourceLoadBalancerListenerFlattenDefaultPool(result vpcv1.LoadBalancerPoolReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerListenerDefaultPoolToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerListenerDefaultPoolToMap(defaultPoolItem vpcv1.LoadBalancerPoolReference) (defaultPoolMap map[string]interface{}) {
	defaultPoolMap = map[string]interface{}{}

	if defaultPoolItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerDefaultPoolDeletedToMap(*defaultPoolItem.Deleted)
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

func dataSourceLoadBalancerListenerDefaultPoolDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerFlattenHTTPSRedirect(result vpcv1.LoadBalancerListenerHTTPSRedirect) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceLoadBalancerListenerHTTPSRedirectToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceLoadBalancerListenerHTTPSRedirectToMap(httpsRedirectItem vpcv1.LoadBalancerListenerHTTPSRedirect) (httpsRedirectMap map[string]interface{}) {
	httpsRedirectMap = map[string]interface{}{}

	if httpsRedirectItem.HTTPStatusCode != nil {
		httpsRedirectMap["http_status_code"] = httpsRedirectItem.HTTPStatusCode
	}
	if httpsRedirectItem.Listener != nil {
		listenerList := []map[string]interface{}{}
		listenerMap := dataSourceLoadBalancerListenerHTTPSRedirectListenerToMap(*httpsRedirectItem.Listener)
		listenerList = append(listenerList, listenerMap)
		httpsRedirectMap["listener"] = listenerList
	}
	if httpsRedirectItem.URI != nil {
		httpsRedirectMap["uri"] = httpsRedirectItem.URI
	}

	return httpsRedirectMap
}

func dataSourceLoadBalancerListenerHTTPSRedirectListenerToMap(listenerItem vpcv1.LoadBalancerListenerReference) (listenerMap map[string]interface{}) {
	listenerMap = map[string]interface{}{}

	if listenerItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerListenerDeletedToMap(*listenerItem.Deleted)
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

func dataSourceLoadBalancerListenerListenerDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceLoadBalancerListenerFlattenPolicies(result []vpcv1.LoadBalancerListenerPolicyReference) (policies []map[string]interface{}) {
	for _, policiesItem := range result {
		policies = append(policies, dataSourceLoadBalancerListenerPoliciesToMap(policiesItem))
	}

	return policies
}

func dataSourceLoadBalancerListenerPoliciesToMap(policiesItem vpcv1.LoadBalancerListenerPolicyReference) (policiesMap map[string]interface{}) {
	policiesMap = map[string]interface{}{}

	if policiesItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceLoadBalancerListenerPoliciesDeletedToMap(*policiesItem.Deleted)
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

func dataSourceLoadBalancerListenerPoliciesDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
