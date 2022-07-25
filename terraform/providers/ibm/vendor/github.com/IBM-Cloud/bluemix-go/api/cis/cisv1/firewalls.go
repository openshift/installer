package cisv1

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/bluemix-go/client"
)

//FirewallRecord ...
type FirewallRecord struct {
	ID             string          `json:"id"`
	Description    string          `json:"description,omitempty"`
	Urls           []string        `json:"urls,omitempty"`
	Configurations []Configuration `json:"configurations,omitempty"`
	Paused         bool            `json:"paused,omitempty"`
	CreatedOn      *time.Time      `json:"created_on,omitempty"`
	ModifiedOn     *time.Time      `json:"modified_on,omitempty"`
	Mode           string          `json:"mode,omitempty"`
	Notes          string          `json:"notes,omitempty"`
	Configuration  *Configuration  `json:"configuration,omitempty"`
	Priority       int             `json:"priority,omitempty"`
}

//Configuration ...
type Configuration struct {
	Target string `json:"target,omitempty"`
	Value  string `json:"value,omitempty"`
}

//FirewallResults ...
type FirewallResults struct {
	FirewallList []FirewallRecord `json:"result"`
	ResultsInfo  ResultsCount     `json:"result_info"`
	Success      bool             `json:"success"`
	Errors       []Error          `json:"errors"`
}

//FirewallResult ...
type FirewallResult struct {
	Firewall FirewallRecord `json:"result"`
	Success  bool           `json:"success"`
	Errors   []Error        `json:"errors"`
	Messages []string       `json:"messages"`
}

//FirewallBody ...
type FirewallBody struct {
	Description    string          `json:"description,omitempty"`
	Urls           []string        `json:"urls,omitempty"`
	Configurations []Configuration `json:"configurations,omitempty"`
	Paused         bool            `json:"paused,omitempty"`
	Mode           string          `json:"mode,omitempty"`
	Notes          string          `json:"notes,omitempty"`
	Configuration  *Configuration  `json:"configuration,omitempty"`
	Priority       int             `json:"priority,omitempty"`
}

//Firewall ...
type Firewall interface {
	ListFirewall(cisID string, zoneID string, firewallType string) ([]FirewallRecord, error)
	GetFirewall(cisID string, zoneID string, firewallType string, firewallID string) (*FirewallRecord, error)
	CreateFirewall(cisID string, zoneID string, firewallType string, firewallBody FirewallBody) (*FirewallRecord, error)
	DeleteFirewall(cisID string, zoneID string, firewallType string, firewallID string) error
	UpdateFirewall(cisID string, zoneID string, firewallType string, firewallID string, firewallBody FirewallBody) (*FirewallRecord, error)
}

//firewall ...
type firewall struct {
	client *client.Client
}

func newFirewallAPI(c *client.Client) Firewall {
	return &firewall{
		client: c,
	}
}

func (r *firewall) ListFirewall(cisID string, zoneID string, firewallType string) ([]FirewallRecord, error) {
	firewallResults := FirewallResults{}

	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s", cisID, zoneID, firewallType)
	if firewallType == "access_rules" {
		rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/rules", cisID, zoneID, firewallType)
	}

	_, err := r.client.Get(rawURL, &firewallResults, nil)
	if err != nil {
		return nil, err
	}
	return firewallResults.FirewallList, err
}

func (r *firewall) GetFirewall(cisID string, zoneID string, firewallType string, firewallID string) (*FirewallRecord, error) {
	firewallResult := FirewallResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/%s", cisID, zoneID, firewallType, firewallID)
	if firewallType == "access_rules" {
		rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/rules/%s", cisID, zoneID, firewallType, firewallID)
	}
	_, err := r.client.Get(rawURL, &firewallResult, nil)
	if err != nil {
		return nil, err
	}
	return &firewallResult.Firewall, nil
}

func (r *firewall) DeleteFirewall(cisID string, zoneID string, firewallType string, firewallID string) error {
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/%s", cisID, zoneID, firewallType, firewallID)
	if firewallType == "access_rules" {
		rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/rules/%s", cisID, zoneID, firewallType, firewallID)
	}
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *firewall) CreateFirewall(cisID string, zoneID string, firewallType string, firewallBody FirewallBody) (*FirewallRecord, error) {
	firewallResult := FirewallResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s", cisID, zoneID, firewallType)
	if firewallType == "access_rules" {
		rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/rules", cisID, zoneID, firewallType)
	}
	log.Printf(">>>> rawURL : %s\n", rawURL)
	_, err := r.client.Post(rawURL, &firewallBody, &firewallResult)
	if err != nil {
		return nil, err
	}
	return &firewallResult.Firewall, nil
}

func (r *firewall) UpdateFirewall(cisID string, zoneID string, firewallType string, firewallID string, firewallBody FirewallBody) (*FirewallRecord, error) {
	firewallResult := FirewallResult{}
	var rawURL string
	rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/%s", cisID, zoneID, firewallType, firewallID)
	if firewallType == "access_rules" {
		rawURL = fmt.Sprintf("/v1/%s/zones/%s/firewall/%s/rules/%s", cisID, zoneID, firewallType, firewallID)
	}
	_, err := r.client.Put(rawURL, &firewallBody, &firewallResult)
	if err != nil {
		return nil, err
	}
	return &firewallResult.Firewall, nil
}
