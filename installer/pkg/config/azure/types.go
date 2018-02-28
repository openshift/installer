package azure

type component struct {
	StorageType string `yaml:"StorageType,omitempty"`
	VMSize      string `yaml:"VMSize,omitempty"`
}

type external struct {
	DNSZoneID      string `yaml:"DNSZoneID,omitempty"`
	MasterSubnetID string `yaml:"MasterSubnetID,omitempty"`
	NSG            nsg    `yaml:"NSG,omitempty"`
	ResourceGroup  string `yaml:"ResourceGroup,omitempty"`
	VNetID         string `yaml:"VNetID,omitempty"`
	WorkerSubnetID string `yaml:"WorkerSubnetID,omitempty"`
}

type network struct {
	External string `yaml:"External,omitempty"`
	Internal string `yaml:"Internal,omitempty"`
}

type nsg struct {
	MasterID string `yaml:"MasterID,omitempty"`
	WorkerID string `yaml:"WorkerID,omitempty"`
}

type ssh struct {
	Key     string  `yaml:"Key,omitempty"`
	Network network `yaml:"Network,omitempty"`
}
