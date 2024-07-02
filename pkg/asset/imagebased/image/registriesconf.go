package image

import (
	"fmt"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/pelletier/go-toml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/types"
)

// RegistriesConf generates the registries.conf data.
type RegistriesConf struct {
	Data   []byte
	Config *sysregistriesv2.V2RegistriesConf
}

// Name returns a human friendly name for the asset.
func (*RegistriesConf) Name() string {
	return "Mirror Registries Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*RegistriesConf) Dependencies() []asset.Asset {
	return []asset.Asset{
		&ImageBasedInstallationConfig{},
	}
}

// Generate generates the registries.conf data from the
// imagebased-installation-config.yaml.
func (i *RegistriesConf) Generate(dependencies asset.Parents) error {
	ibiConfig := &ImageBasedInstallationConfig{}
	dependencies.Get(ibiConfig)

	var imageDigestSources []types.ImageDigestSource
	var deprecatedImageContentSources []types.ImageContentSource

	imageDigestSources = ibiConfig.Config.ImageDigestSources
	deprecatedImageContentSources = ibiConfig.Config.DeprecatedImageContentSources

	if len(imageDigestSources) == 0 && len(deprecatedImageContentSources) == 0 {
		return nil
	}

	err := i.generateRegistriesConf(imageDigestSources, deprecatedImageContentSources)
	if err != nil {
		return err
	}

	i.Data, err = toml.Marshal(i.Config)
	if err != nil {
		return err
	}

	return nil
}

func (i *RegistriesConf) generateRegistriesConf(imageDigestSources []types.ImageDigestSource, deprecatedImageContentSources []types.ImageContentSource) error {
	if len(deprecatedImageContentSources) != 0 && len(imageDigestSources) != 0 {
		return fmt.Errorf("invalid imagebased-installation-config.yaml, cannot set imageContentSources and imageDigestSources at the same time")
	}

	digestMirrorSources := []types.ImageDigestSource{}
	if len(deprecatedImageContentSources) > 0 {
		digestMirrorSources = bootstrap.ContentSourceToDigestMirror(deprecatedImageContentSources)
	} else if len(imageDigestSources) > 0 {
		digestMirrorSources = append(digestMirrorSources, imageDigestSources...)
	}

	registries := &sysregistriesv2.V2RegistriesConf{
		Registries: []sysregistriesv2.Registry{},
	}
	for _, group := range bootstrap.MergedMirrorSets(digestMirrorSources) {
		if len(group.Mirrors) == 0 {
			continue
		}

		registry := sysregistriesv2.Registry{}
		registry.Endpoint.Location = group.Source
		registry.MirrorByDigestOnly = true
		for _, mirror := range group.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
		}
		registries.Registries = append(registries.Registries, registry)
	}
	i.Config = registries

	return nil
}
