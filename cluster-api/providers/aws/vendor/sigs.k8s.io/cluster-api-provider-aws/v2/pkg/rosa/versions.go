package rosa

import (
	"fmt"
	"strings"
	"time"

	"github.com/blang/semver"
	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var MinSupportedVersion = semver.MustParse("4.14.0")

// IsVersionSupported checks whether the input version is supported for ROSA clusters.
func (c *RosaClient) IsVersionSupported(versionID string) (bool, error) {
	parsedVersion, err := semver.Parse(versionID)
	if err != nil {
		return false, err
	}
	if parsedVersion.LT(MinSupportedVersion) {
		return false, nil
	}

	filter := fmt.Sprintf("raw_id='%s' AND channel_group = '%s'", versionID, "stable")
	response, err := c.ocm.ClustersMgmt().V1().
		Versions().
		List().
		Search(filter).
		Page(1).Size(1).
		Parameter("product", "hcp").
		Send()
	if err != nil {
		return false, handleErr(response.Error(), err)
	}
	if response.Total() == 0 {
		return false, nil
	}

	version := response.Items().Get(0)
	return version.ROSAEnabled() && version.HostedControlPlaneEnabled(), nil
}

// CheckExistingScheduledUpgrade checks and returns the current upgrade schedule if any.
func (c *RosaClient) CheckExistingScheduledUpgrade(cluster *cmv1.Cluster) (*cmv1.ControlPlaneUpgradePolicy, error) {
	upgradePolicies, err := c.getControlPlaneUpgradePolicies(cluster.ID())
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
func (c *RosaClient) ScheduleControlPlaneUpgrade(cluster *cmv1.Cluster, version string, nextRun time.Time) (*cmv1.ControlPlaneUpgradePolicy, error) {
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

	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(cluster.ID()).
		ControlPlane().
		UpgradePolicies().
		Add().Body(upgradePolicy).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}

	return response.Body(), nil
}

func (c *RosaClient) getControlPlaneUpgradePolicies(clusterID string) (controlPlaneUpgradePolicies []*cmv1.ControlPlaneUpgradePolicy, err error) {
	collection := c.ocm.ClustersMgmt().V1().
		Clusters().
		Cluster(clusterID).
		ControlPlane().
		UpgradePolicies()
	page := 1
	size := 100
	for {
		response, err := collection.List().
			Page(page).
			Size(size).
			Send()
		if err != nil {
			return nil, handleErr(response.Error(), err)
		}
		controlPlaneUpgradePolicies = append(controlPlaneUpgradePolicies, response.Items().Slice()...)
		if response.Size() < size {
			break
		}
		page++
	}
	return
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

const versionPrefix = "openshift-v"

// RawVersionID returns the rawID from the provided OCM version object.
func RawVersionID(version *cmv1.Version) string {
	rawID := version.RawID()
	if rawID != "" {
		return rawID
	}

	rawID = strings.TrimPrefix(version.ID(), versionPrefix)
	channelSeparator := strings.LastIndex(rawID, "-")
	if channelSeparator > 0 {
		return rawID[:channelSeparator]
	}
	return rawID
}

// VersionID construcuts and returns an OCM versionID from the provided rawVersionID.
func VersionID(rawVersionID string) string {
	return fmt.Sprintf("%s%s", versionPrefix, rawVersionID)
}
