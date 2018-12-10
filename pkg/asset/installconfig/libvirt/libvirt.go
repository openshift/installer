// Package libvirt collects libvirt-specific configuration.
package libvirt

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/libvirt"
)

const (
	defaultNetworkIfName  = "tt0"
	defaultNetworkIPRange = "192.168.126.0/24"
)

// Platform collects libvirt-specific configuration.
func Platform() (*libvirt.Platform, error) {
	uri, err := asset.GenerateUserProvidedAsset(
		"Libvirt Connection URI",
		&survey.Question{
			Prompt: &survey.Input{
				Message: "Libvirt Connection URI",
				Help:    "The libvirt connection URI to be used. This must be accessible from the running cluster.",
				Default: "qemu+tcp://192.168.122.1/system",
			},
			Validate: survey.ComposeValidators(survey.Required, uriValidator),
		},
		"OPENSHIFT_INSTALL_LIBVIRT_URI",
	)
	if err != nil {
		return nil, err
	}

	qcowImage, ok := os.LookupEnv("OPENSHIFT_INSTALL_LIBVIRT_IMAGE")
	if ok {
		err = validURI(qcowImage)
		if err != nil {
			return nil, errors.Wrap(err, "resolve OPENSHIFT_INSTALL_LIBVIRT_IMAGE")
		}
	} else {
		qcowImage, err = rhcos.QEMU(context.TODO(), rhcos.DefaultChannel)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch QEMU image URL")
		}
	}

   byo_raw, ok := os.LookupEnv("OPENSHIFT_INSTALL_BYO")
   byo := false
   if ok {
		   byo, err = strconv.ParseBool(byo_raw)
		   if err != nil {
			   return nil, errors.Wrap(err, "unable to parse boolean for OPENSHIFT_INSTALL_BYO")
		   }
	}

	return &libvirt.Platform{
		Network: libvirt.Network{
			IfName:  defaultNetworkIfName,
			IPRange: defaultNetworkIPRange,
		},
		DefaultMachinePlatform: &libvirt.MachinePool{
			Image: qcowImage,
		},
		URI: uri,
		BYO: byo,
	}, nil
}

// uriValidator validates if the answer provided in prompt is a valid
// url and has non-empty scheme.
func uriValidator(ans interface{}) error {
	return validURI(ans.(string))
}

// validURI validates if the URI is a valid URI with a non-empty scheme.
func validURI(uri string) error {
	parsed, err := url.Parse(uri)
	if err != nil {
		return err
	}
	if parsed.Scheme == "" {
		return fmt.Errorf("invalid URI %q (no scheme)", uri)
	}
	return nil
}
