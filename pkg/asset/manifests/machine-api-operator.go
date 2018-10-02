package manifests

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	maoTargetNamespace = "openshift-cluster-api"
	// DefaultChannel is the default RHCOS channel for the cluster.
	DefaultChannel = "tested"
)

// machineAPIOperator generates the network-operator-*.yml files
type machineAPIOperator struct {
	installConfigAsset asset.Asset
	installConfig      *types.InstallConfig
	aggregatorCA       asset.Asset
}

var _ asset.Asset = (*machineAPIOperator)(nil)

// maoOperatorConfig contains configuration for mao managed stack
// TODO(enxebre): move up to github.com/coreos/tectonic-config (to install-config? /rchopra)
type maoOperatorConfig struct {
	metav1.TypeMeta `json:",inline"`
	TargetNamespace string           `json:"targetNamespace"`
	APIServiceCA    string           `json:"apiServiceCA"`
	Provider        string           `json:"provider"`
	AWS             *awsConfig       `json:"aws"`
	Libvirt         *libvirtConfig   `json:"libvirt"`
	OpenStack       *openstackConfig `json:"openstack"`
}

type libvirtConfig struct {
	ClusterName string `json:"clusterName"`
	URI         string `json:"uri"`
	NetworkName string `json:"networkName"`
	IPRange     string `json:"iprange"`
	Replicas    int    `json:"replicas"`
}

type awsConfig struct {
	ClusterName      string `json:"clusterName"`
	ClusterID        string `json:"clusterID"`
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone"`
	Image            string `json:"image"`
	Replicas         int    `json:"replicas"`
}

type openstackConfig struct {
	ClusterName string `json:"clusterName"`
	ClusterID   string `json:"clusterID"`
	Region      string `json:"region"`
	Replicas    int    `json:"replicas"`
}

// Name returns a human friendly name for the operator
func (mao *machineAPIOperator) Name() string {
	return "Machine API Operator"
}

// Dependencies returns all of the dependencies directly needed by an
// machineAPIOperator asset.
func (mao *machineAPIOperator) Dependencies() []asset.Asset {
	return []asset.Asset{
		mao.installConfigAsset,
		mao.aggregatorCA,
	}
}

// Generate generates the network-operator-config.yml and network-operator-manifest.yml files
func (mao *machineAPIOperator) Generate(dependencies map[asset.Asset]*asset.State) (*asset.State, error) {
	ic, err := installconfig.GetInstallConfig(mao.installConfigAsset, dependencies)
	if err != nil {
		return nil, err
	}
	mao.installConfig = ic

	// installconfig is ready, we can create the mao config from it now
	maoConfig, err := mao.maoConfig(dependencies)
	if err != nil {
		return nil, err
	}

	state := &asset.State{
		Contents: []asset.Content{
			{
				Name: "machine-api-operator-config.yml",
				Data: []byte(maoConfig),
			},
		},
	}
	return state, nil
}

func (mao *machineAPIOperator) maoConfig(dependencies map[asset.Asset]*asset.State) (string, error) {
	cfg := maoOperatorConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "machineAPIOperatorConfig",
		},

		TargetNamespace: maoTargetNamespace,
	}

	ca := dependencies[mao.aggregatorCA].Contents[certIndex].Data
	cfg.APIServiceCA = string(ca)
	cfg.Provider = tectonicCloudProvider(mao.installConfig.Platform)

	if mao.installConfig.Platform.AWS != nil {
		var ami string

		ami, err := rhcos.AMI(context.TODO(), DefaultChannel, mao.installConfig.Platform.AWS.Region)
		if err != nil {
			return "", fmt.Errorf("failed to lookup RHCOS AMI: %v", err)
		}

		cfg.AWS = &awsConfig{
			ClusterName:      mao.installConfig.Name,
			ClusterID:        mao.installConfig.ClusterID,
			Region:           mao.installConfig.Platform.AWS.Region,
			AvailabilityZone: "",
			Image:            ami,
			Replicas:         int(*mao.installConfig.Machines[1].Replicas),
		}
	} else if mao.installConfig.Platform.Libvirt != nil {
		cfg.Libvirt = &libvirtConfig{
			ClusterName: mao.installConfig.Name,
			URI:         mao.installConfig.Platform.Libvirt.URI,
			NetworkName: mao.installConfig.Platform.Libvirt.Network.Name,
			IPRange:     mao.installConfig.Platform.Libvirt.Network.IPRange,
			Replicas:    int(*mao.installConfig.Machines[1].Replicas),
		}
	} else if mao.installConfig.Platform.OpenStack != nil {
		cfg.OpenStack = &openstackConfig{
			ClusterName: mao.installConfig.Name,
			ClusterID:   mao.installConfig.ClusterID,
			Region:      mao.installConfig.Platform.OpenStack.Region,
			Replicas:    int(*mao.installConfig.Machines[1].Replicas),
		}
	} else {
		return "", fmt.Errorf("unknown provider for machine-api-operator")
	}

	return marshalYAML(cfg)
}
