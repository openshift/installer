// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMPICloudConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPICloudConnectionRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_CloudConnectionName: {
				Description:  "The cloud connection name to be used.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_ClassicEnabled: {
				Computed:    true,
				Description: "Enable classic endpoint destination.",
				Type:        schema.TypeBool,
			},
			Attr_ConnectionMode: {
				Computed:    true,
				Description: "Type of service the gateway is attached to.",
				Type:        schema.TypeString,
			},
			Attr_GlobalRouting: {
				Computed:    true,
				Description: "Enable global routing for this cloud connection.",
				Type:        schema.TypeBool,
			},
			Attr_GreDestinationAddress: {
				Computed:    true,
				Description: "GRE destination IP address.",
				Type:        schema.TypeString,
			},
			Attr_GreSourceAddress: {
				Computed:    true,
				Description: "GRE auto-assigned source IP address.",
				Type:        schema.TypeString,
			},
			Attr_IBMIPAddress: {
				Computed:    true,
				Description: "The IBM IP address.",
				Type:        schema.TypeString,
			},
			Attr_Metered: {
				Computed:    true,
				Description: "Enable metering for this cloud connection.",
				Type:        schema.TypeBool,
			},
			Attr_Networks: {
				Computed:    true,
				Description: "Set of Networks attached to this cloud connection",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_Port: {
				Computed:    true,
				Description: "Port.",
				Type:        schema.TypeString,
			},
			Attr_Speed: {
				Computed:    true,
				Description: "Speed of the cloud connection (speed in megabits per second)",
				Type:        schema.TypeInt,
			},
			Attr_Status: {
				Computed:    true,
				Description: "Link status.",
				Type:        schema.TypeString,
			},
			Attr_UserIPAddress: {
				Computed:    true,
				Description: "User IP address.",
				Type:        schema.TypeString,
			},
			Attr_VPCCRNs: {
				Computed:    true,
				Description: "Set of VPCs attached to this cloud connection.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Type:        schema.TypeSet,
			},
			Attr_VPCEnabled: {
				Computed:    true,
				Description: "Enable VPC for this cloud connection.",
				Type:        schema.TypeBool,
			},
		},
	}
}

func dataSourceIBMPICloudConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	cloudConnectionName := d.Get(Arg_CloudConnectionName).(string)
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

	d.Set(Arg_CloudInstanceID, cloudInstanceID)
	d.Set(Arg_CloudConnectionName, cloudConnection.Name)

	d.Set(Attr_GlobalRouting, cloudConnection.GlobalRouting)
	d.Set(Attr_Metered, cloudConnection.Metered)
	d.Set(Attr_IBMIPAddress, cloudConnection.IbmIPAddress)
	d.Set(Attr_UserIPAddress, cloudConnection.UserIPAddress)
	d.Set(Attr_Status, cloudConnection.LinkStatus)
	d.Set(Attr_Port, cloudConnection.Port)
	d.Set(Attr_Speed, cloudConnection.Speed)
	d.Set(Attr_ConnectionMode, cloudConnection.ConnectionMode)
	if cloudConnection.Networks != nil {
		networks := make([]string, len(cloudConnection.Networks))
		for i, ccNetwork := range cloudConnection.Networks {
			if ccNetwork != nil {
				networks[i] = *ccNetwork.NetworkID
			}
		}
		d.Set(Attr_Networks, networks)
	}
	if cloudConnection.Classic != nil {
		d.Set(Attr_ClassicEnabled, cloudConnection.Classic.Enabled)
		if cloudConnection.Classic.Gre != nil {
			d.Set(Attr_GreDestinationAddress, cloudConnection.Classic.Gre.DestIPAddress)
			d.Set(Attr_GreSourceAddress, cloudConnection.Classic.Gre.SourceIPAddress)
		}
	}
	if cloudConnection.Vpc != nil {
		d.Set(Attr_VPCEnabled, cloudConnection.Vpc.Enabled)
		if cloudConnection.Vpc.Vpcs != nil && len(cloudConnection.Vpc.Vpcs) > 0 {
			vpcCRNs := make([]string, len(cloudConnection.Vpc.Vpcs))
			for i, vpc := range cloudConnection.Vpc.Vpcs {
				vpcCRNs[i] = *vpc.VpcID
			}
			d.Set(Attr_VPCCRNs, vpcCRNs)
		}
	}

	return nil
}
