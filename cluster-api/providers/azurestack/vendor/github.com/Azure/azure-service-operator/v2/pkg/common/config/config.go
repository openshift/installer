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
	// If the certificate is password protected,  use the 'AzureClientCertificatePassword' for password.
	AzureClientCertificate = "AZURE_CLIENT_CERTIFICATE"
	// AzureClientCertificatePassword is the password used to protect the AzureClientCertificate.
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
	// Durations are specified as "1h", "15m", or "60s". Specify the special value "never" to prevent
	// syncing. See https://pkg.go.dev/time#ParseDuration for more details.
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
	// UserAgentSuffix is appended to the default User-Agent for Azure HTTP clients.
	UserAgentSuffix = "AZURE_USER_AGENT_SUFFIX"
	// MaxConcurrentReconciles is the number of threads/goroutines dedicated to reconciling each resource type.
	// If not specified, the default is 1.
	// IMPORTANT: Having MaxConcurrentReconciles set to N does not mean that ASO is limited to N interactions with
	// Azure at any given time, because the control loop yields to another resource while it is not actively issuing HTTP
	// calls to Azure. Any single resource only blocks the control-loop for its resource-type for as long as it takes to issue
	// an HTTP call to Azure, view the result, and make a decision. In most cases the time taken to perform these actions
	// (and thus how long the loop is blocked and preventing other resources from being acted upon) is a few hundred
	// milliseconds to at most a second or two. In a typical 60s period, many hundreds or even thousands of resources
	// can be managed with this set to 1.
	// MaxConcurrentReconciles applies to every registered resource type being watched/managed by ASO.
	MaxConcurrentReconciles = "MAX_CONCURRENT_RECONCILES"
	// RateLimitMode configures the internal rate-limiting mode.
	// Valid values are [disabled, bucket]
	// * disabled: No ASO-controlled rate-limiting occurs. ASO will attempt to communicate with Azure and
	//   kube-apiserver as much as needed based on load. It will back off based on throttling from
	//   either kube-apiserver or Azure, but will not artificially limit its throughput.
	// * bucket: Uses a token-bucket algorithm to rate-limit reconciliations. Note that this limits how often
	//   the operator performs a reconciliation, but not every reconciliation triggers a call to kube-apiserver
	//   or Azure (though many do). Since this controls reconciles it can be used to coarsely control throughput
	//   and CPU usage of the operator, as well as the number of requests that the operator issues to Azure.
	//   Keep in mind that the Azure throttling limits (defined at
	//   https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/request-limits-and-throttling)
	//   differentiate between request types. Since a given reconcile for a resource may result in polling (a GET) or
	//   modification (a PUT) it's not possible to entirely avoid Azure throttling by tuning these bucket limits.
	//   We don't recommend enabling this mode by default.
	//   If enabling this mode, we strongly recommend doing some experimentation to tune these values to something to
	//   works for your specific need.
	RateLimitMode = "RATE_LIMIT_MODE"
	// RateLimitQPS is the rate (per second) that the bucket is refilled. This value only has an effect if RateLimitMode is 'bucket'.
	RateLimitQPS = "RATE_LIMIT_QPS"
	// RateLimitBucketSize is the size of the bucket. This value only has an effect if RateLimitMode is 'bucket'.
	RateLimitBucketSize = "RATE_LIMIT_BUCKET_SIZE"
)
