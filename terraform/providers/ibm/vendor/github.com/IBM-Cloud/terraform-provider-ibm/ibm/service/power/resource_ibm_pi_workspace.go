package power

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMPIWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMPIWorkspaceCreate,
		DeleteContext: resourceIBMPIWorkspaceDelete,
		ReadContext:   resourceIBMPIWorkspaceRead,
		UpdateContext: resourceIBMPIWorkspaceUpdate,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
				return flex.ResourcePowerUserTagsCustomizeDiff(diff)
			},
		),

		Schema: map[string]*schema.Schema{
			// Arguments
			Arg_Datacenter: {
				Description:  "Target location or environment to create the resource instance.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Name: {
				Description:  "A descriptive name used to identify the workspace.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_Plan: {
				Default:      Public,
				Description:  "Plan associated with the offering; Valid values are public or private.",
				ForceNew:     true,
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{Private, Public}),
			},
			Arg_ResourceGroupID: {
				Description:  "The ID of the resource group where you want to create the workspace. You can retrieve the value from data source ibm_resource_group.",
				ForceNew:     true,
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
			Arg_UserTags: {
				Computed:    true,
				Description: "List of user tags attached to the resource.",
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Set:         schema.HashString,
				Type:        schema.TypeSet,
			},

			// Attributes
			Attr_CRN: {
				Computed:    true,
				Description: "The Workspace crn.",
				Type:        schema.TypeString,
			},
			Attr_WorkspaceDetails: {
				Computed:    true,
				Deprecated:  "This field is deprecated, use crn instead.",
				Description: "Workspace information.",
				Type:        schema.TypeMap,
			},
		},
	}
}

func resourceIBMPIWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get(Arg_Name).(string)
	datacenter := d.Get(Arg_Datacenter).(string)
	resourceGroup := d.Get(Arg_ResourceGroupID).(string)
	plan := d.Get(Arg_Plan).(string)

	// No need for cloudInstanceID because we are creating a workspace
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, "")
	controller, _, err := client.Create(name, datacenter, resourceGroup, plan)
	if err != nil {
		log.Printf("[DEBUG] create workspace failed %v", err)
		return diag.FromErr(err)
	}

	cloudInstanceID := *controller.GUID
	d.SetId(cloudInstanceID)

	_, err = waitForResourceWorkspaceCreate(ctx, client, cloudInstanceID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	// Add user tags for newly created workspace
	if tags, ok := d.GetOk(Arg_UserTags); ok {
		if len(flex.FlattenSet(tags.(*schema.Set))) > 0 {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *controller.CRN, "", UserTagType)
			if err != nil {
				log.Printf("Error on creation of workspace (%s) pi_user_tags: %s", *controller.CRN, err)
			}
		}
	}
	return resourceIBMPIWorkspaceRead(ctx, d, meta)
}

func waitForResourceWorkspaceCreate(ctx context.Context, client *instance.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_InProgress, State_Inactive, State_Provisioning},
		Target:     []string{State_Active},
		Refresh:    isIBMPIWorkspaceCreateRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Minute,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIWorkspaceCreateRefreshFunc(client *instance.IBMPIWorkspacesClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		controller, _, err := client.GetRC(id)
		if err != nil {
			return nil, "", err
		}
		if *controller.State == State_Failed {
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
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	controller, _, err := client.GetRC(cloudInstanceID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set(Arg_Name, controller.Name)
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *controller.CRN, "", UserTagType)
	if err != nil {
		log.Printf("Error on get of workspace (%s) pi_user_tags: %s", cloudInstanceID, err)
	}
	d.Set(Arg_UserTags, tags)

	d.Set(Attr_CRN, controller.CRN)

	// Deprecated Workspace Details Set
	wsDetails := map[string]interface{}{
		Attr_CreationDate: controller.CreatedAt,
		Attr_CRN:          controller.CRN,
	}
	d.Set(Attr_WorkspaceDetails, flex.Flatten(wsDetails))

	return nil
}

func resourceIBMPIWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).IBMPISession()
	if err != nil {
		return diag.FromErr(err)
	}

	cloudInstanceID := d.Id()
	client := instance.NewIBMPIWorkspacesClient(ctx, sess, cloudInstanceID)
	response, err := client.Delete(cloudInstanceID)
	if err != nil && response != nil && response.StatusCode == 410 {
		return nil
	}
	_, err = waitForResourceWorkspaceDelete(ctx, client, cloudInstanceID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

func waitForResourceWorkspaceDelete(ctx context.Context, client *instance.IBMPIWorkspacesClient, id string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{State_InProgress, State_Inactive, State_Active},
		Target:     []string{State_Removed, State_PendingReclamation},
		Refresh:    isIBMPIResourceDeleteRefreshFunc(client, id),
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Second,
		Timeout:    timeout,
	}
	return stateConf.WaitForStateContext(ctx)
}

func isIBMPIResourceDeleteRefreshFunc(client *instance.IBMPIWorkspacesClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		controller, response, err := client.GetRC(id)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return controller, State_Active, nil
			}
			return nil, "", err
		}
		if controller == nil {
			return controller, State_Removed, nil
		} else {
			if *controller.State == State_Failed {
				return controller, *controller.State, fmt.Errorf("[ERROR] The resource instance %s failed to delete", id)
			}
			return controller, *controller.State, nil
		}
	}
}

func resourceIBMPIWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange(Arg_UserTags) {
		if crn, ok := d.GetOk(Attr_CRN); ok {
			oldList, newList := d.GetChange(Arg_UserTags)
			err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, crn.(string), "", UserTagType)
			if err != nil {
				log.Printf("Error on update of workspace (%s) pi_user_tags: %s", crn, err)
			}
		}
	}
	return resourceIBMPIWorkspaceRead(ctx, d, meta)
}
