package karbon

type PrivateRegistryListResponse []PrivateRegistryResponse

type PrivateRegistryResponse struct {
	Name     *string `json:"name" mapstructure:"name, omitempty"`
	Endpoint *string `json:"endpoint" mapstructure:"endpoint, omitempty"`
	UUID     *string `json:"uuid" mapstructure:"uuid, omitempty"`
}

type PrivateRegistryIntentInput struct {
	Name     *string `json:"name" mapstructure:"name, omitempty"`
	Cert     *string `json:"cert" mapstructure:"cert, omitempty"`
	URL      *string `json:"url" mapstructure:"url, omitempty"`
	Port     *int64  `json:"port" mapstructure:"port, omitempty"`
	Username *string `json:"username" mapstructure:"username, omitempty"`
	Password *string `json:"password" mapstructure:"password, omitempty"`
}

type PrivateRegistryOperationResponse struct {
	RegistryName *string `json:"registry_name" mapstructure:"registry_name, omitempty"`
}

type PrivateRegistryOperationIntentInput struct {
	RegistryName *string `json:"registry_name" mapstructure:"registry_name, omitempty"`
}
