/*
Copyright 2026 The Kubernetes Authors.

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

package powervs

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-cos-sdk-go/aws"
	cosSession "github.com/IBM/ibm-cos-sdk-go/aws/session"
	tgapiv1 "github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	resourcemanagerv2 "github.com/IBM/platform-services-go-sdk/resourcemanagerv2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/endpoints"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/authenticator"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/cos"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/powervs"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcecontroller"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/resourcemanager"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/transitgateway"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/services/vpc"
)

// ClientOptions contains generic configurations required to build IBM Cloud clients.
type ClientOptions struct {
	Authenticator   core.Authenticator
	Zone            string
	WorkspaceID     string
	VPCRegion       string
	ServiceEndpoint []endpoints.ServiceEndpoint
	Debug           bool
}

// COSClientOptions carries the parameters needed to construct a COS client.
// These are kept separate from ClientOptions because they depend on runtime
// Status fields (COS instance ID, bucket region) that are not known at scope
// construction time for non-machine scopes.
type COSClientOptions struct {
	InstanceID      string
	BucketRegion    string
	ServiceEndpoint []endpoints.ServiceEndpoint
}

// HMACCOSClientOptions carries the parameters needed to construct an HMAC-credential
// COS client suitable for generating SigV4 pre-signed URLs.
// The HMAC access_key_id and secret_access_key are read from the Kubernetes Secret
// whose name is stored in IBMPowerVSCluster.Status.COSInstance.HMACSecretName.
type HMACCOSClientOptions struct {
	// Client is the Kubernetes client used to read the HMAC Secret.
	Client              client.Client
	HMACSecretName      string
	HMACSecretNamespace string
	BucketRegion        string
	ServiceEndpoint     []endpoints.ServiceEndpoint
}

// ClientBuilder defines the contract for constructing IBM Cloud service clients.
// This interface enables clean dependency injection and robust mocking for tests.
type ClientBuilder interface {
	GetAuthenticator(ctx context.Context) (core.Authenticator, error)
	GetPowerVSClient(ctx context.Context, options ClientOptions) (powervs.PowerVS, error)
	GetVPCClient(ctx context.Context, options ClientOptions) (vpc.Vpc, error)
	GetTransitGatewayClient(ctx context.Context, options ClientOptions) (transitgateway.TransitGateway, error)
	GetResourceControllerClient(ctx context.Context, options ClientOptions) (resourcecontroller.ResourceController, error)
	GetResourceManagerClient(ctx context.Context, options ClientOptions) (resourcemanager.ResourceManager, error)
	GetCOSClient(ctx context.Context, options COSClientOptions) (cos.Cos, error)
	GetHMACCOSClient(ctx context.Context, options HMACCOSClientOptions) (cos.Cos, error)
}

// ProdClientBuilder is the production implementation of the ClientBuilder interface.
type ProdClientBuilder struct{}

// GetAuthenticator returns an IBM Cloud authenticator using default env/file credentials.
func (b ProdClientBuilder) GetAuthenticator(_ context.Context) (core.Authenticator, error) {
	return authenticator.GetAuthenticator()
}

// GetPowerVSClient constructs a production PowerVS client for the given options.
func (b ProdClientBuilder) GetPowerVSClient(ctx context.Context, opts ClientOptions) (powervs.PowerVS, error) {
	log := ctrl.LoggerFrom(ctx)

	piOptions := powervs.ServiceOptions{
		IBMPIOptions: &ibmpisession.IBMPIOptions{
			Debug:         opts.Debug,
			Zone:          opts.Zone,
			Authenticator: opts.Authenticator,
		},
		WorkspaceID: opts.WorkspaceID,
	}

	powerVSServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.PowerVS), opts.ServiceEndpoint)
	if powerVSServiceEndpoint != "" {
		log.V(3).Info("Overriding the default PowerVS endpoint", "endpoint", powerVSServiceEndpoint)
		piOptions.URL = powerVSServiceEndpoint
	}

	return powervs.NewService(ctx, piOptions)
}

// GetVPCClient constructs a production VPC client for the given options.
func (b ProdClientBuilder) GetVPCClient(_ context.Context, opts ClientOptions) (vpc.Vpc, error) {
	if opts.Debug {
		core.SetLoggingLevel(core.LevelDebug)
	}

	if opts.VPCRegion == "" {
		return nil, fmt.Errorf("failed to create VPC client: VPC region is not set")
	}

	svcEndpoint := endpoints.FetchVPCEndpoint(opts.VPCRegion, opts.ServiceEndpoint)
	return vpc.NewService(svcEndpoint)
}

// GetTransitGatewayClient constructs a production Transit Gateway client for the given options.
func (b ProdClientBuilder) GetTransitGatewayClient(ctx context.Context, opts ClientOptions) (transitgateway.TransitGateway, error) {
	log := ctrl.LoggerFrom(ctx)

	tgOptions := &tgapiv1.TransitGatewayApisV1Options{
		Authenticator: opts.Authenticator,
	}

	tgServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.TransitGateway), opts.ServiceEndpoint)
	if tgServiceEndpoint != "" {
		log.V(3).Info("Overriding the default TransitGateway endpoint", "endpoint", tgServiceEndpoint)
		tgOptions.URL = tgServiceEndpoint
	}

	return transitgateway.NewService(tgOptions)
}

// GetResourceControllerClient constructs a production Resource Controller client for the given options.
func (b ProdClientBuilder) GetResourceControllerClient(ctx context.Context, opts ClientOptions) (resourcecontroller.ResourceController, error) {
	log := ctrl.LoggerFrom(ctx)

	rcOptions := resourcecontroller.ServiceOptions{
		ResourceControllerV2Options: &resourcecontrollerv2.ResourceControllerV2Options{
			Authenticator: opts.Authenticator,
		},
	}

	rcEndpoint := endpoints.FetchEndpoints(string(endpoints.RC), opts.ServiceEndpoint)
	if rcEndpoint != "" {
		log.V(3).Info("Overriding the default Resource Controller endpoint", "endpoint", rcEndpoint)
		rcOptions.URL = rcEndpoint
	}

	return resourcecontroller.NewService(rcOptions)
}

// GetResourceManagerClient constructs a production Resource Manager client for the given options.
func (b ProdClientBuilder) GetResourceManagerClient(ctx context.Context, opts ClientOptions) (resourcemanager.ResourceManager, error) {
	log := ctrl.LoggerFrom(ctx)

	rmOptions := &resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: opts.Authenticator,
	}

	rmEndpoint := endpoints.FetchEndpoints(string(endpoints.RM), opts.ServiceEndpoint)
	if rmEndpoint != "" {
		log.Info("Overriding the default Resource Manager endpoint:", "endpoint", rmEndpoint)
		rmOptions.URL = rmEndpoint
	}

	return resourcemanager.NewService(rmOptions)
}

// GetCOSClient constructs a production COS client for the given options.
func (b ProdClientBuilder) GetCOSClient(ctx context.Context, opts COSClientOptions) (cos.Cos, error) {
	log := ctrl.LoggerFrom(ctx)

	props, err := authenticator.GetProperties()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service properties: %w", err)
	}
	apiKey := props["APIKEY"]
	if apiKey == "" {
		return nil, fmt.Errorf("IBM Cloud API key is not provided, set IBMCLOUD_API_KEY environmental variable")
	}

	serviceEndpoint := fmt.Sprintf("s3.%s.%s", opts.BucketRegion, cosURLDomain)

	cosServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.COS), opts.ServiceEndpoint)
	if cosServiceEndpoint != "" {
		log.V(3).Info("Overriding the default COS endpoint", "cosEndpoint", cosServiceEndpoint)
		serviceEndpoint = cosServiceEndpoint
	}

	cosOptions := cos.ServiceOptions{
		Options: &cosSession.Options{
			Config: aws.Config{
				Endpoint: ptr.To(serviceEndpoint),
				Region:   ptr.To(opts.BucketRegion),
			},
		},
	}

	return cos.NewService(cosOptions, apiKey, opts.InstanceID)
}

// GetHMACCOSClient constructs a COS client that authenticates via HMAC SigV4 credentials,
// enabling self-authenticating pre-signed URL generation via PresignedURL().
// The HMAC access_key_id / secret_access_key are read from the named Kubernetes Secret.
func (b ProdClientBuilder) GetHMACCOSClient(ctx context.Context, opts HMACCOSClientOptions) (cos.Cos, error) {
	log := ctrl.LoggerFrom(ctx)

	if opts.HMACSecretName == "" {
		return nil, fmt.Errorf("COS HMAC Secret name is not yet populated in cluster status. Waiting for cluster reconciler")
	}

	secret := &corev1.Secret{}
	if err := opts.Client.Get(ctx, types.NamespacedName{
		Namespace: opts.HMACSecretNamespace,
		Name:      opts.HMACSecretName,
	}, secret); err != nil {
		return nil, fmt.Errorf("failed to fetch COS HMAC Secret %q: %w", opts.HMACSecretName, err)
	}

	accessKeyID := string(secret.Data[cosHMACAccessKeyField])
	secretAccessKey := string(secret.Data[cosHMACSecretKeyField])
	if accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("COS HMAC Secret %q is missing %s or %s", opts.HMACSecretName, cosHMACAccessKeyField, cosHMACSecretKeyField)
	}

	serviceEndpoint := fmt.Sprintf("s3.%s.%s", opts.BucketRegion, cosURLDomain)
	cosServiceEndpoint := endpoints.FetchEndpoints(string(endpoints.COS), opts.ServiceEndpoint)
	if cosServiceEndpoint != "" {
		log.V(3).Info("Overriding the default COS endpoint", "cosEndpoint", cosServiceEndpoint)
		serviceEndpoint = cosServiceEndpoint
	}

	cosOptions := cos.ServiceOptions{
		Options: &cosSession.Options{
			Config: aws.Config{
				Endpoint: ptr.To(serviceEndpoint),
				Region:   ptr.To(opts.BucketRegion),
			},
		},
	}

	return cos.NewServiceWithHMAC(cosOptions, accessKeyID, secretAccessKey)
}
