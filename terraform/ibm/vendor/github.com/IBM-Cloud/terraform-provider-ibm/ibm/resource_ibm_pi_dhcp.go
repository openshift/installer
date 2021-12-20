// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

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
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_service_d_h_c_p"
)

const (
	PIDhcpStatusBuilding = "Building"
	PIDhcpStatusActive   = "ACTIVE"
	PIDhcpDeleting       = "Deleting"
	PIDhcpDeleted        = "Deleted"
	PIDhcpId             = "dhcp_id"
	PIDhcpStatus         = "status"
	PIDhcpNetwork        = "network"
	PIDhcpLeases         = "leases"
	PIDhcpInstanceIp     = "instance_ip"
	PIDhcpInstanceMac    = "instance_mac"
)

func resourceIBMPIDhcp() *schema.Resource {
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
			// Required Attributes
			helpers.PICloudInstanceId: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "PI cloud instance ID",
				ForceNew:    true,
			},
			//Computed Attributes
			PIDhcpId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the DHCP Server",
			},
			PIDhcpStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the DHCP Server",
			},
			PIDhcpNetwork: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The DHCP Server private network",
			},
			PIDhcpLeases: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of DHCP Server PVM Instance leases",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PIDhcpInstanceIp: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP of the PVM Instance",
						},
						PIDhcpInstanceMac: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC Address of the PVM Instance",
						},
					},
				},
			},
		},
	}
}

func resourceIBMPIDhcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)

	client := st.NewIBMPIDhcpClient(ctx, sess, cloudInstanceID)
	dhcpServer, err := client.Create()
	if err != nil {
		log.Printf("[DEBUG] create DHCP failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", cloudInstanceID, *dhcpServer.ID))

	_, err = waitForIBMPIDhcpStatus(ctx, client, *dhcpServer.ID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		diag.FromErr(err)
	}

	return resourceIBMPIDhcpRead(ctx, d, meta)
}

func resourceIBMPIDhcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, dhcpID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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

	d.Set(PIDhcpId, *dhcpServer.ID)
	d.Set(PIDhcpStatus, *dhcpServer.Status)

	if dhcpServer.Network != nil {
		d.Set(PIDhcpNetwork, dhcpServer.Network.ID)
	}
	if dhcpServer.Leases != nil {
		leaseList := make([]map[string]string, len(dhcpServer.Leases))
		for i, lease := range dhcpServer.Leases {
			leaseList[i] = map[string]string{
				PIDhcpInstanceIp:  *lease.InstanceIP,
				PIDhcpInstanceMac: *lease.InstanceMacAddress,
			}
		}
		d.Set(PIDhcpLeases, leaseList)
	}

	return nil
}
func resourceIBMPIDhcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID, dhcpID, err := splitID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

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
	_, err = waitForIBMPIDhcpDeleted(ctx, client, dhcpID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func waitForIBMPIDhcpStatus(ctx context.Context, client *st.IBMPIDhcpClient, dhcpID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{PIDhcpStatusBuilding},
		Target:  []string{PIDhcpStatusActive},
		Refresh: func() (interface{}, string, error) {
			dhcpServer, err := client.Get(dhcpID)
			if err != nil {
				log.Printf("[DEBUG] get DHCP failed %v", err)
				return nil, "", err
			}
			if *dhcpServer.Status != PIDhcpStatusActive {
				return dhcpServer, PIDhcpStatusBuilding, nil
			}
			return dhcpServer, *dhcpServer.Status, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func waitForIBMPIDhcpDeleted(ctx context.Context, client *st.IBMPIDhcpClient, dhcpID string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{PIDhcpDeleting},
		Target:  []string{PIDhcpDeleted},
		Refresh: func() (interface{}, string, error) {
			dhcpServer, err := client.Get(dhcpID)
			if err != nil {
				log.Printf("[DEBUG] dhcp does not exist %v", err)
				return dhcpServer, PIDhcpDeleted, nil
			}
			return dhcpServer, PIDhcpDeleting, nil
		},
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}
