package configimage

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	utilrand "k8s.io/apimachinery/pkg/util/rand"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

const (
	randomLen = 5
)

// ClusterID is the unique ID of the cluster, immutable during the cluster's life.
type ClusterID struct {
	installconfig.ClusterID
}

var _ asset.Asset = (*ClusterID)(nil)

// Dependencies returns install-config.
func (i *ClusterID) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate generates a new ClusterID.
func (i *ClusterID) Generate(_ context.Context, dep asset.Parents) error {
	ica := &InstallConfig{}
	dep.Get(ica)

	if ica.Config == nil {
		return fmt.Errorf("missing install-config.yaml")
	}

	// resource using InfraID usually have suffixes like `[-/_][a-z]{3,4}` eg. `_int`, `-ext` or `-ctlp`
	// and the maximum length for most resources is approx 32.
	maxLen := 27

	// add random chars to the end to randomize
	i.InfraID = generateInfraID(ica.Config.ObjectMeta.Name, maxLen)
	i.UUID = uuid.New()
	return nil
}

// Name returns the human-friendly name of the asset.
func (i *ClusterID) Name() string {
	return "Cluster ID"
}

// generateInfraID take base and returns a ID that
// - is of length maxLen
// - only contains `alphanum` or `-`.
func generateInfraID(base string, maxLen int) string {
	maxBaseLen := maxLen - (randomLen + 1)

	// replace all characters that are not `alphanum` or `-` with `-`
	re := regexp.MustCompile("[^A-Za-z0-9-]")
	base = re.ReplaceAllString(base, "-")

	// replace all multiple dashes in a sequence with single one.
	re = regexp.MustCompile(`-{2,}`)
	base = re.ReplaceAllString(base, "-")

	// truncate to maxBaseLen
	if len(base) > maxBaseLen {
		logrus.Warnf("Length of cluster name %q is %d which is greater than the max %d allowed. The name will be truncated to %q", base, len(base), maxBaseLen, strings.TrimRight(base[:maxBaseLen], "-"))
		base = base[:maxBaseLen]
	}
	base = strings.TrimRight(base, "-")

	// add random chars to the end to randomize
	return fmt.Sprintf("%s-%s", base, utilrand.String(randomLen))
}
