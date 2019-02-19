package installconfig

import (
	"fmt"

	utilrand "k8s.io/apimachinery/pkg/util/rand"

	"github.com/openshift/installer/pkg/asset"
)

const (
	// AWS load balancers have a maximum name length of 32. The load balancers
	// have suffixes of "-int" and "-ext", which are 4 characters.
	maxNameLen = 32 - 4
	randomLen  = 5
	maxBaseLen = maxNameLen - randomLen - 1
)

// UniqueClusterName is the unique name of the cluster. This combines the name of
// the cluster supplied by the user with random characters.
type UniqueClusterName struct {
	ClusterName string
}

var _ asset.Asset = (*UniqueClusterName)(nil)

// Dependencies returns no dependencies.
func (a *UniqueClusterName) Dependencies() []asset.Asset {
	return []asset.Asset{
		&InstallConfig{},
	}
}

// Generate generates a random, unique cluster name
func (a *UniqueClusterName) Generate(parents asset.Parents) error {
	ic := &InstallConfig{}
	parents.Get(ic)

	a.ClusterName = generateName(ic.Config.ObjectMeta.Name)

	return nil
}

// Name returns the human-friendly name of the asset.
func (a *UniqueClusterName) Name() string {
	return "Unique Cluster Name"
}

func generateName(base string) string {
	if len(base) > maxBaseLen {
		base = base[:maxBaseLen]
	}
	return fmt.Sprintf("%s-%s", base, utilrand.String(randomLen))
}
