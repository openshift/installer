/*
Copyright 2022 Nutanix

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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	prismgoclient "github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/environment"
	"github.com/nutanix-cloud-native/prism-go-client/environment/credentials"
	kubernetesEnv "github.com/nutanix-cloud-native/prism-go-client/environment/providers/kubernetes"
	envTypes "github.com/nutanix-cloud-native/prism-go-client/environment/types"
	nutanixClientV3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	coreinformers "k8s.io/client-go/informers/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
)

const (
	defaultEndpointPort = "9440"

	capxNamespaceKey = "POD_NAMESPACE"
)

const (
	configPath = "/etc/nutanix/config/prismCentral"
)

var (
	ErrPrismAddressNotSet = fmt.Errorf("cannot get credentials if Prism Address is not set")
	ErrPrismPortNotSet    = fmt.Errorf("cannot get credentials if Prism Port is not set")

	ErrPrismIUsernameNotSet = fmt.Errorf("could not create client because username was not set")
	ErrPrismIPasswordNotSet = fmt.Errorf("could not create client because password was not set")

	ErrCredentialRefNotSet = fmt.Errorf("credentialRef must be set on CAPX manager")
)

type NutanixClientHelper struct {
	secretInformer    coreinformers.SecretInformer
	configMapInformer coreinformers.ConfigMapInformer

	managerNutanixPrismEndpointReader func() (*credentials.NutanixPrismEndpoint, error)
}

func NewHelper(secretInformer coreinformers.SecretInformer, cmInformer coreinformers.ConfigMapInformer) *NutanixClientHelper {
	return &NutanixClientHelper{
		secretInformer:                    secretInformer,
		configMapInformer:                 cmInformer,
		managerNutanixPrismEndpointReader: readManagerNutanixPrismEndpointFromDefaultFile,
	}
}

func (n *NutanixClientHelper) withCustomNutanixPrismEndpointReader(getter func() (*credentials.NutanixPrismEndpoint, error)) *NutanixClientHelper {
	n.managerNutanixPrismEndpointReader = getter
	return n
}

// BuildClientForNutanixClusterWithFallback builds a Nutanix Client from the information provided in nutanixCluster.
func (n *NutanixClientHelper) BuildClientForNutanixClusterWithFallback(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) (*nutanixClientV3.Client, error) {
	me, err := n.buildManagementEndpoint(ctx, nutanixCluster)
	if err != nil {
		return nil, err
	}
	creds := prismgoclient.Credentials{
		URL:      me.Address.Host,
		Endpoint: me.Address.Host,
		Insecure: me.Insecure,
		Username: me.ApiCredentials.Username,
		Password: me.ApiCredentials.Password,
	}
	return Build(creds, me.AdditionalTrustBundle)
}

// buildManagementEndpoint takes in a NutanixCluster and constructs a ManagementEndpoint with all the information provided.
// If required information is not set, it will fallback to using information from /etc/nutanix/config/prismCentral,
// which is expected to be mounted in the Pod.
func (n *NutanixClientHelper) buildManagementEndpoint(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) (*envTypes.ManagementEndpoint, error) {
	log := ctrl.LoggerFrom(ctx)

	// Create an empty list of env providers
	providers := make([]envTypes.Provider, 0)

	// Attempt to build a provider from the NutanixCluster object
	providerForNutanixCluster, err := n.buildProviderFromNutanixCluster(nutanixCluster)
	if err != nil {
		return nil, fmt.Errorf("error building an environment provider from NutanixCluster: %w", err)
	}
	if providerForNutanixCluster != nil {
		providers = append(providers, providerForNutanixCluster)
	} else {
		log.Info(fmt.Sprintf("[WARNING] prismCentral attribute was not set on NutanixCluster %s in namespace %s. Defaulting to CAPX manager credentials", nutanixCluster.Name, nutanixCluster.Namespace))
		// Fallback to building a provider using prism central information from the CAPX management cluster
		// using information from /etc/nutanix/config/prismCentral
		providerForLocalFile, err := n.buildProviderFromFile()
		if err != nil {
			return nil, fmt.Errorf("error building an environment provider from file: %w", err)
		}
		if providerForLocalFile != nil {
			providers = append(providers, providerForLocalFile)
		}
	}

	// Initialize environment with providers
	env := environment.NewEnvironment(providers...)
	// GetManagementEndpoint will return the first valid endpoint from the list of providers
	me, err := env.GetManagementEndpoint(envTypes.Topology{})
	if err != nil {
		return nil, fmt.Errorf("failed to get management endpoint object: %w", err)
	}
	return me, nil
}

// buildProviderFromNutanixCluster will return an envTypes.Provider with info from the provided NutanixCluster.
// It will return nil if nutanixCluster.Spec.PrismCentral is nil.
// It will return an error if required information is missing.
func (n *NutanixClientHelper) buildProviderFromNutanixCluster(nutanixCluster *infrav1.NutanixCluster) (envTypes.Provider, error) {
	prismCentralInfo := nutanixCluster.Spec.PrismCentral
	if prismCentralInfo == nil {
		return nil, nil
	}

	// PrismCentral is set, build a provider and fixup missing information
	if prismCentralInfo.Address == "" {
		return nil, ErrPrismAddressNotSet
	}
	if prismCentralInfo.Port == 0 {
		return nil, ErrPrismPortNotSet
	}
	credentialRef, err := nutanixCluster.GetPrismCentralCredentialRef()
	if err != nil {
		//nolint:wrapcheck // error is already wrapped
		return nil, err
	}
	// If namespace is empty, use the cluster namespace
	if credentialRef.Namespace == "" {
		credentialRef.Namespace = nutanixCluster.Namespace
	}
	additionalTrustBundleRef := prismCentralInfo.AdditionalTrustBundle
	if additionalTrustBundleRef != nil &&
		additionalTrustBundleRef.Kind == credentials.NutanixTrustBundleKindConfigMap &&
		additionalTrustBundleRef.Namespace == "" {
		additionalTrustBundleRef.Namespace = nutanixCluster.Namespace
	}

	return kubernetesEnv.NewProvider(*prismCentralInfo, n.secretInformer, n.configMapInformer), nil
}

// buildProviderFromFile will return an envTypes.Provider with info from the provided file.
// It will return an error if required information is missing.
func (n *NutanixClientHelper) buildProviderFromFile() (envTypes.Provider, error) {
	npe, err := n.managerNutanixPrismEndpointReader()
	if err != nil {
		return nil, fmt.Errorf("failed to create prism endpoint: %w", err)
	}
	// If namespaces is not set, set it to the namespace of the CAPX manager
	if npe.CredentialRef.Namespace == "" {
		capxNamespace := os.Getenv(capxNamespaceKey)
		if capxNamespace == "" {
			return nil, fmt.Errorf("failed to retrieve capx-namespace. Make sure %s env variable is set", capxNamespaceKey)
		}
		npe.CredentialRef.Namespace = capxNamespace
	}
	if npe.AdditionalTrustBundle != nil && npe.AdditionalTrustBundle.Namespace == "" {
		capxNamespace := os.Getenv(capxNamespaceKey)
		if capxNamespace == "" {
			return nil, fmt.Errorf("failed to retrieve capx-namespace. Make sure %s env variable is set", capxNamespaceKey)
		}
		npe.AdditionalTrustBundle.Namespace = capxNamespace
	}

	return kubernetesEnv.NewProvider(*npe, n.secretInformer, n.configMapInformer), nil
}

func Build(creds prismgoclient.Credentials, additionalTrustBundle string) (*nutanixClientV3.Client, error) {
	cli, err := buildClientFromCredentials(creds, additionalTrustBundle)
	if err != nil {
		return nil, err
	}
	// Check if the client is working
	_, err = cli.V3.GetCurrentLoggedInUser(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get current logged in user with client: %w", err)
	}
	return cli, nil
}

func buildClientFromCredentials(creds prismgoclient.Credentials, additionalTrustBundle string) (*nutanixClientV3.Client, error) {
	if creds.Username == "" {
		return nil, ErrPrismIUsernameNotSet
	}
	if creds.Password == "" {
		return nil, ErrPrismIPasswordNotSet
	}
	if creds.Port == "" {
		creds.Port = defaultEndpointPort
	}
	if creds.URL == "" {
		creds.URL = fmt.Sprintf("%s:%s", creds.Endpoint, creds.Port)
	}

	clientOpts := make([]nutanixClientV3.ClientOption, 0)
	if additionalTrustBundle != "" {
		clientOpts = append(clientOpts, nutanixClientV3.WithPEMEncodedCertBundle([]byte(additionalTrustBundle)))
	}
	// Build the client with the creds and possibly an additional TrustBundle
	cli, err := nutanixClientV3.NewV3Client(creds, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new nutanix client: %w", err)
	}

	return cli, nil
}

// readManagerNutanixPrismEndpoint reads the default config file and unmarshalls it into NutanixPrismEndpoint.
// Returns an error if the file does not exist and other read or unmarshalling errors.
func readManagerNutanixPrismEndpointFromDefaultFile() (*credentials.NutanixPrismEndpoint, error) {
	return readManagerNutanixPrismEndpointFromFile(configPath)
}

// this function is primarily here to make writing unit tests simpler
// readManagerNutanixPrismEndpointFromDefaultFile should be used outside of tests
func readManagerNutanixPrismEndpointFromFile(configFile string) (*credentials.NutanixPrismEndpoint, error) {
	// fail on all errors including NotExist error
	config, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read prism config in manager: %w", err)
	}
	npe := &credentials.NutanixPrismEndpoint{}
	if err = json.Unmarshal(config, npe); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	if npe.CredentialRef == nil {
		return nil, ErrCredentialRefNotSet
	}
	return npe, nil
}
