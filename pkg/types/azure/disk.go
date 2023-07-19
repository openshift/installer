package azure

import "fmt"

// ToID creates an Azure resource ID for the disk encryption set.
// It is possible to return a non-valid ID when SubscriptionID is empty. This
// should never happen since if SubscriptionID is empty, it is set to the
// current subscription. Also, should it somehow be empty and this returns an
// invalid ID, the validation code will produce an error when checked  against
// the validation.RxDiskEncryptionSetID regular expression.
func (d *DiskEncryptionSet) ToID() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s",
		d.SubscriptionID, d.ResourceGroup, d.Name)
}

// SecurityEncryptionTypes represents the Encryption Type when the Azure Virtual Machine is a
// Confidential VM.
type SecurityEncryptionTypes string

const (
	// SecurityEncryptionTypesVMGuestStateOnly disables OS disk confidential encryption.
	SecurityEncryptionTypesVMGuestStateOnly SecurityEncryptionTypes = "VMGuestStateOnly"
	// SecurityEncryptionTypesDiskWithVMGuestState enables OS disk confidential encryption with
	// a platform-managed key (PMK) or a customer-managed key (CMK).
	SecurityEncryptionTypesDiskWithVMGuestState SecurityEncryptionTypes = "DiskWithVMGuestState"
)

// OSDisk defines the disk for machines on Azure.
type OSDisk struct {
	// DiskSizeGB defines the size of disk in GB.
	//
	// +kubebuilder:validation:Minimum=0
	DiskSizeGB int32 `json:"diskSizeGB"`
	// DiskType defines the type of disk.
	// For control plane nodes, the valid values are Premium_LRS and StandardSSD_LRS.
	// Default is Premium_LRS.
	// +optional
	// +kubebuilder:validation:Enum=Standard_LRS;Premium_LRS;StandardSSD_LRS
	DiskType string `json:"diskType"`
	// DiskEncryptionSet defines a disk encryption set.
	//
	// +optional
	*DiskEncryptionSet `json:"diskEncryptionSet,omitempty"`
	// SecurityProfile specifies the security profile for the managed disk.
	// +optional
	SecurityProfile *VMDiskSecurityProfile `json:"securityProfile,omitempty"`
}

// DiskEncryptionSet defines the configuration for a disk encryption set.
type DiskEncryptionSet struct {
	// SubscriptionID defines the Azure subscription the disk encryption
	// set is in.
	SubscriptionID string `json:"subscriptionId"`
	// ResourceGroup defines the Azure resource group used by the disk
	// encryption set.
	ResourceGroup string `json:"resourceGroup"`
	// Name is the name of the disk encryption set.
	Name string `json:"name"`
}

// VMDiskSecurityProfile specifies the security profile settings for the managed disk.
// It can be set only for Confidential VMs.
type VMDiskSecurityProfile struct {
	// DiskEncryptionSet specifies the customer managed disk encryption set resource id for the
	// managed disk that is used for Customer Managed Key encrypted ConfidentialVM OS Disk and
	// VMGuestState blob.
	// +optional
	DiskEncryptionSet *DiskEncryptionSet `json:"diskEncryptionSet,omitempty"`
	// SecurityEncryptionType specifies the encryption type of the managed disk.
	// It is set to DiskWithVMGuestState to encrypt the managed disk along with the VMGuestState
	// blob, and to VMGuestStateOnly to encrypt the VMGuestState blob only.
	// When set to VMGuestStateOnly, the VTpmEnabled should be set to true.
	// When set to DiskWithVMGuestState, both SecureBootEnabled and VTpmEnabled should be set to true.
	// It can be set only for Confidential VMs.
	// +kubebuilder:validation:Enum=VMGuestStateOnly;DiskWithVMGuestState
	// +optional
	SecurityEncryptionType SecurityEncryptionTypes `json:"securityEncryptionType,omitempty"`
}

// DefaultDiskType holds the default Azure disk type used by the VMs.
const DefaultDiskType string = "Premium_LRS"
