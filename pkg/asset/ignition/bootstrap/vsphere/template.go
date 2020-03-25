package vsphere

import (
	"github.com/openshift/installer/pkg/types/vsphere"
)

// TemplateData holds data specific to templates used for the vsphere platform.
type TemplateData struct {
	// UserProvidedIPs specifies whether users provided IP addresses in the install config.
	UserProvidedVIPs bool
}

// GetTemplateData returns platform-specific data for bootstrap templates.
func GetTemplateData(config *vsphere.Platform) *TemplateData {
	var templateData TemplateData

	templateData.UserProvidedVIPs = config.APIVIP != ""

	return &templateData
}
