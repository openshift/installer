package libvirt

import (
	"context"
	"os"

	"github.com/openshift/installer/pkg/installerassets"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/pkg/errors"
)

func getImage(ctx context.Context) ([]byte, error) {
	value := os.Getenv("OPENSHIFT_INSTALL_LIBVIRT_IMAGE")
	if value == "" {
		value, err := rhcos.QEMU(ctx, rhcos.DefaultChannel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch QEMU image URL")
		}
		return []byte(value), nil
	}

	err := validURI(value)
	if err != nil {
		return nil, errors.Wrap(err, "resolve OPENSHIFT_INSTALL_LIBVIRT_IMAGE")
	}

	return []byte(value), nil
}

func init() {
	installerassets.Defaults["libvirt/image"] = getImage
}
