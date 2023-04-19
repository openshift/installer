package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	"github.com/openshift/assisted-service/pkg/executer"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/mirror"
)

const (
	agentISOFilename = "agent.%s.iso"
)

// AgentImage is an asset that generates the bootable image used to install clusters.
type AgentImage struct {
	cpuArch      string
	rendezvousIP string

	tmpPath  string
	volumeID string
}

var _ asset.WritableAsset = (*AgentImage)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *AgentImage) Dependencies() []asset.Asset {
	return []asset.Asset{
		&Ignition{},
		&BaseIso{},
		&manifests.AgentManifests{},
		&mirror.RegistriesConf{},
	}
}

// Generate generates the image file for to ISO asset.
func (a *AgentImage) Generate(dependencies asset.Parents) error {
	ignition := &Ignition{}
	baseImage := &BaseIso{}
	agentManifests := &manifests.AgentManifests{}
	registriesConf := &mirror.RegistriesConf{}

	dependencies.Get(ignition, baseImage, agentManifests, registriesConf)

	ignitionByte, err := json.Marshal(ignition.Config)
	if err != nil {
		return err
	}

	agentTuiFiles, err := a.fetchAgentTuiFiles(agentManifests.ClusterImageSet.Spec.ReleaseImage, agentManifests.GetPullSecretData(), registriesConf.MirrorConfig)
	if err != nil {
		return err
	}

	err = a.prepareAgentISO(baseImage.File.Filename, ignitionByte, agentTuiFiles)
	if err != nil {
		return err
	}

	a.cpuArch = ignition.CPUArch
	a.rendezvousIP = ignition.RendezvousIP

	return nil
}

func (a *AgentImage) fetchAgentTuiFiles(releaseImage string, pullSecret string, mirrorConfig []mirror.RegistriesConfig) ([]string, error) {
	release := NewRelease(&executer.CommonExecuter{},
		Config{MaxTries: OcDefaultTries, RetryDelay: OcDefaultRetryDelay},
		releaseImage, pullSecret, mirrorConfig)

	agentTuiFilenames := []string{"/usr/bin/agent-tui", "/usr/lib64/libnmstate.so.*"}
	files := []string{}

	for _, srcFile := range agentTuiFilenames {
		extracted, err := release.ExtractFile("agent-installer-node-agent", srcFile)
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

func (a *AgentImage) prepareAgentISO(iso string, ignition []byte, additionalFiles []string) error {
	// Create a tmp folder to store all the pieces required to generate the agent ISO.
	tmpPath, err := os.MkdirTemp("", "agent")
	if err != nil {
		return err
	}
	a.tmpPath = tmpPath

	err = isoeditor.Extract(iso, a.tmpPath)
	if err != nil {
		return err
	}

	err = a.updateIgnitionImg(ignition)
	if err != nil {
		return err
	}

	err = a.appendAgentFilesToInitrd(additionalFiles)
	if err != nil {
		return err
	}

	volumeID, err := isoeditor.VolumeIdentifier(iso)
	if err != nil {
		return err
	}
	a.volumeID = volumeID

	return nil
}

func (a *AgentImage) updateIgnitionImg(ignition []byte) error {
	ca := NewCpioArchive()
	err := ca.StoreBytes("config.ign", ignition, 0o644)
	if err != nil {
		return err
	}
	ignitionBuff, err := ca.SaveBuffer()
	if err != nil {
		return err
	}

	ignitionImgPath := filepath.Join(a.tmpPath, "images", "ignition.img")
	fi, err := os.Stat(ignitionImgPath)
	if err != nil {
		return err
	}

	// Verify that the current compressed ignition archive does not exceed the
	// embed area (usually 256 Kb)
	if len(ignitionBuff) > int(fi.Size()) {
		return fmt.Errorf("ignition content length (%d) exceeds embed area size (%d)", len(ignitionBuff), fi.Size())
	}

	ignitionImg, err := os.OpenFile(ignitionImgPath, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer ignitionImg.Close()

	_, err = ignitionImg.Write(ignitionBuff)
	if err != nil {
		return err
	}

	return nil
}

func (a *AgentImage) appendAgentFilesToInitrd(additionalFiles []string) error {
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
	initrdImgPath := filepath.Join(a.tmpPath, "images", "pxeboot", "initrd.img")
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

// PersistToFile writes the iso image in the assets folder
func (a *AgentImage) PersistToFile(directory string) error {
	defer os.RemoveAll(a.tmpPath)

	// If the volumeId or tmpPath are not set then it means that either one of the AgentImage
	// dependencies or the asset itself failed for some reason
	if a.tmpPath == "" || a.volumeID == "" {
		return errors.New("cannot generate ISO image due to configuration errors")
	}

	agentIsoFile := filepath.Join(directory, fmt.Sprintf(agentISOFilename, a.cpuArch))

	// Remove symlink if it exists
	os.Remove(agentIsoFile)

	err := isoeditor.Create(agentIsoFile, a.tmpPath, a.volumeID)
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(directory, "rendezvousIP"), []byte(a.rendezvousIP), 0o644) //nolint:gosec // no sensitive info
	if err != nil {
		return err
	}

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *AgentImage) Name() string {
	return "Agent Installer ISO"
}

// Load returns the ISO from disk.
func (a *AgentImage) Load(f asset.FileFetcher) (bool, error) {
	// The ISO will not be needed by another asset so load is noop.
	// This is implemented because it is required by WritableAsset
	return false, nil
}

// Files returns the files generated by the asset.
func (a *AgentImage) Files() []*asset.File {
	// Return empty array because File will never be loaded.
	return []*asset.File{}
}
