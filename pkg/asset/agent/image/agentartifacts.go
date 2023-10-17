package image

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	"github.com/openshift/assisted-service/models"
	"github.com/openshift/installer/pkg/asset"
	config "github.com/openshift/installer/pkg/asset/agent/agentconfig"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
)

const (
	// bootArtifactsPath is the path where boot files are created.
	// e.g. initrd, kernel and rootfs.
	bootArtifactsPath = "boot-artifacts"
)

// AgentArtifacts is an asset that generates all the artifacts that could be used
// for a subsequent generation of an ISO image or PXE files, starting from the
// content of the rhcos image enriched with agent specific files.
type AgentArtifacts struct {
	CPUArch              string
	RendezvousIP         string
	TmpPath              string
	IgnitionByte         []byte
	Kargs                []byte
	ISOPath              string
	BootArtifactsBaseURL string
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
	}
}

// Generate generates the configurations for the agent ISO image and PXE assets.
func (a *AgentArtifacts) Generate(dependencies asset.Parents) error {
	ignition := &Ignition{}
	kargs := &Kargs{}
	baseIso := &BaseIso{}
	agentManifests := &manifests.AgentManifests{}
	agentClusterInstall := &manifests.AgentClusterInstall{}
	registriesConf := &mirror.RegistriesConf{}
	agentconfig := &config.AgentConfig{}

	dependencies.Get(ignition, kargs, baseIso, agentManifests, agentClusterInstall, registriesConf, agentconfig)

	ignitionByte, err := json.Marshal(ignition.Config)
	if err != nil {
		return err
	}

	a.CPUArch = ignition.CPUArch
	a.RendezvousIP = ignition.RendezvousIP
	a.IgnitionByte = ignitionByte
	a.ISOPath = baseIso.File.Filename
	a.Kargs = kargs.KernelCmdLine()

	if agentconfig.Config != nil {
		a.BootArtifactsBaseURL = strings.Trim(agentconfig.Config.BootArtifactsBaseURL, "/")
	}

	var agentTuiFiles []string
	if agentClusterInstall.GetExternalPlatformName() != string(models.PlatformTypeOci) {
		agentTuiFiles, err = a.fetchAgentTuiFiles(agentManifests.ClusterImageSet.Spec.ReleaseImage, agentManifests.GetPullSecretData(), registriesConf.MirrorConfig)
		if err != nil {
			return err
		}
	}
	err = a.prepareAgentArtifacts(a.ISOPath, agentTuiFiles)
	if err != nil {
		return err
	}

	return nil
}

func (a *AgentArtifacts) fetchAgentTuiFiles(releaseImage string, pullSecret string, mirrorConfig []mirror.RegistriesConfig) ([]string, error) {
	release := NewRelease(
		Config{MaxTries: OcDefaultTries, RetryDelay: OcDefaultRetryDelay},
		releaseImage, pullSecret, mirrorConfig)

	agentTuiFilenames := []string{"/usr/bin/agent-tui", "/usr/lib64/libnmstate.so.*"}
	files := []string{}

	for _, srcFile := range agentTuiFilenames {
		extracted, err := release.ExtractFile("agent-installer-utils", srcFile)
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

// Files returns the files generated by the asset.
func (a *AgentArtifacts) Files() []*asset.File {
	// Return empty array because File will never be loaded.
	return []*asset.File{}
}

func createDir(bootArtifactsFullPath string) error {
	os.RemoveAll(bootArtifactsFullPath)

	err := os.Mkdir(bootArtifactsFullPath, 0750)
	if err != nil {
		return err
	}
	return nil
}

func extractRootFS(bootArtifactsFullPath, agentISOPath, arch string) error {
	agentRootfsimgFile := filepath.Join(bootArtifactsFullPath, fmt.Sprintf("agent.%s-rootfs.img", arch))
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
