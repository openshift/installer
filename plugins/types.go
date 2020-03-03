package plugins

import (
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"

	"github.com/sirupsen/logrus"
)

type Plugin interface {
	Init()
	NewUninstaller(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error)
}
