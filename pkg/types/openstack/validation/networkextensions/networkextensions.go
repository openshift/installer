package networkextensions

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
)

const (
	// ErrMissingStandardAttrTag is returned when standard-attr-tag is not found in the cloud
	ErrMissingStandardAttrTag Error = "openstack platform does not have the required standard-attr-tag network extension"

	// ErrInvalidStandardAttrTag is returned when standard-attr-tag is too old and missing the required tag and tag-ext extensions
	ErrInvalidStandardAttrTag Error = "openstack platform's neutron extension standard-attr-tag is too old and missing the required tag and tag-ext extensions"
)

// Get returns a slice of the extensions available in the Neutron instance
// targeted by networkClient, or a non-nil error.
func Get(ctx context.Context, networkClient *gophercloud.ServiceClient) ([]extensions.Extension, error) {
	allPages, err := extensions.List(networkClient).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list network extensions: %w", err)
	}
	allExtensions, err := extensions.ExtractExtensions(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response with network extensions list: %w", err)
	}

	return allExtensions, nil
}

// Validate returns a non-nil error if the available extensions do not match
// the Installer requirements.
//
// Network resource tagging is a bit complex with OpenStack:
//   - up until OpenStack Ocata / OSP 11, there were two extensions for
//     tagging neutron resources, `tag` for network only and `tag-ext`
//     for subnet, subnetpool, port, and router.
//   - OpenStack Pike / OSP 12 introduced the `standard-attr-tag`
//     extension for trunk, policy, security_group, and floatingip.
//   - with OpenStack Rocky / OSP 14, everything went under the
//     `standard-attr-tag` extension.
//
// We need to check that:
//  1. `standard-attr-tag` extension is enabled
//  2. `standard-attr-tag` covers all the necessary resources (from the
//     extension description) or that the `tag` and `tag-ext`
//     extensions are enabled as well
func Validate(availableExtensions []extensions.Extension) error {
	var (
		standardAttrTagEnabled      = false
		standardAttrTagAllResources = false
		tagEnabled                  = false
		tagExtEnabled               = false
	)

	for _, extension := range availableExtensions {
		if extension.Alias == "standard-attr-tag" {
			standardAttrTagEnabled = true
			// this should be good enough to ensure OpenStack version is >= Rocky
			if strings.Contains(extension.Name, "subnet") {
				standardAttrTagAllResources = true
			}
		}
		if extension.Alias == "tag" {
			tagEnabled = true
		}
		if extension.Alias == "tag-ext" {
			tagExtEnabled = true
		}
	}

	if !standardAttrTagEnabled {
		return ErrMissingStandardAttrTag
	}

	if !(standardAttrTagAllResources || (tagEnabled && tagExtEnabled)) {
		return ErrInvalidStandardAttrTag
	}

	return nil
}
