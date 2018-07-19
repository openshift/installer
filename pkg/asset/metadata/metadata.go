package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

const (
	// MetadataFilename is name of the file where clustermetadata is stored.
	MetadataFilename  = "metadata.json"
	metadataAssetName = "Cluster Metadata"
)

// Metadata depends on cluster and installconfig,
type Metadata struct {
	installConfig asset.Asset
	cluster       asset.Asset
}

var _ asset.Asset = (*Metadata)(nil)

// Name returns the human-friendly name of the asset.
func (m *Metadata) Name() string {
	return metadataAssetName
}

// Dependencies returns the dependency of the MetaData.
func (m *Metadata) Dependencies() []asset.Asset {
	return []asset.Asset{m.installConfig, m.cluster}
}

// Generate generates the metadata.yaml file.
func (m *Metadata) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	installCfg, err := installconfig.GetInstallConfig(m.installConfig, parents)
	if err != nil {
		return nil, err
	}

	cm := &types.ClusterMetadata{
		ClusterName: installCfg.Name,
	}
	switch {
	case installCfg.Platform.AWS != nil:
		cm.ClusterPlatformMetadata.AWS = &types.ClusterAWSPlatformMetadata{
			Region: installCfg.Platform.AWS.Region,
			Identifier: map[string]string{
				"tectonicClusterID": installCfg.ClusterID,
			},
		}
	case installCfg.Platform.OpenStack != nil:
		cm.ClusterPlatformMetadata.OpenStack = &types.ClusterOpenStackPlatformMetadata{
			Region: installCfg.Platform.OpenStack.Region,
			Identifier: map[string]string{
				"tectonicClusterID": installCfg.ClusterID,
			},
		}
	case installCfg.Platform.Libvirt != nil:
		cm.ClusterPlatformMetadata.Libvirt = &types.ClusterLibvirtPlatformMetadata{
			URI: installCfg.Platform.Libvirt.URI,
		}
	default:
		return nil, fmt.Errorf("no known platform")
	}

	data, err := json.Marshal(cm)
	if err != nil {
		return nil, err
	}
	return &asset.State{
		Contents: []asset.Content{
			{
				Name: MetadataFilename,
				Data: []byte(data),
			},
		},
	}, nil
}
