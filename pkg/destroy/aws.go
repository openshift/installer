package destroy

import (
	atd "github.com/openshift/hive/contrib/pkg/awstagdeprovision"
	"github.com/openshift/installer/pkg/types"
	"github.com/sirupsen/logrus"
)

// NewAWS returns an AWS destroyer from ClusterMetadata.
func NewAWS(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (Destroyer, error) {
	filters := make([]atd.AWSFilter, 0, len(metadata.ClusterPlatformMetadata.AWS.Identifier))
	for _, filter := range metadata.ClusterPlatformMetadata.AWS.Identifier {
		filters = append(filters, filter)
	}

	return &atd.ClusterUninstaller{
		Filters:     filters,
		Region:      metadata.ClusterPlatformMetadata.AWS.Region,
		ClusterName: metadata.ClusterName,
		Logger:      logger,
	}, nil
}

func init() {
	Registry["aws"] = NewAWS
}
