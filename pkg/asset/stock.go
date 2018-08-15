package asset

import (
	"os"
)

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

	directory string
}

func EstablishStock(directory string) *Stock {
	s := &Stock{
		directory: directory,
	}

	s.InstallConfig = &InstallConfig{assetStock: s}

	s.password = &UserProvided{Prompt: "Password: "}
	s.baseDomain = &UserProvided{Prompt: "Base Domain: "}
	s.clusterName = &UserProvided{Prompt: "Cluster Name: "}
	s.license = &UserProvided{Prompt: "License: "}
	s.pullSecret = &UserProvided{Prompt: "Pull Secret: "}
	s.platform = &Platform{}
	s.emailAddress = &UserProvided{Prompt: "Email Address: "}

	return s
}

func (s *Stock) createAssetDirectory() error {
	return os.MkdirAll(s.directory, 0755)
}
