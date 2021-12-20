package cisv1

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/client"
)

type DnsRecord struct {
	Id         string      `json:"id"`
	Name       string      `json:"name,omitempty"`
	DnsType    string      `json:"type"`
	Content    string      `json:"content"`
	ZoneId     string      `json:"zone_id"`
	ZoneName   string      `json:"zone_name"`
	CreatedOn  *time.Time  `json:"created_on,omitempty"`
	ModifiedOn *time.Time  `json:"modified_on,omitempty"`
	Proxiable  bool        `json:"proxiable"`
	Proxied    bool        `json:"proxied"`
	Ttl        int         `json:"ttl"`
	Priority   int         `json:"priority,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

type DnsResults struct {
	DnsList     []DnsRecord  `json:"result"`
	ResultsInfo ResultsCount `json:"result_info"`
	Success     bool         `json:"success"`
	Errors      []Error      `json:"errors"`
}

type DnsResult struct {
	Dns      DnsRecord `json:"result"`
	Success  bool      `json:"success"`
	Errors   []Error   `json:"errors"`
	Messages []string  `json:"messages"`
}

type DnsBody struct {
	Name     string      `json:"name,omitempty"`
	DnsType  string      `json:"type"`
	Content  string      `json:"content,omitempty"`
	Priority int         `json:"priority,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Proxied  bool        `json:"proxied,omitempty"`
	Ttl      int         `json:"ttl,omitempty"`
}

type Dns interface {
	ListDns(cisId string, zoneId string) ([]DnsRecord, error)
	GetDns(cisId string, zoneId string, dnsId string) (*DnsRecord, error)
	CreateDns(cisId string, zoneId string, dnsBody DnsBody) (*DnsRecord, error)
	DeleteDns(cisId string, zoneId string, dnsId string) error
	UpdateDns(cisId string, zoneId string, dnsId string, dnsBody DnsBody) (*DnsRecord, error)
}

type dns struct {
	client *client.Client
}

func newDnsAPI(c *client.Client) Dns {
	return &dns{
		client: c,
	}
}

func (r *dns) ListDns(cisId string, zoneId string) ([]DnsRecord, error) {
	var records []DnsRecord
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/dns_records?page=1", cisId, zoneId)
	if _, err := r.client.GetPaginated(rawURL, NewDNSPaginatedResources(DnsRecord{}), func(resource interface{}) bool {
		if dns, ok := resource.(DnsRecord); ok {
			records = append(records, dns)
			return true
		}
		return false
	}); err != nil {
		return nil, fmt.Errorf("failed to list paginated dns records: %s", err)
	}
	return records, nil
}

func (r *dns) GetDns(cisId string, zoneId string, dnsId string) (*DnsRecord, error) {
	dnsResult := DnsResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/dns_records/%s", cisId, zoneId, dnsId)
	_, err := r.client.Get(rawURL, &dnsResult, nil)
	if err != nil {
		return nil, err
	}
	return &dnsResult.Dns, nil
}

func (r *dns) DeleteDns(cisId string, zoneId string, dnsId string) error {
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/dns_records/%s", cisId, zoneId, dnsId)
	_, err := r.client.Delete(rawURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *dns) CreateDns(cisId string, zoneId string, dnsBody DnsBody) (*DnsRecord, error) {
	dnsResult := DnsResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/dns_records", cisId, zoneId)
	_, err := r.client.Post(rawURL, &dnsBody, &dnsResult)
	if err != nil {
		return nil, err
	}
	return &dnsResult.Dns, nil
}

func (r *dns) UpdateDns(cisId string, zoneId string, dnsId string, dnsBody DnsBody) (*DnsRecord, error) {
	dnsResult := DnsResult{}
	rawURL := fmt.Sprintf("/v1/%s/zones/%s/dns_records/%s", cisId, zoneId, dnsId)
	_, err := r.client.Put(rawURL, &dnsBody, &dnsResult)
	if err != nil {
		return nil, err
	}
	return &dnsResult.Dns, nil
}
