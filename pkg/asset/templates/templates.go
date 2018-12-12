// Package templates deals with creating template assets that will be used by other assets
package templates

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	"github.com/openshift/installer/pkg/asset/templates/content/openshift"
)

// Templates are the targeted assets for generating the dependent unrendered
// template files.
var Templates = []asset.WritableAsset{
	&bootkube.KubeCloudConfig{},
	&bootkube.MachineConfigServerTLSSecret{},
	&bootkube.OpenshiftServiceCertSignerSecret{},
	&bootkube.Pull{},
	&bootkube.CVOOverrides{},
	&bootkube.HostEtcdServiceEndpointsKubeSystem{},
	&bootkube.KubeSystemConfigmapEtcdServingCA{},
	&bootkube.KubeSystemConfigmapRootCA{},
	&bootkube.KubeSystemSecretEtcdClient{},
	&bootkube.OpenshiftMachineConfigOperator{},
	&bootkube.OpenshiftClusterAPINamespace{},
	&bootkube.OpenshiftServiceCertSignerNamespace{},
	&bootkube.EtcdServiceKubeSystem{},
	&bootkube.HostEtcdServiceKubeSystem{},
	&openshift.BindingDiscovery{},
	&openshift.CloudCredsSecret{},
	&openshift.KubeadminPasswordSecret{},
	&openshift.RoleCloudCredsSecretReader{},
}
