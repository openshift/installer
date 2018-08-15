package aws

import (
	"context"
	"time"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/pkg/errors"
)

func amiRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "aws/ami",
		RebuildHelper: amiRebuilder,
	}

	parents, err := asset.GetParents(ctx, getByName, "aws/region")
	if err != nil {
		return nil, err
	}

	amiContext, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	ami, err := rhcos.AMI(amiContext, rhcos.DefaultChannel, string(parents["aws/region"].Data))
	if err != nil {
		return nil, errors.Wrap(err, "failed to determine default AMI")
	}

	asset.Data = []byte(ami)
	return asset, nil
}

func init() {
	installerassets.Rebuilders["aws/ami"] = amiRebuilder
}
