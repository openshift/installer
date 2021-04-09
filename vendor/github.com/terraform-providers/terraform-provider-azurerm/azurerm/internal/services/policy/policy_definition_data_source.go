package policy

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicyDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ExactlyOneOf: []string{"name", "display_name"},
			},

			"management_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"management_group_name"},
				Deprecated:    "Deprecated in favour of `management_group_name`", // TODO -- remove this in next major version
			},

			"management_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"management_group_id"},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"policy_rule": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	displayName := d.Get("display_name").(string)
	name := d.Get("name").(string)
	managementGroupName := ""
	if v, ok := d.GetOk("management_group_name"); ok {
		managementGroupName = v.(string)
	}
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupName = v.(string)
	}

	var policyDefinition policy.Definition
	var err error
	// one of display_name and name must be non-empty, this is guaranteed by schema
	if displayName != "" {
		policyDefinition, err = getPolicyDefinitionByDisplayName(ctx, client, displayName, managementGroupName)
		if err != nil {
			return fmt.Errorf("reading Policy Definition (Display Name %q): %+v", displayName, err)
		}
	}
	if name != "" {
		policyDefinition, err = getPolicyDefinitionByName(ctx, client, name, managementGroupName)
		if err != nil {
			return fmt.Errorf("reading Policy Definition %q: %+v", name, err)
		}
	}

	d.SetId(*policyDefinition.ID)
	d.Set("name", policyDefinition.Name)
	d.Set("display_name", policyDefinition.DisplayName)
	d.Set("description", policyDefinition.Description)
	d.Set("type", policyDefinition.Type)
	d.Set("policy_type", policyDefinition.PolicyType)

	policyRule := policyDefinition.PolicyRule.(map[string]interface{})
	if policyRuleStr := flattenJSON(policyRule); policyRuleStr != "" {
		d.Set("policy_rule", policyRuleStr)
	} else {
		return fmt.Errorf("flattening Policy Definition Rule %q: %+v", name, err)
	}

	if metadataStr := flattenJSON(policyDefinition.Metadata); metadataStr != "" {
		d.Set("metadata", metadataStr)
	}

	if parametersStr, err := flattenParameterDefinitionsValueToString(policyDefinition.Parameters); err == nil {
		d.Set("parameters", parametersStr)
	} else {
		return fmt.Errorf("failed to flatten Policy Parameters %q: %+v", name, err)
	}

	return nil
}
