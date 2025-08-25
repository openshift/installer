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

var (
	// IMPORTANT: In v1 sdk, a constant EndpointsID is exported in each service subpackage to
	// look up custom service endpoints.
	// In v2 SDK, these constants are no longer available.
	// For backwards compatibility, we copy those constants from the SDK v1
	// and map it to ServiceID in SDK v2.
	// v1Tov2ServiceIDMap maps v1 service ID to its v2 equivalent.
	v1Tov2ServiceIDMap = map[string]string{
		"ec2": ec2.ServiceID,

		// In v1 SDK, elb and elbv2 uses the same identifier,
		// thus the same endpoint. We define a simple -v2 equivalent
		// for elbv2 to distinguish between two services.
		"elasticloadbalancing":   elb.ServiceID,
		"elasticloadbalancingv2": elbv2.ServiceID,

		"iam":           iam.ServiceID,
		"route53":       route53.ServiceID,
		"s3":            s3.ServiceID,
		"sts":           sts.ServiceID,
		"tagging":       resourcegroupstaggingapi.ServiceID,
		"servicequotas": servicequotas.ServiceID,
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

	logger logrus.FieldLogger
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
		logger:          logrus.StandardLogger(),
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
		s.logger.Debug("Custom EC2 endpoint not found, using default endpoint")
		return ec2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom EC2 endpoint found, using custom endpoint: %s", endpoint.URL)
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
		s.logger.Debug("Custom IAM endpoint not found, using default endpoint")
		return iam.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom IAM endpoint found, using custom endpoint: %s", endpoint.URL)
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
		s.logger.Debug("Custom STS endpoint not found, using default endpoint")
		return sts.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom STS endpoint found, using custom endpoint: %s", endpoint.URL)
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
		s.logger.Debug("Custom Route53 endpoint not found, using default endpoint")
		return route53.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom Route53 endpoint found, using custom endpoint: %s", endpoint.URL)
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
		s.logger.Debug("Custom Service Quotas endpoint not found, using default endpoint")
		return servicequotas.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom Service Quotas endpoint found, using custom endpoint: %s", endpoint.URL)
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
		s.logger.Debug("Custom S3 endpoint not found, using default endpoint")
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom S3 endpoint found, using custom endpoint: %s", endpoint.URL)
	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// ELBV2EndpointResolver implements EndpointResolverV2 interface for ELBV2.
type ELBV2EndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for ELBV2.
func (s *ELBV2EndpointResolver) ResolveEndpoint(ctx context.Context, params elbv2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	endpoint, ok := s.endpoints[elbv2.ServiceID]
	if !ok {
		// For backwards compatibility, we use the elb endpoint if
		// no elbv2 endpoint is defined.
		endpoint, ok = s.endpoints[elb.ServiceID]
		// If custom endpoint not found, return default endpoint for the service.
		if !ok {
			s.logger.Debug("Custom ELBV2 endpoint not found, using default endpoint")
			return elbv2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
		}

		s.logger.Debugf("Using ELB endpoint for ELBV2: %s", endpoint.URL)
		params.Endpoint = aws.String(endpoint.URL)
		params.Region = aws.String(s.endpointOptions.Region)
		return elbv2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom ELBV2 endpoint found, using custom endpoint: %s", endpoint.URL)
	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return elbv2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// TaggingEndpointResolver implements EndpointResolverV2 interface for Resource Groups Tagging API.
type TaggingEndpointResolver struct {
	*ServiceEndpointResolver
}

// ResolveEndpoint for Resource Groups Tagging API.
func (s *TaggingEndpointResolver) ResolveEndpoint(ctx context.Context, params resourcegroupstaggingapi.EndpointParameters) (smithyendpoints.Endpoint, error) {
	params.UseDualStack = aws.Bool(s.endpointOptions.UseDualStack)
	params.UseFIPS = aws.Bool(s.endpointOptions.UseFIPS)

	// If custom endpoint not found, return default endpoint for the service.
	endpoint, ok := s.endpoints[resourcegroupstaggingapi.ServiceID]
	if !ok {
		s.logger.Debug("Custom Resource Groups Tagging API endpoint not found, using default endpoint")
		return resourcegroupstaggingapi.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	s.logger.Debugf("Custom Resource Groups Tagging API endpoint found, using custom endpoint: %s", endpoint.URL)
	params.Endpoint = aws.String(endpoint.URL)
	params.Region = aws.String(s.endpointOptions.Region)
	return resourcegroupstaggingapi.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
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
