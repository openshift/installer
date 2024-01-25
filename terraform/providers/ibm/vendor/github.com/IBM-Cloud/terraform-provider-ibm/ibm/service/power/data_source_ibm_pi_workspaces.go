// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package power

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	Workspaces = "workspaces"
)

func DatasourceIBMPIWorkspaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPIWorkspacesRead,
		Schema: map[string]*schema.Schema{
			Arg_CloudInstanceID: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			Workspaces: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						Attr_WorkspaceCapabilities: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Workspace Capabilities",
							Elem: &schema.Schema{
								Type: schema.TypeBool,
							},
						},
						Attr_WorkspaceDetails: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Workspace information",
						},
						Attr_WorkspaceID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace ID",
						},
						Attr_WorkspaceLocation: {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Workspace location",
						},
						Attr_WorkspaceName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace name",
						},
						Attr_WorkspaceStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace status",
						},
						Attr_WorkspaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace type",
						},
					},
				},
			},
		},
	}
}
func dataSourceIBMPIWorkspacesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Get(helpers.PICloudInstanceId).(string)
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	wsData, err := client.GetAll()
	if err != nil {
		return diag.FromErr(err)
	}
	workspaces := make([]map[string]interface{}, 0, len(wsData.Workspaces))
	for _, ws := range wsData.Workspaces {
		if ws != nil {
			workspace := map[string]interface{}{
				Attr_WorkspaceName:         ws.Name,
				Attr_WorkspaceID:           ws.ID,
				Attr_WorkspaceStatus:       ws.Status,
				Attr_WorkspaceType:         ws.Type,
				Attr_WorkspaceCapabilities: ws.Capabilities,
				Attr_WorkspaceDetails: map[string]interface{}{
					WorkspaceCreationDate: ws.Details.CreationDate.String(),
					WorkspaceCRN:          *ws.Details.Crn,
				},
				Attr_WorkspaceLocation: map[string]interface{}{
					WorkspaceRegion: *ws.Location.Region,
					WorkspaceType:   ws.Location.Type,
					WorkspaceUrl:    ws.Location.URL,
				},
			}
			workspaces = append(workspaces, workspace)
		}
	}
	var clientgenU, _ = uuid.GenerateUUID()
	d.SetId(clientgenU)
	d.Set(Workspaces, workspaces)
	return nil
}
