// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

/*
Datasource to get the list of Cloud Connections in a power instance
*/

const PICloudConnections = "connections"

func DataSourceIBMPICloudConnections() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceIBMPICloudConnectionsRead,
		Schema: map[string]*schema.Schema{
			helpers.PICloudInstanceId: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			// Computed Attributes
			PICloudConnections: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						PICloudConnectionId: {
							Type:     schema.TypeString,
							Computed: true,
						},
						PICloudConnectionName: {
							Type:     schema.TypeString,
							Computed: true,
						},
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
				},
			},
		},
	}
}

func dataSourceIBMPICloudConnectionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	client := st.NewIBMPICloudConnectionClient(ctx, sess, cloudInstanceID)

	cloudConnections, err := client.GetAll()
	if err != nil {
		log.Printf("[DEBUG] get cloud connections failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(cloudConnections.CloudConnections))
	for _, cloudConnection := range cloudConnections.CloudConnections {
		cc := map[string]interface{}{
			PICloudConnectionId:             *cloudConnection.CloudConnectionID,
			PICloudConnectionName:           *cloudConnection.Name,
			PICloudConnectionGlobalRouting:  *cloudConnection.GlobalRouting,
			PICloudConnectionMetered:        *cloudConnection.Metered,
			PICloudConnectionIBMIPAddress:   *cloudConnection.IbmIPAddress,
			PICloudConnectionUserIPAddress:  *cloudConnection.UserIPAddress,
			PICloudConnectionStatus:         *cloudConnection.LinkStatus,
			PICloudConnectionPort:           *cloudConnection.Port,
			PICloudConnectionSpeed:          *cloudConnection.Speed,
			PICloudConnectionConnectionMode: cloudConnection.ConnectionMode,
		}

		if cloudConnection.Networks != nil {
			networks := make([]string, len(cloudConnection.Networks))
			for i, ccNetwork := range cloudConnection.Networks {
				if ccNetwork != nil {
					networks[i] = *ccNetwork.NetworkID
				}
			}
			cc[PICloudConnectionNetworks] = networks
		}
		if cloudConnection.Classic != nil {
			cc[PICloudConnectionClassicEnabled] = cloudConnection.Classic.Enabled
			if cloudConnection.Classic.Gre != nil {
				cc[PICloudConnectionClassicGreDest] = cloudConnection.Classic.Gre.DestIPAddress
				cc[PICloudConnectionClassicGreSource] = cloudConnection.Classic.Gre.SourceIPAddress
			}
		}
		if cloudConnection.Vpc != nil {
			cc[PICloudConnectionVPCEnabled] = cloudConnection.Vpc.Enabled
			if cloudConnection.Vpc.Vpcs != nil && len(cloudConnection.Vpc.Vpcs) > 0 {
				vpcCRNs := make([]string, len(cloudConnection.Vpc.Vpcs))
				for i, vpc := range cloudConnection.Vpc.Vpcs {
					vpcCRNs[i] = *vpc.VpcID
				}
				cc[PICloudConnectionVPCCRNs] = vpcCRNs
			}
		}

		result = append(result, cc)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(PICloudConnections, result)

	return nil
}
