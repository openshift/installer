package targets

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
	"github.com/openshift/installer/pkg/asset/tls"
)

var (
	// InstallConfig are the install-config targeted assets.
	InstallConfig = []asset.WritableAsset{
		&installconfig.InstallConfig{},
	}

	// Manifests are the manifests targeted assets.
	Manifests = []asset.WritableAsset{
		&machines.Master{},
		&machines.Worker{},
		&manifests.Manifests{},
		&manifests.Openshift{},
	}

	// ManifestTemplates are the manifest-templates targeted assets.
	ManifestTemplates = []asset.WritableAsset{
		&bootkube.KubeCloudConfig{},
		&bootkube.MachineConfigServerTLSSecret{},
		&bootkube.CVOOverrides{},
		&bootkube.KubeSystemConfigmapRootCA{},
		&bootkube.OpenshiftMachineConfigOperator{},
		&bootkube.OpenshiftConfigSecretPullSecret{},
		&bootkube.KubevirtInfraNamespace{},
		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
		&openshift.AzureCloudProviderSecret{},
	}

	// IgnitionConfigs are the ignition-configs targeted assets.
	IgnitionConfigs = []asset.WritableAsset{
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&machine.Master{},
		&machine.Worker{},
		&bootstrap.Bootstrap{},
		&cluster.Metadata{},
	}

	// SingleNodeIgnitionConfig is the bootstrap-in-place ignition-config targeted assets.
	SingleNodeIgnitionConfig = []asset.WritableAsset{
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&bootstrap.SingleNodeBootstrapInPlace{},
		&cluster.Metadata{},
	}

	// Cluster are the cluster targeted assets.
	Cluster = []asset.WritableAsset{
		&cluster.Metadata{},
		&machine.MasterIgnitionCustomizations{},
		&machine.WorkerIgnitionCustomizations{},
		&cluster.TerraformVariables{},
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&tls.JournalCertKey{},
		&cluster.Cluster{},
	}
)
