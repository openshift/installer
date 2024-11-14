package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

// OrganizationsApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK ("github.com/aws/aws-sdk-go-v2/service/organizations")
// The AIM is to only contain methods that are defined in the AWS SDK's Organizations
// Client.
// For the cases where logic is desired to be implemented combining Organizations calls
// and other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's Organizations Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS Organizations service.
//

type OrganizationsApiClient interface {
	CloseAccount(ctx context.Context,
		params *organizations.CloseAccountInput, optFns ...func(*organizations.Options),
	) (*organizations.CloseAccountOutput, error)

	CreateOrganization(ctx context.Context,
		params *organizations.CreateOrganizationInput, optFns ...func(*organizations.Options),
	) (*organizations.CreateOrganizationOutput, error)

	CreatePolicy(ctx context.Context,
		params *organizations.CreatePolicyInput, optFns ...func(*organizations.Options),
	) (*organizations.CreatePolicyOutput, error)

	DeletePolicy(ctx context.Context,
		params *organizations.DeletePolicyInput, optFns ...func(*organizations.Options),
	) (*organizations.DeletePolicyOutput, error)

	DeleteResourcePolicy(ctx context.Context,
		params *organizations.DeleteResourcePolicyInput, optFns ...func(*organizations.Options),
	) (*organizations.DeleteResourcePolicyOutput, error)

	ListPolicies(ctx context.Context,
		params *organizations.ListPoliciesInput, optFns ...func(*organizations.Options),
	) (*organizations.ListPoliciesOutput, error)

	ListTagsForResource(ctx context.Context,
		params *organizations.ListTagsForResourceInput, optFns ...func(*organizations.Options),
	) (*organizations.ListTagsForResourceOutput, error)

	PutResourcePolicy(ctx context.Context,
		params *organizations.PutResourcePolicyInput, optFns ...func(*organizations.Options),
	) (*organizations.PutResourcePolicyOutput, error)

	TagResource(ctx context.Context,
		params *organizations.TagResourceInput, optFns ...func(*organizations.Options),
	) (*organizations.TagResourceOutput, error)

	UntagResource(ctx context.Context,
		params *organizations.UntagResourceInput, optFns ...func(*organizations.Options),
	) (*organizations.UntagResourceOutput, error)
}

// interface guard to ensure that all methods defined in the OrganizationsApiClient
// interface are implemented by the real AWS Organizations client. This interface
// guard should always compile
var _ OrganizationsApiClient = (*organizations.Client)(nil)
