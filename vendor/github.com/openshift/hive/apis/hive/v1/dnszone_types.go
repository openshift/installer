package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/hive/apis/hive/v1/aws"
	"github.com/openshift/hive/apis/hive/v1/azure"
)

const (
	// FinalizerDNSZone is used on DNSZones to ensure we successfully deprovision
	// the cloud objects before cleaning up the API object.
	FinalizerDNSZone string = "hive.openshift.io/dnszone"

	// FinalizerDNSEndpoint is used on DNSZones to ensure we successfully
	// delete the parent-link records before cleaning up the API object.
	FinalizerDNSEndpoint string = "hive.openshift.io/dnsendpoint"
)

// DNSZoneSpec defines the desired state of DNSZone
type DNSZoneSpec struct {
	// Zone is the DNS zone to host
	Zone string `json:"zone"`

	// LinkToParentDomain specifies whether DNS records should
	// be automatically created to link this DNSZone with a
	// parent domain.
	// +optional
	LinkToParentDomain bool `json:"linkToParentDomain,omitempty"`

	// PreserveOnDelete allows the user to disconnect a DNSZone from Hive without deprovisioning it.
	// This can also be used to abandon ongoing DNSZone deprovision.
	// Typically set automatically due to PreserveOnDelete being set on a ClusterDeployment.
	// +optional
	PreserveOnDelete bool `json:"preserveOnDelete,omitempty"`

	// AWS specifies AWS-specific cloud configuration
	// +optional
	AWS *AWSDNSZoneSpec `json:"aws,omitempty"`

	// GCP specifies GCP-specific cloud configuration
	// +optional
	GCP *GCPDNSZoneSpec `json:"gcp,omitempty"`

	// Azure specifes Azure-specific cloud configuration
	// +optional
	Azure *AzureDNSZoneSpec `json:"azure,omitempty"`
}

// AWSDNSZoneSpec contains AWS-specific DNSZone specifications
type AWSDNSZoneSpec struct {
	// CredentialsSecretRef contains a reference to a secret that contains AWS credentials
	// for CRUD operations
	// +optional
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef,omitempty"`

	// CredentialsAssumeRole refers to the IAM role that must be assumed to obtain
	// AWS account access for the DNS CRUD operations.
	// +optional
	CredentialsAssumeRole *aws.AssumeRole `json:"credentialsAssumeRole,omitempty"`

	// AdditionalTags is a set of additional tags to set on the DNS hosted zone. In addition
	// to these tags,the DNS Zone controller will set a hive.openhsift.io/hostedzone tag
	// identifying the HostedZone record that it belongs to.
	AdditionalTags []AWSResourceTag `json:"additionalTags,omitempty"`

	// Region is the AWS region to use for route53 operations.
	// This defaults to us-east-1.
	// For AWS China, use cn-northwest-1.
	// +optional
	Region string `json:"region,omitempty"`
}

// AWSResourceTag represents a tag that is applied to an AWS cloud resource
type AWSResourceTag struct {
	// Key is the key for the tag
	Key string `json:"key"`
	// Value is the value for the tag
	Value string `json:"value"`
}

// GCPDNSZoneSpec contains GCP-specific DNSZone specifications
type GCPDNSZoneSpec struct {
	// CredentialsSecretRef references a secret that will be used to authenticate with
	// GCP CloudDNS. It will need permission to create and manage CloudDNS Hosted Zones.
	// Secret should have a key named 'osServiceAccount.json'.
	// The credentials must specify the project to use.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`
}

// AzureDNSZoneSpec contains Azure-specific DNSZone specifications
type AzureDNSZoneSpec struct {
	// CredentialsSecretRef references a secret that will be used to authenticate with
	// Azure CloudDNS. It will need permission to create and manage CloudDNS Hosted Zones.
	// Secret should have a key named 'osServicePrincipal.json'.
	// The credentials must specify the project to use.
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// ResourceGroupName specifies the Azure resource group in which the Hosted Zone should be created.
	ResourceGroupName string `json:"resourceGroupName"`

	// CloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK
	// with the appropriate Azure API endpoints.
	// If empty, the value is equal to "AzurePublicCloud".
	// +optional
	CloudName azure.CloudEnvironment `json:"cloudName,omitempty"`
}

// DNSZoneStatus defines the observed state of DNSZone
type DNSZoneStatus struct {
	// LastSyncTimestamp is the time that the zone was last sync'd.
	// +optional
	LastSyncTimestamp *metav1.Time `json:"lastSyncTimestamp,omitempty"`

	// LastSyncGeneration is the generation of the zone resource that was last sync'd. This is used to know
	// if the Object has changed and we should sync immediately.
	// +optional
	LastSyncGeneration int64 `json:"lastSyncGeneration,omitempty"`

	// NameServers is a list of nameservers for this DNS zone
	// +optional
	NameServers []string `json:"nameServers,omitempty"`

	// AWSDNSZoneStatus contains status information specific to AWS
	// +optional
	AWS *AWSDNSZoneStatus `json:"aws,omitempty"`

	// GCPDNSZoneStatus contains status information specific to GCP
	// +optional
	GCP *GCPDNSZoneStatus `json:"gcp,omitempty"`

	// AzureDNSZoneStatus contains status information specific to Azure
	Azure *AzureDNSZoneStatus `json:"azure,omitempty"`

	// Conditions includes more detailed status for the DNSZone
	// +optional
	Conditions []DNSZoneCondition `json:"conditions,omitempty"`
}

// AWSDNSZoneStatus contains status information specific to AWS DNS zones
type AWSDNSZoneStatus struct {
	// ZoneID is the ID of the zone in AWS
	// +optional
	ZoneID *string `json:"zoneID,omitempty"`
}

// AzureDNSZoneStatus contains status information specific to Azure DNS zones
type AzureDNSZoneStatus struct {
}

// GCPDNSZoneStatus contains status information specific to GCP Cloud DNS zones
type GCPDNSZoneStatus struct {
	// ZoneName is the name of the zone in GCP Cloud DNS
	// +optional
	ZoneName *string `json:"zoneName,omitempty"`
}

// DNSZoneCondition contains details for the current condition of a DNSZone
type DNSZoneCondition struct {
	// Type is the type of the condition.
	Type DNSZoneConditionType `json:"type"`
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

// DNSZoneConditionType is a valid value for DNSZoneCondition.Type
type DNSZoneConditionType string

// ConditionType satisfies the conditions.Condition interface
func (c DNSZoneCondition) ConditionType() ConditionType {
	return c.Type
}

// String satisfies the conditions.ConditionType interface
func (t DNSZoneConditionType) String() string {
	return string(t)
}

const (
	// ZoneAvailableDNSZoneCondition is true if the DNSZone is responding to DNS queries
	ZoneAvailableDNSZoneCondition DNSZoneConditionType = "ZoneAvailable"
	// ParentLinkCreatedCondition is true if the parent link has been created
	ParentLinkCreatedCondition DNSZoneConditionType = "ParentLinkCreated"
	// DomainNotManaged is true if we try to reconcile a DNSZone and the HiveConfig
	// does not contain a ManagedDNS entry for the domain in the DNSZone
	DomainNotManaged DNSZoneConditionType = "DomainNotManaged"
	// InsufficientCredentialsCondition is true when credentials cannot be used to create a
	// DNS zone because of insufficient permissions
	InsufficientCredentialsCondition DNSZoneConditionType = "InsufficientCredentials"
	// AuthenticationFailureCondition is true when credentials cannot be used to create a
	// DNS zone because they fail authentication
	AuthenticationFailureCondition DNSZoneConditionType = "AuthenticationFailure"
	// APIOptInRequiredCondition is true when the user account used for managing DNS
	// needs to enable the DNS apis.
	APIOptInRequiredCondition DNSZoneConditionType = "APIOptInRequired"
	// GenericDNSErrorsCondition is true when there's some DNS Zone related error that isn't related to
	// authentication or credentials, and needs to be bubbled up to ClusterDeployment
	GenericDNSErrorsCondition DNSZoneConditionType = "DNSError"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSZone is the Schema for the dnszones API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
type DNSZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DNSZoneSpec   `json:"spec,omitempty"`
	Status DNSZoneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DNSZoneList contains a list of DNSZone
type DNSZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DNSZone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DNSZone{}, &DNSZoneList{})
}
