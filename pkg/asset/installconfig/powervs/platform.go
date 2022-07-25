package powervs

import (
	"os"

	"github.com/openshift/installer/pkg/types/powervs"
)

// Platform collects powervs-specific configuration.
func Platform() (*powervs.Platform, error) {

	bxCli, err := NewBxClient()
	if err != nil {
		return nil, err
	}

	err = bxCli.NewPISession()
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

	p.Region = bxCli.PISession.Options.Region
	p.Zone = bxCli.PISession.Options.Zone
	p.UserID = bxCli.PISession.Options.UserAccount

	return &p, nil
}
