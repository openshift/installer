package aws

import (
	"fmt"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	ec2v2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	iamv2 "github.com/aws/aws-sdk-go-v2/service/iam"
	tagsv2 "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	route53v2 "github.com/aws/aws-sdk-go-v2/service/route53"
	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	stsv2 "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/sirupsen/logrus"
)

type endpointOptionInterface interface {
	GetResolvedRegion() string
	GetDisableHTTPS() bool
	GetUseDualStackEndpoint() awsv2.DualStackEndpointState
	GetUseFIPSEndpoint() awsv2.FIPSEndpointState
}

type endpointOptions struct {
	ResolvedRegion       string
	DisableHTTPS         bool
	UseDualStackEndpoint awsv2.DualStackEndpointState
	UseFIPSEndpoint      awsv2.FIPSEndpointState
}

func (e endpointOptions) GetResolvedRegion() string {
	return e.ResolvedRegion
}
func (e endpointOptions) GetDisableHTTPS() bool { return e.DisableHTTPS }
func (e endpointOptions) GetUseDualStackEndpoint() awsv2.DualStackEndpointState {
	return e.UseDualStackEndpoint
}
func (e endpointOptions) GetUseFIPSEndpoint() awsv2.FIPSEndpointState {
	return e.UseFIPSEndpoint
}

// getDefaultServiceEndpoint will get the default service endpoint for a service and region.
func getDefaultServiceEndpoint(region, service string, opts endpointOptions) (*awsv2.Endpoint, error) { //nolint: staticcheck
	var err error
	var endpoint awsv2.Endpoint //nolint: staticcheck
	var options endpointOptionInterface = opts

	switch service {
	case "ec2":
		resolver := ec2v2.NewDefaultEndpointResolver()
		r, _ := options.(ec2v2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "elasticloadbalancing":
		resolver := elbv2.NewDefaultEndpointResolver()
		r, _ := options.(elbv2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "iam":
		resolver := iamv2.NewDefaultEndpointResolver()
		r, _ := options.(iamv2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "route53":
		resolver := route53v2.NewDefaultEndpointResolver()
		r, _ := options.(route53v2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "s3":
		resolver := s3v2.NewDefaultEndpointResolver()
		r, _ := options.(s3v2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "sts":
		resolver := stsv2.NewDefaultEndpointResolver()
		r, _ := options.(stsv2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	case "tagging":
		resolver := tagsv2.NewDefaultEndpointResolver()
		r, _ := options.(tagsv2.EndpointResolverOptions)
		endpoint, err = resolver.ResolveEndpoint(region, r)
	default:
		logrus.Debugf("Could not determine default endpoint for service: %s", service)
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to resolve aws %s service endpoint: %w", service, err)
	}

	return &endpoint, nil
}
