package openstack

import (
	"encoding/json"
	"fmt"
	"strings"

	ignition "github.com/coreos/ignition/config/v2_4/types"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
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

	hostnameConfigFile := ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       "/etc/hostname",
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &fileMode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes([]byte(contents)),
			},
		},
	}

	// Openstack Ca Cert file
	openstackCAFile := ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       "/opt/openshift/tls/cloud-ca-cert.pem",
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &fileMode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes([]byte(userCA)),
			},
		},
	}

	security := ignition.Security{}
	if userCA != "" {
		security = ignition.Security{
			TLS: ignition.TLS{
				CertificateAuthorities: []ignition.CaReference{{
					Source: dataurl.EncodeBytes([]byte(userCA)),
				}},
			},
		}
	}

	headers := []ignition.HTTPHeader{
		{
			Name:  "X-Auth-Token",
			Value: tokenID,
		},
	}

	ign := ignition.Config{
		Ignition: ignition.Ignition{
			Version:  ignition.MaxVersion.String(),
			Security: security,
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{
					{
						Source:      bootstrapConfigURL,
						HTTPHeaders: headers,
					},
				},
			},
		},
		Storage: ignition.Storage{
			Files: []ignition.File{
				hostnameConfigFile,
				openstackCAFile,
			},
		},
	}

	data, err := json.Marshal(ign)
	if err != nil {
		return "", err
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
