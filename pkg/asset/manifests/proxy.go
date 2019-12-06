package manifests

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/vsphere"
)

var proxyCfgFilename = filepath.Join(manifestDir, "cluster-proxy-01-config.yaml")

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
		&Networking{},
	}
}

// Generate generates the Proxy config and its CRD.
func (p *Proxy) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	network := &Networking{}
	dependencies.Get(installConfig, network)

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

		if installConfig.Config.AdditionalTrustBundle != "" {
			p.Config.Spec.TrustedCA = configv1.ConfigMapNameReference{
				Name: additionalTrustBundleConfigMapName,
			}
		}
	}

	if p.Config.Spec.HTTPProxy != "" || p.Config.Spec.HTTPSProxy != "" {
		noProxy, err := createNoProxy(installConfig, network)
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
// localhost, 127.0.0.1, api.clusterdomain, api-int.clusterdomain, etcd-idx.clusterdomain
// If platform is not vSphere or None add 169.254.169.254 to the list of NO_PROXY addresses.
// If platform is AWS, add ".ec2.internal" for region us-east-1 or for all other regions add
// ".<aws_region>.compute.internal" to the list of NO_PROXY addresses. We should not proxy
// the instance metadata services:
// https://docs.openstack.org/nova/latest/user/metadata.html
// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service
// https://cloud.google.com/compute/docs/storing-retrieving-metadata
func createNoProxy(installConfig *installconfig.InstallConfig, network *Networking) (string, error) {
	internalAPIServer, err := url.Parse(getInternalAPIServerURL(installConfig.Config))
	if err != nil {
		return "", errors.New("failed parsing internal API server when creating Proxy manifest")
	}

	set := sets.NewString(
		"127.0.0.1",
		"localhost",
		".svc",
		".cluster.local",
		network.Config.Spec.ServiceNetwork[0],
		internalAPIServer.Hostname(),
		installConfig.Config.Networking.MachineCIDR.String(),
	)
	platform := installConfig.Config.Platform.Name()

	if platform != vsphere.Name && platform != none.Name {
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

	// From https://cloud.google.com/vpc/docs/special-configurations add GCP metadata.
	// "metadata.google.internal." added due to https://bugzilla.redhat.com/show_bug.cgi?id=1754049
	if platform == gcp.Name {
		set.Insert("metadata", "metadata.google.internal", "metadata.google.internal.")
	}

	for i := int64(0); i < *installConfig.Config.ControlPlane.Replicas; i++ {
		etcdHost := fmt.Sprintf("etcd-%d.%s", i, installConfig.Config.ClusterDomain())
		set.Insert(etcdHost)
	}

	for _, clusterNetwork := range network.Config.Spec.ClusterNetwork {
		set.Insert(clusterNetwork.CIDR)
	}

	for _, userValue := range strings.Split(installConfig.Config.Proxy.NoProxy, ",") {
		if userValue != "" {
			set.Insert(userValue)
		}
	}

	return strings.Join(set.List(), ","), nil
}

// Files returns the files generated by the asset.
func (p *Proxy) Files() []*asset.File {
	return p.FileList
}

// Load loads the already-rendered files back from disk.
func (p *Proxy) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
