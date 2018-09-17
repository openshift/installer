package installconfig

import (
	"bufio"

	"github.com/openshift/installer/pkg/asset"
)

// Stock is the stock of InstallConfig assets that can be generated.
type Stock interface {
	// InstallConfig is the asset that generates install-config.yml.
	InstallConfig() asset.Asset
	// ClusterID is the asset that generates a UUID for the cluster.
	ClusterID() asset.Asset
	// EmailAddress is the asset that queries the user for the admin email address.
	EmailAddress() asset.Asset
	// Password is the asset that queries the user for the admin password.
	Password() asset.Asset
	// SSHKey is the asset that queries the user for the ssh public key in a string format.
	SSHKey() asset.Asset
	// BaseDomain is the asset that queries the user for the base domain to use
	// for the cluster.
	BaseDomain() asset.Asset
	// ClusterName is the asset that queries the user for the name of the cluster.
	ClusterName() asset.Asset
	// PullSecret is the asset that queries the user for the pull secret.
	PullSecret() asset.Asset
	// Platform is the asset that queries the user for the platform on which
	// to create the cluster.
	Platform() asset.Asset
}

// StockImpl implements the Stock interface for installconfig and user inputs.
type StockImpl struct {
	installConfig asset.Asset
	clusterID     asset.Asset
	emailAddress  asset.Asset
	password      asset.Asset
	sshKey        asset.Asset
	baseDomain    asset.Asset
	clusterName   asset.Asset
	pullSecret    asset.Asset
	platform      asset.Asset
}

// EstablishStock establishes the stock of assets in the specified directory.
func (s *StockImpl) EstablishStock(directory string, inputReader *bufio.Reader) {
	s.installConfig = &installConfig{
		assetStock: s,
		directory:  directory,
	}
	s.clusterID = &clusterID{}
	s.emailAddress = &asset.UserProvided{
		Prompt:      "Email Address:",
		InputReader: inputReader,
	}
	s.password = &password{
		InputReader: inputReader,
	}
	s.sshKey = &sshPublicKey{
		inputReader: inputReader,
	}
	s.baseDomain = &asset.UserProvided{
		Prompt:      "Base Domain:",
		InputReader: inputReader,
	}
	s.clusterName = &asset.UserProvided{
		Prompt:      "Cluster Name:",
		InputReader: inputReader,
	}
	s.pullSecret = &asset.UserProvided{
		Prompt:      "Pull Secret:",
		InputReader: inputReader,
	}
	s.platform = &Platform{InputReader: inputReader}
}

// ClusterID is the asset that generates a UUID for the cluster.
func (s *StockImpl) ClusterID() asset.Asset {
	return s.clusterID
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

// SSHKey is the asset that queries the user for the ssh public key in a string format.
func (s *StockImpl) SSHKey() asset.Asset {
	return s.sshKey
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

// PullSecret is the asset that queries the user for the pull secret.
func (s *StockImpl) PullSecret() asset.Asset {
	return s.pullSecret
}

// Platform is the asset that queries the user for the platform on which
// to create the cluster.
func (s *StockImpl) Platform() asset.Asset {
	return s.platform
}
