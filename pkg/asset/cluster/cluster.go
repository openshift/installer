package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types/config"
)

const (
	stateFileName = "terraform.state"
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

	// Copy the terraform.tfvars to a temp directory where the terraform will be invoked within.
	tmpDir, err := ioutil.TempDir(os.TempDir(), "openshift-install-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := ioutil.WriteFile(filepath.Join(tmpDir, state.Contents[0].Name), state.Contents[0].Data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write terraform.tfvars file: %v", err)
	}

	var tfvars config.Cluster
	if err := json.Unmarshal(state.Contents[0].Data, &tfvars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal terraform tfvars file: %v", err)
	}

	platform := string(tfvars.Platform)
	if err := data.Unpack(tmpDir, platform); err != nil {
		return nil, err
	}

	if err := data.Unpack(filepath.Join(tmpDir, "config.tf"), "config.tf"); err != nil {
		return nil, err
	}

	logrus.Infof("Using Terraform to create cluster...")

	// This runs the terraform in a temp directory, the tfstate file will be returned
	// to the asset store to persist it on the disk.
	if err := terraform.Init(tmpDir); err != nil {
		return nil, err
	}

	stateFile, err := terraform.Apply(tmpDir)
	if err != nil {
		err = fmt.Errorf("terraform failed: %v", err)
	}

	data, err2 := ioutil.ReadFile(stateFile)
	if err2 != nil {
		if err == nil {
			err = err2
		} else {
			logrus.Errorf("Failed to read tfstate (%q): %v", stateFile, err2)
		}
	}

	return &asset.State{
		Contents: []asset.Content{
			{
				Name: stateFileName,
				Data: data,
			},
		},
	}, err
}
