package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK ("github.com/aws/aws-sdk-go-v2/service/secretsmanager")
// The aim is to only contain methods that are defined in the AWS SDK's Secrets Manager
// Client.
// For the cases where logic is desired to be implemened combining Secrets Manager calls
// and other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's Secrets Manager Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS Secrets Manager service.
//

type SecretsManagerApiClient interface {
	GetSecretValue(ctx context.Context,
		params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.GetSecretValueOutput, error)

	DescribeSecret(ctx context.Context,
		params *secretsmanager.DescribeSecretInput, optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.DescribeSecretOutput, error)

	DeleteSecret(ctx context.Context,
		params *secretsmanager.DeleteSecretInput, optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.DeleteSecretOutput, error)

	CreateSecret(ctx context.Context,
		params *secretsmanager.CreateSecretInput, optFns ...func(*secretsmanager.Options),
	) (*secretsmanager.CreateSecretOutput, error)
}

// interface guard to ensure that all methods defined in the SecretsManagerApiClient
// interface are implemented by the real AWS Secrets Manager client. This interface
// guard should always compile
var _ SecretsManagerApiClient = (*secretsmanager.Client)(nil)
