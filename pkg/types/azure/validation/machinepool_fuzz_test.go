package validation

import (
	"math"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/validation/field"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

func FuzzValidateDataDisk(f *testing.F) {
	f.Add("worker", "etcd", "swap", int32(0), int32(1), int32(1), int32(10), true, true, "Enabled", "", "", "")
	f.Add("master", "etcd", "etcd", int32(-1), int32(0), int32(0), int32(1), true, true, "", "", "", "")
	f.Add("worker", "disk", "other", int32(64), int32(63), int32(-1), int32(1), false, true, "Disabled", "UltraSSD_LRS", "ReadWrite", "")
	f.Add("worker", "disk", "disk2", int32(0), int32(1), int32(1), int32(1), true, true, "", "UltraSSD_LRS", "", "Enabled")
	f.Add("edge", "disk", "disk2", int32(0), int32(1), int32(1), int32(1), true, true, "", "", "", "")
	f.Add("", "a", "b", int32(0), int32(0), int32(0), int32(0), false, false, "", "", "", "")

	f.Fuzz(func(t *testing.T, poolName, suffix1, suffix2 string, lun1, lun2, size1, size2 int32, hasLun1, hasLun2 bool, capability, storageType, caching, defaultCapability string) {
		pool := &azure.MachinePool{
			UltraSSDCapability: capability,
			DataDisks: []capz.DataDisk{
				dataDisk(suffix1, lun1, size1, hasLun1, storageType, caching),
				dataDisk(suffix2, lun2, size2, hasLun2, storageType, caching),
			},
		}
		platform := &azure.Platform{}
		if defaultCapability != "" {
			platform.DefaultMachinePlatform = &azure.MachinePool{UltraSSDCapability: defaultCapability}
		}
		errs := validateDataDisk(pool, poolName, platform, field.NewPath("dataDisks"))
		assertDataDiskErrors(t, poolName, pool, capability, defaultCapability, errs)
	})
}

func FuzzValidateDataDiskSetup(f *testing.F) {
	f.Add(1, 1, "etcd", "etcd", "etcd")
	f.Add(0, 1, "swap", "swap", "other")
	f.Add(2, 1, "user-defined", "userdisk", "userdisk")
	f.Add(1, 2, "etcd", "mismatch", "disk")
	f.Add(2, 2, "other", "id", "suffix")

	f.Fuzz(func(t *testing.T, nDisks, nSetups int, diskType, platformDiskID, nameSuffix string) {
		nDisks = clampCount(nDisks)
		nSetups = clampCount(nSetups)
		azurePool := &azure.MachinePool{DataDisks: make([]capz.DataDisk, nDisks)}
		for i := range azurePool.DataDisks {
			azurePool.DataDisks[i].NameSuffix = nameSuffix
		}
		pool := &types.MachinePool{DiskSetup: make([]types.Disk, nSetups)}
		for i := range pool.DiskSetup {
			pool.DiskSetup[i] = diskSetup(diskType, platformDiskID)
		}
		errs := validateDataDiskSetup(azurePool, pool, field.NewPath("dataDisks"))
		if nDisks < nSetups {
			if len(errs) != 1 || errs[0] == nil || errs[0].Type != field.ErrorTypeInvalid {
				t.Fatalf("diskSetup length %d dataDisks length %d = %v, want one Invalid error", nSetups, nDisks, errs)
			}
			return
		}
		var want int
		if diskSetupMismatches(diskType, platformDiskID, nameSuffix) {
			want = nSetups
		}
		if len(errs) != want {
			t.Fatalf("diskSetup %q id %q suffix %q errors=%d, want %d (%v)", diskType, platformDiskID, nameSuffix, len(errs), want, errs)
		}
		for _, err := range errs {
			if err == nil || err.Type != field.ErrorTypeInvalid {
				t.Fatalf("unexpected diskSetup error %v", err)
			}
		}
	})
}

func FuzzValidateOSImage(f *testing.F) {
	f.Add("", "", "", "", "")
	f.Add(string(azure.ImageWithPurchasePlan), "pub", "offer", "sku", "1.0")
	f.Add(string(azure.ImageNoPurchasePlan), "", "", "", "")
	f.Add("OtherPlan", "pub", "", "sku", "")
	f.Add("", "pub", "offer", "sku", "1.0")

	f.Fuzz(func(t *testing.T, plan, publisher, offer, sku, version string) {
		pool := &azure.MachinePool{OSImage: azure.OSImage{
			Plan:      azure.ImagePurchasePlan(plan),
			Publisher: publisher,
			Offer:     offer,
			SKU:       sku,
			Version:   version,
		}}
		errs := validateOSImage(pool, field.NewPath("pool"))
		want := expectedOSImageErrors(plan, publisher, offer, sku, version)
		if len(errs) != want {
			t.Fatalf("osImage plan=%q publisher=%q offer=%q sku=%q version=%q errors=%d, want %d (%v)", plan, publisher, offer, sku, version, len(errs), want, errs)
		}
		for _, err := range errs {
			if err == nil {
				t.Fatal("validateOSImage returned a nil error")
			}
			if err.Type != field.ErrorTypeRequired && err.Type != field.ErrorTypeNotSupported {
				t.Fatalf("unexpected osImage error type %q", err.Type)
			}
		}
	})
}

func FuzzValidateSecurityProfile(f *testing.F) {
	// Settings absent and an empty security profile used to panic.
	f.Add(string(azure.PublicCloud), "", "", "", "", false, true, false, false, false, false)
	f.Add(string(azure.StackCloud), string(azure.SecurityTypesConfidentialVM), string(azure.SecurityEncryptionTypesVMGuestStateOnly), "Enabled", "Enabled", true, true, false, true, true, false)
	f.Add(string(azure.PublicCloud), string(azure.SecurityTypesConfidentialVM), "", "", "", true, false, false, false, false, false)
	f.Add(string(azure.PublicCloud), string(azure.SecurityTypesTrustedLaunch), "", "Enabled", "Enabled", true, false, false, false, false, false)
	f.Add(string(azure.PublicCloud), "", string(azure.SecurityEncryptionTypesDiskWithVMGuestState), "", "", false, true, true, false, false, false)
	f.Add(string(azure.PublicCloud), string(azure.SecurityTypesConfidentialVM), string(azure.SecurityEncryptionTypesDiskWithVMGuestState), "Disabled", "Disabled", true, true, true, true, true, false)

	f.Fuzz(func(t *testing.T, cloud, securityType, encryptionType, vtpm, secureBoot string, hasSettings, hasProfile, encryptionAtHost, hasConfidential, hasUEFI, hasTrustedLaunch bool) {
		pool := &azure.MachinePool{EncryptionAtHost: encryptionAtHost}
		if hasSettings {
			pool.Settings = &azure.SecuritySettings{SecurityType: azure.SecurityTypes(securityType)}
			if hasConfidential {
				pool.Settings.ConfidentialVM = &azure.ConfidentialVM{}
				if hasUEFI {
					pool.Settings.ConfidentialVM.UEFISettings = uefiSettings(vtpm, secureBoot)
				}
			}
			if hasTrustedLaunch {
				pool.Settings.TrustedLaunch = &azure.TrustedLaunch{}
				if hasUEFI {
					pool.Settings.TrustedLaunch.UEFISettings = uefiSettings(vtpm, secureBoot)
				}
			}
		}
		if hasProfile {
			pool.OSDisk.SecurityProfile = &azure.VMDiskSecurityProfile{
				SecurityEncryptionType: azure.SecurityEncryptionTypes(encryptionType),
			}
		}
		errs := validateSecurityProfile(pool, azure.CloudEnvironment(cloud), field.NewPath("pool"))
		if len(errs) > 1 {
			t.Fatalf("got %d security profile errors, want at most 1: %v", len(errs), errs)
		}
		want := expectedSecurityProfileError(cloud, securityType, encryptionType, vtpm, secureBoot, hasSettings, hasProfile, encryptionAtHost, hasConfidential, hasUEFI, hasTrustedLaunch)
		if want == "" {
			if len(errs) != 0 {
				t.Fatalf("security profile accepted inputs got %v", errs)
			}
			return
		}
		if len(errs) != 1 || errs[0] == nil || errs[0].Detail != want {
			t.Fatalf("security profile error = %v, want %q", errs, want)
		}
	})
}

func dataDisk(suffix string, lun, size int32, hasLun bool, storageType, caching string) capz.DataDisk {
	disk := capz.DataDisk{NameSuffix: suffix, DiskSizeGB: size, CachingType: caching}
	if hasLun {
		lunCopy := lun
		disk.Lun = &lunCopy
	}
	if storageType != "" {
		disk.ManagedDisk = &capz.ManagedDiskParameters{StorageAccountType: storageType}
	}
	return disk
}

func diskSetup(diskType, platformDiskID string) types.Disk {
	setup := types.Disk{Type: types.DiskType(diskType)}
	switch setup.Type {
	case types.Etcd:
		setup.Etcd = &types.DiskEtcd{PlatformDiskID: platformDiskID}
	case types.Swap:
		setup.Swap = &types.DiskSwap{PlatformDiskID: platformDiskID}
	case types.UserDefined:
		setup.UserDefined = &types.DiskUserDefined{PlatformDiskID: platformDiskID}
	}
	return setup
}

func uefiSettings(vtpm, secureBoot string) *azure.UEFISettings {
	settings := &azure.UEFISettings{}
	if vtpm != "" {
		settings.VirtualizedTrustedPlatformModule = &vtpm
	}
	if secureBoot != "" {
		settings.SecureBoot = &secureBoot
	}
	return settings
}

func assertDataDiskErrors(t *testing.T, poolName string, pool *azure.MachinePool, capability, defaultCapability string, errs field.ErrorList) {
	t.Helper()
	if poolName != types.MachinePoolControlPlaneRoleName && poolName != types.MachinePoolComputeRoleName {
		if len(errs) != 0 {
			t.Fatalf("pool %q is not a control-plane or compute pool, got %v", poolName, errs)
		}
		return
	}
	for _, err := range errs {
		if err == nil {
			t.Fatal("validateDataDisk returned a nil error")
		}
	}
	if len(errs) > 8*len(pool.DataDisks) {
		t.Fatalf("got %d errors for %d disks", len(errs), len(pool.DataDisks))
	}

	resolvedCapability := capability
	if resolvedCapability == "" {
		resolvedCapability = defaultCapability
	}
	var missingLun, lunRange, duplicateLun, duplicateName, badSize, ultraCapability, ultraCache int
	seenSuffix := map[string]int{}
	seenLun := map[int32]int{}
	for _, disk := range pool.DataDisks {
		seenSuffix[disk.NameSuffix]++
		if disk.Lun == nil {
			missingLun++
		} else {
			if *disk.Lun < 0 || *disk.Lun > 63 {
				lunRange++
			}
			seenLun[*disk.Lun]++
		}
		if disk.DiskSizeGB <= 0 {
			badSize++
		}
		if disk.ManagedDisk != nil && disk.ManagedDisk.StorageAccountType == ultraSSDStorageAccountType {
			if resolvedCapability != enabled {
				ultraCapability++
			}
			if disk.CachingType != "" && disk.CachingType != cachingTypeNone {
				ultraCache++
			}
		}
	}
	for _, count := range seenSuffix {
		if count > 1 {
			duplicateName += count - 1
		}
	}
	for _, count := range seenLun {
		if count > 1 {
			duplicateLun += count - 1
		}
	}
	got := countDataDiskKinds(errs)
	if got.missingLun != missingLun || got.lunRange != lunRange || got.duplicateLun != duplicateLun || got.duplicateName != duplicateName || got.badSize != badSize || got.ultraCapability != ultraCapability || got.ultraCache != ultraCache {
		t.Fatalf("data disk errors = %+v, want missing=%d range=%d dupLun=%d dupName=%d size=%d ultra=%d cache=%d (%v)",
			got, missingLun, lunRange, duplicateLun, duplicateName, badSize, ultraCapability, ultraCache, errs)
	}
}

type dataDiskKindCounts struct {
	missingLun, lunRange, duplicateLun, duplicateName, badSize, ultraCapability, ultraCache int
}

func countDataDiskKinds(errs field.ErrorList) dataDiskKindCounts {
	var got dataDiskKindCounts
	for _, err := range errs {
		switch {
		case err.Type == field.ErrorTypeDuplicate:
			got.duplicateName++
		case err.Detail == "dataDisk must have a unique lun number":
			got.duplicateLun++
		case strings.Contains(err.Detail, "must have lun id between 0 and 63"):
			got.lunRange++
		case strings.Contains(err.Detail, "must have lun id"):
			got.missingLun++
		case err.Detail == "diskSizeGB must be greater than zero":
			got.badSize++
		case strings.Contains(err.Detail, "ultraSSDCapability must be"):
			got.ultraCapability++
		case strings.Contains(err.Detail, "cachingType must be"):
			got.ultraCache++
		}
	}
	return got
}

func diskSetupMismatches(diskType, platformDiskID, nameSuffix string) bool {
	switch types.DiskType(diskType) {
	case types.Etcd, types.Swap, types.UserDefined:
		return platformDiskID != nameSuffix
	default:
		return false
	}
}

func expectedOSImageErrors(plan, publisher, offer, sku, version string) int {
	if plan == "" && publisher == "" && offer == "" && sku == "" && version == "" {
		return 0
	}
	n := 0
	if plan != "" && plan != string(azure.ImageWithPurchasePlan) && plan != string(azure.ImageNoPurchasePlan) {
		n++
	}
	if publisher == "" {
		n++
	}
	if offer == "" {
		n++
	}
	if sku == "" {
		n++
	}
	if version == "" {
		n++
	}
	return n
}

func expectedSecurityProfileError(cloud, securityType, encryptionType, vtpm, secureBoot string, hasSettings, hasProfile, encryptionAtHost, hasConfidential, hasUEFI, hasTrustedLaunch bool) string {
	if !hasSettings {
		if hasProfile && encryptionType != "" {
			return "settings should be set when osDisk.securityProfile.securityEncryptionType is defined."
		}
		return ""
	}
	if cloud == string(azure.StackCloud) {
		return "the securityType field is not supported on AzureStackCloud."
	}
	switch azure.SecurityTypes(securityType) {
	case azure.SecurityTypesConfidentialVM:
		if !hasProfile || encryptionType == "" {
			return "securityEncryptionType should be set when securityType is set to ConfidentialVM."
		}
		if encryptionType != string(azure.SecurityEncryptionTypesVMGuestStateOnly) && encryptionType != string(azure.SecurityEncryptionTypesDiskWithVMGuestState) {
			return `supported values: "DiskWithVMGuestState", "VMGuestStateOnly"`
		}
		if !hasConfidential {
			return "confidentialVM should be set when securityType is set to ConfidentialVM."
		}
		if !hasUEFI {
			return "uefiSettings should be set when securityType is set to ConfidentialVM."
		}
		if vtpm != "" && vtpm != enabled {
			return "virtualizedTrustedPlatformModule should be enabled when securityType is set to ConfidentialVM."
		}
		if encryptionType == string(azure.SecurityEncryptionTypesDiskWithVMGuestState) {
			if encryptionAtHost {
				return "encryptionAtHost cannot be set to true when securityEncryptionType is set to DiskWithVMGuestState."
			}
			if secureBoot != "" && secureBoot != enabled {
				return "secureBoot should be enabled when securityEncryptionType is set to DiskWithVMGuestState."
			}
		}
		return ""
	case azure.SecurityTypesTrustedLaunch:
		if !hasTrustedLaunch {
			return "trustedLaunch should be set when securityType is set to TrustedLaunch."
		}
		return ""
	default:
		if hasProfile && encryptionType != "" {
			return "securityType should be set to ConfidentialVM when securityEncryptionType is defined."
		}
		if hasTrustedLaunch && hasUEFI && (secureBoot == enabled || vtpm == enabled) {
			return "securityType should be set to TrustedLaunch when uefiSettings are enabled."
		}
		return ""
	}
}

func clampCount(n int) int {
	// Negating math.MinInt overflows and stays negative, which panics in make.
	if n < -math.MaxInt {
		n = math.MaxInt
	} else if n < 0 {
		n = -n
	}
	return n % 3
}
