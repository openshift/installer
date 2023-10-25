// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

const (
	// AzureClientSecret is the client secret of the Azure Service Principal used to authenticate with Azure.
	// NOTE: This is required when using Service Principal authentication.
	// #nosec
	AzureClientSecret = "AZURE_CLIENT_SECRET"
	// AzureSubscriptionID is the Azure Subscription the operator will act against.
	AzureSubscriptionID = "AZURE_SUBSCRIPTION_ID"
	// AzureTenantID is the AAD tenant that the subscription is in
	AzureTenantID = "AZURE_TENANT_ID"
	// AzureClientID is the client ID of the Azure Service Principal or Managed Identity to use to authenticate with Azure.
	AzureClientID = "AZURE_CLIENT_ID"
	// AzureClientCertificate is a PEM or PKCS12 certificate string including the private key for Azure Credential Authentication.
	// If certificate is password protected,  use 'AzureClientCertificatePassword' for password.
	AzureClientCertificate = "AZURE_CLIENT_CERTIFICATE"
	// AzureClientCertificatePassword is password used to protect the AzureClientCertificate.
	// #nosec
	AzureClientCertificatePassword = "AZURE_CLIENT_CERTIFICATE_PASSWORD"
	// TargetNamespaces lists the namespaces the operator will watch
	// for Azure resources (if the mode includes running watchers). If
	// it's empty the operator will watch all namespaces.
	TargetNamespaces = "AZURE_TARGET_NAMESPACES"
	// OperatorMode determines whether the operator should run
	// watchers, webhooks or both.
	OperatorMode = "AZURE_OPERATOR_MODE"
	// SyncPeriod is the frequency at which resources are re-reconciled with Azure
	// when there have been no triggering changes in the Kubernetes resources. This sync
	// exists to detect and correct changes that happened in Azure that Kubernetes is not
	// aware about. BE VERY CAREFUL setting this value low - even a modest number of resources
	// can cause subscription level throttling if they are re-synced frequently.
	// If nil, no sync is performed. Durations are specified as "1h", "15m", or "60s". See
	// https://pkg.go.dev/time#ParseDuration for more details.
	SyncPeriod = "AZURE_SYNC_PERIOD"
	// ResourceManagerEndpoint is the Azure Resource Manager endpoint.
	// If not specified, the default is the Public cloud resource manager endpoint.
	// See https://docs.microsoft.com/cli/azure/manage-clouds-azure-cli#list-available-clouds for details
	// about how to find available resource manager endpoints for your cloud. Note that the resource manager
	// endpoint is referred to as "resourceManager" in the Azure CLI.
	ResourceManagerEndpoint = "AZURE_RESOURCE_MANAGER_ENDPOINT"
	// ResourceManagerAudience is the Azure Resource Manager AAD audience.
	// If not specified, the default is the Public cloud resource manager audience https://management.core.windows.net/.
	// See https://docs.microsoft.com/cli/azure/manage-clouds-azure-cli#list-available-clouds for details
	// about how to find available resource manager audiences for your cloud. Note that the resource manager
	// audience is referred to as "activeDirectoryResourceId" in the Azure CLI.
	ResourceManagerAudience = "AZURE_RESOURCE_MANAGER_AUDIENCE"
	// AzureAuthorityHost is the URL of the AAD authority. If not specified, the default
	// is the AAD URL for the public cloud: https://login.microsoftonline.com/. See
	// https://docs.microsoft.com/azure/active-directory/develop/authentication-national-cloud
	AzureAuthorityHost = "AZURE_AUTHORITY_HOST"
	// PodNamespace is the namespace the operator pods are running in.
	PodNamespace = "POD_NAMESPACE"
	// UseWorkloadIdentityAuth boolean is used to determine if we're using Workload Identity authentication for global credential
	UseWorkloadIdentityAuth = "USE_WORKLOAD_IDENTITY_AUTH"
)
