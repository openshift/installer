package manifests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content/bootkube"
	contentmanifests "github.com/openshift/installer/pkg/asset/templates/content/manifests"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// TestGenerateBootKubeManifestsPullSecretIRIAuth verifies that the global
// pull secret manifest carries the IRI registry credentials when they were
// generated, and stays as it is when they were not.
func TestGenerateBootKubeManifestsPullSecretIRIAuth(t *testing.T) {
	// The bootkube templates are read through data.Assets, which in
	// development builds is the data directory relative to the working
	// directory, so the test has to run from the data directory itself.
	t.Chdir("../../../data")

	const testPassword = "test-iri-password"

	newParents := func(withIRIAuth bool, pullSecret string) asset.Parents {
		installConfig := &installconfig.InstallConfig{
			AssetBase: installconfig.AssetBase{
				Config: &types.InstallConfig{
					TypeMeta: metav1.TypeMeta{
						APIVersion: types.InstallConfigVersion,
					},
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-cluster",
					},
					BaseDomain: "example.com",
					Platform: types.Platform{
						BareMetal: &baremetal.Platform{
							APIVIPs: []string{"192.168.111.5"},
						},
					},
					PullSecret: pullSecret,
				},
			},
		}

		pullSecretAsset := &bootkube.OpenshiftConfigSecretPullSecret{}
		if err := pullSecretAsset.Generate(context.Background(), asset.Parents{}); err != nil {
			t.Fatalf("failed to generate the pull secret asset: %v", err)
		}

		// The credentials asset is always in the graph; without an internal
		// release image it merely carries no credentials.
		iriAuth := &tls.IRIRegistryCredentials{}
		if withIRIAuth {
			iriAuth.Username = tls.IRIRegistryUsername
			iriAuth.Password = testPassword
		}

		parents := asset.Parents{}
		parents.Add(
			installConfig,
			pullSecretAsset,
			iriAuth,
			// Every asset generateBootKubeManifests reads has to be in the
			// parents; the ones this test does not care about stay empty.
			&installconfig.ClusterID{},
			&tls.MCSCertKey{},
			&tls.RootCA{},
			&contentmanifests.InternalReleaseImage{},
			&bootkube.CVOOverrides{},
			&bootkube.KubeCloudConfig{},
			&bootkube.MachineConfigServerCASecret{},
			&bootkube.MachineConfigServerCAConfigMap{},
			&bootkube.KubeSystemConfigmapRootCA{},
			&bootkube.MachineConfigServerTLSSecret{},
		)
		return parents
	}

	existingPullSecret := `{"auths":{"quay.io":{"auth":"` +
		base64.StdEncoding.EncodeToString([]byte("existing:secret")) + `"}}}`

	tests := []struct {
		name        string
		withIRIAuth bool
	}{
		{name: "with IRI registry credentials", withIRIAuth: true},
		{name: "without IRI registry credentials", withIRIAuth: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manifests{}
			files, err := m.generateBootKubeManifests(newParents(tt.withIRIAuth, existingPullSecret))
			if !assert.NoError(t, err) {
				return
			}

			var pullSecretFile *asset.File
			for _, f := range files {
				if f.Filename == "manifests/openshift-config-secret-pull-secret.yaml" {
					pullSecretFile = f
					break
				}
			}
			if !assert.NotNil(t, pullSecretFile, "the global pull secret manifest should be generated") {
				return
			}

			var secret corev1.Secret
			if !assert.NoError(t, yaml.Unmarshal(pullSecretFile.Data, &secret), "the manifest should parse as a Secret") {
				return
			}
			decoded, ok := secret.Data[".dockerconfigjson"]
			if !assert.True(t, ok, "the pull secret manifest should carry the docker config json") {
				return
			}

			var dockerConfig map[string]interface{}
			if !assert.NoError(t, json.Unmarshal(decoded, &dockerConfig)) {
				return
			}
			auths, ok := dockerConfig["auths"].(map[string]interface{})
			if !assert.True(t, ok, "the docker config json should carry auths") {
				return
			}

			registryHost := fmt.Sprintf("api-int.test-cluster.example.com:%d", tls.IRIRegistryPort)
			if !tt.withIRIAuth {
				assert.NotContains(t, auths, registryHost, "the pull secret should not carry the IRI registry auth")
				return
			}

			iriAuth, ok := auths[registryHost].(map[string]interface{})
			if !assert.True(t, ok, "the pull secret should carry the IRI registry auth") {
				return
			}
			expected := base64.StdEncoding.EncodeToString([]byte(tls.IRIRegistryUsername + ":" + testPassword))
			assert.Equal(t, expected, iriAuth["auth"], "the IRI registry auth should carry the credentials")

			if quayAuth, ok := auths["quay.io"].(map[string]interface{}); ok {
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("existing:secret")), quayAuth["auth"],
					"the existing auth should survive the merge")
			} else {
				t.Error("the existing quay.io auth should survive the merge")
			}
		})
	}
}
