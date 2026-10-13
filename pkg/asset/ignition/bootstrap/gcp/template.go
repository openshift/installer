package gcp

import (
	"github.com/openshift/installer/pkg/types/dns"
	"github.com/openshift/installer/pkg/types/gcp"
)

// TemplateData holds data specific to templates used for the gcp platform.
type TemplateData struct {
	// UserProvisionedDNS indicates whether this feature has been enabled on GCP
	UserProvisionedDNS bool
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *gcp.Platform) *TemplateData {
	var templateData TemplateData

	templateData.UserProvisionedDNS = (config.UserProvisionedDNS == dns.UserProvisionedDNSEnabled)

	return &templateData
}
