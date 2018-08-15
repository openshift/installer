package installerassets

import (
	"context"
	"fmt"
	"os"

	"github.com/openshift/installer/pkg/assets"
)

func clusterRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "cluster",
		RebuildHelper: clusterRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"platform",
		// unused in this rebuilder, but we want these before launching the cluster
		"metadata.json",
		"terraform/terraform.tfvars",
	)
	if err != nil {
		return nil, err
	}

	platform := string(parents["platform"].Data)
	perPlatformName := fmt.Sprintf("terraform/%s-terraform.auto.tfvars", platform)
	parents, err = asset.GetParents(ctx, getByName, perPlatformName)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["cluster"] = clusterRebuilder
}
