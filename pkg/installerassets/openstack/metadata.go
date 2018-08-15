package openstack

import (
	"context"
	"encoding/json"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/openshift/installer/pkg/types/openstack"
)

func metadataRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "openstack/metadata.json",
		RebuildHelper: metadataRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"openstack/cloud",
		"openstack/region",
		"cluster-id",
	)
	if err != nil {
		return nil, err
	}

	cloud := string(parents["openstack/cloud"].Data)
	region := string(parents["openstack/region"].Data)
	clusterID := string(parents["cluster-id"].Data)

	metadata := &openstack.Metadata{
		Cloud:  cloud,
		Region: region,
		Identifier: map[string]string{
			"tectonicClusterID": clusterID,
		},
	}

	asset.Data, err = json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["openstack/metadata.json"] = metadataRebuilder
}
