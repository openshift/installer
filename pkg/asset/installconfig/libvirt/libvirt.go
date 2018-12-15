// Package libvirt collects libvirt-specific configuration.
package libvirt

import (
	"context"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/validate"
)

const (
	defaultNetworkIfName = "tt0"
)

var (
	defaultNetworkIPRange = ipnet.MustParseCIDR("192.168.126.0/24")
)

// Platform collects libvirt-specific configuration.
func Platform() (*libvirt.Platform, error) {
	var uri string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: "qemu+tcp://192.168.122.1/system",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &uri)
	if err != nil {
		return nil, err
	}

	qcowImage, err := rhcos.QEMU(context.TODO(), rhcos.DefaultChannel)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch QEMU image URL")
	}

	return &libvirt.Platform{
		Network: libvirt.Network{
			IfName:  defaultNetworkIfName,
			IPRange: *defaultNetworkIPRange,
		},
		DefaultMachinePlatform: &libvirt.MachinePool{
			Image: qcowImage,
		},
		URI: uri,
	}, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}
