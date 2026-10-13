/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package endpoints

import (
	"context"
	"errors"
	"net/url"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	rgapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyendpoints "github.com/aws/smithy-go/endpoints"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
)

var (
	errServiceEndpointFormat             = errors.New("must be formatted as ${ServiceID}=${URL}")
	errServiceEndpointSigningRegion      = errors.New("must be formatted as ${SigningRegion}:${ServiceID1}=${URL1},${ServiceID2}=${URL2...}")
	errServiceEndpointURL                = errors.New("must use a valid URL as a service-endpoint")
	errServiceEndpointServiceID          = errors.New("must use a valid serviceID from the AWS GO SDK")
	errServiceEndpointDuplicateServiceID = errors.New("same serviceID defined twice for signing region")
	serviceEndpointsMap                  = map[string]serviceEndpoint{}
	compatServiceIDMap                   = map[string]string{
		"s3":                   s3.ServiceID,
		"elasticloadbalancing": elb.ServiceID,
		"ec2":                  ec2.ServiceID,
		"tagging":              rgapi.ServiceID,
		"sqs":                  sqs.ServiceID,
		"events":               eventbridge.ServiceID,
		"eks":                  eks.ServiceID,
		"ssm":                  ssm.ServiceID,
		"sts":                  sts.ServiceID,
		"secretsmanager":       secretsmanager.ServiceID,
	}
)

// serviceEndpoint contains AWS Service resolution information for SDK V2.
type serviceEndpoint struct {
	ServiceID     string
	URL           string
	SigningRegion string
}

// ParseFlag parses the command line flag of service endponts in the format ${SigningRegion1}:${ServiceID1}=${URL1},${ServiceID2}=${URL2}...;${SigningRegion2}...
// returning a set of ServiceEndpoints.
func ParseFlag(serviceEndpoints string) error {
	if serviceEndpoints == "" {
		return nil
	}

	// There is no Enum for serviceID in V2, so we will directly use the provided endpoint
	// If the custom endpoint has any issue, EndpointResolverV2 falls back to default endpoint of specific service.
	signingRegionConfigs := strings.Split(serviceEndpoints, ";")
	for _, regionConfig := range signingRegionConfigs {
		components := strings.SplitN(regionConfig, ":", 2)
		if len(components) != 2 {
			return errServiceEndpointSigningRegion
		}
		signingRegion := components[0]
		servicePairs := strings.Split(components[1], ",")
		seenServices := []string{}
		for _, servicePair := range servicePairs {
			kv := strings.Split(servicePair, "=")
			if len(kv) != 2 {
				return errServiceEndpointFormat
			}
			serviceID := kv[0]
			if serviceID == "" {
				return errServiceEndpointServiceID
			}
			if slices.Contains(seenServices, serviceID) {
				return errServiceEndpointDuplicateServiceID
			}
			seenServices = append(seenServices, serviceID)

			// In v1 sdk, a constant EndpointsID is exported in each service to look up the custom service endpoint.
			// For example: https://github.com/aws/aws-sdk-go/blob/070853e88d22854d2355c2543d0958a5f76ad407/service/resourcegroupstaggingapi/service.go#L33-L34
			// In v2 SDK, these constants are no longer available.
			// For backwards compatibility, we copy those constants from the SDK v1 and map it to ServiceID in SDK v2.
			if v2serviceID, ok := compatServiceIDMap[serviceID]; ok {
				serviceID = v2serviceID
			}

			URL, err := url.ParseRequestURI(kv[1])
			if err != nil {
				return errServiceEndpointURL
			}
			endpoint := serviceEndpoint{
				ServiceID:     serviceID,
				URL:           URL.String(),
				SigningRegion: signingRegion,
			}
			serviceEndpointsMap[serviceID] = endpoint
		}

		// In v1 SDK, elb and elbv2 uses the same identifier, thus the same endpoint.
		// elbv2: https://github.com/aws/aws-sdk-go/blob/070853e88d22854d2355c2543d0958a5f76ad407/service/elbv2/service.go#L32-L33
		// elb: https://github.com/aws/aws-sdk-go/blob/070853e88d22854d2355c2543d0958a5f76ad407/service/elb/service.go#L32-L33
		// For backwards compatibility, if elbv2 endpoint is undefined, the elbv2 endpoint resolver should fall back to elb endpoint if any.
		if _, ok := serviceEndpointsMap[elbv2.ServiceID]; !ok {
			if elbEp, ok := serviceEndpointsMap[elb.ServiceID]; ok {
				serviceEndpointsMap[elbv2.ServiceID] = serviceEndpoint{
					ServiceID:     elbv2.ServiceID,
					URL:           elbEp.URL,
					SigningRegion: elbEp.SigningRegion,
				}
			}
		}
	}
	return nil
}

// GetPartitionFromRegion returns the cluster partition.
func GetPartitionFromRegion(region string) string {
	if partition := GetPartition(region); partition != nil {
		return partition.Name
	}

	return defaultPartition
}

// Custom EndpointResolverV2 ResolveEndpoint handlers.

// MultiServiceEndpointResolver implements EndpointResolverV2 interface for services.
type MultiServiceEndpointResolver struct {
	endpoints map[string]serviceEndpoint
}

// NewMultiServiceEndpointResolver returns new MultiServiceEndpointResolver.
func NewMultiServiceEndpointResolver() *MultiServiceEndpointResolver {
	return &MultiServiceEndpointResolver{
		endpoints: serviceEndpointsMap,
	}
}

// S3EndpointResolver implements EndpointResolverV2 interface for S3.
type S3EndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for S3.
func (s *S3EndpointResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[s3.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// ELBEndpointResolver implements EndpointResolverV2 interface for ELB.
type ELBEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for ELB.
func (s *ELBEndpointResolver) ResolveEndpoint(ctx context.Context, params elb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[elb.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return elb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return elb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// ELBV2EndpointResolver implements EndpointResolverV2 interface for ELBV2.
type ELBV2EndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for ELBV2.
func (s *ELBV2EndpointResolver) ResolveEndpoint(ctx context.Context, params elbv2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[elbv2.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return elbv2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return elbv2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// EC2EndpointResolver implements EndpointResolverV2 interface for EC2.
type EC2EndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for ELBV2.
func (s *EC2EndpointResolver) ResolveEndpoint(ctx context.Context, params ec2.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[ec2.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return ec2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return ec2.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// RGAPIEndpointResolver implements EndpointResolverV2 interface for RGAPI.
type RGAPIEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for RGAPI.
func (s *RGAPIEndpointResolver) ResolveEndpoint(ctx context.Context, params rgapi.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[rgapi.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return rgapi.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return rgapi.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// SQSEndpointResolver implements EndpointResolverV2 interface for SQS.
type SQSEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for SQS.
func (s *SQSEndpointResolver) ResolveEndpoint(ctx context.Context, params sqs.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[sqs.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return sqs.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return sqs.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// EventBridgeEndpointResolver implements EndpointResolverV2 interface for EventBridge.
type EventBridgeEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for EventBridge.
func (s *EventBridgeEndpointResolver) ResolveEndpoint(ctx context.Context, params eventbridge.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[eventbridge.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return eventbridge.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return eventbridge.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// EKSEndpointResolver implements EndpointResolverV2 interface for EKS.
type EKSEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for EKS.
func (s *EKSEndpointResolver) ResolveEndpoint(ctx context.Context, params eks.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[eks.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return eks.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return eks.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// SSMEndpointResolver implements EndpointResolverV2 interface for SSM.
type SSMEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for SSM.
func (s *SSMEndpointResolver) ResolveEndpoint(ctx context.Context, params ssm.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[ssm.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return ssm.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return ssm.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// STSEndpointResolver implements EndpointResolverV2 interface for STS.
type STSEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for STS.
func (s *STSEndpointResolver) ResolveEndpoint(ctx context.Context, params sts.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[sts.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return sts.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return sts.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

// SecretsManagerEndpointResolver implements EndpointResolverV2 interface for Secrets Manager.
type SecretsManagerEndpointResolver struct {
	*MultiServiceEndpointResolver
}

// ResolveEndpoint for Secrets Manager.
func (s *SecretsManagerEndpointResolver) ResolveEndpoint(ctx context.Context, params secretsmanager.EndpointParameters) (smithyendpoints.Endpoint, error) {
	// If custom endpoint not found, return default endpoint for the service
	log := logger.FromContext(ctx)
	endpoint, ok := s.endpoints[secretsmanager.ServiceID]

	if !ok {
		log.Debug("Custom endpoint not found, using default endpoint")
		return secretsmanager.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
	}

	log.Debug("Custom endpoint found, using custom endpoint", "endpoint", endpoint.URL)
	params.Endpoint = &endpoint.URL
	params.Region = &endpoint.SigningRegion
	return secretsmanager.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}
