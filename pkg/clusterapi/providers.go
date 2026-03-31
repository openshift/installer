package clusterapi

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	dictFile = "dict"
)

var (
	// ClusterAPI is the core provider for cluster-api.
	ClusterAPI = Provider{
		Name:    "cluster-api",
		Sources: sets.New("cluster-api"),
	}

	// EnvTest is the provider for the local control plane.
	EnvTest = Provider{
		Name:    "envtest",
		Sources: sets.New("kube-apiserver", "etcd"),
	}

	// AWS is the provider for creating resources in AWS.
	AWS = infrastructureProvider("aws")
	// Azure is the provider for creating resources in Azure.
	Azure = infrastructureProvider("azure")
	// AzureStack is the provider for creating resources in AzureStack.
	// The AzureStack provider is maintained in an OpenShift fork of CAPZ.
	AzureStack = infrastructureProvider("azurestack")
	// AzureASO is a companion component to Azure that is used to create resources declaratively.
	AzureASO = infrastructureProvider("azureaso")
	// GCP is the provider for creating resources in GCP.
	GCP = infrastructureProvider("gcp")
	// IBMCloud is the provider for creating resources in IBM Cloud and powervs.
	IBMCloud = infrastructureProvider("ibmcloud")
	// Nutanix is the provider for creating resources in Nutanix.
	Nutanix = infrastructureProvider("nutanix")
	// OpenStack is the provider for creating resources in OpenStack.
	OpenStack = infrastructureProvider("openstack")
	// OpenStackORC is a companion component to OpenStack that is used to create resources declaratively.
	OpenStackORC = infrastructureProvider("openstackorc")
	// VSphere is the provider for creating resources in vSphere.
	VSphere = infrastructureProvider("vsphere")
)

// Provider is a Cluster API provider.
type Provider struct {
	// Name of the provider.
	Name string
	// Sources of the provider.
	Sources sets.Set[string]
}

// infrastructureProvider configures a infrastructureProvider built locally.
func infrastructureProvider(name string) Provider {
	return Provider{
		Name: name,
		Sources: sets.New(
			fmt.Sprintf("cluster-api-provider-%s", name),
		),
	}
}

// Mirror is the embedded data for the providers.
//
//go:embed mirror/*
var Mirror embed.FS

// Extract extracts the provider from the embedded data into the specified directory.
func (p Provider) Extract(dir string) error {
	// Load the shared dictionary.
	dictData, err := Mirror.ReadFile(path.Join("mirror", dictFile))
	if err != nil {
		return fmt.Errorf("failed to read zstd dictionary from mirror: %w", err)
	}

	// Ensure the directory exists.
	logrus.Debugf("Creating %s directory", dir)
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return fmt.Errorf("could not make directory for the %s provider: %w", p.Name, err)
	}

	// Extract only the files needed for this provider.
	for source := range p.Sources {
		zstFile := source + ".zst"
		f, err := Mirror.Open(path.Join("mirror", zstFile))
		if err != nil {
			return fmt.Errorf("failed to open %s from mirror: %w", zstFile, err)
		}

		destPath := filepath.Join(dir, source)
		logrus.Debugf("Extracting %s", destPath)
		if err := decompressFile(f, destPath, dictData); err != nil {
			f.Close()
			return fmt.Errorf("failed to extract %s: %w", source, err)
		}
		f.Close()
	}
	return nil
}

func decompressFile(src io.Reader, destPath string, dictData []byte) error {
	decoder, err := zstd.NewReader(src, zstd.WithDecoderDicts(dictData))
	if err != nil {
		return fmt.Errorf("failed to create zstd decoder: %w", err)
	}
	defer decoder.Close()

	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, decoder); err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}
	return nil
}

// UnpackClusterAPIBinary unpacks the cluster-api binary from the embedded data so that it can be run to create the
// infrastructure for the cluster.
func UnpackClusterAPIBinary(dir string) error {
	return ClusterAPI.Extract(dir)
}

// UnpackEnvtestBinaries unpacks the envtest binaries from the embedded data so that it can be run to create the
// infrastructure for the cluster.
func UnpackEnvtestBinaries(dir string) error {
	return EnvTest.Extract(dir)
}
