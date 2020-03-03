package libvirt

import (
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/plugins/libvirt/loader"

	"github.com/sirupsen/logrus"
)

func init() {
	providers.Registry["libvirt"] = func(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
		lvp, err := loader.LoadPlugin()
		if err != nil {
			return nil, err
		}
		return lvp.NewUninstaller(logger, metadata)
	}
}
