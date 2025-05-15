// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/pkg/common/config"
)

// These are hardcoded because the init function that initializes them in azcore isn't in /cloud it's in /arm which
// we don't import.

var (
	DefaultEndpoint                = "https://management.azure.com"
	DefaultAudience                = "https://management.core.windows.net/"
	DefaultAADAuthorityHost        = "https://login.microsoftonline.com/"
	DefaultMaxConcurrentReconciles = 1
)

// NOTE: Changes to documentation or available values here should be documented in Helm values.yaml as well

// Values stores configuration values that are set for the operator.
type Values struct {
	// SubscriptionID is the Azure subscription the operator will use
	// for ARM communication.
	SubscriptionID string

	// TenantID is the Azure tenantID the operator will use
	// for ARM communication.
	TenantID string

	// ClientID is the Azure clientID the operator will use
	// for ARM communication.
	ClientID string

	// PodNamespace is the namespace the operator pods are running in.
	PodNamespace string

	// OperatorMode determines whether the operator should run
	// watchers, webhooks or both.
	OperatorMode OperatorMode

	// TargetNamespaces lists the namespaces the operator will watch
	// for Azure resources (if the mode includes running watchers). If
	// it's empty the operator will watch all namespaces.
	TargetNamespaces []string

	// SyncPeriod is the frequency at which resources are re-reconciled with Azure
	// when there have been no triggering changes in the Kubernetes resources. This sync
	// exists to detect and correct changes that happened in Azure that Kubernetes is not
	// aware about. BE VERY CAREFUL setting this value low - even a modest number of resources
	// can cause subscription level throttling if they are re-synced frequently.
	// If nil, no sync is performed. Durations are specified as "1h", "15m", or "60s". See
	// https://pkg.go.dev/time#ParseDuration for more details.
	//
	// Specify the special value "never" for AZURE_SYNC_PERIOD to prevent syncing.
	SyncPeriod *time.Duration

	// ResourceManagerEndpoint is the Azure Resource Manager endpoint.
	// If not specified, the default is the Public cloud resource manager endpoint.
	// See https://docs.microsoft.com/cli/azure/manage-clouds-azure-cli#list-available-clouds for details
	// about how to find available resource manager endpoints for your cloud. Note that the resource manager
	// endpoint is referred to as "resourceManager" in the Azure CLI.
	ResourceManagerEndpoint string

	// ResourceManagerAudience is the Azure Resource Manager AAD audience.
	// If not specified, the default is the Public cloud resource manager audience https://management.core.windows.net/.
	// See https://docs.microsoft.com/cli/azure/manage-clouds-azure-cli#list-available-clouds for details
	// about how to find available resource manager audiences for your cloud. Note that the resource manager
	// audience is referred to as "activeDirectoryResourceId" in the Azure CLI.
	ResourceManagerAudience string

	// AzureAuthorityHost is the URL of the AAD authority. If not specified, the default
	// is the AAD URL for the public cloud: https://login.microsoftonline.com/. See
	// https://docs.microsoft.com/azure/active-directory/develop/authentication-national-cloud
	AzureAuthorityHost string

	// UseWorkloadIdentityAuth boolean is used to determine if we're using Workload Identity authentication for global credential
	UseWorkloadIdentityAuth bool

	// UserAgentSuffix is appended to the default User-Agent for Azure HTTP clients.
	UserAgentSuffix string

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
	MaxConcurrentReconciles int

	RateLimit RateLimit
}

type RateLimitMode string

const (
	RateLimitModeDisabled = RateLimitMode("disabled")
	RateLimitModeBucket   = RateLimitMode("bucket")
)

func ParseRateLimitMode(s string) (RateLimitMode, error) {
	switch s {
	case string(RateLimitModeDisabled):
		return RateLimitModeDisabled, nil
	case string(RateLimitModeBucket):
		return RateLimitModeBucket, nil
	default:
		return "", errors.Errorf("invalid rate limit mode %q", s)
	}
}

type RateLimit struct {
	// Mode configures the internal rate-limiting mode.
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
	Mode RateLimitMode

	// QPS is the rate (per second) that the bucket is refilled. This value only has an effect if Mode is 'bucket'.
	QPS float64

	// BucketSize is the size of the bucket. This value only has an effect if Mode is 'bucket'.
	BucketSize int
}

func (r RateLimit) String() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("Mode:%s", r.Mode))

	// Don't log anything other than disabled when mode is disabled
	if r.Mode != RateLimitModeDisabled {
		builder.WriteString(fmt.Sprintf("/QPS:%f/", r.QPS))
		builder.WriteString(fmt.Sprintf("BucketSize:%d", r.BucketSize))
	}

	return builder.String()
}

var _ fmt.Stringer = Values{}

// Returns the configuration as a string
func (v Values) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("SubscriptionID:%s/", v.SubscriptionID))
	builder.WriteString(fmt.Sprintf("TenantID:%s/", v.TenantID))
	builder.WriteString(fmt.Sprintf("ClientID:%s/", v.ClientID))
	builder.WriteString(fmt.Sprintf("PodNamespace:%s/", v.PodNamespace))
	builder.WriteString(fmt.Sprintf("OperatorMode:%s/", v.OperatorMode))
	builder.WriteString(fmt.Sprintf("TargetNamespaces:%s/", strings.Join(v.TargetNamespaces, "|")))
	builder.WriteString(fmt.Sprintf("SyncPeriod:%s/", v.SyncPeriod))
	builder.WriteString(fmt.Sprintf("ResourceManagerEndpoint:%s/", v.ResourceManagerEndpoint))
	builder.WriteString(fmt.Sprintf("ResourceManagerAudience:%s/", v.ResourceManagerAudience))
	builder.WriteString(fmt.Sprintf("AzureAuthorityHost:%s/", v.AzureAuthorityHost))
	builder.WriteString(fmt.Sprintf("UseWorkloadIdentityAuth:%t/", v.UseWorkloadIdentityAuth))
	builder.WriteString(fmt.Sprintf("UserAgentSuffix:%s/", v.UserAgentSuffix))
	builder.WriteString(fmt.Sprintf("MaxConcurrentReconciles:%d/", v.MaxConcurrentReconciles))
	builder.WriteString(fmt.Sprintf("RateLimit:[%s]", v.RateLimit.String()))

	return builder.String()
}

// Cloud returns the cloud the configuration is using
func (v Values) Cloud() cloud.Configuration {
	// Special handling if we've got all the defaults just return the official public cloud
	// configuration
	hasDefaultAzureAuthorityHost := v.AzureAuthorityHost == "" || v.AzureAuthorityHost == DefaultAADAuthorityHost
	hasDefaultResourceManagerEndpoint := v.ResourceManagerEndpoint == "" || v.ResourceManagerEndpoint == DefaultEndpoint
	hasDefaultResourceManagerAudience := v.ResourceManagerAudience == "" || v.ResourceManagerAudience == DefaultAudience

	if hasDefaultResourceManagerEndpoint && hasDefaultResourceManagerAudience && hasDefaultAzureAuthorityHost {
		return cloud.AzurePublic
	}

	// We default here too to more easily support empty Values objects
	azureAuthorityHost := v.AzureAuthorityHost
	resourceManagerEndpoint := v.ResourceManagerEndpoint
	resourceManagerAudience := v.ResourceManagerAudience
	if azureAuthorityHost == "" {
		azureAuthorityHost = DefaultAADAuthorityHost
	}
	if resourceManagerAudience == "" {
		resourceManagerAudience = DefaultAudience
	}
	if resourceManagerEndpoint == "" {
		resourceManagerEndpoint = DefaultEndpoint
	}

	return cloud.Configuration{
		ActiveDirectoryAuthorityHost: azureAuthorityHost,
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Endpoint: resourceManagerEndpoint,
				Audience: resourceManagerAudience,
			},
		},
	}
}

// ReadFromEnvironment loads configuration values from the AZURE_*
// environment variables.
func ReadFromEnvironment() (Values, error) {
	var result Values
	modeValue := os.Getenv(config.OperatorMode)
	if modeValue == "" {
		result.OperatorMode = OperatorModeBoth
	} else {
		mode, err := ParseOperatorMode(modeValue)
		if err != nil {
			return Values{}, err
		}
		result.OperatorMode = mode
	}

	var err error

	result.SubscriptionID = os.Getenv(config.AzureSubscriptionID)
	result.PodNamespace = os.Getenv(config.PodNamespace)
	result.TargetNamespaces = parseTargetNamespaces(os.Getenv(config.TargetNamespaces))
	result.SyncPeriod, err = parseSyncPeriod()
	if err != nil {
		return result, errors.Wrapf(err, "parsing %q", config.SyncPeriod)
	}

	result.ResourceManagerEndpoint = envOrDefault(config.ResourceManagerEndpoint, DefaultEndpoint)
	result.ResourceManagerAudience = envOrDefault(config.ResourceManagerAudience, DefaultAudience)
	result.AzureAuthorityHost = envOrDefault(config.AzureAuthorityHost, DefaultAADAuthorityHost)
	result.ClientID = os.Getenv(config.AzureClientID)
	result.TenantID = os.Getenv(config.AzureTenantID)
	result.MaxConcurrentReconciles, err = envParseOrDefault(config.MaxConcurrentReconciles, DefaultMaxConcurrentReconciles)
	if err != nil {
		return result, err
	}

	// Ignoring error here, as any other value or empty value means we should default to false
	result.UseWorkloadIdentityAuth, _ = strconv.ParseBool(os.Getenv(config.UseWorkloadIdentityAuth))
	result.UserAgentSuffix = os.Getenv(config.UserAgentSuffix)
	result.RateLimit.Mode, err = ParseRateLimitMode(envOrDefault(config.RateLimitMode, string(RateLimitModeDisabled)))
	if err != nil {
		return result, err
	}
	result.RateLimit.QPS, err = envParseOrDefault(config.RateLimitQPS, 5.0)
	if err != nil {
		return result, err
	}
	result.RateLimit.BucketSize, err = envParseOrDefault(config.RateLimitBucketSize, 100)
	if err != nil {
		return result, err
	}

	// Not calling validate here to support using from tests where we
	// don't require consistent settings.
	return result, nil
}

// ReadAndValidate loads the configuration values and checks that
// they're consistent.
func ReadAndValidate() (Values, error) {
	result, err := ReadFromEnvironment()
	if err != nil {
		return Values{}, err
	}
	err = result.Validate()
	if err != nil {
		return Values{}, err
	}
	return result, nil
}

// Validate checks whether the configuration settings are consistent.
func (v Values) Validate() error {
	if v.PodNamespace == "" {
		return errors.Errorf("missing value for %s", config.PodNamespace)
	}
	if !v.OperatorMode.IncludesWatchers() && len(v.TargetNamespaces) > 0 {
		return errors.Errorf("%s must include watchers to specify target namespaces", config.TargetNamespaces)
	}
	if v.MaxConcurrentReconciles <= 0 {
		return errors.Errorf("%s must be at least 1", config.MaxConcurrentReconciles)
	}
	return nil
}

// parseTargetNamespaces splits a comma-separated string into a slice
// of strings with spaces trimmed.
func parseTargetNamespaces(fromEnv string) []string {
	if len(strings.TrimSpace(fromEnv)) == 0 {
		return nil
	}
	items := strings.Split(fromEnv, ",")
	// Remove any whitespace used to separate items.
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}

// parseSyncPeriod parses the sync period from the environment
func parseSyncPeriod() (*time.Duration, error) {
	syncPeriodStr := envOrDefault(config.SyncPeriod, "1h")
	if syncPeriodStr == "never" { // magical string that means no sync
		return nil, nil
	}

	syncPeriod, err := time.ParseDuration(syncPeriodStr)
	if err != nil {
		return nil, err
	}
	return &syncPeriod, nil
}

func envParseOrDefault[T int | string | float64](env string, def T) (T, error) {
	str, specified := os.LookupEnv(env)
	if !specified {
		return def, nil
	}
	if str == "" {
		return def, nil
	}

	var result T
	switch any(def).(type) {
	case int:
		parsedVal, err := strconv.Atoi(str)
		if err != nil {
			return def, errors.Wrapf(err, "failed to parse value %q for %q", str, env)
		}
		result = any(parsedVal).(T)
	case string:
		result = any(str).(T)
	case float64:
		parsedVal, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return def, errors.Wrapf(err, "failed to parse value %q for %q", str, env)
		}
		result = any(parsedVal).(T)
	default:
		return def, errors.Errorf("can't read unsupported type %T from env", def)
	}

	return result, nil
}

// envOrDefault returns the value of the specified env variable or the default value if
// the env variable was not set.
func envOrDefault(env string, def string) string {
	result, specified := os.LookupEnv(env)
	if !specified {
		return def
	}
	if result == "" {
		return def
	}

	return result
}
