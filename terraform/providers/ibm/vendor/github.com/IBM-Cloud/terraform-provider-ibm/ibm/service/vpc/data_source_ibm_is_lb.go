// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"fmt"
	"log"
	"strconv"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	name                      = "name"
	poolAlgorithm             = "algorithm"
	href                      = "href"
	poolProtocol              = "protocol"
	poolCreatedAt             = "created_at"
	poolProvisioningStatus    = "provisioning_status"
	healthMonitor             = "health_monitor"
	instanceGroup             = "instance_group"
	members                   = "members"
	sessionPersistence        = "session_persistence"
	crnInstance               = "crn"
	sessionType               = "type"
	healthMonitorType         = "type"
	healthMonitorDelay        = "delay"
	healthMonitorMaxRetries   = "max_retries"
	healthMonitorPort         = "port"
	healthMonitorTimeout      = "timeout"
	healthMonitorURLPath      = "url_path"
	isLBPrivateIPDetail       = "private_ip"
	isLBPrivateIpAddress      = "address"
	isLBPrivateIpHref         = "href"
	isLBPrivateIpName         = "name"
	isLBPrivateIpId           = "reserved_ip"
	isLBPrivateIpResourceType = "resource_type"
)

func DataSourceIBMISLB() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISLBRead,

		Schema: map[string]*schema.Schema{
			isLBName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Load Balancer name",
			},

			isLBType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer type",
			},

			isLBUdpSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this load balancer supports UDP.",
			},

			isLBStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer status",
			},

			isLBRouteMode: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether route mode is enabled for this load balancer",
			},

			isLBCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this Load Balancer",
			},

			isLBOperatingStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer operating status",
			},

			isLBPublicIPs: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Load Balancer Public IPs",
			},

			isLBPrivateIPs: {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Load Balancer private IPs",
			},

			isLBSubnets: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Load Balancer subnets list",
			},

			isLBSecurityGroups: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Load Balancer securitygroups list",
			},

			isLBSecurityGroupsSupported: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Security Group Supported for this Load Balancer",
			},

			isLBTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "Tags associated to Load Balancer",
			},

			isLBResourceGroup: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer Resource group",
			},

			isLBHostName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer Host Name",
			},

			isLBLogging: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Logging of Load Balancer",
			},

			isLBListeners: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Load Balancer Listeners list",
			},
			isLBPools: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Load Balancer Pools list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						poolAlgorithm: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancing algorithm.",
						},
						healthMonitor: {
							Description: "The health monitor of this pool.",
							Computed:    true,
							Type:        schema.TypeMap,
						},

						instanceGroup: {
							Description: "The instance group that is managing this pool.",
							Computed:    true,
							Type:        schema.TypeMap,
						},

						members: {
							Description: "The backend server members of the pool.",
							Computed:    true,
							Type:        schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The member's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer pool member.",
									},
								},
							},
						},
						sessionPersistence: {
							Description: "The session persistence of this pool.",
							Computed:    true,
							Type:        schema.TypeMap,
						},
						poolCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this pool was created.",
						},
						href: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The pool's canonical URL.",
						},
						ID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this load balancer pool",
						},
						name: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this load balancer pool",
						},
						poolProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol used for this load balancer pool.",
						},
						poolProvisioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this pool.",
						},
					},
				},
			},
			isLBPrivateIPDetail: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The private IP addresses assigned to this load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isLBPrivateIpAddress: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address to reserve, which must not already be reserved on the subnet.",
						},
						isLBPrivateIpHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP",
						},
						isLBPrivateIpName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this reserved IP. If unspecified, the name will be a hyphenated list of randomly-selected words. Names must be unique within the subnet the reserved IP resides in. ",
						},
						isLBPrivateIpId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Identifies a reserved IP by a unique property.",
						},
						isLBPrivateIpResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
					},
				},
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			flex.ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			flex.ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func dataSourceIBMISLBRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get(isLBName).(string)
	err := lbGetByName(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func lbGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}

	start := ""
	allrecs := []vpcv1.LoadBalancer{}
	for {
		listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
		if start != "" {
			listLoadBalancersOptions.Start = &start
		}
		lbs, response, err := sess.ListLoadBalancers(listLoadBalancersOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Load Balancers %s\n%s", err, response)
		}
		start = flex.GetNext(lbs.Next)
		allrecs = append(allrecs, lbs.LoadBalancers...)
		if start == "" {
			break
		}
	}

	for _, lb := range allrecs {
		if *lb.Name == name {
			d.SetId(*lb.ID)
			d.Set(isLBName, *lb.Name)
			if lb.Logging != nil && lb.Logging.Datapath != nil {
				d.Set(isLBLogging, *lb.Logging.Datapath.Active)
			}
			if *lb.IsPublic {
				d.Set(isLBType, "public")
			} else {
				d.Set(isLBType, "private")
			}
			d.Set(isLBStatus, *lb.ProvisioningStatus)
			if lb.RouteMode != nil {
				d.Set(isLBRouteMode, *lb.RouteMode)
			}
			if lb.UDPSupported != nil {
				d.Set(isLBUdpSupported, *lb.UDPSupported)
			}
			d.Set(isLBCrn, *lb.CRN)
			d.Set(isLBOperatingStatus, *lb.OperatingStatus)
			publicIpList := make([]string, 0)
			if lb.PublicIps != nil {
				for _, ip := range lb.PublicIps {
					if ip.Address != nil {
						pubip := *ip.Address
						publicIpList = append(publicIpList, pubip)
					}
				}
			}
			d.Set(isLBPublicIPs, publicIpList)
			privateIpList := make([]string, 0)
			privateIpDetailList := make([]map[string]interface{}, 0)
			if lb.PrivateIps != nil {
				for _, ip := range lb.PrivateIps {
					if ip.Address != nil {
						prip := *ip.Address
						privateIpList = append(privateIpList, prip)
					}
					currentPriIp := map[string]interface{}{}

					if ip.Address != nil {
						currentPriIp[isLBPrivateIpAddress] = ip.Address
					}
					if ip.Href != nil {
						currentPriIp[isLBPrivateIpHref] = ip.Href
					}
					if ip.Name != nil {
						currentPriIp[isLBPrivateIpName] = ip.Name
					}
					if ip.ID != nil {
						currentPriIp[isLBPrivateIpId] = ip.ID
					}
					if ip.ResourceType != nil {
						currentPriIp[isLBPrivateIpResourceType] = ip.ResourceType
					}
					privateIpDetailList = append(privateIpDetailList, currentPriIp)

				}
			}
			d.Set(isLBPrivateIPDetail, privateIpDetailList)
			d.Set(isLBPrivateIPs, privateIpList)
			if lb.Subnets != nil {
				subnetList := make([]string, 0)
				for _, subnet := range lb.Subnets {
					if subnet.ID != nil {
						sub := *subnet.ID
						subnetList = append(subnetList, sub)
					}
				}
				d.Set(isLBSubnets, subnetList)
			}

			d.Set(isLBSecurityGroupsSupported, false)
			if lb.SecurityGroups != nil {
				securitygroupList := make([]string, 0)
				for _, securityGroup := range lb.SecurityGroups {
					if securityGroup.ID != nil {
						securityGroupID := *securityGroup.ID
						securitygroupList = append(securitygroupList, securityGroupID)
					}
				}
				d.Set(isLBSecurityGroups, securitygroupList)
				d.Set(isLBSecurityGroupsSupported, true)
			}

			if lb.Listeners != nil {
				listenerList := make([]string, 0)
				for _, listener := range lb.Listeners {
					if listener.ID != nil {
						lis := *listener.ID
						listenerList = append(listenerList, lis)
					}
				}
				d.Set(isLBListeners, listenerList)
			}
			listLoadBalancerPoolsOptions := &vpcv1.ListLoadBalancerPoolsOptions{}
			listLoadBalancerPoolsOptions.SetLoadBalancerID(*lb.ID)
			poolsResult, _, _ := sess.ListLoadBalancerPools(listLoadBalancerPoolsOptions)
			if poolsResult != nil {
				poolsInfo := make([]map[string]interface{}, 0)

				for _, p := range poolsResult.Pools {
					//	log.Printf("******* p ******** : (%+v)", p)
					pool := make(map[string]interface{})
					pool[poolAlgorithm] = *p.Algorithm
					pool[ID] = *p.ID
					pool[href] = *p.Href
					pool[poolProtocol] = *p.Protocol
					pool[poolCreatedAt] = p.CreatedAt.String()
					pool[poolProvisioningStatus] = *p.ProvisioningStatus
					pool["name"] = *p.Name
					if p.HealthMonitor != nil {
						healthMonitorInfo := make(map[string]interface{})
						delayfinal := strconv.FormatInt(*(p.HealthMonitor.Delay), 10)
						healthMonitorInfo[healthMonitorDelay] = delayfinal
						maxRetriesfinal := strconv.FormatInt(*(p.HealthMonitor.MaxRetries), 10)
						timeoutfinal := strconv.FormatInt(*(p.HealthMonitor.Timeout), 10)
						healthMonitorInfo[healthMonitorMaxRetries] = maxRetriesfinal
						healthMonitorInfo[healthMonitorTimeout] = timeoutfinal
						if p.HealthMonitor.URLPath != nil {
							healthMonitorInfo[healthMonitorURLPath] = *(p.HealthMonitor.URLPath)
						}
						healthMonitorInfo[healthMonitorType] = *(p.HealthMonitor.Type)
						pool[healthMonitor] = healthMonitorInfo
					}

					if p.SessionPersistence != nil {
						sessionPersistenceInfo := make(map[string]interface{})
						sessionPersistenceInfo[sessionType] = *p.SessionPersistence.Type
						pool[sessionPersistence] = sessionPersistenceInfo
					}
					if p.Members != nil {
						memberList := make([]map[string]interface{}, len(p.Members))
						for j, m := range p.Members {
							member := make(map[string]interface{})
							member[ID] = *m.ID
							member[href] = *m.Href
							memberList[j] = member
						}
						pool[members] = memberList
					}

					if p.InstanceGroup != nil {
						instanceGroupInfo := make(map[string]interface{})
						instanceGroupInfo[ID] = *(p.InstanceGroup.ID)
						instanceGroupInfo[crnInstance] = *(p.InstanceGroup.CRN)
						instanceGroupInfo[href] = *(p.InstanceGroup.Href)
						instanceGroupInfo[name] = *(p.InstanceGroup.Name)
						pool[instanceGroup] = instanceGroupInfo
					}
					poolsInfo = append(poolsInfo, pool)
				} //for
				d.Set(isLBPools, poolsInfo)
			}

			d.Set(isLBResourceGroup, *lb.ResourceGroup.ID)
			d.Set(isLBHostName, *lb.Hostname)
			tags, err := flex.GetTagsUsingCRN(meta, *lb.CRN)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
			}
			d.Set(isLBTags, tags)
			controller, err := flex.GetBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(flex.ResourceControllerURL, controller+"/vpc-ext/network/loadBalancers")
			d.Set(flex.ResourceName, *lb.Name)
			if lb.ResourceGroup != nil {
				d.Set(flex.ResourceGroupName, *lb.ResourceGroup.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("[ERROR] No Load balancer found with name %s", name)
}
