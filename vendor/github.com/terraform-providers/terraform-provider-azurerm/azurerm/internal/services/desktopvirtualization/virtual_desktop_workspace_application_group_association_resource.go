package desktopvirtualization

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/desktopvirtualization/mgmt/2019-12-10-preview/desktopvirtualization"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVirtualDesktopWorkspaceApplicationGroupAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate,
		Read:   resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead,
		Delete: resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.WorkspaceApplicationGroupAssociationID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    migration.WorkspaceApplicationGroupAssociationUpgradeV0Schema().CoreConfigSchema().ImpliedType(),
				Upgrade: migration.WorkspaceApplicationGroupAssociationUpgradeV0ToV1,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"application_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ApplicationGroupID,
			},
		},
	}
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Virtual Desktop Workspace <-> Application Group Association creation.")
	workspaceId, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return err
	}
	applicationGroupId, err := parse.ApplicationGroupID(d.Get("application_group_id").(string))
	if err != nil {
		return err
	}
	associationId := parse.NewWorkspaceApplicationGroupAssociationId(*workspaceId, *applicationGroupId).ID()

	locks.ByName(workspaceId.Name, workspaceResourceType)
	defer locks.UnlockByName(workspaceId.Name, workspaceResourceType)

	locks.ByName(applicationGroupId.Name, applicationGroupType)
	defer locks.UnlockByName(applicationGroupId.Name, applicationGroupType)

	workspace, err := client.Get(ctx, workspaceId.ResourceGroup, workspaceId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(workspace.Response) {
			return fmt.Errorf("Virtual Desktop Workspace %q (Resource Group %q) was not found", workspaceId.Name, workspaceId.ResourceGroup)
		}

		return fmt.Errorf("retrieving Virtual Desktop Workspace for Association %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, err)
	}

	applicationGroupAssociations := []string{}
	if props := workspace.WorkspaceProperties; props != nil && props.ApplicationGroupReferences != nil {
		applicationGroupAssociations = *props.ApplicationGroupReferences
	}

	applicationGroupIdStr := applicationGroupId.ID()
	if associationExists(workspace.WorkspaceProperties, applicationGroupIdStr) {
		return tf.ImportAsExistsError("azurerm_virtual_desktop_workspace_application_group_association", associationId)
	}
	applicationGroupAssociations = append(applicationGroupAssociations, applicationGroupIdStr)

	workspace.WorkspaceProperties.ApplicationGroupReferences = &applicationGroupAssociations

	if _, err = client.CreateOrUpdate(ctx, workspaceId.ResourceGroup, workspaceId.Name, workspace); err != nil {
		return fmt.Errorf("creating association between Virtual Desktop Workspace %q (Resource Group %q) and Application Group %q (Resource Group %q): %+v", workspaceId.Name, workspaceId.ResourceGroup, applicationGroupId.Name, applicationGroupId.ResourceGroup, err)
	}

	d.SetId(associationId)
	return resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead(d, meta)
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name)
	if err != nil {
		if utils.ResponseWasNotFound(workspace.Response) {
			log.Printf("[DEBUG] Virtual Desktop Workspace %q was not found in Resource Group %q - removing from state!", id.Workspace.Name, id.Workspace.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Desktop Desktop Workspace %q (Resource Group %q): %+v", id.Workspace.Name, id.Workspace.ResourceGroup, err)
	}

	applicationGroupId := id.ApplicationGroup.ID()
	exists := associationExists(workspace.WorkspaceProperties, applicationGroupId)
	if !exists {
		log.Printf("[DEBUG] Association between Virtual Desktop Workspace %q (Resource Group %q) and Application Group %q (Resource Group %q) was not found - removing from state!", id.Workspace.Name, id.Workspace.ResourceGroup, id.ApplicationGroup.Name, id.ApplicationGroup.ResourceGroup)
		d.SetId("")
		return nil
	}

	d.Set("workspace_id", id.Workspace.ID())
	d.Set("application_group_id", applicationGroupId)

	return nil
}

func resourceVirtualDesktopWorkspaceApplicationGroupAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DesktopVirtualization.WorkspacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.WorkspaceApplicationGroupAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Workspace.Name, workspaceResourceType)
	defer locks.UnlockByName(id.Workspace.Name, workspaceResourceType)

	locks.ByName(id.ApplicationGroup.Name, applicationGroupType)
	defer locks.UnlockByName(id.ApplicationGroup.Name, applicationGroupType)

	workspace, err := client.Get(ctx, id.Workspace.ResourceGroup, id.Workspace.Name)
	if err != nil {
		if utils.ResponseWasNotFound(workspace.Response) {
			return fmt.Errorf("Virtual Desktop Workspace %q (Resource Group %q) was not found", id.Workspace.Name, id.Workspace.ResourceGroup)
		}

		return fmt.Errorf("retrieving Virtual Desktop Workspace %q (Resource Group %q): %+v", id.Workspace.Name, id.Workspace.ResourceGroup, err)
	}

	applicationGroupReferences := []string{}
	applicationGroupId := id.ApplicationGroup.ID()
	if workspace.WorkspaceProperties != nil && workspace.WorkspaceProperties.ApplicationGroupReferences != nil {
		for _, referenceId := range *workspace.WorkspaceProperties.ApplicationGroupReferences {
			if strings.EqualFold(referenceId, applicationGroupId) {
				continue
			}

			applicationGroupReferences = append(applicationGroupReferences, referenceId)
		}
	}

	workspace.WorkspaceProperties.ApplicationGroupReferences = &applicationGroupReferences

	if _, err = client.CreateOrUpdate(ctx, id.Workspace.ResourceGroup, id.Workspace.Name, workspace); err != nil {
		return fmt.Errorf("removing association between Virtual Desktop Workspace %q (Resource Group %q) and Application Group %q (Resource Group %q): %+v", id.Workspace.Name, id.Workspace.ResourceGroup, id.ApplicationGroup.Name, id.ApplicationGroup.ResourceGroup, err)
	}

	return nil
}

func associationExists(props *desktopvirtualization.WorkspaceProperties, applicationGroupId string) bool {
	if props == nil || props.ApplicationGroupReferences == nil {
		return false
	}

	for _, id := range *props.ApplicationGroupReferences {
		if strings.EqualFold(id, applicationGroupId) {
			return true
		}
	}

	return false
}
