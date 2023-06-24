package openstack

type terraformFixedIP struct {
	SubnetID  string `json:"subnet_id"`
	IPAddress string `json:"ip_address"`
}

type terraformPort struct {
	NetworkID string             `json:"network_id"`
	FixedIP   []terraformFixedIP `json:"fixed_ips"`
}
