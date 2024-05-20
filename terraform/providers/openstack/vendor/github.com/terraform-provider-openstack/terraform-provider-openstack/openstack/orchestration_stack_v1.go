package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/orchestration/v1/stacks"
)

func buildTE(t map[string]interface{}) stacks.TE {
	log.Printf("[DEBUG] Start to build TE structure")
	te := stacks.TE{}
	if t["Bin"] != nil {
		te.Bin = []byte(t["Bin"].(string))
	}
	if t["URL"] != nil {
		te.URL = t["URL"].(string)
	}
	if t["Files"] != nil {
		te.Files = t["Files"].(map[string]string)
	}
	log.Printf("[DEBUG] TE structure builded")
	return te
}

func buildTemplateOpts(d *schema.ResourceData) *stacks.Template {
	log.Printf("[DEBUG] Start building TemplateOpts")
	template := &stacks.Template{}
	template.TE = buildTE(d.Get("template_opts").(map[string]interface{}))
	log.Printf("[DEBUG] Return TemplateOpts")
	return template
}

func buildEnvironmentOpts(d *schema.ResourceData) *stacks.Environment {
	log.Printf("[DEBUG] Start building EnvironmentOpts")
	environment := &stacks.Environment{}
	if d.Get("environment_opts") != nil {
		t := d.Get("environment_opts").(map[string]interface{})
		environment.TE = buildTE(t)
		log.Printf("[DEBUG] Return EnvironmentOpts")
		return environment
	}
	return nil
}

func orchestrationStackV1StateRefreshFunc(client *gophercloud.ServiceClient, stackID string, isdelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Refresh Stack status %s", stackID)
		stack, err := stacks.Find(client, stackID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok && isdelete {
				return stack, "DELETE_COMPLETE", nil
			}

			return nil, "", err
		}

		if strings.Contains(stack.Status, "FAILED") {
			return stack, stack.Status, fmt.Errorf("The stack is in error status. " +
				"Please check with your cloud admin or check the orchestration " +
				"API logs to see why this error occurred.")
		}

		return stack, stack.Status, nil
	}
}
