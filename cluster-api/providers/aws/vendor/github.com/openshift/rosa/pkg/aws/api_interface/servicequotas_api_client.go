package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
)

// ServiceQuotasApiClient is an interface that defines the methods that we want to use
// from the Client type in the AWS SDK ("github.com/aws/aws-sdk-go-v2/service/servicequotas")
// The aim is to only contain methods that are defined in the AWS SDK's Service Quotas
// Client.
// For the cases where logic is desired to be implemened combining Service Quotas calls
// and other logic use the pkg/aws.Client type.
// If you need to use a method provided by the AWS SDK's Service Quotas Client but it
// is not defined in this interface then it has to be added and all
// the types implementing this interface have to implement the new method.
// The reason this interface has been defined is so we can perform unit testing
// on methods that make use of the AWS Service Quotas service.
//

type ServiceQuotasApiClient interface {
	GetServiceQuota(ctx context.Context,
		params *servicequotas.GetServiceQuotaInput, optFns ...func(*servicequotas.Options),
	) (*servicequotas.GetServiceQuotaOutput, error)

	ListServiceQuotas(ctx context.Context,
		params *servicequotas.ListServiceQuotasInput, optFns ...func(*servicequotas.Options),
	) (*servicequotas.ListServiceQuotasOutput, error)
}

var _ ServiceQuotasApiClient = (*servicequotas.Client)(nil)
