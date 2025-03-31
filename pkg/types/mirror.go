package types

import (
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	operatorv1alpha1 "github.com/openshift/api/operator/v1alpha1"
)

// Mirror holds the mirror list for a registry.
type Mirror struct {
	Location string   `json:"location"`
	Mirrors  []string `json:"mirrors,omitempty"`
}

// MirrorConfig holds the registry mirror data.
type MirrorConfig []Mirror

// BuildMirrorConfig populates a MirrorConfig from a given InstallConfig.
func BuildMirrorConfig(ic *InstallConfig) MirrorConfig {
	var mc MirrorConfig

	if len(ic.ImageDigestSources) > 0 {
		for _, src := range ic.ImageDigestSources {
			mc = append(mc, Mirror{
				Location: src.Source,
				Mirrors:  src.Mirrors,
			})
		}
	} else if len(ic.DeprecatedImageContentSources) > 0 {
		for _, src := range ic.DeprecatedImageContentSources {
			mc = append(mc, Mirror{
				Location: src.Source,
				Mirrors:  src.Mirrors,
			})
		}
	}

	return mc
}

// HasMirrors returns whether there are any mirrors configured.
func (mc MirrorConfig) HasMirrors() bool {
	return len(mc) > 0
}

// GetICSPContents converts the data in registries.conf into ICSP format.
func (mc MirrorConfig) GetICSPContents() ([]byte, error) {
	icsp := operatorv1alpha1.ImageContentSourcePolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorv1alpha1.SchemeGroupVersion.String(),
			Kind:       "ImageContentSourcePolicy",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "image-policy",
			// not namespaced
		},
	}

	icsp.Spec.RepositoryDigestMirrors = make([]operatorv1alpha1.RepositoryDigestMirrors, len(mc))
	for i, mirrorRegistries := range mc {
		icsp.Spec.RepositoryDigestMirrors[i] = operatorv1alpha1.RepositoryDigestMirrors{Source: mirrorRegistries.Location, Mirrors: mirrorRegistries.Mirrors}
	}

	// Convert to json first so json tags are handled
	jsonData, err := json.Marshal(&icsp)
	if err != nil {
		return nil, err
	}
	contents, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return nil, err
	}

	return contents, nil
}
