package openstack

const (
	// DefaultNetworkCIDRBlock is the default CIDR range for an OpenStack network.
	DefaultNetworkCIDRBlock = "10.0.0.0/16"
	// DefaultRegion is the default OpenStack region for the cluster.
	DefaultRegion = "regionOne"
)

// OpenStack converts OpenStack related config.
type OpenStack struct {
	BaseImage        string `json:"tectonic_openstack_base_image,omitempty" yaml:"baseImage,omitempty"`
	Credentials      `json:",inline" yaml:"credentials,omitempty"`
	External         `json:",inline" yaml:"external,omitempty"`
	ExternalNetwork  string            `json:"tectonic_openstack_external_network,omitempty" yaml:"externalNetwork,omitempty"`
	ExtraTags        map[string]string `json:"tectonic_openstack_extra_tags,omitempty" yaml:"extraTags,omitempty"`
	Master           `json:",inline" yaml:"master,omitempty"`
	Region           string `json:"tectonic_openstack_region,omitempty" yaml:"region,omitempty"`
	NetworkCIDRBlock string `json:"tectonic_openstack_network_cidr_block,omitempty" yaml:"networkCIDRBlock,omitempty"`
}

// External converts external related config.
type External struct {
	MasterSubnetIDs []string `json:"tectonic_openstack_external_master_subnet_ids,omitempty" yaml:"masterSubnetIDs,omitempty"`
}

// Master converts master related config.
type Master struct {
	FlavorName string   `json:"tectonic_openstack_master_flavor_name,omitempty" yaml:"flavorName,omitempty"`
	ExtraSGIDs []string `json:"tectonic_openstack_master_extra_sg_ids,omitempty" yaml:"extraSGIDs,omitempty"`
}

// Credentials converts credentials related config.
type Credentials struct {
	AuthURL           string `json:"tectonic_openstack_credentials_auth_url,omitempty" yaml:"authUrl,omitempty"`
	Cert              string `json:"tectonic_openstack_credentials_cert,omitempty" yaml:"cert,omitempty"`
	Cloud             string `json:"tectonic_openstack_credentials_cloud,omitempty" yaml:"cloud,omitempty"`
	DomainID          string `json:"tectonic_openstack_credentials_domain_id,omitempty" yaml:"domainId,omitempty"`
	DomainName        string `json:"tectonic_openstack_credentials_domain_name,omitempty" yaml:"domainName,omitempty"`
	EndpointType      string `json:"tectonic_openstack_credentials_endpoint_type,omitempty" yaml:"endpointType,omitempty"`
	Insecure          bool   `json:"tectonic_openstack_credentials_insecure,omitempty" yaml:"insecure,omitempty"`
	Key               string `json:"tectonic_openstack_credentials_key,omitempty" yaml:"key,omitempty"`
	Password          string `json:"tectonic_openstack_credentials_password,omitempty" yaml:"password,omitempty"`
	ProjectDomainID   string `json:"tectonic_openstack_credentials_project_domain_id,omitempty" yaml:"projectDomainId,omitempty"`
	ProjectDomainName string `json:"tectonic_openstack_credentials_project_domain_name,omitempty" yaml:"projectDomainName,omitempty"`
	Region            string `json:"tectonic_openstack_credentials_region,omitempty" yaml:"region,omitempty"`
	Swauth            bool   `json:"tectonic_openstack_credentials_swauth,omitempty" yaml:"swauth,omitempty"`
	TenantID          string `json:"tectonic_openstack_credentials_tenant_id,omitempty" yaml:"tenantId,omitempty"`
	TenantName        string `json:"tectonic_openstack_credentials_tenant_name,omitempty" yaml:"tenantName,omitempty"`
	Token             string `json:"tectonic_openstack_credentials_token,omitempty" yaml:"token,omitempty"`
	UseOctavia        bool   `json:"tectonic_openstack_credentials_use_octavia,omitempty" yaml:"useOctavia,omitempty"`
	UserDomainID      string `json:"tectonic_openstack_credentials_user_domain_id,omitempty" yaml:"userDomainId,omitempty"`
	UserDomainName    string `json:"tectonic_openstack_credentials_user_domain_name,omitempty" yaml:"userDomainName,omitempty"`
	UserID            string `json:"tectonic_openstack_credentials_user_id,omitempty" yaml:"userId,omitempty"`
	UserName          string `json:"tectonic_openstack_credentials_user_name,omitempty" yaml:"userName,omitempty"`
}
