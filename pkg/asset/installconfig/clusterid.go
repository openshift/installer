package installconfig

import (
	"fmt"
	"regexp"

	"github.com/pborman/uuid"
	utilrand "k8s.io/apimachinery/pkg/util/rand"

	"github.com/openshift/installer/pkg/asset"
)

const (
	// resource using InfraID usually have suffixes like `[-/_][a-z]{3,4}` eg. `_int`, `-ext` or `-ctlp`
	maxNameLen = 32 - 5
	randomLen  = 5
	maxBaseLen = maxNameLen - randomLen - 1
)

// ClusterID is the unique ID of the cluster, immutable during the cluster's life
type ClusterID struct {
	// UUID is a globally unique identifier.
	UUID string

	// InfraID is an identifier for the cluster that is more human friendly.
	// This does not have
	InfraID string
}

var _ asset.Asset = (*ClusterID)(nil)

// Dependencies returns install-config.
func (a *ClusterID) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate generates a new ClusterID
func (a *ClusterID) Generate(dep asset.Parents) error {
	ica := &InstallConfig{}
	dep.Get(ica)

	// add random chars to the end to randomize
	a.InfraID = generateInfraID(ica.Config.ObjectMeta.Name)
	a.UUID = uuid.New()
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *ClusterID) Name() string {
	return "Cluster ID"
}

// generateInfraID take base and returns a ID that
// - is of length maxNameLen
// - only contains `alphanum` or `-`
func generateInfraID(base string) string {
	// truncate to maxBaseLen
	if len(base) > maxBaseLen {
		base = base[:maxBaseLen]
	}

	// replace all characters that are not `alphanum` or `-` with `-`
	re := regexp.MustCompile("[^A-Za-z0-9-]")
	base = re.ReplaceAllString(base, "-")

	// add random chars to the end to randomize
	return fmt.Sprintf("%s-%s", base, utilrand.String(randomLen))
}
