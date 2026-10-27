package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/azure/defaults"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name          string
		azurePlatform azure.CloudEnvironment
		pool          *types.MachinePool
		// Optional platform-wide defaults for testing fallback and pool overrides.
		// Nil means no defaultMachinePlatform is configured.
		defaultMachinePlatform *azure.MachinePool
		expected               string
	}{
		{
			name:          "empty",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				},
			},
		},
		{
			name:          "valid iops",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 120,
						},
					},
				},
			},
		},
		{
			name:          "invalid iops",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: -120,
						},
					},
				},
			},
			expected: `^test-path\.diskSizeGB: Invalid value: -120: Storage DiskSizeGB must be positive$`,
		},
		{
			name:          "valid disk",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Premium_LRS",
						},
					},
				},
			},
		},
		{
			name:          "invalid disk",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "LRS",
						},
					},
				},
			},
			expected: `^test-path\.diskType: Unsupported value: "LRS": supported values: "Premium_LRS", "StandardSSD_LRS", "Standard_LRS"$`,
		},
		{
			name:          "multiple disk and setup missing lun id",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: "etcd",
					Etcd: &types.DiskEtcd{
						PlatformDiskID: "etcd",
					},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         nil,
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.lun: Required value: \"etcd\" must have lun id$`,
		},
		{
			name:          "lun id must be below 64",
			azurePlatform: azure.PublicCloud,

			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: "etcd",
					Etcd: &types.DiskEtcd{
						PlatformDiskID: "etcd",
					},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         ptr.To(int32(64)),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.lun: Invalid value: 64: \"etcd\" must have lun id between 0 and 63$`,
		},
		{
			name:          "multiple disk and setup PlatformDiskID does not match",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: "etcd",
					Etcd: &types.DiskEtcd{
						PlatformDiskID: "etcd",
					},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "foo",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
						},
						},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Invalid value: \"foo\": does not match etcd platformDiskID \"etcd\"$`,
		},
		{
			name:          "lun id must be above 0",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: "etcd",
					Etcd: &types.DiskEtcd{
						PlatformDiskID: "etcd",
					},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         ptr.To(int32(-1)),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.lun: Invalid value: -1: \"etcd\" must have lun id between 0 and 63$`,
		},
		{
			name:          "multiple disk size must be greater than zero",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: "etcd",
					Etcd: &types.DiskEtcd{
						PlatformDiskID: "etcd",
					},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  0,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.diskSizeGB: Invalid value: 0: diskSizeGB must be greater than zero$`,
		},
		{
			name:          "datadisks in default machine pool is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  0,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks: Invalid value: \"etcd\": not allowed on default machine pool, use dataDisks compute and controlPlane only$`,
		},
		{
			name:          "lun id must be unique",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{
					{
						Type: "etcd",
						Etcd: &types.DiskEtcd{
							PlatformDiskID: "etcd",
						},
					},
					{
						Type: "user-defined",
						UserDefined: &types.DiskUserDefined{
							PlatformDiskID: "containers",
							MountPath:      "/var/lib/containers",
						},
					},
				},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{
							{
								NameSuffix:  "etcd",
								DiskSizeGB:  32,
								ManagedDisk: nil,
								Lun:         pointer.Int32(0),
								CachingType: "",
							},
							{
								NameSuffix:  "containers",
								DiskSizeGB:  32,
								ManagedDisk: nil,
								Lun:         pointer.Int32(0),
								CachingType: "",
							},
						},
					},
				},
			},
			expected: `^test-path\.dataDisks\[1\]\.lun: Duplicate value: 0$`,
		},
		{
			name:          "empty nameSuffix is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Required value: nameSuffix is required$`,
		},
		{
			name:          "nameSuffix with characters Azure rejects in disk names is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "bad/name",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Invalid value: "bad/name": may only contain alphanumerics, '_', '\.' and '-', and must end with an alphanumeric or '_'$`,
		},
		{
			name:          "nameSuffix ending in a period is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "data.",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Invalid value: "data\.": may only contain alphanumerics, '_', '\.' and '-', and must end with an alphanumeric or '_'$`,
		},
		{
			name:          "nameSuffix longer than the Azure disk name limit is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "aaaaaaaaaaaaaaaa",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Too long: may not be more than 15 bytes$`,
		},
		{
			// A diskSetup entry carries its data disk's nameSuffix as the
			// platformDiskID, so an over-length ID surfaces through the
			// nameSuffix bound the entry is required to match.
			name:          "etcd diskSetup platformDiskID longer than the Azure disk name limit is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{{
					Type: types.Etcd,
					Etcd: &types.DiskEtcd{PlatformDiskID: "etcddisketcddisk1"},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcddisketcddisk1",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Too long: may not be more than 15 bytes$`,
		},
		{
			name:          "swap diskSetup platformDiskID longer than the Azure disk name limit is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				DiskSetup: []types.Disk{{
					Type: types.Swap,
					Swap: &types.DiskSwap{PlatformDiskID: "swapdiskswapdisk1"},
				}},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "swapdiskswapdisk1",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.nameSuffix: Too long: may not be more than 15 bytes$`,
		},
		{
			name:          "nameSuffix at the maximum length is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "aaaaaaaaaaaaaaa",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
		},
		{
			name:          "nameSuffix with the full Azure-allowed character set is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name:      "master",
				DiskSetup: []types.Disk{},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "my_disk-1",
							DiskSizeGB:  32,
							ManagedDisk: nil,
							Lun:         pointer.Int32(0),
							CachingType: "",
						}},
					},
				},
			},
		},
		{
			name:          "unsupported disk master",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Standard_LRS",
						},
					},
				},
			},
			expected: `^test-path\.diskType: Unsupported value: "Standard_LRS": supported values: "Premium_LRS", "StandardSSD_LRS"$`,
		},
		{
			name:          "unsupported disk default pool",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Standard_LRS",
						},
					},
				},
			},
			expected: `^test-path\.diskType: Unsupported value: "Standard_LRS": supported values: "Premium_LRS", "StandardSSD_LRS"$`,
		},
		{
			name:          "supported disk worker",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskType: "Standard_LRS",
						},
					},
				},
			},
		},
		{
			name: "valid OS image",
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
							Version:   "test-version",
							Plan:      "NoPurchasePlan",
						},
					},
				},
			},
		},
		{
			name:          "valid OS image with purchase plan omitted",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
							Version:   "test-version",
						},
					},
				},
			},
		},
		{
			name:          "OS image missing publisher",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Offer:   "test-offer",
							SKU:     "test-sku",
							Version: "test-version",
						},
					},
				},
			},
			expected: `^test-path\.osImage.publisher: Required value: must specify publisher for the OS image$`,
		},
		{
			name:          "OS image missing offer",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							SKU:       "test-sku",
							Version:   "test-version",
						},
					},
				},
			},
			expected: `^test-path\.osImage.offer: Required value: must specify offer for the OS image$`,
		},
		{
			name:          "OS image missing sku",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							Version:   "test-version",
						},
					},
				},
			},
			expected: `^test-path\.osImage.sku: Required value: must specify SKU for the OS image$`,
		},
		{
			name:          "OS image missing version",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
						},
					},
				},
			},
			expected: `^test-path\.osImage.version: Required value: must specify version for the OS image$`,
		},
		{
			name:          "OS image with invalid plan",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
							Version:   "test-version",
							Plan:      "test-plan",
						},
					},
				},
			},
			expected: `^test-path\.osImage.plan: Unsupported value: ".*": supported values: "NoPurchasePlan", "WithPurchasePlan"$`,
		},
		{
			name:          "OS image for master",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
							Version:   "test-version",
						},
					},
				},
			},
		},
		{
			name:          "OS image for default pool",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSImage: azure.OSImage{
							Publisher: "test-publisher",
							Offer:     "test-offer",
							SKU:       "test-sku",
							Version:   "test-version",
						},
					},
				},
			},
		},
		{
			name:          "insufficient disk size azurestack",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 120,
						},
					},
				},
			},
			expected: `^test-path.diskSizeGB: Invalid value: 120: Storage DiskSizeGB must be between 128 and 1023 inclusive for Azure Stack$`,
		},
		{
			name:          "excessive disk size azurestack",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
						},
					},
				},
			},
			expected: `^test-path.diskSizeGB: Invalid value: 1200: Storage DiskSizeGB must be between 128 and 1023 inclusive for Azure Stack$`,
		},
		{
			name:          "empty disk size azurestack",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				},
			},
		},
		{
			name:          "empty settings and securityEncryptionType",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
						},
					},
				},
			},
		},
		{
			name:          "undefined settings when securityEncryptionType is defined",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings: Required value: settings should be set when osDisk.securityProfile.securityEncryptionType is defined.$`,
		},
		{
			name:          "securityType set to ConfidentialVM and platform to AzureStackCloud",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 128,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.securityType: Invalid value: "ConfidentialVM": the securityType field is not supported on AzureStackCloud.$`,
		},
		{
			name:          "securityType set to TrustedLaunch and platform to AzureStackCloud",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 128,
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesTrustedLaunch,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.securityType: Invalid value: "TrustedLaunch": the securityType field is not supported on AzureStackCloud.$`,
		},
		{
			name:          "securityType set to ConfidentialVM but securityEncryptionType is empty",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.osDisk.securityProfile.securityEncryptionType: Required value: securityEncryptionType should be set when securityType is set to ConfidentialVM.$`,
		},
		{
			name:          "securityType set to ConfidentialVM but securityEncryptionType is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypes("invalidSecurityEncryptionType"),
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.osDisk.securityProfile.securityEncryptionType: Unsupported value: "invalidSecurityEncryptionType": supported values: "DiskWithVMGuestState", "VMGuestStateOnly"$`,
		},
		{
			name:          "securityType set to ConfidentialVM, securityEncryptionType is set but confidentialVM section is empty",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.confidentialVM: Required value: confidentialVM should be set when securityType is set to ConfidentialVM.$`,
		},
		{
			name:          "securityType set to ConfidentialVM, securityEncryptionType is set but uefiSettings is not set",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType:   azure.SecurityTypesConfidentialVM,
							ConfidentialVM: &azure.ConfidentialVM{},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.confidentialVM.uefiSettings: Required value: uefiSettings should be set when securityType is set to ConfidentialVM.$`,
		},
		{
			name:          "securityType set to ConfidentialVM, securityEncryptionType is set but virtualizedTrustedPlatformModule is not enabled",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
							ConfidentialVM: &azure.ConfidentialVM{
								UEFISettings: &azure.UEFISettings{
									VirtualizedTrustedPlatformModule: pointer.String("Disabled"),
								},
							},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.confidentialVM.uefiSettings.virtualizedTrustedPlatformModule: Invalid value: "Disabled": virtualizedTrustedPlatformModule should be enabled when securityType is set to ConfidentialVM.$`,
		},
		{
			name:          "encryptionAtHost cannot be true when securityEncryptionType is set to DiskWithVMGuestState",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						EncryptionAtHost: true,
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesDiskWithVMGuestState,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
							ConfidentialVM: &azure.ConfidentialVM{
								UEFISettings: &azure.UEFISettings{
									VirtualizedTrustedPlatformModule: pointer.String("Enabled"),
								},
							},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.encryptionAtHost: Invalid value: true: encryptionAtHost cannot be set to true when securityEncryptionType is set to DiskWithVMGuestState.$`,
		},
		{
			name:          "encryptionAtHost cannot be true when securityEncryptionType is set to DiskWithVMGuestState",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesDiskWithVMGuestState,
							},
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesConfidentialVM,
							ConfidentialVM: &azure.ConfidentialVM{
								UEFISettings: &azure.UEFISettings{
									VirtualizedTrustedPlatformModule: pointer.String("Enabled"),
									SecureBoot:                       pointer.String("Disabled"),
								},
							},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.confidentialVM.uefiSettings.secureBoot: Invalid value: "Disabled": secureBoot should be enabled when securityEncryptionType is set to DiskWithVMGuestState.$`,
		},
		{
			name:          "securityType is TrustedLaunch but trustedLaunch section is not defined",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
						},
						Settings: &azure.SecuritySettings{
							SecurityType: azure.SecurityTypesTrustedLaunch,
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.trustedLaunch: Required value: trustedLaunch should be set when securityType is set to TrustedLaunch.$`,
		},
		{
			name:          "securityEncryptionType is set but securityType is not set to ConfidentialVM",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
							SecurityProfile: &azure.VMDiskSecurityProfile{
								SecurityEncryptionType: azure.SecurityEncryptionTypesVMGuestStateOnly,
							},
						},
						Settings: &azure.SecuritySettings{},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.securityType: Invalid value: "": securityType should be set to ConfidentialVM when securityEncryptionType is defined.$`,
		},
		{
			name:          "securityEncryptionType is set but securityType is not set to ConfidentialVM",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						OSDisk: azure.OSDisk{
							DiskSizeGB: 1200,
						},
						Settings: &azure.SecuritySettings{
							TrustedLaunch: &azure.TrustedLaunch{
								UEFISettings: &azure.UEFISettings{
									VirtualizedTrustedPlatformModule: pointer.String("Enabled"),
								},
							},
						},
					},
				},
			},
			expected: `^test-path.defaultMachinePlatform.settings.securityType: Invalid value: "": securityType should be set to TrustedLaunch when uefiSettings are enabled.$`,
		},
		{
			name:          "master data disk without security profile passes",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 64,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			name:          "master data disk with SecurityProfile and SecurityEncryptionType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 64,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
								SecurityProfile: &capz.VMDiskSecurityProfile{
									SecurityEncryptionType: capz.SecurityEncryptionTypeVMGuestStateOnly,
								},
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.securityProfile: Forbidden: data disk security profiles are not yet supported by the Machine API and would be silently ignored`,
		},
		{
			name:          "master data disk with SecurityProfile but no SecurityEncryptionType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 64,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
								SecurityProfile:    &capz.VMDiskSecurityProfile{},
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.securityProfile: Forbidden: data disk security profiles are not yet supported by the Machine API and would be silently ignored`,
		},
		{
			name:          "worker data disk with SecurityProfile and SecurityEncryptionType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 128,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
								SecurityProfile: &capz.VMDiskSecurityProfile{
									SecurityEncryptionType: capz.SecurityEncryptionTypeVMGuestStateOnly,
								},
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.securityProfile: Forbidden: data disk security profiles are not yet supported by the Machine API and would be silently ignored`,
		},
		{
			name:          "worker data disk with SecurityProfile but no SecurityEncryptionType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 128,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
								SecurityProfile:    &capz.VMDiskSecurityProfile{},
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.securityProfile: Forbidden: data disk security profiles are not yet supported by the Machine API and would be silently ignored`,
		},
		{
			name:          "master data disk with managedDisk but missing storageAccountType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "etcd",
							DiskSizeGB:  128,
							ManagedDisk: &capz.ManagedDiskParameters{},
							Lun:         ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Required value: storageAccountType is required when managedDisk is specified`,
		},
		{
			name:          "worker data disk with managedDisk but missing storageAccountType is rejected",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 128,
							ManagedDisk: &capz.ManagedDiskParameters{
								DiskEncryptionSet: &capz.DiskEncryptionSetParameters{
									ID: "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Compute/diskEncryptionSets/des",
								},
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Required value: storageAccountType is required when managedDisk is specified`,
		},
		{
			name:          "data disk without managedDisk passes without storageAccountType",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 128,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			name:          "data disks with no security settings pass validation",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "data",
							DiskSizeGB:  128,
							ManagedDisk: nil,
							Lun:         ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			name:          "data disks on Azure Stack Hub are rejected",
			azurePlatform: azure.StackCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 128,
							ManagedDisk: &capz.ManagedDiskParameters{
								StorageAccountType: "Premium_LRS",
							},
							Lun: ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `test-path\.dataDisks: Forbidden: data disks are not supported on AzureStackCloud`,
		},
		{
			name:          "azure VM Identity is unrecognized type",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						Identity: &azure.VMIdentity{
							Type: "unrecognized",
						},
					},
				},
			},
			expected: `^test-path.identity.type: Unsupported value: "unrecognized": supported values: "None", "UserAssigned"$`,
		},
		{
			name:          "azure VM SystemAssignedIdentity is not allowed",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						Identity: &azure.VMIdentity{
							Type: "SystemAssigned",
						},
					},
				},
			},
			expected: `^test-path.identity.type: Unsupported value: "SystemAssigned": supported values: "None", "UserAssigned"$`,
		},
		{
			name:          "azure VM identity cannot mismatch type and field",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						Identity: &azure.VMIdentity{
							Type:                   capz.VMIdentityNone,
							UserAssignedIdentities: []azure.UserAssignedIdentity{},
						},
					},
				},
			},
			expected: `^test-path.identity.type: Invalid value: "None": userAssignedIdentities may only be used with type: UserAssigned$`,
		},
		{
			name:          "more diskSetup entries than dataDisks",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{
					{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "etcd"}},
					{Type: types.UserDefined, UserDefined: &types.DiskUserDefined{PlatformDiskID: "containers", MountPath: "/var/lib/containers"}},
				},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "etcd",
							DiskSizeGB: 32,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks: Invalid value: 1: each diskSetup entry configures the dataDisk at the same index; diskSetup length is 2 but dataDisks length is 1, so the extra diskSetup entries would be silently ignored$`,
		},
		{
			name:          "diskSetup with no dataDisks at all",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				DiskSetup: []types.Disk{
					{Type: types.Etcd, Etcd: &types.DiskEtcd{PlatformDiskID: "etcd"}},
				},
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				},
			},
			expected: `^test-path\.dataDisks: Invalid value: 0: each diskSetup entry configures the dataDisk at the same index; diskSetup length is 1 but dataDisks length is 0, so the extra diskSetup entries would be silently ignored$`,
		},
		{
			name:          "dataDisks without diskSetup are still checked for a lun id",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "containers",
							DiskSizeGB: 32,
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.lun: Required value: \"containers\" must have lun id$`,
		},
		{
			name:          "dataDisks without diskSetup are still checked for size",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "containers",
							DiskSizeGB: 0,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.diskSizeGB: Invalid value: 0: diskSizeGB must be greater than zero$`,
		},
		{
			name:          "duplicate dataDisk name suffix",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{
							{
								NameSuffix: "containers",
								DiskSizeGB: 32,
								Lun:        ptr.To(int32(0)),
							},
							{
								NameSuffix: "containers",
								DiskSizeGB: 32,
								Lun:        ptr.To(int32(1)),
							},
						},
					},
				},
			},
			expected: `^test-path\.dataDisks\[1\]\.nameSuffix: Duplicate value: \"containers\"$`,
		},
		{
			name:          "unsupported dataDisk storage account type",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "containers",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "Premium_LRS_v9"},
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Unsupported value: "Premium_LRS_v9": supported values: .+$`,
		},
		{
			name:          "supported dataDisk storage account type",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "containers",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(1)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "StandardSSD_LRS"},
						}},
					},
				},
			},
		},
		{
			name:          "ultra ssd dataDisk without the capability",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Invalid value: \"UltraSSD_LRS\": ultraSSDCapability must be \"Enabled\" on this pool or on platform\.azure\.defaultMachinePlatform to attach an UltraSSD_LRS data disk$`,
		},
		{
			name:          "ultra ssd dataDisk with the capability disabled",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Disabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Invalid value: \"UltraSSD_LRS\": ultraSSDCapability must be \"Enabled\" on this pool or on platform\.azure\.defaultMachinePlatform to attach an UltraSSD_LRS data disk$`,
		},
		{
			name:          "ultra ssd dataDisk with the capability enabled on the pool",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Enabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
		},
		{
			name:                   "ultra ssd dataDisk with the capability enabled only on defaultMachinePlatform",
			azurePlatform:          azure.PublicCloud,
			defaultMachinePlatform: &azure.MachinePool{UltraSSDCapability: "Enabled"},
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
		},
		{
			// Same precedence MachinePool.Set applies: a value on the pool wins,
			// so this stays an error despite defaultMachinePlatform enabling it.
			name:                   "ultra ssd dataDisk with the pool overriding defaultMachinePlatform",
			azurePlatform:          azure.PublicCloud,
			defaultMachinePlatform: &azure.MachinePool{UltraSSDCapability: "Enabled"},
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Disabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.managedDisk\.storageAccountType: Invalid value: \"UltraSSD_LRS\": ultraSSDCapability must be \"Enabled\" on this pool or on platform\.azure\.defaultMachinePlatform to attach an UltraSSD_LRS data disk$`,
		},
		{
			name:          "unsupported dataDisk caching type",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "containers",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							CachingType: "readwrite",
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.cachingType: Unsupported value: "readwrite": supported values: .+$`,
		},
		{
			name:          "ultra ssd dataDisk cannot be cached",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Enabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							CachingType: "ReadWrite",
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.cachingType: Invalid value: \"ReadWrite\": cachingType must be \"None\" or left unset for an UltraSSD_LRS data disk$`,
		},
		{
			name:          "ultra ssd dataDisk with caching explicitly disabled",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Enabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							CachingType: "None",
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
		},
		{
			// An unset caching type is defaulted to None, so an empty value
			// must not be reported as a conflict.
			name:          "ultra ssd dataDisk with caching left unset",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						UltraSSDCapability: "Enabled",
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "ultra",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
						}},
					},
				},
			},
		},
		{
			name:          "caching is allowed on a non-ultra dataDisk",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix:  "containers",
							DiskSizeGB:  32,
							Lun:         ptr.To(int32(0)),
							CachingType: "ReadWrite",
							ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "Premium_LRS"},
						}},
					},
				},
			},
		},
		{
			// CAPZ creates the control plane machines and rejects sizes outside 4-32767,
			// only after infrastructure provisioning has started.
			name:          "control plane dataDisk below the CAPZ minimum size is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 3,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.diskSizeGB: Invalid value: 3: must be between 4 and 32767 for control plane data disks$`,
		},
		{
			name:          "control plane dataDisk above the CAPZ maximum size is invalid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 32768,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
			expected: `^test-path\.dataDisks\[0\]\.diskSizeGB: Invalid value: 32768: must be between 4 and 32767 for control plane data disks$`,
		},
		{
			name:          "control plane dataDisk at the CAPZ minimum size is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 4,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			name:          "control plane dataDisk at the CAPZ maximum size is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "master",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 32767,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			// Compute machines are created by the Machine API, not CAPZ, so the CAPZ
			// bounds do not apply.
			name:          "compute dataDisk below the CAPZ minimum size is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 1,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
		},
		{
			name:          "compute dataDisk above the CAPZ maximum size is valid",
			azurePlatform: azure.PublicCloud,
			pool: &types.MachinePool{
				Name: "worker",
				Platform: types.MachinePoolPlatform{
					Azure: &azure.MachinePool{
						DataDisks: []capz.DataDisk{{
							NameSuffix: "data",
							DiskSizeGB: 65536,
							Lun:        ptr.To(int32(0)),
						}},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			azurePlatform := &azure.Platform{CloudName: tc.azurePlatform, DefaultMachinePlatform: tc.defaultMachinePlatform}
			err := ValidateMachinePool(tc.pool.Platform.Azure, tc.pool.Name, azurePlatform, tc.pool, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

// TestDefaultedDataDiskCachingPassesUltraSSDValidation pins the invariant
// between defaulting and validation: an omitted cachingType is defaulted to
// the one caching mode UltraSSD_LRS admits, so a defaulted pool must pass
// the UltraSSD check. Defaulting and the check live in different packages
// and share the constant through defaults.CachingTypeNone; if either side
// ever moves to a different value this test fails.
func TestDefaultedDataDiskCachingPassesUltraSSDValidation(t *testing.T) {
	azurePool := &azure.MachinePool{
		UltraSSDCapability: enabled,
		DataDisks: []capz.DataDisk{{
			NameSuffix:  "etcd",
			DiskSizeGB:  32,
			Lun:         ptr.To(int32(0)),
			ManagedDisk: &capz.ManagedDiskParameters{StorageAccountType: "UltraSSD_LRS"},
		}},
	}
	defaults.SetMachinePoolDefaults(azurePool)

	err := ValidateMachinePool(azurePool, "master", &azure.Platform{CloudName: azure.PublicCloud}, nil, field.NewPath("test-path")).ToAggregate()
	assert.NoError(t, err, "the defaulted caching type must pass the UltraSSD validation")
}
