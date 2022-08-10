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

func DataSourceIBMIsVPNServerClient() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNServerClientRead,

		Schema: map[string]*schema.Schema{
			"vpn_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN server identifier.",
			},
			"identifier": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN client identifier.",
			},
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
	}
}

func dataSourceIBMIsVPNServerClientRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getVPNServerClientOptions := &vpcv1.GetVPNServerClientOptions{}

	getVPNServerClientOptions.SetVPNServerID(d.Get("vpn_server").(string))
	getVPNServerClientOptions.SetID(d.Get("identifier").(string))

	vpnServerClient, response, err := sess.GetVPNServerClientWithContext(context, getVPNServerClientOptions)
	if err != nil {
		log.Printf("[DEBUG] GetVPNServerClientWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("[ERROR] GetVPNServerClientWithContext failed %s\n%s", err, response))
	}

	d.SetId(*vpnServerClient.ID)

	if vpnServerClient.ClientIP != nil {
		err = d.Set("client_ip", dataSourceVPNServerClientFlattenClientIP(*vpnServerClient.ClientIP))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting client_ip %s", err))
		}
	}
	if err = d.Set("common_name", vpnServerClient.CommonName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting common_name: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(vpnServerClient.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("disconnected_at", flex.DateTimeToString(vpnServerClient.DisconnectedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting disconnected_at: %s", err))
	}
	if err = d.Set("href", vpnServerClient.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}

	if vpnServerClient.RemoteIP != nil {
		err = d.Set("remote_ip", dataSourceVPNServerClientFlattenRemoteIP(*vpnServerClient.RemoteIP))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting remote_ip %s", err))
		}
	}
	if err = d.Set("remote_port", flex.IntValue(vpnServerClient.RemotePort)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting remote_port: %s", err))
	}
	if err = d.Set("resource_type", vpnServerClient.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("status", vpnServerClient.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if err = d.Set("username", vpnServerClient.Username); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting username: %s", err))
	}

	return nil
}

func dataSourceVPNServerClientFlattenClientIP(result vpcv1.IP) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerClientClientIPToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerClientClientIPToMap(clientIPItem vpcv1.IP) (clientIPMap map[string]interface{}) {
	clientIPMap = map[string]interface{}{}

	if clientIPItem.Address != nil {
		clientIPMap["address"] = clientIPItem.Address
	}

	return clientIPMap
}

func dataSourceVPNServerClientFlattenRemoteIP(result vpcv1.IP) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceVPNServerClientRemoteIPToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceVPNServerClientRemoteIPToMap(remoteIPItem vpcv1.IP) (remoteIPMap map[string]interface{}) {
	remoteIPMap = map[string]interface{}{}

	if remoteIPItem.Address != nil {
		remoteIPMap["address"] = remoteIPItem.Address
	}

	return remoteIPMap
}
