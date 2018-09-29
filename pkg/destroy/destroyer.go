package destroy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	atd "github.com/openshift/hive/contrib/pkg/awstagdeprovision"
	log "github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/metadata"
	"github.com/openshift/installer/pkg/types"
)

// Destroyer allows multiple implementations of destroy
// for different platforms.
type Destroyer interface {
	Run() error
}

// NewDestroyer returns Destroyer based on `metadata.json` in `rootDir`.
func NewDestroyer(level log.Level, rootDir string) (Destroyer, error) {
	raw, err := ioutil.ReadFile(filepath.Join(rootDir, metadata.MetadataFilename))
	if err != nil {
		return nil, err
	}

	var cmetadata types.ClusterMetadata
	if err := json.Unmarshal(raw, &cmetadata); err != nil {
		return nil, err
	}

	var ret Destroyer
	switch {
	case cmetadata.ClusterPlatformMetadata.AWS != nil:
		ret = NewAWSDestroyer(level, &cmetadata)
	case cmetadata.ClusterPlatformMetadata.Libvirt != nil:
		// ret = NewLibvirtDestroyer(level, &cmetadata)
		return nil, fmt.Errorf("libvirt destroyer is not yet supported")
	default:
		return nil, fmt.Errorf("couldn't find Destroyer for %q", metadata.MetadataFilename)
	}
	return ret, nil
}

// // NewLibvirtDestroyer returns libvirt Uninstaller from ClusterMetadata.
// func NewLibvirtDestroyer(level log.Level, metadata *types.ClusterMetadata) *lpd.ClusterUninstaller {
// 	return &lpd.ClusterUninstaller{
// 		LibvirtURI: metadata.ClusterPlatformMetadata.Libvirt.URI,
// 		Filter:     lpd.AlwaysTrueFilter(), //TODO: change to ClusterNamePrefixFilter when all resources are prefixed.
// 		Logger: log.NewEntry(&log.Logger{
// 			Out: os.Stdout,
// 			Formatter: &log.TextFormatter{
// 				FullTimestamp: true,
// 			},
// 			Hooks: make(log.LevelHooks),
// 			Level: level,
// 		}),
// 	}
// }

// NewAWSDestroyer returns aws Uninstaller from ClusterMetadata.
func NewAWSDestroyer(level log.Level, metadata *types.ClusterMetadata) *atd.ClusterUninstaller {
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
	}
}
