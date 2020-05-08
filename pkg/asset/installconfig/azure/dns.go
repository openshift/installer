package azure

import (
	"context"
	"fmt"
	"time"

	azdns "github.com/Azure/azure-sdk-for-go/profiles/latest/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

//DNSConfig exposes functions to choose the DNS settings
type DNSConfig struct {
	Session *Session
}

//ZonesGetter fetches the DNS zones available for the installer
type ZonesGetter interface {
	GetAllPublicZones() (map[string]string, error)
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
		config.Session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

//GetPrivateDNSZoneID returns the Azure Private DNS zone resourceID
//by interpolating the subscriptionID, the resource group and the zone name
func (config DNSConfig) GetPrivateDNSZoneID(rgName string, zoneName string) string {
	return fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateDnsZones/%s",
		config.Session.Credentials.SubscriptionID,
		rgName,
		zoneName)
}

//GetDNSZone returns a DNS zone selected by survey
func (config DNSConfig) GetDNSZone() (*Zone, error) {
	//call azure api using the session to retrieve available base domain
	zonesClient := newZonesClient(config.Session)
	allZones, _ := zonesClient.GetAllPublicZones()
	if len(allZones) == 0 {
		return nil, errors.New("no public dns zone found in your subscription")
	}
	zoneNames := []string{}
	for zoneName := range allZones {
		zoneNames = append(zoneNames, zoneName)
	}

	var zoneName string
	err := survey.Ask([]*survey.Question{
		{
			Prompt: &survey.Select{
				Message: "Base Domain",
				Help:    "The base domain of the cluster. All DNS records will be sub-domains of this base and will also include the cluster name.\n\nIf you don't see you intended base-domain listed, create a new Azure DNS Zone and rerun the installer.",
				Options: zoneNames,
			},
		},
	}, &zoneName)
	if err != nil {
		return nil, err
	}

	return &Zone{
		ID:   allZones[zoneName],
		Name: zoneName,
	}, nil

}

//GetDNSRecordSet gets a record set for the zone identified by publicZoneID
func (config DNSConfig) GetDNSRecordSet(rgName string, zoneName string, relativeRecordSetName string, recordType azdns.RecordType) (*azdns.RecordSet, error) {
	recordsetsClient := newRecordSetsClient(config.Session)
	return recordsetsClient.GetRecordSet(rgName, zoneName, relativeRecordSetName, recordType)
}

//NewDNSConfig returns a new DNSConfig struct that helps configuring the DNS
//by querying your subscription and letting you choose
//which domain you wish to use for the cluster
func NewDNSConfig() (*DNSConfig, error) {
	session, err := GetSession()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve session information")
	}
	return &DNSConfig{Session: session}, nil
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

//GetAllPublicZones get all public zones from the current subscription
func (client *ZonesClient) GetAllPublicZones() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	allZones := map[string]string{}
	for zonesPage, err := client.azureClient.List(ctx, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}
		//TODO: filter out private zone and show only public zones.
		//the property is present in the REST api response, but not mapped yet in the stable SDK (present in preview)
		//https://github.com/Azure/azure-sdk-for-go/blob/07f918ba2d513bbc5b75bc4caac845e10f27449e/services/preview/dns/mgmt/2018-03-01-preview/dns/models.go#L857
		for _, zone := range zonesPage.Values() {
			allZones[to.String(zone.Name)] = to.String(zone.ID)
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
