package aws

// AWS converts AWS related config.
type AWS struct {
	EC2AMIOverride string            `json:"aws_ec2_ami_override,omitempty"`
	ExtraTags      map[string]string `json:"aws_extra_tags,omitempty"`
	ControlPlane   `json:",inline"`
	Region         string `json:"aws_region,omitempty"`
	Compute        `json:",inline"`
}

// ControlPlane converts control plane related config.
type ControlPlane struct {
	EC2Type                string `json:"aws_controlplane_ec2_type,omitempty"`
	IAMRoleName            string `json:"aws_controlplane_iam_role_name,omitempty"`
	ControlPlaneRootVolume `json:",inline"`
}

// ControlPlaneRootVolume converts control plane rool volume related config.
type ControlPlaneRootVolume struct {
	IOPS int    `json:"aws_controlplane_root_volume_iops,omitempty"`
	Size int    `json:"aws_controlplane_root_volume_size,omitempty"`
	Type string `json:"aws_controlplane_root_volume_type,omitempty"`
}

// Compute converts compute related config.
type Compute struct {
	IAMRoleName string `json:"aws_compute_iam_role_name,omitempty"`
}
