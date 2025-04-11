package ibmcloud

import (
	"fmt"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset/ignition"
)

// GenerateIgnitionShimWithCredentials creates an Ignition Config shim, directing additional configuration request to the provided URL, typically a COS object. A provided IAM access token is embedded within the Ignition Config as an HTTP header.
func GenerateIgnitionShimWithCredentials(url string, iamToken string) ([]byte, error) {
	config := &igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
			Config: igntypes.IgnitionConfig{
				Replace: igntypes.Resource{
					Source: ptr.To(url),
					// NOTE(cjschaef): Replace authorization with Service ID credentials.
					HTTPHeaders: igntypes.HTTPHeaders{
						{
							Name:  "Authorization",
							Value: ptr.To(fmt.Sprintf("Bearer %s", iamToken)),
						},
					},
				},
			},
		},
	}

	return ignition.Marshal(config)
}

// GetIgnitionBucketName returns the name for the COS Bucket designed to hold temporary bootstrap Ignition data.
func GetIgnitionBucketName(infraID string) string {
	return fmt.Sprintf("%s-bootstrap-ignition", infraID)
}

// GetIgnitionFileName returns the name of the file in COS which holds the bootstrap Ignition config.
func GetIgnitionFileName() string {
	return "bootstrap.ign"
}
