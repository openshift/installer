package installerassets

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/openshift/installer/pkg/assets"
	"github.com/pkg/errors"
)

type terraformConfig struct {
	ClusterID  string `json:"tectonic_cluster_id,omitempty"`
	Name       string `json:"tectonic_cluster_name,omitempty"`
	BaseDomain string `json:"tectonic_base_domain,omitempty"`
	Masters    int    `json:"tectonic_master_count,omitempty"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`
}

func terraformRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "terraform/terraform.tfvars",
		RebuildHelper: terraformRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"base-domain",
		"cluster-id",
		"cluster-name",
		"ignition/bootstrap.ign",
		"ignition/master.ign",
		"machines/master-count",
	)
	if err != nil {
		return nil, err
	}

	masterCount, err := strconv.ParseUint(string(parents["machines/master-count"].Data), 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "parse master count")
	}

	config := &terraformConfig{
		ClusterID:         string(parents["cluster-id"].Data),
		Name:              string(parents["cluster-name"].Data),
		BaseDomain:        string(parents["base-domain"].Data),
		Masters:           int(masterCount),
		IgnitionBootstrap: string(parents["ignition/bootstrap.ign"].Data),
		IgnitionMaster:    string(parents["ignition/master.ign"].Data),
	}

	asset.Data, err = json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	Rebuilders["terraform/terraform.tfvars"] = terraformRebuilder
}
