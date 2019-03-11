package rhcos

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// QEMU fetches the URL of the latest Red Hat Enterprise Linux CoreOS release.
func QEMU(ctx context.Context) (string, error) {
	meta, err := fetchLatestMetadata(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch RHCOS metadata")
	}

	return fmt.Sprintf("%s/%s/%s/%s", baseURL, getChannel(), meta.OSTreeVersion, meta.Images.QEMU.Path), nil
}
