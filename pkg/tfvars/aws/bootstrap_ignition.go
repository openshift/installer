package aws

import (
	"bytes"
	"encoding/pem"
	"fmt"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset/ignition"
)

// parseCertificateBundle loads each certificate in the bundle to the Ingition
// carrier type, ignoring any invisible character before, after and in between
// certificates.
func parseCertificateBundle(userCA []byte) ([]igntypes.Resource, error) {
	userCA = bytes.TrimSpace(userCA)

	var carefs []igntypes.Resource
	for len(userCA) > 0 {
		var block *pem.Block
		block, userCA = pem.Decode(userCA)
		if block == nil {
			return nil, fmt.Errorf("unable to parse certificate, please check the certificates")
		}

		carefs = append(carefs, igntypes.Resource{Source: ignutil.StrToPtr(dataurl.EncodeBytes(pem.EncodeToMemory(block)))})

		userCA = bytes.TrimSpace(userCA)
	}

	return carefs, nil
}

// generateIgnitionShim is used to generate an ignition file that contains a user ca bundle
// in its Security section.
func generateIgnitionShim(bootstrapConfigURL string, userCA string) (string, error) {
	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Replace: igntypes.Resource{
					Source: ignutil.StrToPtr(bootstrapConfigURL),
				},
			},
		},
	}

	carefs, err := parseCertificateBundle([]byte(userCA))
	if err != nil {
		return "", err
	}
	if len(carefs) > 0 {
		ign.Ignition.Security = igntypes.Security{
			TLS: igntypes.TLS{
				CertificateAuthorities: carefs,
			},
		}
	}

	data, err := ignition.Marshal(ign)
	if err != nil {
		return "", err
	}

	// Check the size of the raw ignition stub is less than 16KB for aws user-data
	// see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-add-user-data.html
	if len(data) > 16000 {
		return "", fmt.Errorf("rendered bootstrap ignition shim exceeds the 16KB limit for AWS user data -- try reducing the size of your CA cert bundle")
	}

	return string(data), nil
}
