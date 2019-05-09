package installconfig

import (
	utilrand "k8s.io/apimachinery/pkg/util/rand"

	"github.com/openshift/installer/pkg/asset"
)

// IgnitionAuth gates access to the Machine Config Server
type IgnitionAuth struct {
	Master string
	Worker string
}

var _ asset.Asset = (*IgnitionAuth)(nil)

// Dependencies returns nothing.
func (a *IgnitionAuth) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates a new IgnitionAuth
func (a *IgnitionAuth) Generate(dep asset.Parents) error {
	a.Master = utilrand.String(64)
	a.Worker = utilrand.String(64)
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *IgnitionAuth) Name() string {
	return "Ignition Auth"
}
