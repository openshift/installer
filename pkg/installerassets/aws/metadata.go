package aws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/openshift/installer/pkg/types/aws"
)

func metadataRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "aws/metadata.json",
		RebuildHelper: metadataRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"aws/region",
		"cluster-id",
		"cluster-name",
	)
	if err != nil {
		return nil, err
	}

	region := string(parents["aws/region"].Data)
	clusterID := string(parents["cluster-id"].Data)
	clusterName := string(parents["cluster-name"].Data)

	metadata := &aws.Metadata{
		Region: region,
		Identifier: []map[string]string{
			{
				"tectonicClusterID":                                  clusterID,
				fmt.Sprintf("kubernetes.io/cluster/%s", clusterName): "owned",
			},
		},
	}

	asset.Data, err = json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["aws/metadata.json"] = metadataRebuilder
}
