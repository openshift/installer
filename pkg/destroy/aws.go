package destroy

import (
	"github.com/openshift/installer/pkg/destroy/aws"
	"github.com/openshift/installer/pkg/types"
	"github.com/sirupsen/logrus"
)

// NewAWS returns an AWS destroyer from ClusterMetadata.
func NewAWS(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (Destroyer, error) {
	filters := make([]aws.Filter, 0, len(metadata.ClusterPlatformMetadata.AWS.Identifier))
	for _, filter := range metadata.ClusterPlatformMetadata.AWS.Identifier {
		filters = append(filters, filter)
	}

	return &aws.ClusterUninstaller{
		Filters:   filters,
		Region:    metadata.ClusterPlatformMetadata.AWS.Region,
		Logger:    logger,
		ClusterID: metadata.InfraID,
	}, nil
}

func init() {
	Registry["aws"] = NewAWS
}
