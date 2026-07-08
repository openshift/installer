package credentialsrequest

import (
	"k8s.io/apimachinery/pkg/runtime"

	ccov1 "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	awstypes "github.com/openshift/installer/pkg/types/aws"
)

// AWSProviderSpec holds the decoded AWS-specific provider spec
// from a CredentialsRequest.
type AWSProviderSpec struct {
	StatementEntries []ccov1.StatementEntry
}

func init() {
	registerProviderSpecDecoder(awstypes.Name, decodeAWSProviderSpec)
}

func decodeAWSProviderSpec(raw *runtime.RawExtension) (interface{}, error) {
	awsSpec := &ccov1.AWSProviderSpec{}
	if err := ccov1.Codec.DecodeProviderSpec(raw, awsSpec); err != nil {
		return nil, err
	}
	return &AWSProviderSpec{
		StatementEntries: awsSpec.StatementEntries,
	}, nil
}
