package nodejoiner

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/configimage"
	"github.com/openshift/installer/pkg/asset/agent/image"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	workflowreport "github.com/openshift/installer/pkg/asset/agent/workflow/report"
	"github.com/openshift/installer/pkg/asset/store"
)

const (
	addNodesResultFile = "exit_code"
)

// NewAddNodesCommand creates a new command for add nodes.
func NewAddNodesCommand(directory string, kubeConfig string, generatePXE bool, generateConfigISO bool) error {
	if generatePXE && generateConfigISO {
		return fmt.Errorf("invalid configuration found")
	}

	err := saveParams(directory, kubeConfig)
	if err != nil {
		return err
	}

	assets := []asset.WritableAsset{
		&workflow.AgentWorkflowAddNodes{},
	}
	var targetAsset asset.WritableAsset
	switch {
	case generatePXE:
		targetAsset = &image.AgentPXEFiles{}
	case generateConfigISO:
		targetAsset = &configimage.ConfigImage{}
	default:
		targetAsset = &image.AgentImage{}
	}
	assets = append(assets, targetAsset)

	ctx := workflowreport.Context(string(workflow.AgentWorkflowTypeAddNodes), directory)

	fetcher := store.NewAssetsFetcher(directory)
	err = fetcher.FetchAndPersist(ctx, assets)

	if reportErr := workflowreport.GetReport(ctx).Complete(err); reportErr != nil {
		return reportErr
	}

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
