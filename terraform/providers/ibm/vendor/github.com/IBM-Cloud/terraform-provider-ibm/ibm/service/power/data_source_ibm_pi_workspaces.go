// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DatasourceIBMPIWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIWorkspacesRead,
		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_CloudInstanceID: {
				Description:  "The GUID of the service instance associated with an account.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},

			// Attributes
			Attr_Workspaces: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						Attr_WorkspaceID: {
							Computed:    true,
							Description: "Workspace ID.",
							Type:        schema.TypeString,
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
				},
			},
		},
	}
}

func dataSourceIBMPIWorkspacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(Arg_CloudInstanceID).(string)
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	wsData, err := client.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	workspaces := make([]map[string]interface{}, 0, len(wsData.Workspaces))
	for _, ws := range wsData.Workspaces {
		if ws != nil {
			wsDetails := []map[string]interface{}{}
			detailsData := make(map[string]interface{})
			detailsData[Attr_CreationDate] = ws.Details.CreationDate.String()
			detailsData[Attr_CRN] = *ws.Details.Crn

			if ws.Details.PowerEdgeRouter != nil {
				wsPowerEdge := map[string]interface{}{
					Attr_MigrationStatus: ws.Details.PowerEdgeRouter.MigrationStatus,
					Attr_State:           *ws.Details.PowerEdgeRouter.State,
					Attr_Type:            *ws.Details.PowerEdgeRouter.Type,
				}
				detailsData[Attr_PowerEdgeRouter] = []map[string]interface{}{wsPowerEdge}
				wsDetails = append(wsDetails, detailsData)
			}
			if ws.Details.NetworkSecurityGroups != nil {
				wsNSG := map[string]interface{}{
					Attr_State: *ws.Details.NetworkSecurityGroups.State,
				}
				detailsData[Attr_NetworkSecurityGroups] = []map[string]interface{}{wsNSG}
				wsDetails = append(wsDetails, detailsData)
			}

			workspace := map[string]interface{}{
				Attr_WorkspaceCapabilities: ws.Capabilities,
				Attr_WorkspaceDetails:      wsDetails,
				Attr_WorkspaceID:           ws.ID,
				Attr_WorkspaceLocation: map[string]interface{}{
					Attr_Region: *ws.Location.Region,
					Attr_Type:   ws.Location.Type,
					Attr_URL:    ws.Location.URL,
				},
				Attr_WorkspaceName:   ws.Name,
				Attr_WorkspaceStatus: ws.Status,
				Attr_WorkspaceType:   ws.Type,
			}
			workspaces = append(workspaces, workspace)
		}
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Attr_Workspaces, workspaces)
	return nil
}
