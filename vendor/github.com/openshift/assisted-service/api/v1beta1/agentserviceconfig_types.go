/*
Copyright 2021.

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

package v1beta1

import (
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OSImage defines an Operating System image and the OpenShift version it
// is associated with.
type OSImage struct {
	// OpenshiftVersion is the Major.Minor version of OpenShift that this image
	// is to be associated with.
	OpenshiftVersion string `json:"openshiftVersion"`
	// Version is the Operating System version of the image.
	Version string `json:"version"`
	// Url specifies the path to the Operating System image.
	Url string `json:"url"`
	// rootFSUrl specifies the path to the root filesystem.
	// +optional
	// Deprecated: this field is ignored (will be removed in a future release).
	RootFSUrl string `json:"rootFSUrl"`
	// The CPU architecture of the image (x86_64/arm64/etc).
	// +optional
	CPUArchitecture string `json:"cpuArchitecture"`
}

type MustGatherImage struct {
	// OpenshiftVersion is the Major.Minor version of OpenShift that this image
	// is to be associated with.
	OpenshiftVersion string `json:"openshiftVersion"`
	// Name specifies the name of the component (e.g. operator)
	// that the image is used to collect information about.
	Name string `json:"name"`
	// Url specifies the path to the Operating System image.
	Url string `json:"url"`
}

// AgentServiceConfigSpec defines the desired state of AgentServiceConfig.
type AgentServiceConfigSpec struct {
	// FileSystemStorage defines the spec of the PersistentVolumeClaim to be
	// created for the assisted-service's filesystem (logs, etc).
	// With respect to the resource requests, the amount of filesystem storage
	// consumed will depend largely on the number of clusters created (~200MB
	// per cluster and ~2-3GiB per supported OpenShift version). Minimum 100GiB
	// recommended.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage for service filesystem"
	FileSystemStorage corev1.PersistentVolumeClaimSpec `json:"filesystemStorage"`
	// DatabaseStorage defines the spec of the PersistentVolumeClaim to be
	// created for the database's filesystem.
	// With respect to the resource requests, minimum 10GiB is recommended.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage for database"
	DatabaseStorage corev1.PersistentVolumeClaimSpec `json:"databaseStorage"`
	// ImageStorage defines the spec of the PersistentVolumeClaim to be
	// created for each replica of the image service.
	// If a PersistentVolumeClaim is provided 2GiB per OSImage entry is required
	// +optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage for images"
	ImageStorage *corev1.PersistentVolumeClaimSpec `json:"imageStorage"`
	// MirrorRegistryRef is the reference to the configmap that contains mirror registry configuration
	// In case no configuration is need, this field will be nil. ConfigMap must contain to entries:
	// ca-bundle.crt - hold the contents of mirror registry certificate/s
	// registries.conf - holds the content of registries.conf file configured with mirror registries
	// +optional
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Mirror Registry and Certificate ConfigMap Name"
	MirrorRegistryRef *corev1.LocalObjectReference `json:"mirrorRegistryRef,omitempty"`

	// OSImages defines a collection of Operating System images (ie. RHCOS images)
	// that the assisted-service should use as the base when generating discovery ISOs.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Operating System Images"
	OSImages []OSImage `json:"osImages,omitempty"`

	// MustGatherImages defines a collection of operator related must-gather images
	// that are used if one the operators fails to be successfully deployed
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Must-Gather Images"
	MustGatherImages []MustGatherImage `json:"mustGatherImages,omitempty"`

	// IPXEHTTPRoute is controlling whether the operator is creating plain HTTP routes
	// iPXE hosts may not work with router cyphers and may access artifacts via HTTP only
	// This setting accepts "enabled,disabled", defaults to disabled. Empty value defaults to disabled
	// The following endpoints would be exposed via http:
	// * api/assisted-installer/v2/infra-envs/<id>/downloads/files?file_name=ipxe-script in assisted-service
	// * boot-artifacts/ and images/<infra-enf id>/pxe-initrd in -image-service
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Expose IPXE HTTP route"
	// +kubebuilder:validation:Enum=enabled;disabled
	// +kubebuilder:validation:default:=disabled
	// +optional
	IPXEHTTPRoute string `json:"iPXEHTTPRoute,omitempty"`
	// UnauthenticatedRegistries is a list of registries from which container images can be pulled
	// without authentication. They will be appended to the default list (quay.io,
	// registry.ci.openshift.org). Any registry on this list will not require credentials
	// to be in the pull secret validated by the assisted-service.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="List of container registries without authentication"
	// +optional
	UnauthenticatedRegistries []string `json:"unauthenticatedRegistries,omitempty"`

	// OSImageCACertRef is a reference to a config map containing a certificate authority certificate
	// this is an optional certificate to allow a user to add a certificate authority for a HTTPS source of images
	// this certificate will be used by the assisted-image-service when pulling OS images.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="OS Image CA Cert ConfigMap reference"
	// +optional
	OSImageCACertRef *corev1.LocalObjectReference `json:"OSImageCACertRef,omitempty"`

	// OSImageAdditionalParamsRef is a reference to a secret containing a headers and query parameters to be used during OS image fetch.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="OS Images additional parameters reference"
	// +optional
	OSImageAdditionalParamsRef *corev1.LocalObjectReference `json:"OSImageAdditionalParamsRef,omitempty"`

	// Ingress contains configuration for the ingress resources.
	// Has no effect when running on an OpenShift cluster.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ingress"
	// +optional
	Ingress *Ingress `json:"ingress,omitempty"`
}

type Ingress struct {
	// AssistedServiceHostname is the hostname to be assigned to the assisted-service ingress.
	// Has no effect when running on an OpenShift cluster.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Assisted Service hostname"
	AssistedServiceHostname string `json:"assistedServiceHostname"`

	// ImageServiceHostname is the hostname to be assigned to the assisted-image-service ingress.
	// Has no effect when running on an OpenShift cluster.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Assisted Image Service hostname"
	ImageServiceHostname string `json:"imageServiceHostname"`

	// ClassName is the name of the ingress class to be used when configuring ingress resources.
	// Has no effect when running on an OpenShift cluster.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Class Name"
	// +optional
	ClassName *string `json:"className,omitempty"`
}

// ConditionType related to our reconcile loop in addition to all the reasons
// why ConditionStatus could be true or false.
const (
	// ConditionReconcileCompleted reports on whether or not the local cluster is managed.
	ConditionLocalClusterManaged conditionsv1.ConditionType = "LocalClusterManaged"
	// ConditionReconcileCompleted reports whether reconcile completed without error.
	ConditionReconcileCompleted conditionsv1.ConditionType = "ReconcileCompleted"
	// ConditionDeploymentsHealthy reports whether deployments are healthy.
	ConditionDeploymentsHealthy conditionsv1.ConditionType = "DeploymentsHealthy"
	// ReasonLocalClusterImportNotEnabled when the import of local cluster is not enabled.
	ReasonLocalClusterImportNotEnabled string = "Local cluster import is not enabled"
	// ReasonLocalClusterEntitiesCreated when the local cluster is managed.
	ReasonLocalClusterManaged string = "Local cluster is managed."
	// ReasonLocalClusterEntitiesRemoved when the local cluster is not managed.
	ReasonLocalClusterNotManaged string = "Local cluster is not managed."
	// ReasonUnableToDetermineLocalClusterManagedStatus when unable to determine the status of local cluster entities.
	ReasonUnableToDetermineLocalClusterManagedStatus string = "Unable to determine local cluster managed status."
	// ReasonReconcileSucceeded when the reconcile completes all operations without error.
	ReasonReconcileSucceeded string = "ReconcileSucceeded"
	// ReasonDeploymentSucceeded when configuring/deploying the assisted-service deployment completed without errors.
	ReasonDeploymentSucceeded string = "DeploymentSucceeded"
	// ReasonStorageFailure when there was a failure configuring/deploying storage.
	ReasonStorageFailure string = "StorageFailure"
	// ReasonImageHandlerServiceFailure when there was a failure related to the assisted-image-service's service.
	ReasonImageHandlerServiceFailure string = "ImageHandlerServiceFailure"
	// ReasonAgentServiceFailure when there was a failure related to the assisted-service's service.
	ReasonAgentServiceFailure string = "AgentServiceFailure"
	// ReasonAgentServiceFailure when there was a failure related to generating/deploying the service monitor.
	ReasonAgentServiceMonitorFailure string = "AgentServiceMonitorFailure"
	// ReasonImageHandlerRouteFailure when there was a failure configuring/deploying the assisted-image-service's route.
	ReasonImageHandlerRouteFailure string = "ImageHandlerRouteFailure"
	// ReasonAgentRouteFailure when there was a failure configuring/deploying the assisted-service's route.
	ReasonAgentRouteFailure string = "AgentRouteFailure"
	// ReasonAgentLocalAuthSecretFailure when there was a failure generating/deploying the local auth key pair secret.
	ReasonAgentLocalAuthSecretFailure string = "AgentLocalAuthSecretFailure" // #nosec
	// ReasonPostgresSecretFailure when there was a failure generating/deploying the database secret.
	ReasonPostgresSecretFailure string = "PostgresSecretFailure"
	// ReasonImageHandlerServiceAccountFailure when there was a failure related to the assisted-image-service's service account.
	ReasonImageHandlerServiceAccountFailure string = "ImageHandlerServiceAccountFailure"
	// ReasonIngressCertFailure when there was a failure generating/deploying the ingress cert configmap.
	ReasonIngressCertFailure string = "IngressCertFailure"
	// ReasonConfigFailure when there was a failure configuring/deploying the assisted-service configmap.
	ReasonConfigFailure string = "ConfigFailure"
	// ReasonImageHandlerStatefulSetFailure when there was a failure configuring/deploying the assisted-image-service stateful set.
	ReasonImageHandlerStatefulSetFailure string = "ImageHandlerStatefulSetFailure"
	// ReasonDeploymentFailure when there was a failure configuring/deploying the assisted-service deployment.
	ReasonDeploymentFailure string = "DeploymentFailure"
	// ReasonStorageFailure when there was a failure configuring/deploying the validating webhook.
	ReasonValidatingWebHookFailure string = "ValidatingWebHookFailure"
	// ReasonStorageFailure when there was a failure configuring/deploying the validating webhook.
	ReasonMutatingWebHookFailure string = "MutatingWebHookFailure"
	// ReasonWebHookServiceFailure when there was a failure related to the webhook's service.
	ReasonWebHookServiceFailure string = "ReasonWebHookServiceFailure"
	// ReasonWebHookDeploymentFailure when there was a failure configuring/deploying the webhook deployment.
	ReasonWebHookDeploymentFailure string = "ReasonWebHookDeploymentFailure"
	// ReasonWebReasonWebHookClusterRoleBindingFailureHookDeploymentFailure when there was a failure configuring/deploying the webhook cluster role binding.
	ReasonWebHookClusterRoleBindingFailure string = "ReasonWebHookClusterRoleBindingFailure"
	// ReasonWebHookClusterRoleFailure when there was a failure configuring/deploying the webhook cluster role.
	ReasonWebHookClusterRoleFailure string = "ReasonWebHookClusterRoleFailure"
	// ReasonRBACConfigurationFailure when there was a failure configuring/deploying RBAC entities on hosted clusters.
	ReasonRBACConfigurationFailure string = "ReasonRBACConfigurationFailure"
	// ReasonWebHookServiceAccountFailure when there was a failure related to the webhook's service account.
	ReasonWebHookServiceAccountFailure string = "ReasonWebHookServiceAccountFailure"
	// ReasonWebHookAPIServiceFailure when there was a failure related to the webhook's API service.
	ReasonWebHookAPIServiceFailure string = "ReasonWebHookAPIServiceFailure"
	//ReasonWebHookEndpointFailure when there was a failure related to configuring the endpoint that routes to the admission service.
	ReasonWebHookEndpointFailure string = "ReasonWebHookEndpointFailure"
	// ReasonServiceServiceAccount when there was a failure configuring/deploying the assisted-service service-account.
	ReasonServiceServiceAccount string = "ServiceServiceAccount"
	// ReasonNamespaceCreationFailure when there was a failure creating the namespace.
	ReasonNamespaceCreationFailure string = "NamespaceCreationFailure"
	// ReasonSpokeClusterCRDsSyncFailure when there was a failure syncing spoke cluster CRDs.
	ReasonSpokeClusterCRDsSyncFailure string = "SpokeClusterCRDsSyncFailure"
	// ReasonKubeconfigSecretFetchFailure when there was a failure fetching kubeconfig secret.
	ReasonKubeconfigSecretFetchFailure string = "ReasonKubeconfigSecretFetchFailure"
	// ReasonSpokeClientCreationFailure when there was a failure creating spoke client.
	ReasonSpokeClientCreationFailure string = "ReasonSpokeClientCreationFailure"
	// ReasonKonnectivityAgentFailure when there was a failure creating the namespace.
	ReasonKonnectivityAgentFailure string = "KonnectivityAgentFailure"
	// ReasonOSImageCACertRefFailure when there has been a failure resolving the OS image CA using OSImageCACertRef.
	ReasonOSImageCACertRefFailure string = "OSImageCACertRefFailure"
	// ReasonMonitoringFailure indicates there was a failure monitoring operand status
	ReasonMonitoringFailure string = "MonitoringFailure"
	// ReasonKubernetesIngressMissing indicates the user has not provided the required configuration for kubernetes ingress
	ReasonKubernetesIngressMissing string = "KubernetesIngressConfigMissing"
	// ReasonCertificateFailure indicates that the required certificates could not be created
	ReasonCertificateFailure string = "CertificateConfigurationFailure"

	// IPXEHTTPRouteEnabled is expected value in IPXEHTTPRoute to enable the route
	IPXEHTTPRouteEnabled string = "enabled"
	// IPXEHTTPRouteEnabled is expected value in IPXEHTTPRoute to disable the route
	IPXEHTTPRouteDisabled string = "disabled"
	// ReasonOSImageAdditionalParamsRefFailure when there has been a failure resolving the OS image additional params secret using OSImageAdditionalParamsRef.
	ReasonOSImageAdditionalParamsRefFailure string = "ReasonOSImageAdditionalParamsRefFailure"
)

// AgentServiceConfigStatus defines the observed state of AgentServiceConfig
type AgentServiceConfigStatus struct {
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// AgentServiceConfig represents an Assisted Service deployment.
// Only an AgentServiceConfig with name="agent" will be reconciled. All other
// names will be rejected.
// +operator-sdk:csv:customresourcedefinitions:displayName="Agent Service Config"
// +operator-sdk:csv:customresourcedefinitions:order=1
type AgentServiceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentServiceConfigSpec   `json:"spec,omitempty"`
	Status AgentServiceConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AgentServiceConfigList contains a list of AgentServiceConfig
type AgentServiceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AgentServiceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AgentServiceConfig{}, &AgentServiceConfigList{})
}
