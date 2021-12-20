package containerv2

import (
	"encoding/json"
	"fmt"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type nlbdns struct {
	client *client.Client
}

//Clusters interface
type Nlbdns interface {
	GetNLBDNSList(clusterNameOrID string) ([]NlbVPCListConfig, error)
	GetLocationNLBDNSList(location string) ([]NlbVPCListConfig, error)
}
type NlbVPCListConfig struct {
	// ExtendedNlbVPCConfig is the response body for the get v2 vpc apis.
	Nlb ExtendedNlbVPCConfig `json:"Nlb,omitempty"`

	SecretName string `json:"secretName,omitempty"`

	SecretStatus string `json:"secretStatus,omitempty"`
}
type ExtendedNlbVPCConfig struct {
	Cluster string `json:"cluster,omitempty"`

	DnsType string `json:"dnsType,omitempty"`

	LbHostname string `json:"lbHostname,omitempty"`

	NlbIPArray []interface{} `json:"nlbIPArray,omitempty"`

	NlbSubdomain string `json:"nlbSubdomain,omitempty"`

	SecretNamespace string `json:"secretNamespace,omitempty"`
	NlbMonitorState string `json:"nlbMonitorState,omitempty"`

	Type string `json:"type,omitempty"`
}

func newNlbdnsAPI(c *client.Client) Nlbdns {
	return &nlbdns{
		client: c,
	}
}

// GetNLBDNSList returns the list of NLBDNS available for cluster
func (r *nlbdns) GetNLBDNSList(clusterNameOrID string) ([]NlbVPCListConfig, error) {
	var success []NlbVPCListConfig
	var successV []interface{}
	rawURL := fmt.Sprintf("/v2/nlb-dns/getNlbDNSList?cluster=%s", clusterNameOrID)
	_, err := r.client.Get(rawURL, &successV)
	if err != nil {
		return success, err
	}
	if len(successV) > 0 {
		if _, isVpc := successV[0].(map[string]interface{})["Nlb"]; isVpc {
			bodyBytes, _ := json.Marshal(successV)
			json.Unmarshal(bodyBytes, &success)
		} else {
			nlb := NlbVPCListConfig{}
			for _, s := range successV {
				b := s.(map[string]interface{})
				nlb.SecretName = b["nlbSslSecretName"].(string)
				nlb.SecretStatus = b["nlbSslSecretStatus"].(string)
				nlb.Nlb.Cluster = b["clusterID"].(string)
				nlb.Nlb.Type = b["nlbType"].(string)
				nlb.Nlb.DnsType = b["nlbDnsType"].(string)
				nlb.Nlb.NlbSubdomain = b["nlbHost"].(string)
				nlb.Nlb.SecretNamespace = b["secretNamespace"].(string)
				nlb.Nlb.NlbIPArray = b["nlbIPArray"].([]interface{})
				nlb.Nlb.NlbMonitorState = b["nlbMonitorState"].(string)
				success = append(success, nlb)
			}
		}
	}
	return success, err
}
func (r *nlbdns) GetLocationNLBDNSList(location string) ([]NlbVPCListConfig, error) {
	var success []NlbVPCListConfig
	var successV []interface{}
	rawURL := fmt.Sprintf("/v2/nlb-dns/getSatLocationSubdomains?controller=%s", location)
	_, err := r.client.Get(rawURL, &successV)
	if err != nil {
		return success, err
	}
	if len(successV) > 0 {
		if _, isVpc := successV[0].(map[string]interface{})["Nlb"]; isVpc {
			bodyBytes, _ := json.Marshal(successV)
			json.Unmarshal(bodyBytes, &success)
		} else {
			nlb := NlbVPCListConfig{}
			for _, s := range successV {
				b := s.(map[string]interface{})
				nlb.SecretName = b["nlbSslSecretName"].(string)
				nlb.SecretStatus = b["nlbSslSecretStatus"].(string)
				nlb.Nlb.Cluster = b["clusterID"].(string)
				nlb.Nlb.Type = b["nlbType"].(string)
				nlb.Nlb.DnsType = b["nlbDnsType"].(string)
				nlb.Nlb.NlbSubdomain = b["nlbHost"].(string)
				nlb.Nlb.SecretNamespace = b["secretNamespace"].(string)
				nlb.Nlb.NlbIPArray = b["nlbIPArray"].([]interface{})
				nlb.Nlb.NlbMonitorState = b["nlbMonitorState"].(string)
				success = append(success, nlb)
			}
		}
	}
	return success, err
}
