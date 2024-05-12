package utils

import (
	"path/filepath"

	"github.com/openshift-kni/lifecycle-agent/internal/common"
)

const (
	IBUWorkspacePath string = common.LCAConfigDir + "/workspace"
	// IBUName defines the valid name of the CR for the controller to reconcile
	IBUName     string = "upgrade"
	IBUFilePath string = common.LCAConfigDir + "/ibu.json"

	ManualCleanupAnnotation    string = "lca.openshift.io/manual-cleanup-done"
	TriggerReconcileAnnotation string = "lca.openshift.io/trigger-reconcile"

	// SeedGenName defines the valid name of the CR for the controller to reconcile
	SeedGenName          string = "seedimage"
	SeedGenSecretName    string = "seedgen"
	SeedgenWorkspacePath string = common.LCAConfigDir + "/ibu-seedgen-orch" // The LCAConfigDir folder is excluded from the var.tgz backup in seed image creation
)

var (
	SeedGenStoredCR       = filepath.Join(SeedgenWorkspacePath, "seedgen-cr.json")
	SeedGenStoredSecretCR = filepath.Join(SeedgenWorkspacePath, "seedgen-secret.json")

	StoredPullSecret = filepath.Join(SeedgenWorkspacePath, "pull-secret.json")
)
