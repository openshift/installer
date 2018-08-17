package asset

import (
	"bufio"
	"os"
)

// Stock is the stock of assets that can be generated.
type Stock struct {
	// Targetable assets
	InstallConfig Asset

	// Non-targetable assets
	password     Asset
	baseDomain   Asset
	clusterName  Asset
	license      Asset
	pullSecret   Asset
	platform     Asset
	emailAddress Asset

	directory   string
	inputReader *bufio.Reader
}

// EstablishStock establishes the stock of assets in the specified directory.
func EstablishStock(directory string) *Stock {
	s := &Stock{
		directory:   directory,
		inputReader: bufio.NewReader(os.Stdin),
	}

	s.InstallConfig = &installConfig{assetStock: s}

	s.emailAddress = &userProvided{prompt: "Email Address:", inputReader: s.inputReader}
	s.password = &userProvided{prompt: "Password:", inputReader: s.inputReader}
	s.baseDomain = &userProvided{prompt: "Base Domain:", inputReader: s.inputReader}
	s.clusterName = &userProvided{prompt: "Cluster Name:", inputReader: s.inputReader}
	s.license = &userProvided{prompt: "License:", inputReader: s.inputReader}
	s.pullSecret = &userProvided{prompt: "Pull Secret:", inputReader: s.inputReader}
	s.platform = newPlatform(s.inputReader)

	return s
}

func (s *Stock) createAssetDirectory() error {
	return os.MkdirAll(s.directory, 0755)
}
