package alibabacloud

// Platform stores all the global configuration that all machinesets use.
type Platform struct {
	// Region specifies the Alibaba Cloud region where the cluster will be created.
	Region            string `json:"region"`
	ResourceGroupName string `json:"resourceGroupName"`
}

func (p *Platform) SetBaseDomain() error {
	// pass
	return nil
}
