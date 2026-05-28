package isoeditor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/openshift/assisted-image-service/internal/common"
	log "github.com/sirupsen/logrus"
)

const (
	RamDiskPaddingLength        = uint64(1024 * 1024) // 1MB
	NmstatectlPathInRamdisk     = "/usr/bin/nmstatectl"
	ramDiskImagePath            = "/images/assisted_installer_custom.img"
	nmstateDiskImagePath        = "/images/nmstate.img"
	MinimalVersionForNmstatectl = "4.18.0-ec.0"
	RootfsImagePath             = "images/pxeboot/rootfs.img"
)

// transformKernelArgs applies the standard kernel argument transformations:
// 1. Remove coreos.liveiso parameter
// 2. Add coreos.live.rootfs_url parameter at the specified insertion point
func transformKernelArgs(content string, insertionPattern string, rootFSURL string, fileEntry *kargsFileEntry) (string, error) {
	// Validate rootfs URL
	if strings.Contains(rootFSURL, "$") {
		return "", fmt.Errorf("invalid rootfs URL: contains invalid character '$'")
	}
	if strings.Contains(rootFSURL, "\\") {
		return "", fmt.Errorf("invalid rootfs URL: contains invalid character '\\'")
	}

	var err error

	// Remove the coreos.liveiso parameter
	content, err = editString(content, `\b(?P<replace>coreos\.liveiso=\S+ ?)`, "", fileEntry)
	if err != nil {
		return "", err
	}

	// Add the rootfs_url parameter at the specified insertion point
	replacement := " coreos.live.rootfs_url=\"" + rootFSURL + "\""
	content, err = editString(content, insertionPattern, replacement, fileEntry)
	if err != nil {
		return "", err
	}

	return content, nil
}

//go:generate mockgen -package=isoeditor -destination=mock_editor.go . Editor
type Editor interface {
	CreateMinimalISOTemplate(fullISOPath, rootFSURL, arch, minimalISOPath, openshiftVersion, nmstatectlPath string) error
}

type rhcosEditor struct {
	workDir        string
	nmstateHandler NmstateHandler
}

func NewEditor(dataDir string, nmstateHandler NmstateHandler) Editor {
	return &rhcosEditor{
		workDir:        dataDir,
		nmstateHandler: nmstateHandler,
	}
}

// CreateMinimalISO Creates the minimal iso by removing the rootfs and adding the url
func CreateMinimalISO(extractDir, volumeID, rootFSURL, arch, minimalISOPath string) error {
	if err := os.Remove(filepath.Join(extractDir, RootfsImagePath)); err != nil {
		return err
	}

	if err := embedInitrdPlaceholders(extractDir); err != nil {
		log.WithError(err).Warnf("Failed to embed initrd placeholders")
		return err
	}

	var includeNmstateRamDisk bool
	if _, err := os.Stat(filepath.Join(extractDir, nmstateDiskImagePath)); err == nil {
		includeNmstateRamDisk = true
	}

	if err := updateKargs(extractDir, rootFSURL, includeNmstateRamDisk, arch); err != nil {
		log.WithError(err).Warnf("Failed to update kargs offsets and sizes")
		return err
	}

	if err := Create(minimalISOPath, extractDir, volumeID); err != nil {
		return err
	}
	return nil
}

// CreateMinimalISOTemplate Creates the template minimal iso by removing the rootfs and adding the url
func (e *rhcosEditor) CreateMinimalISOTemplate(fullISOPath, rootFSURL, arch, minimalISOPath, openshiftVersion, nmstatectlPath string) error {
	extractDir, err := os.MkdirTemp(e.workDir, "isoutil")
	if err != nil {
		return err
	}

	if err = Extract(fullISOPath, extractDir); err != nil {
		return err
	}

	volumeID, err := VolumeIdentifier(fullISOPath)
	if err != nil {
		return err
	}

	ramDiskPath := filepath.Join(extractDir, nmstateDiskImagePath)

	versionOK, err := common.VersionGreaterOrEqual(openshiftVersion, MinimalVersionForNmstatectl)
	if err != nil {
		return err
	}

	if versionOK {
		var compressedCpio []byte
		var readErr error

		if _, err = os.Stat(nmstatectlPath); err == nil {
			// Read and return the cached content
			compressedCpio, readErr = os.ReadFile(nmstatectlPath)
			if readErr != nil {
				return fmt.Errorf("failed to read cached nmstatectl: %v", readErr)
			}
		} else if os.IsNotExist(err) {
			// File doesn't exist - this should be an error condition
			return fmt.Errorf("nmstatectl cache file not found: %s", nmstatectlPath)
		} else {
			// Some other error occurred
			return fmt.Errorf("failed to stat nmstatectl cache file: %v", err)
		}

		err = os.WriteFile(ramDiskPath, compressedCpio, 0755) //nolint:gosec
		if err != nil {
			return err
		}
	}

	err = CreateMinimalISO(extractDir, volumeID, rootFSURL, arch, minimalISOPath)
	if err != nil {
		return err
	}

	return nil
}

func embedInitrdPlaceholders(extractDir string) error {
	f, err := os.Create(filepath.Join(extractDir, ramDiskImagePath))
	if err != nil {
		return err
	}
	defer func() {
		if deferErr := f.Sync(); deferErr != nil {
			log.WithError(deferErr).Error("Failed to sync disk image placeholder file")
		}
		if deferErr := f.Close(); deferErr != nil {
			log.WithError(deferErr).Error("Failed to close disk image placeholder file")
		}
	}()

	err = f.Truncate(int64(RamDiskPaddingLength))
	if err != nil {
		return err
	}

	return nil
}

// fixGrubConfig modifies grub.cfg and updates kargs config in place
func fixGrubConfig(rootFSURL, extractDir string, includeNmstateRamDisk bool, kargs *kargsConfig) error {
	availableGrubPaths := []string{"EFI/redhat/grub.cfg", "EFI/fedora/grub.cfg", "boot/grub/grub.cfg", "EFI/centos/grub.cfg"}
	var foundGrubPath string
	var fileEntry *kargsFileEntry
	for _, pathSection := range availableGrubPaths {
		path := filepath.Join(extractDir, pathSection)
		if _, err := os.Stat(path); err == nil {
			foundGrubPath = path
			fileEntry = kargs.FindFileByPath(pathSection)
			break
		}
	}
	if len(foundGrubPath) == 0 {
		return fmt.Errorf("no grub.cfg found, possible paths are %v", availableGrubPaths)
	}

	// Read the file content
	content, err := os.ReadFile(foundGrubPath)
	if err != nil {
		return err
	}
	contentStr := string(content)

	// Typical grub.cfg lines we're modifying:
	//
	//	linux /images/pxeboot/vmlinuz rw  coreos.liveiso=rhcos-9.6.20250523-0 ignition.firstboot ignition.platform.id=metal
	//	initrd /images/pxeboot/initrd.img /images/ignition.img

	// Apply standard kernel argument transformations (remove coreos.liveiso, add rootfs_url)
	contentStr, err = transformKernelArgs(contentStr, `(?m)^(\s+linux .+)(?P<replace>$)`, rootFSURL, fileEntry)
	if err != nil {
		return err
	}

	// Edit config to add custom ramdisk image to initrd - capture the end-of-line position to append our images
	var initrdReplacement string
	if includeNmstateRamDisk {
		initrdReplacement = " " + ramDiskImagePath + " " + nmstateDiskImagePath
	} else {
		initrdReplacement = " " + ramDiskImagePath
	}
	contentStr, err = editString(contentStr, `(?m)^(\s+initrd .+)(?P<replace>$)`, initrdReplacement, fileEntry)
	if err != nil {
		return err
	}

	// Write the modified content back to the file
	return os.WriteFile(foundGrubPath, []byte(contentStr), 0600)
}

// fixIsolinuxConfig modifies isolinux.cfg and updates kargs config in place
func fixIsolinuxConfig(rootFSURL, extractDir string, includeNmstateRamDisk bool, kargs *kargsConfig) error {
	relativeIsolinuxPath := strings.TrimPrefix(defaultIsolinuxFilePath, "/")
	isolinuxPath := filepath.Join(extractDir, relativeIsolinuxPath)

	// Read the file content
	content, err := os.ReadFile(isolinuxPath)
	if err != nil {
		return err
	}
	contentStr := string(content)

	// Typical isolinux.cfg line we're modifying:
	//
	//	append initrd=/images/pxeboot/initrd.img,/images/ignition.img rw  coreos.liveiso=rhcos-9.6.20250523-0 ignition.firstboot ignition.platform.id=metal

	// Find the kargs file entry for this file
	fileEntry := kargs.FindFileByPath(relativeIsolinuxPath)

	// Apply standard kernel argument transformations (remove coreos.liveiso, add rootfs_url)
	contentStr, err = transformKernelArgs(contentStr, `(?m)^(\s+append .+)(?P<replace>$)`, rootFSURL, fileEntry)
	if err != nil {
		return err
	}

	// Add ramdisk images to initrd specification - capture the position right after the initrd argument to append our images
	var initrdReplacement string
	if includeNmstateRamDisk {
		initrdReplacement = "," + ramDiskImagePath + "," + nmstateDiskImagePath
	} else {
		initrdReplacement = "," + ramDiskImagePath
	}
	contentStr, err = editString(contentStr, `(?m)^(\s+append.*initrd=\S+)(?P<replace>)`, initrdReplacement, fileEntry)
	if err != nil {
		return err
	}

	// Write the modified content back to the file
	return os.WriteFile(isolinuxPath, []byte(contentStr), 0600)
}

// editString applies a regex replacement to a string and returns the modified string
// It looks for a named capture group called "replace" and replaces only that content, using precise string manipulation
func editString(content string, reString string, replacement string, fileEntry *kargsFileEntry) (string, error) {
	re := regexp.MustCompile(reString)

	// Get the index of the "replace" named capture group
	replaceIndex := re.SubexpIndex("replace")
	if replaceIndex == -1 {
		return "", fmt.Errorf("regex must have a named capture group called 'replace'")
	}

	// Find the first match with subgroups
	submatchIndexes := re.FindStringSubmatchIndex(content)
	if submatchIndexes == nil {
		// No match found
		return content, nil
	}

	// submatchIndexes contains [fullMatchStart, fullMatchEnd, group1Start, group1End, group2Start, group2End, ...]
	if len(submatchIndexes) < (replaceIndex+1)*2 {
		return "", fmt.Errorf("regex match does not contain the 'replace' capture group")
	}

	replaceStart := submatchIndexes[replaceIndex*2]
	replaceEnd := submatchIndexes[replaceIndex*2+1]

	// Replace only the "replace" capturing group
	newContent := content[:replaceStart] + replacement + content[replaceEnd:]

	if content == newContent {
		return content, nil
	}

	if fileEntry == nil || fileEntry.Offset == nil {
		return newContent, nil
	}

	embedStart := *fileEntry.Offset
	fileSizeChange := int64(len(newContent)) - int64(len(content))

	replaceStartPos := int64(replaceStart)
	replaceEndPos := int64(replaceEnd)

	// Add boundary crossing check to ensure no replacements span across the embed area start boundary
	if replaceStartPos < embedStart && replaceEndPos > embedStart {
		return "", fmt.Errorf("replacement spans across embed area boundary (replace: %d-%d, embed starts: %d)", replaceStartPos, replaceEndPos, embedStart)
	}

	if replaceEndPos <= embedStart {
		// Change is before embed area - affects offset
		*fileEntry.Offset += fileSizeChange
	}

	return newContent, nil
}

type kargsFileEntry struct {
	End    *string `json:"end,omitempty"`
	Offset *int64  `json:"offset,omitempty"`
	Pad    *string `json:"pad,omitempty"`
	Path   *string `json:"path,omitempty"`
}

type kargsConfig struct {
	Default string           `json:"default"`
	Files   []kargsFileEntry `json:"files"`
	Size    int64            `json:"size"`
}

// FindFileByPath searches for a file entry by its path and returns a pointer to it.
// Returns nil if no file with the specified path is found.
func (k *kargsConfig) FindFileByPath(path string) *kargsFileEntry {
	if k == nil {
		return nil
	}
	for i := range k.Files {
		if k.Files[i].Path != nil && *k.Files[i].Path == path {
			return &k.Files[i]
		}
	}
	return nil
}

// updateDefaultKargs modifies the default kernel arguments to match bootloader modifications
// and applies the embed area size change directly to config.Size
func updateDefaultKargs(config *kargsConfig, rootFSURL string) error {
	originalDefault := config.Default

	// Apply the same transformations we make to bootloader configs
	// For default kargs, we append at the end (using $ insertion pattern)
	// Pass nil for fileEntry since we don't track offsets for the default string
	var err error
	config.Default, err = transformKernelArgs(config.Default, `(?P<replace>$)`, rootFSURL, nil)
	if err != nil {
		return err
	}

	// Calculate and apply the embed area size change from the default kargs transformation
	embedSizeChange := int64(len(config.Default)) - int64(len(originalDefault))
	config.Size += embedSizeChange

	return nil
}

// updateKargs reads kargs.json, applies fixes with embed area awareness, and updates kargs.json
func updateKargs(extractDir, rootFSURL string, includeNmstateRamDisk bool, arch string) error {
	kargsPath := filepath.Join(extractDir, "coreos/kargs.json")

	var config *kargsConfig

	if _, err := os.Stat(kargsPath); !os.IsNotExist(err) {
		kargsData, err := os.ReadFile(kargsPath)
		if err != nil {
			return fmt.Errorf("failed to read kargs.json: %v", err)
		}

		var kargsStruct kargsConfig
		if err := json.Unmarshal(kargsData, &kargsStruct); err != nil {
			return fmt.Errorf("failed to unmarshal kargs.json: %v", err)
		}
		config = &kargsStruct
	}

	// Apply bootloader config changes (without tracking size changes)
	if err := fixGrubConfig(rootFSURL, extractDir, includeNmstateRamDisk, config); err != nil {
		return fmt.Errorf("failed to fix grub config: %v", err)
	}

	// ignore isolinux.cfg for ppc64le because it doesn't exist
	if arch != "ppc64le" {
		if err := fixIsolinuxConfig(rootFSURL, extractDir, includeNmstateRamDisk, config); err != nil {
			return fmt.Errorf("failed to fix isolinux config: %v", err)
		}
	}

	if config != nil {
		// Update the default kernel arguments and apply embed area size change
		if err := updateDefaultKargs(config, rootFSURL); err != nil {
			return fmt.Errorf("failed to update default kargs: %v", err)
		}

		updatedData, err := json.MarshalIndent(*config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal updated kargs.json: %v", err)
		}

		if err := os.WriteFile(kargsPath, updatedData, 0600); err != nil {
			return fmt.Errorf("failed to write updated kargs.json: %v", err)
		}
	}

	return nil
}
