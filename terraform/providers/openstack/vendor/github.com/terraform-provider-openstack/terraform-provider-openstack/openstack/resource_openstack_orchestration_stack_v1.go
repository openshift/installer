package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
)

func resourceOrchestrationStackV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrchestrationStackV1Create,
		ReadContext:   resourceOrchestrationStackV1Read,
		UpdateContext: resourceOrchestrationStackV1Update,
		DeleteContext: resourceOrchestrationStackV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"template_opts": {
				Type:     schema.TypeMap,
				Required: true,
			},

			"disable_rollback": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"environment_opts": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Below are schemas for stack read
			"capabilities": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"creation_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"updated_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"notification_topics": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"outputs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"output_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"output_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"status_reason": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"template_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
		},
	}
}

func resourceOrchestrationStackV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Prepare for create openstack_orchestration_stack_v1")
	config := meta.(*Config)
	orchestrationClient, err := config.OrchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack Orchestration client: %s", err)
	}
	createOpts := &stacks.CreateOpts{
		Name:         d.Get("name").(string),
		TemplateOpts: buildTemplateOpts(d),
	}
	if d.Get("disable_rollback") != nil {
		disableRollback := d.Get("disable_rollback").(bool)
		createOpts.DisableRollback = &disableRollback
	}
	env := buildEnvironmentOpts(d)
	if env != nil {
		createOpts.EnvironmentOpts = env
	}
	if d.Get("parameters") != nil {
		createOpts.Parameters = d.Get("parameters").(map[string]interface{})
	}
	if d.Get("tags") != nil {
		t := d.Get("tags").([]interface{})
		tags := make([]string, len(t))
		for _, tag := range t {
			tags = append(tags, tag.(string))
		}
		createOpts.Tags = tags
	}
	if d.Get("timeout") != nil {
		createOpts.Timeout = d.Get("timeout").(int)
	}

	log.Printf("[DEBUG] Creating openstack_orchestration_stack_v1")
	stack, err := stacks.Create(orchestrationClient, createOpts).Extract()
	if err != nil {
		log.Printf("[DEBUG] openstack_orchestration_stack_v1 error occurred during Create: %s", err)
		return diag.Errorf("Error creating openstack_orchestration_stack_v1: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE_IN_PROGRESS", "INIT_COMPLETE"},
		Target:     []string{"CREATE_COMPLETE", "UPDATE_COMPLETE", "UPDATE_IN_PROGRESS"},
		Refresh:    orchestrationStackV1StateRefreshFunc(orchestrationClient, stack.ID, false),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_orchestration_stack_v1 %s to become ready: %s", stack.ID, err)
	}

	// Store the ID now
	d.SetId(stack.ID)
	log.Printf("[INFO] openstack_orchestration_stack_v1 %s create complete", stack.ID)

	return resourceOrchestrationStackV1Read(ctx, d, meta)
}

func resourceOrchestrationStackV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	orchestrationClient, err := config.OrchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack Orchestration client: %s", err)
	}

	log.Printf("[DEBUG] Fetch openstack_orchestration_stack_v1 information: %s", d.Id())
	stack, err := stacks.Find(orchestrationClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_orchestration_stack_v1"))
	}

	d.Set("name", stack.Name)
	d.Set("capabilities", stack.Capabilities)
	d.Set("description", stack.Description)
	d.Set("disable_rollback", stack.DisableRollback)
	d.Set("notification_topics", stack.NotificationTopics)
	d.Set("status", stack.Status)
	d.Set("status_reason", stack.StatusReason)
	d.Set("template_description", stack.TemplateDescription)
	d.Set("timeout", stack.Timeout)

	// Set the outputs
	outputs := make([]map[string]interface{}, 0, len(stack.Outputs))
	for _, o := range stack.Outputs {
		output := make(map[string]interface{})
		output["description"] = o["description"]
		output["output_key"] = o["output_key"]
		output["output_value"] = o["output_value"]

		outputs = append(outputs, output)
	}
	d.Set("outputs", outputs)

	params := stack.Parameters
	if stack.Parameters != nil {
		removeList := []string{"OS::project_id", "OS::stack_id", "OS::stack_name"}
		for _, v := range removeList {
			_, ok := params[v]
			if ok {
				delete(params, v)
			}
		}
		d.Set("parameters", stack.Parameters)
	}

	if len(stack.Tags) > 0 {
		var tags = []string{}
		for _, v := range stack.Tags {
			if v != "" {
				tags = append(tags, v)
			}
		}
		d.Set("tags", tags)
	}

	if err := d.Set("creation_time", stack.CreationTime.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_orchestration_stack_v1 creation_time: %s", err)
	}
	if err := d.Set("updated_time", stack.UpdatedTime.Format(time.RFC3339)); err != nil {
		log.Printf("[DEBUG] Unable to set openstack_orchestration_stack_v1 updated_at: %s", err)
	}

	log.Printf("[DEBUG] openstack_orchestration_stack_v1 information fetched: %s", d.Id())
	return nil
}

func resourceOrchestrationStackV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Prepare information for update openstack_orchestration_stack_v1")

	config := meta.(*Config)
	orchestrationClient, err := config.OrchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack Orchestration client: %s", err)
	}

	updateOpts := &stacks.UpdateOpts{
		TemplateOpts: buildTemplateOpts(d),
	}
	env := buildEnvironmentOpts(d)
	if env != nil {
		updateOpts.EnvironmentOpts = env
	}
	if d.Get("parameters") != nil {
		updateOpts.Parameters = d.Get("parameters").(map[string]interface{})
	}
	if d.Get("timeout") != nil {
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.Get("tags") != nil {
		t := d.Get("tags").([]interface{})
		tags := make([]string, len(t))
		for _, tag := range t {
			tags = append(tags, tag.(string))
		}
		updateOpts.Tags = tags
	}

	stack, err := stacks.Find(orchestrationClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving openstack_orchestration_stack_v1 %s before Update:  %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Updating openstack_orchestration_stack_v1")
	result := stacks.Update(orchestrationClient, stack.Name, d.Id(), updateOpts)
	if result.Err != nil {
		return diag.Errorf("Error updating openstack_orchestration_stack_v1 %s: %s", d.Id(), result.Err)
	}

	log.Printf("[INFO] openstack_orchestration_stack_v1 %s update complete", d.Id())
	return resourceOrchestrationStackV1Read(ctx, d, meta)
}

func resourceOrchestrationStackV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Prepare for delete openstack_orchestration_stack_v1")
	config := meta.(*Config)
	orchestrationClient, err := config.OrchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack Orchestration client: %s", err)
	}

	stack, err := stacks.Find(orchestrationClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_orchestration_stack_v1"))
	}

	if stack.Status != "DELETE_IN_PROGRESS" {
		log.Printf("[DEBUG] Deleting openstack_orchestration_stack_v1: %s", d.Id())
		if err := stacks.Delete(orchestrationClient, stack.Name, d.Id()).ExtractErr(); err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_orchestration_stack_v1"))
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETE_IN_PROGRESS"},
		Target:     []string{"DELETE_COMPLETE"},
		Refresh:    orchestrationStackV1StateRefreshFunc(orchestrationClient, d.Id(), true),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_orchestration_stack_v1 %s to Delete:  %s", d.Id(), err)
	}

	log.Printf("[INFO] openstack_orchestration_stack_v1 %s delete complete", d.Id())
	return nil
}
