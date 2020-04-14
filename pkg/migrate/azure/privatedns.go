package azure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/sirupsen/logrus"

	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	azdns "github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	azprivatedns "github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	azconfig "github.com/openshift/installer/pkg/asset/installconfig/azure"
)

type legacyDNSZone struct {
	zone       *azdns.Zone
	recordsets []azdns.RecordSet
}

type legacyDNSClient struct {
	resourceGroup    string
	zonesClient      azdns.ZonesClient
	recordsetsClient azdns.RecordSetsClient
}

func newLegacyDNSClient(session *azconfig.Session, resourceGroup string) *legacyDNSClient {
	zonesClient := azdns.NewZonesClient(session.Credentials.SubscriptionID)
	zonesClient.Authorizer = session.Authorizer

	recordsetsClient := azdns.NewRecordSetsClient(session.Credentials.SubscriptionID)
	recordsetsClient.Authorizer = session.Authorizer

	return &legacyDNSClient{resourceGroup, zonesClient, recordsetsClient}
}

// Takes a subscription ID and parses the resource group out of it.
// A subscription ID has the format "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/xxxx-xxxxx-rg/providers/...".
// Splitting the string on '/' gives us the following slice:
// parts[0] = ''
// parts[1] = 'subscriptions'
// parts[2] = 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
// parts[3] = 'resourceGroups'
// parts[4] = 'xxxx-xxxxx-rg'
// parts[..] = ... the rest
// So if the length of the split is at least 5 and index 3 is "resourcegroups",
// we can safely assume the resource group is in the correct place.
func idToResourceGroup(id string) string {
	rg := ""
	parts := strings.Split(id, "/")
	if len(parts) >= 5 && strings.ToLower(parts[3]) == "resourcegroups" {
		rg = parts[4]
	}
	return rg
}

// Gets a single legacy zone and its recordsets
func (client *legacyDNSClient) getZone(legacyZone string) (*legacyDNSZone, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	zone, err := client.zonesClient.Get(ctx, client.resourceGroup, legacyZone)
	zone.Response.Response = nil
	if err != nil {
		return nil, err
	}

	if zone.ZoneProperties.ZoneType != azdns.Private {
		return nil, errors.New("not a private zone")
	}

	legacyDNSZone := legacyDNSZone{}
	legacyDNSZone.zone = &zone

	for recordsetsPage, err := client.recordsetsClient.ListAllByDNSZone(ctx, client.resourceGroup, to.String(zone.Name), to.Int32Ptr(100), ""); recordsetsPage.NotDone(); err = recordsetsPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}

		for _, rs := range recordsetsPage.Values() {
			legacyDNSZone.recordsets = append(legacyDNSZone.recordsets, rs)
		}
	}

	return &legacyDNSZone, nil
}

// Gets all legacy zones
func (client *legacyDNSClient) getZones() ([]azdns.Zone, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	var legacyDNSZones []azdns.Zone
	for zonesPage, err := client.zonesClient.List(ctx, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}

		for _, zone := range zonesPage.Values() {
			zone.Response.Response = nil
			if zone.ZoneProperties.ZoneType != azdns.Private {
				continue
			}

			legacyDNSZones = append(legacyDNSZones, zone)
		}
	}

	return legacyDNSZones, nil
}

type privateDNSClient struct {
	resourceGroup             string
	vnetResourceGroup         string
	virtualNetwork            string
	zonesClient               azprivatedns.PrivateZonesClient
	recordsetsClient          azprivatedns.RecordSetsClient
	virtualNetworkLinksClient azprivatedns.VirtualNetworkLinksClient
	virtualNetworksClient     aznetwork.VirtualNetworksClient
}

func newPrivateDNSClient(session *azconfig.Session, resourceGroup string, virtualNetwork string, vnetResourceGroup string) *privateDNSClient {
	zonesClient := azprivatedns.NewPrivateZonesClient(session.Credentials.SubscriptionID)
	zonesClient.Authorizer = session.Authorizer

	recordsetsClient := azprivatedns.NewRecordSetsClient(session.Credentials.SubscriptionID)
	recordsetsClient.Authorizer = session.Authorizer

	virtualNetworkLinksClient := azprivatedns.NewVirtualNetworkLinksClient(session.Credentials.SubscriptionID)
	virtualNetworkLinksClient.Authorizer = session.Authorizer

	virtualNetworksClient := aznetwork.NewVirtualNetworksClient(session.Credentials.SubscriptionID)
	virtualNetworksClient.Authorizer = session.Authorizer

	return &privateDNSClient{resourceGroup, vnetResourceGroup, virtualNetwork, zonesClient, recordsetsClient, virtualNetworkLinksClient, virtualNetworksClient}
}

// convert a legacy SOA record to a private SOA record
func legacySoaRecordToPrivate(legacySoaRecord *azdns.SoaRecord) *azprivatedns.SoaRecord {
	var soaRecord *azprivatedns.SoaRecord = nil

	if legacySoaRecord != nil {
		soaRecord = &azprivatedns.SoaRecord{
			Host:         legacySoaRecord.Host,
			Email:        legacySoaRecord.Email,
			SerialNumber: legacySoaRecord.SerialNumber,
			RefreshTime:  legacySoaRecord.RefreshTime,
			RetryTime:    legacySoaRecord.RetryTime,
			ExpireTime:   legacySoaRecord.ExpireTime,
			MinimumTTL:   legacySoaRecord.MinimumTTL,
		}
	}

	return soaRecord
}

// convert a legacy MX record to a private MX record
func legacyMxRecordToPrivate(legacyMxRecord *azdns.MxRecord) *azprivatedns.MxRecord {
	var mxRecord *azprivatedns.MxRecord = nil

	if legacyMxRecord != nil {
		mxRecord = &azprivatedns.MxRecord{
			Preference: legacyMxRecord.Preference,
			Exchange:   legacyMxRecord.Exchange,
		}
	}

	return mxRecord
}

// convert a legacy A record to a private A record
func legacyARecordToPrivate(legacyARecord *azdns.ARecord) *azprivatedns.ARecord {
	var aRecord *azprivatedns.ARecord = nil

	if legacyARecord != nil {
		aRecord = &azprivatedns.ARecord{
			Ipv4Address: legacyARecord.Ipv4Address,
		}
	}
	return aRecord
}

// convert a legacy AAAA record to a private AAAA record
func legacyAaaaRecordToPrivate(legacyAaaaRecord *azdns.AaaaRecord) *azprivatedns.AaaaRecord {
	var aaaaRecord *azprivatedns.AaaaRecord = nil

	if legacyAaaaRecord != nil {
		aaaaRecord = &azprivatedns.AaaaRecord{
			Ipv6Address: legacyAaaaRecord.Ipv6Address,
		}
	}

	return aaaaRecord
}

// convert a legacy CNAME record to a private CNAME record
func legacyCnameRecordToPrivate(legacyCnameRecord *azdns.CnameRecord) *azprivatedns.CnameRecord {
	var cnameRecord *azprivatedns.CnameRecord = nil

	if legacyCnameRecord != nil {
		cnameRecord = &azprivatedns.CnameRecord{
			Cname: legacyCnameRecord.Cname,
		}
	}

	return cnameRecord
}

// convert a legacy PTR record to a private PTR record
func legacyPtrRecordToPrivate(legacyPtrRecord *azdns.PtrRecord) *azprivatedns.PtrRecord {
	var ptrRecord *azprivatedns.PtrRecord = nil

	if legacyPtrRecord != nil {
		ptrRecord = &azprivatedns.PtrRecord{
			Ptrdname: legacyPtrRecord.Ptrdname,
		}
	}

	return ptrRecord
}

// convert a legacy SRV record to a private SRV record
func legacySrvRecordToPrivate(legacySrvRecord *azdns.SrvRecord) *azprivatedns.SrvRecord {
	var srvRecord *azprivatedns.SrvRecord = nil

	if legacySrvRecord != nil {
		srvRecord = &azprivatedns.SrvRecord{
			Priority: legacySrvRecord.Priority,
			Weight:   legacySrvRecord.Weight,
			Port:     legacySrvRecord.Port,
			Target:   legacySrvRecord.Target,
		}
	}

	return srvRecord
}

// convert a legacy TXT record to a private TXT record
func legacyTxtRecordToPrivate(legacyTxtRecord *azdns.TxtRecord) *azprivatedns.TxtRecord {
	var txtRecord *azprivatedns.TxtRecord = nil

	if legacyTxtRecord != nil {
		txtRecord = &azprivatedns.TxtRecord{
			Value: legacyTxtRecord.Value,
		}
	}

	return txtRecord
}

// convert an array of legacy MX records to an array of private MX records
func legacyMxRecordsToPrivate(legacyMxRecords *[]azdns.MxRecord) *[]azprivatedns.MxRecord {
	var mxRecords []azprivatedns.MxRecord = nil

	if legacyMxRecords != nil {
		mxRecords = make([]azprivatedns.MxRecord, len(*legacyMxRecords))
		for _, legacyMxRecord := range *legacyMxRecords {
			mxRecord := legacyMxRecordToPrivate(&legacyMxRecord)
			if mxRecord != nil {
				mxRecords = append(mxRecords, *mxRecord)
			}
		}
	}

	return &mxRecords
}

// convert an array of legacy A records to an array of private A records
func legacyARecordsToPrivate(legacyARecords *[]azdns.ARecord) *[]azprivatedns.ARecord {
	var aRecords []azprivatedns.ARecord = nil

	if legacyARecords != nil {
		for _, legacyARecord := range *legacyARecords {
			aRecord := legacyARecordToPrivate(&legacyARecord)
			if aRecord != nil {
				aRecords = append(aRecords, *aRecord)
			}

		}
	}

	return &aRecords
}

// convert an array of legacy AAAA records to an array of private AAAA records
func legacyAaaaRecordsToPrivate(legacyAaaaRecords *[]azdns.AaaaRecord) *[]azprivatedns.AaaaRecord {
	var aaaaRecords []azprivatedns.AaaaRecord = nil

	if legacyAaaaRecords != nil {
		for _, legacyAaaaRecord := range *legacyAaaaRecords {
			aaaaRecord := legacyAaaaRecordToPrivate(&legacyAaaaRecord)
			if aaaaRecord != nil {
				aaaaRecords = append(aaaaRecords, *aaaaRecord)
			}
		}
	}

	return &aaaaRecords
}

// convert an array of legacy PTR records to an array of private PTR records
func legacyPtrRecordsToPrivate(legacyPtrRecords *[]azdns.PtrRecord) *[]azprivatedns.PtrRecord {
	var ptrRecords []azprivatedns.PtrRecord = nil

	if legacyPtrRecords != nil {
		for _, legacyPtrRecord := range *legacyPtrRecords {
			ptrRecord := legacyPtrRecordToPrivate(&legacyPtrRecord)
			if ptrRecord != nil {
				ptrRecords = append(ptrRecords, *ptrRecord)
			}
		}
	}

	return &ptrRecords
}

// convert an array of legacy SRV records to an array of private SRV records
func legacySrvRecordsToPrivate(legacySrvRecords *[]azdns.SrvRecord) *[]azprivatedns.SrvRecord {
	var srvRecords []azprivatedns.SrvRecord = nil

	if legacySrvRecords != nil {
		for _, legacySrvRecord := range *legacySrvRecords {
			srvRecord := legacySrvRecordToPrivate(&legacySrvRecord)
			if srvRecord != nil {
				srvRecords = append(srvRecords, *srvRecord)
			}
		}
	}

	return &srvRecords
}

// convert an array of legacy TXT records to an array of private TXT records
func legacyTxtRecordsToPrivate(legacyTxtRecords *[]azdns.TxtRecord) *[]azprivatedns.TxtRecord {
	var txtRecords []azprivatedns.TxtRecord = nil

	if legacyTxtRecords != nil {
		for _, legacyTxtRecord := range *legacyTxtRecords {
			txtRecord := legacyTxtRecordToPrivate(&legacyTxtRecord)
			if txtRecord != nil {
				txtRecords = append(txtRecords, *txtRecord)
			}
		}
	}

	return &txtRecords
}

// Transforms a legacy zone to a private zone
func (client *privateDNSClient) migrateLegacyZone(legacyDNSZone *legacyDNSZone, link bool) error {
	legacyZone := legacyDNSZone.zone

	// Setup the private zone to create
	privateZone := azprivatedns.PrivateZone{}
	privateZone.Tags = legacyZone.Tags
	privateZone.Location = legacyZone.Location
	privateZone.Name = legacyZone.Name
	privateZone.PrivateZoneProperties = &azprivatedns.PrivateZoneProperties{
		MaxNumberOfRecordSets: legacyZone.ZoneProperties.MaxNumberOfRecordSets,
	}

	legacyRecordSets := legacyDNSZone.recordsets

	// Setup the associated recordsets to create
	privateRecordSets := []*azprivatedns.RecordSet{}
	for _, legacyRecordSet := range legacyRecordSets {
		recordType := strings.Replace(*legacyRecordSet.Type, "/dnszones/", "/privateDnsZones/", 1)

		// NS not supported in private zones
		if strings.TrimPrefix(recordType, "Microsoft.Network/privateDnsZones/") == string(azdns.NS) {
			continue
		}

		privateRecordSet := azprivatedns.RecordSet{
			Name: legacyRecordSet.Name,
			RecordSetProperties: &azprivatedns.RecordSetProperties{
				Metadata:    legacyRecordSet.Metadata,
				TTL:         legacyRecordSet.TTL,
				Fqdn:        legacyRecordSet.Fqdn,
				SoaRecord:   legacySoaRecordToPrivate(legacyRecordSet.SoaRecord),
				CnameRecord: legacyCnameRecordToPrivate(legacyRecordSet.CnameRecord),
				MxRecords:   legacyMxRecordsToPrivate(legacyRecordSet.MxRecords),
				ARecords:    legacyARecordsToPrivate(legacyRecordSet.ARecords),
				AaaaRecords: legacyAaaaRecordsToPrivate(legacyRecordSet.AaaaRecords),
				PtrRecords:  legacyPtrRecordsToPrivate(legacyRecordSet.PtrRecords),
				SrvRecords:  legacySrvRecordsToPrivate(legacyRecordSet.SrvRecords),
				TxtRecords:  legacyTxtRecordsToPrivate(legacyRecordSet.TxtRecords),
			},
			Type: &recordType,
		}
		privateRecordSets = append(privateRecordSets, &privateRecordSet)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 300*time.Second)
	defer cancel()

	// Create/Update the Zone
	logrus.Infof("zone: %s ... ", *privateZone.Name)
	zoneFuture, err := client.zonesClient.CreateOrUpdate(ctx, client.resourceGroup, *privateZone.Name, privateZone, "", "")
	if err != nil {
		return err
	}

	// Wait for zone creation to complete
	err = zoneFuture.WaitForCompletionRef(ctx, client.zonesClient.Client)
	if err != nil {
		return err
	}

	// Read back the newly created zone to verify creation
	_, err = client.zonesClient.Get(ctx, client.resourceGroup, *privateZone.Name)
	if err != nil {
		return err
	}
	logrus.Info("ok.")

	for _, recordSet := range privateRecordSets {
		recordType := azprivatedns.RecordType(strings.TrimPrefix(*recordSet.Type, "Microsoft.Network/privateDnsZones/"))
		relativeRecordSetName := *recordSet.Name
		recordSet.Type = nil

		// Create/Update the record
		logrus.Infof("record: %s %s ... ", recordType, relativeRecordSetName)
		_, err := client.recordsetsClient.CreateOrUpdate(ctx, client.resourceGroup, *privateZone.Name, recordType, relativeRecordSetName, *recordSet, "", "")
		if err != nil {
			return err
		}

		// Read back the newly created record to verify creation
		_, err = client.recordsetsClient.Get(ctx, client.resourceGroup, *privateZone.Name, recordType, relativeRecordSetName)
		if err != nil {
			return err
		}
		logrus.Info("ok.")
	}

	// Do we link, or not?
	if link == false || client.virtualNetwork == "" {
		return nil
	}

	// Get the virtual network so we have some parameters for the link creation
	virtualNetwork, err := client.virtualNetworksClient.Get(ctx, client.vnetResourceGroup, client.virtualNetwork, "")
	if err != nil {
		return err
	}

	virtualNetworkLinkName := fmt.Sprintf("%s-network-link", strings.Replace(client.vnetResourceGroup, "-rg", "", 1))

	virtualNetworkLink := azprivatedns.VirtualNetworkLink{
		Location: to.StringPtr("global"),
		VirtualNetworkLinkProperties: &azprivatedns.VirtualNetworkLinkProperties{
			VirtualNetwork: &azprivatedns.SubResource{
				ID: virtualNetwork.ID,
			},
			RegistrationEnabled: to.BoolPtr(false),
		},
	}

	// Create the virtual network link to DNS
	logrus.Infof("link: %s ... ", virtualNetworkLinkName)
	linkFuture, err := client.virtualNetworkLinksClient.CreateOrUpdate(ctx, client.resourceGroup, *privateZone.Name, virtualNetworkLinkName, virtualNetworkLink, "", "")
	if err != nil {
		return err
	}

	// Wait for the link creation to complete
	if err = linkFuture.WaitForCompletionRef(ctx, client.virtualNetworkLinksClient.Client); err != nil {
		return err
	}

	// Read back the newly created link to verify creation
	_, err = client.virtualNetworkLinksClient.Get(ctx, client.resourceGroup, *privateZone.Name, virtualNetworkLinkName)
	if err != nil {
		return err
	}
	logrus.Info("ok.")

	return nil
}

// Migrate does a migration from a legacy zone to a private zone
func Migrate(resourceGroup string, migrateZone string, virtualNetwork string, vnetResourceGroup string, link bool) error {
	session, err := azconfig.GetSession()
	if err != nil {
		return err
	}

	legacyDNSClient := newLegacyDNSClient(session, resourceGroup)
	privateDNSClient := newPrivateDNSClient(session, resourceGroup, virtualNetwork, vnetResourceGroup)

	legacyZone, err := legacyDNSClient.getZone(migrateZone)
	if err != nil {
		return err
	}

	// create new private zone
	err = privateDNSClient.migrateLegacyZone(legacyZone, link)
	if err != nil {
		return err
	}

	return nil
}

// Eligible shows legacy zones that are eligible for migrating to private zones
func Eligible() error {
	session, err := azconfig.GetSession()
	if err != nil {
		return err
	}

	legacyDNSClient := newLegacyDNSClient(session, "")

	zones, err := legacyDNSClient.getZones()
	if err != nil {
		return err
	}

	for _, zone := range zones {
		logrus.Infof("legacy zone=%s resourceGroup=%s", *zone.Name, idToResourceGroup(*zone.ID))
	}

	return nil
}
