package azure

import (
	"context"
	"errors"
	"fmt"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
)

// DNSConfig exposes functions to choose the DNS settings
type DNSConfig struct {
	session *Session
}

// ZonesGetter fetches the DNS zones available for the installer
type ZonesGetter interface {
	GetAllPublicZones() (map[string]string, error)
}

// ZonesClient wraps the azure ZonesClient internal
type ZonesClient struct {
	azureClient *armdns.ZonesClient
}

// RecordSetsClient wraps the azure RecordSetsClient internal
type RecordSetsClient struct {
	azureClient *armdns.RecordSetsClient
}

// Zone represents an Azure DNS Zone
type Zone struct {
	ID   string
	Name string
}

func (z Zone) String() string {
	return z.Name
}

// GetDNSZoneID returns the Azure DNS zone resourceID
// by interpolating the subscriptionID, the resource group and the zone name
func (config DNSConfig) GetDNSZoneID(rgName string, zoneName string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s",
		config.session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

// GetPrivateDNSZoneID returns the Azure Private DNS zone resourceID
// by interpolating the subscriptionID, the resource group and the zone name
func (config DNSConfig) GetPrivateDNSZoneID(rgName string, zoneName string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s",
		config.session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

// GetDNSZone returns a DNS zone selected by survey
func (config DNSConfig) GetDNSZone() (*Zone, error) {
	//call azure api using the session to retrieve available base domain
	zonesClient, err := newZonesClient(config.session)
	if err != nil {
		return nil, err
	}
	allZones, err := zonesClient.GetAllPublicZones()
	if err != nil {
		return nil, fmt.Errorf("failed to get public zones: %w", err)
	}
	if len(allZones) == 0 {
		return nil, errors.New("no public dns zone found in your subscription")
	}
	zoneNames := []string{}
	for zoneName := range allZones {
		zoneNames = append(zoneNames, zoneName)
	}

	var zoneName string
	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new Azure DNS Zone and rerun the installer.",
				Options: zoneNames,
			},
		},
	}, &zoneName)
	if err != nil {
		return nil, fmt.Errorf("failed UserInput: %w", err)
	}

	return &Zone{
		ID:   allZones[zoneName],
		Name: zoneName,
	}, nil

}

// GetDNSRecordSet gets a record set for the zone identified by publicZoneID
func (config DNSConfig) GetDNSRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType armdns.RecordType) (*armdns.RecordSet, error) {
	recordsetsClient, err := newRecordSetsClient(config.session)
	if err != nil {
		return nil, err
	}
	return recordsetsClient.GetRecordSet(rgName, zoneName, relativeRecordSetName, recordType)
}

// NewDNSConfig returns a new DNSConfig struct that helps configuring the DNS
// by querying your subscription and letting you choose
// which domain you wish to use for the cluster
func NewDNSConfig(ssn *Session) *DNSConfig {
	return &DNSConfig{session: ssn}
}

func newZonesClient(session *Session) (ZonesGetter, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}
	azureClient, err := armdns.NewZonesClient(session.Credentials.SubscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create zones client: %w", err)
	}
	return &ZonesClient{azureClient: azureClient}, nil
}

func newRecordSetsClient(session *Session) (*RecordSetsClient, error) {
	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: session.CloudConfig,
		},
	}
	azureClient, err := armdns.NewRecordSetsClient(session.Credentials.SubscriptionID, session.TokenCreds, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create record sets client: %w", err)
	}
	return &RecordSetsClient{azureClient: azureClient}, nil
}

// GetAllPublicZones get all public zones from the current subscription
func (client *ZonesClient) GetAllPublicZones() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	allZones := map[string]string{}
	pager := client.azureClient.NewListPager(&armdns.ZonesClientListOptions{Top: to.Ptr(int32(100))})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, zone := range page.Value {
			if zone.Name != nil && zone.ID != nil {
				allZones[*zone.Name] = *zone.ID
			}
		}
	}
	return allZones, nil
}

// GetRecordSet gets an Azure DNS recordset by zone, name and recordset type
func (client *RecordSetsClient) GetRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType armdns.RecordType) (*armdns.RecordSet, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	resp, err := client.azureClient.Get(ctx, rgName, zoneName, relativeRecordSetName, recordType, nil)
	if err != nil {
		return nil, err
	}

	return &resp.RecordSet, nil
}
