package destroy

import (
	session "github.com/openshift/installer/pkg/asset/installconfig/aws"
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

	awsSession, err := session.GetSession()
	if err != nil {
		return nil, err
	}

	return &aws.ClusterUninstaller{
		Filters:   filters,
		Region:    metadata.ClusterPlatformMetadata.AWS.Region,
		Logger:    logger,
		ClusterID: metadata.InfraID,
		Session:   awsSession,
	}, nil
}

func init() {
	Registry["aws"] = NewAWS
}
