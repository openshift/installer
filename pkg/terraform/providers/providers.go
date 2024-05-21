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
	// AWS is the provider for creating resources in AWS.
	AWS = provider("aws")
	// AzureRM is the provider for creating resources in the Azure clouds.
	AzureRM = provider("azurerm")
	// AzureStack is the provider for creating resources in Azure Stack.
	AzureStack = provider("azurestack")
	// Google is the provider for creating resources in GCP.
	Google = provider("google")
	// IBM is the provider for creating resources in IBM Cloud.
	IBM = provider("ibm")
	// Ignition is the provider for creating ignition config files.
	Ignition = provider("ignition")
	// Libvirt is the provider for provisioning VMs on a libvirt host.
	Libvirt = provider("libvirt")
	// Local is the provider for creating local files.
	Local = provider("local")
	// Nutanix is the provider for creating resources in Nutanix.
	Nutanix = provider("nutanix")
	// OpenStack is the provider for creating resources in OpenStack.
	OpenStack = provider("openstack")
	// OVirt is the provider for creating resources in oVirt.
	OVirt = provider("ovirt")
	// Time is the provider for adding create and sleep requirements for resources.
	Time = provider("time")
)

// Provider is a terraform provider.
type Provider struct {
	// Name of the provider.
	Name string
	// Source of the provider.
	Source string
}

// provider configures a provider built locally.
func provider(name string) Provider {
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
