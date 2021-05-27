// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const cisRangeApps = "range_apps"

func dataSourceIBMCISRangeApps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISRangeAppsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisRangeApps: {
				Type:        schema.TypeList,
				Description: "Collection of range application detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "range app id",
						},
						cisRangeAppID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application identifier",
						},
						cisRangeAppProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Defines the protocol and port for this application",
						},
						cisRangeAppDNS: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DNS record for this application",
						},
						cisRangeAppDNSType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the DNS record for this application",
						},
						cisRangeAppOriginDirect: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "IP address and port of the origin for this Range application.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						// cisRangeAppOriginDNS: {
						// 	Type:     schema.TypeString,
						// 	Computed: true,
						// },
						// cisRangeAppOriginPort: {
						// 	Type:     schema.TypeInt,
						// 	Computed: true,
						// },
						cisRangeAppIPFirewall: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enables the IP Firewall for this application. Only available for TCP applications.",
						},
						cisRangeAppProxyProtocol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Allows for the true client IP to be passed to the service.",
						},
						cisRangeAppEdgeIPsType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of edge IP configuration.",
						},
						cisRangeAppEdgeIPsConnectivity: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the IP version.",
						},
						cisRangeAppTrafficType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configure how traffic is handled at the edge.",
						},
						cisRangeAppTLS: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configure if and how TLS connections are terminated at the edge.",
						},
						cisRangeAppCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "created on date",
						},
						cisRangeAppModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "modified on date",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCISRangeAppsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRangeAppClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewListRangeAppsOptions()
	result, resp, err := cisClient.ListRangeApps(opt)
	if err != nil {
		return fmt.Errorf("Failed to list range applications: %v", resp)
	}
	apps := make([]map[string]interface{}, 0)
	for _, i := range result.Result {
		app := map[string]interface{}{}
		app["id"] = convertCisToTfThreeVar(*i.ID, zoneID, crn)
		app[cisRangeAppID] = *i.ID
		app[cisRangeAppProtocol] = *i.Protocol
		app[cisRangeAppDNS] = *i.Dns.Name
		app[cisRangeAppDNSType] = *i.Dns.Type
		app[cisRangeAppOriginDirect] = flattenStringList(i.OriginDirect)
		app[cisRangeAppIPFirewall] = *i.IpFirewall
		app[cisRangeAppProxyProtocol] = *i.ProxyProtocol
		app[cisRangeAppEdgeIPsType] = *i.EdgeIps.Type
		app[cisRangeAppEdgeIPsConnectivity] = *i.EdgeIps.Connectivity
		app[cisRangeAppTLS] = *i.Tls
		app[cisRangeAppTrafficType] = *i.TrafficType
		app[cisRangeAppCreatedOn] = (*i.CreatedOn).String()
		app[cisRangeAppModifiedOn] = (*i.ModifiedOn).String()
		apps = append(apps, app)

	}
	d.SetId(dataSourceIBMCISRangeAppsID(d))
	d.Set(cisRangeApps, apps)
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	return nil
}

func dataSourceIBMCISRangeAppsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
