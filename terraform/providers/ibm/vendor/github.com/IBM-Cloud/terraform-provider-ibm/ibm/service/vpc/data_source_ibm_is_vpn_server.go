// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMIsVPNServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServerRead,

		Schema: map[string]*schema.Schema{

			"identifier": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The unique identifier for this VPN server",
			},

			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"name", "identifier"},
				Description:  "The unique user-defined name for this VPN server",
			},

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
				Type:        schema.TypeList,
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
	}
}

func dataSourceIBMIsVPNServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	var vpnServer *vpcv1.VPNServer

	if v, ok := d.GetOk("identifier"); ok {

		getVPNServerOptions := &vpcv1.GetVPNServerOptions{}
		getVPNServerOptions.SetID(v.(string))
		vpnServerInfo, response, err := sess.GetVPNServerWithContext(context, getVPNServerOptions)
		if err != nil {
			log.Printf("[DEBUG] GetVPNServerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerWithContext failed %s\n%s", err, response))
		}
		vpnServer = vpnServerInfo
	} else if v, ok := d.GetOk("name"); ok {

		name := v.(string)
		start := ""
		allrecs := []vpcv1.VPNServer{}

		for {
			listVPNServersOptions := &vpcv1.ListVPNServersOptions{}
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

		for _, vpnServerInfo := range allrecs {
			if *vpnServerInfo.Name == name {
				vpnServer = &vpnServerInfo
				break
			}
		}
		if vpnServer == nil {
			log.Printf("[DEBUG] No vpnServer found with name %s", name)
			return diag.FromErr(fmt.Errorf("[ERROR] No vpn server found with name %s", name))
		}
	}

	d.SetId(fmt.Sprintf("%s", *vpnServer.ID))
	err = d.Set("identifier", vpnServer.ID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting identifier %s", err))
	}

	if vpnServer.Certificate != nil {
		err = d.Set("certificate", dataSourceVPNServerFlattenCertificate(*vpnServer.Certificate))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting certificate %s", err))
		}
	}

	vpnServerAuthenticationPrototypeArray := make([]interface{}, len(vpnServer.ClientAuthentication))
	if vpnServer.ClientAuthentication != nil {
		for i, clientAuthenticationItem := range vpnServer.ClientAuthentication {
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
	}
	err = d.Set("client_authentication", vpnServerAuthenticationPrototypeArray)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_authentication %s", err))
	}

	if err = d.Set("client_auto_delete", vpnServer.ClientAutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_auto_delete: %s", err))
	}
	if err = d.Set("client_auto_delete_timeout", flex.IntValue(vpnServer.ClientAutoDeleteTimeout)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_auto_delete_timeout: %s", err))
	}

	if vpnServer.ClientDnsServerIps != nil {
		err = d.Set("client_dns_server_ips", dataSourceVPNServerFlattenClientDnsServerIps(vpnServer.ClientDnsServerIps))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_dns_server_ips %s", err))
		}
	}
	if err = d.Set("client_idle_timeout", flex.IntValue(vpnServer.ClientIdleTimeout)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_idle_timeout: %s", err))
	}
	if err = d.Set("client_ip_pool", vpnServer.ClientIPPool); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_ip_pool: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(vpnServer.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("crn", vpnServer.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting crn: %s", err))
	}
	if err = d.Set("enable_split_tunneling", vpnServer.EnableSplitTunneling); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enable_split_tunneling: %s", err))
	}
	if err = d.Set("health_state", vpnServer.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_state: %s", err))
	}
	if vpnServer.HealthReasons != nil {
		if err := d.Set("health_reasons", resourceVPNServerFlattenHealthReasons(vpnServer.HealthReasons)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting health_reasons: %s", err))
		}
	}
	if err = d.Set("hostname", vpnServer.Hostname); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting hostname: %s", err))
	}
	if err = d.Set("href", vpnServer.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", vpnServer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if vpnServer.LifecycleReasons != nil {
		if err := d.Set("lifecycle_reasons", resourceVPNServerFlattenLifecycleReasons(vpnServer.LifecycleReasons)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_reasons: %s", err))
		}
	}
	if err = d.Set("name", vpnServer.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("port", flex.IntValue(vpnServer.Port)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting port: %s", err))
	}

	if vpnServer.PrivateIps != nil {
		err = d.Set("private_ips", dataSourceVPNServerFlattenPrivateIps(vpnServer.PrivateIps))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting private_ips %s", err))
		}
	}
	if err = d.Set("protocol", vpnServer.Protocol); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting protocol: %s", err))
	}

	if vpnServer.ResourceGroup != nil {
		err = d.Set("resource_group", dataSourceVPNServerFlattenResourceGroup(*vpnServer.ResourceGroup))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", vpnServer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}

	if vpnServer.SecurityGroups != nil {
		err = d.Set("security_groups", dataSourceVPNServerFlattenSecurityGroups(vpnServer.SecurityGroups))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting security_groups %s", err))
		}
	}

	if vpnServer.Subnets != nil {
		err = d.Set("subnets", dataSourceVPNServerFlattenSubnets(vpnServer.Subnets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnets %s", err))
		}
	}

	if vpnServer.VPC != nil {
		err = d.Set("vpc", dataSourceVPNServerFlattenVpcReference(vpnServer.VPC))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting the vpc: %s", err))
		}
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *vpnServer.CRN, "", isVPNServerAccessTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of vpn server (%s) access tags: %s", d.Id(), err)
	}
	d.Set(isVPNServerAccessTags, accesstags)

	return nil
}

func dataSourceVPNServerFlattenCertificate(result vpcv1.CertificateInstanceReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerCertificateToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerCertificateToMap(certificateItem vpcv1.CertificateInstanceReference) (certificateMap map[string]interface{}) {
	certificateMap = map[string]interface{}{}

	if certificateItem.CRN != nil {
		certificateMap["crn"] = certificateItem.CRN
	}

	return certificateMap
}

func dataSourceVPNServerClientAuthenticationClientCaToMap(clientCaItem vpcv1.CertificateInstanceReference) (clientCaMap map[string]interface{}) {
	clientCaMap = map[string]interface{}{}

	if clientCaItem.CRN != nil {
		clientCaMap["crn"] = clientCaItem.CRN
	}

	return clientCaMap
}

func dataSourceVPNServerFlattenClientDnsServerIps(result []vpcv1.IP) (clientDnsServerIps []map[string]interface{}) {
	for _, clientDnsServerIpsItem := range result {
		clientDnsServerIps = append(clientDnsServerIps, dataSourceVPNServerClientDnsServerIpsToMap(clientDnsServerIpsItem))
	}

	return clientDnsServerIps
}

func dataSourceVPNServerClientDnsServerIpsToMap(clientDnsServerIpsItem vpcv1.IP) (clientDnsServerIpsMap map[string]interface{}) {
	clientDnsServerIpsMap = map[string]interface{}{}

	if clientDnsServerIpsItem.Address != nil {
		clientDnsServerIpsMap["address"] = clientDnsServerIpsItem.Address
	}

	return clientDnsServerIpsMap
}

func dataSourceVPNServerFlattenPrivateIps(result []vpcv1.ReservedIPReference) (privateIps []map[string]interface{}) {
	for _, privateIpsItem := range result {
		privateIps = append(privateIps, dataSourceVPNServerPrivateIpsToMap(privateIpsItem))
	}

	return privateIps
}

func dataSourceVPNServerPrivateIpsToMap(privateIpsItem vpcv1.ReservedIPReference) (privateIpsMap map[string]interface{}) {
	privateIpsMap = map[string]interface{}{}

	if privateIpsItem.Address != nil {
		privateIpsMap["address"] = privateIpsItem.Address
	}
	if privateIpsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerPrivateIpsDeletedToMap(*privateIpsItem.Deleted)
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

func dataSourceVPNServerPrivateIpsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNServerFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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

func dataSourceVPNServerFlattenSecurityGroups(result []vpcv1.SecurityGroupReference) (securityGroups []map[string]interface{}) {
	for _, securityGroupsItem := range result {
		securityGroups = append(securityGroups, dataSourceVPNServerSecurityGroupsToMap(securityGroupsItem))
	}

	return securityGroups
}

func dataSourceVPNServerSecurityGroupsToMap(securityGroupsItem vpcv1.SecurityGroupReference) (securityGroupsMap map[string]interface{}) {
	securityGroupsMap = map[string]interface{}{}

	if securityGroupsItem.CRN != nil {
		securityGroupsMap["crn"] = securityGroupsItem.CRN
	}
	if securityGroupsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerSecurityGroupsDeletedToMap(*securityGroupsItem.Deleted)
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

func dataSourceVPNServerSecurityGroupsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceVPNServerFlattenSubnets(result []vpcv1.SubnetReference) (subnets []map[string]interface{}) {
	for _, subnetsItem := range result {
		subnets = append(subnets, dataSourceVPNServerSubnetsToMap(subnetsItem))
	}

	return subnets
}

func dataSourceVPNServerSubnetsToMap(subnetsItem vpcv1.SubnetReference) (subnetsMap map[string]interface{}) {
	subnetsMap = map[string]interface{}{}

	if subnetsItem.CRN != nil {
		subnetsMap["crn"] = subnetsItem.CRN
	}
	if subnetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerSubnetsDeletedToMap(*subnetsItem.Deleted)
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

	return subnetsMap
}

func dataSourceVPNServerSubnetsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

// dataSourceVPNServerFlattenVpcRefrence
func dataSourceVPNServerFlattenVpcReference(result *vpcv1.VPCReference) (vpcs []map[string]interface{}) {
	vpcs = append(vpcs, dataSourceVPNServerVpcToMap(*result))
	return vpcs
}
func dataSourceVPNServerVpcToMap(vpcItem vpcv1.VPCReference) (vpcsMap map[string]interface{}) {
	vpcsMap = map[string]interface{}{}

	if vpcItem.CRN != nil {
		vpcsMap["crn"] = vpcItem.CRN
	}
	if vpcItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceVPNServerVpcsDeletedToMap(*vpcItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		vpcsMap["deleted"] = deletedList
	}
	if vpcItem.Href != nil {
		vpcsMap["href"] = vpcItem.Href
	}
	if vpcItem.ID != nil {
		vpcsMap["id"] = vpcItem.ID
	}
	if vpcItem.Name != nil {
		vpcsMap["name"] = vpcItem.Name
	}

	return vpcsMap
}

func dataSourceVPNServerVpcsDeletedToMap(deletedItem vpcv1.Deleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
