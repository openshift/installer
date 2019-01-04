package openstack

// OpenStack converts OpenStack related config.
type OpenStack struct {
	BaseImage        string `json:"openstack_base_image,omitempty"`
	Credentials      `json:",inline"`
	External         `json:",inline"`
	ExternalNetwork  string            `json:"openstack_external_network,omitempty"`
	ExtraTags        map[string]string `json:"openstack_extra_tags,omitempty"`
	Master           `json:",inline"`
	Region           string `json:"openstack_region,omitempty"`
	NetworkCIDRBlock string `json:"openstack_network_cidr_block,omitempty"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"openstack_external_master_subnet_ids,omitempty"`
}

// Master converts master related config.
type Master struct {
	FlavorName string   `json:"openstack_master_flavor_name,omitempty"`
	ExtraSGIDs []string `json:"openstack_master_extra_sg_ids,omitempty"`
}

// Credentials converts credentials related config.
type Credentials struct {
	AuthURL           string `json:"openstack_credentials_auth_url,omitempty"`
	Cert              string `json:"openstack_credentials_cert,omitempty"`
	Cloud             string `json:"openstack_credentials_cloud,omitempty"`
	DomainID          string `json:"openstack_credentials_domain_id,omitempty"`
	DomainName        string `json:"openstack_credentials_domain_name,omitempty"`
	EndpointType      string `json:"openstack_credentials_endpoint_type,omitempty"`
	Insecure          bool   `json:"openstack_credentials_insecure,omitempty"`
	Key               string `json:"openstack_credentials_key,omitempty"`
	Password          string `json:"openstack_credentials_password,omitempty"`
	ProjectDomainID   string `json:"openstack_credentials_project_domain_id,omitempty"`
	ProjectDomainName string `json:"openstack_credentials_project_domain_name,omitempty"`
	Region            string `json:"openstack_credentials_region,omitempty"`
	Swauth            bool   `json:"openstack_credentials_swauth,omitempty"`
	TenantID          string `json:"openstack_credentials_tenant_id,omitempty"`
	TenantName        string `json:"openstack_credentials_tenant_name,omitempty"`
	Token             string `json:"openstack_credentials_token,omitempty"`
	UseOctavia        bool   `json:"openstack_credentials_use_octavia,omitempty"`
	UserDomainID      string `json:"openstack_credentials_user_domain_id,omitempty"`
	UserDomainName    string `json:"openstack_credentials_user_domain_name,omitempty"`
	UserID            string `json:"openstack_credentials_user_id,omitempty"`
	UserName          string `json:"openstack_credentials_user_name,omitempty"`
}
