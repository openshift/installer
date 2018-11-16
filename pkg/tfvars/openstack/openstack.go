package openstack

const (
	// DefaultNetworkCIDRBlock is the default CIDR range for an OpenStack network.
	DefaultNetworkCIDRBlock = "10.0.0.0/16"
	// DefaultRegion is the default OpenStack region for the cluster.
	DefaultRegion = "regionOne"
)

// OpenStack converts OpenStack related config.
type OpenStack struct {
	BaseImage        string `json:"tectonic_openstack_base_image,omitempty"`
	Credentials      `json:",inline"`
	External         `json:",inline"`
	ExternalNetwork  string            `json:"tectonic_openstack_external_network,omitempty"`
	ExtraTags        map[string]string `json:"tectonic_openstack_extra_tags,omitempty"`
	Master           `json:",inline"`
	Region           string `json:"tectonic_openstack_region,omitempty"`
	NetworkCIDRBlock string `json:"tectonic_openstack_network_cidr_block,omitempty"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"tectonic_openstack_external_master_subnet_ids,omitempty"`
}

// Master converts master related config.
type Master struct {
	FlavorName string   `json:"tectonic_openstack_master_flavor_name,omitempty"`
	ExtraSGIDs []string `json:"tectonic_openstack_master_extra_sg_ids,omitempty"`
}

// Credentials converts credentials related config.
type Credentials struct {
	AuthURL           string `json:"tectonic_openstack_credentials_auth_url,omitempty"`
	Cert              string `json:"tectonic_openstack_credentials_cert,omitempty"`
	Cloud             string `json:"tectonic_openstack_credentials_cloud,omitempty"`
	DomainID          string `json:"tectonic_openstack_credentials_domain_id,omitempty"`
	DomainName        string `json:"tectonic_openstack_credentials_domain_name,omitempty"`
	EndpointType      string `json:"tectonic_openstack_credentials_endpoint_type,omitempty"`
	Insecure          bool   `json:"tectonic_openstack_credentials_insecure,omitempty"`
	Key               string `json:"tectonic_openstack_credentials_key,omitempty"`
	Password          string `json:"tectonic_openstack_credentials_password,omitempty"`
	ProjectDomainID   string `json:"tectonic_openstack_credentials_project_domain_id,omitempty"`
	ProjectDomainName string `json:"tectonic_openstack_credentials_project_domain_name,omitempty"`
	Region            string `json:"tectonic_openstack_credentials_region,omitempty"`
	Swauth            bool   `json:"tectonic_openstack_credentials_swauth,omitempty"`
	TenantID          string `json:"tectonic_openstack_credentials_tenant_id,omitempty"`
	TenantName        string `json:"tectonic_openstack_credentials_tenant_name,omitempty"`
	Token             string `json:"tectonic_openstack_credentials_token,omitempty"`
	UseOctavia        bool   `json:"tectonic_openstack_credentials_use_octavia,omitempty"`
	UserDomainID      string `json:"tectonic_openstack_credentials_user_domain_id,omitempty"`
	UserDomainName    string `json:"tectonic_openstack_credentials_user_domain_name,omitempty"`
	UserID            string `json:"tectonic_openstack_credentials_user_id,omitempty"`
	UserName          string `json:"tectonic_openstack_credentials_user_name,omitempty"`
}
