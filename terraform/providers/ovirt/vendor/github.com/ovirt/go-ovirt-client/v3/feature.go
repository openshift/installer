package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// Feature is a specialized type for feature flags. These can be checked for support by using SupportsFeature in
// FeatureClient.
type Feature string

const (
	// FeatureAutoPinning is a feature flag for autopinning support in the oVirt Engine supported since 4.4.5.
	FeatureAutoPinning Feature = "autopinning"

	// FeaturePlacementPolicy is a feature flag to indicate placement policy support in the oVirt Engine.
	FeaturePlacementPolicy Feature = "placement_policy"
)

// FeatureClient provides the functions to determine the capabilities of the oVirt Engine.
type FeatureClient interface {
	// SupportsFeature checks the features supported by the oVirt Engine.
	SupportsFeature(feature Feature, retries ...RetryStrategy) (bool, error)
}

func (o *oVirtClient) SupportsFeature(feature Feature, retries ...RetryStrategy) (result bool, err error) {
	var minimumVersion *ovirtsdk.Version
	switch feature {
	case FeatureAutoPinning:
		minimumVersion = ovirtsdk.NewVersionBuilder().
			Major(4).
			Minor(4).
			Build_(5).
			Revision(0).
			MustBuild()
	case FeaturePlacementPolicy:
		minimumVersion = ovirtsdk.NewVersionBuilder().
			Major(4).
			Minor(4).
			Build_(5).
			Revision(0).
			MustBuild()
	default:
		return false, newError(EBug, "unknown feature: %s", feature)
	}

	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		"fetching engine version",
		o.logger,
		retries,
		func() error {
			systemGetResponse, err := o.conn.SystemService().Get().Send()
			if err != nil {
				return err
			}
			engineVer := systemGetResponse.MustApi().MustProductInfo().MustVersion()
			versionCompareResult := versionCompare(engineVer, minimumVersion)
			result = versionCompareResult >= 0
			return nil
		})
	return result, err
}

func versionCompare(v *ovirtsdk.Version, other *ovirtsdk.Version) int64 {
	if result := v.MustMajor() - other.MustMajor(); result != 0 {
		return result
	}
	if result := v.MustMinor() - other.MustMinor(); result != 0 {
		return result
	}

	if result := v.MustBuild() - other.MustBuild(); result != 0 {
		return result
	}
	return v.MustRevision() - other.MustRevision()
}

func (m *mockClient) SupportsFeature(_ Feature, _ ...RetryStrategy) (bool, error) {
	return true, nil
}
