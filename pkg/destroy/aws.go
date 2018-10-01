package destroy

import (
	"os"

	atd "github.com/openshift/hive/contrib/pkg/awstagdeprovision"
	"github.com/openshift/installer/pkg/types"
	log "github.com/sirupsen/logrus"
)

// NewAWS returns an AWS destroyer from ClusterMetadata.
func NewAWS(level log.Level, metadata *types.ClusterMetadata) (Destroyer, error) {
	return &atd.ClusterUninstaller{
		Filters:     metadata.ClusterPlatformMetadata.AWS.Identifier,
		Region:      metadata.ClusterPlatformMetadata.AWS.Region,
		ClusterName: metadata.ClusterName,
		Logger: log.NewEntry(&log.Logger{
			Out: os.Stdout,
			Formatter: &log.TextFormatter{
				FullTimestamp: true,
			},
			Hooks: make(log.LevelHooks),
			Level: level,
		}),
	}, nil
}

func init() {
	Registry["aws"] = NewAWS
}
