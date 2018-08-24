package installconfig

import (
	"bufio"

	"github.com/openshift/installer/pkg/asset"
)

// Stock is the stock of InstallConfig assets that can be generated.
type Stock interface {
	// InstallConfig is the asset that generates install-config.yml.
	InstallConfig() asset.Asset
	// EmailAddress is the asset that queries the user for the admin email address.
	EmailAddress() asset.Asset
	// Password is the asset that queries the user for the admin password.
	Password() asset.Asset
	// BaseDomain is the asset that queries the user for the base domain to use
	// for the cluster.
	BaseDomain() asset.Asset
	// ClusterName is the asset that queries the user for the name of the cluster.
	ClusterName() asset.Asset
	// License is the asset that queries the user for the OpenShift license.
	License() asset.Asset
	// PullSecret is the asset that queries the user for the pull secret.
	PullSecret() asset.Asset
	// Platform is the asset that queries the user for the platform on which
	// to create the cluster.
	Platform() asset.Asset
}

// StockImpl is the
type StockImpl struct {
	installConfig asset.Asset
	emailAddress  asset.Asset
	password      asset.Asset
	baseDomain    asset.Asset
	clusterName   asset.Asset
	license       asset.Asset
	pullSecret    asset.Asset
	platform      asset.Asset
}

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(directory string, inputReader *bufio.Reader) {
	s.installConfig = &installConfig{
		assetStock:  s,
		directory:   directory,
		inputReader: inputReader,
	}
	s.emailAddress = &asset.UserProvided{
		Prompt:      "Email Address:",
		InputReader: inputReader,
	}
	s.password = &asset.UserProvided{
		Prompt:      "Password:",
		InputReader: inputReader,
	}
	s.baseDomain = &asset.UserProvided{
		Prompt:      "Base Domain:",
		InputReader: inputReader,
	}
	s.clusterName = &asset.UserProvided{
		Prompt:      "Cluster Name:",
		InputReader: inputReader,
	}
	s.license = &asset.UserProvided{
		Prompt:      "License:",
		InputReader: inputReader,
	}
	s.pullSecret = &asset.UserProvided{
		Prompt:      "Pull Secret:",
		InputReader: inputReader,
	}
	s.platform = &Platform{InputReader: inputReader}
}

// InstallConfig is the asset that generates install-config.yml.
func (s *StockImpl) InstallConfig() asset.Asset {
	return s.installConfig
}

// EmailAddress is the asset that queries the user for the admin email address.
func (s *StockImpl) EmailAddress() asset.Asset {
	return s.emailAddress
}

// Password is the asset that queries the user for the admin password.
func (s *StockImpl) Password() asset.Asset {
	return s.password
}

// BaseDomain is the asset that queries the user for the base domain to use
// for the cluster.
func (s *StockImpl) BaseDomain() asset.Asset {
	return s.baseDomain
}

// ClusterName is the asset that queries the user for the name of the cluster.
func (s *StockImpl) ClusterName() asset.Asset {
	return s.clusterName
}

// License is the asset that queries the user for the OpenShift license.
func (s *StockImpl) License() asset.Asset {
	return s.license
}

// PullSecret is the asset that queries the user for the pull secret.
func (s *StockImpl) PullSecret() asset.Asset {
	return s.pullSecret
}

// Platform is the asset that queries the user for the platform on which
// to create the cluster.
func (s *StockImpl) Platform() asset.Asset {
	return s.platform
}
