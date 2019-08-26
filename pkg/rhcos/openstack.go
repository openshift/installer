package rhcos

import (
	"context"
	"math/rand"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// OpenStack fetches the URL of the Red Hat Enterprise Linux CoreOS release,
// for the openstack platform
func OpenStack(ctx context.Context) (string, error) {
	meta, err := fetchRHCOSBuild(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	base, err := url.Parse(meta.BaseURI)
	if err != nil {
		return "", err
	}

	relOpenStack, err := url.Parse(meta.Images.OpenStack.Path)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(relOpenStack).String(), nil
}

// OpenStackGlanceImageName generates the name for RHCOS image in Glance
// in format "<random_8_chars>-rhcos".
// Example: fbekjlba-rhcos
func OpenStackGlanceImageName() string {
	prefix := make([]byte, 8)
	for i := range prefix {
		prefix[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(prefix) + "-rhcos"
}
