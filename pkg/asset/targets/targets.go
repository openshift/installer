package targets

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
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
		&manifests.Manifests{},
		&manifests.Openshift{},
	}

	// ManifestTemplates are the manifest-templates targeted assets.
	ManifestTemplates = []asset.WritableAsset{
		&bootkube.KubeCloudConfig{},
		&bootkube.MachineConfigServerTLSSecret{},
		&bootkube.Pull{},
		&bootkube.CVOOverrides{},
		&bootkube.HostEtcdServiceEndpointsKubeSystem{},
		&bootkube.KubeSystemConfigmapEtcdServingCA{},
		&bootkube.KubeSystemConfigmapRootCA{},
		&bootkube.KubeSystemSecretEtcdClient{},
		&bootkube.OpenshiftMachineConfigOperator{},
		&bootkube.EtcdServiceKubeSystem{},
		&bootkube.HostEtcdServiceKubeSystem{},
		&openshift.BindingDiscovery{},
		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
	}

	// IgnitionConfigs are the ignition-configs targeted assets.
	IgnitionConfigs = []asset.WritableAsset{
		&kubeconfig.Admin{},
		&machine.Master{},
		&machine.Worker{},
		&bootstrap.Bootstrap{},
		&cluster.Metadata{},
	}

	// Cluster are the cluster targeted assets.
	Cluster = []asset.WritableAsset{
		&cluster.TerraformVariables{},
		&kubeconfig.Admin{},
		&tls.JournalCertKey{},
		&cluster.Metadata{},
		&cluster.Cluster{},
	}
)
