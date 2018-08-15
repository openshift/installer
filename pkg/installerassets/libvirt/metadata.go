package libvirt

import (
	"context"
	"encoding/json"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/openshift/installer/pkg/types/libvirt"
)

func metadataRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "libvirt/metadata.json",
		RebuildHelper: metadataRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "libvirt/uri")
	if err != nil {
		return nil, err
	}

	uri := string(parents["libvirt/uri"].Data)

	metadata := &libvirt.Metadata{
		URI: uri,
	}

	asset.Data, err = json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["libvirt/metadata.json"] = metadataRebuilder
}
