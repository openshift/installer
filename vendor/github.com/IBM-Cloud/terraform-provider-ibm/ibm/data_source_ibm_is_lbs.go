// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcclassicv1"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	loadBalancers      = "load_balancers"
	CRN                = "crn"
	CreatedAt          = "created_at"
	isLbProfile        = "profile"
	ProvisioningStatus = "provisioning_status"
)

func dataSourceIBMISLBS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMISLBSRead,
		Schema: map[string]*schema.Schema{
			loadBalancers: {
				Type:        schema.TypeList,
				Description: "Collection of load balancers",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						CRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The load balancer's CRN",
						},
						CreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this pool was created.",
						},
						ProvisioningStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this load balancer",
						},
						isLBName: {
							Type:        schema.TypeString,
							Computed:    true,
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
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer subnets list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subnet's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer subnet",
									},
									name: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this load balancer subnet",
									},
									CRN: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet",
									},
								},
							},
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

						isLBListeners: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer Listeners list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									href: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The listener's canonical URL.",
									},
									ID: {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this load balancer listener",
									},
								},
							},
						},
						isLbProfile: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The profile to use for this load balancer",
						},

						isLBPools: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Load Balancer Pools list",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceIBMISLBSRead(d *schema.ResourceData, meta interface{}) error {
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return err
	}
	if userDetails.generation == 1 {
		err := classiclbs(d, meta)
		if err != nil {
			return err
		}
		fmt.Println("classics")
	} else {
		err := getLbs(d, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func classiclbs(d *schema.ResourceData, meta interface{}) error {
	sess, err := classicVpcClient(meta)
	if err != nil {
		return err
	}
	listLoadBalancersOptions := &vpcclassicv1.ListLoadBalancersOptions{}
	lbs, response, err := sess.ListLoadBalancers(listLoadBalancersOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Load Balancers %s\n%s", err, response)
	}

	lbList := make([]map[string]interface{}, 0)

	for _, lb := range lbs.LoadBalancers {
		lbInfo := make(map[string]interface{})
		//	log.Printf("******* lb ******** : (%+v)", lb)
		lbInfo[ID] = *lb.ID
		lbInfo[isLBName] = *lb.Name
		lbInfo[CRN] = *lb.CRN
		lbInfo[ProvisioningStatus] = *lb.ProvisioningStatus
		if *lb.IsPublic {
			lbInfo[isLBType] = "public"
		} else {
			lbInfo[isLBType] = "private"
		}
		lbInfo[isLBStatus] = *lb.ProvisioningStatus
		lbInfo[isLBOperatingStatus] = *lb.OperatingStatus
		publicIpList := make([]string, 0)
		if lb.PublicIps != nil {
			for _, ip := range lb.PublicIps {
				if ip.Address != nil {
					pubip := *ip.Address
					publicIpList = append(publicIpList, pubip)
				}
			}
		}

		lbInfo[isLBPublicIPs] = publicIpList
		privateIpList := make([]string, 0)
		if lb.PrivateIps != nil {
			for _, ip := range lb.PrivateIps {
				if ip.Address != nil {
					prip := *ip.Address
					privateIpList = append(privateIpList, prip)
				}
			}
		}
		lbInfo[isLBPrivateIPs] = privateIpList
		//log.Printf("*******isLBPrivateIPs %+v", lbInfo[isLBPrivateIPs])

		if lb.Subnets != nil {
			subnetList := make([]map[string]interface{}, 0)
			for _, subnet := range lb.Subnets {
				//log.Printf("*******subnet %+v", subnet)
				sub := make(map[string]interface{})
				sub[ID] = *subnet.ID
				sub[href] = *subnet.Href
				if subnet.CRN != nil {
					sub[CRN] = *subnet.CRN
				}
				sub[name] = *subnet.Name
				subnetList = append(subnetList, sub)

			}
			lbInfo[isLBSubnets] = subnetList
			//log.Printf("*******isLBSubnets %+v", lbInfo[isLBSubnets])

		}
		if lb.Listeners != nil {
			listenerList := make([]map[string]interface{}, 0)
			for _, listener := range lb.Listeners {
				lis := make(map[string]interface{})
				lis[ID] = *listener.ID
				lis[href] = *listener.Href
				listenerList = append(listenerList, lis)
			}
			lbInfo[isLBListeners] = listenerList
		}
		//log.Printf("*******isLBListeners %+v", lbInfo[isLBListeners])

		if lb.Pools != nil {
			poolList := make([]map[string]interface{}, 0)

			for _, p := range lb.Pools {
				pool := make(map[string]interface{})
				pool[name] = *p.Name
				pool[ID] = *p.ID
				pool[href] = *p.Href
				poolList = append(poolList, pool)

			}
			lbInfo[isLBPools] = poolList
		}
		lbInfo[isLBResourceGroup] = *lb.ResourceGroup.ID
		lbInfo[isLBHostName] = *lb.Hostname
		tags, err := GetTagsUsingCRN(meta, *lb.CRN)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
		}
		lbInfo[isLBTags] = tags
		//log.Printf("*******tags %+v", tags)

		controller, err := getBaseController(meta)
		if err != nil {
			return err
		}
		lbInfo[ResourceControllerURL] = controller + "/vpc-ext/network/loadBalancers"
		lbInfo[ResourceName] = *lb.Name
		//log.Printf("*******lbInfo %+v", lbInfo)

		if lb.ResourceGroup != nil {
			lbInfo[ResourceGroupName] = *lb.ResourceGroup.ID
		}
		lbList = append(lbList, lbInfo)
		//	log.Printf("*******lbList %+v", lbList)

	}
	d.SetId(dataSourceIBMISLBsID(d))
	d.Set(loadBalancers, lbList)

	return nil
}

func getLbs(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		return err
	}
	listLoadBalancersOptions := &vpcv1.ListLoadBalancersOptions{}
	lbs, response, err := sess.ListLoadBalancers(listLoadBalancersOptions)
	if err != nil {
		return fmt.Errorf("Error Fetching Load Balancers %s\n%s", err, response)
	}
	lbList := make([]map[string]interface{}, 0)

	for _, lb := range lbs.LoadBalancers {
		lbInfo := make(map[string]interface{})
		//	log.Printf("******* lb ******** : (%+v)", lb)
		lbInfo[ID] = *lb.ID
		lbInfo[isLBName] = *lb.Name
		lbInfo[CRN] = *lb.CRN
		lbInfo[ProvisioningStatus] = *lb.ProvisioningStatus

		lbInfo[CreatedAt] = lb.CreatedAt.String()
		if *lb.IsPublic {
			lbInfo[isLBType] = "public"
		} else {
			lbInfo[isLBType] = "private"
		}
		lbInfo[isLBStatus] = *lb.ProvisioningStatus
		lbInfo[isLBOperatingStatus] = *lb.OperatingStatus
		publicIpList := make([]string, 0)
		if lb.PublicIps != nil {
			for _, ip := range lb.PublicIps {
				if ip.Address != nil {
					pubip := *ip.Address
					publicIpList = append(publicIpList, pubip)
				}
			}
		}

		lbInfo[isLBPublicIPs] = publicIpList
		privateIpList := make([]string, 0)
		if lb.PrivateIps != nil {
			for _, ip := range lb.PrivateIps {
				if ip.Address != nil {
					prip := *ip.Address
					privateIpList = append(privateIpList, prip)
				}
			}
		}
		lbInfo[isLBPrivateIPs] = privateIpList
		//log.Printf("*******isLBPrivateIPs %+v", lbInfo[isLBPrivateIPs])

		if lb.Subnets != nil {
			subnetList := make([]map[string]interface{}, 0)
			for _, subnet := range lb.Subnets {
				//	log.Printf("*******subnet %+v", subnet)
				sub := make(map[string]interface{})
				sub[ID] = *subnet.ID
				sub[href] = *subnet.Href
				if subnet.CRN != nil {
					sub[CRN] = *subnet.CRN
				}
				sub[name] = *subnet.Name
				subnetList = append(subnetList, sub)

			}
			lbInfo[isLBSubnets] = subnetList
			//	log.Printf("*******isLBSubnets %+v", lbInfo[isLBSubnets])

		}
		if lb.Listeners != nil {
			listenerList := make([]map[string]interface{}, 0)
			for _, listener := range lb.Listeners {
				lis := make(map[string]interface{})
				lis[ID] = *listener.ID
				lis[href] = *listener.Href
				listenerList = append(listenerList, lis)
			}
			lbInfo[isLBListeners] = listenerList
		}
		//log.Printf("*******isLBListeners %+v", lbInfo[isLBListeners])
		if lb.Pools != nil {
			poolList := make([]map[string]interface{}, 0)

			for _, p := range lb.Pools {
				pool := make(map[string]interface{})
				pool[name] = *p.Name
				pool[ID] = *p.ID
				pool[href] = *p.Href
				poolList = append(poolList, pool)

			}
			lbInfo[isLBPools] = poolList
		}
		if lb.Profile != nil {
			lbProfile := make(map[string]interface{})
			lbProfile[name] = *lb.Profile.Name
			lbProfile[href] = *lb.Profile.Href
			lbProfile["family"] = *lb.Profile.Family
			lbInfo[isLbProfile] = lbProfile
		}
		lbInfo[isLBResourceGroup] = *lb.ResourceGroup.ID
		lbInfo[isLBHostName] = *lb.Hostname
		tags, err := GetTagsUsingCRN(meta, *lb.CRN)
		if err != nil {
			log.Printf(
				"Error on get of resource vpc Load Balancer (%s) tags: %s", d.Id(), err)
		}
		lbInfo[isLBTags] = tags
		controller, err := getBaseController(meta)
		if err != nil {
			return err
		}
		lbInfo[ResourceControllerURL] = controller + "/vpc-ext/network/loadBalancers"
		lbInfo[ResourceName] = *lb.Name
		//log.Printf("*******lbInfo %+v", lbInfo)

		if lb.ResourceGroup != nil {
			lbInfo[ResourceGroupName] = *lb.ResourceGroup.ID
		}
		lbList = append(lbList, lbInfo)
		//	log.Printf("*******lbList %+v", lbList)

	}
	//log.Printf("*******lbList %+v", lbList)
	d.SetId(dataSourceIBMISLBsID(d))
	d.Set(loadBalancers, lbList)
	return nil
}

// dataSourceIBMISLBsID returns a reasonable ID for a transit gateways list.
func dataSourceIBMISLBsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
