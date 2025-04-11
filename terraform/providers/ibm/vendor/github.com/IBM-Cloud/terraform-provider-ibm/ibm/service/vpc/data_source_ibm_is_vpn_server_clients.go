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

func DataSourceIBMIsVPNServerClients() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServerClientsRead,

		Schema: map[string]*schema.Schema{
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN server identifier.",
			},
			"clients": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPN clients.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP address assigned to this VPN client from `client_ip_pool`.",
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
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The common name of client certificate that the VPN client provided when connecting to the server.This property will be present only when the `certificate` client authentication method is enabled on the VPN server.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the VPN client was created.",
						},
						"disconnected_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the VPN client was disconnected.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPN client.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPN client.",
						},
						"remote_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The remote IP address of this VPN client.",
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
						"remote_port": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The remote port of this VPN client.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the VPN client:- `connected`: the VPN client is `connected` to this VPN server.- `disconnected`: the VPN client is `disconnected` from this VPN server.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the VPN client on which the unexpected property value was encountered.",
						},
						"username": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username that this VPN client provided when connecting to the VPN server.This property will be present only when  the`username` client authentication method is enabled on the VPN server.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsVPNServerClientsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	start := ""
	allrecs := []vpcv1.VPNServerClient{}

	for {
		listVPNServerClientsOptions := &vpcv1.ListVPNServerClientsOptions{}
		listVPNServerClientsOptions.SetVPNServerID(d.Get("vpn_server").(string))
		if start != "" {
			listVPNServerClientsOptions.Start = &start
		}
		vpnServerClientCollection, response, err := sess.ListVPNServerClientsWithContext(context, listVPNServerClientsOptions)
		if err != nil {
			log.Printf("[DEBUG] ListVPNServerClientsWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("[ERROR] ListVPNServerClientsWithContext failed %s\n%s", err, response))
		}
		start = flex.GetNext(vpnServerClientCollection.Next)
		allrecs = append(allrecs, vpnServerClientCollection.Clients...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIBMIsVPNServerClientsID(d))

	if allrecs != nil {
		err = d.Set("clients", dataSourceVPNServerClientCollectionFlattenClients(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting clients %s", err))
		}
	}
	return nil
}

// dataSourceIBMIsVPNServerClientsID returns a reasonable ID for the list.
func dataSourceIBMIsVPNServerClientsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceVPNServerClientCollectionFlattenClients(result []vpcv1.VPNServerClient) (clients []map[string]interface{}) {
	for _, clientsItem := range result {
		clients = append(clients, dataSourceVPNServerClientCollectionClientsToMap(clientsItem))
	}

	return clients
}

func dataSourceVPNServerClientCollectionClientsToMap(clientsItem vpcv1.VPNServerClient) (clientsMap map[string]interface{}) {
	clientsMap = map[string]interface{}{}

	if clientsItem.ClientIP != nil {
		clientIPList := []map[string]interface{}{}
		clientIPMap := dataSourceVPNServerClientCollectionClientsClientIPToMap(*clientsItem.ClientIP)
		clientIPList = append(clientIPList, clientIPMap)
		clientsMap["client_ip"] = clientIPList
	}
	if clientsItem.CommonName != nil {
		clientsMap["common_name"] = clientsItem.CommonName
	}
	if clientsItem.CreatedAt != nil {
		clientsMap["created_at"] = clientsItem.CreatedAt.String()
	}
	if clientsItem.DisconnectedAt != nil {
		clientsMap["disconnected_at"] = clientsItem.DisconnectedAt.String()
	}
	if clientsItem.Href != nil {
		clientsMap["href"] = clientsItem.Href
	}
	if clientsItem.ID != nil {
		clientsMap["id"] = clientsItem.ID
	}
	if clientsItem.RemoteIP != nil {
		remoteIPList := []map[string]interface{}{}
		remoteIPMap := dataSourceVPNServerClientCollectionClientsRemoteIPToMap(*clientsItem.RemoteIP)
		remoteIPList = append(remoteIPList, remoteIPMap)
		clientsMap["remote_ip"] = remoteIPList
	}
	if clientsItem.RemotePort != nil {
		clientsMap["remote_port"] = clientsItem.RemotePort
	}
	if clientsItem.ResourceType != nil {
		clientsMap["resource_type"] = clientsItem.ResourceType
	}
	if clientsItem.Status != nil {
		clientsMap["status"] = clientsItem.Status
	}
	if clientsItem.Username != nil {
		clientsMap["username"] = clientsItem.Username
	}

	return clientsMap
}

func dataSourceVPNServerClientCollectionClientsClientIPToMap(clientIPItem vpcv1.IP) (clientIPMap map[string]interface{}) {
	clientIPMap = map[string]interface{}{}

	if clientIPItem.Address != nil {
		clientIPMap["address"] = clientIPItem.Address
	}

	return clientIPMap
}

func dataSourceVPNServerClientCollectionClientsRemoteIPToMap(remoteIPItem vpcv1.IP) (remoteIPMap map[string]interface{}) {
	remoteIPMap = map[string]interface{}{}

	if remoteIPItem.Address != nil {
		remoteIPMap["address"] = remoteIPItem.Address
	}

	return remoteIPMap
}

func dataSourceVPNServerClientCollectionFlattenFirst(result vpcv1.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerClientCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerClientCollectionFirstToMap(firstItem vpcv1.PageLink) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceVPNServerClientCollectionFlattenNext(result vpcv1.PageLink) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerClientCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerClientCollectionNextToMap(nextItem vpcv1.PageLink) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}
