package power

import (
	"context"
	"fmt"
	"log"
	"time"

	st "github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMPIWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIWorkspaceCreate,
		ReadContext:   resourceIBMPIWorkspaceRead,
		DeleteContext: resourceIBMPIWorkspaceDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			PIWorkspaceName: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A descriptive name used to identify the workspace.",
			},
			PIWorkspaceDatacenter: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target location or environment to create the resource instance.",
			},
			PIWorkspaceResourceGroup: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the resource group where you want to create the workspace. You can retrieve the value from data source ibm_resource_group.",
			},
			PIWorkspacePlan: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Plan associated with the offering; Valid values are public or private.",
			},
		},
	}
}

func resourceIBMPIWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(PIWorkspaceName).(string)
	datacenter := d.Get(PIWorkspaceDatacenter).(string)
	resourceGroup := d.Get(PIWorkspaceResourceGroup).(string)
	plan := d.Get(PIWorkspacePlan).(string)

	// No need for cloudInstanceID because we are creating a workspace
	client := st.NewIBMPIWorkspacesClient(ctx, sess, "")
	controller, _, err := client.Create(name, datacenter, resourceGroup, plan)
	if err != nil {
		log.Printf("[DEBUG] create workspace failed %v", err)
		return diag.FromErr(err)
	}

	d.SetId(*controller.GUID)
	_, err = waitForResourceInstanceCreate(ctx, client, *controller.GUID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIBMPIWorkspaceRead(ctx, d, meta)
}

func waitForResourceInstanceCreate(ctx context.Context, client *st.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in progress", "inactive", "provisioning"},
		Target:     []string{"active"},
		Refresh:    isIBMPIWorkspaceCreateRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIWorkspaceCreateRefreshFunc(client *st.IBMPIWorkspacesClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		controller, _, err := client.GetRC(id)
		if err != nil {
			return nil, "", err
		}
		if *controller.State == "failed" {
			return controller, *controller.State, fmt.Errorf("[ERROR] The resource instance %s failed to create", id)
		}
		return controller, *controller.State, nil
	}
}

func resourceIBMPIWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// session
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Id()
	client := st.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	controller, _, err := client.GetRC(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(PIWorkspaceName, controller.Name)

	return nil
}

func resourceIBMPIWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Id()
	client := st.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	response, err := client.Delete(cloudInstanceID)
	if err != nil && response != nil && response.StatusCode == 410 {
		return nil
	}
	_, err = waitForResourceInstanceDelete(ctx, client, cloudInstanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func waitForResourceInstanceDelete(ctx context.Context, client *st.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"in progress", "inactive", "active"},
		Target:     []string{"removed", "pending_reclamation"},
		Refresh:    isIBMPIResourceDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Second,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIResourceDeleteRefreshFunc(client *st.IBMPIWorkspacesClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		controller, response, err := client.GetRC(id)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return controller, "active", nil
			}
			return nil, "", err
		}
		if controller == nil {
			return controller, "removed", nil
		} else {
			if *controller.State == "failed" {
				return controller, *controller.State, fmt.Errorf("[ERROR] The resource instance %s failed to delete", id)
			}
			return controller, *controller.State, nil
		}
	}
}
