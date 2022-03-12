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

	// @TODO: The way we're using this (a precreated boot image in a Power VS Service instance) doesn't
	// align with the installer's definition of this. We need a new var here and in the install config.
	// This should be done before code cutoff in a followon PR.
	if osOverride := os.Getenv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); len(osOverride) != 0 {
		p.ClusterOSImage = osOverride
	}

	p.Region = ssn.Session.Options.Region
	p.Zone = ssn.Session.Options.Zone
	p.UserID = ssn.Session.Options.UserAccount

	return &p, nil
}
