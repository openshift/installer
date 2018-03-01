package azure

type component struct {
	StorageType string `yaml:"storageType,omitempty"`
	VMSize      string `yaml:"vmSize,omitempty"`
}

type external struct {
	DNSZoneID      string `yaml:"dnsZoneID,omitempty"`
	MasterSubnetID string `yaml:"masterSubnetID,omitempty"`
	NSG            nsg    `yaml:"nsg,omitempty"`
	ResourceGroup  string `yaml:"resourceGroup,omitempty"`
	VNetID         string `yaml:"vNetID,omitempty"`
	WorkerSubnetID string `yaml:"workerSubnetID,omitempty"`
}

type network struct {
	External string `yaml:"external,omitempty"`
	Internal string `yaml:"internal,omitempty"`
}

type nsg struct {
	MasterID string `yaml:"masterID,omitempty"`
	WorkerID string `yaml:"workerID,omitempty"`
}

type ssh struct {
	Key     string  `yaml:"key,omitempty"`
	Network network `yaml:"network,omitempty"`
}
