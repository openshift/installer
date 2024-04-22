package features

import (
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

type ReleaseFeatureDiffInfo struct {
	featureInfo map[string]map[string]*FeatureGateDiffInfo
}

func (i *ReleaseFeatureDiffInfo) AllFeatureInfo() []*FeatureGateDiffInfo {
	ret := []*FeatureGateDiffInfo{}
	for _, forClusterProfile := range i.featureInfo {
		for featureSet := range forClusterProfile {
			ret = append(ret, forClusterProfile[featureSet])
		}
	}
	return ret
}

func (i *ReleaseFeatureDiffInfo) AllFeatureSets() sets.Set[string] {
	allFeatureSets := sets.Set[string]{}
	for _, forClusterProfile := range i.featureInfo {
		for _, forFeatureSet := range forClusterProfile {
			allFeatureSets.Insert(forFeatureSet.FeatureSet)
		}
	}
	return allFeatureSets
}

func (i *ReleaseFeatureDiffInfo) AllClusterProfiles() sets.Set[string] {
	allClusterProfiles := sets.Set[string]{}
	for _, forClusterProfile := range i.featureInfo {
		for _, forFeatureSet := range forClusterProfile {
			allClusterProfiles.Insert(forFeatureSet.ClusterProfile)
		}
	}
	return allClusterProfiles
}

func (i *ReleaseFeatureDiffInfo) FeatureInfoFor(clusterProfile, featureSet string) *FeatureGateDiffInfo {
	forClusterProfile, ok := i.featureInfo[clusterProfile]
	if !ok {
		forClusterProfile, ok = i.featureInfo[""]
		if !ok {
			return nil
		}
	}
	forFeatureGate, ok := forClusterProfile[featureSet]
	if !ok {
		forFeatureGate, ok = forClusterProfile[""]
		if !ok {
			return nil
		}
	}
	return forFeatureGate
}

type FeatureGateDiffInfo struct {
	ClusterProfile string
	FeatureSet     string

	ChangedFeatureGates map[string]string
}

func (i *ReleaseFeatureDiffInfo) GetOrderedFeatureGates() []string {
	allInfo := i.AllFeatureInfo()
	counts := map[string]int{}
	for _, curr := range allInfo {
		for featureGate, changedFeatureGate := range curr.ChangedFeatureGates {
			if strings.HasSuffix(changedFeatureGate, "Unconditional") {
				counts[featureGate] = counts[featureGate] - 100
			}
			if strings.HasSuffix(changedFeatureGate, "(Changed)") {
				switch curr.FeatureSet {
				case "Default":
					counts[featureGate] = counts[featureGate] - 100
				}
				counts[featureGate] = counts[featureGate] - 10
			}
			if strings.HasSuffix(changedFeatureGate, "Enabled (New)") {
				switch curr.FeatureSet {
				case "Default":
					counts[featureGate] = counts[featureGate] - 10
				}
				counts[featureGate] = counts[featureGate] - 1
			}
			if strings.HasSuffix(changedFeatureGate, "Disabled (New)") {
				switch curr.FeatureSet {
				case "Default":
					counts[featureGate] = counts[featureGate] - 1
				}
				counts[featureGate] = counts[featureGate] - 1
			}
		}
	}

	toSort := []StringCount{}
	for name, count := range counts {
		toSort = append(toSort, StringCount{
			name:  name,
			count: count,
		})
	}

	sort.Sort(ByCount(toSort))
	ret := []string{}
	for _, curr := range toSort {
		ret = append(ret, curr.name)
	}

	return ret
}
