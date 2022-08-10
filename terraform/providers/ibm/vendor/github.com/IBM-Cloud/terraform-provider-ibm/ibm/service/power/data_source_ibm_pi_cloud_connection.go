// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
)

const (
	PICloudConnectionId               = "cloud_connection_id"
	PICloudConnectionName             = "name"
	PICloudConnectionSpeed            = "speed"
	PICloudConnectionGlobalRouting    = "global_routing"
	PICloudConnectionMetered          = "metered"
	PICloudConnectionStatus           = "status"
	PICloudConnectionClassicEnabled   = "classic_enabled"
	PICloudConnectionUserIPAddress    = "user_ip_address"
	PICloudConnectionIBMIPAddress     = "ibm_ip_address"
	PICloudConnectionPort             = "port"
	PICloudConnectionNetworks         = "networks"
	PICloudConnectionClassicGreDest   = "gre_destination_address"
	PICloudConnectionClassicGreSource = "gre_source_address"
	PICloudConnectionVPCEnabled       = "vpc_enabled"
	PICloudConnectionVPCCRNs          = "vpc_crns"
	PICloudConnectionConnectionMode   = "connection_mode"
)

func DataSourceIBMPICloudConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPICloudConnectionRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			helpers.PICloudConnectionName: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Cloud Connection Name to be used",
				ValidateFunc: validation.NoZeroValues,
			},

			// Computed Attributes
			PICloudConnectionSpeed: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			PICloudConnectionGlobalRouting: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			PICloudConnectionMetered: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			PICloudConnectionStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PICloudConnectionIBMIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PICloudConnectionUserIPAddress: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PICloudConnectionPort: {
				Type:     schema.TypeString,
				Computed: true,
			},
			PICloudConnectionNetworks: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of Networks attached to this cloud connection",
			},
			PICloudConnectionClassicEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable classic endpoint destination",
			},
			PICloudConnectionClassicGreDest: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GRE destination IP address",
			},
			PICloudConnectionClassicGreSource: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GRE auto-assigned source IP address",
			},
			PICloudConnectionVPCEnabled: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Enable VPC for this cloud connection",
			},
			PICloudConnectionVPCCRNs: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of VPCs attached to this cloud connection",
			},
			PICloudConnectionConnectionMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of service the gateway is attached to",
			},
		},
	}
}

func dataSourceIBMPICloudConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	cloudConnectionName := d.Get(helpers.PICloudConnectionName).(string)
	client := instance.NewIBMPICloudConnectionClient(ctx, sess, cloudInstanceID)

	// Get API does not work with name for Cloud Connection hence using GetAll (max 2)
	// TODO: Uncomment Get call below when avaiable and remove GetAll
	// cloudConnectionD, err := client.GetWithContext(ctx, cloudConnectionName, cloudInstanceID)
	// if err != nil {
	// 	log.Printf("[DEBUG] get cloud connection failed %v", err)
	// 	return diag.Errorf(errors.GetCloudConnectionOperationFailed, cloudConnectionName, err)
	// }
	cloudConnections, err := client.GetAll()
	if err != nil {
		log.Printf("[DEBUG] get cloud connections failed %v", err)
		return diag.FromErr(err)
	}
	var cloudConnection *models.CloudConnection
	if cloudConnections != nil {
		for _, cc := range cloudConnections.CloudConnections {
			if cloudConnectionName == *cc.Name {
				cloudConnection = cc
				break
			}
		}
	}
	if cloudConnection == nil {
		log.Printf("[DEBUG] cloud connection not found")
		return diag.Errorf("failed to perform get cloud connection operation for name %s", cloudConnectionName)
	}

	cloudConnection, err = client.Get(*cloudConnection.CloudConnectionID)
	if err != nil {
		log.Printf("[DEBUG] get cloud connection failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*cloudConnection.CloudConnectionID)
	d.Set(helpers.PICloudConnectionName, cloudConnection.Name)
	d.Set(PICloudConnectionGlobalRouting, cloudConnection.GlobalRouting)
	d.Set(PICloudConnectionMetered, cloudConnection.Metered)
	d.Set(PICloudConnectionIBMIPAddress, cloudConnection.IbmIPAddress)
	d.Set(PICloudConnectionUserIPAddress, cloudConnection.UserIPAddress)
	d.Set(PICloudConnectionStatus, cloudConnection.LinkStatus)
	d.Set(PICloudConnectionPort, cloudConnection.Port)
	d.Set(PICloudConnectionSpeed, cloudConnection.Speed)
	d.Set(helpers.PICloudInstanceId, cloudInstanceID)
	d.Set(PICloudConnectionConnectionMode, cloudConnection.ConnectionMode)
	if cloudConnection.Networks != nil {
		networks := make([]string, len(cloudConnection.Networks))
		for i, ccNetwork := range cloudConnection.Networks {
			if ccNetwork != nil {
				networks[i] = *ccNetwork.NetworkID
			}
		}
		d.Set(PICloudConnectionNetworks, networks)
	}
	if cloudConnection.Classic != nil {
		d.Set(PICloudConnectionClassicEnabled, cloudConnection.Classic.Enabled)
		if cloudConnection.Classic.Gre != nil {
			d.Set(PICloudConnectionClassicGreDest, cloudConnection.Classic.Gre.DestIPAddress)
			d.Set(PICloudConnectionClassicGreSource, cloudConnection.Classic.Gre.SourceIPAddress)
		}
	}
	if cloudConnection.Vpc != nil {
		d.Set(PICloudConnectionVPCEnabled, cloudConnection.Vpc.Enabled)
		if cloudConnection.Vpc.Vpcs != nil && len(cloudConnection.Vpc.Vpcs) > 0 {
			vpcCRNs := make([]string, len(cloudConnection.Vpc.Vpcs))
			for i, vpc := range cloudConnection.Vpc.Vpcs {
				vpcCRNs[i] = *vpc.VpcID
			}
			d.Set(PICloudConnectionVPCCRNs, vpcCRNs)
		}
	}
	return nil
}
