package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/servicequotas"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/sirupsen/logrus"

	typesaws "github.com/openshift/installer/pkg/types/aws"
)

// Partition identifiers.
const (
	AwsPartitionID      = "aws"        // AWS Standard partition.
	AwsCnPartitionID    = "aws-cn"     // AWS China partition.
	AwsUsGovPartitionID = "aws-us-gov" // AWS GovCloud (US) partition.
	AwsIsoPartitionID   = "aws-iso"    // AWS ISO (US) partition.
	AwsIsoBPartitionID  = "aws-iso-b"  // AWS ISOB (US) partition.
)

var (
	// v1Tov2ServiceIDMap maps v1 service ID to its v2 equivalent.
	v1Tov2ServiceIDMap = map[string]string{
		"ec2":                      ec2.ServiceID,
		"elasticloadbalancing":     elb.ServiceID,
		"elasticloadbalancingv2":   elbv2.ServiceID,
		"iam":                      iam.ServiceID,
		"route53":                  route53.ServiceID,
		"s3":                       s3.ServiceID,
		"sts":                      sts.ServiceID,
		"resourcegroupstaggingapi": resourcegroupstaggingapi.ServiceID,
		"servicequotas":            servicequotas.ServiceID,
	}
)

// resolveServiceID converts a service ID in the SDK from v1 to v2.
// If the service ID is not recognized, return as-is.
func resolveServiceID(serviceID string) string {
	if v2serviceID, ok := v1Tov2ServiceIDMap[serviceID]; ok {
		return v2serviceID
	}
	return serviceID
}

// EndpointOptions defines options to configure the service endpoint resolver.
type EndpointOptions struct {
	Endpoints []typesaws.ServiceEndpoint
	Region    string

	UseDualStack bool
	UseFIPS      bool
}

// ServiceEndpointResolver implements EndpointResolverV2 interface for services.
type ServiceEndpointResolver struct {
	// endpoints is the map of provided service endpoints
	// indexed by service ID.
	endpoints       map[string]typesaws.ServiceEndpoint
	endpointOptions EndpointOptions
}

// NewServiceEndpointResolver returns a new ServiceEndpointResolver.
func NewServiceEndpointResolver(opts EndpointOptions) *ServiceEndpointResolver {
	endpointMap := make(map[string]typesaws.ServiceEndpoint, len(opts.Endpoints))
	for _, endpoint := range opts.Endpoints {
		endpointMap[resolveServiceID(endpoint.Name)] = endpoint
	}
	return &ServiceEndpointResolver{
		endpoints:       endpointMap,
		endpointOptions: opts,
	}
}

// EC2EndpointResolver implements EndpointResolverV2 interface for EC2.
type EC2EndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for EC2.
func (s *EC2EndpointResolver) ResolveEndpoint(ctx context.Context, params ec2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[ec2.ServiceID]
	if !ok {
		return ec2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return ec2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// IAMEndpointResolver implements EndpointResolverV2 interface for IAM.
type IAMEndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for IAM.
func (s *IAMEndpointResolver) ResolveEndpoint(ctx context.Context, params iam.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[iam.ServiceID]
	if !ok {
		return iam.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return iam.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// STSEndpointResolver implements EndpointResolverV2 interface for STS.
type STSEndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for STS.
func (s *STSEndpointResolver) ResolveEndpoint(ctx context.Context, params sts.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service
	endpoint, ok := s.endpoints[sts.ServiceID]
	if !ok {
		return sts.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return sts.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// Route53EndpointResolver implements EndpointResolverV2 interface for Route 53.
type Route53EndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for Route 53.
func (s *Route53EndpointResolver) ResolveEndpoint(ctx context.Context, params route53.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[route53.ServiceID]
	if !ok {
		return route53.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return route53.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// ServiceQuotasEndpointResolver implements EndpointResolverV2 interface for Service Quotas.
type ServiceQuotasEndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for Service Quotas.
func (s *ServiceQuotasEndpointResolver) ResolveEndpoint(ctx context.Context, params servicequotas.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[servicequotas.ServiceID]
	if !ok {
		return servicequotas.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return servicequotas.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// S3EndpointResolver implements EndpointResolverV2 interface for S3.
type S3EndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for S3.
func (s *S3EndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[s3.ServiceID]
	if !ok {
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// GetDefaultServiceEndpoint will get the default service endpoint for a service and region.
// Note: This uses the v1 EndpointResolver, which exposes the partition ID.
func GetDefaultServiceEndpoint(ctx context.Context, service string, opts EndpointOptions) (aws.Endpoint, error) { //nolint: staticcheck
	region := opts.Region
	useFIPs := aws.FIPSEndpointStateDisabled
	if opts.UseFIPS {
		useFIPs = aws.FIPSEndpointStateEnabled
	}
	useDualstack := aws.DualStackEndpointStateDisabled
	if opts.UseDualStack {
		useDualstack = aws.DualStackEndpointStateEnabled
	}

	var endpoint aws.Endpoint //nolint: staticcheck
	var err error
	switch service {
	case ec2.ServiceID:
		options := ec2.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = ec2.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case elb.ServiceID:
		options := elb.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = elb.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case elbv2.ServiceID:
		options := elbv2.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = elbv2.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case iam.ServiceID:
		options := iam.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = iam.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case route53.ServiceID:
		options := route53.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = route53.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case s3.ServiceID:
		options := s3.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = s3.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case sts.ServiceID:
		options := sts.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = sts.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	case resourcegroupstaggingapi.ServiceID:
		options := resourcegroupstaggingapi.EndpointResolverOptions{UseFIPSEndpoint: useFIPs, UseDualStackEndpoint: useDualstack}
		endpoint, err = resourcegroupstaggingapi.NewDefaultEndpointResolver().ResolveEndpoint(region, options)
	default:
		logrus.Debugf("Could not determine default endpoint for unknown service: %s", service)
		return endpoint, nil
	}

	if err != nil {
		return endpoint, fmt.Errorf("failed to resolve aws %s service endpoint: %w", service, err)
	}
	return endpoint, nil
}
