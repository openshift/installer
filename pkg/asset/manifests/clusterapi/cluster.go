package clusterapi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/manifests/aws"
	"github.com/openshift/installer/pkg/asset/manifests/azure"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/manifests/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/openstack"
	"github.com/openshift/installer/pkg/asset/manifests/vsphere"
	"github.com/openshift/installer/pkg/asset/openshiftinstall"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/clusterapi"
	awstypes "github.com/openshift/installer/pkg/types/aws"
	azuretypes "github.com/openshift/installer/pkg/types/azure"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
	openstacktypes "github.com/openshift/installer/pkg/types/openstack"
	vsphereplatform "github.com/openshift/installer/pkg/types/vsphere"
)

var _ asset.WritableRuntimeAsset = (*Cluster)(nil)

// Cluster generates manifests for target cluster
// creation using CAPI.
type Cluster struct {
	FileList []*asset.RuntimeFile
}

// Name returns a human friendly name for the operator.
func (c *Cluster) Name() string {
	return "Cluster API Manifests"
}

// Dependencies returns all of the dependencies directly needed by the
// ClusterAPI asset.
func (c *Cluster) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&installconfig.ClusterID{},
		&openshiftinstall.Config{},
		&manifests.FeatureGate{},
		new(rhcos.Image),
	}
}

// Generate generates the respective operator config.yml files.
func (c *Cluster) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	clusterID := &installconfig.ClusterID{}
	openshiftInstall := &openshiftinstall.Config{}
	featureGate := &manifests.FeatureGate{}
	rhcosImage := new(rhcos.Image)
	dependencies.Get(installConfig, clusterID, openshiftInstall, featureGate, rhcosImage)

	// If the feature gate is not enabled, do not generate any manifests.
	if !capiutils.IsEnabled(installConfig) {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(capiutils.ManifestDir), 0755); err != nil {
		return err
	}

	c.FileList = []*asset.RuntimeFile{}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: capiutils.Namespace,
		},
	}
	c.FileList = append(c.FileList, &asset.RuntimeFile{Object: namespace, File: asset.File{Filename: "000_capi-namespace.yaml"}})

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: clusterv1.ClusterSpec{},
	}
	c.FileList = append(c.FileList, &asset.RuntimeFile{Object: cluster, File: asset.File{Filename: "01_capi-cluster.yaml"}})

	var out *capiutils.GenerateClusterAssetsOutput
	switch platform := installConfig.Config.Platform.Name(); platform {
	case awstypes.Name:
		var err error
		out, err = aws.GenerateClusterAssets(installConfig, clusterID)
		if err != nil {
			return errors.Wrap(err, "failed to generate AWS manifests")
		}
	case azuretypes.Name:
		var err error
		out, err = azure.GenerateClusterAssets(installConfig, clusterID)
		if err != nil {
			return errors.Wrap(err, "failed to generate Azure manifests")
		}
	case gcptypes.Name:
		var err error
		out, err = gcp.GenerateClusterAssets(installConfig, clusterID)
		if err != nil {
			return fmt.Errorf("failed to generate GCP manifests: %w", err)
		}
	case vsphereplatform.Name:
		var err error
		out, err = vsphere.GenerateClusterAssets(installConfig, clusterID)
		if err != nil {
			return fmt.Errorf("failed to generate vSphere manifests %w", err)
		}
	case openstacktypes.Name:
		var err error
		out, err = openstack.GenerateClusterAssets(installConfig, clusterID)
		if err != nil {
			return errors.Wrap(err, "failed to generate OpenStack manifests")
		}
	default:
		return fmt.Errorf("unsupported platform %q", platform)
	}

	// Set the infrastructure reference in the Cluster object.
	cluster.Spec.InfrastructureRef = out.InfrastructureRef
	if cluster.Spec.InfrastructureRef == nil {
		return fmt.Errorf("failed to generate manifests: cluster.Spec.InfrastructureRef was never set")
	}

	// Append the infrastructure manifests.
	c.FileList = append(c.FileList, out.Manifests...)

	// Create the infrastructure manifests.
	for _, m := range c.FileList {
		objData, err := yaml.Marshal(m.Object)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal infrastructure manifest %s", m.Filename)
		}
		m.Data = objData

		// If the filename is already a path, do not append the manifest dir.
		if filepath.Dir(m.Filename) == capiutils.ManifestDir {
			continue
		}
		m.Filename = filepath.Join(capiutils.ManifestDir, m.Filename)
	}

	asset.SortManifestFiles(c.FileList)
	return nil
}

// Files returns the files generated by the asset.
func (c *Cluster) Files() []*asset.File {
	files := []*asset.File{}
	for _, f := range c.FileList {
		files = append(files, &f.File)
	}
	return files
}

// RuntimeFiles returns the files generated by the asset.
func (c *Cluster) RuntimeFiles() []*asset.RuntimeFile {
	return c.FileList
}

// Load returns the openshift asset from disk.
func (c *Cluster) Load(f asset.FileFetcher) (bool, error) {
	yamlFileList, err := f.FetchByPattern(filepath.Join(capiutils.ManifestDir, "*.yaml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yaml files")
	}
	ymlFileList, err := f.FetchByPattern(filepath.Join(capiutils.ManifestDir, "*.yml"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.yml files")
	}
	jsonFileList, err := f.FetchByPattern(filepath.Join(capiutils.ManifestDir, "*.json"))
	if err != nil {
		return false, errors.Wrap(err, "failed to load *.json files")
	}
	fileList := append(yamlFileList, ymlFileList...) //nolint:gocritic
	fileList = append(fileList, jsonFileList...)

	for _, file := range fileList {
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(file.Data, u); err != nil {
			return false, errors.Wrap(err, "failed to unmarshal file")
		}
		obj, err := clusterapi.Scheme.New(u.GroupVersionKind())
		if err != nil {
			return false, errors.Wrap(err, "failed to create object")
		}
		if err := clusterapi.Scheme.Convert(u, obj, nil); err != nil {
			return false, errors.Wrap(err, "failed to convert object")
		}
		c.FileList = append(c.FileList, &asset.RuntimeFile{
			File: asset.File{
				Filename: file.Filename,
				Data:     file.Data},
			Object: obj.(client.Object),
		})
	}

	asset.SortManifestFiles(c.FileList)
	return len(c.FileList) > 0, nil
}
