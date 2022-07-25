package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v2"
)

const affinityNegative = "negative"
const affinityPositive = "positive"

var affinityRuleSchema = map[string]*schema.Schema{
	"affinity": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Positive or negative affinity.",
		ForceNew:    true,
		ValidateDiagFunc: validateEnum(
			[]string{
				affinityPositive,
				affinityNegative,
			},
		),
		Default: affinityNegative,
	},
	"enforcing": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true VMs will fail to start if they cannot observe this affintiy group.",
		Default:     false,
	},
}

var affinityGroupSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"cluster_id": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "ID of the cluster to use for affinity group creation.",
		ForceNew:         true,
		ValidateDiagFunc: validateUUID,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The name of the affinity group",
		ForceNew:         true,
		ValidateDiagFunc: validateNonEmpty,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The description of the affinity group",
		ForceNew:    true,
	},
	"priority": {
		Type:        schema.TypeFloat,
		Optional:    true,
		Description: "Order in which the affinity group should be applied.",
		ForceNew:    true,
	},
	"enforcing": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "If set to true, VMs will fail to start if the affinity group cannot be observed.",
		ForceNew:    true,
		Default:     false,
	},
	"hosts_rule": {
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		MinItems: 0,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: affinityRuleSchema,
		},
	},
	"vms_rule": {
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		MinItems: 0,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: affinityRuleSchema,
		},
	},
}

func (p *provider) affinityGroupResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.affinityGroupCreate,
		ReadContext:   p.affinityGroupRead,
		DeleteContext: p.affinityGroupDelete,
		Schema:        affinityGroupSchema,
		Description:   "The ovirt_affinity_group resource creates affinity groups in oVirt.",
	}
}

func (p *provider) affinityGroupCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	params := ovirtclient.CreateAffinityGroupParams()

	if hostsRule, ok := data.GetOk("hosts_rule"); ok {
		rules := hostsRule.(*schema.Set).List()
		if len(rules) == 1 {
			rule := rules[0].(map[string]interface{})
			affinity := ovirtclient.Affinity(rule["affinity"].(string) == affinityPositive)
			enforcing := rule["enforcing"].(bool)
			_, err := params.WithHostsRuleParameters(true, affinity, enforcing)
			if err != nil {
				return errorToDiags("create affinity group", err)
			}
		}
	}
	if vmsRule, ok := data.GetOk("vms_rule"); ok {
		rules := vmsRule.(*schema.Set).List()
		if len(rules) == 1 {
			rule := rules[0].(map[string]interface{})
			affinity := ovirtclient.Affinity(rule["affinity"].(string) == affinityPositive)
			enforcing := rule["enforcing"].(bool)
			_, err := params.WithVMsRuleParameters(true, affinity, enforcing)
			if err != nil {
				return errorToDiags("create affinity group", err)
			}
		}
	}

	clusterID := ovirtclient.ClusterID(data.Get("cluster_id").(string))
	name := data.Get("name").(string)
	if description, ok := data.GetOk("description"); ok {
		var err error
		params, err = params.WithDescription(description.(string))
		if err != nil {
			return errorToDiags("add description to affinity group", err)
		}
	}

	ag, err := client.CreateAffinityGroup(clusterID, name, params)
	if err != nil {
		return errorToDiags("create affinity group", err)
	}

	return affinityGroupToData(ag, data)
}

func affinityGroupToData(ag ovirtclient.AffinityGroup, data *schema.ResourceData) diag.Diagnostics {
	data.SetId(string(ag.ID()))
	diags := diag.Diagnostics{}
	if err := data.Set("name", ag.Name()); err != nil {
		diags = append(diags, errorToDiag("set name", err))
	}
	if err := data.Set("description", ag.Description()); err != nil {
		diags = append(diags, errorToDiag("set description", err))
	}

	if err := data.Set("enforcing", ag.Enforcing()); err != nil {
		diags = append(diags, errorToDiag("set enforcing", err))
	}

	diags = appendDiags(diags, "set hosts_rule", data.Set("hosts_rule", affinityRuleToData(ag.HostsRule())))
	diags = appendDiags(diags, "set vms_rule", data.Set("vms_rule", affinityRuleToData(ag.VMsRule())))

	return diags
}

func affinityRuleToData(rule ovirtclient.AffinityVMsRule) []map[string]interface{} {
	var result []map[string]interface{}
	if rule.Enabled() {
		result = append(
			result, map[string]interface{}{
				"affinity":  convertAffinity(rule.Affinity()),
				"enforcing": rule.Enforcing(),
			},
		)
	}
	return result
}

func convertAffinity(affinity ovirtclient.Affinity) string {
	var a string
	if affinity == ovirtclient.AffinityPositive {
		a = affinityPositive
	} else {
		a = affinityNegative
	}
	return a
}

func (p *provider) affinityGroupRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	id := data.Id()
	clusterID, ok := data.GetOk("cluster_id")
	if !ok {
		data.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       "Cluster ID not set",
				Detail:        "The cluster_id field is not set on the affinity group, cannot read from the oVirt engine.",
				AttributePath: nil,
			},
		}
	}
	ag, err := client.GetAffinityGroup(
		ovirtclient.ClusterID(clusterID.(string)),
		ovirtclient.AffinityGroupID(id),
	)
	if err != nil {
		return errorToDiags(fmt.Sprintf("read affinity group %s", id), err)
	}
	return affinityGroupToData(ag, data)
}

func (p *provider) affinityGroupDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	clusterID, ok := data.GetOk("cluster_id")
	if !ok {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity:      diag.Warning,
				Summary:       "Cluster ID not set",
				Detail:        "The cluster_id field is not set on the affinity group, cannot remove from the oVirt engine.",
				AttributePath: nil,
			},
		}
	}
	if err := client.RemoveAffinityGroup(
		ovirtclient.ClusterID(clusterID.(string)),
		ovirtclient.AffinityGroupID(data.Id()),
	); err != nil {
		return errorToDiags(
			fmt.Sprintf("delete affinity group %s", data.Id()),
			err,
		)
	}
	return nil
}
