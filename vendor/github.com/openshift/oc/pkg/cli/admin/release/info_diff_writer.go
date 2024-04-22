package release

import (
	"fmt"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/library-go/pkg/features"
	"github.com/openshift/library-go/pkg/markdown"
	"k8s.io/apimachinery/pkg/util/sets"
	"strings"
)

func produceDiffMarkdown(testReporting *configv1.TestReporting, releaseFeatureDiffInfo *features.ReleaseFeatureDiffInfo) ([]byte, error) {
	allFeatureSets := releaseFeatureDiffInfo.AllFeatureSets()
	allFeatureSets.Delete("LatencySensitive") // this was a dead-end featureset we removed.
	allClusterProfiles := releaseFeatureDiffInfo.AllClusterProfiles()

	cols := []features.ColumnTuple{}
	md := markdown.NewMarkdown("FeatureGate Diff")
	md.NextTableColumn()
	md.Exact("FeatureGate ")
	for _, featureSet := range sets.List(allFeatureSets) {
		for _, clusterProfile := range sets.List(allClusterProfiles) {
			cols = append(cols, features.ColumnTuple{
				ClusterProfile: clusterProfile,
				FeatureSet:     featureSet,
			})
			md.NextTableColumn()
			md.Exact(fmt.Sprintf("%v<br/>%v ", featureSet, clusterProfile))
		}
	}
	md.EndTableRow()
	md.NextTableColumn()
	md.Exact(":------ ")
	for i := 0; i < len(cols); i++ {
		md.NextTableColumn()
		md.Exact(":---: ")
	}
	md.EndTableRow()

	orderedFeatureGates := releaseFeatureDiffInfo.GetOrderedFeatureGates()
	for _, featureGate := range orderedFeatureGates {
		md.NextTableColumn()
		md.Exact(fmt.Sprintf("%s <br/> %s", featureGate, testMarkdownForFeatureGate(testReporting, featureGate)))
		for _, col := range cols {
			change := releaseFeatureDiffInfo.FeatureInfoFor(col.ClusterProfile, col.FeatureSet).ChangedFeatureGates[featureGate]
			md.NextTableColumn()

			if change == "Disabled (New)" {
				md.Exact(" ")
			} else {
				md.Exact(strings.ReplaceAll(change, " (", "<br/>("))
			}
		}
		md.EndTableRow()
	}
	return md.ExactBytes(), nil
}

func testMarkdownForFeatureGate(testReport *configv1.TestReporting, featureGate string) string {
	var testDetails *configv1.FeatureGateTests
	for _, curr := range testReport.Spec.TestsForFeatureGates {
		if curr.FeatureGate == featureGate {
			testDetails = &curr
			break
		}
	}

	if testDetails == nil {
		return "(0 tests)"
	}

	// TODO add link to these tests on ComponentReadiness
	return fmt.Sprintf("(%d tests)", len(testDetails.Tests))
}
