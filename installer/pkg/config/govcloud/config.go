package govcloud

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config/aws"
)

// Config defines the GovCloud configuraiton for a cluster.
type Config struct {
	AWS         aws.Config `yaml:",inline"`
	DNSServerIP string     `yaml:"DNSServerIP,omitempty"`
}
