package bootstrap

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
)

// MergedMirrorSets consolidates a list of ImageDigestSources so that each
// source appears only once.
func MergedMirrorSets(sources []types.ImageDigestSource) []types.ImageDigestSource {
	sourceSet := make(map[string][]string)
	mirrorSet := make(map[string]sets.String)
	orderedSources := []string{}

	for _, group := range sources {
		if _, ok := sourceSet[group.Source]; !ok {
			orderedSources = append(orderedSources, group.Source)
			sourceSet[group.Source] = nil
			mirrorSet[group.Source] = sets.NewString()
		}
		for _, mirror := range group.Mirrors {
			if !mirrorSet[group.Source].Has(mirror) {
				sourceSet[group.Source] = append(sourceSet[group.Source], mirror)
				mirrorSet[group.Source].Insert(mirror)
			}
		}
	}

	out := []types.ImageDigestSource{}
	for _, source := range orderedSources {
		out = append(out, types.ImageDigestSource{Source: source, Mirrors: sourceSet[source]})
	}
	return out
}

// ContentSourceToDigestMirror creates the ImageContentSource to ImageDigestSource struct
// ImageContentSource is deprecated, use ImageDigestSource.
func ContentSourceToDigestMirror(sources []types.ImageContentSource) []types.ImageDigestSource {
	digestSources := []types.ImageDigestSource{}
	for _, s := range sources {
		digestSources = append(digestSources, types.ImageDigestSource(s))
	}
	return digestSources
}
