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
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	elb "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/identity"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/throttle"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/logger"
	"sigs.k8s.io/cluster-api-provider-aws/v2/util/system"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

const (
	notPermittedError = "Namespace is not permitted to use %s: %s"
)

var sessionCache sync.Map
var providerCache sync.Map

type sessionCacheEntry struct {
	session         *aws.Config
	serviceLimiters throttle.ServiceLimiters
}

// ChainCredentialsProvider defines custom CredentialsProvider chain
// NewChainCredentialsProvider can be used to initialize this struct.
type ChainCredentialsProvider struct {
	providers []aws.CredentialsProvider
}

func sessionForRegion(region string) (*aws.Config, throttle.ServiceLimiters, error) {
	if s, ok := sessionCache.Load(region); ok {
		entry := s.(*sessionCacheEntry)
		return entry.session, entry.serviceLimiters, nil
	}

	ns, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))

	if err != nil {
		return nil, nil, err
	}

	sl := newServiceLimiters()
	sessionCache.Store(region, &sessionCacheEntry{
		session:         &ns,
		serviceLimiters: sl,
	})
	return &ns, sl, nil
}

func sessionForClusterWithRegion(k8sClient client.Client, clusterScoper cloud.SessionMetadata, region string, log logger.Wrapper) (*aws.Config, throttle.ServiceLimiters, error) {
	log = log.WithName("identity")
	log.Trace("Creating an AWS Session")

	providers, err := getProvidersForCluster(context.Background(), k8sClient, clusterScoper, region, log)
	if err != nil {
		// could not get providers and retrieve the credentials
		conditions.MarkFalse(clusterScoper.InfraCluster(), infrav1.PrincipalCredentialRetrievedCondition, infrav1.PrincipalCredentialRetrievalFailedReason, clusterv1.ConditionSeverityError, "%s", err.Error())
		return nil, nil, errors.Wrap(err, "Failed to get providers for cluster")
	}

	isChanged := false
	awsProviders := make([]aws.CredentialsProvider, len(providers))
	for i, provider := range providers {
		// load an existing matching providers from the cache if such a providers exists
		providerHash, err := provider.Hash()
		if err != nil {
			return nil, nil, errors.Wrap(err, "Failed to calculate provider hash")
		}
		cachedProvider, ok := providerCache.Load(providerHash)
		if ok {
			provider = cachedProvider.(identity.AWSPrincipalTypeProvider)
		} else {
			isChanged = true
			// add this provider to the cache
			providerCache.Store(providerHash, provider)
		}
		awsProviders[i] = provider.(aws.CredentialsProvider)
	}

	if !isChanged {
		if s, ok := sessionCache.Load(getSessionName(region, clusterScoper)); ok {
			entry := s.(*sessionCacheEntry)
			return entry.session, entry.serviceLimiters, nil
		}
	}

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}

	if len(providers) > 0 {
		// Check if identity credentials can be retrieved. One reason this will fail is that source identity is not authorized for assume role.
		_, err := providers[0].Retrieve(context.Background())
		if err != nil {
			conditions.MarkUnknown(clusterScoper.InfraCluster(), infrav1.PrincipalCredentialRetrievedCondition, infrav1.CredentialProviderBuildFailedReason, "%s", err.Error())

			// delete the existing session from cache. Otherwise, we give back a defective session on next method invocation with same cluster scope
			sessionCache.Delete(getSessionName(region, clusterScoper))

			return nil, nil, errors.Wrap(err, "Failed to retrieve identity credentials")
		}
		chainProvider := NewChainCredentialsProvider(awsProviders)
		optFns = append(optFns, config.WithCredentialsProvider(chainProvider))
	}

	conditions.MarkTrue(clusterScoper.InfraCluster(), infrav1.PrincipalCredentialRetrievedCondition)

	ns, err := config.LoadDefaultConfig(context.Background(), optFns...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create a new AWS session")
	}
	sl := newServiceLimiters()
	sessionCache.Store(getSessionName(region, clusterScoper), &sessionCacheEntry{
		session:         &ns,
		serviceLimiters: sl,
	})

	return &ns, sl, nil
}

func getSessionName(region string, clusterScoper cloud.SessionMetadata) string {
	return fmt.Sprintf("%s-%s-%s-%s", region, clusterScoper.ControllerName(), clusterScoper.InfraClusterName(), clusterScoper.Namespace())
}

func newServiceLimiters() throttle.ServiceLimiters {
	return throttle.ServiceLimiters{
		ec2.ServiceID:                      newEC2ServiceLimiter(),
		elb.ServiceID:                      newGenericServiceLimiter(),
		elbv2.ServiceID:                    newGenericServiceLimiter(),
		resourcegroupstaggingapi.ServiceID: newGenericServiceLimiter(),
		secretsmanager.ServiceID:           newGenericServiceLimiter(),
	}
}

func newGenericServiceLimiter() *throttle.ServiceLimiter {
	return &throttle.ServiceLimiter{
		{
			Operation:  throttle.NewMultiOperationMatch("Describe", "Get", "List"),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation:  ".*",
			RefillRate: 5.0,
			Burst:      200,
		},
	}
}

func newEC2ServiceLimiter() *throttle.ServiceLimiter {
	return &throttle.ServiceLimiter{
		{
			Operation:  throttle.NewMultiOperationMatch("Describe", "Get"),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation: throttle.NewMultiOperationMatch(
				"AuthorizeSecurityGroupIngress",
				"CancelSpotInstanceRequests",
				"CreateKeyPair",
				"RequestSpotInstances",
			),
			RefillRate: 20.0,
			Burst:      100,
		},
		{
			Operation:  "RunInstances",
			RefillRate: 2.0,
			Burst:      5,
		},
		{
			Operation:  "StartInstances",
			RefillRate: 2.0,
			Burst:      5,
		},
		{
			Operation:  ".*",
			RefillRate: 5.0,
			Burst:      200,
		},
	}
}

func buildProvidersForRef(
	ctx context.Context,
	providers []identity.AWSPrincipalTypeProvider,
	k8sClient client.Client,
	clusterScoper cloud.SessionMetadata,
	ref *infrav1.AWSIdentityReference,
	region string,
	log logger.Wrapper) ([]identity.AWSPrincipalTypeProvider, error) {
	if ref == nil {
		log.Trace("AWSCluster does not have a IdentityRef specified")
		return providers, nil
	}

	var provider identity.AWSPrincipalTypeProvider
	identityObjectKey := client.ObjectKey{Name: ref.Name}
	log = log.WithValues("identityKey", identityObjectKey)
	log.Trace("Getting identity")

	switch ref.Kind {
	case infrav1.ControllerIdentityKind:
		err := buildAWSClusterControllerIdentity(ctx, identityObjectKey, k8sClient, clusterScoper)
		if err != nil {
			return providers, err
		}
		// returning empty provider list to default to Controller Principal.
		return []identity.AWSPrincipalTypeProvider{}, nil
	case infrav1.ClusterStaticIdentityKind:
		provider, err := buildAWSClusterStaticIdentity(ctx, identityObjectKey, k8sClient, clusterScoper)
		if err != nil {
			return providers, err
		}
		providers = append(providers, provider)
	case infrav1.ClusterRoleIdentityKind:
		roleIdentity := &infrav1.AWSClusterRoleIdentity{}
		err := k8sClient.Get(ctx, identityObjectKey, roleIdentity)
		if err != nil {
			return providers, err
		}
		log.Trace("Principal retrieved")
		canUse, err := isClusterPermittedToUsePrincipal(k8sClient, roleIdentity.Spec.AllowedNamespaces, clusterScoper.Namespace())
		if err != nil {
			return providers, err
		}
		if !canUse {
			setPrincipalUsageNotAllowedCondition(infrav1.ClusterRoleIdentityKind, identityObjectKey, clusterScoper)
			return providers, errors.Errorf(notPermittedError, infrav1.ClusterRoleIdentityKind, roleIdentity.Name)
		}
		setPrincipalUsageAllowedCondition(clusterScoper)

		if roleIdentity.Spec.SourceIdentityRef != nil {
			providers, err = buildProvidersForRef(ctx, providers, k8sClient, clusterScoper, roleIdentity.Spec.SourceIdentityRef, region, log)
			if err != nil {
				return providers, err
			}
		}
		var sourceProvider identity.AWSPrincipalTypeProvider
		if len(providers) > 0 {
			sourceProvider = providers[len(providers)-1]
			// Remove last provider
			if len(providers) > 0 {
				providers = providers[:len(providers)-1]
			}
		}

		provider = identity.NewAWSRolePrincipalTypeProvider(roleIdentity, sourceProvider, region, log)
		providers = append(providers, provider)
	default:
		return providers, errors.Errorf("No such provider known: '%s'", ref.Kind)
	}
	conditions.MarkTrue(clusterScoper.InfraCluster(), infrav1.PrincipalUsageAllowedCondition)
	return providers, nil
}

func setPrincipalUsageAllowedCondition(clusterScoper cloud.SessionMetadata) {
	conditions.MarkTrue(clusterScoper.InfraCluster(), infrav1.PrincipalUsageAllowedCondition)
}

func setPrincipalUsageNotAllowedCondition(kind infrav1.AWSIdentityKind, identityObjectKey client.ObjectKey, clusterScoper cloud.SessionMetadata) {
	errMsg := fmt.Sprintf(notPermittedError, kind, identityObjectKey.Name)

	if clusterScoper.IdentityRef().Name == identityObjectKey.Name {
		conditions.MarkFalse(clusterScoper.InfraCluster(), infrav1.PrincipalUsageAllowedCondition, infrav1.PrincipalUsageUnauthorizedReason, clusterv1.ConditionSeverityError, "%s", errMsg)
	} else {
		conditions.MarkFalse(clusterScoper.InfraCluster(), infrav1.PrincipalUsageAllowedCondition, infrav1.SourcePrincipalUsageUnauthorizedReason, clusterv1.ConditionSeverityError, "%s", errMsg)
	}
}

func buildAWSClusterStaticIdentity(ctx context.Context, identityObjectKey client.ObjectKey, k8sClient client.Client, clusterScoper cloud.SessionMetadata) (*identity.AWSStaticPrincipalTypeProvider, error) {
	staticPrincipal := &infrav1.AWSClusterStaticIdentity{}
	err := k8sClient.Get(ctx, identityObjectKey, staticPrincipal)
	if err != nil {
		return nil, err
	}
	secret := &corev1.Secret{}
	err = k8sClient.Get(ctx, client.ObjectKey{Name: staticPrincipal.Spec.SecretRef, Namespace: system.GetManagerNamespace()}, secret)
	if err != nil {
		return nil, err
	}

	// Set ClusterStaticPrincipal as Secret's owner reference for 'clusterctl move'.
	patchHelper, err := patch.NewHelper(secret, k8sClient)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to init patch helper for secret name:%s namespace:%s", secret.Name, secret.Namespace)
	}

	secret.OwnerReferences = util.EnsureOwnerRef(secret.OwnerReferences, metav1.OwnerReference{
		APIVersion: infrav1.GroupVersion.String(),
		Kind:       string(infrav1.ClusterStaticIdentityKind),
		Name:       staticPrincipal.Name,
		UID:        staticPrincipal.UID,
	})

	if err := patchHelper.Patch(ctx, secret); err != nil {
		return nil, errors.Wrapf(err, "failed to patch secret name:%s namespace:%s", secret.Name, secret.Namespace)
	}

	canUse, err := isClusterPermittedToUsePrincipal(k8sClient, staticPrincipal.Spec.AllowedNamespaces, clusterScoper.Namespace())
	if err != nil {
		return nil, err
	}
	if !canUse {
		setPrincipalUsageNotAllowedCondition(infrav1.ClusterStaticIdentityKind, identityObjectKey, clusterScoper)
		return nil, errors.Errorf(notPermittedError, infrav1.ClusterStaticIdentityKind, identityObjectKey.Name)
	}
	setPrincipalUsageAllowedCondition(clusterScoper)

	return identity.NewAWSStaticPrincipalTypeProvider(staticPrincipal, secret), nil
}

func buildAWSClusterControllerIdentity(ctx context.Context, identityObjectKey client.ObjectKey, k8sClient client.Client, clusterScoper cloud.SessionMetadata) error {
	controllerIdentity := &infrav1.AWSClusterControllerIdentity{}
	controllerIdentity.Kind = string(infrav1.ControllerIdentityKind)

	// Enforce the singleton again for depth
	if identityObjectKey.Name != infrav1.AWSClusterControllerIdentityName {
		return errors.Errorf("Expected AWSClusterControllerIdentity of name %s, got %s", infrav1.AWSClusterControllerIdentityName, identityObjectKey.Name)
	}

	err := k8sClient.Get(ctx, client.ObjectKey{Name: identityObjectKey.Name}, controllerIdentity)
	if err != nil {
		return err
	}

	canUse, err := isClusterPermittedToUsePrincipal(k8sClient, controllerIdentity.Spec.AllowedNamespaces, clusterScoper.Namespace())
	if err != nil {
		return err
	}
	if !canUse {
		setPrincipalUsageNotAllowedCondition(infrav1.ControllerIdentityKind, identityObjectKey, clusterScoper)
		return errors.Errorf(notPermittedError, infrav1.ControllerIdentityKind, controllerIdentity.Name)
	}
	setPrincipalUsageAllowedCondition(clusterScoper)
	return nil
}

func getProvidersForCluster(ctx context.Context, k8sClient client.Client, clusterScoper cloud.SessionMetadata, region string, log logger.Wrapper) ([]identity.AWSPrincipalTypeProvider, error) {
	providers := make([]identity.AWSPrincipalTypeProvider, 0)
	providers, err := buildProvidersForRef(ctx, providers, k8sClient, clusterScoper, clusterScoper.IdentityRef(), region, log)
	if err != nil {
		return nil, err
	}

	return providers, nil
}

func isClusterPermittedToUsePrincipal(k8sClient client.Client, allowedNs *infrav1.AllowedNamespaces, clusterNamespace string) (bool, error) {
	// nil value does not match with any namespaces
	if allowedNs == nil {
		return false, nil
	}

	// empty value matches with all namespaces
	if cmp.Equal(*allowedNs, infrav1.AllowedNamespaces{}) {
		return true, nil
	}

	for _, v := range allowedNs.NamespaceList {
		if v == clusterNamespace {
			return true, nil
		}
	}

	// Check if clusterNamespace is in the namespaces selected by the identity's allowedNamespaces selector.
	namespaces := &corev1.NamespaceList{}
	selector, err := metav1.LabelSelectorAsSelector(&allowedNs.Selector)
	if err != nil {
		return false, errors.Wrap(err, "failed to get label selector from spec selector")
	}

	// If a Selector has a nil or empty selector, it should match nothing, not everything.
	if selector.Empty() {
		return false, nil
	}

	if err := k8sClient.List(context.Background(), namespaces, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return false, errors.Wrap(err, "failed to list namespaces")
	}

	for i := range namespaces.Items {
		n := &namespaces.Items[i]
		if n.Name == clusterNamespace {
			return true, nil
		}
	}
	return false, nil
}

// NewChainCredentialsProvider initializes a new ChainCredentialsProvider.
func NewChainCredentialsProvider(providers []aws.CredentialsProvider) *ChainCredentialsProvider {
	return &ChainCredentialsProvider{
		providers: providers,
	}
}

// Retrieve implements aws.CredentialsProvider for custom list of credenetials providers.
// The first provider in the list without error will be used to return credentials.
func (c *ChainCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	var lastErr error
	for _, provider := range c.providers {
		creds, err := provider.Retrieve(ctx)
		if err != nil {
			lastErr = err
			continue
		}
		if creds.AccessKeyID != "" && creds.SecretAccessKey != "" {
			return creds, nil
		}
	}
	return aws.Credentials{}, lastErr
}
