package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
)

var (
	regionUSSouth = "us-south"

	cisURL                = "test-cis-url.com"
	cosURL                = "test-cos-url.com"
	iamURL                = "test-iam-url.com"
	globalTaggingURL      = "test-gt-url.com"
	keyProtectURL         = "test-kp-url.com"
	resourceControllerURL = "test-rc-url.com"
	resourceManagerURL    = "test-rm-url.com"
	vpcURL                = "test-vpc-url.com"
)

func TestGetRegionAndEndpointsFlag(t *testing.T) {
	testCases := []struct {
		name           string
		region         string
		endpoints      []configv1.IBMCloudServiceEndpoint
		expectedResult string
	}{
		{
			name:           "no service endpionts",
			region:         regionUSSouth,
			endpoints:      []configv1.IBMCloudServiceEndpoint{},
			expectedResult: "",
		},
		{
			name:   "ignore iam service endpoint",
			region: regionUSSouth,
			endpoints: []configv1.IBMCloudServiceEndpoint{
				{
					Name: configv1.IBMCloudServiceIAM,
					URL:  iamURL,
				},
			},
			expectedResult: "",
		},
		{
			name:   "single service endpoint flag",
			region: regionUSSouth,
			endpoints: []configv1.IBMCloudServiceEndpoint{
				{
					Name: configv1.IBMCloudServiceCOS,
					URL:  cosURL,
				},
			},
			expectedResult: fmt.Sprintf("%s:%s=%s", regionUSSouth, "cos", cosURL),
		},
		{
			name:   "all supported capi service endpoints flag",
			region: regionUSSouth,
			endpoints: []configv1.IBMCloudServiceEndpoint{
				{
					Name: configv1.IBMCloudServiceCIS,
					URL:  cisURL,
				},
				{
					Name: configv1.IBMCloudServiceCOS,
					URL:  cosURL,
				},
				{
					Name: configv1.IBMCloudServiceIAM,
					URL:  iamURL,
				},
				{
					Name: configv1.IBMCloudServiceGlobalTagging,
					URL:  globalTaggingURL,
				},
				{
					Name: configv1.IBMCloudServiceKeyProtect,
					URL:  keyProtectURL,
				},
				{
					Name: configv1.IBMCloudServiceResourceController,
					URL:  resourceControllerURL,
				},
				{
					Name: configv1.IBMCloudServiceResourceManager,
					URL:  resourceManagerURL,
				},
				{
					Name: configv1.IBMCloudServiceVPC,
					URL:  vpcURL,
				},
			},
			expectedResult: fmt.Sprintf("%s:%s=%s,%s=%s,%s=%s,%s=%s,%s=%s", regionUSSouth, "cos", cosURL, "globaltagging", globalTaggingURL, "rc", resourceControllerURL, "rm", resourceManagerURL, "vpc", vpcURL),
		},
	}

	for _, tCase := range testCases {
		m := Metadata{
			Region:           tCase.region,
			ServiceEndpoints: tCase.endpoints,
		}
		assert.Equal(t, tCase.expectedResult, m.GetRegionAndEndpointsFlag())
	}
}
