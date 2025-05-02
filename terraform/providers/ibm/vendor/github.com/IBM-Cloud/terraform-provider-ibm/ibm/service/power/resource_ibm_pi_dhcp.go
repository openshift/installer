// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/errors"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_service_d_h_c_p"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

func ResourceIBMPIDhcp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIDhcpCreate,
		ReadContext:   resourceIBMPIDhcpRead,
		DeleteContext: resourceIBMPIDhcpDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{

			// Required Arguments
			Arg_CloudInstanceID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
				ForceNew:    true,
			},

			// Optional Arguments
			Arg_DhcpCidr: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional cidr for DHCP private network",
				ForceNew:    true,
			},
			Arg_DhcpCloudConnectionID: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional cloud connection uuid to connect with DHCP private network",
				ForceNew:    true,
			},
			Arg_DhcpDnsServer: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional DNS Server for DHCP service",
				ForceNew:    true,
			},
			Arg_DhcpName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional name of DHCP Service (will be prefixed by DHCP identifier)",
				ForceNew:    true,
			},
			Arg_DhcpSnatEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates if SNAT will be enabled for the DHCP service",
				ForceNew:    true,
			},

			// Attributes
			Attr_DhcpID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server",
			},
			Attr_DhcpLeases: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DHCP Server PVM Instance leases",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_DhcpLeaseInstanceIP: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of the PVM Instance",
						},
						Attr_DhcpLeaseInstanceMac: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC Address of the PVM Instance",
						},
					},
				},
			},
			Attr_DhcpNetworkDeprecated: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server private network (deprecated - replaced by network_id)",
			},
			Attr_DhcpNetworkID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server private network",
			},
			Attr_DhcpNetworkName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the DHCP Server private network",
			},
			Attr_DhcpStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the DHCP Server",
			},
		},
	}
}

func resourceIBMPIDhcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// dhcp create object
	body := &models.DHCPServerCreate{}

	// arguments
	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	if cidr, ok := d.GetOk(Arg_DhcpCidr); ok {
		c := cidr.(string)
		body.Cidr = &c
	}
	if cloudConnectionID, ok := d.GetOk(Arg_DhcpCloudConnectionID); ok {
		c := cloudConnectionID.(string)
		body.CloudConnectionID = &c
	}
	if dnsServer, ok := d.GetOk(Arg_DhcpDnsServer); ok {
		d := dnsServer.(string)
		body.DNSServer = &d
	}
	if name, ok := d.GetOk(Arg_DhcpName); ok {
		n := name.(string)
		body.Name = &n
	}
	snatEnabled := d.Get(Arg_DhcpSnatEnabled).(bool)
	body.SnatEnabled = &snatEnabled

	// create dhcp
	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServer, err := client.Create(body)
	if err != nil {
		log.Printf("[DEBUG] create DHCP failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *dhcpServer.ID))

	// wait for creation
	_, err = waitForIBMPIDhcpStatus(ctx, client, *dhcpServer.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		diag.FromErr(err)
	}

	return resourceIBMPIDhcpRead(ctx, d, meta)
}

func resourceIBMPIDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID, dhcpID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// get dhcp
	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServer, err := client.Get(dhcpID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_service_d_h_c_p.PcloudDhcpGetNotFound:
			log.Printf("[DEBUG] dhcp does not exist %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] get DHCP failed %v", err)
		return diag.FromErr(err)
	}

	// set attributes
	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *dhcpServer.ID))
	d.Set(Attr_DhcpID, *dhcpServer.ID)
	d.Set(Attr_DhcpStatus, *dhcpServer.Status)

	if dhcpServer.Network != nil {
		dhcpNetwork := dhcpServer.Network
		if dhcpNetwork.ID != nil {
			d.Set(Attr_DhcpNetworkDeprecated, *dhcpNetwork.ID)
			d.Set(Attr_DhcpNetworkID, *dhcpNetwork.ID)
		}
		if dhcpNetwork.Name != nil {
			d.Set(Attr_DhcpNetworkName, *dhcpNetwork.Name)
		}
	}

	if dhcpServer.Leases != nil {
		leaseList := make([]map[string]string, len(dhcpServer.Leases))
		for i, lease := range dhcpServer.Leases {
			leaseList[i] = map[string]string{
				Attr_DhcpLeaseInstanceIP:  *lease.InstanceIP,
				Attr_DhcpLeaseInstanceMac: *lease.InstanceMacAddress,
			}
		}
		d.Set(Attr_DhcpLeases, leaseList)
	}

	return nil
}
func resourceIBMPIDhcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	// arguments
	cloudInstanceID, dhcpID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// delete dhcp
	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	err = client.Delete(dhcpID)
	if err != nil {
		uErr := errors.Unwrap(err)
		switch uErr.(type) {
		case *p_cloud_service_d_h_c_p.PcloudDhcpDeleteNotFound:
			log.Printf("[DEBUG] dhcp does not exist %v", err)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] delete DHCP failed %v", err)
		return diag.FromErr(err)
	}

	// wait for deletion
	_, err = waitForIBMPIDhcpDeleted(ctx, client, dhcpID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func waitForIBMPIDhcpStatus(ctx context.Context, client *st.IBMPIDhcpClient, dhcpID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"building"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			dhcpServer, err := client.Get(dhcpID)
			if err != nil {
				log.Printf("[DEBUG] get DHCP failed %v", err)
				return nil, "", err
			}
			if *dhcpServer.Status != StatusActive {
				return dhcpServer, "building", nil
			}
			return dhcpServer, "active", nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func waitForIBMPIDhcpDeleted(ctx context.Context, client *st.IBMPIDhcpClient, dhcpID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (interface{}, string, error) {
			dhcpServer, err := client.Get(dhcpID)
			if err != nil {
				log.Printf("[DEBUG] dhcp does not exist %v", err)
				return dhcpServer, "deleted", nil
			}
			return dhcpServer, "deleting", nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}
