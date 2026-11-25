package bootstrap

import (
	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

// MergedMirrorSets consolidates a list of ImageDigestSources so that each
// source appears only once, with the union of its mirrors in first-seen order.
//
// Entries for the same source are merged regardless of their SourcePolicy, so
// that registries.conf never contains two entries for one location with
// conflicting blocked settings, which containers/image rejects. If any entry
// for a source is NeverContactSource the merged source is NeverContactSource;
// otherwise it is AllowContactingSource if any entry sets that explicitly.
// This matches how the Machine Config Operator merges ImageDigestMirrorSets
// (mergedDigestMirrorSets in github.com/openshift/runtime-utils/pkg/registries).
func MergedMirrorSets(sources []types.ImageDigestSource) []types.ImageDigestSource {
	sourceSet := make(map[string][]string)
	mirrorSet := make(map[string]sets.Set[string])
	policySet := make(map[string]configv1.MirrorSourcePolicy)
	orderedSources := []string{}

	for _, group := range sources {
		key := group.Source
		if _, ok := sourceSet[key]; !ok {
			orderedSources = append(orderedSources, key)
			sourceSet[key] = nil
			mirrorSet[key] = sets.New[string]()
		}
		switch group.SourcePolicy {
		case configv1.NeverContactSource:
			policySet[key] = configv1.NeverContactSource
		case configv1.AllowContactingSource:
			if policySet[key] != configv1.NeverContactSource {
				policySet[key] = configv1.AllowContactingSource
			}
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
		out = append(out, types.ImageDigestSource{Source: source, Mirrors: sourceSet[source], SourcePolicy: policySet[source]})
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
