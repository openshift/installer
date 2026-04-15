// Package openstack generates Machine objects for openstack.
package openstack

import (
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/types/openstack"
)

// ResolveBootstrapFlavor determines the effective flavor to use for the
// bootstrap instance. If the platform explicitly sets BootstrapFlavor, that
// value is returned. Otherwise the control plane flavor is used as a fallback.
func ResolveBootstrapFlavor(platform *openstack.Platform, controlPlaneFlavor string) string {
	if platform != nil && platform.BootstrapFlavor != "" {
		logrus.Infof("Using explicitly configured bootstrap flavor: %q", platform.BootstrapFlavor)
		return platform.BootstrapFlavor
	}

	logrus.Infof("Bootstrap flavor not set; inheriting control plane flavor: %q", controlPlaneFlavor)
	return controlPlaneFlavor
}
