package azure

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	azuresession "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	SubscriptionID string
	Authorizer     autorest.Authorizer

	InfraID string

	log *logrus.Entry

	resourceGroupsClient resources.GroupsGroupClient
	zonesClient          dns.ZonesClient
	recordsClient        dns.RecordSetsClient
}

func (o *ClusterUninstaller) configureClients() {
	o.resourceGroupsClient = resources.NewGroupsGroupClient(o.SubscriptionID)
	o.resourceGroupsClient.Authorizer = o.Authorizer

	o.zonesClient = dns.NewZonesClient(o.SubscriptionID)
	o.zonesClient.Authorizer = o.Authorizer

	o.recordsClient = dns.NewRecordSetsClient(o.SubscriptionID)
	o.recordsClient.Authorizer = o.Authorizer
}

// New returns an Azure destroyer from ClusterMetadata.
func New(log *logrus.Entry, metadata *types.ClusterMetadata) (destroy.Destroyer, error) {
	session, err := azuresession.GetSession(log)
	if err != nil {
		return nil, err
	}

	return &ClusterUninstaller{
		SubscriptionID: session.Credentials.SubscriptionID,
		Authorizer:     session.Authorizer,
		InfraID:        metadata.InfraID,
		log:            log,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	o.configureClients()
	group := o.InfraID + "-rg"
	o.log.Debug("deleting public records")
	if err := deletePublicRecords(context.TODO(), o.log, o.zonesClient, o.recordsClient, group); err != nil {
		o.log.Debug(err)
		return errors.Wrap(err, "failed to delete public DNS records")
	}
	o.log.Debug("deleting resource group")
	if err := deleteResourceGroup(context.TODO(), o.log, o.resourceGroupsClient, group); err != nil {
		o.log.Debug(err)
		return errors.Wrap(err, "failed to delete resource group")
	}

	return nil
}

func deletePublicRecords(ctx context.Context, log *logrus.Entry, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, rgName string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// collect records from private zones in rgName
	var errs []error
	for zonesPage, err := dnsClient.ListByResourceGroup(ctx, rgName, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to list private zone"))
			continue
		}
		for _, zone := range zonesPage.Values() {
			if zone.ZoneType == dns.Private {
				if err := deletePublicRecordsForZone(ctx, log, dnsClient, recordsClient, rgName, to.String(zone.Name)); err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to delete public records for %s", to.String(zone.Name)))
					continue
				}
			}
		}
	}
	return utilerrors.NewAggregate(errs)
}

func deletePublicRecordsForZone(ctx context.Context, log *logrus.Entry, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, zoneGroup, zoneName string) error {
	// collect all the records from the zoneName
	allPrivateRecords := sets.NewString()
	for recordPages, err := recordsClient.ListByDNSZone(ctx, zoneGroup, zoneName, to.Int32Ptr(100), ""); recordPages.NotDone(); err = recordPages.NextWithContext(ctx) {
		if err != nil {
			return err
		}
		for _, record := range recordPages.Values() {
			if t := toRecordType(to.String(record.Type)); t == dns.SOA || t == dns.NS {
				continue
			}
			allPrivateRecords.Insert(fmt.Sprintf("%s.%s", to.String(record.Name), zoneName))
		}
	}

	sharedZones, err := getSharedDNSZones(ctx, dnsClient, zoneName)
	if err != nil {
		return errors.Wrapf(err, "failed to find shared zone for %s", zoneName)
	}
	for _, sharedZone := range sharedZones {
		log.Debugf("removing matching private records from %s", sharedZone.Name)
		for recordPages, err := recordsClient.ListByDNSZone(ctx, sharedZone.Group, sharedZone.Name, to.Int32Ptr(100), ""); recordPages.NotDone(); err = recordPages.NextWithContext(ctx) {
			if err != nil {
				return err
			}
			for _, record := range recordPages.Values() {
				if allPrivateRecords.Has(fmt.Sprintf("%s.%s", to.String(record.Name), sharedZone.Name)) {
					resp, err := recordsClient.Delete(ctx, sharedZone.Group, sharedZone.Name, to.String(record.Name), toRecordType(to.String(record.Type)), "")
					if err != nil {
						if wasNotFound(resp.Response) {
							log.WithField("record", to.String(record.Name)).Debug("already deleted")
							continue
						}
						return errors.Wrapf(err, "failed to delete record %s in zone %s", to.String(record.Name), sharedZone.Name)
					}
					log.WithField("record", to.String(record.Name)).Info("deleted")
				}
			}
		}
	}
	return nil
}

// getSharedDNSZones returns the all parent public dns zones for privZoneName in decreasing order of closeness.
func getSharedDNSZones(ctx context.Context, client dns.ZonesClient, privZoneName string) ([]dnsZone, error) {
	domain := privZoneName
	parents := sets.NewString(domain)
	for {
		idx := strings.Index(domain, ".")
		if idx == -1 {
			break
		}
		if len(domain[idx+1:]) > 0 {
			parents.Insert(domain[idx+1:])
		}
		domain = domain[idx+1:]
	}

	allPublicZones := []dnsZone{}
	for zonesPage, err := client.List(ctx, to.Int32Ptr(100)); zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			return nil, err
		}
		for _, zone := range zonesPage.Values() {
			if zone.ZoneType == dns.Public && parents.Has(to.String(zone.Name)) {
				allPublicZones = append(allPublicZones, dnsZone{Name: to.String(zone.Name), ID: to.String(zone.ID), Group: groupFromID(to.String(zone.ID)), Public: true})
				continue
			}
		}
	}
	sort.Slice(allPublicZones, func(i, j int) bool { return len(allPublicZones[i].Name) > len(allPublicZones[j].Name) })
	return allPublicZones, nil
}

type dnsZone struct {
	Name   string
	ID     string
	Group  string
	Public bool
}

func groupFromID(id string) string {
	return strings.Split(id, "/")[4]
}

func toRecordType(t string) dns.RecordType {
	return dns.RecordType(strings.TrimPrefix(t, "Microsoft.Network/dnszones/"))
}

func deleteResourceGroup(ctx context.Context, log *logrus.Entry, client resources.GroupsGroupClient, name string) error {
	log = log.WithField("resource group", name)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	delFuture, err := client.Delete(ctx, name)
	if err != nil {
		return err
	}

	err = delFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if wasNotFound(delFuture.Response()) {
			log.Debug("already deleted")
			return nil
		}
		return errors.Wrapf(err, "failed to delete %s", name)
	}
	log.Info("deleted")
	return nil
}

func wasNotFound(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}
