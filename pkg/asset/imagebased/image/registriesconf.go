package image

import (
	"context"

	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/pelletier/go-toml"

	configv1 "github.com/openshift/api/config/v1"
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

// Generate generates the registries.conf data from the image-based-installation-config.
func (i *RegistriesConf) Generate(_ context.Context, dependencies asset.Parents) error {
	ibiConfig := &ImageBasedInstallationConfig{}
	dependencies.Get(ibiConfig)

	if ibiConfig.Config == nil {
		return nil
	}

	imageDigestSources := ibiConfig.Config.ImageDigestSources

	if len(imageDigestSources) == 0 {
		return nil
	}

	err := i.generateRegistriesConf(imageDigestSources)
	if err != nil {
		return err
	}

	i.Data, err = toml.Marshal(i.Config)
	if err != nil {
		return err
	}

	return nil
}

func (i *RegistriesConf) generateRegistriesConf(imageDigestSources []types.ImageDigestSource) error {
	digestMirrorSources := []types.ImageDigestSource{}
	digestMirrorSources = append(digestMirrorSources, imageDigestSources...)

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
		registry.Blocked = group.SourcePolicy == configv1.NeverContactSource
		for _, mirror := range group.Mirrors {
			registry.Mirrors = append(registry.Mirrors, sysregistriesv2.Endpoint{Location: mirror})
		}
		registries.Registries = append(registries.Registries, registry)
	}
	i.Config = registries

	return nil
}
