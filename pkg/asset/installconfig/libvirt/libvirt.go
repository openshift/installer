// Package libvirt collects libvirt-specific configuration.
package libvirt

import (
	"os/exec"

	survey "gopkg.in/AlecAivazis/survey.v1"

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

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validate.URI(ans.(string))
}

// AddWildcardDNS configures the host's dnsmasq, so that all apps can
// found by Ingress Operator
func AddWildcardDNS(domain string) error {
	cmd := exec.Command("hack/configure-dnsmasq.sh", domain)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
