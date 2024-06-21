package networkextensions_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"

	"github.com/openshift/installer/pkg/types/openstack/validation/networkextensions"
)

func TestValidate(t *testing.T) {
	type checkFunc func(e error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	ext := func(alias, name string) (ext extensions.Extension) {
		return extensions.Extension{Alias: alias, Name: name}
	}

	noError := func(have error) error {
		if have != nil {
			return fmt.Errorf("expected nil error, got: %q", have)
		}
		return nil
	}

	networkextensionsError := func(have error) error {
		if have == nil {
			return fmt.Errorf("expected networkextensions.Error error, got nil")
		}
		var want networkextensions.Error
		if !errors.As(have, &want) {
			return fmt.Errorf("expected error to be of networkextensions.Error type")
		}
		return nil
	}

	for _, tc := range [...]struct {
		name       string
		extensions []extensions.Extension
		checks     []checkFunc
	}{
		{
			"compliant OSP 14+",
			[]extensions.Extension{
				ext("standard-attr-tag", "Tag support for resources with standard attribute: port, subnet, subnetpool, network, router, floatingip, policy, security_group, trunk, network_segment_range"),
			},
			check(noError),
		},
		{
			"compliant OSP 12",
			[]extensions.Extension{
				ext("standard-attr-tag", "Tag support for resources with standard attribute: trunk, policy, security_group, floatingip"),
				ext("tag", "Tag support"),
				ext("tag-ext", "Tag support for resources: subnet, subnetpool, port, router"),
			},
			check(noError),
		},
		{
			"OSP 11",
			[]extensions.Extension{
				ext("tag", "Tag support"),
				ext("tag-ext", "Tag support for resources: subnet, subnetpool, port, router"),
			},
			check(networkextensionsError),
		},
		{
			"no tag extension",
			[]extensions.Extension{},
			check(networkextensionsError),
		},
		{
			"OSP 12, no standard-attr-tag",
			[]extensions.Extension{
				ext("tag", "Tag support"),
				ext("tag-ext", "Tag support for resources: subnet, subnetpool, port, router"),
			},
			check(networkextensionsError),
		},
		{
			"OSP 12, no tag",
			[]extensions.Extension{
				ext("standard-attr-tag", "Tag support for resources with standard attribute: trunk, policy, security_group, floatingip"),
				ext("tag-ext", "Tag support for resources: subnet, subnetpool, port, router"),
			},
			check(networkextensionsError),
		},
		{
			"OSP 12, no tag-ext",
			[]extensions.Extension{
				ext("standard-attr-tag", "Tag support for resources with standard attribute: trunk, policy, security_group, floatingip"),
				ext("tag", "Tag support"),
			},
			check(networkextensionsError),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := networkextensions.Validate(tc.extensions)
			for _, check := range tc.checks {
				if e := check(err); e != nil {
					t.Error(e)
				}
			}
		})
	}
}
