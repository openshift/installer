package gcp

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/file/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"
)

// ServiceNameGCP is the name of the GCP Service.
type ServiceNameGCP string

const (
	// ServiceNameGCPCompute is the name used for the GCP Compute Service endpoint.
	ServiceNameGCPCompute ServiceNameGCP = "compute"

	// ServiceNameGCPContainer is the name used for the GCP Container Service endpoint.
	ServiceNameGCPContainer ServiceNameGCP = "container"

	// ServiceNameGCPCloudResource is the name used for the GCP Resource Manager Service endpoint.
	ServiceNameGCPCloudResource ServiceNameGCP = "cloudresourcemanager"

	// ServiceNameGCPDNS is the name used for the GCP DNS Service endpoint.
	ServiceNameGCPDNS ServiceNameGCP = "dns"

	// ServiceNameGCPFile is the name used for the GCP File Service endpoint.
	ServiceNameGCPFile ServiceNameGCP = "file"

	// ServiceNameGCPIAM is the name used for the GCP IAM Service endpoint.
	ServiceNameGCPIAM ServiceNameGCP = "iam"

	// ServiceNameGCPServiceUsage is the name used for the GCP Service Usage Service endpoint.
	ServiceNameGCPServiceUsage ServiceNameGCP = "serviceusage"

	// ServiceNameGCPStorage is the name used for the GCP Storage Service endpoint.
	ServiceNameGCPStorage ServiceNameGCP = "storage"
)

// CreateServiceEndpoint creates a string endpoint for a service from the endpoint name.
func CreateServiceEndpoint(endpointName string, service ServiceNameGCP) string {
	baseEndpoint := fmt.Sprintf("https://%s-%s.p.googleapis.com/", string(service), endpointName)
	switch service {
	case ServiceNameGCPCompute, ServiceNameGCPContainer, ServiceNameGCPStorage:
		baseEndpoint = fmt.Sprintf("%s%s/v1/", baseEndpoint, string(service))
	}
	return baseEndpoint
}

// CreateEndpointOption creates an Endpoint Option for a service, overriding the base/default endpoint.
func CreateEndpointOption(endpointName string, service ServiceNameGCP) option.ClientOption {
	endpoint := CreateServiceEndpoint(endpointName, service)
	return option.WithEndpoint(endpoint)
}

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
func GetComputeService(ctx context.Context, options ...option.ClientOption) (*compute.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get compute service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := compute.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	return svc, nil
}

// GetDNSService creates the dns service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetDNSService(ctx context.Context, options ...option.ClientOption) (*dns.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dns service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := dns.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create dns service: %w", err)
	}

	return svc, nil
}

// GetCloudResourceService creates the cloud resource service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetCloudResourceService(ctx context.Context, options ...option.ClientOption) (*cloudresourcemanager.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud resource service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := cloudresourcemanager.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud resource service: %w", err)
	}

	return svc, nil
}

// GetServiceUsageService creates the service usage service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetServiceUsageService(ctx context.Context, options ...option.ClientOption) (*serviceusage.APIService, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get service usage service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := serviceusage.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create service usage service: %w", err)
	}

	return svc, nil
}

// GetIAMService creates the iam service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetIAMService(ctx context.Context, options ...option.ClientOption) (*iam.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get IAM service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := iam.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create IAM service: %w", err)
	}

	return svc, nil
}

// GetStorageService creates the storage service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetStorageService(ctx context.Context, options ...option.ClientOption) (*storage.Client, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := storage.NewClient(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage service: %w", err)
	}
	return svc, nil
}

// GetFileService creates the file service. The service is created with credentials and any service
// endpoint overrides entered by the user in the installconfig.
func GetFileService(ctx context.Context, options ...option.ClientOption) (*file.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file service options: %w", err)
	}

	options = append(options, genOptions...)
	svc, err := file.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create file service: %w", err)
	}

	return svc, nil
}
