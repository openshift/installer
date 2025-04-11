package aws

import (
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/dns"
)

// TemplateData holds data specific to templates used for the AWS platform.
type TemplateData struct {
	// UserProvisionedDNS indicates whether this feature has been enabled on AWS
	UserProvisionedDNS bool
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *aws.Platform) *TemplateData {
	var templateData TemplateData

	templateData.UserProvisionedDNS = (config.UserProvisionedDNS == dns.UserProvisionedDNSEnabled)

	return &templateData
}
