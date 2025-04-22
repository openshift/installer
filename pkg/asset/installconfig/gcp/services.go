package gcp

import (
	"context"
	"fmt"
	"net/url"

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

type ServiceNameGCP string
type ServiceVersionGCP string

const (
	CloudResourceManagerNameGCP ServiceNameGCP = "cloudresourcemanager"
	ComputeServiceNameGCP       ServiceNameGCP = "compute"
	ContainerServiceNameGCP     ServiceNameGCP = "container"
	DNSServiceNameGCP           ServiceNameGCP = "dns"
	FileServiceNameGCP          ServiceNameGCP = "file"
	IAMServiceNameGCP           ServiceNameGCP = "iam"
	ServiceUsageNameGCP         ServiceNameGCP = "serviceusage"
	StorageServiceNameGCP       ServiceNameGCP = "storage"

	ServiceVersionGCP1 ServiceVersionGCP = "v1"
	ServiceVersionGCP3 ServiceVersionGCP = "v3"
	// TODO: should this be betav1
	ServiceVersionGCPBeta ServiceVersionGCP = "beta"
)

// FormatGCPEndpointList will format the list of GCP Service Endpoints to match the expected url
// for WithEndpoint or BasePath override endpoint options
func FormatGCPEndpointList(endpoints []configv1.GCPServiceEndpoint) ([]configv1.GCPServiceEndpoint, error) {
	// The endpoints are modified to include the path
	modifiedEndpoints := []configv1.GCPServiceEndpoint{}
	for _, se := range endpoints {
		var err error
		var formattedURL string
		switch se.Name {
		case configv1.GCPServiceEndpointNameCloudResource:
			formattedURL, err = FormatGCPEndpoint(se.URL, CloudResourceManagerNameGCP, ServiceVersionGCP3)
		case configv1.GCPServiceEndpointNameCompute:
			formattedURL, err = FormatGCPEndpoint(se.URL, ComputeServiceNameGCP, ServiceVersionGCP1)
		case configv1.GCPServiceEndpointNameDNS:
			formattedURL, err = FormatGCPEndpoint(se.URL, DNSServiceNameGCP, ServiceVersionGCP1)
		//case configv1.GCPServiceEndpointNameFile:
		//	formattedURL, err = FormatGCPEndpoint(se.URL, FileServiceNameGCP, ServiceVersionGCP1)
		case configv1.GCPServiceEndpointNameIAM:
			formattedURL, err = FormatGCPEndpoint(se.URL, IAMServiceNameGCP, ServiceVersionGCP1)
		//case configv1.GCPServiceEndpointNameServiceUsage:
		//	formattedURL, err = FormatGCPEndpoint(se.URL, ServiceUsageNameGCP, ServiceVersionGCPBeta)
		case configv1.GCPServiceEndpointNameStorage:
			formattedURL, err = FormatGCPEndpoint(se.URL, StorageServiceNameGCP, ServiceVersionGCP1)
		}

		if err != nil || formattedURL == "" {
			return nil, fmt.Errorf("failed to format GCP Service Endpoint URL %s: %w", se.URL, err)
		}
		modifiedEndpoints = append(modifiedEndpoints, configv1.GCPServiceEndpoint{Name: se.Name, URL: formattedURL})
	}
	return modifiedEndpoints, nil
}

// FormatGCPEndpoint will format the endpoint to ensure that the string is in the format that would be
// accepted by both options (WithEndpoint and BasePath override).
func FormatGCPEndpoint(endpoint string, service ServiceNameGCP, version ServiceVersionGCP) (string, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse GCP Service Endpoint URL %s: %w", endpoint, err)
	}

	endpointURL.Scheme = "https"
	endpointURL.Path = fmt.Sprintf("/%s/%s/", service, version)
	return endpointURL.String(), nil
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
