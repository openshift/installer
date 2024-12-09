// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"
	"log"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Datasource to list Cloud Connections in a power instance
func DataSourceIBMPICloudConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPICloudConnectionsRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Connections: {
				Computed:    true,
				Description: "List of all the Cloud Connections.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_ClassicEnabled: {
							Computed:    true,
							Description: "Enable classic endpoint destination.",
							Type:        schema.TypeBool,
						},
						Attr_CloudConnectionID: {
							Computed:    true,
							Description: "The unique identifier of the cloud connection.",
							Type:        schema.TypeString,
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
							Description: "IBM IP address.",
							Type:        schema.TypeString,
						},
						Attr_Metered: {
							Computed:    true,
							Description: "Enable metering for this cloud connection.",
							Type:        schema.TypeBool,
						},
						Attr_Name: {
							Computed:    true,
							Description: "Name of the cloud connection.",
							Type:        schema.TypeString,
						},
						Attr_Networks: {
							Computed:    true,
							Description: "Set of Networks attached to this cloud connection.",
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
							Description: "Speed of the cloud connection (speed in megabits per second).",
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
				},
				Type: schema.TypeList,
			},
		},
	}
}

func dataSourceIBMPICloudConnectionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPICloudConnectionClient(ctx, sess, cloudInstanceID)

	cloudConnections, err := client.GetAll()
	if err != nil {
		log.Printf("[DEBUG] get cloud connections failed %v", err)
		return diag.FromErr(err)
	}

	result := make([]map[string]interface{}, 0, len(cloudConnections.CloudConnections))
	for _, cloudConnection := range cloudConnections.CloudConnections {
		cc := map[string]interface{}{
			Attr_CloudConnectionID: *cloudConnection.CloudConnectionID,
			Attr_ConnectionMode:    cloudConnection.ConnectionMode,
			Attr_GlobalRouting:     *cloudConnection.GlobalRouting,
			Attr_IBMIPAddress:      *cloudConnection.IbmIPAddress,
			Attr_Metered:           *cloudConnection.Metered,
			Attr_Name:              *cloudConnection.Name,
			Attr_Port:              *cloudConnection.Port,
			Attr_Speed:             *cloudConnection.Speed,
			Attr_Status:            *cloudConnection.LinkStatus,
			Attr_UserIPAddress:     *cloudConnection.UserIPAddress,
		}

		if cloudConnection.Networks != nil {
			networks := make([]string, len(cloudConnection.Networks))
			for i, ccNetwork := range cloudConnection.Networks {
				if ccNetwork != nil {
					networks[i] = *ccNetwork.NetworkID
				}
			}
			cc[Attr_Networks] = networks
		}
		if cloudConnection.Classic != nil {
			cc[Attr_ClassicEnabled] = cloudConnection.Classic.Enabled
			if cloudConnection.Classic.Gre != nil {
				cc[Attr_GreDestinationAddress] = cloudConnection.Classic.Gre.DestIPAddress
				cc[Attr_GreSourceAddress] = cloudConnection.Classic.Gre.SourceIPAddress
			}
		}
		if cloudConnection.Vpc != nil {
			cc[Attr_VPCEnabled] = cloudConnection.Vpc.Enabled
			if cloudConnection.Vpc.Vpcs != nil && len(cloudConnection.Vpc.Vpcs) > 0 {
				vpcCRNs := make([]string, len(cloudConnection.Vpc.Vpcs))
				for i, vpc := range cloudConnection.Vpc.Vpcs {
					vpcCRNs[i] = *vpc.VpcID
				}
				cc[Attr_VPCCRNs] = vpcCRNs
			}
		}

		result = append(result, cc)
	}

	var genID, _ = uuid.GenerateUUID()
	d.SetId(genID)
	d.Set(Attr_Connections, result)

	return nil
}
