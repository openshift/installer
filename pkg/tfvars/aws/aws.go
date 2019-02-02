package aws

// AWS converts AWS related config.
type AWS struct {
	EC2AMIOverride string            `json:"aws_ec2_ami_override,omitempty"`
	ExtraTags      map[string]string `json:"aws_extra_tags,omitempty"`
	Master         `json:",inline"`
	Region         string `json:"aws_region,omitempty"`
}

// Master converts master related config.
type Master struct {
	EC2Type          string `json:"aws_master_ec2_type,omitempty"`
	MasterRootVolume `json:",inline"`
}

// MasterRootVolume converts master rool volume related config.
type MasterRootVolume struct {
	IOPS int    `json:"aws_master_root_volume_iops,omitempty"`
	Size int    `json:"aws_master_root_volume_size,omitempty"`
	Type string `json:"aws_master_root_volume_type,omitempty"`
}
