package nodejoiner

import (
	"context"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	"github.com/openshift/installer/pkg/asset/store"
)

const (
	addNodesResultFile = "exit_code"
)

// NewAddNodesCommand creates a new command for add nodes.
func NewAddNodesCommand(directory string, kubeConfig string) error {
	err := saveParams(directory, kubeConfig)
	if err != nil {
		return err
	}

	fetcher := store.NewAssetsFetcher(directory)
	err = fetcher.FetchAndPersist(context.Background(), []asset.WritableAsset{
		&workflow.AgentWorkflowAddNodes{},
		&image.AgentImage{},
	})

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

func saveParams(directory, kubeConfig string) error {
	// Store the current parameters into the assets folder, so
	// that they could be retrieved later by the assets
	params := joiner.Params{
		Kubeconfig: kubeConfig,
	}
	return params.Save(directory)
}
