package validation

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsvalidation "github.com/openshift/installer/pkg/types/aws/validation"
	"github.com/openshift/installer/pkg/types/azure"
	azurevalidation "github.com/openshift/installer/pkg/types/azure/validation"
	"github.com/openshift/installer/pkg/types/baremetal"
	baremetalvalidation "github.com/openshift/installer/pkg/types/baremetal/validation"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpvalidation "github.com/openshift/installer/pkg/types/gcp/validation"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	ibmcloudvalidation "github.com/openshift/installer/pkg/types/ibmcloud/validation"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackvalidation "github.com/openshift/installer/pkg/types/openstack/validation"
	"github.com/openshift/installer/pkg/types/ovirt"
	ovirtvalidation "github.com/openshift/installer/pkg/types/ovirt/validation"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/powervs"
	powervsvalidation "github.com/openshift/installer/pkg/types/powervs/validation"
	"github.com/openshift/installer/pkg/types/vsphere"
	vspherevalidation "github.com/openshift/installer/pkg/types/vsphere/validation"
)

var (
	validHyperthreadingModes = map[types.HyperthreadingMode]bool{
		types.HyperthreadingDisabled: true,
		types.HyperthreadingEnabled:  true,
	}

	validHyperthreadingModeValues = func() []string {
		v := make([]string, 0, len(validHyperthreadingModes))
		for m := range validHyperthreadingModes {
			v = append(v, string(m))
		}
		return v
	}()

	validArchitectures = map[types.Architecture]bool{
		types.ArchitectureAMD64:   true,
		types.ArchitectureS390X:   true,
		types.ArchitecturePPC64LE: true,
		types.ArchitectureARM64:   true,
	}

	validArchitectureValues = func() []string {
		v := make([]string, 0, len(validArchitectures))
		for m := range validArchitectures {
			v = append(v, string(m))
		}
		return v
	}()

	validDiskTypeValues = []string{
		string(types.Etcd),
		string(types.Swap),
		string(types.UserDefined),
	}

	// diskSetupFieldNames maps a disk type to the install-config field holding its
	// configuration. The two are not interchangeable: the user-defined type is
	// spelled "user-defined" but its field is userDefined, so building a path out
	// of the type string would point the user at a field that does not exist.
	diskSetupFieldNames = map[types.DiskType]string{
		types.Etcd:        "etcd",
		types.Swap:        "swap",
		types.UserDefined: "userDefined",
	}

	// mountPathRegexp matches the characters allowed in a user-defined disk's
	// mount path. The path is written verbatim into the Where= setting of the
	// generated mount unit, where systemd expands "%" specifiers and a line
	// break would start a new directive, and its escaped form names the unit.
	// Restricting it to the characters common in paths keeps Where= and the
	// unit name in agreement.
	mountPathRegexp = regexp.MustCompile(`^[A-Za-z0-9/._-]+$`)
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(platform *types.Platform, p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Replicas != nil {
		if *p.Replicas < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("replicas"), p.Replicas, "number of replicas must not be negative"))
		}
	} else {
		allErrs = append(allErrs, field.Required(fldPath.Child("replicas"), "replicas is required"))
	}
	if !validHyperthreadingModes[p.Hyperthreading] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("hyperthreading"), p.Hyperthreading, validHyperthreadingModeValues))
	}
	if !validArchitectures[p.Architecture] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("architecture"), p.Architecture, validArchitectureValues))
	}
	if platform.AWS != nil {
		allErrs = append(allErrs, awsvalidation.ValidateMachinePoolArchitecture(p, fldPath.Child("architecture"))...)
	}

	allErrs = append(allErrs, validateDiskSetup(p, fldPath.Child("diskSetup"))...)
	allErrs = append(allErrs, validateMachineManagement(platform, p, fldPath.Child("management"))...)
	allErrs = append(allErrs, validateMachinePoolPlatform(platform, &p.Platform, p, fldPath.Child("platform"))...)
	return allErrs
}

func validateDiskSetup(p *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	foundEtcd := false
	foundSwap := false
	// platformDiskID is what pairs an entry with a platform data disk. Two entries
	// claiming the same ID resolve to the same underlying disk.
	seenDiskIDs := sets.New[string]()
	// The resolved label names the generated MachineConfig. Colliding names
	// silently overwrite each other and one disk is never configured.
	seenLabels := sets.New[string]()
	// The mount path names the generated systemd mount unit. Two entries
	// mounting at the same path ship conflicting units on different devices.
	seenMountPaths := sets.New[string]()
	for i, ds := range p.DiskSetup {
		// Errors about a single entry point at its index, so that the user can
		// tell which of several entries of the same type is at fault.
		entryPath := fldPath.Index(i)
		// outputting the yaml to make recognizing the issue easier for the user
		dsBytes, err := yaml.Marshal(ds)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(fldPath, err))
		}
		dsYaml := string(dsBytes)
		switch ds.Type {
		case types.UserDefined:
			if ds.UserDefined == nil {
				allErrs = append(allErrs, field.Invalid(entryPath.Child("userDefined"), dsYaml, "userDefined configuration must be created"))
				continue
			}
			userDefinedPath := entryPath.Child("userDefined")
			if len(ds.UserDefined.PlatformDiskID) > 12 {
				allErrs = append(allErrs, field.Invalid(userDefinedPath.Child("platformDiskID"), dsYaml, "cannot be longer than 12 characters"))
				continue
			}
			// The partition label is the platform disk ID stripped of every
			// non-alphanumeric character, so an ID carrying none at all resolves
			// to an empty label: a MachineConfig named after nothing and a mount
			// unit pointing at /dev/disk/by-partlabel/ with no disk behind it.
			if types.SanitizeDiskLabel(ds.UserDefined.PlatformDiskID) == "" {
				allErrs = append(allErrs, field.Invalid(userDefinedPath.Child("platformDiskID"), ds.UserDefined.PlatformDiskID,
					"must contain at least one alphanumeric character"))
				continue
			}
			// The mount path must be absolute and normalized: it names the
			// mount point in the generated MachineConfig, and both the systemd
			// unit name and the comparisons against installer-managed paths
			// (e.g. /var) are derived from it verbatim. The unit name derived
			// from it must also be one systemd can load.
			mountPath := ds.UserDefined.MountPath
			mountPathFld := userDefinedPath.Child("mountPath")
			switch {
			case mountPath == "":
				allErrs = append(allErrs, field.Required(mountPathFld, "a mount path is required for user-defined disks"))
				continue
			case !path.IsAbs(mountPath):
				allErrs = append(allErrs, field.Invalid(mountPathFld, mountPath, "must be an absolute path"))
				continue
			case !mountPathRegexp.MatchString(mountPath):
				allErrs = append(allErrs, field.Invalid(mountPathFld, mountPath, "may only contain alphanumerics, '/', '.', '_' and '-'"))
				continue
			case path.Clean(mountPath) != mountPath:
				allErrs = append(allErrs, field.Invalid(mountPathFld, mountPath,
					fmt.Sprintf("must be a normalized path without trailing '/', '.', '..' or repeated '/' (did you mean %q?)", path.Clean(mountPath))))
				continue
			case mountPath == "/":
				allErrs = append(allErrs, field.Invalid(mountPathFld, mountPath, "a data disk cannot be mounted at the root filesystem"))
				continue
			case len(types.MountUnitName(mountPath)) > types.MaxUnitNameLength:
				allErrs = append(allErrs, field.Invalid(mountPathFld, mountPath,
					fmt.Sprintf("is too long: the name of its systemd mount unit, %q, exceeds the %d characters systemd allows",
						types.MountUnitName(mountPath), types.MaxUnitNameLength)))
				continue
			}
		case types.Etcd:
			if foundEtcd {
				allErrs = append(allErrs, field.TooMany(entryPath.Child("etcd"), 2, 1))
				continue
			}
			if ds.Etcd == nil {
				allErrs = append(allErrs, field.Invalid(entryPath.Child("etcd"), dsYaml, "etcd configuration must be created"))
				continue
			}
			// etcd should only be setup on control plane, not any other machine type.
			if p.Name != "master" {
				allErrs = append(allErrs, field.Invalid(entryPath.Child("etcd"), dsYaml, "cannot specify etcd on worker machine pools"))
				continue
			}
			foundEtcd = true
			// The etcd platformDiskID names the platform data disk the entry
			// is paired with, so an empty ID leaves the entry unpaired. Length
			// limits are left to the platform validations, which know the
			// naming rules of the disks the ID refers to.
			if ds.Etcd.PlatformDiskID == "" {
				allErrs = append(allErrs, field.Required(entryPath.Child("etcd", "platformDiskID"), "a platform disk ID is required"))
				continue
			}
		case types.Swap:
			if p.Name == "master" {
				allErrs = append(allErrs, field.Invalid(entryPath.Child("swap"), dsYaml, "swap is unsupported on control plane nodes"))
				continue
			}
			if foundSwap {
				allErrs = append(allErrs, field.TooMany(entryPath.Child("swap"), 2, 1))
				continue
			}
			if ds.Swap == nil {
				allErrs = append(allErrs, field.Invalid(entryPath.Child("swap"), dsYaml, "swap configuration must be created"))
				continue
			}
			foundSwap = true
			// The swap platformDiskID names the platform data disk the entry
			// is paired with, so an empty ID leaves the entry unpaired. Length
			// limits are left to the platform validations, which know the
			// naming rules of the disks the ID refers to.
			if ds.Swap.PlatformDiskID == "" {
				allErrs = append(allErrs, field.Required(entryPath.Child("swap", "platformDiskID"), "a platform disk ID is required"))
				continue
			}
		default:
			allErrs = append(allErrs, field.NotSupported(entryPath.Child("type"), ds.Type, validDiskTypeValues))
			continue
		}

		// Only entries that are otherwise well-formed can be checked for collisions,
		// since a malformed entry has no disk ID or label to compare.
		diskID, ok := ds.PlatformDiskID()
		if !ok {
			continue
		}
		if seenDiskIDs.Has(diskID) {
			allErrs = append(allErrs, field.Duplicate(entryPath.Child(diskSetupFieldNames[ds.Type], "platformDiskID"), diskID))
		} else {
			seenDiskIDs.Insert(diskID)
		}

		// Labels are compared case-insensitively because the MachineConfig name is
		// lower-cased, so "Data" and "data" name the same object.
		label := strings.ToLower(ds.MachineConfigLabel())
		if seenLabels.Has(label) {
			allErrs = append(allErrs, field.Invalid(entryPath, dsYaml,
				fmt.Sprintf("resolves to disk label %q, which is already used by an earlier diskSetup entry in this pool; "+
					"etcd and swap disks are labelled after their type, user-defined disks after their platformDiskID stripped of non-alphanumeric characters", label)))
		} else {
			seenLabels.Insert(label)
		}

		// Mount paths are compared verbatim: they are already normalized by
		// the per-entry checks, so entries that reach this point resolve to
		// the exact path their mount unit will be named after. The etcd entry
		// participates too, so a user-defined disk cannot claim /var/lib/etcd.
		if mountPath, ok := ds.MountPath(); ok && mountPath != "" {
			if seenMountPaths.Has(mountPath) {
				// The etcd entry has no mountPath field: its type implies the path.
				mountPathFld := entryPath.Child("userDefined", "mountPath")
				if ds.Type == types.Etcd {
					mountPathFld = entryPath.Child("etcd")
				}
				allErrs = append(allErrs, field.Duplicate(mountPathFld, mountPath))
			} else {
				seenMountPaths.Insert(mountPath)
			}
		}
	}

	return allErrs
}

func validateMachinePoolPlatform(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	platformName := platform.Name()
	validate := func(n string, value interface{}, validation func(*field.Path) field.ErrorList) {
		f := fldPath.Child(n)
		if platformName == n || (platformName == powervc.Name && n == openstack.Name) {
			allErrs = append(allErrs, validation(f)...)
		} else {
			allErrs = append(allErrs, field.Invalid(f, value, fmt.Sprintf("cannot specify %q for machine pool when cluster is using %q", n, platformName)))
		}
	}
	if platform.AWS != nil {
		allErrs = append(allErrs, awsvalidation.ValidateAMIID(platform.AWS, p.AWS, fldPath.Child("aws"))...)
	}
	if p.AWS != nil {
		validate(aws.Name, p.AWS, func(f *field.Path) field.ErrorList {
			return awsvalidation.ValidateMachinePool(platform.AWS, p.AWS, pool.Name, f)
		})
	}
	if p.Azure != nil {
		validate(azure.Name, p.Azure, func(f *field.Path) field.ErrorList {
			return azurevalidation.ValidateMachinePool(p.Azure, pool.Name, platform.Azure, pool, f)
		})
	}
	if platform.GCP != nil {
		allErrs = append(allErrs, gcpvalidation.ValidateOSImageForSovereignCloud(platform.GCP, p.GCP, fldPath.Child("gcp"))...)
	}
	if p.GCP != nil {
		validate(gcp.Name, p.GCP, func(f *field.Path) field.ErrorList { return validateGCPMachinePool(platform, p, pool, f) })
	}
	if p.IBMCloud != nil {
		validate(ibmcloud.Name, p.IBMCloud, func(f *field.Path) field.ErrorList {
			return ibmcloudvalidation.ValidateMachinePool(platform.IBMCloud, p.IBMCloud, f)
		})
	}
	if p.BareMetal != nil {
		validate(baremetal.Name, p.BareMetal, func(f *field.Path) field.ErrorList { return baremetalvalidation.ValidateMachinePool(p.BareMetal, f) })
	}
	if p.VSphere != nil {
		validate(vsphere.Name, p.VSphere, func(f *field.Path) field.ErrorList {
			return vspherevalidation.ValidateMachinePool(platform.VSphere, pool, f)
		})
	}
	if p.Ovirt != nil {
		validate(ovirt.Name, p.Ovirt, func(f *field.Path) field.ErrorList { return ovirtvalidation.ValidateMachinePool(p.Ovirt, f) })
	}
	if p.PowerVS != nil {
		validate(powervs.Name, p.PowerVS, func(f *field.Path) field.ErrorList {
			return powervsvalidation.ValidateMachinePool(platform.PowerVS, p.PowerVS, f)
		})
	}
	if p.OpenStack != nil {
		validate(openstack.Name, p.OpenStack, func(f *field.Path) field.ErrorList {
			return openstackvalidation.ValidateMachinePool(platform.OpenStack, p.OpenStack, pool.Name, f)
		})
	}

	// A diskSetup entry is paired with a data disk from the pool's platform
	// stanza, and with no stanza there is no data disk anywhere for the
	// entries to pair with. Catching that here fails create install-config
	// instead of deferring the problem to machine generation during create
	// manifests. Defaulting runs before validation and materializes the
	// stanza from platform.<plat>.defaultMachinePlatform where the platform
	// supports that (Azure does, vSphere does not), and Azure rejects
	// dataDisks on defaultMachinePlatform outright, so a stanza still absent
	// at this point leaves the entries with nothing to pair with on either
	// platform.
	if len(pool.DiskSetup) > 0 {
		switch platformName {
		case azure.Name:
			if p.Azure == nil {
				allErrs = append(allErrs, field.Required(fldPath.Child(azure.Name, "dataDisks"),
					"diskSetup entries must be paired with data disks: set dataDisks on this machine pool's platform.azure stanza"))
			}
		case vsphere.Name:
			if p.VSphere == nil {
				allErrs = append(allErrs, field.Required(fldPath.Child(vsphere.Name, "dataDisks"),
					"diskSetup entries must be paired with data disks: set dataDisks on this machine pool's platform.vsphere stanza"))
			}
		}
	}
	return allErrs
}

func validateGCPMachinePool(platform *types.Platform, p *types.MachinePoolPlatform, pool *types.MachinePool, f *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, gcpvalidation.ValidateMachinePool(platform.GCP, p.GCP, f)...)
	allErrs = append(allErrs, gcpvalidation.ValidateMasterDiskType(pool, f)...)
	allErrs = append(allErrs, gcpvalidation.ValidateServiceAccount(platform.GCP, pool, f)...)

	return allErrs
}
