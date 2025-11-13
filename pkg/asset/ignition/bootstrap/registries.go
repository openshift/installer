package bootstrap

import (
	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

type SourceSetKey struct {
	Source       string
	SourcePolicy configv1.MirrorSourcePolicy
}

// MergedMirrorSets consolidates a list of ImageDigestSources so that each
// source appears only once.
func MergedMirrorSets(sources []types.ImageDigestSource) []types.ImageDigestSource {
	sourceSet := make(map[SourceSetKey][]string)
	mirrorSet := make(map[SourceSetKey]sets.String)
	orderedSources := []SourceSetKey{}

	for _, group := range sources {
		key := SourceSetKey{Source: group.Source, SourcePolicy: group.SourcePolicy}
		if _, ok := sourceSet[key]; !ok {
			orderedSources = append(orderedSources, key)
			sourceSet[key] = nil
			mirrorSet[key] = sets.NewString()
		}
		for _, mirror := range group.Mirrors {
			if !mirrorSet[key].Has(mirror) {
				sourceSet[key] = append(sourceSet[key], mirror)
				mirrorSet[key].Insert(mirror)
			}
		}
	}

	out := []types.ImageDigestSource{}
	for _, source := range orderedSources {
		out = append(out, types.ImageDigestSource{Source: source.Source, Mirrors: sourceSet[source], SourcePolicy: source.SourcePolicy})
	}
	return out
}

// ContentSourceToDigestMirror creates the ImageContentSource to ImageDigestSource struct
// ImageContentSource is deprecated, use ImageDigestSource.
func ContentSourceToDigestMirror(sources []types.ImageContentSource) []types.ImageDigestSource {
	digestSources := []types.ImageDigestSource{}
	for _, s := range sources {
		digestSources = append(digestSources, types.ImageDigestSource{Source: s.Source, Mirrors: s.Mirrors})
	}
	return digestSources
}
