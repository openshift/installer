package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// CloudFormationApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK (github.com/aws/aws-sdk-go-v2/service/cloudformation)
// The AIM is to only contain methods that are defined in the AWS SDK's EC2
// Client.
// For the cases where logic is desired to be implemented combining cloudformation calls and
// other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's cloudformation Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS cloudformation service.
//

type CloudFormationApiClient interface {
	CreateChangeSet(ctx context.Context,
		params *cloudformation.CreateChangeSetInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.CreateChangeSetOutput, error)

	CreateStack(ctx context.Context,
		params *cloudformation.CreateStackInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.CreateStackOutput, error)

	DeleteStack(ctx context.Context,
		params *cloudformation.DeleteStackInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DeleteStackOutput, error)

	DescribeChangeSet(ctx context.Context,
		params *cloudformation.DescribeChangeSetInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DescribeChangeSetOutput, error)

	DescribeStackInstance(ctx context.Context,
		params *cloudformation.DescribeStackInstanceInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DescribeStackInstanceOutput, error)

	DescribeStackResources(ctx context.Context,
		params *cloudformation.DescribeStackResourcesInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DescribeStackResourcesOutput, error)

	DescribeStackSetOperation(ctx context.Context,
		params *cloudformation.DescribeStackSetOperationInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DescribeStackSetOperationOutput, error)

	DescribeStacks(ctx context.Context,
		params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.DescribeStacksOutput, error)

	ExecuteChangeSet(ctx context.Context,
		params *cloudformation.ExecuteChangeSetInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ExecuteChangeSetOutput, error)

	GetTemplateSummary(ctx context.Context,
		params *cloudformation.GetTemplateSummaryInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.GetTemplateSummaryOutput, error)

	ListStackInstanceResourceDrifts(ctx context.Context,
		params *cloudformation.ListStackInstanceResourceDriftsInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ListStackInstanceResourceDriftsOutput, error)

	ListStackInstances(ctx context.Context,
		params *cloudformation.ListStackInstancesInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ListStackInstancesOutput, error)

	ListStackSetOperationResults(ctx context.Context,
		params *cloudformation.ListStackSetOperationResultsInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ListStackSetOperationResultsOutput, error)

	ListStackSetOperations(ctx context.Context,
		params *cloudformation.ListStackSetOperationsInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ListStackSetOperationsOutput, error)

	ListStacks(ctx context.Context,
		params *cloudformation.ListStacksInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.ListStacksOutput, error)

	RollbackStack(ctx context.Context,
		params *cloudformation.RollbackStackInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.RollbackStackOutput, error)

	UpdateStack(ctx context.Context,
		params *cloudformation.UpdateStackInput, optFns ...func(*cloudformation.Options),
	) (*cloudformation.UpdateStackOutput, error)
}

// interface guard to ensure that all methods defined in the cloudformationApiClient
// interface are implemented by the real AWS cloudformation client. This interface
// guard should always compile
var _ CloudFormationApiClient = (*cloudformation.Client)(nil)
