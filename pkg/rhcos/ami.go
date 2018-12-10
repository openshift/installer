package rhcos

import (
	"github.com/pkg/errors"
)

// AMI fetches the HVM AMI ID of the latest Red Hat CoreOS release.
func AMI(build Build, region string) (string, error) {
	for _, ami := range build.Meta.AMIs {
		if ami.Name == region {
			return ami.HVM, nil
		}
	}

	return "", errors.Errorf("RHCOS build %s does not have AMIs in region %s", build.Meta.OSTreeVersion, region)
}
