// Package libvirt collects libvirt-specific configuration.
package libvirt

import (
	"context"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/libvirt"
	libvirtdefaults "github.com/openshift/installer/pkg/types/libvirt/defaults"
	"github.com/openshift/installer/pkg/validate"
)

// Platform collects libvirt-specific configuration.
func Platform() (*libvirt.Platform, error) {
	var uri string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: libvirtdefaults.DefaultURI,
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
	}, &uri)
	if err != nil {
		return nil, err
	}

	return &libvirt.Platform{
		URI: uri,
	}, nil
}

// SetLatestImage sets the image to use to the latest image.
func SetLatestImage(p *libvirt.Platform) error {
	if p.DefaultMachinePlatform == nil {
		p.DefaultMachinePlatform = &libvirt.MachinePool{}
	}
	if p.DefaultMachinePlatform.Image != "" {
		return nil
	}
	qcowImage, err := rhcos.QEMU(context.TODO(), rhcos.DefaultChannel)
	if err != nil {
		return errors.Wrap(err, "failed to fetch QEMU image URL")
	}
	p.DefaultMachinePlatform.Image = qcowImage
	return nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}
