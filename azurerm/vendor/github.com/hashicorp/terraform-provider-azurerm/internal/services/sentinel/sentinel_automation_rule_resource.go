package sentinel

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	securityinsight "github.com/tombuildsstuff/kermit/sdk/securityinsights/2022-10-01-preview/securityinsights"
)

func resourceSentinelAutomationRule() *pluginsdk.Resource {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"order": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 1000),
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"triggers_on": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(securityinsight.TriggersOnIncidents),
			ValidateFunc: validation.StringInSlice([]string{
				string(securityinsight.TriggersOnIncidents),
				string(securityinsight.TriggersOnAlerts),
			}, false),
		},

		"triggers_when": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(securityinsight.TriggersWhenCreated),
			ValidateFunc: validation.StringInSlice([]string{
				string(securityinsight.TriggersWhenCreated),
				string(securityinsight.TriggersWhenUpdated),
			}, false),
		},

		"expiration": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"condition_json": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// We can't use the pluginsdk.SuppressJsonDiff here as the "condition_json" is always an array, while that function assume its input is an object.
			// Once https://github.com/hashicorp/terraform-plugin-sdk/pull/1102 is merged, we can switch to pluginsdk.SuppressJsonDiff.
			DiffSuppressFunc: func(_, old, new string, _ *pluginsdk.ResourceData) bool {
				return utils.NormalizeJson(old) == utils.NormalizeJson(new)
			},
			ValidateFunc: validation.StringIsJSON,
		},

		"action_incident": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityinsight.IncidentStatusActive),
							string(securityinsight.IncidentStatusClosed),
							string(securityinsight.IncidentStatusNew),
						}, false),
					},

					"classification": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityinsight.IncidentClassificationUndetermined),
							string(securityinsight.IncidentClassificationBenignPositive) + "_" + string(securityinsight.IncidentClassificationReasonSuspiciousButExpected),
							string(securityinsight.IncidentClassificationFalsePositive) + "_" + string(securityinsight.IncidentClassificationReasonIncorrectAlertLogic),
							string(securityinsight.IncidentClassificationFalsePositive) + "_" + string(securityinsight.IncidentClassificationReasonInaccurateData),
							string(securityinsight.IncidentClassificationTruePositive) + "_" + string(securityinsight.IncidentClassificationReasonSuspiciousActivity),
						}, false),
					},

					"classification_comment": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"owner_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"severity": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityinsight.IncidentSeverityHigh),
							string(securityinsight.IncidentSeverityInformational),
							string(securityinsight.IncidentSeverityLow),
							string(securityinsight.IncidentSeverityMedium),
						}, false),
					},
				},
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},

		"action_playbook": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"logic_app_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"tenant_id": {
						Type: pluginsdk.TypeString,
						// We'll use the current tenant id if this property is absent.
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
			AtLeastOneOf: []string{"action_incident", "action_playbook"},
		},
	}

	if !features.FourPointOhBeta() {
		schema["condition"] = &pluginsdk.Schema{
			Deprecated: "This is deprecated in favor of `condition_json`",
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"property": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountAadTenantID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountAadUserID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountNTDomain),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountObjectGUID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountPUID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountSid),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAccountUPNSuffix),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAzureResourceResourceID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyAzureResourceSubscriptionID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyCloudApplicationAppName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyDNSDomainName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyFileDirectory),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyFileHashValue),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyFileName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyHostAzureID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyHostNTDomain),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyHostName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyHostNetBiosName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyHostOSVersion),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIPAddress),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentDescription),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentProviderName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentRelatedAnalyticRuleIds),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentSeverity),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentStatus),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentTactics),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIncidentTitle),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceModel),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceOperatingSystem),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceType),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyIoTDeviceVendor),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryAction),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageDeliveryLocation),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageP1Sender),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageP2Sender),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageRecipient),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageSenderIP),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailMessageSubject),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailboxDisplayName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailboxPrimaryAddress),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMailboxUPN),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMalwareCategory),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyMalwareName),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyProcessCommandLine),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyProcessID),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyRegistryKey),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyRegistryValueData),
							string(securityinsight.AutomationRulePropertyConditionSupportedPropertyURL),
						}, false),
					},

					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorContains),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorEndsWith),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorEquals),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorNotContains),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorNotEndsWith),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorNotEquals),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorNotStartsWith),
							string(securityinsight.AutomationRulePropertyConditionSupportedOperatorStartsWith),
						}, false),
					},

					"values": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			ConflictsWith: []string{"condition_json"},
		}
		schema["condition_json"].Computed = true
		schema["condition_json"].ConflictsWith = []string{"condition"}
	}

	return &pluginsdk.Resource{
		Create: resourceSentinelAutomationRuleCreateOrUpdate,
		Read:   resourceSentinelAutomationRuleRead,
		Update: resourceSentinelAutomationRuleCreateOrUpdate,
		Delete: resourceSentinelAutomationRuleDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SentinelAutomationRuleV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomationRuleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(5 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: schema,
	}
}

func resourceSentinelAutomationRuleCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	tenantId := meta.(*clients.Client).Account.TenantId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, err := workspaces.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewAutomationRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_sentinel_automation_rule", id.ID())
		}
	}

	actions, err := expandAutomationRuleActions(d, tenantId)
	if err != nil {
		return err
	}
	params := securityinsight.AutomationRule{
		AutomationRuleProperties: &securityinsight.AutomationRuleProperties{
			DisplayName: utils.String(d.Get("display_name").(string)),
			Order:       utils.Int32(int32(d.Get("order").(int))),
			TriggeringLogic: &securityinsight.AutomationRuleTriggeringLogic{
				IsEnabled:    utils.Bool(d.Get("enabled").(bool)),
				TriggersOn:   securityinsight.TriggersOn(d.Get("triggers_on").(string)),
				TriggersWhen: securityinsight.TriggersWhen(d.Get("triggers_when").(string)),
				Conditions:   expandAutomationRuleConditions(d.Get("condition").([]interface{})),
			},
			Actions: actions,
		},
	}

	if v, ok := d.GetOk("condition_json"); ok {
		conditions, err := expandAutomationRuleConditionsFromJSON(v.(string))
		if err != nil {
			return fmt.Errorf("expanding `condition_json`: %v", err)
		}
		params.AutomationRuleProperties.TriggeringLogic.Conditions = conditions
	} else if !features.FourPointOhBeta() {
		params.AutomationRuleProperties.TriggeringLogic.Conditions = expandAutomationRuleConditions(d.Get("condition").([]interface{}))
	}

	if expiration := d.Get("expiration").(string); expiration != "" {
		t, _ := time.Parse(time.RFC3339, expiration)
		params.AutomationRuleProperties.TriggeringLogic.ExpirationTimeUtc = &date.Time{Time: t}
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.Name, &params)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAutomationRuleRead(d, meta)
}

func resourceSentinelAutomationRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	if prop := resp.AutomationRuleProperties; prop != nil {
		d.Set("display_name", prop.DisplayName)

		var order int
		if prop.Order != nil {
			order = int(*prop.Order)
		}
		d.Set("order", order)

		if tl := prop.TriggeringLogic; tl != nil {
			var enabled bool
			if tl.IsEnabled != nil {
				enabled = *tl.IsEnabled
			}
			d.Set("enabled", enabled)

			d.Set("triggers_on", string(tl.TriggersOn))
			d.Set("triggers_when", string(tl.TriggersWhen))

			var expiration string
			if tl.ExpirationTimeUtc != nil {
				expiration = tl.ExpirationTimeUtc.Format(time.RFC3339)
			}
			d.Set("expiration", expiration)

			if !features.FourPointOhBeta() {
				if err := d.Set("condition", flattenAutomationRuleConditions(tl.Conditions)); err != nil {
					return fmt.Errorf("setting `condition`: %v", err)
				}
			}

			conditionJSON, err := flattenAutomationRuleConditionsToJSON(tl.Conditions)
			if err != nil {
				return fmt.Errorf("flattening `condition_json`: %v", err)
			}
			d.Set("condition_json", conditionJSON)
		}

		actionIncident, actionPlaybook := flattenAutomationRuleActions(prop.Actions)

		if err := d.Set("action_incident", actionIncident); err != nil {
			return fmt.Errorf("setting `action_incident`: %v", err)
		}
		if err := d.Set("action_playbook", actionPlaybook); err != nil {
			return fmt.Errorf("setting `action_playbook`: %v", err)
		}
	}

	return nil
}

func resourceSentinelAutomationRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationRuleID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandAutomationRuleConditions(input []interface{}) *[]securityinsight.BasicAutomationRuleCondition {
	if len(input) == 0 {
		return nil
	}

	out := make([]securityinsight.BasicAutomationRuleCondition, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		out = append(out, &securityinsight.PropertyConditionProperties{
			ConditionProperties: &securityinsight.AutomationRulePropertyValuesCondition{
				PropertyName:   securityinsight.AutomationRulePropertyConditionSupportedProperty(b["property"].(string)),
				Operator:       securityinsight.AutomationRulePropertyConditionSupportedOperator(b["operator"].(string)),
				PropertyValues: utils.ExpandStringSlice(b["values"].([]interface{})),
			},
			ConditionType: securityinsight.ConditionTypeBasicAutomationRuleConditionConditionTypeProperty,
		})
	}
	return &out
}

func flattenAutomationRuleConditions(conditions *[]securityinsight.BasicAutomationRuleCondition) interface{} {
	if conditions == nil {
		return nil
	}

	out := make([]interface{}, 0, len(*conditions))
	for _, condition := range *conditions {
		// "condition" only applies to the Property condition
		condition, ok := condition.(securityinsight.PropertyConditionProperties)
		if !ok {
			continue
		}

		var (
			property string
			operator string
			values   []interface{}
		)
		if p := condition.ConditionProperties; p != nil {
			property = string(p.PropertyName)
			operator = string(p.Operator)
			values = utils.FlattenStringSlice(p.PropertyValues)
		}

		out = append(out, map[string]interface{}{
			"property": property,
			"operator": operator,
			"values":   values,
		})
	}
	return out
}

func expandAutomationRuleConditionsFromJSON(input string) (*[]securityinsight.BasicAutomationRuleCondition, error) {
	if input == "" {
		return nil, nil
	}
	triggerLogic := &securityinsight.AutomationRuleTriggeringLogic{}
	err := triggerLogic.UnmarshalJSON([]byte(fmt.Sprintf(`{ "conditions": %s }`, input)))
	if err != nil {
		return nil, err
	}
	return triggerLogic.Conditions, nil
}

func flattenAutomationRuleConditionsToJSON(input *[]securityinsight.BasicAutomationRuleCondition) (string, error) {
	if input == nil {
		return "", nil
	}
	result, err := json.Marshal(input)
	return string(result), err
}

func expandAutomationRuleActions(d *pluginsdk.ResourceData, defaultTenantId string) (*[]securityinsight.BasicAutomationRuleAction, error) {
	actionIncident, err := expandAutomationRuleActionIncident(d.Get("action_incident").([]interface{}))
	if err != nil {
		return nil, err
	}
	actionPlaybook, err := expandAutomationRuleActionPlaybook(d.Get("action_playbook").([]interface{}), defaultTenantId)
	if err != nil {
		return nil, err
	}

	if len(actionIncident)+len(actionPlaybook) == 0 {
		return nil, nil
	}

	out := make([]securityinsight.BasicAutomationRuleAction, 0, len(actionIncident)+len(actionPlaybook))
	out = append(out, actionIncident...)
	out = append(out, actionPlaybook...)
	return &out, nil
}

func flattenAutomationRuleActions(input *[]securityinsight.BasicAutomationRuleAction) (actionIncident []interface{}, actionPlaybook []interface{}) {
	if input == nil {
		return nil, nil
	}

	actionIncident = make([]interface{}, 0)
	actionPlaybook = make([]interface{}, 0)

	for _, action := range *input {
		switch action := action.(type) {
		case securityinsight.AutomationRuleModifyPropertiesAction:
			actionIncident = append(actionIncident, flattenAutomationRuleActionIncident(action))
		case securityinsight.AutomationRuleRunPlaybookAction:
			actionPlaybook = append(actionPlaybook, flattenAutomationRuleActionPlaybook(action))
		}
	}

	return
}

func expandAutomationRuleActionIncident(input []interface{}) ([]securityinsight.BasicAutomationRuleAction, error) {
	if len(input) == 0 {
		return nil, nil
	}

	out := make([]securityinsight.BasicAutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		status := securityinsight.IncidentStatus(b["status"].(string))
		l := strings.Split(b["classification"].(string), "_")
		classification, clr := l[0], ""
		if len(l) == 2 {
			clr = l[1]
		}
		classificationComment := b["classification_comment"].(string)

		// sanity check on classification
		if status == securityinsight.IncidentStatusClosed && classification == "" {
			return nil, fmt.Errorf("`classification` is required when `status` is set to `Closed`")
		}
		if status != securityinsight.IncidentStatusClosed {
			if classification != "" {
				return nil, fmt.Errorf("`classification` can't be set when `status` is not set to `Closed`")
			}
			if classificationComment != "" {
				return nil, fmt.Errorf("`classification_comment` can't be set when `status` is not set to `Closed`")
			}
		}

		var labelsPtr *[]securityinsight.IncidentLabel
		if labelStrsPtr := utils.ExpandStringSlice(b["labels"].([]interface{})); labelStrsPtr != nil && len(*labelStrsPtr) > 0 {
			labels := make([]securityinsight.IncidentLabel, 0, len(*labelStrsPtr))
			for _, label := range *labelStrsPtr {
				labels = append(labels, securityinsight.IncidentLabel{
					LabelName: utils.String(label),
				})
			}
			labelsPtr = &labels
		}

		var ownerPtr *securityinsight.IncidentOwnerInfo
		if ownerIdStr := b["owner_id"].(string); ownerIdStr != "" {
			ownerId, err := uuid.FromString(ownerIdStr)
			if err != nil {
				return nil, fmt.Errorf("getting `owner_id`: %v", err)
			}
			ownerPtr = &securityinsight.IncidentOwnerInfo{
				ObjectID: &ownerId,
			}
		}

		severity := b["severity"].(string)

		// sanity check on the whole incident action
		if severity == "" && ownerPtr == nil && labelsPtr == nil && status == "" {
			return nil, fmt.Errorf("at least one of `severity`, `owner_id`, `labels` or `status` should be specified")
		}

		out = append(out, securityinsight.AutomationRuleModifyPropertiesAction{
			ActionType: securityinsight.ActionTypeBasicAutomationRuleActionActionTypeModifyProperties,
			Order:      utils.Int32(int32(b["order"].(int))),
			ActionConfiguration: &securityinsight.IncidentPropertiesAction{
				Status:                status,
				Classification:        securityinsight.IncidentClassification(classification),
				ClassificationComment: &classificationComment,
				ClassificationReason:  securityinsight.IncidentClassificationReason(clr),
				Labels:                labelsPtr,
				Owner:                 ownerPtr,
				Severity:              securityinsight.IncidentSeverity(severity),
			},
		})
	}

	return out, nil
}

func flattenAutomationRuleActionIncident(input securityinsight.AutomationRuleModifyPropertiesAction) map[string]interface{} {
	var order int
	if input.Order != nil {
		order = int(*input.Order)
	}

	var (
		status      string
		clsf        string
		clsfComment string
		clsfReason  string
		labels      []interface{}
		owner       string
		severity    string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		status = string(cfg.Status)
		clsf = string(cfg.Classification)
		if cfg.ClassificationComment != nil {
			clsfComment = *cfg.ClassificationComment
		}
		clsfReason = string(cfg.ClassificationReason)

		if cfg.Labels != nil {
			for _, label := range *cfg.Labels {
				if label.LabelName != nil {
					labels = append(labels, *label.LabelName)
				}
			}
		}

		if cfg.Owner != nil && cfg.Owner.ObjectID != nil {
			owner = cfg.Owner.ObjectID.String()
		}

		severity = string(cfg.Severity)
	}

	classification := clsf
	if clsfReason != "" {
		classification = classification + "_" + clsfReason
	}

	return map[string]interface{}{
		"order":                  order,
		"status":                 status,
		"classification":         classification,
		"classification_comment": clsfComment,
		"labels":                 labels,
		"owner_id":               owner,
		"severity":               severity,
	}
}

func expandAutomationRuleActionPlaybook(input []interface{}, defaultTenantId string) ([]securityinsight.BasicAutomationRuleAction, error) {
	if len(input) == 0 {
		return nil, nil
	}

	out := make([]securityinsight.BasicAutomationRuleAction, 0, len(input))
	for _, b := range input {
		b := b.(map[string]interface{})

		tid := defaultTenantId
		if t := b["tenant_id"].(string); t != "" {
			tid = t
		}

		tenantId, err := uuid.FromString(tid)
		if err != nil {
			return nil, fmt.Errorf("getting `tenant_id`: %v", err)
		}

		out = append(out, securityinsight.AutomationRuleRunPlaybookAction{
			ActionType: securityinsight.ActionTypeBasicAutomationRuleActionActionTypeRunPlaybook,
			Order:      utils.Int32(int32(b["order"].(int))),
			ActionConfiguration: &securityinsight.PlaybookActionProperties{
				LogicAppResourceID: utils.String(b["logic_app_id"].(string)),
				TenantID:           &tenantId,
			},
		})
	}
	return out, nil
}

func flattenAutomationRuleActionPlaybook(input securityinsight.AutomationRuleRunPlaybookAction) map[string]interface{} {
	var order int

	if input.Order != nil {
		order = int(*input.Order)
	}

	var (
		logicAppId string
		tenantId   string
	)

	if cfg := input.ActionConfiguration; cfg != nil {
		if cfg.LogicAppResourceID != nil {
			logicAppId = *cfg.LogicAppResourceID
		}

		if cfg.TenantID != nil {
			tenantId = cfg.TenantID.String()
		}
	}

	return map[string]interface{}{
		"order":        order,
		"logic_app_id": logicAppId,
		"tenant_id":    tenantId,
	}
}
