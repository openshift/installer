package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

func TestValidateMachinePool(t *testing.T) {
	cases := []struct {
		name          string
		azurePlatform azure.CloudEnvironment
		pool          *types.MachinePool
		expected      string
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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			azurePlatform := &azure.Platform{CloudName: tc.azurePlatform}
			err := ValidateMachinePool(tc.pool.Platform.Azure, tc.pool.Name, azurePlatform, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
