package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/file/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"
	"google.golang.org/api/storage/v1"

	configv1 "github.com/openshift/api/config/v1"
)

// getOptions creates the options for use during service creation.
func getOptions(ctx context.Context) ([]option.ClientOption, error) {
	ssn, err := GetSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	options := []option.ClientOption{
		option.WithCredentials(ssn.Credentials),
	}
	return options, nil
}

// GetComputeService creates the compute service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetComputeService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*compute.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get compute service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameCompute {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := compute.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameCompute {
	//		svc.BasePath = endpoint.URL
	//	}
	//}

	return svc, nil
}

// GetDNSService creates the dns service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetDNSService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*dns.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dns service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameDNS {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := dns.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameDNS {
	//		svc.BasePath = endpoint.URL
	//	}
	//}
	return svc, nil
}

// GetCloudResourceService creates the cloud resource service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetCloudResourceService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*cloudresourcemanager.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud resource service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameCloudResource {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := cloudresourcemanager.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud resource service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameCloudResource {
	//		svc.BasePath = endpoint.URL
	//	}
	//}
	return svc, nil
}

// GetServiceUsageService creates the service usage service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetServiceUsageService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*serviceusage.APIService, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get service usage service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameServiceUsage {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := serviceusage.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create service usage service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameServiceUsage {
	//		svc.BasePath = endpoint.URL
	//	}
	//}
	return svc, nil
}

// GetIAMService creates the iam service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetIAMService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*iam.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get IAM service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameIAM {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := iam.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameIAM {
	//		svc.BasePath = endpoint.URL
	//	}
	//}
	return svc, nil
}

// GetStorageService creates the storage service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetStorageService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*storage.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameStorage {
			genOptions = append(genOptions, option.WithEndpoint(endpoint.URL))
		}
	}

	options = append(options, genOptions...)
	svc, err := storage.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage service: %w", err)
	}

	//for _, endpoint := range serviceEndpoints {
	//	if endpoint.Name == configv1.GCPServiceEndpointNameStorage {
	//		svc.BasePath = endpoint.URL
	//	}
	//}
	return svc, nil
}

// GetFileService creates the file service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetFileService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*file.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := file.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create file service: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameFile {
			svc.BasePath = endpoint.URL
		}
	}
	return svc, nil
}
