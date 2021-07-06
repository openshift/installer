// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"strconv"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	name                    = "name"
	poolAlgorithm           = "algorithm"
	href                    = "href"
	poolProtocol            = "protocol"
	poolCreatedAt           = "created_at"
	poolProvisioningStatus  = "provisioning_status"
	healthMonitor           = "health_monitor"
	instanceGroup           = "instance_group"
	members                 = "members"
	sessionPersistence      = "session_persistence"
	crnInstance             = "crn"
	sessionType             = "type"
	healthMonitorType       = "type"
	healthMonitorDelay      = "delay"
	healthMonitorMaxRetries = "max_retries"
	healthMonitorPort       = "port"
	healthMonitorTimeout    = "timeout"
	healthMonitorURLPath    = "url_path"
)

func dataSourceIBMISLB() *schema.Resource {
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

			isLBStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load Balancer status",
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
				Set:         resourceIBMVPCHash,
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
			ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this instance",
			},

			ResourceName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the resource",
			},

			ResourceGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name in which resource is provisioned",
			},
		},
	}
}

func dataSourceIBMISLBRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	name := d.Get(isLBName).(string)
	if userDetails.generation == 1 {
		err := classiclbGetbyName(d, meta, name)
		if err != nil {
			return err
		}
	} else {
		err := lbGetByName(d, meta, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func classiclbGetbyName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	listLoadBalancersOptions := &vpcclassicv1.ListLoadBalancersOptions{}
	lbs, response, err := sess.ListLoadBalancers(listLoadBalancersOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Load Balancers %s\n%s", err, response)
	}
	for _, lb := range lbs.LoadBalancers {
		if *lb.Name == name {
			d.SetId(*lb.ID)
			d.Set(isLBName, *lb.Name)
			if *lb.IsPublic {
				d.Set(isLBType, "public")
			} else {
				d.Set(isLBType, "private")
			}
			d.Set(isLBStatus, *lb.ProvisioningStatus)
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
			if lb.PrivateIps != nil {
				for _, ip := range lb.PrivateIps {
					if ip.Address != nil {
						prip := *ip.Address
						privateIpList = append(privateIpList, prip)
					}
				}
			}
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
			listLoadBalancerPoolsOptions := &vpcclassicv1.ListLoadBalancerPoolsOptions{}
			listLoadBalancerPoolsOptions.SetLoadBalancerID(*lb.ID)
			poolsResult, _, _ := sess.ListLoadBalancerPools(listLoadBalancerPoolsOptions)
			if poolsResult != nil {
				poolsInfo := make([]map[string]interface{}, 0)
				for _, p := range poolsResult.Pools {
					//log.Printf("******* p ******** : (%+v)", p)
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
					poolsInfo = append(poolsInfo, pool)
				} //for
				d.Set(isLBPools, poolsInfo)

			}
			d.Set(isLBResourceGroup, *lb.ResourceGroup.ID)
			d.Set(isLBHostName, *lb.Hostname)
			tags, err := GetTagsUsingCRN(meta, *lb.CRN)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
			}
			d.Set(isLBTags, tags)
			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc/network/loadBalancers")
			d.Set(ResourceName, *lb.Name)
			if lb.ResourceGroup != nil {
				d.Set(ResourceGroupName, *lb.ResourceGroup.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("No Load balancer found with name %s", name)
}

func lbGetByName(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
	lbs, response, err := sess.ListLoadBalancers(listLoadBalancersOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Load Balancers %s\n%s", err, response)
	}
	for _, lb := range lbs.LoadBalancers {
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
			if lb.PrivateIps != nil {
				for _, ip := range lb.PrivateIps {
					if ip.Address != nil {
						prip := *ip.Address
						privateIpList = append(privateIpList, prip)
					}
				}
			}
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
			tags, err := GetTagsUsingCRN(meta, *lb.CRN)
			if err != nil {
				log.Printf(
					"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
			}
			d.Set(isLBTags, tags)
			controller, err := getBaseController(meta)
			if err != nil {
				return err
			}
			d.Set(ResourceControllerURL, controller+"/vpc-ext/network/loadBalancers")
			d.Set(ResourceName, *lb.Name)
			if lb.ResourceGroup != nil {
				d.Set(ResourceGroupName, *lb.ResourceGroup.ID)
			}
			return nil
		}
	}
	return fmt.Errorf("No Load balancer found with name %s", name)
}
