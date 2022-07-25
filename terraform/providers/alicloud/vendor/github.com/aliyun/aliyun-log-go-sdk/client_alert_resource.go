package sls

import "encoding/json"

const (
	ResourceNameAlertPolicy        = "sls.alert.alert_policy"
	ResourceNameActionPolicy       = "sls.alert.action_policy"
	ResourceNameUser               = "sls.common.user"
	ResourceNameUserGroup          = "sls.common.user_group"
	ResourceNameContentTemplate    = "sls.alert.content_template"
	ResourceNameGlobalConfig       = "sls.alert.global_config"
	ResourceNameWebhookIntegration = "sls.alert.webhook_application"
)

type (
	// Notified users.
	ResourceUser struct {
		UserId       string   `json:"user_id"`
		UserName     string   `json:"user_name"`
		Enabled      bool     `json:"enabled"`
		CountryCode  string   `json:"country_code"`
		Phone        string   `json:"phone"`
		Email        []string `json:"email"`
		SmsEnabled   bool     `json:"sms_enabled"`
		VoiceEnabled bool     `json:"voice_enabled"`
	}

	// ResourceUserGroup is a collection of users.
	ResourceUserGroup struct {
		Id      string   `json:"user_group_id"`
		Name    string   `json:"user_group_name"`
		Enabled bool     `json:"enabled"`
		Members []string `json:"members"`
	}

	// ResourceAlertPolicy defines how alerts should be grouped, inhibited and silenced.
	ResourceAlertPolicy struct {
		PolicyId      string `json:"policy_id"`
		PolicyName    string `json:"policy_name"`
		Parent        string `json:"parent_id"`
		IsDefault     bool   `json:"is_default"`
		GroupPolicy   string `json:"group_script"`
		InhibitPolicy string `json:"inhibit_script"`
		SilencePolicy string `json:"silence_script"`
	}

	// ResourceActionPolicy defines how to send alert notifications.
	ResourceActionPolicy struct {
		ActionPolicyId              string            `json:"action_policy_id"`
		ActionPolicyName            string            `json:"action_policy_name"`
		IsDefault                   bool              `json:"is_default"`
		PrimaryPolicyScript         string            `json:"primary_policy_script"`
		SecondaryPolicyScript       string            `json:"secondary_policy_script"`
		EscalationStartEnabled      bool              `json:"escalation_start_enabled"`
		EscalationStartTimeout      string            `json:"escalation_start_timeout"`
		EscalationInprogressEnabled bool              `json:"escalation_inprogress_enabled"`
		EscalationInprogressTimeout string            `json:"escalation_inprogress_timeout"`
		EscalationEnabled           bool              `json:"escalation_enabled"`
		EscalationTimeout           string            `json:"escalation_timeout"`
		Labels                      map[string]string `json:"labels"`
	}

	// ContentTemplate
	ResourceTemplate struct {
		Content  string `json:"content"`
		Locale   string `json:"locale"`
		Title    string `json:"title"`
		Subject  string `json:"subject"`
		SendType string `json:"send_type"`
		Limit    int    `json:"limit"`
	}

	ResourceTemplates struct {
		Sms           ResourceTemplate `json:"sms"`
		Voice         ResourceTemplate `json:"voice"`
		Email         ResourceTemplate `json:"email"`
		Dingtalk      ResourceTemplate `json:"dingtalk"`
		Webhook       ResourceTemplate `json:"webhook"`
		MessageCenter ResourceTemplate `json:"message_center"`
		Wechat        ResourceTemplate `json:"wechat"`
		Lark          ResourceTemplate `json:"lark"`
		Slack         ResourceTemplate `json:"slack"`
	}
	ResourceContentTemplate struct {
		TemplateId   string            `json:"template_id"`
		TemplateName string            `json:"template_name"`
		IsDefault    bool              `json:"is_default"`
		Templates    ResourceTemplates `json:"templates"`
	}

	// WebhookIntegration is a wrap of webhook notification config.
	ResourceWebhookHeader struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	WebhookIntegration struct {
		Id      string                   `json:"id"`
		Name    string                   `json:"name"`
		Method  string                   `json:"method"`
		Url     string                   `json:"url"`
		Type    string                   `json:"type"`
		Headers []*ResourceWebhookHeader `json:"headers"`
	}

	// GlobalConfig is the global configuration for alerts.
	GlobalConfig struct {
		ConfigId     string `json:"config_id"`
		ConfigName   string `json:"config_name"`
		ConfigDetail struct {
			AlertCenterLog struct {
				Region string `json:"region"`
			} `json:"alert_center_log"`
		} `json:"config_detail"`
	}
)

func JsonMarshal(v interface{}) string {
	vBytes, _ := json.Marshal(v)
	return string(vBytes)
}
