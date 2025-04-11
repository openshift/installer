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
	"time"

	"github.com/go-logr/logr"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	osclient "github.com/gophercloud/utils/v2/client"
	"github.com/gophercloud/utils/v2/openstack/clientconfig"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/clients"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/hash"
	"sigs.k8s.io/cluster-api-provider-openstack/version"
)

const (
	CloudsSecretKey = "clouds.yaml"
	CASecretKey     = "cacert"
)

type providerScopeFactory struct {
	clientCache *cache.LRUExpireCache
}

func (f *providerScopeFactory) NewClientScopeFromObject(ctx context.Context, ctrlClient client.Client, defaultCACert []byte, logger logr.Logger, objects ...infrav1.IdentityRefProvider) (Scope, error) {
	var namespace *string
	var identityRef *infrav1.OpenStackIdentityReference

	for _, o := range objects {
		namespace, identityRef = o.GetIdentityRef()
		if namespace != nil || identityRef != nil {
			break
		}
	}

	if namespace == nil || identityRef == nil {
		return nil, fmt.Errorf("unable to get identityRef from provided objects")
	}

	var cloud clientconfig.Cloud
	var caCert []byte

	var err error
	cloud, caCert, err = getCloudFromSecret(ctx, ctrlClient, *namespace, identityRef.Name, identityRef.CloudName)
	if err != nil {
		return nil, err
	}

	if caCert == nil {
		caCert = defaultCACert
	}

	if f.clientCache == nil {
		return NewProviderScope(cloud, identityRef.Region, caCert, logger)
	}

	return NewCachedProviderScope(f.clientCache, cloud, identityRef.Region, caCert, logger)
}

func getScopeCacheKey(cloud clientconfig.Cloud) (string, error) {
	key, err := hash.ComputeSpewHash(cloud)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", key), nil
}

type providerScope struct {
	providerClient     *gophercloud.ProviderClient
	providerClientOpts *clientconfig.ClientOpts
	projectID          string
}

func NewProviderScope(cloud clientconfig.Cloud, regionName string, caCert []byte, logger logr.Logger) (Scope, error) {
	providerClient, clientOpts, projectID, err := NewProviderClient(cloud, regionName, caCert, logger)
	if err != nil {
		return nil, err
	}

	return &providerScope{
		providerClient:     providerClient,
		providerClientOpts: clientOpts,
		projectID:          projectID,
	}, nil
}

func NewCachedProviderScope(cache *cache.LRUExpireCache, cloud clientconfig.Cloud, regionName string, caCert []byte, logger logr.Logger) (Scope, error) {
	key, err := getScopeCacheKey(cloud)
	if err != nil {
		return nil, fmt.Errorf("compute cloud config cache key: %w", err)
	}

	if scope, found := cache.Get(key); found {
		logger.V(6).Info("Using scope from cache")
		return scope.(Scope), nil
	}

	scope, err := NewProviderScope(cloud, regionName, caCert, logger)
	if err != nil {
		return nil, err
	}

	token, err := scope.ExtractToken()
	if err != nil {
		return nil, err
	}

	// compute the token expiration time
	expiry := time.Until(token.ExpiresAt) / 2

	cache.Add(key, scope, expiry)
	return scope, nil
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

func (s *providerScope) ExtractToken() (*tokens.Token, error) {
	client, err := openstack.NewIdentityV3(s.providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, fmt.Errorf("create new identity service client: %w", err)
	}
	return tokens.Get(context.TODO(), client, s.providerClient.Token()).ExtractToken()
}

func NewProviderClient(cloud clientconfig.Cloud, regionName string, caCert []byte, logger logr.Logger) (*gophercloud.ProviderClient, *clientconfig.ClientOpts, string, error) {
	clientOpts := new(clientconfig.ClientOpts)

	// We explicitly disable reading auth data from env variables by setting an invalid EnvPrefix.
	// By doing this, we make sure that the data from clouds.yaml is enough to authenticate.
	// For more information: https://github.com/gophercloud/utils/v2/blob/8677e053dcf1f05d0fa0a616094aace04690eb94/openstack/clientconfig/requests.go#L508
	clientOpts.EnvPrefix = "NO_ENV_VARIABLES_"
	if regionName == "" {
		regionName = cloud.RegionName
	}
	if cloud.AuthInfo != nil {
		clientOpts.AuthInfo = cloud.AuthInfo
		clientOpts.AuthType = cloud.AuthType
		clientOpts.RegionName = regionName
		clientOpts.EndpointType = cloud.EndpointType
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
		ok := config.RootCAs.AppendCertsFromPEM(caCert)
		if !ok {
			// If no certificates were successfully parsed, set RootCAs to nil to use the host's root CA
			config.RootCAs = nil
		}
	}

	provider.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config}
	if klog.V(6).Enabled() {
		provider.HTTPClient.Transport = &osclient.RoundTripper{
			Rt:     provider.HTTPClient.Transport,
			Logger: &gophercloudLogger{logger},
		}
	}
	err = openstack.Authenticate(context.TODO(), provider, *opts)
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

	content, ok := secret.Data[CloudsSecretKey]
	if !ok {
		return emptyCloud, nil, fmt.Errorf("OpenStack credentials secret %v did not contain key %v",
			secretName, CloudsSecretKey)
	}
	var clouds clientconfig.Clouds
	if err = yaml.Unmarshal(content, &clouds); err != nil {
		return emptyCloud, nil, fmt.Errorf("failed to unmarshal clouds credentials stored in secret %v: %v", secretName, err)
	}

	// get caCert
	caCert, ok := secret.Data[CASecretKey]
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
