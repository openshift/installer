package providers

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// AliCloud is the provider for creating resources in the alibaba cloud.
	AliCloud = remoteProvider("hashicorp/alicloud")
	// AWS is the provider for creating resources in AWS.
	AWS = remoteProvider("hashicorp/aws")
	// AzurePrivateDNS is an internal provider for creating private DNS zones in Azure.
	AzurePrivateDNS = localProvider("azureprivatedns")
	// AzureRM is the provider for creating resources in the Azure clouds.
	AzureRM = remoteProvider("hashicorp/azurerm")
	// AzureStack is the provider for creating resources in Azure Stack.
	AzureStack = localProvider("azurestack")
	// Google is the provider for creating resources in GCP.
	Google = remoteProvider("hashicorp/google")
	// IBM is the provider for creating resources in IBM Cloud.
	IBM = localProvider("ibm")
	// Ignition is the provider for creating ignition config files.
	Ignition = localProvider("ignition")
	// Ironic is the provider for provisioning baremetal machines managed by Ironic.
	Ironic = localProvider("ironic")
	// Libvirt is the provider for provisioning VMs on a libvirt host.
	Libvirt = remoteProvider("dmacvicar/libvirt")
	// Local is the provider for creating local files.
	Local = remoteProvider("hashicorp/local")
	// OpenStack is the provider for creating resources in OpenStack.
	OpenStack = remoteProvider("terraform-provider-openstack/openstack")
	// OVirt is the provider for creating resources in oVirt.
	OVirt = localProvider("ovirt")
	// Random is the provider for generating randomness.
	Random = remoteProvider("hashicorp/random")
	// VSphere is the provider for creating resource in vSphere.
	VSphere = localProvider("vsphere")
	// VSpherePrivate is an internal provider augmenting the VSphere provider by adding functionality.
	VSpherePrivate = localProvider("vsphereprivate")
)

// Provider is a terraform provider.
type Provider struct {
	// Name of the provider.
	Name string
	// Source of the provider.
	Source string
}

// remoteProvider configures a provider downloaded from a remote registry.
func remoteProvider(source string) Provider {
	sourceParts := strings.Split(source, "/")
	switch len(sourceParts) {
	case 1:
		source = "hashicorp/" + source
		fallthrough
	case 2:
		source = "registry.terraform.io/" + source
	}
	name := sourceParts[len(sourceParts)-1]
	return Provider{
		Name:   name,
		Source: source,
	}
}

// localProvider configures a provider built locally.
func localProvider(name string) Provider {
	return Provider{
		Name:   name,
		Source: fmt.Sprintf("openshift/local/%s", name),
	}
}

//go:embed mirror/*
var mirror embed.FS

// Extract extracts the provider from the embedded data into the specified directory.
func (p Provider) Extract(dir string) error {
	providerDir := filepath.Join(strings.Split(p.Source, "/")...)
	destProviderDir := filepath.Join(dir, providerDir)
	destDir := destProviderDir
	srcDir := filepath.Join("mirror", providerDir)
	logrus.Debugf("creating %s directory", destDir)
	if err := os.MkdirAll(destDir, 0777); err != nil {
		return errors.Wrapf(err, "could not make directory for the %s provider", p.Name)
	}
	if err := unpack(srcDir, destDir); err != nil {
		return errors.Wrapf(err, "could not unpack the directory for the %s provider", p.Name)
	}
	return nil
}

func unpack(srcDir, destDir string) error {
	entries, err := mirror.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			childSrcDir := filepath.Join(srcDir, entry.Name())
			childDestDir := filepath.Join(destDir, entry.Name())
			logrus.Debugf("creating %s directory", childDestDir)
			if err := os.Mkdir(childDestDir, 0777); err != nil {
				return err
			}
			if err := unpack(childSrcDir, childDestDir); err != nil {
				return err
			}
			continue
		}
		logrus.Debugf("creating %s file", filepath.Join(destDir, entry.Name()))
		if err := unpackFile(filepath.Join(srcDir, entry.Name()), filepath.Join(destDir, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func unpackFile(srcPath, destPath string) error {
	srcFile, err := mirror.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer destFile.Close()
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	return nil
}

// UnpackTerraformBinary unpacks the terraform binary from the embedded data so that it can be run to create the
// infrastructure for the cluster.
func UnpackTerraformBinary(dir string) error {
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	return unpack("mirror/terraform", dir)
}
