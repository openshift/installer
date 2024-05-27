package preprovision

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/gophercloud/gophercloud/v2"
	gophercloud_openstack "github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imagedata"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// ReplaceBootstrapIgnitionInTFVars replaces the ignition file in the terraform variables.
func ReplaceBootstrapIgnitionInTFVars(ctx context.Context, tfvarsFile *asset.File, installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) error {
	if tfvarsFile == nil {
		return fmt.Errorf("missing tfvars file")
	}

	var tfvars map[string]json.RawMessage
	if err := json.Unmarshal(tfvarsFile.Data, &tfvars); err != nil {
		return fmt.Errorf("unable to decode tfvars: %w", err)
	}

	bootstrapIgnitionJSON, ok := tfvars["openstack_bootstrap_shim_ignition"]
	if !ok {
		return fmt.Errorf("bootstrap's Ignition file not found")
	}

	var bootstrapIgnition string
	if err := json.Unmarshal(bootstrapIgnitionJSON, &bootstrapIgnition); err != nil {
		return fmt.Errorf("failed to decode bootstrap's Ignition")
	}

	logrus.Debugf("Uploading Ignition to Glance")
	bootstrapShim, err := UploadIgnitionAndBuildShim(ctx, installConfig.Config.Platform.OpenStack.Cloud, clusterID.InfraID, installConfig.Config.Proxy, []byte(bootstrapIgnition))
	if err != nil {
		return fmt.Errorf("failed to build bootstrap's Ignition shim: %w", err)
	}

	bootstrapShimJSON, err := json.Marshal(bootstrapShim)
	if err != nil {
		return fmt.Errorf("failed to encode bootstrap's Ignition shim: %w", err)
	}

	logrus.Debugf("Replacing the Ignition file in the Terraform variables")
	tfvars["openstack_bootstrap_shim_ignition"] = bootstrapShimJSON

	b, err := json.MarshalIndent(tfvars, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to encode tfvars: %w", err)
	}
	tfvarsFile.Data = b
	return nil
}

// UploadIgnitionAndBuildShim uploads the bootstrap Ignition config in Glance.
func UploadIgnitionAndBuildShim(ctx context.Context, cloud string, infraID string, proxy *types.Proxy, bootstrapIgn []byte) ([]byte, error) {
	opts := openstackdefaults.DefaultClientOpts(cloud)
	conn, err := openstackdefaults.NewServiceClient(ctx, "image", opts)
	if err != nil {
		return nil, err
	}

	var userCA []byte
	{
		cloudConfig, err := clientconfig.GetCloudFromYAML(opts)
		if err != nil {
			return nil, err
		}
		// Get the ca-cert-bundle key if there is a value for cacert in clouds.yaml
		if caPath := cloudConfig.CACertFile; caPath != "" {
			userCA, err = os.ReadFile(caPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read clouds.yaml ca-cert from disk: %w", err)
			}
		}
	}

	// we need to obtain Glance public endpoint that will be used by Ignition to download bootstrap ignition files.
	// By design this should be done by using https://www.terraform.io/docs/providers/openstack/d/identity_endpoint_v3.html
	// but OpenStack default policies forbid to use this API for regular users.
	// On the other hand when a user authenticates in OpenStack (i.e. gets a token), it includes the whole service
	// catalog in the output json. So we are able to parse the data and get the endpoint from there
	// https://docs.openstack.org/api-ref/identity/v3/?expanded=token-authentication-with-scoped-authorization-detail#token-authentication-with-scoped-authorization
	// Unfortunately this feature is not currently supported by Terraform, so we had to implement it here.
	var glancePublicURL string
	{
		// Authenticate in OpenStack, get the token and extract the service catalog
		var serviceCatalog *tokens.ServiceCatalog
		{
			authResult := conn.GetAuthResult()
			auth, ok := authResult.(tokens.CreateResult)
			if !ok {
				return nil, fmt.Errorf("unable to extract service catalog")
			}

			var err error
			serviceCatalog, err = auth.ExtractServiceCatalog()
			if err != nil {
				return nil, err
			}
		}
		clientConfigCloud, err := clientconfig.GetCloudFromYAML(openstackdefaults.DefaultClientOpts(cloud))
		if err != nil {
			return nil, err
		}
		glancePublicURL, err = gophercloud_openstack.V3EndpointURL(serviceCatalog, gophercloud.EndpointOpts{
			Type:         "image",
			Availability: gophercloud.AvailabilityPublic,
			Region:       clientConfigCloud.RegionName,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot retrieve Glance URL from the service catalog: %w", err)
		}
	}

	// upload the bootstrap Ignition config in Glance and save its location
	var bootstrapConfigURL string
	{
		img, err := images.Create(ctx, conn, images.CreateOpts{
			Name:            infraID + "-ignition",
			ContainerFormat: "bare",
			DiskFormat:      "raw",
			Tags:            []string{"openshiftClusterID=" + infraID},
		}).Extract()
		if err != nil {
			return nil, fmt.Errorf("unable to create a Glance image for the bootstrap server's Ignition file: %w", err)
		}

		if res := imagedata.Upload(ctx, conn, img.ID, bytes.NewReader(bootstrapIgn)); res.Err != nil {
			return nil, fmt.Errorf("unable to upload a Glance image for the bootstrap server's Ignition file: %w", res.Err)
		}

		bootstrapConfigURL = glancePublicURL + img.File
	}

	// To allow Ignition to download its config on the bootstrap machine from a location secured by a
	// self-signed certificate, we have to provide it a valid custom ca bundle.
	// To do so we generate a small ignition config that contains just Security section with the bundle
	// and later append it to the main ignition config.
	tokenID, err := conn.GetAuthResult().ExtractTokenID()
	if err != nil {
		return nil, fmt.Errorf("unable to extract an OpenStack token: %w", err)
	}

	caRefs, err := parseCertificateBundle(userCA)
	if err != nil {
		return nil, err
	}

	var ignProxy igntypes.Proxy
	if proxy != nil {
		if proxy.HTTPProxy != "" {
			ignProxy.HTTPProxy = &proxy.HTTPProxy
		}
		if proxy.HTTPSProxy != "" {
			ignProxy.HTTPSProxy = &proxy.HTTPSProxy
		}
		if proxy.NoProxy != "" {
			noProxy := strings.Split(proxy.NoProxy, ",")
			ignProxy.NoProxy = make([]igntypes.NoProxyItem, len(noProxy))
			for i := range noProxy {
				ignProxy.NoProxy[i] = igntypes.NoProxyItem(noProxy[i])
			}
		}
	}

	data, err := ignition.Marshal(igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Timeouts: igntypes.Timeouts{
				HTTPResponseHeaders: ptr.To(120),
			},
			Security: igntypes.Security{
				TLS: igntypes.TLS{
					CertificateAuthorities: caRefs,
				},
			},
			Config: igntypes.IgnitionConfig{
				Merge: []igntypes.Resource{
					{
						Source: &bootstrapConfigURL,
						HTTPHeaders: []igntypes.HTTPHeader{
							{
								Name:  "X-Auth-Token",
								Value: &tokenID,
							},
						},
					},
				},
			},
			Proxy: ignProxy,
		},
		Storage: igntypes.Storage{
			Files: []igntypes.File{
				{
					Node: igntypes.Node{
						Path:      "/etc/hostname",
						Overwrite: ptr.To(true),
					},
					FileEmbedded1: igntypes.FileEmbedded1{
						Mode: ptr.To(420),
						Contents: igntypes.Resource{
							Source: ptr.To(dataurl.EncodeBytes([]byte(infraID + "bootstrap"))),
						},
					},
				},
				{
					Node: igntypes.Node{
						Path:      "/opt/openshift/tls/cloud-ca-cert.pem",
						Overwrite: ptr.To(true),
					},
					FileEmbedded1: igntypes.FileEmbedded1{
						Mode: ptr.To(420),
						Contents: igntypes.Resource{
							Source: ptr.To(dataurl.EncodeBytes(userCA)),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to encode the Ignition shim: %w", err)
	}

	// Check the size of the base64-rendered ignition shim isn't to big for nova
	// https://docs.openstack.org/nova/latest/user/metadata.html#user-data
	if len(base64.StdEncoding.EncodeToString(data)) > 65535 {
		return nil, fmt.Errorf("rendered bootstrap ignition shim exceeds the 64KB limit for nova user data -- try reducing the size of your CA cert bundle")
	}
	return data, nil
}

// ParseCertificateBundle loads each certificate in the bundle to the Ignition
// carrier type, ignoring any invisible character before, after and in between
// certificates.
func parseCertificateBundle(userCA []byte) ([]igntypes.Resource, error) {
	var caRefs []igntypes.Resource
	userCA = bytes.TrimSpace(userCA)
	for len(userCA) > 0 {
		var block *pem.Block
		block, userCA = pem.Decode(userCA)
		if block == nil {
			return nil, fmt.Errorf("unable to parse certificate, please check the cacert section of clouds.yaml")
		}
		caRefs = append(caRefs, igntypes.Resource{Source: ptr.To(dataurl.EncodeBytes(pem.EncodeToMemory(block)))})
		userCA = bytes.TrimSpace(userCA)
	}
	return caRefs, nil
}
