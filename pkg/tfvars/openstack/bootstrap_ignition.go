package openstack

import (
	"bytes"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	ignutil "github.com/coreos/ignition/v2/config/util"
	igntypes "github.com/coreos/ignition/v2/config/v3_1/types"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset/ignition"
)

// Starting from OpenShift 4.4 we store bootstrap Ignition configs in Glance.

// uploadBootstrapConfig uploads the bootstrap Ignition config in Glance and returns its location
func uploadBootstrapConfig(cloud string, bootstrapIgn string, clusterID string) (string, error) {
	logrus.Debugln("Creating a Glance image for your bootstrap ignition config...")
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		return "", err
	}

	imageCreateOpts := images.CreateOpts{
		Name:            fmt.Sprintf("%s-ignition", clusterID),
		ContainerFormat: "bare",
		DiskFormat:      "raw",
		Tags:            []string{fmt.Sprintf("openshiftClusterID=%s", clusterID)},
		// TODO(mfedosin): add Description when gophercloud supports it.
	}

	img, err := images.Create(conn, imageCreateOpts).Extract()
	if err != nil {
		return "", err
	}
	logrus.Debugf("Image %s was created.", img.Name)

	logrus.Debugf("Uploading bootstrap config to the image %v with ID %v", img.Name, img.ID)
	res := imagedata.Upload(conn, img.ID, strings.NewReader(bootstrapIgn))
	if res.Err != nil {
		return "", res.Err
	}
	logrus.Debugf("The config was uploaded.")

	// img.File contains location of the uploaded data
	return img.File, nil
}

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
			return nil, fmt.Errorf("unable to parse certificate, please check the cacert section of clouds.yaml")
		}

		carefs = append(carefs, igntypes.Resource{Source: ignutil.StrToPtr(dataurl.EncodeBytes(pem.EncodeToMemory(block)))})

		userCA = bytes.TrimSpace(userCA)
	}

	return carefs, nil
}

// To allow Ignition to download its config on the bootstrap machine from a location secured by a
// self-signed certificate, we have to provide it a valid custom ca bundle.
// To do so we generate a small ignition config that contains just Security section with the bundle
// and later append it to the main ignition config.
// We can't do it directly in Terraform, because Ignition provider suppors only 2.1 version, but
// Security section was added in 2.2 only.

// generateIgnitionShim is used to generate an ignition file that contains a user ca bundle
// in its Security section.
func generateIgnitionShim(userCA string, clusterID string, bootstrapConfigURL string, tokenID string) (string, error) {
	fileMode := 420

	// Hostname Config
	contents := fmt.Sprintf("%s-bootstrap", clusterID)

	hostnameConfigFile := igntypes.File{
		Node: igntypes.Node{
			Path:      "/etc/hostname",
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &fileMode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes([]byte(contents))),
			},
		},
	}

	// Openstack Ca Cert file
	openstackCAFile := igntypes.File{
		Node: igntypes.Node{
			Path:      "/opt/openshift/tls/cloud-ca-cert.pem",
			Overwrite: ignutil.BoolToPtr(true),
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &fileMode,
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(dataurl.EncodeBytes([]byte(userCA))),
			},
		},
	}

	carefs, err := parseCertificateBundle([]byte(userCA))
	if err != nil {
		return "", err
	}
	security := igntypes.Security{
		TLS: igntypes.TLS{
			CertificateAuthorities: carefs,
		},
	}

	headers := []igntypes.HTTPHeader{
		{
			Name:  "X-Auth-Token",
			Value: ignutil.StrToPtr(tokenID),
		},
	}

	ign := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version:  igntypes.MaxVersion.String(),
			Security: security,
			Config: igntypes.IgnitionConfig{
				Merge: []igntypes.Resource{
					{
						Source:      ignutil.StrToPtr(bootstrapConfigURL),
						HTTPHeaders: headers,
					},
				},
			},
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				hostnameConfigFile,
				openstackCAFile,
			},
		},
	}

	data, err := ignition.Marshal(ign)
	if err != nil {
		return "", err
	}

	// Check the size of the base64-rendered ignition shim isn't to big for nova
	// https://docs.openstack.org/nova/latest/user/metadata.html#user-data
	if len(base64.StdEncoding.EncodeToString(data)) > 65535 {
		return "", fmt.Errorf("rendered bootstrap ignition shim exceeds the 64KB limit for nova user data -- try reducing the size of your CA cert bundle")
	}

	return string(data), nil
}

// getAuthToken fetches valid OpenStack authentication token ID
func getAuthToken(cloud string) (string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return "", err
	}

	token, err := conn.GetAuthResult().ExtractTokenID()
	if err != nil {
		return "", err
	}

	return token, nil
}
