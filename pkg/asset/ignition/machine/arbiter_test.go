package machine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestArbiterGenerate tests generating the arbiter asset.
func TestArbiterGenerate(t *testing.T) {
	testCases := []struct {
		expectedIgnitionConfigNames []string
		installConfig               *installconfig.InstallConfig
		description                 string
	}{
		{
			description: "should generate with arbiter config",
			expectedIgnitionConfigNames: []string{
				"arbiter.ign",
			},
			installConfig: installconfig.MakeAsset(
				&types.InstallConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cluster",
					},
					BaseDomain: "test-domain",
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
					},
					Platform: types.Platform{
						AWS: &aws.Platform{
							Region: "us-east",
						},
					},
					Arbiter: &types.MachinePool{
						Name:     "arbiter",
						Replicas: ptr.To(int64(1)),
					},
				}),
		},
		{
			description:                 "should not generate arbiter ignition when no arbiter",
			expectedIgnitionConfigNames: []string{},
			installConfig: installconfig.MakeAsset(
				&types.InstallConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cluster",
					},
					BaseDomain: "test-domain",
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
					},
					Platform: types.Platform{
						AWS: &aws.Platform{
							Region: "us-east",
						},
					},
				}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			rootCA := &tls.RootCA{}
			err := rootCA.Generate(context.Background(), nil)
			assert.NoError(t, err, "unexpected error generating root CA")

			parents := asset.Parents{}
			parents.Add(tc.installConfig, rootCA)

			arbiter := &Arbiter{}
			err = arbiter.Generate(context.Background(), parents)
			assert.NoError(t, err, "unexpected error generating arbiter asset")

			actualFiles := arbiter.Files()
			actualIgnitionConfigNames := make([]string, len(actualFiles))
			for i, f := range actualFiles {
				actualIgnitionConfigNames[i] = f.Filename
			}
			assert.Equal(t, tc.expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for arbiter ignition configs")
		})
	}
}
