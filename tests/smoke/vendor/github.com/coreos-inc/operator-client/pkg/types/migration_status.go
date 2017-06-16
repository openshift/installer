package types

import (
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/unversioned"
)

// MigrationStage represents where in the process of an upgrade
// the migration is ran.
type MigrationStage string

const (
	// MigrationStatusKind is the Kind of the MigrationStatus TPR.
	MigrationStatusKind = "MigrationStatus"
	// MigrationAPIGroup is the API Group for MigrationStatus TPR.
	MigrationAPIGroup = "coreos.com"
	// MigrationTPRVersion is the version for MigrationStatus TPR.
	MigrationTPRVersion = "v1"

	// MigrationStageBefore is Migrations ran before update.
	MigrationStageBefore MigrationStage = "lastBeforeMigrationRan"
	// MigrationStageAfter is Migrations ran after update.
	MigrationStageAfter MigrationStage = "lastAfterMigrationRan"
)

// MigrationStatus represents the 3rd Party API Object.
type MigrationStatus struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`

	Versions MigrationVersions `json:"versions,omitempty"`
}

// MigrationVersions represents the migrations that have already
// been run for a specific version. The format of the map is:
// version -> "before/after" -> N (where N is the last
// migration that was run successfully).
type MigrationVersions map[string]*MigrationVersion

// MigrationVersion represents a migration for a given version.
// Eeach field holds the last migration that was successfully ran.
type MigrationVersion struct {
	LastBeforeMigrationRan *int `json:"lastBeforeMigrationRan,omitempty"`
	LastAfterMigrationRan  *int `json:"lastAfterMigrationRan,omitempty"`
}
