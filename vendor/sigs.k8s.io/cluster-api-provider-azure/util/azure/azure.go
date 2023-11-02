/*
Copyright 2022 The Kubernetes Authors.

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

package azure

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	azureautorest "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/jongio/azidext/go/azidext"
	"github.com/pkg/errors"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AzureSystemNodeLabelPrefix is a standard node label prefix for Azure features, e.g., kubernetes.azure.com/scalesetpriority.
const AzureSystemNodeLabelPrefix = "kubernetes.azure.com"

const (
	// ProviderIDPrefix will be appended to the beginning of Azure resource IDs to form the Kubernetes Provider ID.
	// NOTE: this format matches the 2 slashes format used in cloud-provider and cluster-autoscaler.
	ProviderIDPrefix = "azure://"
)

// IsAzureSystemNodeLabelKey is a helper function that determines whether a node label key is an Azure "system" label.
func IsAzureSystemNodeLabelKey(labelKey string) bool {
	return strings.HasPrefix(labelKey, AzureSystemNodeLabelPrefix)
}

func getCloudConfig(environment azureautorest.Environment) cloud.Configuration {
	var config cloud.Configuration
	switch environment.Name {
	case "AzureStackCloud":
		config = cloud.Configuration{
			ActiveDirectoryAuthorityHost: environment.ActiveDirectoryEndpoint,
			Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
				cloud.ResourceManager: {
					Audience: environment.TokenAudience,
					Endpoint: environment.ResourceManagerEndpoint,
				},
			},
		}
	case "AzureChinaCloud":
		config = cloud.AzureChina
	case "AzureUSGovernmentCloud":
		config = cloud.AzureGovernment
	default:
		config = cloud.AzurePublic
	}
	return config
}

// GetAuthorizer returns an autorest.Authorizer-compatible object from MSAL.
func GetAuthorizer(settings auth.EnvironmentSettings) (autorest.Authorizer, error) {
	// azidentity uses different envvars for certificate authentication:
	//  azidentity: AZURE_CLIENT_CERTIFICATE_{PATH,PASSWORD}
	//  autorest: AZURE_CERTIFICATE_{PATH,PASSWORD}
	// Let's set them according to the envvars used by autorest, in case they are present
	_, azidSet := os.LookupEnv("AZURE_CLIENT_CERTIFICATE_PATH")
	path, autorestSet := os.LookupEnv("AZURE_CERTIFICATE_PATH")
	if !azidSet && autorestSet {
		os.Setenv("AZURE_CLIENT_CERTIFICATE_PATH", path)
		os.Setenv("AZURE_CLIENT_CERTIFICATE_PASSWORD", os.Getenv("AZURE_CERTIFICATE_PASSWORD"))
	}

	options := azidentity.DefaultAzureCredentialOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: getCloudConfig(settings.Environment),
		},
	}
	cred, err := azidentity.NewDefaultAzureCredential(&options)
	if err != nil {
		return nil, err
	}

	// We must use TokenAudience for StackCloud, otherwise we get an
	// AADSTS500011 error from the API
	scope := settings.Environment.TokenAudience
	if !strings.HasSuffix(scope, "/.default") {
		scope += "/.default"
	}
	return azidext.NewTokenCredentialAdapter(cred, []string{scope}), nil
}

// FindParentMachinePool finds the parent MachinePool for the AzureMachinePool.
func FindParentMachinePool(ampName string, cli client.Client) (*expv1.MachinePool, error) {
	ctx := context.Background()
	machinePoolList := &expv1.MachinePoolList{}
	if err := cli.List(ctx, machinePoolList); err != nil {
		return nil, errors.Wrapf(err, "failed to list MachinePools for %s", ampName)
	}
	for _, mp := range machinePoolList.Items {
		if mp.Spec.Template.Spec.InfrastructureRef.Name == ampName {
			return &mp, nil
		}
	}
	return nil, errors.Errorf("failed to get MachinePool for %s", ampName)
}

// FindParentMachinePoolWithRetry finds the parent MachinePool for the AzureMachinePool with retry.
func FindParentMachinePoolWithRetry(ampName string, cli client.Client, maxAttempts int) (*expv1.MachinePool, error) {
	for i := 1; ; i++ {
		p, err := FindParentMachinePool(ampName, cli)
		if err != nil {
			if i >= maxAttempts {
				return nil, errors.Wrap(err, "failed to find parent MachinePool")
			}
			time.Sleep(1 * time.Second)
			continue
		}
		return p, nil
	}
}

// ParseResourceID parses a string to an *arm.ResourceID, first removing any "azure://" prefix.
func ParseResourceID(id string) (*arm.ResourceID, error) {
	return arm.ParseResourceID(strings.TrimPrefix(id, ProviderIDPrefix))
}
