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
	DiskType string `json:"diskType,omitempty"`

	// DiskEncryptionSet defines a disk encryption set.
	//
	// +optional
	*DiskEncryptionSet `json:"diskEncryptionSet,omitempty"`
}

// DiskEncryptionSet defines the configuration for a disk encryption set.
type DiskEncryptionSet struct {
	// SubscriptionID defines the Azure subscription the disk encryption
	// set is in.
	SubscriptionID string `json:"subscriptionId,omitempty"`
	// ResourceGroup defines the Azure resource group used by the disk
	// encryption set.
	ResourceGroup string `json:"resourceGroup"`
	// Name is the name of the disk encryption set.
	Name string `json:"name"`
}

// DefaultDiskType holds the default Azure disk type used by the VMs.
const DefaultDiskType string = "Premium_LRS"
