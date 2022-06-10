package bootstrap

import (
	"encoding/json"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests"
)

var (
	_ asset.WritableAsset = (*CVOIgnore)(nil)
)

const (
	cvoOverridesFilename      = "manifests/cvo-overrides.yaml"
	originalOverridesFilename = "original_cvo_overrides.patch"
)

// CVOIgnore adds bootstrap files needed to inform CVO to ignore resources for which the installer is providing manifests.
type CVOIgnore struct {
	asset.DefaultFileListWriter
}

// Name returns a human friendly name for the operator
func (a *CVOIgnore) Name() string {
	return "CVO Ignore"
}

// Dependencies returns all of the dependencies directly needed by the CVOIgnore asset
func (a *CVOIgnore) Dependencies() []asset.Asset {
	return []asset.Asset{
		&manifests.Manifests{},
		&manifests.Openshift{},
	}
}

// Generate generates the respective operator config.yml files
func (a *CVOIgnore) Generate(dependencies asset.Parents) error {
	operators := &manifests.Manifests{}
	openshiftManifests := &manifests.Openshift{}
	dependencies.Get(operators, openshiftManifests)

	var clusterVersion *unstructured.Unstructured
	var ignoredResources []interface{}
	var files []*asset.File
	files = append(files, operators.FileList...)
	files = append(files, openshiftManifests.FileList...)
	for _, file := range files {
		u := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(file.Data, u); err != nil {
			return errors.Wrapf(err, "could not unmarshal %q", file.Filename)
		}
		if file.Filename == cvoOverridesFilename {
			clusterVersion = u
			continue
		}
		ignoredResources = append(ignoredResources,
			configv1.ComponentOverride{
				Kind:      u.GetKind(),
				Group:     u.GetObjectKind().GroupVersionKind().Group,
				Namespace: u.GetNamespace(),
				Name:      u.GetName(),
				Unmanaged: true,
			})
	}

	specAsInterface, ok := clusterVersion.Object["spec"]
	if !ok {
		specAsInterface = map[string]interface{}{}
		clusterVersion.Object["spec"] = specAsInterface
	}
	spec, ok := specAsInterface.(map[string]interface{})
	if !ok {
		return errors.Errorf("unexpected type (%T) for .spec in clusterversion", specAsInterface)
	}
	originalOverridesAsInterface := spec["overrides"]
	originalOverrides, ok := originalOverridesAsInterface.([]interface{})
	if !ok && originalOverridesAsInterface != nil {
		return errors.Errorf("unexpected type (%T) for .spec.overrides in clusterversion", originalOverridesAsInterface)
	}
	originalOverridesPatch := map[string]interface{}{
		"spec": map[string]interface{}{
			"overrides": originalOverrides,
		},
	}
	spec["overrides"] = append(ignoredResources, originalOverrides...)

	cvData, err := yaml.Marshal(clusterVersion)
	if err != nil {
		return errors.Wrap(err, "error marshalling clusterversion")
	}
	a.FileList = append(a.FileList, &asset.File{
		Filename: cvoOverridesFilename,
		Data:     cvData,
	})

	origOverrideData, err := json.Marshal(originalOverridesPatch)
	if err != nil {
		return errors.Wrap(err, "error marshalling original overrides")
	}
	a.FileList = append(a.FileList, &asset.File{
		Filename: originalOverridesFilename,
		Data:     origOverrideData,
	})

	return nil
}

// Load does nothing as the file should not be loaded from disk.
func (a *CVOIgnore) Load(f asset.FileFetcher) (bool, error) {
	return false, nil
}
