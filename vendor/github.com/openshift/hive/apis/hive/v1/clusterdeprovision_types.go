package v1

import (
	"github.com/openshift/hive/apis/hive/v1/aws"
	"github.com/openshift/hive/apis/hive/v1/azure"
	"github.com/openshift/hive/apis/hive/v1/nutanix"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterDeprovisionSpec defines the desired state of ClusterDeprovision
type ClusterDeprovisionSpec struct {
	// InfraID is the identifier generated during installation for a cluster. It is used for tagging/naming resources in cloud providers.
	InfraID string `json:"infraID"`

	// ClusterID is a globally unique identifier for the cluster to deprovision. It will be used if specified.
	ClusterID string `json:"clusterID,omitempty"`

	// ClusterName is the friendly name of the cluster. It is used for subdomains,
	// some resource tagging, and other instances where a friendly name for the
	// cluster is useful.
	ClusterName string `json:"clusterName,omitempty"`

	// BaseDomain is the DNS base domain.
	BaseDomain string `json:"baseDomain,omitempty"`

	// MetadataJSONSecretRef references the secret containing the metadata.json emitted by the
	// installer, potentially scrubbed for sensitive data.
	MetadataJSONSecretRef *corev1.LocalObjectReference `json:"metadataJSONSecretRef,omitempty"`

	// Platform contains platform-specific configuration for a ClusterDeprovision
	Platform ClusterDeprovisionPlatform `json:"platform,omitempty"`
}

// ClusterDeprovisionStatus defines the observed state of ClusterDeprovision
type ClusterDeprovisionStatus struct {
	// Completed is true when the uninstall has completed successfully
	Completed bool `json:"completed,omitempty"`

	// Conditions includes more detailed status for the cluster deprovision
	// +optional
	Conditions []ClusterDeprovisionCondition `json:"conditions,omitempty"`
}

// ClusterDeprovisionPlatform contains platform-specific configuration for the
// deprovision
type ClusterDeprovisionPlatform struct {
	// AWS contains AWS-specific deprovision settings
	AWS *AWSClusterDeprovision `json:"aws,omitempty"`
	// Azure contains Azure-specific deprovision settings
	Azure *AzureClusterDeprovision `json:"azure,omitempty"`
	// GCP contains GCP-specific deprovision settings
	GCP *GCPClusterDeprovision `json:"gcp,omitempty"`
	// OpenStack contains OpenStack-specific deprovision settings
	OpenStack *OpenStackClusterDeprovision `json:"openstack,omitempty"`
	// VSphere contains VMWare vSphere-specific deprovision settings
	VSphere *VSphereClusterDeprovision `json:"vsphere,omitempty"`
	// IBMCloud contains IBM Cloud specific deprovision settings
	IBMCloud *IBMClusterDeprovision `json:"ibmcloud,omitempty"`
	// Nutanix contains Nutanix-specific deprovision settings
	Nutanix *NutanixClusterDeprovision `json:"nutanix,omitempty"`
}

// AWSClusterDeprovision contains AWS-specific configuration for a ClusterDeprovision
type AWSClusterDeprovision struct {
	// Region is the AWS region for this deprovisioning
	Region string `json:"region"`

	// CredentialsSecretRef is the AWS account credentials to use for deprovisioning the cluster
	// +optional
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// CredentialsAssumeRole refers to the IAM role that must be assumed to obtain
	// AWS account access for deprovisioning the cluster.
	// +optional
	CredentialsAssumeRole *aws.AssumeRole `json:"credentialsAssumeRole,omitempty"`

	// HostedZoneRole is the role to assume when performing operations
	// on a hosted zone owned by another account.
	// +optional
	HostedZoneRole *string `json:"hostedZoneRole,omitempty"`
}

// AzureClusterDeprovision contains Azure-specific configuration for a ClusterDeprovision
type AzureClusterDeprovision struct {
	// CredentialsSecretRef is the Azure account credentials to use for deprovisioning the cluster
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
	// cloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK
	// with the appropriate Azure API endpoints.
	// If empty, the value is equal to "AzurePublicCloud".
	// +optional
	CloudName *azure.CloudEnvironment `json:"cloudName,omitempty"`
	// ResourceGroupName is the name of the resource group where the cluster was installed.
	// Required for new deprovisions (schema notwithstanding).
	// +optional
	ResourceGroupName *string `json:"resourceGroupName,omitempty"`
	// BaseDomainResourceGroupName is the name of the resource group where the cluster's DNS records
	// were created, if different from the default or the custom ResourceGroupName.
	// +optional
	BaseDomainResourceGroupName *string `json:"baseDomainResourceGroupName,omitempty"`
}

// GCPClusterDeprovision contains GCP-specific configuration for a ClusterDeprovision
type GCPClusterDeprovision struct {
	// Region is the GCP region for this deprovision
	Region string `json:"region"`
	// CredentialsSecretRef is the GCP account credentials to use for deprovisioning the cluster
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// NetworkProjectID is used for shared VPC setups
	// +optional
	NetworkProjectID *string `json:"networkProjectID,omitempty"`
}

// OpenStackClusterDeprovision contains OpenStack-specific configuration for a ClusterDeprovision
type OpenStackClusterDeprovision struct {
	// Cloud is the secion in the clouds.yaml secret below to use for auth/connectivity.
	Cloud string `json:"cloud"`
	// CredentialsSecretRef is the OpenStack account credentials to use for deprovisioning the cluster
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
	// CertificatesSecretRef refers to a secret that contains CA certificates
	// necessary for communicating with the OpenStack.
	//
	// +optional
	CertificatesSecretRef *corev1.LocalObjectReference `json:"certificatesSecretRef,omitempty"`
}

// VSphereClusterDeprovision contains VMware vSphere-specific configuration for a ClusterDeprovision
type VSphereClusterDeprovision struct {
	// CredentialsSecretRef is the vSphere account credentials to use for deprovisioning the cluster
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// CertificatesSecretRef refers to a secret that contains the vSphere CA certificates
	// necessary for communicating with the VCenter.
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`
	// VCenter is the vSphere vCenter hostname.
	VCenter string `json:"vCenter"`
}

// IBMClusterDeprovision contains IBM Cloud specific configuration for a ClusterDeprovision
type IBMClusterDeprovision struct {
	// CredentialsSecretRef is the IBM Cloud credentials to use for deprovisioning the cluster
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// Region specifies the IBM Cloud region
	Region string `json:"region"`
	// BaseDomain is the DNS base domain.
	// TODO: Use the non-platform-specific BaseDomain field.
	BaseDomain string `json:"baseDomain"`
}

// NutanixClusterDeprovision contains Nutanix-specific configuration for a ClusterDeprovision
type NutanixClusterDeprovision struct {
	// PrismCentral is the endpoint (address and port) to connect to the Prism Central.
	// This serves as the default Prism-Central.
	PrismCentral nutanix.PrismEndpoint `json:"prismCentral"`

	// CredentialsSecretRef refers to a secret that contains the Nutanix account access
	// credentials.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// CertificatesSecretRef refers to a secret that contains the Nutanix CA certificates
	// necessary for communicating with the Prism Central.
	// +optional
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterDeprovision is the Schema for the clusterdeprovisions API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="InfraID",type="string",JSONPath=".spec.infraID"
// +kubebuilder:printcolumn:name="ClusterID",type="string",JSONPath=".spec.clusterID"
// +kubebuilder:printcolumn:name="Completed",type="boolean",JSONPath=".status.completed"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=clusterdeprovisions,shortName=cdr,scope=Namespaced
type ClusterDeprovision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterDeprovisionSpec   `json:"spec,omitempty"`
	Status ClusterDeprovisionStatus `json:"status,omitempty"`
}

// ClusterDeprovisionCondition contains details for the current condition of a ClusterDeprovision
type ClusterDeprovisionCondition struct {
	// Type is the type of the condition.
	Type ClusterDeprovisionConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

// ClusterDeprovisionConditionType is a valid value for ClusterDeprovisionCondition.Type
type ClusterDeprovisionConditionType string

// ConditionType satisfies the conditions.Condition interface
func (c ClusterDeprovisionCondition) ConditionType() ConditionType {
	return c.Type
}

// String satisfies the conditions.ConditionType interface
func (t ClusterDeprovisionConditionType) String() string {
	return string(t)
}

const (
	// AuthenticationFailureClusterDeprovisionCondition is true when credentials cannot be used because of authentication failure
	AuthenticationFailureClusterDeprovisionCondition ClusterDeprovisionConditionType = "AuthenticationFailure"

	// DeprovisionFailedClusterDeprovisionCondition is true when deprovision attempt failed
	DeprovisionFailedClusterDeprovisionCondition ClusterDeprovisionConditionType = "DeprovisionFailed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterDeprovisionList contains a list of ClusterDeprovision
type ClusterDeprovisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterDeprovision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterDeprovision{}, &ClusterDeprovisionList{})
}
