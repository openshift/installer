package bootstrap

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
)

// MergedMirrorSets consolidates a list of ImageContentSources so that each
// source appears only once.
func MergedMirrorSets(sources []types.ImageContentSource) []types.ImageContentSource {
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

	out := []types.ImageContentSource{}
	for _, source := range orderedSources {
		out = append(out, types.ImageContentSource{Source: source, Mirrors: sourceSet[source]})
	}
	return out
}
