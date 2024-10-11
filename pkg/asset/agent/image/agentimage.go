package image

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/openshift/assisted-image-service/pkg/isoeditor"
	hiveext "github.com/openshift/assisted-service/api/hiveextension/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/gencrypto"
	"github.com/openshift/installer/pkg/asset/agent/joiner"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
)

const (
	agentISOFilename         = "agent.%s.iso"
	agentAddNodesISOFilename = "node.%s.iso"
	iso9660Level1ExtLen      = 3
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
	isoFilename          string
	imageExpiresAt       string
}

var _ asset.WritableAsset = (*AgentImage)(nil)

// Dependencies returns the assets on which the Bootstrap asset depends.
func (a *AgentImage) Dependencies() []asset.Asset {
	return []asset.Asset{
		&workflow.AgentWorkflow{},
		&joiner.ClusterInfo{},
		&AgentArtifacts{},
		&manifests.AgentManifests{},
		&BaseIso{},
		&gencrypto.AuthConfig{},
	}
}

// Generate generates the image file for to ISO asset.
func (a *AgentImage) Generate(ctx context.Context, dependencies asset.Parents) error {
	agentWorkflow := &workflow.AgentWorkflow{}
	clusterInfo := &joiner.ClusterInfo{}
	agentArtifacts := &AgentArtifacts{}
	agentManifests := &manifests.AgentManifests{}
	baseIso := &BaseIso{}
	dependencies.Get(agentArtifacts, agentManifests, baseIso, agentWorkflow, clusterInfo)

	switch agentWorkflow.Workflow {
	case workflow.AgentWorkflowTypeInstall:
		a.platform = agentManifests.AgentClusterInstall.Spec.PlatformType
		a.isoFilename = agentISOFilename

	case workflow.AgentWorkflowTypeAddNodes:
		authConfig := &gencrypto.AuthConfig{}
		dependencies.Get(authConfig)

		a.platform = clusterInfo.PlatformType
		a.isoFilename = agentAddNodesISOFilename
		a.imageExpiresAt = authConfig.AuthTokenExpiry

	default:
		return fmt.Errorf("AgentWorkflowType value not supported: %s", agentWorkflow.Workflow)
	}

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

	if a.platform == hiveext.ExternalPlatformType {
		// when the bootArtifactsBaseURL is specified, construct the custom rootfs URL
		if a.bootArtifactsBaseURL != "" {
			a.rootFSURL = fmt.Sprintf("%s/%s", a.bootArtifactsBaseURL, fmt.Sprintf("agent.%s-rootfs.img", a.cpuArch))
			logrus.Debugf("Using custom rootfs URL: %s", a.rootFSURL)
		} else {
			// Default to the URL from the RHCOS streams file
			defaultRootFSURL, err := baseIso.getRootFSURL(ctx, a.cpuArch)
			if err != nil {
				return err
			}
			a.rootFSURL = defaultRootFSURL
			logrus.Debugf("Using default rootfs URL: %s", a.rootFSURL)
		}
	}

	// Update Ignition images
	err = a.updateIgnitionContent(agentArtifacts)
	if err != nil {
		return err
	}

	err = a.appendKargs(agentArtifacts.Kargs)
	if err != nil {
		return err
	}

	return nil
}

// updateIgnitionContent updates the ignition data into the corresponding images in the ISO.
func (a *AgentImage) updateIgnitionContent(agentArtifacts *AgentArtifacts) error {
	ignitionc := &isoeditor.IgnitionContent{}
	ignitionc.Config = agentArtifacts.IgnitionByte
	fileInfo, err := isoeditor.NewIgnitionImageReader(a.isoPath, ignitionc)
	if err != nil {
		return err
	}

	return a.overwriteFileData(fileInfo)
}

func (a *AgentImage) overwriteFileData(fileInfo []isoeditor.FileData) error {
	var errs []error
	for _, fileData := range fileInfo {
		defer fileData.Data.Close()

		filename := filepath.Join(a.tmpPath, fileData.Filename)
		file, err := os.Create(filename)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer file.Close()

		_, err = io.Copy(file, fileData.Data)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (a *AgentImage) appendKargs(kargs string) error {
	if kargs == "" {
		return nil
	}

	fileInfo, err := isoeditor.NewKargsReader(a.isoPath, kargs)
	if err != nil {
		return err
	}
	return a.overwriteFileData(fileInfo)
}

// normalizeFilesExtension scans the extracted ISO files and trims
// the file extensions longer than three chars.
func (a *AgentImage) normalizeFilesExtension() error {
	var skipFiles = map[string]bool{
		"boot.catalog": true, // Required for arm64 iso
	}

	return filepath.WalkDir(a.tmpPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(p)
		// ext includes also the dot separator
		if len(ext) > iso9660Level1ExtLen+1 {
			b := filepath.Base(p)
			if _, ok := skipFiles[filepath.Base(b)]; ok {
				return nil
			}

			// Replaces file extensions longer than three chars
			np := p[:len(p)-len(ext)] + ext[:iso9660Level1ExtLen+1]
			err = os.Rename(p, np)

			if err != nil {
				return err
			}
		}

		return nil
	})
}

// PersistToFile writes the iso image in the assets folder
func (a *AgentImage) PersistToFile(directory string) error {
	defer os.RemoveAll(a.tmpPath)

	// If the volumeId or tmpPath are not set then it means that either one of the AgentImage
	// dependencies or the asset itself failed for some reason
	if a.tmpPath == "" || a.volumeID == "" {
		return errors.New("cannot generate ISO image due to configuration errors")
	}

	agentIsoFile := filepath.Join(directory, fmt.Sprintf(a.isoFilename, a.cpuArch))

	// Remove symlink if it exists
	os.Remove(agentIsoFile)

	err := a.normalizeFilesExtension()
	if err != nil {
		return err
	}

	var msg string
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
		msg = fmt.Sprintf("Generated minimal ISO at %s", agentIsoFile)
	} else {
		// Generate full ISO
		err = isoeditor.Create(agentIsoFile, a.tmpPath, a.volumeID)
		if err != nil {
			return err
		}
		msg = fmt.Sprintf("Generated ISO at %s.", agentIsoFile)
	}
	if a.imageExpiresAt != "" {
		msg = fmt.Sprintf("%s The ISO is valid up to %s", msg, a.imageExpiresAt)
	}
	logrus.Info(msg)

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
