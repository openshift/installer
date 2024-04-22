package features

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
	kyaml "sigs.k8s.io/yaml"
)

type ReleaseFeatureInfo struct {
	featureInfo map[string]map[string]*FeatureGateInfo
}

func (i *ReleaseFeatureInfo) AllFeatureInfo() []*FeatureGateInfo {
	ret := []*FeatureGateInfo{}
	for _, forClusterProfile := range i.featureInfo {
		for featureSet := range forClusterProfile {
			ret = append(ret, forClusterProfile[featureSet])
		}
	}
	return ret
}

func (i *ReleaseFeatureInfo) FeatureInfoFor(clusterProfile, featureSet string) *FeatureGateInfo {
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

func (i *ReleaseFeatureInfo) AllFeatureSets() sets.Set[string] {
	allFeatureSets := sets.Set[string]{}
	for _, forClusterProfile := range i.featureInfo {
		for _, forFeatureSet := range forClusterProfile {
			allFeatureSets.Insert(forFeatureSet.FeatureSet)
		}
	}
	return allFeatureSets
}

func (i *ReleaseFeatureInfo) AllClusterProfiles() sets.Set[string] {
	allClusterProfiles := sets.Set[string]{}
	for _, forClusterProfile := range i.featureInfo {
		for _, forFeatureSet := range forClusterProfile {
			allClusterProfiles.Insert(forFeatureSet.ClusterProfile)
		}
	}
	return allClusterProfiles
}

func (i *ReleaseFeatureInfo) AllFeatureGates() sets.Set[string] {
	allFeatureGates := sets.Set[string]{}
	for _, forClusterProfile := range i.featureInfo {
		for _, forFeatureSet := range forClusterProfile {
			allFeatureGates.Insert(sets.List(forFeatureSet.Enabled)...)
			allFeatureGates.Insert(sets.List(forFeatureSet.Disabled)...)
		}
	}
	return allFeatureGates
}

type FeatureGateInfo struct {
	ClusterProfile string
	FeatureSet     string

	Enabled         sets.Set[string]
	Disabled        sets.Set[string]
	AllFeatureGates map[string]bool
}

type FilenameToContent map[string][]byte

func (i *ReleaseFeatureInfo) CalculateDiff(ctx context.Context, from *ReleaseFeatureInfo) *ReleaseFeatureDiffInfo {
	to := i
	allFeatureSets := from.AllFeatureSets()
	allFeatureSets.Insert(sets.List(to.AllFeatureSets())...)
	allClusterProfiles := from.AllClusterProfiles()
	allClusterProfiles.Insert(sets.List(to.AllClusterProfiles())...)

	releaseFeatureDiffInfo := &ReleaseFeatureDiffInfo{
		featureInfo: map[string]map[string]*FeatureGateDiffInfo{},
	}
	changedFeatureGates := getChangedFeatureGates(from, to)
	for _, clusterProfile := range sets.List(allClusterProfiles) {
		if len(clusterProfile) == 0 {
			continue
		}
		releaseFeatureDiffInfo.featureInfo[clusterProfile] = map[string]*FeatureGateDiffInfo{}

		for _, featureSet := range sets.List(allFeatureSets) {
			currDiffInfo := &FeatureGateDiffInfo{
				ClusterProfile:      clusterProfile,
				FeatureSet:          featureSet,
				ChangedFeatureGates: map[string]string{},
			}
			for _, featureGate := range sets.List(changedFeatureGates) {
				fromFeatureInfo := from.FeatureInfoFor(clusterProfile, featureSet)
				toFeatureInfo := to.FeatureInfoFor(clusterProfile, featureSet)
				switch {
				case toFeatureInfo == nil && fromFeatureInfo != nil:
					currDiffInfo.ChangedFeatureGates[featureGate] = "Not Available (Changed)"
					continue
				case toFeatureInfo == nil && fromFeatureInfo == nil:
					currDiffInfo.ChangedFeatureGates[featureGate] = "Not Available"
					continue
				case toFeatureInfo != nil && fromFeatureInfo == nil:
					toEnabled, toOk := toFeatureInfo.AllFeatureGates[featureGate]
					switch {
					case !toOk:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Unconditional (New)"
					case toEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Enabled (New)"
					case !toEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Disabled (New)"
					}
					continue
				case toFeatureInfo != nil && fromFeatureInfo != nil:
				}

				fromEnabled, fromOk := fromFeatureInfo.AllFeatureGates[featureGate]
				toEnabled, toOk := toFeatureInfo.AllFeatureGates[featureGate]
				switch {
				case toOk && !fromOk:
					switch {
					case toEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Enabled (New)"
					case !toEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Disabled (New)"
					}
				case toOk && fromOk:
					switch {
					case toEnabled && !fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Enabled (Changed)"
					case toEnabled && fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Enabled"
					case !toEnabled && fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Disabled (Changed)"
						changedFeatureGates.Insert(featureGate)
					case !toEnabled && !fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Disabled"
					}
				case !toOk && fromOk:
					switch {
					case fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Unconditional (Changed)"
					case !fromEnabled:
						currDiffInfo.ChangedFeatureGates[featureGate] = "Unconditional (Changed)"
					}
				case !toOk && !fromOk:
					currDiffInfo.ChangedFeatureGates[featureGate] = "Unconditional"
				}
			}
			releaseFeatureDiffInfo.featureInfo[clusterProfile][featureSet] = currDiffInfo
		}
	}

	return releaseFeatureDiffInfo
}

func getChangedFeatureGates(from, to *ReleaseFeatureInfo) sets.Set[string] {
	allFeatureSets := from.AllFeatureSets()
	allFeatureSets.Insert(sets.List(to.AllFeatureSets())...)
	allClusterProfiles := from.AllClusterProfiles()
	allClusterProfiles.Insert(sets.List(to.AllClusterProfiles())...)
	allFeatureGates := from.AllFeatureGates()
	allFeatureGates.Insert(sets.List(to.AllFeatureGates())...)

	changedFeatureGates := sets.Set[string]{}
	for _, featureGate := range sets.List(allFeatureGates) {
		for _, clusterProfile := range sets.List(allClusterProfiles) {
			for _, featureSet := range sets.List(allFeatureSets) {
				fromFeatureInfo := from.FeatureInfoFor(clusterProfile, featureSet)
				if fromFeatureInfo == nil {
					continue
				}
				toFeatureInfo := to.FeatureInfoFor(clusterProfile, featureSet)
				if toFeatureInfo == nil {
					continue
				}
				fromEnabled, fromOk := fromFeatureInfo.AllFeatureGates[featureGate]
				toEnabled, toOk := toFeatureInfo.AllFeatureGates[featureGate]
				switch {
				case toOk && !fromOk:
					changedFeatureGates.Insert(featureGate)
				case toOk && fromOk:
					switch {
					case toEnabled && !fromEnabled:
						changedFeatureGates.Insert(featureGate)
					case !toEnabled && fromEnabled:
						changedFeatureGates.Insert(featureGate)
					case toEnabled && fromEnabled:
					case !toEnabled && !fromEnabled:
					}
				case !toOk && fromOk:
					changedFeatureGates.Insert(featureGate)
				case !toOk && !fromOk:
				}
			}
		}
	}

	return changedFeatureGates
}

func (i *ReleaseFeatureInfo) GetOrderedFeatureGates() []string {
	allInfo := i.AllFeatureInfo()
	counts := map[string]int{}
	for _, curr := range allInfo {
		for _, featureGate := range sets.List(curr.Enabled) {
			counts[featureGate] = counts[featureGate] + 1
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

type StringCount struct {
	name  string
	count int
}
type ByCount []StringCount

func (a ByCount) Len() int      { return len(a) }
func (a ByCount) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCount) Less(i, j int) bool {
	if a[i].count < a[j].count {
		return true
	}
	if a[i].count > a[j].count {
		return false
	}
	if strings.Compare(a[i].name, a[j].name) < 0 {
		return true
	}
	return false
}

type ColumnTuple struct {
	ClusterProfile string
	FeatureSet     string
}

func ReadReleaseFeatureInfoFromDir(ctx context.Context, dir string) (*ReleaseFeatureInfo, error) {
	files := FilenameToContent{}
	featureSetManifestFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read FeatureSetManifestDir: %w", err)
	}
	for _, currFeatureSetManifestFile := range featureSetManifestFiles {
		if !strings.Contains(currFeatureSetManifestFile.Name(), "featureGate-") {
			continue
		}
		featureGateFilename := filepath.Join(dir, currFeatureSetManifestFile.Name())
		featureGateBytes, err := os.ReadFile(featureGateFilename)
		if err != nil {
			return nil, fmt.Errorf("unable to read %q: %w", featureGateFilename, err)
		}
		files[currFeatureSetManifestFile.Name()] = featureGateBytes
	}

	return ReadReleaseFeatureInfo(ctx, files)
}

func ReadReleaseFeatureInfo(ctx context.Context, files FilenameToContent) (*ReleaseFeatureInfo, error) {
	ret := &ReleaseFeatureInfo{
		featureInfo: map[string]map[string]*FeatureGateInfo{},
	}

	for featureGateFilename, featureGateBytes := range files {
		if !strings.Contains(featureGateFilename, "featureGate-") {
			continue
		}

		// use unstructured to pull this information to avoid vendoring openshift/api
		featureGateMap := map[string]interface{}{}
		if err := kyaml.Unmarshal(featureGateBytes, &featureGateMap); err != nil {
			return nil, fmt.Errorf("unable to parse featuregate %q: %w", featureGateFilename, err)
		}
		uncastFeatureGate := unstructured.Unstructured{
			Object: featureGateMap,
		}

		clusterProfiles := clusterOperatorClusterProfilesFrom(uncastFeatureGate.GetAnnotations())
		if len(clusterProfiles) > 1 {
			return nil, fmt.Errorf("expected at most one clusterProfile from %q: %v", featureGateFilename, sets.List(clusterProfiles))
		}
		shortClusterProfile := ""
		if len(clusterProfiles) > 0 {
			shortClusterProfile = ClusterProfileToShortName(sets.List(clusterProfiles)[0])
		}

		currFeatureSet, _, _ := unstructured.NestedString(uncastFeatureGate.Object, "spec", "featureSet")
		if len(currFeatureSet) == 0 {
			currFeatureSet = "Default"
		}
		if currFeatureSet == "CustomNoUpgrade" {
			continue
		}

		uncastFeatureGateSlice, _, err := unstructured.NestedSlice(uncastFeatureGate.Object, "status", "featureGates")
		if err != nil {
			return nil, fmt.Errorf("no slice found %w", err)
		}

		featureGateValues := map[string]bool{}
		enabledGates := sets.Set[string]{}
		disabledGates := sets.Set[string]{}
		enabledFeatureGates, _, err := unstructured.NestedSlice(uncastFeatureGateSlice[0].(map[string]interface{}), "enabled")
		if err != nil {
			return nil, fmt.Errorf("no enabled found %w", err)
		}
		for _, currGate := range enabledFeatureGates {
			featureGateName, _, err := unstructured.NestedString(currGate.(map[string]interface{}), "name")
			if err != nil {
				return nil, fmt.Errorf("no gate name found %w", err)
			}
			enabledGates.Insert(featureGateName)
			featureGateValues[featureGateName] = true
		}

		disabledFeatureGates, _, err := unstructured.NestedSlice(uncastFeatureGateSlice[0].(map[string]interface{}), "disabled")
		if err != nil {
			return nil, fmt.Errorf("no enabled found %w", err)
		}
		for _, currGate := range disabledFeatureGates {
			featureGateName, _, err := unstructured.NestedString(currGate.(map[string]interface{}), "name")
			if err != nil {
				return nil, fmt.Errorf("no gate name found %w", err)
			}
			disabledGates.Insert(featureGateName)
			featureGateValues[featureGateName] = false
		}

		featureGateInfo := &FeatureGateInfo{
			ClusterProfile:  shortClusterProfile,
			FeatureSet:      currFeatureSet,
			Enabled:         enabledGates,
			Disabled:        disabledGates,
			AllFeatureGates: featureGateValues,
		}
		existing, ok := ret.featureInfo[shortClusterProfile]
		if !ok {
			existing = map[string]*FeatureGateInfo{}
			ret.featureInfo[shortClusterProfile] = existing
		}
		existing[currFeatureSet] = featureGateInfo
		ret.featureInfo[currFeatureSet] = existing
	}

	return ret, nil
}

func clusterOperatorClusterProfilesFrom(annotations map[string]string) sets.Set[string] {
	ret := sets.Set[string]{}
	for k, v := range annotations {
		if strings.HasPrefix(k, "include.release.openshift.io/") && v != "false" {
			ret.Insert(k)
		}
	}
	return ret
}

var (
	clusterProfileToShortName = map[string]string{
		"include.release.openshift.io/ibm-cloud-managed":              "Hypershift",
		"include.release.openshift.io/self-managed-high-availability": "SelfManagedHA",
		"include.release.openshift.io/single-node-developer":          "SingleNode",
	}
)

func ClusterProfileToShortName(annotation string) string {
	return clusterProfileToShortName[annotation]
}
