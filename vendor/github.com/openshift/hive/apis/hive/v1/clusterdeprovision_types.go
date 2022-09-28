package v1

import (
	"github.com/openshift/hive/apis/hive/v1/aws"
	"github.com/openshift/hive/apis/hive/v1/azure"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterDeprovisionSpec defines the desired state of ClusterDeprovision
type ClusterDeprovisionSpec struct {
	// InfraID is the identifier generated during installation for a cluster. It is used for tagging/naming resources in cloud providers.
	InfraID string `json:"infraID"`

	// ClusterID is a globally unique identifier for the cluster to deprovision. It will be used if specified.
	ClusterID string `json:"clusterID,omitempty"`

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
	// Ovirt contains oVirt-specific deprovision settings
	Ovirt *OvirtClusterDeprovision `json:"ovirt,omitempty"`
	// IBMCloud contains IBM Cloud specific deprovision settings
	IBMCloud *IBMClusterDeprovision `json:"ibmcloud,omitempty"`
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
}

// GCPClusterDeprovision contains GCP-specific configuration for a ClusterDeprovision
type GCPClusterDeprovision struct {
	// Region is the GCP region for this deprovision
	Region string `json:"region"`
	// CredentialsSecretRef is the GCP account credentials to use for deprovisioning the cluster
	CredentialsSecretRef *corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`
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

// OvirtClusterDeprovision contains oVirt-specific configuration for a ClusterDeprovision
type OvirtClusterDeprovision struct {
	// The oVirt cluster ID
	ClusterID string `json:"clusterID"`
	// CredentialsSecretRef is the oVirt account credentials to use for deprovisioning the cluster
	// secret fields: ovirt_url, ovirt_username, ovirt_password, ovirt_ca_bundle
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// CertificatesSecretRef refers to a secret that contains the oVirt CA certificates
	// necessary for communicating with the oVirt.
	CertificatesSecretRef corev1.LocalObjectReference `json:"certificatesSecretRef"`
}

// IBMClusterDeprovision contains IBM Cloud specific configuration for a ClusterDeprovision
type IBMClusterDeprovision struct {
	// CredentialsSecretRef is the IBM Cloud credentials to use for deprovisioning the cluster
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
	// AccountID is the IBM Cloud Account ID
	AccountID string `json:"accountID"`
	// CISInstanceCRN is the IBM Cloud Internet Services Instance CRN
	CISInstanceCRN string `json:"cisInstanceCRN"`
	// Region specifies the IBM Cloud region
	Region string `json:"region"`
	// BaseDomain is the DNS base domain
	BaseDomain string `json:"baseDomain"`
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
