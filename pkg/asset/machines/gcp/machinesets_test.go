package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	machineapi "github.com/openshift/api/machine/v1beta1"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func TestMachineSetsBootImages(t *testing.T) {
	cases := []struct {
		name        string
		customImage *gcptypes.OSImage
		wantImage   string
	}{
		{name: "default image"},
		{
			name: "custom image",
			customImage: &gcptypes.OSImage{
				Project: "custom-project",
				Name:    "custom-image",
			},
			wantImage: "projects/custom-project/global/images/custom-image",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config, pool := gcpImageFixture("worker")
			pool.Platform.GCP.OSImage = tc.customImage

			machineSets, err := MachineSets("cluster", config, pool, "worker", "worker-user-data")
			if !assert.NoError(t, err) || !assert.Len(t, machineSets, 2) {
				return
			}
			for _, machineSet := range machineSets {
				provider := machineSet.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
				assert.Equal(t, tc.wantImage, provider.Disks[0].Image)
			}
		})
	}
}
