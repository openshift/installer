//nolint:revive // package name matches the directory; renaming would break test access to internals
package tls

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/manifests"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

func TestIRIRegistryCredentialsGenerate(t *testing.T) {
	tests := []struct {
		name           string
		iriManifest    bool
		shouldGenerate bool
	}{
		{
			name:           "Generate with IRI manifest present",
			iriManifest:    true,
			shouldGenerate: true,
		},
		{
			name:           "Skip without IRI manifest",
			iriManifest:    false,
			shouldGenerate: false,
		},
	}

	// The pull secret the merged form of which the asset carries: an entry
	// for the IRI registry is added to it when credentials are generated,
	// and it is carried unchanged when they are not.
	existingPullSecret := `{"auths":{"quay.io":{"auth":"dGVzdDp0ZXN0"}}}` //nolint:gosec // test fixture, not a real credential

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create install config with feature gate
			ic := &installconfig.InstallConfig{
				AssetBase: installconfig.AssetBase{
					Config: &types.InstallConfig{
						TypeMeta: metav1.TypeMeta{
							APIVersion: types.InstallConfigVersion,
						},
						ObjectMeta: metav1.ObjectMeta{
							Name: "test-cluster",
						},
						BaseDomain: "example.com",
						Networking: &types.Networking{
							MachineNetwork: []types.MachineNetworkEntry{
								{CIDR: *ipnet.MustParseCIDR("10.0.0.0/16")},
							},
						},
						ControlPlane: &types.MachinePool{
							Name:     "master",
							Replicas: pointer(int64(3)),
						},
						Compute: []types.MachinePool{
							{
								Name:     "worker",
								Replicas: pointer(int64(3)),
							},
						},
						Platform: types.Platform{
							BareMetal: &baremetal.Platform{
								APIVIPs: []string{"192.168.111.5"},
							},
						},
						PullSecret: existingPullSecret,
					},
				},
			}

			// Create IRI manifest asset
			iri := &manifests.InternalReleaseImage{}
			if tt.iriManifest {
				iri.FileList = []*asset.File{
					{Filename: "manifests/internal-release-image.yaml"},
				}
			}

			// Create IRIRegistryCredentials asset and generate
			auth := &IRIRegistryCredentials{}
			parents := asset.Parents{}
			parents.Add(ic, iri)

			err := auth.Generate(context.Background(), parents)
			if !assert.NoError(t, err) {
				return
			}

			if !tt.shouldGenerate {
				assert.Empty(t, auth.Password, "Password should be empty when generation is skipped")
				assert.Empty(t, auth.HtpasswdContent, "HtpasswdContent should be empty when generation is skipped")
				assert.Equal(t, existingPullSecret, auth.MergedPullSecret, "MergedPullSecret should be the pull secret unchanged when generation is skipped")
				return
			}

			// Verify password was generated
			assert.NotEmpty(t, auth.Password, "Password should not be empty")
			assert.Equal(t, IRIRegistryUsername, auth.Username, "Username should be 'openshift'")

			// Verify password is base64-encoded 32 bytes (256-bit)
			passwordBytes, err := base64.StdEncoding.DecodeString(auth.Password)
			if !assert.NoError(t, err, "Password should be valid base64") {
				return
			}
			assert.Equal(t, PasswordBytes, len(passwordBytes), "Password should be 32 bytes before encoding")

			// Verify htpasswd content format: "openshift:$2y$10$..."
			assert.True(t, strings.HasPrefix(auth.HtpasswdContent, "openshift:$2"), "Htpasswd should start with 'openshift:$2'")
			assert.True(t, strings.HasSuffix(auth.HtpasswdContent, "\n"), "Htpasswd should end with newline")

			// Extract bcrypt hash from htpasswd content
			parts := strings.Split(strings.TrimSpace(auth.HtpasswdContent), ":")
			if !assert.Equal(t, 2, len(parts), "Htpasswd should have format 'username:hash'") {
				return
			}
			hash := parts[1]

			// Verify bcrypt hash validates against password
			err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(auth.Password))
			assert.NoError(t, err, "Bcrypt hash should validate against password")

			// Verify the merged pull secret carries the credentials for the
			// registry behind api-int on the cluster domain, and the entry
			// the install config already carried.
			var merged map[string]interface{}
			if !assert.NoError(t, json.Unmarshal([]byte(auth.MergedPullSecret), &merged)) {
				return
			}
			auths, ok := merged["auths"].(map[string]interface{})
			if !assert.True(t, ok, "MergedPullSecret should carry auths") {
				return
			}
			iriAuth, ok := auths["api-int.test-cluster.example.com:22625"].(map[string]interface{})
			if assert.True(t, ok, "MergedPullSecret should carry the IRI registry auth") {
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(IRIRegistryUsername+":"+auth.Password)), iriAuth["auth"])
			}
			if quayAuth, ok := auths["quay.io"].(map[string]interface{}); ok {
				assert.Equal(t, "dGVzdDp0ZXN0", quayAuth["auth"], "the existing auth should survive the merge")
			} else {
				t.Error("the existing quay.io auth should survive the merge")
			}
		})
	}
}

func TestIRIRegistryCredentialsName(t *testing.T) {
	auth := &IRIRegistryCredentials{}
	assert.Equal(t, "IRI Registry Authentication", auth.Name())
}

// pointer returns a pointer to the given value.
func pointer(i int64) *int64 {
	return &i
}
