package azure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v3"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/msi/armmsi"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	ignutil "github.com/coreos/ignition/v2/config/util"
	ignitionconfig "github.com/coreos/ignition/v2/config/v3_2"
	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/credentialsrequest"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	aztypes "github.com/openshift/installer/pkg/types/azure"
)

// WIFIdentity stores the Azure managed identity info for a credential request.
type WIFIdentity struct {
	ClientID    string
	PrincipalID string
	ResourceID  string
	Name        string
}

// provisionWIF creates the OIDC infrastructure, managed identities,
// federated identity credentials, and role assignments for Azure WIF.
func (p *Provider) provisionWIF(ctx context.Context, in clusterapi.PreProvisionInput) error {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return fmt.Errorf("failed to get Azure session: %w", err)
	}

	platform := in.InstallConfig.Config.Platform.Azure
	infraID := in.InfraID
	subscriptionID := session.Credentials.SubscriptionID
	region := platform.Region
	resourceGroupName := platform.ClusterResourceGroupName(infraID)

	clientOpts := &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}

	tags := map[string]*string{
		fmt.Sprintf("kubernetes.io_cluster.%s", infraID): to.Ptr("owned"),
	}

	// Ensure the resource group exists.
	rgClient, err := armresources.NewResourceGroupsClient(subscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create resource groups client: %w", err)
	}
	_, err = rgClient.CreateOrUpdate(ctx, resourceGroupName, armresources.ResourceGroup{
		Location: to.Ptr(region),
		Tags:     tags,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to ensure resource group %s: %w", resourceGroupName, err)
	}

	if err := p.createOIDCStorage(ctx, session, infraID, resourceGroupName, region, tags, clientOpts, in.BoundSASigningKey); err != nil {
		return err
	}

	if err := p.createManagedIdentities(ctx, session.TokenCreds, subscriptionID, resourceGroupName, region, infraID, tags, clientOpts, in.CredentialsRequests.Requests); err != nil {
		return err
	}

	if err := p.createFederatedCredentials(ctx, session.TokenCreds, subscriptionID, resourceGroupName, clientOpts, in.CredentialsRequests.Requests); err != nil {
		return err
	}

	return p.createRoleAssignments(ctx, session.TokenCreds, subscriptionID, resourceGroupName, infraID, clientOpts, in.CredentialsRequests.Requests)
}

func (p *Provider) createOIDCStorage(ctx context.Context, session *azconfig.Session, infraID, resourceGroupName, region string, tags map[string]*string, clientOpts *arm.ClientOptions, signingKey *tls.BoundSASigningKey) error {
	storageAccountName := aztypes.WIFStorageAccountName(infraID)

	storageClientFactory, err := armstorage.NewClientFactory(session.Credentials.SubscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create storage client factory: %w", err)
	}

	logrus.Infof("Creating WIF OIDC storage account %s", storageAccountName)
	saPoller, err := storageClientFactory.NewAccountsClient().BeginCreate(ctx, resourceGroupName, storageAccountName,
		armstorage.AccountCreateParameters{
			Kind:     to.Ptr(armstorage.KindStorageV2),
			Location: to.Ptr(region),
			SKU: &armstorage.SKU{
				Name: to.Ptr(armstorage.SKUNameStandardLRS),
			},
			Properties: &armstorage.AccountPropertiesCreateParameters{
				AllowBlobPublicAccess: to.Ptr(true),
				AllowSharedKeyAccess:  to.Ptr(true),
				MinimumTLSVersion:     to.Ptr(armstorage.MinimumTLSVersionTLS12),
			},
			Tags: tags,
		}, nil)
	if err != nil {
		return fmt.Errorf("failed to create WIF storage account: %w", err)
	}
	if _, err = saPoller.PollUntilDone(ctx, nil); err != nil {
		return fmt.Errorf("failed waiting for WIF storage account creation: %w", err)
	}

	containerName := aztypes.WIFContainerName
	logrus.Infof("Creating WIF blob container %s", containerName)
	_, err = storageClientFactory.NewBlobContainersClient().Create(ctx, resourceGroupName, storageAccountName, containerName,
		armstorage.BlobContainer{
			ContainerProperties: &armstorage.ContainerProperties{
				PublicAccess: to.Ptr(armstorage.PublicAccessBlob),
			},
		}, nil)
	if err != nil {
		return fmt.Errorf("failed to create WIF blob container: %w", err)
	}

	pubKeyData := signingKey.PublicKeyData()
	if pubKeyData == nil {
		return fmt.Errorf("SA signing public key not available for WIF provisioning")
	}
	saSigningPubKey, err := tls.PemToPublicKey(pubKeyData)
	if err != nil {
		return fmt.Errorf("failed to parse SA signing public key: %w", err)
	}

	jwksBytes, keyID, err := tls.BuildJSONWebKeySet(saSigningPubKey)
	if err != nil {
		return fmt.Errorf("failed to build JWKS: %w", err)
	}
	logrus.Debugf("Built JWKS with key-id=%s", keyID)

	storageSuffix := session.Environment.StorageEndpointSuffix
	issuerURL := fmt.Sprintf("https://%s.blob.%s/%s", storageAccountName, storageSuffix, containerName)
	discoveryDoc := tls.BuildDiscoveryDocument(issuerURL)

	blobBaseURL := fmt.Sprintf("https://%s.blob.%s", storageAccountName, storageSuffix)
	if err := uploadBlob(ctx, blobBaseURL, containerName, ".well-known/openid-configuration", []byte(discoveryDoc), session.TokenCreds); err != nil {
		return fmt.Errorf("failed to upload OIDC discovery document: %w", err)
	}
	if err := uploadBlob(ctx, blobBaseURL, containerName, "keys.json", jwksBytes, session.TokenCreds); err != nil {
		return fmt.Errorf("failed to upload JWKS: %w", err)
	}
	logrus.Infof("Uploaded OIDC discovery document and JWKS to %s", issuerURL)

	p.wifIssuerURL = issuerURL
	return nil
}

func (p *Provider) createManagedIdentities(ctx context.Context, cred azcore.TokenCredential, subscriptionID, resourceGroupName, region, infraID string, tags map[string]*string, clientOpts *arm.ClientOptions, credReqs []credentialsrequest.CredentialRequest) error {
	msiClientFactory, err := armmsi.NewClientFactory(subscriptionID, cred, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create MSI client factory: %w", err)
	}

	p.wifIdentities = make(map[string]WIFIdentity, len(credReqs))
	for _, cr := range credReqs {
		identityName := aztypes.WIFManagedIdentityName(infraID, cr.SecretRefNamespace, cr.SecretRefName)
		key := cr.SecretRefNamespace + "/" + cr.SecretRefName
		logrus.Infof("Creating managed identity %s for %s", identityName, key)

		_, err := msiClientFactory.NewUserAssignedIdentitiesClient().CreateOrUpdate(ctx, resourceGroupName, identityName,
			armmsi.Identity{
				Location: to.Ptr(region),
				Tags:     tags,
			}, nil)
		if err != nil {
			return fmt.Errorf("failed to create managed identity %s: %w", identityName, err)
		}

		identity, err := msiClientFactory.NewUserAssignedIdentitiesClient().Get(ctx, resourceGroupName, identityName, nil)
		if err != nil {
			return fmt.Errorf("failed to get managed identity %s: %w", identityName, err)
		}

		p.wifIdentities[key] = WIFIdentity{
			ClientID:    *identity.Properties.ClientID,
			PrincipalID: *identity.Properties.PrincipalID,
			ResourceID:  *identity.ID,
			Name:        identityName,
		}
		logrus.Debugf("Managed identity %s: clientID=%s", identityName, *identity.Properties.ClientID)
	}
	return nil
}

func (p *Provider) createFederatedCredentials(ctx context.Context, cred azcore.TokenCredential, subscriptionID, resourceGroupName string, clientOpts *arm.ClientOptions, credReqs []credentialsrequest.CredentialRequest) error {
	msiClientFactory, err := armmsi.NewClientFactory(subscriptionID, cred, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create MSI client factory: %w", err)
	}
	ficClient := msiClientFactory.NewFederatedIdentityCredentialsClient()

	for _, cr := range credReqs {
		key := cr.SecretRefNamespace + "/" + cr.SecretRefName
		identity, ok := p.wifIdentities[key]
		if !ok {
			return fmt.Errorf("no managed identity found for %s", key)
		}

		for _, saName := range cr.ServiceAccountNames {
			subject := fmt.Sprintf("system:serviceaccount:%s:%s", cr.SecretRefNamespace, saName)
			ficName := sanitizeFICName(fmt.Sprintf("%s-%s", identity.Name, saName))
			logrus.Infof("Creating federated identity credential %s (subject=%s)", ficName, subject)

			_, err := ficClient.CreateOrUpdate(ctx, resourceGroupName, identity.Name, ficName,
				armmsi.FederatedIdentityCredential{
					Properties: &armmsi.FederatedIdentityCredentialProperties{
						Issuer:    to.Ptr(p.wifIssuerURL),
						Subject:   to.Ptr(subject),
						Audiences: []*string{to.Ptr("openshift")},
					},
				}, nil)
			if err != nil {
				return fmt.Errorf("failed to create FIC %s: %w", ficName, err)
			}
		}
	}
	return nil
}

func (p *Provider) createRoleAssignments(ctx context.Context, cred azcore.TokenCredential, subscriptionID, resourceGroupName, infraID string, clientOpts *arm.ClientOptions, credReqs []credentialsrequest.CredentialRequest) error {
	authClientFactory, err := armauthorization.NewClientFactory(subscriptionID, cred, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create authorization client factory: %w", err)
	}

	roleDefClient := authClientFactory.NewRoleDefinitionsClient()
	roleAssignClient := authClientFactory.NewRoleAssignmentsClient()
	rgScope := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, resourceGroupName)

	for _, cr := range credReqs {
		key := cr.SecretRefNamespace + "/" + cr.SecretRefName
		identity, ok := p.wifIdentities[key]
		if !ok {
			return fmt.Errorf("no managed identity found for %s", key)
		}

		azureSpec, ok := cr.ProviderSpec.(*credentialsrequest.AzureProviderSpec)
		if !ok {
			return fmt.Errorf("credential request %s does not have an Azure provider spec", key)
		}

		// Handle built-in role bindings.
		for _, rb := range azureSpec.RoleBindings {
			roleDefID, err := findRoleDefinitionByName(ctx, roleDefClient, rgScope, rb.Role)
			if err != nil {
				return fmt.Errorf("failed to find role %s for %s: %w", rb.Role, key, err)
			}
			if err := assignRole(ctx, roleAssignClient, rgScope, identity.PrincipalID, roleDefID); err != nil {
				return fmt.Errorf("failed to assign role %s to %s: %w", rb.Role, key, err)
			}
		}

		// Handle fine-grained permissions via custom role definition.
		if len(azureSpec.Permissions) > 0 || len(azureSpec.DataPermissions) > 0 {
			customRoleName := fmt.Sprintf("%s-%s-%s", infraID, cr.SecretRefNamespace, cr.SecretRefName)
			if len(customRoleName) > 64 {
				customRoleName = customRoleName[:64]
			}
			customRoleID := uuid.New().String()
			subScope := fmt.Sprintf("/subscriptions/%s", subscriptionID)

			actions := make([]*string, len(azureSpec.Permissions))
			for i, p := range azureSpec.Permissions {
				actions[i] = to.Ptr(p)
			}
			dataActions := make([]*string, len(azureSpec.DataPermissions))
			for i, d := range azureSpec.DataPermissions {
				dataActions[i] = to.Ptr(d)
			}

			logrus.Infof("Creating custom role definition %s for %s", customRoleName, key)
			_, err := roleDefClient.CreateOrUpdate(ctx, subScope, customRoleID,
				armauthorization.RoleDefinition{
					Properties: &armauthorization.RoleDefinitionProperties{
						RoleName:         to.Ptr(customRoleName),
						Description:      to.Ptr(fmt.Sprintf("Custom role for %s WIF", key)),
						RoleType:         to.Ptr("CustomRole"),
						AssignableScopes: []*string{to.Ptr(subScope)},
						Permissions: []*armauthorization.Permission{{
							Actions:     actions,
							DataActions: dataActions,
						}},
					},
				}, nil)
			if err != nil {
				return fmt.Errorf("failed to create custom role %s: %w", customRoleName, err)
			}

			customRoleDefID := fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", subScope, customRoleID)
			if err := assignRole(ctx, roleAssignClient, rgScope, identity.PrincipalID, customRoleDefID); err != nil {
				return fmt.Errorf("failed to assign custom role to %s: %w", key, err)
			}
		}
	}
	return nil
}

// injectWIFManifests adds credential secrets and the Authentication CR
// to the bootstrap ignition so operators have WIF credentials at first boot.
func (p *Provider) injectWIFManifests(in clusterapi.IgnitionInput, bootstrapIgnData []byte) ([]byte, error) {
	session, err := in.InstallConfig.Azure.Session()
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure session: %w", err)
	}

	ignConfig, _, err := ignitionconfig.Parse(bootstrapIgnData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bootstrap ignition: %w", err)
	}

	subscriptionID := session.Credentials.SubscriptionID
	tenantID := session.Credentials.TenantID

	for key, identity := range p.wifIdentities {
		parts := strings.SplitN(key, "/", 2)
		namespace, name := parts[0], parts[1]

		tokenPath := aztypes.WIFDefaultTokenPath
		if tp, ok := p.wifTokenPaths[key]; ok && tp != "" {
			tokenPath = tp
		}

		secretManifest := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  name: %s
  namespace: %s
stringData:
  azure_subscription_id: %s
  azure_tenant_id: %s
  azure_client_id: %s
  azure_federated_token_file: %s
`, name, namespace, subscriptionID, tenantID, identity.ClientID, tokenPath)

		filename := fmt.Sprintf("99_cloud-creds-secret-%s-%s.yaml", namespace, name)
		addFileToIgnition(&ignConfig, "/opt/openshift/openshift/"+filename, []byte(secretManifest))
	}

	authManifest := fmt.Sprintf(`apiVersion: config.openshift.io/v1
kind: Authentication
metadata:
  name: cluster
spec:
  serviceAccountIssuer: %s
`, p.wifIssuerURL)
	addFileToIgnition(&ignConfig, "/opt/openshift/openshift/cluster-authentication-02-config.yaml", []byte(authManifest))

	logrus.Infof("Injected %d WIF credential secrets and Authentication CR into bootstrap ignition", len(p.wifIdentities))

	updatedData, err := json.Marshal(ignConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updated bootstrap ignition: %w", err)
	}
	return updatedData, nil
}

func addFileToIgnition(config *igntypes.Config, path string, data []byte) {
	encoded := fmt.Sprintf("data:text/plain;charset=utf-8;base64,%s", base64.StdEncoding.EncodeToString(data))
	config.Storage.Files = append(config.Storage.Files, igntypes.File{
		Node: igntypes.Node{Path: path},
		FileEmbedded1: igntypes.FileEmbedded1{
			Contents: igntypes.Resource{
				Source: ignutil.StrToPtr(encoded),
			},
			Mode: ignutil.IntToPtr(0o644),
		},
	})
}

func uploadBlob(ctx context.Context, blobBaseURL, containerName, blobName string, data []byte, cred azcore.TokenCredential) error {
	blobURL := fmt.Sprintf("%s/%s/%s", blobBaseURL, containerName, blobName)
	client, err := blockblob.NewClient(blobURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create block blob client for %s: %w", blobName, err)
	}
	_, err = client.UploadBuffer(ctx, data, nil)
	if err != nil {
		return fmt.Errorf("failed to upload %s: %w", blobName, err)
	}
	return nil
}

func findRoleDefinitionByName(ctx context.Context, client *armauthorization.RoleDefinitionsClient, scope, roleName string) (string, error) {
	pager := client.NewListPager(scope, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return "", fmt.Errorf("failed to list role definitions: %w", err)
		}
		for _, rd := range page.Value {
			if rd.Properties != nil && rd.Properties.RoleName != nil && *rd.Properties.RoleName == roleName {
				return *rd.ID, nil
			}
		}
	}
	return "", fmt.Errorf("role definition %q not found", roleName)
}

func assignRole(ctx context.Context, client *armauthorization.RoleAssignmentsClient, scope, principalID, roleDefinitionID string) error {
	roleAssignmentID := uuid.New().String()
	var err error
	for i := 0; i < retryCount; i++ {
		_, err = client.Create(ctx, scope, roleAssignmentID,
			armauthorization.RoleAssignmentCreateParameters{
				Properties: &armauthorization.RoleAssignmentProperties{
					PrincipalID:      to.Ptr(principalID),
					RoleDefinitionID: to.Ptr(roleDefinitionID),
				},
			}, nil)
		if err == nil {
			return nil
		}
		time.Sleep(retryTime)
	}
	return fmt.Errorf("failed to create role assignment after %d retries: %w", retryCount, err)
}

// sanitizeFICName ensures the name is valid for Azure FIC resources.
// Azure FIC names: 3-120 chars, alphanumeric, dashes, underscores.
func sanitizeFICName(name string) string {
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '-'
	}, name)
	if len(name) > 120 {
		name = name[:120]
	}
	if len(name) < 3 {
		name += "---"[:3-len(name)]
	}
	return name
}
