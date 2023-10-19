package providers

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	zipFile = "cluster-api.zip"
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
	// AzureASO is a companion component to Azure that is used to create resources declaratively.
	AzureASO = infrastructureProvider("azureaso")
	// IBMCloud is the provider for creating resources in IBM Cloud and powervs.
	IBMCloud = infrastructureProvider("ibmcloud")
	// GCP is the provider for creating resources in GCP.
	GCP = infrastructureProvider("gcp")
	// Nutanix is the provider for creating resources in Nutanix.
	Nutanix = infrastructureProvider("nutanix")
	// vSphere is the provider for creating resources in vSphere.
	VSphere = infrastructureProvider("vsphere")
)

// Provider is a terraform provider.
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
	zipFile, err := Mirror.Open(filepath.Join("mirror", zipFile))
	if err != nil {
		return errors.Wrap(err, "failed to open cluster api zip from mirror")
	}
	defer zipFile.Close()
	stat, err := zipFile.Stat()
	if err != nil {
		return errors.Wrap(err, "failed to stat cluster api zip")
	}
	zipReaderAt, ok := zipFile.(io.ReaderAt)
	if !ok {
		return errors.New("zip file does not support seeking")
	}

	// Open a zip archive for reading.
	r, err := zip.NewReader(zipReaderAt, stat.Size())
	if err != nil {
		return errors.Wrap(err, "failed to open cluster api zip")
	}

	// Ensure the directory exists.
	logrus.Debugf("creating %s directory", dir)
	if err := os.MkdirAll(dir, 0o777); err != nil {
		return errors.Wrapf(err, "could not make directory for the %s provider", p.Name)
	}

	// Extract the files.
	for _, f := range r.File {
		if !p.Sources.Has(f.Name) {
			continue
		}
		logrus.Debugf("extracting %s file", filepath.Join(dir, f.Name))
		if err := unpackFile(f, filepath.Join(dir, f.Name)); err != nil {
			return errors.Wrapf(err, "failed to extract %s file", f.Name)
		}
	}
	return nil
}

func unpackFile(f *zip.File, destPath string) error {
	src, err := f.Open()
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", f.Name)
	}
	defer src.Close()
	destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o777)
	if err != nil {
		return err
	}
	defer destFile.Close()
	if _, err := io.Copy(destFile, src); err != nil {
		return err
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
