package featuregates

// source: https://github.com/openshift/cluster-config-operator/blob/636a2dc303037e2561a243ae1ab5c5b953ddad04/pkg/cmd/render/render.go#L153

import (
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	configv1 "github.com/openshift/api/config/v1"
	features "github.com/openshift/api/features"
)

func toFeatureGateNames(in []features.FeatureGateDescription) []configv1.FeatureGateName {
	out := []configv1.FeatureGateName{}
	for _, curr := range in {
		out = append(out, curr.FeatureGateAttributes.Name)
	}

	return out
}

// completeFeatureGates identifies every known feature and ensures that is explicitly on or explicitly off.
func completeFeatureGates(knownFeatureSets map[configv1.FeatureSet]*features.FeatureGateEnabledDisabled, enabled, disabled []configv1.FeatureGateName) ([]configv1.FeatureGateName, []configv1.FeatureGateName) {
	specificallyEnabledFeatureGates := sets.New[configv1.FeatureGateName]()
	specificallyEnabledFeatureGates.Insert(enabled...)

	knownFeatureGates := sets.New[configv1.FeatureGateName]()
	knownFeatureGates.Insert(enabled...)
	knownFeatureGates.Insert(disabled...)
	for _, known := range knownFeatureSets {
		for _, curr := range known.Disabled {
			knownFeatureGates.Insert(curr.FeatureGateAttributes.Name)
		}
		for _, curr := range known.Enabled {
			knownFeatureGates.Insert(curr.FeatureGateAttributes.Name)
		}
	}

	return enabled, knownFeatureGates.Difference(specificallyEnabledFeatureGates).UnsortedList()
}

// FeatureGateFromFeatureSets creates a FeatureGate from the active feature sets.
func FeatureGateFromFeatureSets(knownFeatureSets map[configv1.FeatureSet]*features.FeatureGateEnabledDisabled, fs configv1.FeatureSet, customFS *configv1.CustomFeatureGates) FeatureGate {
	if customFS != nil {
		completeEnabled, completeDisabled := completeFeatureGates(knownFeatureSets, customFS.Enabled, customFS.Disabled)
		return newFeatureGate(completeEnabled, completeDisabled)
	}

	featureSet := knownFeatureSets[fs]

	completeEnabled, completeDisabled := completeFeatureGates(knownFeatureSets, toFeatureGateNames(featureSet.Enabled), toFeatureGateNames(featureSet.Disabled))
	return newFeatureGate(completeEnabled, completeDisabled)
}

// GenerateCustomFeatures generates the custom feature gates from the install config.
func GenerateCustomFeatures(features []string) *configv1.CustomFeatureGates {
	customFeatures := &configv1.CustomFeatureGates{}

	for _, feature := range features {
		featureName, enabled := parseCustomFeatureGate(feature)

		if enabled {
			customFeatures.Enabled = append(customFeatures.Enabled, featureName)
		} else {
			customFeatures.Disabled = append(customFeatures.Disabled, featureName)
		}
	}

	return customFeatures
}

// parseCustomFeatureGates parses the custom feature gate string into the feature name and whether it is enabled.
// The expected format is <FeatureName>=<Enabled>.
func parseCustomFeatureGate(rawFeature string) (configv1.FeatureGateName, bool) {
	var featureName string
	var enabled bool

	featureParts := strings.Split(rawFeature, "=")
	if len(featureParts) != 2 {
		return "", false
	}

	featureName = featureParts[0]

	var err error
	enabled, err = strconv.ParseBool(featureParts[1])
	if err != nil {
		return "", false
	}

	return configv1.FeatureGateName(featureName), enabled
}
