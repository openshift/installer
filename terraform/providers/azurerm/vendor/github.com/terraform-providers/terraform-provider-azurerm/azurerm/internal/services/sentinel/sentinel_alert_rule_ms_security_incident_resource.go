package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/securityinsight/mgmt/2019-01-01-preview/securityinsight"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	loganalyticsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sentinel/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSentinelAlertRuleMsSecurityIncident() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Read:   resourceSentinelAlertRuleMsSecurityIncidentRead,
		Update: resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate,
		Delete: resourceSentinelAlertRuleMsSecurityIncidentDelete,

		Importer: azSchema.ValidateResourceIDPriorToImportThen(func(id string) error {
			_, err := parse.AlertRuleID(id)
			return err
		}, importSentinelAlertRule(securityinsight.AlertRuleKindMicrosoftSecurityIncidentCreation)),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loganalyticsValidate.LogAnalyticsWorkspaceID,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"product_filter": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(securityinsight.MicrosoftCloudAppSecurity),
					string(securityinsight.AzureSecurityCenter),
					string(securityinsight.AzureActiveDirectoryIdentityProtection),
					string(securityinsight.AzureSecurityCenterforIoT),
					string(securityinsight.AzureAdvancedThreatProtection),
				}, false),
			},

			"severity_filter": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(securityinsight.High),
						string(securityinsight.Medium),
						string(securityinsight.Low),
						string(securityinsight.Informational),
					}, false),
				},
			},

			"alert_rule_template_guid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"display_name_filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true, // remove in 3.0
				MinItems:      1,
				ConflictsWith: []string{"text_whitelist"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"display_name_exclude_filter": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"text_whitelist": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true, // remove in 3.0
				MinItems:      1,
				ConflictsWith: []string{"display_name_filter"},
				Deprecated:    "this property has been renamed to display_name_filter to better match the SDK & API",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceID, err := loganalyticsParse.LogAnalyticsWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroup, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.WorkspaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
			}
		}

		id := alertRuleID(resp.Value)
		if id != nil && *id != "" {
			return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_ms_security_incident", *id)
		}
	}

	param := securityinsight.MicrosoftSecurityIncidentCreationAlertRule{
		Kind: securityinsight.KindMicrosoftSecurityIncidentCreation,
		MicrosoftSecurityIncidentCreationAlertRuleProperties: &securityinsight.MicrosoftSecurityIncidentCreationAlertRuleProperties{
			ProductFilter:    securityinsight.MicrosoftSecurityProductName(d.Get("product_filter").(string)),
			DisplayName:      utils.String(d.Get("display_name").(string)),
			Description:      utils.String(d.Get("description").(string)),
			Enabled:          utils.Bool(d.Get("enabled").(bool)),
			SeveritiesFilter: expandAlertRuleMsSecurityIncidentSeverityFilter(d.Get("severity_filter").(*schema.Set).List()),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.MicrosoftSecurityIncidentCreationAlertRuleProperties.AlertRuleTemplateName = utils.String(v.(string))
	}

	if dnf, ok := d.GetOk("display_name_filter"); ok {
		param.DisplayNamesFilter = utils.ExpandStringSlice(dnf.(*schema.Set).List())
	} else if dnf, ok := d.GetOk("text_whitelist"); ok {
		param.DisplayNamesFilter = utils.ExpandStringSlice(dnf.(*schema.Set).List())
	}

	if v, ok := d.GetOk("display_name_exclude_filter"); ok {
		param.DisplayNamesExcludeFilter = utils.ExpandStringSlice(v.(*schema.Set).List())
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.WorkspaceName, name)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
		}

		if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindMicrosoftSecurityIncidentCreation); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
		param.Etag = resp.Value.(securityinsight.MicrosoftSecurityIncidentCreationAlertRule).Etag
	}

	if _, err := client.CreateOrUpdate(ctx, workspaceID.ResourceGroup, "Microsoft.OperationalInsights", workspaceID.WorkspaceName, name, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleMsSecurityIncidentRead(d, meta)
}

func resourceSentinelAlertRuleMsSecurityIncidentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Sentinel Alert Rule Ms Security Incident %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Value, securityinsight.AlertRuleKindMicrosoftSecurityIncidentCreation); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	rule := resp.Value.(securityinsight.MicrosoftSecurityIncidentCreationAlertRule)

	d.Set("name", id.Name)

	workspaceId := loganalyticsParse.NewLogAnalyticsWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	d.Set("log_analytics_workspace_id", workspaceId.ID())
	if prop := rule.MicrosoftSecurityIncidentCreationAlertRuleProperties; prop != nil {
		d.Set("product_filter", string(prop.ProductFilter))
		d.Set("display_name", prop.DisplayName)
		d.Set("description", prop.Description)
		d.Set("enabled", prop.Enabled)
		d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)

		if err := d.Set("text_whitelist", utils.FlattenStringSlice(prop.DisplayNamesFilter)); err != nil {
			return fmt.Errorf(`setting "text_whitelist": %+v`, err)
		}
		if err := d.Set("display_name_filter", utils.FlattenStringSlice(prop.DisplayNamesFilter)); err != nil {
			return fmt.Errorf(`setting "display_name_filter": %+v`, err)
		}
		if err := d.Set("display_name_exclude_filter", utils.FlattenStringSlice(prop.DisplayNamesExcludeFilter)); err != nil {
			return fmt.Errorf(`setting "display_name_exclude_filter": %+v`, err)
		}
		if err := d.Set("severity_filter", flattenAlertRuleMsSecurityIncidentSeverityFilter(prop.SeveritiesFilter)); err != nil {
			return fmt.Errorf(`setting "severity_filter": %+v`, err)
		}
	}

	return nil
}

func resourceSentinelAlertRuleMsSecurityIncidentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, "Microsoft.OperationalInsights", id.WorkspaceName, id.Name); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	return nil
}

func expandAlertRuleMsSecurityIncidentSeverityFilter(input []interface{}) *[]securityinsight.AlertSeverity {
	result := make([]securityinsight.AlertSeverity, 0)

	for _, e := range input {
		result = append(result, securityinsight.AlertSeverity(e.(string)))
	}

	return &result
}

func flattenAlertRuleMsSecurityIncidentSeverityFilter(input *[]securityinsight.AlertSeverity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		output = append(output, string(e))
	}

	return output
}
