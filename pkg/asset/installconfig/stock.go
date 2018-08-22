package installconfig

import (
	"bufio"

	"github.com/openshift/installer/pkg/asset"
)

type Stock interface {
	InstallConfig() asset.Asset
	EmailAddress() asset.Asset
	Password() asset.Asset
	BaseDomain() asset.Asset
	ClusterName() asset.Asset
	License() asset.Asset
	PullSecret() asset.Asset
	Platform() asset.Asset
}

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
	s.platform = newPlatform(inputReader)
}

func (s *StockImpl) InstallConfig() asset.Asset {
	return s.installConfig
}

func (s *StockImpl) EmailAddress() asset.Asset {
	return s.emailAddress
}

func (s *StockImpl) Password() asset.Asset {
	return s.password
}

func (s *StockImpl) BaseDomain() asset.Asset {
	return s.baseDomain
}

func (s *StockImpl) ClusterName() asset.Asset {
	return s.clusterName
}

func (s *StockImpl) License() asset.Asset {
	return s.license
}

func (s *StockImpl) PullSecret() asset.Asset {
	return s.pullSecret
}

func (s *StockImpl) Platform() asset.Asset {
	return s.platform
}
