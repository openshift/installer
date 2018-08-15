package openstack

import (
	"context"

	"github.com/ghodss/yaml"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func cloudConfigRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	clouds, err := clientconfig.LoadCloudsYAML()
	if err != nil {
		return nil, err
	}

	marshalled, err := yaml.Marshal(clouds)
	if err != nil {
		return nil, err
	}

	return installerassets.TemplateRebuilder(
		"files/opt/tectonic/tectonic/openstack/99_cloud-creds-secret.yaml",
		nil,
		map[string]interface{}{
			"Creds": marshalled,
		},
	)(ctx, getByName)
}

func init() {
	installerassets.Rebuilders["files/opt/tectonic/tectonic/openstack/99_cloud-creds-secret.yaml"] = cloudConfigRebuilder
}
