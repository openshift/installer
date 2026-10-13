package ocm

import (
	"time"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// HypershiftUpgrader represents a Hypershift Control Plane or Node Pool Update
type HypershiftUpgrader interface {
	ID() string
	ClusterID() string
	Version() string
	State() *cmv1.UpgradePolicyState
	NextRun() time.Time
	CreationTimestamp() time.Time
	EnableMinorVersionUpgrades() bool
	Schedule() string
	ScheduleType() cmv1.ScheduleType
}

type UpgradeScheduling struct {
	ScheduleDate             string
	ScheduleTime             string
	Schedule                 string
	AllowMinorVersionUpdates bool
	AutomaticUpgrades        bool
	NextRun                  time.Time
}
