/*
Copyright 2018 The Kubernetes Authors.

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

package scope

import (
	"context"
	"fmt"
	"time"

	computerest "cloud.google.com/go/compute/apiv1"
	container "cloud.google.com/go/container/apiv1"
	credentials "cloud.google.com/go/iam/credentials/apiv1"
	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"k8s.io/client-go/pkg/version"
	"k8s.io/client-go/util/flowcontrol"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GCPServices contains all the gcp services used by the scopes.
type GCPServices struct {
	Compute *compute.Service
}

// GCPRateLimiter implements cloud.RateLimiter.
type GCPRateLimiter struct{}

// Accept blocks until the operation can be performed.
func (rl *GCPRateLimiter) Accept(ctx context.Context, key *cloud.RateLimitKey) error {
	if key.Operation == "Get" && key.Service == "Operations" {
		// Wait a minimum amount of time regardless of rate limiter.
		rl := &cloud.MinimumRateLimiter{
			// Convert flowcontrol.RateLimiter into cloud.RateLimiter
			RateLimiter: &cloud.AcceptRateLimiter{
				Acceptor: flowcontrol.NewTokenBucketRateLimiter(5, 5), // 5
			},
			Minimum: time.Second,
		}

		return rl.Accept(ctx, key)
	}
	return nil
}

// Observe does nothing.
func (rl *GCPRateLimiter) Observe(context.Context, error, *cloud.RateLimitKey) {
	// noop
}

func newCloud(project string, service GCPServices) cloud.Cloud {
	return cloud.NewGCE(&cloud.Service{
		GA:            service.Compute,
		ProjectRouter: &cloud.SingleProjectRouter{ID: project},
		RateLimiter:   &GCPRateLimiter{},
	})
}

func defaultClientOptions(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client) ([]option.ClientOption, error) {
	opts := []option.ClientOption{
		option.WithUserAgent(fmt.Sprintf("gcp.cluster.x-k8s.io/%s", version.Get())),
	}

	if credentialsRef != nil {
		rawData, err := getCredentialDataFromRef(ctx, credentialsRef, crClient)
		if err != nil {
			return nil, fmt.Errorf("getting gcp credentials from reference %s: %w", credentialsRef, err)
		}
		opts = append(opts, option.WithCredentialsJSON(rawData))
	}

	return opts, nil
}

func newComputeService(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, endpoints *infrav1.ServiceEndpoints) (*compute.Service, error) {
	opts, err := defaultClientOptions(ctx, credentialsRef, crClient)
	if err != nil {
		return nil, fmt.Errorf("getting default gcp client options: %w", err)
	}

	if endpoints != nil && endpoints.ComputeServiceEndpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoints.ComputeServiceEndpoint))
	}

	computeSvc, err := compute.NewService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating new compute service instance: %w", err)
	}

	return computeSvc, nil
}

func newClusterManagerClient(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, endpoints *infrav1.ServiceEndpoints) (*container.ClusterManagerClient, error) {
	opts, err := defaultClientOptions(ctx, credentialsRef, crClient)
	if err != nil {
		return nil, fmt.Errorf("getting default gcp client options: %w", err)
	}

	if endpoints != nil && endpoints.ContainerServiceEndpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoints.ContainerServiceEndpoint))
	}

	managedClusterClient, err := container.NewClusterManagerClient(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf("failed to create gcp cluster manager client: %v", err)
	}

	return managedClusterClient, nil
}

func newIamCredentialsClient(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, endpoints *infrav1.ServiceEndpoints) (*credentials.IamCredentialsClient, error) {
	opts, err := defaultClientOptions(ctx, credentialsRef, crClient)
	if err != nil {
		return nil, fmt.Errorf("getting default gcp client options: %w", err)
	}

	if endpoints != nil && endpoints.IAMServiceEndpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoints.IAMServiceEndpoint))
	}

	credentialsClient, err := credentials.NewIamCredentialsClient(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf("failed to create gcp ciam credentials client: %v", err)
	}

	return credentialsClient, nil
}

func newInstanceGroupManagerClient(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, endpoints *infrav1.ServiceEndpoints) (*computerest.InstanceGroupManagersClient, error) {
	opts, err := defaultClientOptions(ctx, credentialsRef, crClient)
	if err != nil {
		return nil, fmt.Errorf("getting default gcp client options: %w", err)
	}

	if endpoints != nil && endpoints.ComputeServiceEndpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoints.ComputeServiceEndpoint))
	}

	instanceGroupManagersClient, err := computerest.NewInstanceGroupManagersRESTClient(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf("failed to create gcp instance group managers rest client: %v", err)
	}

	return instanceGroupManagersClient, nil
}

func newTagBindingsClient(ctx context.Context, credentialsRef *infrav1.ObjectReference, crClient client.Client, location string, endpoints *infrav1.ServiceEndpoints) (*resourcemanager.TagBindingsClient, error) {
	opts, err := defaultClientOptions(ctx, credentialsRef, crClient)

	if endpoints != nil && endpoints.ResourceManagerServiceEndpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoints.ResourceManagerServiceEndpoint))
	} else {
		endpoint := location + "-cloudresourcemanager.googleapis.com:443"
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	if err != nil {
		return nil, fmt.Errorf("getting default gcp client options: %w", err)
	}

	client, err := resourcemanager.NewTagBindingsClient(ctx, opts...)
	if err != nil {
		return nil, errors.Errorf("failed to create gcp tag binding client: %v", err)
	}

	return client, nil
}
