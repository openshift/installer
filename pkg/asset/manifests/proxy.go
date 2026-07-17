package manifests

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervc"
)

// BuildNoProxySet creates a set of NoProxy entries from install configuration.
// It includes localhost entries (.svc, .cluster.local, 127.0.0.1, localhost),
// all network CIDRs (cluster, service, machine), the internal API server hostname,
// and user-provided NoProxy values.
// The caller can add platform-specific entries as needed.
// When NoProxy is "*", returns nil and true to signal that all traffic should
// bypass the proxy without computing individual entries.
func BuildNoProxySet(config *types.InstallConfig) (sets.Set[string], bool) {
	if config.Proxy != nil && config.Proxy.NoProxy == "*" {
		return nil, true
	}

	set := sets.New[string](
		"127.0.0.1",
		"localhost",
		".svc",
		".cluster.local",
	)

	for _, network := range config.Networking.ServiceNetwork {
		set.Insert(network.String())
	}

	for _, network := range config.Networking.MachineNetwork {
		set.Insert(network.CIDR.String())
	}

	for _, network := range config.Networking.ClusterNetwork {
		set.Insert(network.CIDR.String())
	}

	// Add internal API server hostname
	set.Insert("api-int." + config.ClusterDomain())

	if config.Proxy != nil {
		for _, userValue := range strings.Split(config.Proxy.NoProxy, ",") {
			trimmed := strings.TrimSpace(userValue)
			if trimmed != "" {
				set.Insert(trimmed)
			}
		}
	}

	return set, false
}

var proxyCfgFilename = path.Join(manifestDir, "cluster-proxy-01-config.yaml")

// Proxy generates the cluster-proxy-*.yml files.
type Proxy struct {
	FileList []*asset.File
	Config   *configv1.Proxy
}

var _ asset.WritableAsset = (*Proxy)(nil)

// Name returns a human-friendly name for the asset.
func (*Proxy) Name() string {
	return "Proxy Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*Proxy) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the Proxy config and its CRD.
func (p *Proxy) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	p.Config = &configv1.Proxy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "Proxy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
	}

	if installConfig.Config.Proxy != nil {
		p.Config.Spec = configv1.ProxySpec{
			HTTPProxy:  installConfig.Config.Proxy.HTTPProxy,
			HTTPSProxy: installConfig.Config.Proxy.HTTPSProxy,
			NoProxy:    installConfig.Config.Proxy.NoProxy,
		}
	}

	if installConfig.Config.AdditionalTrustBundlePolicy == types.PolicyAlways ||
		installConfig.Config.Proxy != nil {
		if installConfig.Config.AdditionalTrustBundle != "" {
			p.Config.Spec.TrustedCA = configv1.ConfigMapNameReference{
				Name: additionalTrustBundleConfigMapName,
			}
		}
	}

	if p.Config.Spec.HTTPProxy != "" || p.Config.Spec.HTTPSProxy != "" {
		noProxy, err := createNoProxy(installConfig)
		if err != nil {
			return err
		}
		p.Config.Status = configv1.ProxyStatus{
			HTTPProxy:  installConfig.Config.Proxy.HTTPProxy,
			HTTPSProxy: installConfig.Config.Proxy.HTTPSProxy,
			NoProxy:    noProxy,
		}
	}

	configData, err := yaml.Marshal(p.Config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", p.Name())
	}

	p.FileList = []*asset.File{
		{
			Filename: proxyCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

// createNoProxy combines user-provided & platform-specific values to create a comma-separated
// list of unique NO_PROXY values. Platform values are: serviceCIDR, podCIDR, machineCIDR,
// localhost, 127.0.0.1, api.clusterdomain, api-int.clusterdomain.
// If platform is AWS, GCP, Azure, or OpenStack add 169.254.169.254 to the list of NO_PROXY addresses.
// If platform is AWS, add ".ec2.internal" for region us-east-1 or for all other regions add
// ".<aws_region>.compute.internal" to the list of NO_PROXY addresses. We should not proxy
// the instance metadata services:
// https://docs.openstack.org/nova/latest/user/metadata.html
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service
// https://cloud.google.com/compute/docs/storing-retrieving-metadata
func createNoProxy(installConfig *installconfig.InstallConfig) (string, error) {
	set, wildcard := BuildNoProxySet(installConfig.Config)
	if wildcard {
		return "*", nil
	}

	platform := installConfig.Config.Platform.Name()

	// FIXME: The cluster-network-operator duplicates this code in pkg/util/proxyconfig/no_proxy.go,
	//  if altering this list of platforms, you must ALSO alter the code in cluster-network-operator.
	switch platform {
	case aws.Name, gcp.Name, azure.Name, openstack.Name, powervc.Name:
		set.Insert("169.254.169.254")
	}

	// TODO: Add support for additional cloud providers.
	if platform == aws.Name {
		region := installConfig.Config.AWS.Region
		if region == "us-east-1" {
			set.Insert(".ec2.internal")
		} else {
			set.Insert(fmt.Sprintf(".%s.compute.internal", region))
		}
	}

	// TODO: IBM[#95]: proxy

	if platform == azure.Name && installConfig.Azure.CloudName != azure.PublicCloud {
		// https://learn.microsoft.com/en-us/azure/virtual-network/what-is-ip-address-168-63-129-16
		set.Insert("168.63.129.16")
		if installConfig.Azure.CloudName == azure.StackCloud {
			set.Insert(installConfig.Config.Azure.ARMEndpoint)
		}
	}

	// From https://cloud.google.com/vpc/docs/special-configurations add GCP metadata.
	// "metadata.google.internal." added due to https://bugzilla.redhat.com/show_bug.cgi?id=1754049
	if platform == gcp.Name {
		set.Insert("metadata", "metadata.google.internal", "metadata.google.internal.")
	}

	return strings.Join(sets.List(set), ","), nil
}

// Files returns the files generated by the asset.
func (p *Proxy) Files() []*asset.File {
	return p.FileList
}

// Load loads the already-rendered files back from disk.
func (p *Proxy) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
