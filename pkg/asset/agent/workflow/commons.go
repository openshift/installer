package workflow

import (
	wr "github.com/openshift/installer/pkg/asset/agent/workflow/report"
)

// AgentWorkflowType defines the supported
// agent workflows.
type AgentWorkflowType string

const (
	// AgentWorkflowTypeInstall identifies the install workflow.
	AgentWorkflowTypeInstall AgentWorkflowType = "install"
	// AgentWorkflowTypeAddNodes identifies the add nodes workflow.
	AgentWorkflowTypeAddNodes AgentWorkflowType = "addnodes"
	// AgentWorkflowTypeInstallInteractiveDisconnected identifies a specific kind of
	// disconnected install workflow. The installation details will be provided through
	// a dedicated UI running on the rendezvous node, and in addition no external registry
	// will be required for an air-gapped deployment.
	AgentWorkflowTypeInstallInteractiveDisconnected = "install-interactive-disconnected"

	agentWorkflowFilename = ".agentworkflow"
)

var (
	// StageClusterInspection represents cluster inspection stage.
	StageClusterInspection wr.StageID = wr.NewStageID("add-nodes-cluster-inspection", "Gathering additional information from the target cluster")

	// StageCreateManifests represents the manifests creation stage.
	StageCreateManifests wr.StageID = wr.NewStageID("create-manifest", "Creating internal configuration manifests")

	// StageIgnition represents the ignition creation stage.
	StageIgnition wr.StageID = wr.NewStageID("ignition", "Rendering ISO ignition")

	// StageFetchBaseISO represents the base image fetching stage.
	StageFetchBaseISO wr.StageID = wr.NewStageID("fetch-base-iso", "Retrieving the base ISO image")
	// StageFetchBaseISOExtract represents the base image extraction substage.
	StageFetchBaseISOExtract wr.StageID = wr.NewStageID("fetch-base-iso.extract-image", "Extracting base image from release payload")
	// StageFetchBaseISOVerify represents the image version verification substage.
	StageFetchBaseISOVerify wr.StageID = wr.NewStageID("fetch-base-iso.verify-version", "Verifying base image version")
	// StageFetchBaseISODownload represents the base image download substage.
	StageFetchBaseISODownload wr.StageID = wr.NewStageID("fetch-base-iso.download-image", "Downloading base ISO image")

	// StageAgentArtifacts represents the agent artifact stage.
	StageAgentArtifacts wr.StageID = wr.NewStageID("create-agent-artifacts", "Creating agent artifacts for the final image")
	// StageAgentArtifactsAgentTUI represents the agent-tui embedding substage.
	StageAgentArtifactsAgentTUI wr.StageID = wr.NewStageID("create-agent-artifacts.agent-tui", "Extracting required artifacts from release payload")
	// StageAgentArtifactsPrepare represents the artifacts preparation substage.
	StageAgentArtifactsPrepare wr.StageID = wr.NewStageID("create-agent-artifacts.prepare", "Preparing artifacts")

	// StageGenerateISO represents the iso assembling stage.
	StageGenerateISO wr.StageID = wr.NewStageID("generate-iso", "Assembling ISO image")

	// StageGeneratePXE represents the pxe assembling stage.
	StageGeneratePXE wr.StageID = wr.NewStageID("generate-pxe", "Assembling PXE files")
)
