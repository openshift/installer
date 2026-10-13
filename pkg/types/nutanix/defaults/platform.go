package defaults

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *nutanix.Platform) {
	if p.DNSRecordsType == "" {
		p.DNSRecordsType = configv1.DNSRecordsTypeInternal
	}
}
