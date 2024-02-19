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

package scope

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	osclient "github.com/gophercloud/utils/client"
	"github.com/gophercloud/utils/openstack/clientconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/version"
)

const (
	cloudsSecretKey = "clouds.yaml"
	caSecretKey     = "cacert"
)

type providerScopeFactory struct{}

func (providerScopeFactory) NewClientScopeFromMachine(ctx context.Context, ctrlClient client.Client, openStackMachine *infrav1.OpenStackMachine, defaultCACert []byte, logger logr.Logger) (Scope, error) {
	var cloud clientconfig.Cloud
	var caCert []byte

	if openStackMachine.Spec.IdentityRef != nil {
		var err error
		cloud, caCert, err = getCloudFromSecret(ctx, ctrlClient, openStackMachine.Namespace, openStackMachine.Spec.IdentityRef.Name, openStackMachine.Spec.CloudName)
		if err != nil {
			return nil, err
		}
	}

	if caCert == nil {
		caCert = defaultCACert
	}

	return NewProviderScope(cloud, caCert, logger)
}

func (providerScopeFactory) NewClientScopeFromCluster(ctx context.Context, ctrlClient client.Client, openStackCluster *infrav1.OpenStackCluster, defaultCACert []byte, logger logr.Logger) (Scope, error) {
	var cloud clientconfig.Cloud
	var caCert []byte

	if openStackCluster.Spec.IdentityRef != nil {
		var err error
		cloud, caCert, err = getCloudFromSecret(ctx, ctrlClient, openStackCluster.Namespace, openStackCluster.Spec.IdentityRef.Name, openStackCluster.Spec.CloudName)
		if err != nil {
			return nil, err
		}
	}

	if caCert == nil {
		caCert = defaultCACert
	}

	return NewProviderScope(cloud, caCert, logger)
}

type providerScope struct {
	providerClient     *gophercloud.ProviderClient
	providerClientOpts *clientconfig.ClientOpts
	projectID          string
	logger             logr.Logger
}

func NewProviderScope(cloud clientconfig.Cloud, caCert []byte, logger logr.Logger) (Scope, error) {
	providerClient, clientOpts, projectID, err := NewProviderClient(cloud, caCert, logger)
	if err != nil {
		return nil, err
	}

	return &providerScope{
		providerClient:     providerClient,
		providerClientOpts: clientOpts,
		projectID:          projectID,
		logger:             logger,
	}, nil
}

func (s *providerScope) Logger() logr.Logger {
	return s.logger
}

func (s *providerScope) ProjectID() string {
	return s.projectID
}

func (s *providerScope) NewComputeClient() (clients.ComputeClient, error) {
	return clients.NewComputeClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewNetworkClient() (clients.NetworkClient, error) {
	return clients.NewNetworkClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewVolumeClient() (clients.VolumeClient, error) {
	return clients.NewVolumeClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewImageClient() (clients.ImageClient, error) {
	return clients.NewImageClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewLbClient() (clients.LbClient, error) {
	return clients.NewLbClient(s.providerClient, s.providerClientOpts)
}

func NewProviderClient(cloud clientconfig.Cloud, caCert []byte, logger logr.Logger) (*gophercloud.ProviderClient, *clientconfig.ClientOpts, string, error) {
	clientOpts := new(clientconfig.ClientOpts)
	if cloud.AuthInfo != nil {
		clientOpts.AuthInfo = cloud.AuthInfo
		clientOpts.AuthType = cloud.AuthType
		clientOpts.RegionName = cloud.RegionName
	}

	opts, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return nil, nil, "", fmt.Errorf("auth option failed for cloud %v: %v", cloud.Cloud, err)
	}
	opts.AllowReauth = true

	provider, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return nil, nil, "", fmt.Errorf("create providerClient err: %v", err)
	}

	ua := gophercloud.UserAgent{}
	ua.Prepend(fmt.Sprintf("cluster-api-provider-openstack/%s", version.Get().String()))
	provider.UserAgent = ua

	config := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	if cloud.Verify != nil {
		config.InsecureSkipVerify = !*cloud.Verify
	}
	if caCert != nil {
		config.RootCAs = x509.NewCertPool()
		config.RootCAs.AppendCertsFromPEM(caCert)
	}

	provider.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	if klog.V(6).Enabled() {
		provider.HTTPClient.Transport = &osclient.RoundTripper{
			Rt:     provider.HTTPClient.Transport,
			Logger: &gophercloudLogger{logger},
		}
	}
	err = openstack.Authenticate(provider, *opts)
	if err != nil {
		return nil, nil, "", fmt.Errorf("providerClient authentication err: %v", err)
	}

	projectID, err := getProjectIDFromAuthResult(provider.GetAuthResult())
	if err != nil {
		return nil, nil, "", err
	}

	return provider, clientOpts, projectID, nil
}

type gophercloudLogger struct {
	logger logr.Logger
}

// Printf is a default Printf method.
func (g gophercloudLogger) Printf(format string, args ...interface{}) {
	g.logger.Info(fmt.Sprintf(format, args...))
}

// getCloudFromSecret extract a Cloud from the given namespace:secretName.
func getCloudFromSecret(ctx context.Context, ctrlClient client.Client, secretNamespace string, secretName string, cloudName string) (clientconfig.Cloud, []byte, error) {
	emptyCloud := clientconfig.Cloud{}

	if secretName == "" {
		return emptyCloud, nil, nil
	}

	if cloudName == "" {
		return emptyCloud, nil, fmt.Errorf("secret name set to %v but no cloud was specified. Please set cloud_name in your machine spec", secretName)
	}

	secret := &corev1.Secret{}
	err := ctrlClient.Get(ctx, types.NamespacedName{
		Namespace: secretNamespace,
		Name:      secretName,
	}, secret)
	if err != nil {
		return emptyCloud, nil, err
	}

	content, ok := secret.Data[cloudsSecretKey]
	if !ok {
		return emptyCloud, nil, fmt.Errorf("OpenStack credentials secret %v did not contain key %v",
			secretName, cloudsSecretKey)
	}
	var clouds clientconfig.Clouds
	if err = yaml.Unmarshal(content, &clouds); err != nil {
		return emptyCloud, nil, fmt.Errorf("failed to unmarshal clouds credentials stored in secret %v: %v", secretName, err)
	}

	// get caCert
	caCert, ok := secret.Data[caSecretKey]
	if !ok {
		return clouds.Clouds[cloudName], nil, nil
	}

	return clouds.Clouds[cloudName], caCert, nil
}

// getProjectIDFromAuthResult handles different auth mechanisms to retrieve the
// current project id. Usually we use the Identity v3 Token mechanism that
// returns the project id in the response to the initial auth request.
func getProjectIDFromAuthResult(authResult gophercloud.AuthResult) (string, error) {
	switch authResult := authResult.(type) {
	case tokens.CreateResult:
		project, err := authResult.ExtractProject()
		if err != nil {
			return "", fmt.Errorf("unable to extract project from CreateResult: %v", err)
		}

		return project.ID, nil

	default:
		return "", fmt.Errorf("unable to get the project id from auth response with type %T", authResult)
	}
}
