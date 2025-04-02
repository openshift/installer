// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DatasourceIBMPIWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIWorkspaceRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_WorkspaceCapabilities: {
				Computed:    true,
				Description: "Workspace Capabilities.",
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
				Type: schema.TypeMap,
			},
			Attr_WorkspaceDetails: {
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						Attr_CreationDate: {
							Computed:    true,
							Description: "Workspace creation date.",
							Type:        schema.TypeString,
						},
						Attr_CRN: {
							Computed:    true,
							Description: "The Workspace crn.",
							Type:        schema.TypeString,
						},
						Attr_NetworkSecurityGroups: {
							Computed:    true,
							Description: "Network security groups configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_State: {
										Computed:    true,
										Description: "The state of a Network Security Groups configuration.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
						},
						Attr_PowerEdgeRouter: {
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									Attr_MigrationStatus: {
										Computed:    true,
										Description: "The migration status of a Power Edge Router.",
										Type:        schema.TypeString,
									},
									Attr_State: {
										Computed:    true,
										Description: "The state of a Power Edge Router.",
										Type:        schema.TypeString,
									},
									Attr_Type: {
										Computed:    true,
										Description: "The Power Edge Router type.",
										Type:        schema.TypeString,
									},
								},
							},
							Type: schema.TypeList,
						},
					},
				},
				Type: schema.TypeList,
			},
			Attr_WorkspaceLocation: {
				Computed:    true,
				Description: "Workspace location.",
				Type:        schema.TypeMap,
			},
			Attr_WorkspaceName: {
				Computed:    true,
				Description: "Workspace name.",
				Type:        schema.TypeString,
			},
			Attr_WorkspaceStatus: {
				Computed:    true,
				Description: "Workspace status, active, critical, failed, provisioning.",
				Type:        schema.TypeString,
			},
			Attr_WorkspaceType: {
				Computed:    true,
				Description: "Workspace type, off-premises or on-premises.",
				Type:        schema.TypeString,
			},
		},
	}
}

func dataSourceIBMPIWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	wsData, err := client.Get(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(Attr_WorkspaceName, wsData.Name)
	d.Set(Attr_WorkspaceStatus, wsData.Status)
	d.Set(Attr_WorkspaceType, wsData.Type)
	d.Set(Attr_WorkspaceCapabilities, wsData.Capabilities)

	wsDetails := []map[string]interface{}{}
	detailsData := make(map[string]interface{})
	detailsData[Attr_CreationDate] = wsData.Details.CreationDate.String()
	detailsData[Attr_CRN] = *wsData.Details.Crn

	if wsData.Details.PowerEdgeRouter != nil {
		wsPowerEdge := map[string]interface{}{
			Attr_MigrationStatus: wsData.Details.PowerEdgeRouter.MigrationStatus,
			Attr_State:           *wsData.Details.PowerEdgeRouter.State,
			Attr_Type:            *wsData.Details.PowerEdgeRouter.Type,
		}
		detailsData[Attr_PowerEdgeRouter] = []map[string]interface{}{wsPowerEdge}
		wsDetails = append(wsDetails, detailsData)
	}
	if wsData.Details.NetworkSecurityGroups != nil {
		wsNSG := map[string]interface{}{
			Attr_State: *wsData.Details.NetworkSecurityGroups.State,
		}
		detailsData[Attr_NetworkSecurityGroups] = []map[string]interface{}{wsNSG}
		wsDetails = append(wsDetails, detailsData)
	}

	d.Set(Attr_WorkspaceDetails, wsDetails)
	wsLocation := map[string]interface{}{
		Attr_Region: *wsData.Location.Region,
		Attr_Type:   wsData.Location.Type,
		Attr_URL:    wsData.Location.URL,
	}
	d.Set(Attr_WorkspaceLocation, flex.Flatten(wsLocation))
	d.SetId(*wsData.ID)
	return nil
}
