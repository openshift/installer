/*
Copyright 2020 The ORC Authors.

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

	orcv1alpha1 "github.com/k-orc/openstack-resource-controller/v2/api/v1alpha1"
	clients "github.com/k-orc/openstack-resource-controller/v2/internal/osclients"
	"github.com/k-orc/openstack-resource-controller/v2/internal/version"
)

type providerScopeFactory struct {
	clientCache   *cache.LRUExpireCache
	defaultCACert []byte
}

func (f *providerScopeFactory) NewClientScopeFromObject(ctx context.Context, ctrlClient client.Client, logger logr.Logger, objects ...orcv1alpha1.CloudCredentialsRefProvider) (Scope, error) {
	namespace, credentialsRef := func() (*string, *orcv1alpha1.CloudCredentialsReference) {
		for _, o := range objects {
			namespace, credentialsRef := o.GetCloudCredentialsRef()

			if namespace != nil && credentialsRef != nil {
				return namespace, credentialsRef
			}
		}
		return nil, nil
	}()

	if namespace == nil || credentialsRef == nil {
		return nil, fmt.Errorf("unable to get credentials from provided objects")
	}

	var cloud clientconfig.Cloud
	var caCert []byte

	var err error
	cloud, caCert, err = getCloudFromSecret(ctx, ctrlClient, *namespace, credentialsRef.SecretName, credentialsRef.CloudName)
	if err != nil {
		return nil, err
	}

	if caCert == nil {
		caCert = f.defaultCACert
	}

	if f.clientCache == nil {
		return NewProviderScope(cloud, caCert, logger)
	}

	return NewCachedProviderScope(f.clientCache, cloud, caCert, logger)
}

func getScopeCacheKey(cloud clientconfig.Cloud) (string, error) {
	key, err := computeSpewHash(cloud)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", key), nil
}

type providerScope struct {
	providerClient     *gophercloud.ProviderClient
	providerClientOpts *clientconfig.ClientOpts
}

func NewProviderScope(cloud clientconfig.Cloud, caCert []byte, logger logr.Logger) (Scope, error) {
	providerClient, clientOpts, err := NewProviderClient(cloud, caCert, logger)
	if err != nil {
		return nil, err
	}

	return &providerScope{
		providerClient:     providerClient,
		providerClientOpts: clientOpts,
	}, nil
}

func NewCachedProviderScope(cache *cache.LRUExpireCache, cloud clientconfig.Cloud, caCert []byte, logger logr.Logger) (Scope, error) {
	key, err := getScopeCacheKey(cloud)
	if err != nil {
		return nil, fmt.Errorf("compute cloud config cache key: %w", err)
	}

	if scope, found := cache.Get(key); found {
		logger.V(6).Info("Using scope from cache")
		return scope.(Scope), nil
	}

	scope, err := NewProviderScope(cloud, caCert, logger)
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

func (s *providerScope) NewComputeClient() (clients.ComputeClient, error) {
	return clients.NewComputeClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewNetworkClient() (clients.NetworkClient, error) {
	return clients.NewNetworkClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewImageClient() (clients.ImageClient, error) {
	return clients.NewImageClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewIdentityClient() (clients.IdentityClient, error) {
	return clients.NewIdentityClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewVolumeClient() (clients.VolumeClient, error) {
	return clients.NewVolumeClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) NewVolumeTypeClient() (clients.VolumeTypeClient, error) {
	return clients.NewVolumeTypeClient(s.providerClient, s.providerClientOpts)
}

func (s *providerScope) ExtractToken() (*tokens.Token, error) {
	client, err := openstack.NewIdentityV3(s.providerClient, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, fmt.Errorf("create new identity service client: %w", err)
	}
	return tokens.Get(context.TODO(), client, s.providerClient.Token()).ExtractToken()
}

func NewProviderClient(cloud clientconfig.Cloud, caCert []byte, logger logr.Logger) (*gophercloud.ProviderClient, *clientconfig.ClientOpts, error) {
	clientOpts := new(clientconfig.ClientOpts)

	// We explicitly disable reading auth data from env variables by setting an invalid EnvPrefix.
	// By doing this, we make sure that the data from clouds.yaml is enough to authenticate.
	// For more information: https://github.com/gophercloud/utils/v2/blob/8677e053dcf1f05d0fa0a616094aace04690eb94/openstack/clientconfig/requests.go#L508
	clientOpts.EnvPrefix = "NO_ENV_VARIABLES_"

	if cloud.AuthInfo != nil {
		clientOpts.AuthInfo = cloud.AuthInfo
		clientOpts.AuthType = cloud.AuthType
		clientOpts.RegionName = cloud.RegionName
		clientOpts.EndpointType = cloud.EndpointType
	}

	opts, err := clientconfig.AuthOptions(clientOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("auth option failed for cloud %v: %v", cloud.Cloud, err)
	}
	opts.AllowReauth = true

	provider, err := openstack.NewClient(opts.IdentityEndpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("create providerClient err: %v", err)
	}

	ua := gophercloud.UserAgent{}
	ua.Prepend(fmt.Sprintf("k-orc/%s", version.Get().String()))
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

	provider.HTTPClient.Transport = &RoundTripper{
		RoundTripper: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: config},
	}
	if klog.V(6).Enabled() {
		provider.HTTPClient.Transport = &osclient.RoundTripper{
			Rt:     provider.HTTPClient.Transport,
			Logger: &gophercloudLogger{logger},
		}
	}
	err = openstack.Authenticate(context.TODO(), provider, *opts)
	if err != nil {
		return nil, nil, fmt.Errorf("providerClient authentication err: %v", err)
	}

	return provider, clientOpts, nil
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

	content, ok := secret.Data[orcv1alpha1.CloudCredentialsConfigSecretKey]
	if !ok {
		return emptyCloud, nil, fmt.Errorf("OpenStack credentials secret %v did not contain key %v",
			secretName, orcv1alpha1.CloudCredentialsConfigSecretKey)
	}
	var clouds clientconfig.Clouds
	if err = yaml.Unmarshal(content, &clouds); err != nil {
		return emptyCloud, nil, fmt.Errorf("failed to unmarshal clouds credentials stored in secret %v: %v", secretName, err)
	}

	// get caCert
	caCert, ok := secret.Data[orcv1alpha1.CloudCredencialsCASecretKey]
	if !ok {
		return clouds.Clouds[cloudName], nil, nil
	}

	return clouds.Clouds[cloudName], caCert, nil
}
