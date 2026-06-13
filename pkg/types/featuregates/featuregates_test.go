package featuregates

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	features "github.com/openshift/api/features"
)

func TestFeatureGateFromFeatureSets_NilKnownFS(t *testing.T) {
	fg := FeatureGateFromFeatureSets(nil, configv1.TechPreviewNoUpgrade, nil)
	if fg == nil {
		t.Fatal("expected non-nil FeatureGate, got nil")
	}
	// An empty gate should panic for any key since nothing is registered.
	assertPanics(t, func() {
		fg.Enabled("SomeRandomFeature")
	})
}

func TestFeatureGateFromFeatureSets_MissingFS(t *testing.T) {
	knownSets := map[configv1.FeatureSet]*features.FeatureGateEnabledDisabled{
		configv1.Default: {}, // only Default exists
	}
	fg := FeatureGateFromFeatureSets(knownSets, configv1.TechPreviewNoUpgrade, nil)
	if fg == nil {
		t.Fatal("expected non-nil FeatureGate, got nil")
	}
	assertPanics(t, func() {
		fg.Enabled("SomeRandomFeature")
	})
}

func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic but none occurred")
		}
	}()
	f()
}
