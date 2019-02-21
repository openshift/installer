package rhcos

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// QEMU fetches the URL of the latest Red Hat Enterprise Linux CoreOS release.
func QEMU(ctx context.Context, channel string, build string) (string, error) {
	meta, err := fetchMetadata(ctx, channel, build)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	return fmt.Sprintf("%s/%s/%s/%s", baseURL, channel, meta.OSTreeVersion, meta.Images.QEMU.Path), nil
}
