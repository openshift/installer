package nodejoiner

import (
	"context"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/configimage"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/store"
)

const (
	addNodesResultFile = "exit_code"
)

// NewAddNodesCommand creates a new command for add nodes.
func NewAddNodesCommand(directory string, kubeConfig string, generatePXE bool) error {
	err := saveParams(directory, kubeConfig, generatePXE)
	if err != nil {
		return err
	}

	assets := []asset.WritableAsset{
		&workflow.AgentWorkflowAddNodes{},
	}
	if generatePXE {
		assets = append(assets, &image.AgentPXEFiles{})
	} else {
		assets = append(assets, &image.AgentImage{})
		assets = append(assets, &configimage.ConfigImage{})
	}

	fetcher := store.NewAssetsFetcher(directory)
	err = fetcher.FetchAndPersist(context.Background(), assets)

	// Save the exit code result
	exitCode := "0"
	if err != nil {
		exitCode = "1"
	}
	if err2 := os.WriteFile(filepath.Join(directory, addNodesResultFile), []byte(exitCode), 0600); err2 != nil {
		return err2
	}

	return err
}

func saveParams(directory, kubeConfig string, generatePXE bool) error {
	// Store the current parameters into the assets folder, so
	// that they could be retrieved later by the assets
	params := joiner.Params{
		Kubeconfig:  kubeConfig,
		GeneratePXE: generatePXE,
	}
	return params.Save(directory)
}
