package openstack

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/kubernetes/apimachinery/pkg/util/rand"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
)

// createBootstrapSwiftObject creates a container and object in swift with the bootstrap ignition config.
func createBootstrapSwiftObject(cloud string, bootstrapIgn string, clusterID string) (string, error) {
	logrus.Debugln("Creating a Swift container for your bootstrap ignition...")
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("object-store", &opts)
	if err != nil {
		return "", err
	}

	containerCreateOpts := containers.CreateOpts{
		ContainerRead: ".r:*",

		// "kubernetes.io/cluster/${var.cluster_id}" = "owned"
		Metadata: map[string]string{
			"Name":               fmt.Sprintf("%s-ignition", clusterID),
			"openshiftClusterID": clusterID,
		},
	}

	_, err = containers.Create(conn, clusterID, containerCreateOpts).Extract()
	if err != nil {
		return "", err
	}
	logrus.Debugf("Container %s was created.", clusterID)

	logrus.Debugf("Creating a Swift object in container %s containing your bootstrap ignition...", clusterID)
	objectCreateOpts := objects.CreateOpts{
		ContentType: "text/plain",
		Content:     strings.NewReader(bootstrapIgn),
		DeleteAfter: 3600,
	}

	objID := rand.String(16)

	_, err = objects.Create(conn, clusterID, objID, objectCreateOpts).Extract()
	if err != nil {
		return "", err
	}
	logrus.Debugf("The object was created.")

	return objID, nil
}

// To allow Ignition to download its config on the bootstrap machine from a location secured by a
// self-signed certificate, we have to provide it a valid custom ca bundle.
// To do so we generate a small ignition config that contains just Security section with the bundle
// and later append it to the main ignition config.
// We can't do it directly in Terraform, because Ignition provider suppors only 2.1 version, but
// Security section was added in 2.2 only.

// generateIgnitionShim is used to generate an ignition file that contains a user ca bundle
// in its Security section.
func generateIgnitionShim(userCA string, clusterID string, swiftObject string) (string, error) {
	fileMode := 420

	// DHCP Config
	contents := `[main]
dhcp=dhclient`

	dhcpConfigFile := ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       "/etc/NetworkManager/conf.d/dhcp-client.conf",
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &fileMode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes([]byte(contents)),
			},
		},
	}

	// DNS Config
	contents = `send dhcp-client-identifier = hardware;
prepend domain-name-servers 127.0.0.1;`

	dnsConfigFile := ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       "/etc/dhcp/dhclient.conf",
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &fileMode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes([]byte(contents)),
			},
		},
	}

	// Hostname Config
	contents = fmt.Sprintf("%s-bootstrap", clusterID)

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
		carefs := []ignition.CaReference{}
		rest := []byte(userCA)

		for {
			var block *pem.Block
			block, rest = pem.Decode(rest)
			if block == nil {
				return "", fmt.Errorf("unable to parse certificate, please check the cacert section of clouds.yaml")
			}

			carefs = append(carefs, ignition.CaReference{Source: dataurl.EncodeBytes(pem.EncodeToMemory(block))})

			if len(rest) == 0 {
				break
			}
		}

		security = ignition.Security{
			TLS: ignition.TLS{
				CertificateAuthorities: carefs,
			},
		}
	}

	ign := ignition.Config{
		Ignition: ignition.Ignition{
			Version:  ignition.MaxVersion.String(),
			Security: security,
			Config: ignition.IgnitionConfig{
				Append: []ignition.ConfigReference{
					{
						Source: swiftObject,
					},
				},
			},
		},
		Storage: ignition.Storage{
			Files: []ignition.File{
				dhcpConfigFile,
				dnsConfigFile,
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
