package image

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	config "github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	workflowreport "github.com/openshift/installer/pkg/asset/agent/workflow/report"
	"github.com/openshift/installer/pkg/asset/rhcos"
)

const (
	// bootArtifactsPath is the path where boot files are created.
	// e.g. initrd, kernel and rootfs.
	bootArtifactsPath = "boot-artifacts"
	// agentFilePrefix is the prefix used for day 1 images.
	agentFilePrefix = "agent"
	// nodeFilePrefix is the prefix used for day 2 images.
	nodeFilePrefix = "node"
)

// AgentArtifacts is an asset that generates all the artifacts that could be used
// for a subsequent generation of an ISO image or PXE files, starting from the
// content of the rhcos image enriched with agent specific files.
type AgentArtifacts struct {
	CPUArch              string
	RendezvousIP         string
	TmpPath              string
	IgnitionByte         []byte
	Kargs                string
	ISOPath              string
	BootArtifactsBaseURL string
	MinimalISO           bool
}

// Dependencies returns the assets on which the AgentArtifacts asset depends.
func (a *AgentArtifacts) Dependencies() []asset.Asset {
	return []asset.Asset{
		&Ignition{},
		&Kargs{},
		&BaseIso{},
		&manifests.AgentManifests{},
		&manifests.AgentClusterInstall{},
		&mirror.RegistriesConf{},
		&config.AgentConfig{},
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
	}
}

// Generate generates the configurations for the agent ISO image and PXE assets.
func (a *AgentArtifacts) Generate(ctx context.Context, dependencies asset.Parents) error {
	ignition := &Ignition{}
	kargs := &Kargs{}
	baseIso := &BaseIso{}
	agentManifests := &manifests.AgentManifests{}
	agentClusterInstall := &manifests.AgentClusterInstall{}
	registriesConf := &mirror.RegistriesConf{}
	agentconfig := &config.AgentConfig{}
	agentWorkflow := &workflow.AgentWorkflow{}
	dependencies.Get(ignition, kargs, baseIso, agentManifests, agentClusterInstall, registriesConf, agentconfig, agentWorkflow)

	if err := workflowreport.GetReport(ctx).Stage(workflow.StageAgentArtifacts); err != nil {
		return err
	}

	ignitionByte, err := json.Marshal(ignition.Config)
	if err != nil {
		return err
	}

	a.CPUArch = ignition.CPUArch
	a.RendezvousIP = ignition.RendezvousIP
	a.IgnitionByte = ignitionByte
	a.ISOPath = baseIso.File.Filename
	a.Kargs = kargs.KernelCmdLine()

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		if agentconfig.Config != nil {
			a.BootArtifactsBaseURL = strings.Trim(agentconfig.Config.BootArtifactsBaseURL, "/")
			// External platform will always create a minimal ISO
			a.MinimalISO = agentconfig.Config.MinimalISO || agentManifests.AgentClusterInstall.Spec.PlatformType == hiveext.ExternalPlatformType
			if agentconfig.Config.MinimalISO {
				logrus.Infof("Minimal ISO will be created based on configuration")
			} else if agentManifests.AgentClusterInstall.Spec.PlatformType == hiveext.ExternalPlatformType {
				logrus.Infof("Minimal ISO will be created for External platform")
			}
		}
	case workflow.AgentWorkflowTypeAddNodes:
		clusterInfo := &joiner.ClusterInfo{}
		dependencies.Get(clusterInfo)
		if clusterInfo.BootArtifactsBaseURL != "" {
			a.BootArtifactsBaseURL = strings.Trim(clusterInfo.BootArtifactsBaseURL, "/")
		}
	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

	var agentTuiFiles []string
	if agentClusterInstall.GetExternalPlatformName() != agent.ExternalPlatformNameOci {
		if err := workflowreport.GetReport(ctx).SubStage(workflow.StageAgentArtifactsAgentTUI); err != nil {
			return err
		}
		agentTuiFiles, err = a.fetchAgentTuiFiles(agentManifests.ClusterImageSet.Spec.ReleaseImage, agentManifests.GetPullSecretData(), registriesConf)
		if err != nil {
			return err
		}
	}

	if err := workflowreport.GetReport(ctx).SubStage(workflow.StageAgentArtifactsPrepare); err != nil {
		return err
	}
	err = a.prepareAgentArtifacts(a.ISOPath, agentTuiFiles)
	if err != nil {
		return err
	}

	return nil
}

func (a *AgentArtifacts) fetchAgentTuiFiles(releaseImage string, pullSecret string, mirrorConfig rhcos.MirrorConfig) ([]string, error) {
	release := rhcos.NewReleasePayload(
		rhcos.ExtractConfig{},
		releaseImage, pullSecret, mirrorConfig)

	agentTuiFilenames := []string{"/usr/bin/agent-tui", "/usr/lib64/libnmstate.so.*"}
	files := []string{}

	for _, srcFile := range agentTuiFilenames {
		extracted, err := release.ExtractFile("agent-installer-utils", srcFile, a.CPUArch)
		if err != nil {
			return nil, err
		}

		for _, f := range extracted {
			// Make sure it could be executed
			err = os.Chmod(f, 0o555)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
	}

	return files, nil
}

func (a *AgentArtifacts) prepareAgentArtifacts(iso string, additionalFiles []string) error {
	// Create a tmp folder to store all the pieces required to generate the agent artifacts.
	tmpPath, err := os.MkdirTemp("", "agent")
	if err != nil {
		return err
	}
	a.TmpPath = tmpPath

	err = isoeditor.Extract(iso, a.TmpPath)
	if err != nil {
		return err
	}

	err = a.appendAgentFilesToInitrd(additionalFiles)
	if err != nil {
		return err
	}

	return nil
}

func (a *AgentArtifacts) appendAgentFilesToInitrd(additionalFiles []string) error {
	ca := NewCpioArchive()

	dstPath := "/agent-files/"
	err := ca.StorePath(dstPath)
	if err != nil {
		return err
	}

	// Add the required agent files to the archive
	for _, f := range additionalFiles {
		err := ca.StoreFile(f, dstPath)
		if err != nil {
			return err
		}
	}

	// Add a dracut hook to copy the files. The $NEWROOT environment variable is exported by
	// dracut during the startup and it refers the mountpoint for the root filesystem.
	dracutHookScript := `#!/bin/sh
cp -R /agent-files/* $NEWROOT/usr/local/bin/
# Fix the selinux label
for i in $(find /agent-files/ -printf "%P\n"); do chcon system_u:object_r:bin_t:s0 $NEWROOT/usr/local/bin/$i; done`

	err = ca.StoreBytes("/usr/lib/dracut/hooks/pre-pivot/99-agent-copy-files.sh", []byte(dracutHookScript), 0o755)
	if err != nil {
		return err
	}

	buff, err := ca.SaveBuffer()
	if err != nil {
		return err
	}

	// Append the archive to initrd.img
	initrdImgPath := filepath.Join(a.TmpPath, "images", "pxeboot", "initrd.img")
	initrdImg, err := os.OpenFile(initrdImgPath, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer initrdImg.Close()

	_, err = initrdImg.Write(buff)
	if err != nil {
		return err
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *AgentArtifacts) Name() string {
	return "Agent Installer Artifacts"
}

func createDir(bootArtifactsFullPath string) error {
	os.RemoveAll(bootArtifactsFullPath)

	err := os.Mkdir(bootArtifactsFullPath, 0750)
	if err != nil {
		return err
	}
	return nil
}

func extractRootFS(bootArtifactsFullPath, agentISOPath, filePrefix, arch string) error {
	agentRootfsimgFile := filepath.Join(bootArtifactsFullPath, fmt.Sprintf("%s.%s-rootfs.img", filePrefix, arch))
	rootfsReader, err := os.Open(filepath.Join(agentISOPath, "images", "pxeboot", "rootfs.img"))
	if err != nil {
		return err
	}
	defer rootfsReader.Close()

	err = copyfile(agentRootfsimgFile, rootfsReader)
	if err != nil {
		return err
	}

	return nil
}
