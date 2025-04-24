package gcp

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/file/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
	serviceusage "google.golang.org/api/serviceusage/v1beta1"

	configv1 "github.com/openshift/api/config/v1"
)

type serviceNameGCP string
type serviceVersionGCP string

const (
	computeServiceNameGCP   serviceNameGCP = "compute"
	containerServiceNameGCP serviceNameGCP = "container"
	storageServiceNameGCP   serviceNameGCP = "storage"

	serviceVersionGCPv1 serviceVersionGCP = "v1"
)

// FormatGCPEndpointInput is the structure containing input variables for formatting the GCP Service Endpoints.
type FormatGCPEndpointInput struct {
	// SkipPath should be set to true when the path should not be added to the
	// formatted endpoint. When the path is added, an example endpoint of
	// https://compute-exampleendpoint.p.googleapis.com would be formatted as
	// https://compute-exampleendpoint.p.googleapis.com/compute/v1/
	SkipPath bool
}

// FormatGCPEndpointList will format the list of GCP Service Endpoints to match the expected url
// for WithEndpoint or BasePath override endpoint options.
func FormatGCPEndpointList(endpoints []configv1.GCPServiceEndpoint, input FormatGCPEndpointInput) ([]configv1.GCPServiceEndpoint, error) {
	// The endpoints are modified to include the path
	modifiedEndpoints := []configv1.GCPServiceEndpoint{}
	for _, se := range endpoints {
		formattedURL, err := FormatGCPEndpoint(se.Name, se.URL, input)
		if err != nil {
			return nil, fmt.Errorf("failed to format GCP service endpoint %s: %w", se.Name, err)
		}

		logrus.Debugf("Formatted GCP service endpoint %s: %s", se.Name, formattedURL)
		modifiedEndpoints = append(modifiedEndpoints, configv1.GCPServiceEndpoint{Name: se.Name, URL: formattedURL})
	}
	return modifiedEndpoints, nil
}

// FormatGCPEndpoint will format the endpoint to ensure that the string is in the format that would be
// accepted by both options (WithEndpoint and BasePath override).
func FormatGCPEndpoint(service configv1.GCPServiceEndpointName, endpoint string, input FormatGCPEndpointInput) (string, error) {
	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse GCP Service Endpoint URL %s: %w", endpoint, err)
	}

	if endpointURL.Host == "" {
		logrus.Debugf("GCP endpoint did not set a host, setting host to %s", endpoint)
		endpointURL.Host = endpoint
		endpointURL.Path = ""
	}

	endpointURL.Path = strings.TrimSuffix(endpointURL.Path, "/")
	endpointURL.Scheme = "https"
	if !input.SkipPath {
		switch service {
		case configv1.GCPServiceEndpointNameCloudResource:
			endpointURL.Path = "/"
		case configv1.GCPServiceEndpointNameCompute:
			endpointURL.Path = fmt.Sprintf("/%s/%s/", computeServiceNameGCP, serviceVersionGCPv1)
		case configv1.GCPServiceEndpointNameContainer:
			endpointURL.Path = fmt.Sprintf("/%s/%s/", containerServiceNameGCP, serviceVersionGCPv1)
		case configv1.GCPServiceEndpointNameDNS:
			endpointURL.Path = "/"
		case configv1.GCPServiceEndpointNameFile:
			endpointURL.Path = "/"
		case configv1.GCPServiceEndpointNameIAM:
			endpointURL.Path = "/"
		case configv1.GCPServiceEndpointNameServiceUsage:
			// This is a beta service but the endpoint override remains the same.
			endpointURL.Path = "/"
		case configv1.GCPServiceEndpointNameStorage:
			endpointURL.Path = fmt.Sprintf("/%s/%s/", storageServiceNameGCP, serviceVersionGCPv1)
		default:
			return "", fmt.Errorf("unknown GCP Service Endpoint name: %s", service)
		}
	}
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
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP compute service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetDNSService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*dns.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dns service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameDNS {
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP dns service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetCloudResourceService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*cloudresourcemanager.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cloud resource service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameCloudResource {
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP resource manager service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetServiceUsageService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*serviceusage.APIService, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get service usage service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameServiceUsage {
			// TODO: using a beta version, this will require further testing.
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP service usage service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetIAMService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*iam.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get IAM service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameIAM {
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP IAM service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetStorageService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*storage.Client, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameStorage {
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP storage service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
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
func GetFileService(ctx context.Context, serviceEndpoints []configv1.GCPServiceEndpoint, options ...option.ClientOption) (*file.Service, error) {
	genOptions, err := getOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get file service options: %w", err)
	}

	for _, endpoint := range serviceEndpoints {
		if endpoint.Name == configv1.GCPServiceEndpointNameFile {
			formattedURL, err := FormatGCPEndpoint(endpoint.Name, endpoint.URL, FormatGCPEndpointInput{SkipPath: false})
			if err != nil {
				return nil, fmt.Errorf("failed to format GCP file service endpoint URL %s: %w", endpoint.URL, err)
			}
			genOptions = append(genOptions, option.WithEndpoint(formattedURL))
		}
	}

	options = append(options, genOptions...)
	svc, err := file.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create file service: %w", err)
	}

	return svc, nil
}
