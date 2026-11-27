package bootstrap

import (
	"context"
	"io"
	"path"
	"strings"
	"unicode"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

const (
	singleNodeBootstrapInPlaceIgnFilename = "bootstrap-in-place-for-live-iso.ign"

	bipInstallToDiskTemplateURI = "bootstrap/bootstrap-in-place/files/usr/local/bin/install-to-disk.sh.template"

	// installationDiskOnePathDetail is the generic rejection for extra arguments
	// or other malformed nonempty installationDisk values that are not the
	// exact legacy "--copy-network" form.
	installationDiskOnePathDetail = "must contain one absolute disk path"

	installationDiskCopyNetworkMigrationDetail = "Remove --copy-network from installationDisk; this installation flow already copies network configuration by default."

	installationDiskCopyNetworkMissingDashNDetail = "Remove --copy-network from installationDisk; the install-to-disk template does not currently pass -n to enable network configuration copying by default."
)

var (
	bootstrapInPlaceEnabledServices = []string{
		"install-to-disk.service",
	}
)

// SingleNodeBootstrapInPlace is an asset that generates the ignition config for single node OpenShift.
type SingleNodeBootstrapInPlace struct {
	Common
}

var _ asset.Asset = (*SingleNodeBootstrapInPlace)(nil)

// Name returns the human-friendly name of the asset.
func (a *SingleNodeBootstrapInPlace) Name() string {
	return "Single node bootstrap In Place Ignition Config"
}

// Generate generates the ignition config for the Bootstrap asset.
func (a *SingleNodeBootstrapInPlace) Generate(_ context.Context, dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)
	if err := verifyBootstrapInPlace(installConfig.Config); err != nil {
		return err
	}
	templateData := a.getTemplateData(dependencies, true)
	if err := a.generateConfig(dependencies, templateData); err != nil {
		return err
	}
	if err := AddStorageFiles(a.Config, "/", "bootstrap/bootstrap-in-place/files", templateData); err != nil {
		return err
	}
	if err := AddSystemdUnits(a.Config, "bootstrap/bootstrap-in-place/systemd/units", templateData, bootstrapInPlaceEnabledServices); err != nil {
		return err
	}
	if err := a.Common.generateFile(singleNodeBootstrapInPlaceIgnFilename); err != nil {
		return err
	}
	return nil
}

// Load returns the bootstrap-in-place ignition from disk.
func (a *SingleNodeBootstrapInPlace) Load(f asset.FileFetcher) (found bool, err error) {
	return a.load(f, singleNodeBootstrapInPlaceIgnFilename)
}

// verifyBootstrapInPlace validate the number of control plane replica is one and that installation disk is set
func verifyBootstrapInPlace(installConfig *types.InstallConfig) error {
	errorList := field.ErrorList{}
	if installConfig.ControlPlane.Replicas == nil {
		errorList = append(errorList, field.Invalid(field.NewPath("controlPlane", "replicas"), installConfig.ControlPlane.Replicas,
			"bootstrap in place requires ControlPlane.Replicas configuration"))
	}
	if *installConfig.ControlPlane.Replicas != 1 {
		errorList = append(errorList, field.Invalid(field.NewPath("controlPlane", "replicas"), installConfig.ControlPlane.Replicas,
			"bootstrap in place requires a single ControlPlane replica"))
	}
	if installConfig.BootstrapInPlace == nil {
		errorList = append(errorList, field.Required(field.NewPath("bootstrapInPlace"), "bootstrapInPlace is required when creating a single node bootstrap-in-place ignition"))
	} else if installConfig.BootstrapInPlace.InstallationDisk == "" {
		errorList = append(errorList, field.Required(field.NewPath("bootstrapInPlace", "installationDisk"),
			"installationDisk must be set the target disk drive for the installation"))
	} else if err := validateBootstrapInPlaceInstallationDisk(
		installConfig.BootstrapInPlace.InstallationDisk,
		field.NewPath("bootstrapInPlace", "installationDisk"),
	); err != nil {
		errorList = append(errorList, err)
	}
	return errorList.ToAggregate()
}

// validateBootstrapInPlaceInstallationDisk enforces that disk is a single absolute
// path suitable as one coreos-installer argument. This is an input-format policy:
// shell-safety for accepted names is provided by quoting in install-to-disk.sh.
// Extra CLI arguments are rejected. The legacy form "<path> --copy-network" is
// recognized only to select a migration hint; the flag is never stripped or forwarded.
func validateBootstrapInPlaceInstallationDisk(disk string, fldPath *field.Path) *field.Error {
	fields := strings.Fields(disk)
	switch {
	case len(fields) == 0:
		return field.Invalid(fldPath, disk, installationDiskOnePathDetail)
	case len(fields) == 1:
		// Reject leading/trailing whitespace or other padding around a single token.
		if fields[0] != disk {
			return field.Invalid(fldPath, disk, installationDiskOnePathDetail)
		}
		return validateSingleBootstrapInPlaceDiskPath(disk, fldPath)
	case len(fields) == 2 && fields[1] == "--copy-network" && isValidSingleBootstrapInPlaceDiskPath(fields[0]):
		// Exact legacy token only (not a substring / lookalike). Still reject.
		return field.Invalid(fldPath, disk, legacyCopyNetworkRejectionDetail())
	default:
		return field.Invalid(fldPath, disk, installationDiskOnePathDetail)
	}
}

func validateSingleBootstrapInPlaceDiskPath(disk string, fldPath *field.Path) *field.Error {
	// Reject non-absolute paths, root, ".." components, and paths that clean to
	// root (for example "/." or "/.."). Keep the original disk string for
	// rendering; do not normalize accepted values.
	if !path.IsAbs(disk) || disk == "/" || path.Clean(disk) == "/" {
		return field.Invalid(fldPath, disk, "must be an absolute disk path")
	}
	for _, part := range strings.Split(disk, "/") {
		if part == ".." {
			return field.Invalid(fldPath, disk, "must be an absolute disk path")
		}
	}
	for _, r := range disk {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return field.Invalid(fldPath, disk, "must be a single path without whitespace or control characters")
		}
	}
	// Reject shell expression syntax that is not needed for documented udev /
	// by-id / by-path names. Backslash is allowed so literal \xNN sequences pass;
	// quoting in the generated script preserves them as one argv.
	if strings.ContainsAny(disk, ";$`|&<>()\"'") {
		return field.Invalid(fldPath, disk, "must not contain shell expression syntax")
	}
	return nil
}

func isValidSingleBootstrapInPlaceDiskPath(disk string) bool {
	return validateSingleBootstrapInPlaceDiskPath(disk, field.NewPath("bootstrapInPlace", "installationDisk")) == nil
}

// legacyCopyNetworkRejectionDetail chooses the migration message only after
// confirming the production install-to-disk template still passes -n. If -n is
// absent, report that difference instead of claiming default network copying.
func legacyCopyNetworkRejectionDetail() string {
	if installToDiskTemplateEnablesNetworkCopy() {
		return installationDiskCopyNetworkMigrationDetail
	}
	return installationDiskCopyNetworkMissingDashNDetail
}

func installToDiskTemplateEnablesNetworkCopy() bool {
	content, err := readBIPInstallToDiskTemplate()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "coreos-installer install") && strings.Contains(trimmed, " -n ") {
			return true
		}
	}
	return false
}

func readBIPInstallToDiskTemplate() ([]byte, error) {
	f, err := data.Assets.Open(bipInstallToDiskTemplateURI)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
