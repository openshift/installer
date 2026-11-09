package rosa

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift/rosa/pkg/ocm"
)

// MinSupportedVersion is the minimum supported version for ROSA.
var MinSupportedVersion = semver.MustParse("4.14.0")

// CheckExistingScheduledUpgrade checks and returns the current upgrade schedule if any.
func CheckExistingScheduledUpgrade(client OCMClient, cluster *cmv1.Cluster) (*cmv1.ControlPlaneUpgradePolicy, error) {
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
func ScheduleControlPlaneUpgrade(client OCMClient, cluster *cmv1.Cluster, version string, nextRun time.Time, ack bool) (*cmv1.ControlPlaneUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// Set our next run request to something slightly longer than 5min to make sure we account for the latency between when we send this
	// request and when the server processes it.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewControlPlaneUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeControlPlane).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		EnableMinorVersionUpgrades(true).
		Build()
	if err != nil {
		return nil, err
	}

	versionGates, err := client.GetMissingGateAgreementsHypershift(cluster.ID(), upgradePolicy)
	if err != nil {
		return nil, err
	}

	if !ack && len(versionGates) > 0 {
		errMess := "version gate acknowledgment required"
		for id := range versionGates {
			errMess = fmt.Sprintf(errMess+"\nid:%s\n %s\n %s\n %s\n", versionGates[id].ID(), versionGates[id].Description(), versionGates[id].DocumentationURL(), versionGates[id].WarningMessage())
		}

		return nil, fmt.Errorf("%s", errMess)
	}

	for id := range versionGates {
		if err = client.AckVersionGate(cluster.ID(), versionGates[id].ID()); err != nil {
			return nil, err
		}
	}

	return client.ScheduleHypershiftControlPlaneUpgrade(cluster.ID(), upgradePolicy)
}

// ScheduleNodePoolUpgrade schedules a new nodePool upgrade to the specified version at the specified time.
func ScheduleNodePoolUpgrade(client OCMClient, clusterID string, nodePool *cmv1.NodePool, version string, nextRun time.Time) (*cmv1.NodePoolUpgradePolicy, error) {
	// earliestNextRun is set to at least 5 min from now by the OCM API.
	// Set our next run request to something slightly longer than 5min to make sure we account for the latency between when we send this
	// request and when the server processes it.
	earliestNextRun := time.Now().Add(time.Minute * 6)
	if nextRun.Before(earliestNextRun) {
		nextRun = earliestNextRun
	}

	upgradePolicy, err := cmv1.NewNodePoolUpgradePolicy().
		UpgradeType(cmv1.UpgradeTypeNodePool).
		NodePoolID(nodePool.ID()).
		ScheduleType(cmv1.ScheduleTypeManual).
		Version(version).
		NextRun(nextRun).
		EnableMinorVersionUpgrades(true).
		Build()
	if err != nil {
		return nil, err
	}

	scheduledUpgrade, err := client.ScheduleNodePoolUpgrade(clusterID, nodePool.ID(), upgradePolicy)
	if err != nil {
		return nil, fmt.Errorf("failed to schedule nodePool upgrade to version %s: %w", version, err)
	}

	return scheduledUpgrade, nil
}

// machinepools can be created with a minimal of two minor versions from the control plane.
const minorVersionsAllowedDeviation = 2

// CoreVersion strips any prerelease and build qualifiers, leaving major.minor.patch.
// The machine pool skew policy is expressed purely in terms of minor versions, and
// semver sorts a prerelease below its own release ("5.0.0-rc.0" < "5.0.0"), so the
// qualifiers have to be dropped before the lower-bound/skew comparison -- otherwise a
// prerelease control plane falls outside the range derived from itself.
func CoreVersion(version semver.Version) semver.Version {
	return semver.Version{Major: version.Major, Minor: version.Minor, Patch: version.Patch}
}

// MachinePoolSupportedVersionsRange returns the supported version range for a
// machine pool given the control plane version. Cross-major lower bound is
// deliberately permissive -- OCM is authoritative for unsupported versions.
func MachinePoolSupportedVersionsRange(controlPlaneVersion string) (*semver.Version, *semver.Version, error) {
	parsed, err := semver.Parse(controlPlaneVersion)
	if err != nil {
		return nil, nil, err
	}

	// Preserve prerelease in the upper bound so a GA pool is rejected against an RC CP.
	maxVersion := parsed

	// Strip prerelease for the lower bound so an RC CP falls within its own range.
	coreVersion := CoreVersion(parsed)

	// Minor is uint64 -- subtraction underflows instead of clamping. Use min()
	// on the subtrahend. Cross-major case admits the whole previous major series.
	var minVersion semver.Version
	if coreVersion.Minor < minorVersionsAllowedDeviation && coreVersion.Major > 0 {
		minVersion = semver.Version{Major: coreVersion.Major - 1, Minor: 0, Patch: 0}
	} else {
		minVersion = semver.Version{
			Major: coreVersion.Major,
			Minor: coreVersion.Minor - min(coreVersion.Minor, minorVersionsAllowedDeviation),
			Patch: 0,
		}
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
