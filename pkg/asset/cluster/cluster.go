package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types/config"
)

// Cluster uses the terraform executable to launch a cluster
// with the given terraform tfvar and generated templates.
type Cluster struct {
	// The root directory of the generated assets.
	rootDir    string
	tfvars     asset.Asset
	kubeconfig asset.Asset
}

var _ asset.Asset = (*Cluster)(nil)

// Name returns the human-friendly name of the asset.
func (c *Cluster) Name() string {
	return "Cluster"
}

// Dependencies returns the direct dependency for launching
// the cluster.
func (c *Cluster) Dependencies() []asset.Asset {
	return []asset.Asset{c.tfvars, c.kubeconfig}
}

// Generate launches the cluster and generates the terraform state file on disk.
func (c *Cluster) Generate(parents map[asset.Asset]*asset.State) (*asset.State, error) {
	state, ok := parents[c.tfvars]
	if !ok {
		return nil, fmt.Errorf("failed to get terraform.tfvar state in the parent asset states")
	}

	var tfvars config.Cluster
	if err := json.Unmarshal(state.Contents[0].Data, &tfvars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal terraform tfvars file: %v", err)
	}

	steps := []string{
		terraform.InfraStep,
		terraform.BootstrapStep,
	}

	dir, err := terraform.BaseLocation()
	if err != nil {
		return nil, err
	}

	var result asset.State

	for _, step := range steps {
		templateDir, err := terraform.FindStepTemplates(dir, step, tfvars.Platform)
		if err != nil {
			return nil, err
		}

		if err := terraform.Init(c.rootDir, templateDir); err != nil {
			return nil, err
		}

		stateFile, err := terraform.Apply(c.rootDir, step, templateDir)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadFile(stateFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read tfstate file %q: %v", stateFile, err)
		}

		result.Contents = append(result.Contents, asset.Content{
			Name: stateFile,
			Data: data,
		})
	}

	// TODO(yifan): Use the kubeconfig to verify the cluster is up.

	return &result, nil
}
