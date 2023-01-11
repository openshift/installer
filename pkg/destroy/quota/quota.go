package quota

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types"
)

const (
	quotaFileName = "quota.json"
)

// WriteQuota writes the cluster quota footprint into the asset directory.
func WriteQuota(dir string, quota *types.ClusterQuota) error {
	path := filepath.Join(dir, quotaFileName)
	logrus.Infof("Writing quota footprint to %s", path)

	raw, err := json.Marshal(quota)
	if err != nil {
		return errors.Wrap(err, "failed to marshal quota")
	}
	if err := os.WriteFile(path, raw, 0o777); err != nil { //nolint:gosec // no sensitive info
		return errors.Wrap(err, "failed to write quota")
	}
	return nil
}
