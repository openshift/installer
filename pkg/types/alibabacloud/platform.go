package alibabacloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region            string `json:"region"`
	ResourceGroupName string `json:"resourceGroupName"`
}

// SetBaseDomain parses the baseDomainID and sets the related fields on alibabacloud.Platform
func (p *Platform) SetBaseDomain() error {
	// pass
	return nil
}
