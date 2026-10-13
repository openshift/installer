package common

import (
	"github.com/hashicorp/go-version"
)

func VersionGreaterOrEqual(version1, version2 string) (bool, error) {
	v1, err := version.NewVersion(version1)
	if err != nil {
		return false, err
	}
	v2, err := version.NewVersion(version2)
	if err != nil {
		return false, err
	}
	return v1.GreaterThanOrEqual(v2), nil
}
