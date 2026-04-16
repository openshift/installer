package prismgoclient

// Credentials can include either username and password for basic authentication
// or an API key for API key-based authentication
type Credentials struct {
	URL                string
	APIKey             string
	Username           string
	Password           string
	Endpoint           string
	Port               string
	Insecure           bool
	SessionAuth        bool
	ProxyURL           string
	FoundationEndpoint string              // Required field for connecting to foundation VM APIs
	FoundationPort     string              // Port for connecting to foundation VM APIs
	RequiredFields     map[string][]string // RequiredFields is client to its required fields mapping for validations and usage in every client
}

// AdditionalFilter specification for client side filters
type AdditionalFilter struct {
	Name   string
	Values []string
}
