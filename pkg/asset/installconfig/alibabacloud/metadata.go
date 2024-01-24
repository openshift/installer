package alibabacloud

// Metadata holds additional metadata for InstallConfig resources that
// does not need to be user-supplied (e.g. because it can be retrieved
// from external APIs).
type Metadata struct {
	client      *Client
	vswitchMaps map[string]string

	VSwitchIDs []string
	Region     string `json:"region,omitempty"`
}

// NewMetadata initializes a new Metadata object.
func NewMetadata(region string, vswitchIDs []string) *Metadata {
	return &Metadata{
		Region:     region,
		VSwitchIDs: vswitchIDs,
	}
}

// Client returns a client used for making API calls to Alibaba Cloud services.
func (m *Metadata) Client() (*Client, error) {
	if m.client == nil {
		client, err := NewClient(m.Region)
		if err != nil {
			return nil, err
		}
		m.client = client
	}
	return m.client, nil
}

// VSwitchMaps retrieves VSwitch and availability zone metadata indexed by VSwitch ID
func (m *Metadata) VSwitchMaps() (vswitchMaps map[string]string, err error) {
	vswitchMaps = map[string]string{}
	if len(m.vswitchMaps) == 0 {
		client, err := m.Client()
		if err != nil {
			return nil, err
		}
		for _, vswitchID := range m.VSwitchIDs {
			response, err := client.ListVSwitches(vswitchID)
			if err != nil {
				return nil, err
			}
			vswitchMaps[response.VSwitches.VSwitch[0].ZoneId] = vswitchID
		}
		m.vswitchMaps = vswitchMaps
	}
	return m.vswitchMaps, nil
}
