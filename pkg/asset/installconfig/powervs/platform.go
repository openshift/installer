package powervs

import (
	"os"

	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {

	ssn, err := GetSession()
	if err != nil {
		return nil, err
	}

	var p powervs.Platform
	if osOverride := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); len(osOverride) != 0 {
		p.BootstrapOSImage = osOverride
		p.ClusterOSImage = osOverride
	}

	p.Region = ssn.Session.Region
	p.Zone = ssn.Session.Zone
	p.APIKey = ssn.Session.IAMToken
	p.UserID = ssn.Session.UserAccount

	return &p, nil
}
