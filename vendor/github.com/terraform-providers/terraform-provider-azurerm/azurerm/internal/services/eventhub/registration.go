package eventhub

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "EventHub"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
<<<<<<< HEAD
		"azurerm_eventhub_namespace": dataSourceEventHubNamespace(),
=======
		"azurerm_eventhub":                              dataSourceEventHub(),
		"azurerm_eventhub_authorization_rule":           dataSourceEventHubAuthorizationRule(),
		"azurerm_eventhub_consumer_group":               dataSourceEventHubConsumerGroup(),
		"azurerm_eventhub_namespace":                    dataSourceEventHubNamespace(),
		"azurerm_eventhub_namespace_authorization_rule": dataSourceEventHubNamespaceAuthorizationRule(),
>>>>>>> 5aa20dd53... vendor: bump terraform-provider-azure to version v2.17.0
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_eventhub_authorization_rule":                 resourceArmEventHubAuthorizationRule(),
		"azurerm_eventhub_cluster":                            resourceArmEventHubCluster(),
		"azurerm_eventhub_consumer_group":                     resourceArmEventHubConsumerGroup(),
		"azurerm_eventhub_namespace_authorization_rule":       resourceArmEventHubNamespaceAuthorizationRule(),
		"azurerm_eventhub_namespace_disaster_recovery_config": resourceArmEventHubNamespaceDisasterRecoveryConfig(),
		"azurerm_eventhub_namespace":                          resourceArmEventHubNamespace(),
		"azurerm_eventhub":                                    resourceArmEventHub()}
}
