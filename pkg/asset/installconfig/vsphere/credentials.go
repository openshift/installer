package vsphere

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	ini "gopkg.in/ini.v1"

	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// DefaultCredentialsDir is the default directory for vSphere credentials.
	DefaultCredentialsDir = ".vsphere"
	// DefaultCredentialsFile is the default filename for vSphere credentials.
	DefaultCredentialsFile = "credentials"
	// CredentialsFileEnvVar is the environment variable for custom credentials file location.
	CredentialsFileEnvVar = "VSPHERE_CREDENTIALS_FILE"
)

// GetCredentialsFilePath returns the path to the credentials file, checking:
// 1. VSPHERE_CREDENTIALS_FILE environment variable
// 2. ~/.vsphere/credentials default location
func GetCredentialsFilePath() (string, error) {
	// Check environment variable first
	if path := os.Getenv(CredentialsFileEnvVar); path != "" {
		logrus.Debugf("Using credentials file from %s: %s", CredentialsFileEnvVar, path)
		return path, nil
	}

	// Use default location
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	path := filepath.Join(home, DefaultCredentialsDir, DefaultCredentialsFile)
	return path, nil
}

// LoadCredentialsFile loads and parses the ~/.vsphere/credentials file.
// Returns a map of vCenter server -> component credentials.
// Returns nil if the file doesn't exist (which is not an error).
func LoadCredentialsFile() (map[string]*vsphere.VCenterComponentCredentials, error) {
	path, err := GetCredentialsFilePath()
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logrus.Debugf("Credentials file does not exist at %s, skipping", path)
		return nil, nil
	}

	// Validate file permissions (must be 0600 or stricter)
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat credentials file: %w", err)
	}
	mode := info.Mode().Perm()
	// Check if group or others have any permissions
	if mode&0077 != 0 {
		return nil, fmt.Errorf("credentials file %s has insecure permissions %o (must be 0600 or stricter)", path, mode)
	}

	// Load INI file
	cfg, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials file %s: %w", path, err)
	}

	logrus.Infof("Loading vSphere credentials from %s", path)

	credentials := make(map[string]*vsphere.VCenterComponentCredentials)

	// Parse each section (vCenter server)
	for _, section := range cfg.Sections() {
		sectionName := section.Name()

		// Skip default section
		if sectionName == ini.DefaultSection {
			continue
		}

		// This is a vCenter server section
		vcenterCreds := &vsphere.VCenterComponentCredentials{}

		// Parse machine-api credentials
		if section.HasKey("machine-api.user") && section.HasKey("machine-api.password") {
			vcenterCreds.MachineAPI = &vsphere.VCenterCredential{
				User:     section.Key("machine-api.user").String(),
				Password: section.Key("machine-api.password").String(),
			}
		}

		// Parse csi-driver credentials
		if section.HasKey("csi-driver.user") && section.HasKey("csi-driver.password") {
			vcenterCreds.CSIDriver = &vsphere.VCenterCredential{
				User:     section.Key("csi-driver.user").String(),
				Password: section.Key("csi-driver.password").String(),
			}
		}

		// Parse cloud-controller credentials
		if section.HasKey("cloud-controller.user") && section.HasKey("cloud-controller.password") {
			vcenterCreds.CloudController = &vsphere.VCenterCredential{
				User:     section.Key("cloud-controller.user").String(),
				Password: section.Key("cloud-controller.password").String(),
			}
		}

		// Parse diagnostics credentials
		if section.HasKey("diagnostics.user") && section.HasKey("diagnostics.password") {
			vcenterCreds.Diagnostics = &vsphere.VCenterCredential{
				User:     section.Key("diagnostics.user").String(),
				Password: section.Key("diagnostics.password").String(),
			}
		}

		// Only add if at least one component credential is defined
		if vcenterCreds.MachineAPI != nil || vcenterCreds.CSIDriver != nil ||
			vcenterCreds.CloudController != nil || vcenterCreds.Diagnostics != nil {
			credentials[sectionName] = vcenterCreds
			logrus.Debugf("Loaded component credentials for vCenter %s", sectionName)
		}
	}

	return credentials, nil
}

// MergeCredentials merges credentials from file with install-config, with install-config taking precedence.
// This modifies the vcenters slice in-place.
func MergeCredentials(vcenters []vsphere.VCenter, fileCredentials map[string]*vsphere.VCenterComponentCredentials) {
	if fileCredentials == nil {
		return
	}

	for i := range vcenters {
		vcenter := &vcenters[i]

		// Check if we have file credentials for this vCenter
		fileCreds, exists := fileCredentials[vcenter.Server]
		if !exists {
			continue
		}

		// Only use file credentials if install-config doesn't have component credentials
		if vcenter.ComponentCredentials == nil {
			vcenter.ComponentCredentials = fileCreds
			logrus.Infof("Using component credentials from file for vCenter %s", vcenter.Server)
		} else {
			// Merge individual components - install-config takes precedence
			if vcenter.ComponentCredentials.MachineAPI == nil && fileCreds.MachineAPI != nil {
				vcenter.ComponentCredentials.MachineAPI = fileCreds.MachineAPI
				logrus.Debugf("Using machine-api credentials from file for vCenter %s", vcenter.Server)
			}
			if vcenter.ComponentCredentials.CSIDriver == nil && fileCreds.CSIDriver != nil {
				vcenter.ComponentCredentials.CSIDriver = fileCreds.CSIDriver
				logrus.Debugf("Using csi-driver credentials from file for vCenter %s", vcenter.Server)
			}
			if vcenter.ComponentCredentials.CloudController == nil && fileCreds.CloudController != nil {
				vcenter.ComponentCredentials.CloudController = fileCreds.CloudController
				logrus.Debugf("Using cloud-controller credentials from file for vCenter %s", vcenter.Server)
			}
			if vcenter.ComponentCredentials.Diagnostics == nil && fileCreds.Diagnostics != nil {
				vcenter.ComponentCredentials.Diagnostics = fileCreds.Diagnostics
				logrus.Debugf("Using diagnostics credentials from file for vCenter %s", vcenter.Server)
			}
		}
	}
}
