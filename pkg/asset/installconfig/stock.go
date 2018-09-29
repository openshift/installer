package installconfig

import (
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/types/config"
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

// EstablishStock establishes the stock of assets.
func (s *StockImpl) EstablishStock() {
	s.installConfig = &installConfig{
		assetStock: s,
	}
	s.clusterID = &clusterID{}
	s.emailAddress = &asset.UserProvided{
		AssetName: "Email Address",
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Email Address",
				Help:    "The email address of the cluster administrator. This will be used to log in to the console.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return config.ValidateEmail(ans.(string))
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_EMAIL_ADDRESS",
	}
	s.password = &asset.UserProvided{
		AssetName: "Password",
		Question: &survey.Question{
			Prompt: &survey.Password{
				Message: "Password",
				Help:    "The password of the cluster administrator. This will be used to log in to the console.",
			},
		},
		EnvVarName: "OPENSHIFT_INSTALL_PASSWORD",
	}
	s.baseDomain = &asset.UserProvided{
		AssetName: "Base Domain",
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return config.ValidateDomainName(ans.(string))
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_BASE_DOMAIN",
	}
	s.clusterName = &asset.UserProvided{
		AssetName: "Cluster Name",
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Cluster Name",
				Help:    "The name of the cluster. This will be used when generating sub-domains.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return config.ValidateDomainName(ans.(string))
			}),
		},
		EnvVarName: "OPENSHIFT_INSTALL_CLUSTER_NAME",
	}
	s.pullSecret = &asset.UserProvided{
		AssetName: "Pull Secret",
		Question: &survey.Question{
			Prompt: &survey.Input{
				Message: "Pull Secret",
				Help:    "The container registry pull secret for this cluster.",
			},
			Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
				return config.ValidateJSON([]byte(ans.(string)))
			}),
		},
		EnvVarName:     "OPENSHIFT_INSTALL_PULL_SECRET",
		PathEnvVarName: "OPENSHIFT_INSTALL_PULL_SECRET_PATH",
	}
	s.platform = &Platform{}
	s.sshKey = &sshPublicKey{}
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
