package azure

import (
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/dns"
)

// TemplateData holds data specific to templates used for the azure platform.
type TemplateData struct {
	// UserProvisionedDNS indicates whether this feature has been enabled on Azure
	UserProvisionedDNS bool
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *azure.Platform) *TemplateData {
	var templateData TemplateData

	templateData.UserProvisionedDNS = (config.UserProvisionedDNS == dns.UserProvisionedDNSEnabled)

	return &templateData
}
