// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	pdnsGLBMonitors = "dns_glb_monitors"
)

func dataSourceIBMPrivateDNSGLBMonitors() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMPrivateDNSGLBMonitorsRead,

		Schema: map[string]*schema.Schema{

			pdnsInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The GUID of the private DNS. ",
			},

			pdnsGLBMonitors: {
				Type:        schema.TypeList,
				Description: "Collection of GLB monitors collectors",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						pdnsGlbMonitorID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitor Id",
						},

						pdnsGlbMonitorName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier of a service instance.",
						},

						pdnsGlbMonitorDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive text of the load balancer monitor",
						},

						pdnsGlbMonitorType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol to use for the health check",
						},

						pdnsGlbMonitorPort: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Port number to connect to for the health check",
						},

						pdnsGlbMonitorInterval: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The interval between each health check",
						},

						pdnsGlbMonitorRetries: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of retries to attempt in case of a timeout before marking the origin as unhealthy",
						},

						pdnsGlbMonitorTimeout: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timeout (in seconds) before marking the health check as failed",
						},

						pdnsGlbMonitorMethod: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The method to use for the health check",
						},

						pdnsGlbMonitorPath: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The endpoint path to health check against",
						},

						pdnsGlbMonitorAllowInsecure: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Do not validate the certificate when monitor use HTTPS. This parameter is currently only valid for HTTPS monitors.",
						},

						pdnsGlbMonitorExpectedCodes: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expected HTTP response code or code range of the health check. This parameter is only valid for HTTP and HTTPS",
						},

						pdnsGlbMonitorExpectedBody: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A case-insensitive sub-string to look for in the response body",
						},

						pdnsGlbMonitorCreatedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Monitor creation date",
						},

						pdnsGlbMonitorModifiedOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "GLB Monitor Modification date",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPrivateDNSGLBMonitorsRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).PrivateDNSClientSession()
	if err != nil {
		return err
	}
	instanceID := d.Get(pdnsInstanceID).(string)
	listDNSGLBMonitorions := sess.NewListMonitorsOptions(instanceID)
	availableGLBMonitors, detail, err := sess.ListMonitors(listDNSGLBMonitorions)
	if err != nil {
		return fmt.Errorf("Error reading list of pdns GLB monitors:%s\n%s", err, detail)
	}

	dnsMonitors := make([]map[string]interface{}, 0)
	for _, instance := range availableGLBMonitors.Monitors {
		dnsMonitor := map[string]interface{}{}
		dnsMonitor[pdnsGlbMonitorID] = *instance.ID
		dnsMonitor[pdnsGlbMonitorName] = *instance.Name
		dnsMonitor[pdnsGlbMonitorType] = *instance.Type
		dnsMonitor[pdnsGlbMonitorCreatedOn] = *instance.CreatedOn
		dnsMonitor[pdnsGlbMonitorModifiedOn] = *instance.ModifiedOn
		dnsMonitor[pdnsGlbMonitorPort] = *instance.Port
		dnsMonitor[pdnsGlbMonitorInterval] = *instance.Interval
		dnsMonitor[pdnsGlbMonitorRetries] = *instance.Retries
		dnsMonitor[pdnsGlbMonitorTimeout] = *instance.Timeout
		dnsMonitor[pdnsGlbMonitorDescription] = *instance.Description
		dnsMonitor[pdnsGlbMonitorMethod] = *instance.Method
		dnsMonitor[pdnsGlbMonitorPath] = *instance.Path
		dnsMonitor[pdnsGlbMonitorExpectedCodes] = *instance.ExpectedCodes
		dnsMonitor[pdnsGlbMonitorExpectedBody] = *instance.ExpectedBody
		dnsMonitor[pdnsGlbMonitorAllowInsecure] = *instance.AllowInsecure

		dnsMonitors = append(dnsMonitors, dnsMonitor)
	}
	d.SetId(dataSourceIBMPrivateDNSGLBMonitorsID(d))
	d.Set(pdnsGLBMonitors, dnsMonitors)
	return nil
}

// dataSourceIBMPrivateDNSGLBMonitorsID returns a reasonable ID  list.
func dataSourceIBMPrivateDNSGLBMonitorsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
