package common

import (
	"context"

	"github.com/google/uuid"

	"github.com/openshift/installer/pkg/asset"
)

// InfraEnvID is an asset that generates infraEnvID.
type InfraEnvID struct {
	ID string
}

var _ asset.Asset = (*InfraEnvID)(nil)

// Dependencies returns the assets on which the InfraEnv asset depends.
func (a *InfraEnvID) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the InfraEnvID for agent installer.
func (a *InfraEnvID) Generate(_ context.Context, dependencies asset.Parents) error {
	a.ID = uuid.New().String()
	return nil
}

// Name returns the human-friendly name of the asset.
func (a *InfraEnvID) Name() string {
	return "Agent Installer InfraEnv ID"
}
