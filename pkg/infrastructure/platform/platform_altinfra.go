//go:build altinfra
// +build altinfra

package platform

import (
	"fmt"

	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/infrastructure/aws"

	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// ProviderForPlatform returns the stages to run to provision the infrastructure for the specified platform.
func ProviderForPlatform(platform, installDir string) ([]infrastructure.Stage, func() error, error) {
	var stages []infrastructure.Stage
	var cleanup func() error
	var err error

	switch platform {
	case awstypes.Name:
		if stages, cleanup, err = aws.InitializeProvider(installDir); err != nil {
			return nil, nil, err
		}
	default:
		panic(fmt.Sprintf("unsupported platform %q", platform))
	}
	return stages, cleanup, err
}
