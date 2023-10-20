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
	"path/filepath"

	prismgoclient "github.com/nutanix-cloud-native/prism-go-client"
	"github.com/nutanix-cloud-native/prism-go-client/environment"
	credentialTypes "github.com/nutanix-cloud-native/prism-go-client/environment/credentials"
	kubernetesEnv "github.com/nutanix-cloud-native/prism-go-client/environment/providers/kubernetes"
	envTypes "github.com/nutanix-cloud-native/prism-go-client/environment/types"
	nutanixClientV3 "github.com/nutanix-cloud-native/prism-go-client/v3"
	coreinformers "k8s.io/client-go/informers/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
)

const (
	defaultEndpointPort = "9440"
	ProviderName        = "nutanix"
	configPath          = "/etc/nutanix/config"
	endpointKey         = "prismCentral"
	capxNamespaceKey    = "POD_NAMESPACE"
)

type NutanixClientHelper struct {
	secretInformer    coreinformers.SecretInformer
	configMapInformer coreinformers.ConfigMapInformer
}

func NewNutanixClientHelper(secretInformer coreinformers.SecretInformer, cmInformer coreinformers.ConfigMapInformer) (*NutanixClientHelper, error) {
	return &NutanixClientHelper{
		secretInformer:    secretInformer,
		configMapInformer: cmInformer,
	}, nil
}

func (n *NutanixClientHelper) GetClientFromEnvironment(ctx context.Context, nutanixCluster *infrav1.NutanixCluster) (*nutanixClientV3.Client, error) {
	log := ctrl.LoggerFrom(ctx)
	// Create a list of env providers
	providers := make([]envTypes.Provider, 0)

	// If PrismCentral is set, add the required env provider
	prismCentralInfo := nutanixCluster.Spec.PrismCentral
	if prismCentralInfo != nil {
		if prismCentralInfo.Address == "" {
			return nil, fmt.Errorf("cannot get credentials if Prism Address is not set")
		}
		if prismCentralInfo.Port == 0 {
			return nil, fmt.Errorf("cannot get credentials if Prism Port is not set")
		}
		credentialRef := prismCentralInfo.CredentialRef
		if credentialRef == nil {
			return nil, fmt.Errorf("credentialRef must be set on prismCentral attribute for cluster %s in namespace %s", nutanixCluster.Name, nutanixCluster.Namespace)
		}
		// If namespace is empty, use the cluster namespace
		if credentialRef.Namespace == "" {
			credentialRef.Namespace = nutanixCluster.Namespace
		}
		additionalTrustBundleRef := prismCentralInfo.AdditionalTrustBundle
		if additionalTrustBundleRef != nil &&
			additionalTrustBundleRef.Kind == credentialTypes.NutanixTrustBundleKindConfigMap &&
			additionalTrustBundleRef.Namespace == "" {
			additionalTrustBundleRef.Namespace = nutanixCluster.Namespace
		}
		providers = append(providers, kubernetesEnv.NewProvider(
			*nutanixCluster.Spec.PrismCentral,
			n.secretInformer,
			n.configMapInformer))
	} else {
		log.Info(fmt.Sprintf("[WARNING] prismCentral attribute was not set on NutanixCluster %s in namespace %s. Defaulting to CAPX manager credentials", nutanixCluster.Name, nutanixCluster.Namespace))
	}

	// Add env provider for CAPX manager
	npe, err := n.getManagerNutanixPrismEndpoint()
	if err != nil {
		return nil, err
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
	providers = append(providers, kubernetesEnv.NewProvider(
		*npe,
		n.secretInformer,
		n.configMapInformer))

	// init env with providers
	env := environment.NewEnvironment(
		providers...,
	)
	// fetch endpoint details
	me, err := env.GetManagementEndpoint(envTypes.Topology{})
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

	return n.GetClient(creds, me.AdditionalTrustBundle)
}

func (n *NutanixClientHelper) GetClient(cred prismgoclient.Credentials, additionalTrustBundle string) (*nutanixClientV3.Client, error) {
	if cred.Username == "" {
		return nil, fmt.Errorf("could not create client because username was not set")
	}
	if cred.Password == "" {
		return nil, fmt.Errorf("could not create client because password was not set")
	}
	if cred.Port == "" {
		cred.Port = defaultEndpointPort
	}
	if cred.URL == "" {
		cred.URL = fmt.Sprintf("%s:%s", cred.Endpoint, cred.Port)
	}
	clientOpts := make([]nutanixClientV3.ClientOption, 0)
	if additionalTrustBundle != "" {
		clientOpts = append(clientOpts, nutanixClientV3.WithPEMEncodedCertBundle([]byte(additionalTrustBundle)))
	}
	cli, err := nutanixClientV3.NewV3Client(cred, clientOpts...)
	if err != nil {
		return nil, err
	}

	// Check if the client is working
	_, err = cli.V3.GetCurrentLoggedInUser(context.Background())
	if err != nil {
		return nil, err
	}

	return cli, nil
}

func (n *NutanixClientHelper) getManagerNutanixPrismEndpoint() (*credentialTypes.NutanixPrismEndpoint, error) {
	npe := &credentialTypes.NutanixPrismEndpoint{}
	config, err := n.readEndpointConfig()
	if err != nil {
		return npe, err
	}
	if err = json.Unmarshal(config, npe); err != nil {
		return npe, err
	}
	if npe.CredentialRef == nil {
		return nil, fmt.Errorf("credentialRef must be set on CAPX manager")
	}
	return npe, nil
}

func (n *NutanixClientHelper) readEndpointConfig() ([]byte, error) {
	if b, err := os.ReadFile(filepath.Join(configPath, endpointKey)); err == nil {
		return b, err
	} else if os.IsNotExist(err) {
		return []byte{}, nil
	} else {
		return []byte{}, err
	}
}

func GetCredentialRefForCluster(nutanixCluster *infrav1.NutanixCluster) (*credentialTypes.NutanixCredentialReference, error) {
	if nutanixCluster == nil {
		return nil, fmt.Errorf("cannot get credential reference if nutanix cluster object is nil")
	}
	prismCentralinfo := nutanixCluster.Spec.PrismCentral
	if prismCentralinfo == nil {
		return nil, nil
	}
	if prismCentralinfo.CredentialRef == nil {
		return nil, fmt.Errorf("credentialRef must be set on prismCentral attribute for cluster %s in namespace %s", nutanixCluster.Name, nutanixCluster.Namespace)
	}
	if prismCentralinfo.CredentialRef.Kind != credentialTypes.SecretKind {
		return nil, nil
	}

	return prismCentralinfo.CredentialRef, nil
}
