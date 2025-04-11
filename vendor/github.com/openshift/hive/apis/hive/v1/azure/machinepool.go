package azure

// MachinePool stores the configuration for a machine pool installed
// on Azure.
type MachinePool struct {
	// Zones is list of availability zones that can be used.
	// eg. ["1", "2", "3"]
	Zones []string `json:"zones,omitempty"`

	// InstanceType defines the azure instance type.
	// eg. Standard_DS_V2
	InstanceType string `json:"type"`

	// OSDisk defines the storage for instance.
	OSDisk `json:"osDisk"`

	// OSImage defines the image to use for the OS.
	// +optional
	OSImage *OSImage `json:"osImage,omitempty"`
}

// OSImage is the image to use for the OS of a machine.
type OSImage struct {
	// Publisher is the publisher of the image.
	Publisher string `json:"publisher"`
	// Offer is the offer of the image.
	Offer string `json:"offer"`
	// SKU is the SKU of the image.
	SKU string `json:"sku"`
	// Version is the version of the image.
	Version string `json:"version"`
}
