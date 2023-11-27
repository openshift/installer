package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
)

const (
	agentISOFilename = "agent.%s.iso"
)

// AgentImage is an asset that generates the bootable image used to install clusters.
type AgentImage struct {
	cpuArch              string
	rendezvousIP         string
	tmpPath              string
	volumeID             string
	isoPath              string
	rootFSURL            string
	bootArtifactsBaseURL string
	platform             hiveext.PlatformType
}

var _ asset.WritableAsset = (*AgentImage)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *AgentImage) Dependencies() []asset.Asset {
	return []asset.Asset{
		&AgentArtifacts{},
		&manifests.AgentManifests{},
		&BaseIso{},
	}
}

// Generate generates the image file for to ISO asset.
func (a *AgentImage) Generate(dependencies asset.Parents) error {
	agentArtifacts := &AgentArtifacts{}
	agentManifests := &manifests.AgentManifests{}
	baseIso := &BaseIso{}
	dependencies.Get(agentArtifacts, agentManifests, baseIso)

	a.cpuArch = agentArtifacts.CPUArch
	a.rendezvousIP = agentArtifacts.RendezvousIP
	a.tmpPath = agentArtifacts.TmpPath
	a.isoPath = agentArtifacts.ISOPath
	a.bootArtifactsBaseURL = agentArtifacts.BootArtifactsBaseURL

	volumeID, err := isoeditor.VolumeIdentifier(a.isoPath)
	if err != nil {
		return err
	}
	a.volumeID = volumeID

	a.platform = agentManifests.AgentClusterInstall.Spec.PlatformType
	if a.platform == hiveext.ExternalPlatformType {
		// when the bootArtifactsBaseURL is specified, construct the custom rootfs URL
		if a.bootArtifactsBaseURL != "" {
			a.rootFSURL = fmt.Sprintf("%s/%s", a.bootArtifactsBaseURL, fmt.Sprintf("agent.%s-rootfs.img", a.cpuArch))
			logrus.Debugf("Using custom rootfs URL: %s", a.rootFSURL)
		} else {
			// Default to the URL from the RHCOS streams file
			defaultRootFSURL, err := baseIso.getRootFSURL(a.cpuArch)
			if err != nil {
				return err
			}
			a.rootFSURL = defaultRootFSURL
			logrus.Debugf("Using default rootfs URL: %s", a.rootFSURL)
		}
	}

	logrus.Infof("Updating the ignition data in the image file")
	err = a.updateIgnitionImg(agentArtifacts.IgnitionByte)
	if err != nil {
		return err
	}

	err = a.appendKargs(agentArtifacts.Kargs)
	if err != nil {
		return err
	}

	return nil
}

// Update the ignition image file with the ignition data
func (a *AgentImage) updateIgnitionImg(ignition []byte) error {

	// Saving the ignition buffer into archive
	ca := NewCpioArchive()
	err := ca.StoreBytes("config.ign", ignition, 0o644)
	if err != nil {
		return err
	}
	ignitionBuff, err := ca.SaveBuffer()
	if err != nil {
		return err
	}

	// Defining a struct for the ignition info
	type IgnInfo struct {
		File   string `json:"file"`
		Length int64  `json:"length"`
		Offset int64  `json:"offset"`
	}
	var ignInfo IgnInfo

	// Reading the ignition details from igninfo.json file into the struct
	ignInfoJSONPath := filepath.Join(a.tmpPath, "coreos", "igninfo.json")
	ignInfoJSONData, err := os.ReadFile(ignInfoJSONPath)
	if err != nil {
		logrus.Debugf("Failed to read json %s", ignInfoJSONPath)
		return err
	}
	if err := json.Unmarshal(ignInfoJSONData, &ignInfo); err != nil {
		logrus.Debugf("Failed to umarshal json: %s", err)
		return err
	}

	// Fetching the ignition details from the struct
	ignImagePath := filepath.Join(a.tmpPath, ignInfo.File)
	ignStartOffset := ignInfo.Offset
	ignMaxLength := ignInfo.Length

	// Verify that the compressed ignition buffer archive does not exceed the embed area (usually 256 Kb)
	if ignMaxLength == 0 {
		// Checking the file size if no ignition length is provided
		fileInfo, err := os.Stat(ignImagePath)
		if err != nil {
			return err
		}

		if len(ignitionBuff) > int(fileInfo.Size()) {
			return fmt.Errorf("ignition content length (%d) exceeds embed area size (%d)", len(ignitionBuff), fileInfo.Size())
		}
	} else {
		if int64(len(ignitionBuff)) > ignMaxLength {
			return fmt.Errorf("Ignition content length (%d) exceeds embed area size (%d)", len(ignitionBuff), ignMaxLength)
		}
	}

	// Open the image file with permissions to seek and write the ignition content
	ignitionImg, err := os.OpenFile(ignImagePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer ignitionImg.Close()

	// Adjusting the ignition image file to the offset of ignition content (in case of cdboot.img)
	_, err = ignitionImg.Seek(ignStartOffset, 0)
	if err != nil {
		return err
	}

	// Write the ignition content into the image file
	_, err = ignitionImg.Write(ignitionBuff)
	if err != nil {
		return err
	}

	// Padding 0's at the end if the ignition content has less bytes in the image file (in case of cdboot.img)
	paddingLength := ignMaxLength - int64(len(ignitionBuff))
	if paddingLength > 0 {
		padding := make([]byte, paddingLength)
		_, err = ignitionImg.Write(padding)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateKargsFile(tmpPath, filename string, embedArea *regexp.Regexp, kargs []byte) error {
	file, err := os.OpenFile(filepath.Join(tmpPath, filename), os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	indices := embedArea.FindSubmatchIndex(content)
	if len(indices) != 4 {
		return fmt.Errorf("failed to find COREOS_KARG_EMBED_AREA in %s", filename)
	}

	if size := (indices[3] - indices[2]); len(kargs) > size {
		return fmt.Errorf("kernel args content length (%d) exceeds embed area size (%d)", len(kargs), size)
	}

	if _, err := file.WriteAt(append(kargs, '\n'), int64(indices[2])); err != nil {
		return err
	}
	return nil
}

func (a *AgentImage) appendKargs(kargs []byte) error {
	if len(kargs) == 0 {
		return nil
	}

	kargsFiles, err := isoeditor.KargsFiles(a.isoPath)
	if err != nil {
		return err
	}

	embedArea := regexp.MustCompile(`(\n#*)# COREOS_KARG_EMBED_AREA`)
	for _, f := range kargsFiles {
		if err := updateKargsFile(a.tmpPath, f, embedArea, kargs); err != nil {
			return err
		}
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

	var err error
	// For external platform when the bootArtifactsBaseURL is specified,
	// output the rootfs file alongside the minimal ISO
	if a.platform == hiveext.ExternalPlatformType {
		if a.bootArtifactsBaseURL != "" {
			bootArtifactsFullPath := filepath.Join(directory, bootArtifactsPath)
			err := createDir(bootArtifactsFullPath)
			if err != nil {
				return err
			}
			err = extractRootFS(bootArtifactsFullPath, a.tmpPath, a.cpuArch)
			if err != nil {
				return err
			}
			logrus.Infof("RootFS file created in: %s. Upload it at %s", bootArtifactsFullPath, a.rootFSURL)
		}
		err = isoeditor.CreateMinimalISO(a.tmpPath, a.volumeID, a.rootFSURL, a.cpuArch, agentIsoFile)
		if err != nil {
			return err
		}
		logrus.Infof("Generated minimal ISO at %s", agentIsoFile)
	} else {
		// Generate full ISO
		err = isoeditor.Create(agentIsoFile, a.tmpPath, a.volumeID)
		if err != nil {
			return err
		}
		logrus.Infof("Generated ISO at %s", agentIsoFile)
	}

	err = os.WriteFile(filepath.Join(directory, "rendezvousIP"), []byte(a.rendezvousIP), 0o644) //nolint:gosec // no sensitive info
	if err != nil {
		return err
	}
	// For external platform OCI, add CCM manifests in the openshift directory.
	if a.platform == hiveext.ExternalPlatformType {
		logrus.Infof("When using %s oci platform, always make sure CCM manifests were added in the %s directory.", hiveext.ExternalPlatformType, manifests.OpenshiftManifestDir())
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
