package quota

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	if err := ioutil.WriteFile(path, raw, 0777); err != nil {
		return errors.Wrap(err, "failed to write quota")
	}
	return nil
}
