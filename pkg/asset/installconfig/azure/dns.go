package azure

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	azdns "github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

//DNSConfig exposes functions to choose the DNS settings
type DNSConfig struct {
	session *Session
}

//ZonesGetter fetches the DNS zones available for the installer
type ZonesGetter interface {
	GetAllPublicZones() (map[string][]string, error)
}

//ZonesClient wraps the azure ZonesClient internal
type ZonesClient struct {
	azureClient azdns.ZonesClient
}

//RecordSetsClient wraps the azure RecordSetsClient internal
type RecordSetsClient struct {
	azureClient azdns.RecordSetsClient
}

//Zone represents an Azure DNS Zone
type Zone struct {
	ID   string
	Name string
}

func (z Zone) String() string {
	return fmt.Sprintf("%s", z.Name)
}

func transformZone(f func(s string) *Zone) survey.Transformer {
	return func(ans interface{}) interface{} {
		// if the answer value passed in is the zero value of the appropriate type
		if "" == ans.(string) {
			return nil
		}

		s, ok := ans.(string)
		if !ok {
			return nil
		}

		return f(s)
	}
}

//GetDNSZoneID returns the Azure DNS zone resourceID
//by interpolating the subscriptionID, the resource group and the zone name
func (config DNSConfig) GetDNSZoneID(rgName string, zoneName string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnszones/%s",
		config.session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

//GetPrivateDNSZoneID returns the Azure Private DNS zone resourceID
//by interpolating the subscriptionID, the resource group and the zone name
func (config DNSConfig) GetPrivateDNSZoneID(rgName string, zoneName string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s",
		config.session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

//GetDNSZone returns a DNS zone selected by survey
func (config DNSConfig) GetDNSZone() (*Zone, error) {
	//call azure api using the session to retrieve available base domain
	zonesClient := newZonesClient(config.Session)
	publicZonesMap, err := zonesClient.GetAllPublicZones()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve base domains")
	}

	if len(publicZonesMap) == 0 {
		return nil, errors.New("no public dns zone found in your subscription")
	}

	publicZones := make([]string, 0, len(publicZonesMap))
	for name, ids := range publicZonesMap {
		for _, id := range ids {
			// A subscription ID has the format
			// "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/xxxx-xxxxx-rg/providers/...".
			// Splitting the string on '/' gives us the following slice:
			// parts[0] = ''
			// parts[1] = 'subscriptions'
			// parts[2] = 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
			// parts[3] = 'resourceGroups'
			// parts[4] = 'xxxx-xxxxx-rg' <- This is the resource group name
			// parts[..] = ... the rest
			parts := strings.Split(id, "/")
			rgName := parts[4]
			publicZones = append(publicZones, fmt.Sprintf("%s (%s)", name, rgName))
		}
	}
	sort.Strings(publicZones)

	var publicZoneNameChoice string
	var publicZoneName string
	var publicZoneID string

	err = survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new Azure DNS Zone and rerun the installer.",
				Options: publicZones,
			},
		},
	}, &publicZoneNameChoice)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(publicZoneNameChoice, " ")
	publicZoneName = parts[0]
	rgName := parts[1][1 : len(parts[1])-1]
	publicZoneID = config.GetDNSZoneID(rgName, publicZoneName)

	return &Zone{
		ID:   publicZoneID,
		Name: publicZoneName,
	}, nil

}

//GetDNSRecordSet gets a record set for the zone identified by publicZoneID
func (config DNSConfig) GetDNSRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType azdns.RecordType) (*azdns.RecordSet, error) {
	recordsetsClient := newRecordSetsClient(config.session)
	return recordsetsClient.GetRecordSet(rgName, zoneName, relativeRecordSetName, recordType)
}

//NewDNSConfig returns a new DNSConfig struct that helps configuring the DNS
//by querying your subscription and letting you choose
//which domain you wish to use for the cluster
func NewDNSConfig(ssn *Session) *DNSConfig {
	return &DNSConfig{session: ssn}
}

func newZonesClient(session *Session) ZonesGetter {
	azureClient := azdns.NewZonesClient(session.Credentials.SubscriptionID)
	azureClient.Authorizer = session.Authorizer
	return &ZonesClient{azureClient: azureClient}
}

func newRecordSetsClient(session *Session) *RecordSetsClient {
	azureClient := azdns.NewRecordSetsClient(session.Credentials.SubscriptionID)
	azureClient.Authorizer = session.Authorizer
	return &RecordSetsClient{azureClient: azureClient}
}

// GetAllPublicZones returns a map of DNS names to zone ID's of all public domains
// in the current subscription.
func (client *ZonesClient) GetAllPublicZones() (map[string][]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	allZones := map[string][]string{}
	for zonesPage, err := client.azureClient.List(ctx, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}
		//TODO: filter out private zone and show only public zones.
		//the property is present in the REST api response, but not mapped yet in the stable SDK (present in preview)
		//https://github.com/Azure/azure-sdk-for-go/blob/07f918ba2d513bbc5b75bc4caac845e10f27449e/services/preview/dns/mgmt/2018-03-01-preview/dns/models.go#L857
		for _, zone := range zonesPage.Values() {
			zoneName := to.String(zone.Name)
			allZones[zoneName] = append(allZones[zoneName], to.String(zone.ID))
		}
	}
	return allZones, nil
}

//GetRecordSet gets an Azure DNS recordset by zone, name and recordset type
func (client *RecordSetsClient) GetRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType azdns.RecordType) (*azdns.RecordSet, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	recordset, err := client.azureClient.Get(ctx, rgName, zoneName, relativeRecordSetName, recordType)
	if err != nil {
		return nil, err
	}

	return &recordset, nil
}
