package featuregates

// source: https://github.com/openshift/library-go/blob/c515269de16e5e239bd6e93e1f9821a976bb460b/pkg/operator/configobserver/featuregates/featuregate.go#L23-L28

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
)

// FeatureGate indicates whether a given feature is enabled or not
// This interface is heavily influenced by k8s.io/component-base, but not exactly compatible.
type FeatureGate interface {
	// Enabled returns true if the key is enabled.
	Enabled(key configv1.FeatureGateName) bool
}

type featureGate struct {
	enabled  sets.Set[configv1.FeatureGateName]
	disabled sets.Set[configv1.FeatureGateName]
}

// GatedInstallConfigFeature contains fields that will be used to validate
// that required feature gates are enabled when gated install config fields
// are used.
// FeatureGateName: openshift/api feature gate required to enable the use of Field
// Condition: the check which determines whether the install config field is used,
// if Condition evaluates to True, FeatureGateName must be enabled  to pass validation.
// Field: the gated install config field, Field is used in the error message.
type GatedInstallConfigFeature struct {
	FeatureGateName configv1.FeatureGateName
	Condition       bool
	Field           *field.Path
}

func newFeatureGate(enabled, disabled []configv1.FeatureGateName) FeatureGate {
	return &featureGate{
		enabled:  sets.New[configv1.FeatureGateName](enabled...),
		disabled: sets.New[configv1.FeatureGateName](disabled...),
	}
}

// Enabled returns true if the provided feature gate is enabled
// in the active feature sets.
func (f *featureGate) Enabled(key configv1.FeatureGateName) bool {
	if f.enabled.Has(key) {
		return true
	}
	if f.disabled.Has(key) {
		return false
	}

	panic(fmt.Errorf("feature %q is not registered in FeatureGates", key))
}
