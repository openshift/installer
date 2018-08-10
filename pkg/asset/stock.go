package asset

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
}

func EstablishStock() *Stock {
	s := &Stock{}

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
