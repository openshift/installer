package rosa

import (
	"time"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"
)

var MinSupportedVersion = semver.MustParse("4.14.0")

// CheckExistingScheduledUpgrade checks and returns the current upgrade schedule if any.
func CheckExistingScheduledUpgrade(client *ocm.Client, cluster *cmv1.Cluster) (*cmv1.ControlPlaneUpgradePolicy, error) {
	upgradePolicies, err := client.GetControlPlaneUpgradePolicies(cluster.ID())
	if err != nil {
		return nil, err
	}
	for _, upgradePolicy := range upgradePolicies {
		if upgradePolicy.UpgradeType() == cmv1.UpgradeTypeControlPlane {
			return upgradePolicy, nil
		}
	}
	return nil, nil
}

// ScheduleControlPlaneUpgrade schedules a new control plane upgrade to the specified version at the specified time.
func ScheduleControlPlaneUpgrade(client *ocm.Client, cluster *cmv1.Cluster, version string, nextRun time.Time) (*cmv1.ControlPlaneUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// we set it to 6 min here to account for latencty.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewControlPlaneUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeControlPlane).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		Build()
	if err != nil {
		return nil, err
	}
	return client.ScheduleHypershiftControlPlaneUpgrade(cluster.ID(), upgradePolicy)
}

// machinepools can be created with a minimal of two minor versions from the control plane.
const minorVersionsAllowedDeviation = 2

func MachinePoolSupportedVersionsRange(controlPlaneVersion string) (*semver.Version, *semver.Version, error) {
	maxVersion, err := semver.Parse(controlPlaneVersion)
	if err != nil {
		return nil, nil, err
	}

	minVersion := semver.Version{
		Major: maxVersion.Major,
		Minor: max(0, maxVersion.Minor-minorVersionsAllowedDeviation),
		Patch: 0,
	}

	if minVersion.LT(MinSupportedVersion) {
		minVersion = MinSupportedVersion
	}

	return &minVersion, &maxVersion, nil
}

// RawVersionID returns the rawID from the provided OCM version object.
func RawVersionID(version *cmv1.Version) string {
	rawID := version.RawID()
	if rawID != "" {
		return rawID
	}

	return ocm.GetRawVersionId(version.ID())
}
