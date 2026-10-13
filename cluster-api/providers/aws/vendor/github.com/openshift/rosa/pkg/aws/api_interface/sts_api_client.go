package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// StsApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK ("github.com/aws/aws-sdk-go-v2/service/sts")
// The AIM is to only contain methods that are defined in the AWS SDK's STS
// Client.
// For the cases where logic is desired to be implemened combining STS calls
// and other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's STS Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS STS service.
//

type StsApiClient interface {
	AssumeRole(ctx context.Context,
		params *sts.AssumeRoleInput, optFns ...func(*sts.Options),
	) (*sts.AssumeRoleOutput, error)

	AssumeRoleWithWebIdentity(ctx context.Context,
		params *sts.AssumeRoleWithWebIdentityInput, optFns ...func(*sts.Options),
	) (*sts.AssumeRoleWithWebIdentityOutput, error)

	GetCallerIdentity(ctx context.Context,
		params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options),
	) (*sts.GetCallerIdentityOutput, error)
}

var _ StsApiClient = (*sts.Client)(nil)
