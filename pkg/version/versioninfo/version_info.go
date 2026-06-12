package versioninfo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/version"
)

// fallbackVersion is the version to be used if we are unable to parse
// the build version. Mostly used for dev/testing scenarios where the
// binary is not built using the build scripts.
var fallbackVersion = Info{Major: 5, Minor: 0, Patch: 0}

// Info represents a parsed semantic version.
type Info struct {
	Major int
	Minor int
	Patch int
}

func (v Info) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// GetInfo returns the build version parsed as an Info.
func GetInfo() Info {
	s, err := version.Version()
	if err != nil {
		logrus.Warnf("unable to determine version, defaulting to: %v. Error: %v", fallbackVersion, err)
		return fallbackVersion
	}
	v, err := parseInfo(s)
	if err != nil {
		logrus.Warnf("unable to parse version, defaulting to: %v. Error: %v", fallbackVersion, err)
		return fallbackVersion
	}
	return v
}

func parseInfo(s string) (Info, error) {
	if idx := strings.Index(s, "-"); idx != -1 {
		s = s[:idx]
	}
	parts := strings.Split(s, ".")
	if len(parts) == 0 {
		return fallbackVersion, fmt.Errorf("invalid version %q", s)
	}
	var v Info
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Info{}, fmt.Errorf("invalid version %q", s)
	}
	v.Major = major
	if len(parts) > 1 {
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			logrus.Warnf("failed to parse minor version from %q: %v", s, err)
		} else {
			v.Minor = minor
		}
	}
	if len(parts) > 2 {
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			logrus.Warnf("failed to parse patch version from %q: %v", s, err)
		} else {
			v.Patch = patch
		}
	}
	return v, nil
}
