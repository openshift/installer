// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServersRead,

		Schema: map[string]*schema.Schema{
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "resource group identifier.",
			},
			"vpn_servers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPN servers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The certificate instance for this VPN server.",
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
						"client_authentication": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The methods used to authenticate VPN clients to this VPN server. VPN clients must authenticate against all provided methods.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of authentication.",
									},
									"identity_provider": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of identity provider to be used by the VPN client.- `iam`: IBM identity and access managementThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered,The type of identity provider to be used by VPN client.",
									},
									"client_ca": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this certificate instance,The certificate instance used for the VPN client certificate authority (CA).",
									},
								},
							},
						},
						"client_auto_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, disconnected VPN clients will be automatically deleted after the `client_auto_delete_timeout` time has passed.",
						},
						"client_auto_delete_timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Hours after which disconnected VPN clients will be automatically deleted. If `0`, disconnected VPN clients will be deleted immediately.",
						},
						"client_dns_server_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The DNS server addresses that will be provided to VPN clients that are connected to this VPN server.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
								},
							},
						},
						"client_idle_timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The seconds a VPN client can be idle before this VPN server will disconnect it.  If `0`, the server will not disconnect idle clients.",
						},
						"client_ip_pool": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VPN client IPv4 address pool, expressed in CIDR format.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the VPN server was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPN server.",
						},
						"enable_split_tunneling": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the split tunneling is enabled on this VPN server.",
						},
						"health_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource.- `ok`: Healthy- `degraded`: Suffering from compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
						},
						"health_reasons": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						"hostname": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Fully qualified domain name assigned to this VPN server.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPN server.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN server.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the VPN server.",
						},
						"lifecycle_reasons": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current lifecycle_state (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
									},

									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},

									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this VPN server.",
						},
						"port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port number used by this VPN server.",
						},
						"private_ips": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reserved IPs bound to this VPN server.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address. This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
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
										Description: "The URL for this reserved IP.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined or system-provided name for this reserved IP.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"protocol": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The transport protocol used by this VPN server.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this VPN server.",
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
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
						"security_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The security groups targeting this VPN server.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The security group's CRN.",
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
										Description: "The security group's canonical URL.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this security group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
									},
								},
							},
						},
						"subnets": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnets this VPN server is part of.",
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
						"vpc": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The VPC this VPN server resides in.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this VPC.",
									},
									"deleted": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this VPC.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this VPC.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this VPC.",
									},
								},
							},
						},

						isVPNServerAccessTags: {
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

func dataSourceIBMIsVPNServersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceGrp := d.Get("resource_group_id").(string)

	start := ""
	allrecs := []vpcv1.VPNServer{}

	for {
		listVPNServersOptions := &vpcv1.ListVPNServersOptions{}
		if resourceGrp != "" {
			listVPNServersOptions.ResourceGroupID = &resourceGrp
		}

		if start != "" {
			listVPNServersOptions.Start = &start
		}
		vpnServerCollection, response, err := sess.ListVPNServersWithContext(context, listVPNServersOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVPNServersWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] ListVPNServersWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(vpnServerCollection.Next)
		allrecs = append(allrecs, vpnServerCollection.VPNServers...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsVPNServersID(d))

	if allrecs != nil {
		err = d.Set("vpn_servers", dataSourceVPNServerCollectionFlattenVPNServers(allrecs, meta))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting vpn_servers %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsVPNServersID returns a reasonable ID for the list.
func dataSourceIBMIsVPNServersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVPNServerCollectionFlattenFirst(result vpcv1.VPNServerCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerCollectionFirstToMap(firstItem vpcv1.VPNServerCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceVPNServerCollectionFlattenNext(result vpcv1.VPNServerCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerCollectionNextToMap(nextItem vpcv1.VPNServerCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceVPNServerCollectionFlattenVPNServers(result []vpcv1.VPNServer, meta interface{}) (vpnServers []map[string]interface{}) {
	for _, vpnServersItem := range result {
		vpnServers = append(vpnServers, dataSourceVPNServerCollectionVPNServersToMap(vpnServersItem, meta))
	}

	return vpnServers
}

func dataSourceVPNServerCollectionVPNServersToMap(vpnServersItem vpcv1.VPNServer, meta interface{}) (vpnServersMap map[string]interface{}) {
	vpnServersMap = map[string]interface{}{}

	if vpnServersItem.Certificate != nil {
		certificateList := []map[string]interface{}{}
		certificateMap := dataSourceVPNServerCollectionVPNServersCertificateToMap(*vpnServersItem.Certificate)
		certificateList = append(certificateList, certificateMap)
		vpnServersMap["certificate"] = certificateList
	}

	vpnServerAuthenticationPrototypeArray := make([]interface{}, len(vpnServersItem.ClientAuthentication))
	for i, clientAuthenticationItem := range vpnServersItem.ClientAuthentication {
		vpnServerAuthenticationPrototype := make(map[string]interface{})
		vpnServerAuthentication := clientAuthenticationItem.(*vpcv1.VPNServerAuthentication)
		if vpnServerAuthentication != nil {
			if vpnServerAuthentication.Method != nil {
				vpnServerAuthenticationPrototype["method"] = *vpnServerAuthentication.Method
				if vpnServerAuthentication.ClientCa != nil && vpnServerAuthentication.ClientCa.CRN != nil {
					vpnServerAuthenticationPrototype["client_ca"] = *vpnServerAuthentication.ClientCa.CRN
				}
				if vpnServerAuthentication.IdentityProvider != nil {
					vpnServerAuthenticationByUsernameIDProvider := vpnServerAuthentication.IdentityProvider.(*vpcv1.VPNServerAuthenticationByUsernameIDProvider)
					vpnServerAuthenticationPrototype["identity_provider"] = *vpnServerAuthenticationByUsernameIDProvider.ProviderType
				}
			}
		}
		vpnServerAuthenticationPrototypeArray[i] = vpnServerAuthenticationPrototype
	}
	vpnServersMap["client_authentication"] = vpnServerAuthenticationPrototypeArray

	if vpnServersItem.ClientAutoDelete != nil {
		vpnServersMap["client_auto_delete"] = vpnServersItem.ClientAutoDelete
	}
	if vpnServersItem.ClientAutoDeleteTimeout != nil {
		vpnServersMap["client_auto_delete_timeout"] = vpnServersItem.ClientAutoDeleteTimeout
	}
	if vpnServersItem.ClientDnsServerIps != nil {
		clientDnsServerIpsList := []map[string]interface{}{}
		for _, clientDnsServerIpsItem := range vpnServersItem.ClientDnsServerIps {
			clientDnsServerIpsList = append(clientDnsServerIpsList, dataSourceVPNServerCollectionVPNServersClientDnsServerIpsToMap(clientDnsServerIpsItem))
		}
		vpnServersMap["client_dns_server_ips"] = clientDnsServerIpsList
	}
	if vpnServersItem.ClientIdleTimeout != nil {
		vpnServersMap["client_idle_timeout"] = vpnServersItem.ClientIdleTimeout
	}
	if vpnServersItem.ClientIPPool != nil {
		vpnServersMap["client_ip_pool"] = vpnServersItem.ClientIPPool
	}
	if vpnServersItem.CreatedAt != nil {
		vpnServersMap["created_at"] = vpnServersItem.CreatedAt.String()
	}
	if vpnServersItem.CRN != nil {
		vpnServersMap["crn"] = vpnServersItem.CRN
	}
	if vpnServersItem.EnableSplitTunneling != nil {
		vpnServersMap["enable_split_tunneling"] = vpnServersItem.EnableSplitTunneling
	}
	if vpnServersItem.HealthState != nil {
		vpnServersMap["health_state"] = vpnServersItem.HealthState
	}
	if vpnServersItem.HealthReasons != nil {
		vpnServersMap["health_reasons"] = resourceVPNServerFlattenHealthReasons(vpnServersItem.HealthReasons)
	}
	if vpnServersItem.Hostname != nil {
		vpnServersMap["hostname"] = vpnServersItem.Hostname
	}
	if vpnServersItem.Href != nil {
		vpnServersMap["href"] = vpnServersItem.Href
	}
	if vpnServersItem.ID != nil {
		vpnServersMap["id"] = vpnServersItem.ID
	}
	if vpnServersItem.LifecycleState != nil {
		vpnServersMap["lifecycle_state"] = vpnServersItem.LifecycleState
	}
	if vpnServersItem.LifecycleReasons != nil {
		vpnServersMap["lifecycle_reasons"] = resourceVPNServerFlattenLifecycleReasons(vpnServersItem.LifecycleReasons)
	}
	if vpnServersItem.Name != nil {
		vpnServersMap["name"] = vpnServersItem.Name
	}
	if vpnServersItem.Port != nil {
		vpnServersMap["port"] = vpnServersItem.Port
	}
	if vpnServersItem.PrivateIps != nil {
		privateIpsList := []map[string]interface{}{}
		for _, privateIpsItem := range vpnServersItem.PrivateIps {
			privateIpsList = append(privateIpsList, dataSourceVPNServerCollectionVPNServersPrivateIpsToMap(privateIpsItem))
		}
		vpnServersMap["private_ips"] = privateIpsList
	}
	if vpnServersItem.Protocol != nil {
		vpnServersMap["protocol"] = vpnServersItem.Protocol
	}
	if vpnServersItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceVPNServerCollectionVPNServersResourceGroupToMap(*vpnServersItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		vpnServersMap["resource_group"] = resourceGroupList
	}
	if vpnServersItem.ResourceType != nil {
		vpnServersMap["resource_type"] = vpnServersItem.ResourceType
	}
	if vpnServersItem.SecurityGroups != nil {
		securityGroupsList := []map[string]interface{}{}
		for _, securityGroupsItem := range vpnServersItem.SecurityGroups {
			securityGroupsList = append(securityGroupsList, dataSourceVPNServerCollectionVPNServersSecurityGroupsToMap(securityGroupsItem))
		}
		vpnServersMap["security_groups"] = securityGroupsList
	}
	if vpnServersItem.Subnets != nil {
		subnetsList := []map[string]interface{}{}
		for _, subnetsItem := range vpnServersItem.Subnets {
			subnetsList = append(subnetsList, dataSourceVPNServerCollectionVPNServersSubnetsToMap(subnetsItem))
		}
		vpnServersMap["subnets"] = subnetsList
	}
	if vpnServersItem.VPC != nil {
		vpcList := []map[string]interface{}{}
		// for _, vpcsItem := range vpnServersItem.VPC {
		vpcList = append(vpcList, dataSourceVPNServerCollectionVPNServersVpcReferenceToMap(vpnServersItem.VPC))
		// }
		vpnServersMap["vpc"] = vpcList
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpnServersItem.CRN, "", isVPNServerAccessTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of vpn server (%s) access tags: %s", *vpnServersItem.ID, err)
	}
	vpnServersMap[isVPNServerAccessTags] = accesstags

	return vpnServersMap
}

func dataSourceVPNServerCollectionVPNServersCertificateToMap(certificateItem vpcv1.CertificateInstanceReference) (certificateMap map[string]interface{}) {
	certificateMap = map[string]interface{}{}

	if certificateItem.CRN != nil {
		certificateMap["crn"] = certificateItem.CRN
	}

	return certificateMap
}

func dataSourceVPNServerCollectionClientAuthenticationIdentityProviderToMap(identityProviderItem vpcv1.VPNServerAuthenticationByUsernameIDProvider) (identityProviderMap map[string]interface{}) {
	identityProviderMap = map[string]interface{}{}

	if identityProviderItem.ProviderType != nil {
		identityProviderMap["provider_type"] = identityProviderItem.ProviderType
	}

	return identityProviderMap
}

func dataSourceVPNServerCollectionClientAuthenticationClientCaToMap(clientCaItem vpcv1.CertificateInstanceReference) (clientCaMap map[string]interface{}) {
	clientCaMap = map[string]interface{}{}

	if clientCaItem.CRN != nil {
		clientCaMap["crn"] = clientCaItem.CRN
	}

	return clientCaMap
}

func dataSourceVPNServerCollectionVPNServersClientDnsServerIpsToMap(clientDnsServerIpsItem vpcv1.IP) (clientDnsServerIpsMap map[string]interface{}) {
	clientDnsServerIpsMap = map[string]interface{}{}

	if clientDnsServerIpsItem.Address != nil {
		clientDnsServerIpsMap["address"] = clientDnsServerIpsItem.Address
	}

	return clientDnsServerIpsMap
}

func dataSourceVPNServerCollectionVPNServersPrivateIpsToMap(privateIpsItem vpcv1.ReservedIPReference) (privateIpsMap map[string]interface{}) {
	privateIpsMap = map[string]interface{}{}

	if privateIpsItem.Address != nil {
		privateIpsMap["address"] = privateIpsItem.Address
	}
	if privateIpsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerCollectionPrivateIpsDeletedToMap(*privateIpsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		privateIpsMap["deleted"] = deletedList
	}
	if privateIpsItem.Href != nil {
		privateIpsMap["href"] = privateIpsItem.Href
	}
	if privateIpsItem.ID != nil {
		privateIpsMap["id"] = privateIpsItem.ID
	}
	if privateIpsItem.Name != nil {
		privateIpsMap["name"] = privateIpsItem.Name
	}
	if privateIpsItem.ResourceType != nil {
		privateIpsMap["resource_type"] = privateIpsItem.ResourceType
	}

	return privateIpsMap
}

func dataSourceVPNServerCollectionPrivateIpsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNServerCollectionVPNServersResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
}

func dataSourceVPNServerCollectionVPNServersSecurityGroupsToMap(securityGroupsItem vpcv1.SecurityGroupReference) (securityGroupsMap map[string]interface{}) {
	securityGroupsMap = map[string]interface{}{}

	if securityGroupsItem.CRN != nil {
		securityGroupsMap["crn"] = securityGroupsItem.CRN
	}
	if securityGroupsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerCollectionSecurityGroupsDeletedToMap(*securityGroupsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		securityGroupsMap["deleted"] = deletedList
	}
	if securityGroupsItem.Href != nil {
		securityGroupsMap["href"] = securityGroupsItem.Href
	}
	if securityGroupsItem.ID != nil {
		securityGroupsMap["id"] = securityGroupsItem.ID
	}
	if securityGroupsItem.Name != nil {
		securityGroupsMap["name"] = securityGroupsItem.Name
	}

	return securityGroupsMap
}

func dataSourceVPNServerCollectionSecurityGroupsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNServerCollectionVPNServersSubnetsToMap(subnetsItem vpcv1.SubnetReference) (subnetsMap map[string]interface{}) {
	subnetsMap = map[string]interface{}{}

	if subnetsItem.CRN != nil {
		subnetsMap["crn"] = subnetsItem.CRN
	}
	if subnetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerCollectionSubnetsDeletedToMap(*subnetsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		subnetsMap["deleted"] = deletedList
	}
	if subnetsItem.Href != nil {
		subnetsMap["href"] = subnetsItem.Href
	}
	if subnetsItem.ID != nil {
		subnetsMap["id"] = subnetsItem.ID
	}
	if subnetsItem.Name != nil {
		subnetsMap["name"] = subnetsItem.Name
	}
	// if subnetsItem.ResourceType != nil {
	// 	subnetsMap["resource_type"] = subnetsItem.ResourceType
	// }

	return subnetsMap
}

func dataSourceVPNServerCollectionSubnetsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

// dataSourceVPNServerCollectionVPNServersVpcReferenceToMap
func dataSourceVPNServerCollectionVPNServersVpcReferenceToMap(vpcsItem *vpcv1.VPCReference) (vpcsMap map[string]interface{}) {
	vpcsMap = map[string]interface{}{}

	if vpcsItem.CRN != nil {
		vpcsMap["crn"] = vpcsItem.CRN
	}
	if vpcsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerCollectionVpcsDeletedToMap(*vpcsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcsMap["deleted"] = deletedList
	}
	if vpcsItem.Href != nil {
		vpcsMap["href"] = vpcsItem.Href
	}
	if vpcsItem.ID != nil {
		vpcsMap["id"] = vpcsItem.ID
	}
	if vpcsItem.Name != nil {
		vpcsMap["name"] = vpcsItem.Name
	}
	// if vpcsItem.ResourceType != nil {
	// 	vpcsMap["resource_type"] = vpcsItem.ResourceType
	// }

	return vpcsMap
}

func dataSourceVPNServerCollectionVpcsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
