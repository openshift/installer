package installconfig

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
	azureconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	gcpconfig "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/sirupsen/logrus"
)

// PlatformCredsCheck is an asset that checks the platform credentials, asks for them or errors out if invalid
// the cluster.
type PlatformCredsCheck struct {
}

var _ asset.Asset = (*PlatformCredsCheck)(nil)

// Dependencies returns the dependencies for PlatformCredsCheck
func (a *PlatformCredsCheck) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate queries for input from the user.
func (a *PlatformCredsCheck) Generate(dependencies asset.Parents) error {
	ctx := context.TODO()
	ic := &InstallConfig{}
	dependencies.Get(ic)

	var err error
	platform := ic.Config.Platform.Name()
	switch platform {
	case aws.Name:
		permissionGroups := []awsconfig.PermissionGroup{awsconfig.PermissionCreateBase, awsconfig.PermissionDeleteBase}
		// If subnets are not provided in install-config.yaml, include network permissions
		if len(ic.Config.AWS.Subnets) == 0 {
			permissionGroups = append(permissionGroups, awsconfig.PermissionCreateNetworking, awsconfig.PermissionDeleteNetworking)
		}

		ssn, err := ic.AWS.Session(ctx)
		if err != nil {
			return err
		}
		err = awsconfig.ValidateCreds(ssn, permissionGroups)
		if err != nil {
			return errors.Wrap(err, "validate AWS credentials")
		}
	case gcp.Name:
		_, err = gcpconfig.GetSession(ctx)
		if err != nil {
			return errors.Wrap(err, "creating GCP session")
		}
	case openstack.Name:
		opts := new(clientconfig.ClientOpts)
		opts.Cloud = ic.Config.Platform.OpenStack.Cloud
		_, err = clientconfig.GetCloudFromYAML(opts)

		// GetCloudFromYAML sequentially checks several places for the presence
		// of clouds.yaml file. Unfortunately, it does not provide an interface
		// to understand from which specific location this file was read. Thus, we
		// have to look through the presence of files in the same order and determine
		// the location by ourselves.
		configPath := determineOpenStackConfigLocation()
		if configPath == "" {
			logrus.Warning("The location of OpenStack config file cannot be determined.")
		} else {
			logrus.Infof("Read OpenStack config from %v.", configPath)
		}
	case baremetal.Name, libvirt.Name, none.Name, vsphere.Name:
		// no creds to check
	case azure.Name:
		_, err = azureconfig.GetSession()
		if err != nil {
			return errors.Wrap(err, "creating Azure session")
		}
	default:
		err = fmt.Errorf("unknown platform type %q", platform)
	}

	return err
}

// Name returns the human-friendly name of the asset.
func (a *PlatformCredsCheck) Name() string {
	return "Platform Credentials Check"
}

// fileExists reports whether the file exists.
func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		return true, nil
	}
	if err != nil && os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// determineOpenStackConfigLocation sequentially searches through the
// possible places where clouds.yaml can be located and returns the first
// location found.
func determineOpenStackConfigLocation() string {
	// Currently the search order is:
	// 1. OS_CLIENT_CONFIG_FILE env variable
	// 2. Current directory.
	// 3. unix-specific user_config_dir (~/.config/openstack/clouds.yaml)
	// 4. unix-specific site_config_dir (/etc/openstack/clouds.yaml)
	//
	// For more information:
	// https://github.com/gophercloud/utils/blob/master/openstack/clientconfig/utils.go#L95-L98

	yamlFile := "clouds.yaml"

	// OS_CLIENT_CONFIG_FILE
	if path := os.Getenv("OS_CLIENT_CONFIG_FILE"); path != "" {
		if ok, _ := fileExists(path); ok {
			return path
		}
	}

	// Current directory.
	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, yamlFile)
	if ok, _ := fileExists(path); ok {
		return path
	}

	// unix user config directory: ~/.config/openstack.
	if currentUser, err := user.Current(); err == nil {
		homeDir := currentUser.HomeDir
		if homeDir != "" {
			path := filepath.Join(homeDir, ".config/openstack/"+yamlFile)
			if ok, _ := fileExists(path); ok {
				return path
			}
		}
	}

	// unix-specific site config directory: /etc/openstack.
	path = "/etc/openstack/" + yamlFile
	if ok, _ := fileExists(path); ok {
		return path
	}

	return ""
}
