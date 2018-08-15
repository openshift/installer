package installerassets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openshift/installer/pkg/assets"
)

func metadataRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "metadata.json",
		RebuildHelper: metadataRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "cluster-name", "platform")
	if err != nil {
		return nil, err
	}

	clusterName := string(parents["cluster-name"].Data)
	platform := string(parents["platform"].Data)

	platformMetadataName := fmt.Sprintf("%s/metadata.json", platform)
	parents, err = asset.GetParents(ctx, getByName, platformMetadataName)
	if err != nil {
		return nil, err
	}

	platformMetadata := json.RawMessage(parents[platformMetadataName].Data)
	metadata := map[string]interface{}{
		"clusterName": clusterName,
		platform:      platformMetadata,
	}

	asset.Data, err = json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["metadata.json"] = metadataRebuilder
}
