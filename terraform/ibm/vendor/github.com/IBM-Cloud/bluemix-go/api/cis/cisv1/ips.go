package cisv1

import (
	"fmt"
	"github.com/IBM-Cloud/bluemix-go/client"
)

type IpsList struct {
	Ipv4 []string `json:"ipv4_cidrs"`
	Ipv6 []string `json:"ipv6_cidrs"`
}

type IpsResults struct {
	IpList      IpsList      `json:"result"`
	ResultsInfo ResultsCount `json:"result_info"`
	Success     bool         `json:"success"`
	Errors      []Error      `json:"errors"`
}

type Ips interface {
	ListIps() (*IpsList, error)
}

type ips struct {
	client *client.Client
}

func newIpsAPI(c *client.Client) Ips {
	return &ips{
		client: c,
	}
}

func (r *ips) ListIps() (*IpsList, error) {
	ipsResults := IpsResults{}
	rawURL := fmt.Sprintf("/v1/ips")
	_, err := r.client.Get(rawURL, &ipsResults)
	if err != nil {
		return nil, err
	}
	return &ipsResults.IpList, err
}
