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
		&bootkube.EtcdHostServiceEndpoints{},
		&bootkube.EtcdServingCAConfigMap{},
		&bootkube.KubeSystemConfigmapRootCA{},
		&bootkube.EtcdClientSecret{},
		&bootkube.OpenshiftMachineConfigOperator{},
		&bootkube.EtcdNamespace{},
		&bootkube.EtcdService{},
		&bootkube.EtcdHostService{},
		&bootkube.EtcdMetricClientSecret{},
		&bootkube.EtcdMetricSignerSecret{},
		&bootkube.EtcdMetricServingCAConfigMap{},
		&bootkube.OpenshiftConfigSecretPullSecret{},
		&openshift.CloudCredsSecret{},
		&openshift.KubeadminPasswordSecret{},
		&openshift.RoleCloudCredsSecretReader{},
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

	// Cluster are the cluster targeted assets.
	Cluster = []asset.WritableAsset{
		&cluster.Metadata{},
		&cluster.TerraformVariables{},
		&kubeconfig.AdminClient{},
		&password.KubeadminPassword{},
		&tls.JournalCertKey{},
		&cluster.Cluster{},
	}
)
