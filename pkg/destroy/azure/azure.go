package azure

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	azuresession "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	SubscriptionID  string
	TenantID        string
	GraphAuthorizer autorest.Authorizer
	Authorizer      autorest.Authorizer
	Environment     azureenv.Environment

	InfraID           string
	ResourceGroupName string

	Logger logrus.FieldLogger

	resourceGroupsClient    resources.GroupsClient
	zonesClient             dns.ZonesClient
	recordsClient           dns.RecordSetsClient
	privateRecordSetsClient privatedns.RecordSetsClient
	privateZonesClient      privatedns.PrivateZonesClient
	serviceprincipalsClient graphrbac.ServicePrincipalsClient
	applicationsClient      graphrbac.ApplicationsClient
}

func (o *ClusterUninstaller) configureClients() {
	o.resourceGroupsClient = resources.NewGroupsClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	o.resourceGroupsClient.Authorizer = o.Authorizer

	o.zonesClient = dns.NewZonesClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	o.zonesClient.Authorizer = o.Authorizer

	o.recordsClient = dns.NewRecordSetsClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	o.recordsClient.Authorizer = o.Authorizer

	o.privateZonesClient = privatedns.NewPrivateZonesClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	o.privateZonesClient.Authorizer = o.Authorizer

	o.privateRecordSetsClient = privatedns.NewRecordSetsClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	o.privateRecordSetsClient.Authorizer = o.Authorizer

	o.serviceprincipalsClient = graphrbac.NewServicePrincipalsClientWithBaseURI(o.Environment.GraphEndpoint, o.TenantID)
	o.serviceprincipalsClient.Authorizer = o.GraphAuthorizer

	o.applicationsClient = graphrbac.NewApplicationsClientWithBaseURI(o.Environment.GraphEndpoint, o.TenantID)
	o.applicationsClient.Authorizer = o.GraphAuthorizer
}

// New returns an Azure destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	cloudName := metadata.Azure.CloudName
	if cloudName == "" {
		cloudName = azure.PublicCloud
	}
	session, err := azuresession.GetSession(cloudName, metadata.Azure.ARMEndpoint)
	if err != nil {
		return nil, err
	}

	group := metadata.Azure.ResourceGroupName
	if len(group) == 0 {
		group = metadata.InfraID + "-rg"
	}

	return &ClusterUninstaller{
		SubscriptionID:    session.Credentials.SubscriptionID,
		TenantID:          session.Credentials.TenantID,
		GraphAuthorizer:   session.GraphAuthorizer,
		Authorizer:        session.Authorizer,
		Environment:       session.Environment,
		InfraID:           metadata.InfraID,
		ResourceGroupName: group,
		Logger:            logger,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() error {
	var errs []error
	var err error

	o.configureClients()

	// 2 hours
	timeout := 120 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	wait.UntilWithContext(
		waitCtx,
		func(ctx context.Context) {
			o.Logger.Debugf("deleting public records")
			err = deletePublicRecords(ctx, o.zonesClient, o.recordsClient, o.privateZonesClient, o.privateRecordSetsClient, o.Logger, o.ResourceGroupName)
			if err != nil {
				o.Logger.Debug(err)
				if isAuthError(err) {
					cancel()
					errs = append(errs, errors.Wrap(err, "unable to authenticate when deleting public DNS records"))
				}
				return
			}
			cancel()
		},
		1*time.Second,
	)
	err = waitCtx.Err()
	if err != nil && err != context.Canceled {
		errs = append(errs, errors.Wrap(err, "failed to delete public DNS records"))
		o.Logger.Debug(err)
	}

	deadline, _ := waitCtx.Deadline()
	diff := time.Until(deadline)
	if diff > 0 {
		waitCtx, cancel = context.WithTimeout(context.Background(), diff)
	}

	wait.UntilWithContext(
		waitCtx,
		func(ctx context.Context) {
			o.Logger.Debugf("deleting resource group")
			err = deleteResourceGroup(ctx, o.resourceGroupsClient, o.Logger, o.ResourceGroupName)
			if err != nil {
				o.Logger.Debug(err)
				if isAuthError(err) {
					cancel()
					errs = append(errs, errors.Wrap(err, "unable to authenticate when deleting resource group"))
				}
				return
			}
			cancel()
		},
		1*time.Second,
	)
	err = waitCtx.Err()
	if err != nil && err != context.Canceled {
		errs = append(errs, errors.Wrap(err, "failed to delete resource group"))
		o.Logger.Debug(err)
	}

	deadline, _ = waitCtx.Deadline()
	diff = time.Until(deadline)
	if diff > 0 {
		waitCtx, cancel = context.WithTimeout(context.Background(), diff)
	}

	wait.UntilWithContext(
		waitCtx,
		func(ctx context.Context) {
			o.Logger.Debugf("deleting application registrations")
			err = deleteApplicationRegistrations(ctx, o.applicationsClient, o.serviceprincipalsClient, o.Logger, o.InfraID)
			if err != nil {
				o.Logger.Debug(err)
				if isAuthError(err) {
					cancel()
					errs = append(errs, errors.Wrap(err, "unable to authenticate when deleting application registrations and their service principals"))
				}
				return
			}
			cancel()
		},
		1*time.Second,
	)
	err = waitCtx.Err()
	if err != nil && err != context.Canceled {
		errs = append(errs, errors.Wrap(err, "failed to delete application registrations and their service principals"))
		o.Logger.Debug(err)
	}

	return utilerrors.NewAggregate(errs)
}

func deletePublicRecords(ctx context.Context, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, privateDNSClient privatedns.PrivateZonesClient, privateRecordsClient privatedns.RecordSetsClient, logger logrus.FieldLogger, rgName string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// collect records from private zones in rgName
	var errs []error

	zonesPage, err := dnsClient.ListByResourceGroup(ctx, rgName, to.Int32Ptr(100))
	if err != nil {
		if zonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			logger.Debug("already deleted")
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, errors.Wrap(err, "failed to list dns zone"))
		if isAuthError(err) {
			return err
		}
	}

	for ; zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to advance to next dns zone"))
			continue
		}

		for _, zone := range zonesPage.Values() {
			if zone.ZoneType == dns.Private {
				if err := deletePublicRecordsForZone(ctx, dnsClient, recordsClient, logger, rgName, to.String(zone.Name)); err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to delete public records for %s", to.String(zone.Name)))
					if isAuthError(err) {
						return err
					}
					continue
				}
			}
		}
	}

	privateZonesPage, err := privateDNSClient.ListByResourceGroup(ctx, rgName, to.Int32Ptr(100))
	if err != nil {
		if privateZonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			logger.Debug("already deleted")
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, errors.Wrap(err, "failed to list private dns zone"))
		if isAuthError(err) {
			return err
		}
	}

	for ; privateZonesPage.NotDone(); err = privateZonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to advance to next dns zone"))
			continue
		}

		for _, zone := range privateZonesPage.Values() {
			if err := deletePublicRecordsForPrivateZone(ctx, privateRecordsClient, dnsClient, recordsClient, logger, rgName, to.String(zone.Name)); err != nil {
				errs = append(errs, errors.Wrapf(err, "failed to delete public records for %s", to.String(zone.Name)))
				if isAuthError(err) {
					return err
				}
				continue
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func deletePublicRecordsForZone(ctx context.Context, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, logger logrus.FieldLogger, zoneGroup, zoneName string) error {
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

	return deletePublicRecordsMatchingZoneName(ctx, dnsClient, recordsClient, logger, allPrivateRecords, zoneName)
}

func deletePublicRecordsForPrivateZone(ctx context.Context, privateRecordsClient privatedns.RecordSetsClient, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, logger logrus.FieldLogger, zoneGroup, zoneName string) error {
	// collect all the records from the zoneName
	allPrivateRecords := sets.NewString()
	for recordPages, err := privateRecordsClient.List(ctx, zoneGroup, zoneName, to.Int32Ptr(100), ""); recordPages.NotDone(); err = recordPages.NextWithContext(ctx) {
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

	return deletePublicRecordsMatchingZoneName(ctx, dnsClient, recordsClient, logger, allPrivateRecords, zoneName)
}

func deletePublicRecordsMatchingZoneName(ctx context.Context, dnsClient dns.ZonesClient, recordsClient dns.RecordSetsClient, logger logrus.FieldLogger, privateRecords sets.String, zoneName string) error {
	sharedZones, err := getSharedDNSZones(ctx, dnsClient, zoneName)
	if err != nil {
		return errors.Wrapf(err, "failed to find shared zone for %s", zoneName)
	}
	for _, sharedZone := range sharedZones {
		logger.Debugf("removing matching private records from %s", sharedZone.Name)
		for recordPages, err := recordsClient.ListByDNSZone(ctx, sharedZone.Group, sharedZone.Name, to.Int32Ptr(100), ""); recordPages.NotDone(); err = recordPages.NextWithContext(ctx) {
			if err != nil {
				return err
			}
			for _, record := range recordPages.Values() {
				if privateRecords.Has(fmt.Sprintf("%s.%s", to.String(record.Name), sharedZone.Name)) {
					resp, err := recordsClient.Delete(ctx, sharedZone.Group, sharedZone.Name, to.String(record.Name), toRecordType(to.String(record.Type)), "")
					if err != nil {
						if wasNotFound(resp.Response) {
							logger.WithField("record", to.String(record.Name)).Debug("already deleted")
							continue
						}
						return errors.Wrapf(err, "failed to delete record %s in zone %s", to.String(record.Name), sharedZone.Name)
					}
					logger.WithField("record", to.String(record.Name)).Info("deleted")
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

func deleteResourceGroup(ctx context.Context, client resources.GroupsClient, logger logrus.FieldLogger, name string) error {
	logger = logger.WithField("resource group", name)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	delFuture, err := client.Delete(ctx, name)
	if err != nil {
		if wasNotFound(delFuture.Response()) {
			logger.Debug("already deleted")
			return nil
		}
		return err
	}

	err = delFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if wasNotFound(delFuture.Response()) {
			logger.Debug("already deleted")
			return nil
		}
		return errors.Wrapf(err, "failed to delete %s", name)
	}
	logger.Info("deleted")
	return nil
}

func wasNotFound(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

func isAuthError(err error) bool {
	if err == nil {
		return false
	}

	var dErr autorest.DetailedError
	if errors.As(err, &dErr) {
		switch statusCode := dErr.StatusCode.(type) {
		case int:
			if statusCode >= 400 && statusCode <= 403 {
				return true
			}
		}
	}

	return false
}

func deleteApplicationRegistrations(ctx context.Context, appClient graphrbac.ApplicationsClient, spClient graphrbac.ServicePrincipalsClient, logger logrus.FieldLogger, infraID string) error {
	errorList := []error{}

	tag := fmt.Sprintf("kubernetes.io_cluster.%s=owned", infraID)
	servicePrincipals, err := getServicePrincipalsByTag(ctx, spClient, tag, infraID)
	if err != nil {
		if isAuthError(err) {
			return err
		}
		return errors.Wrap(err, "failed to gather list of Service Principals by tag")
	}

	for _, sp := range servicePrincipals {
		logger = logger.WithField("appID", *sp.AppID)
		appFilter := fmt.Sprintf("appId eq '%s'", *sp.AppID)
		appResults, err := appClient.List(ctx, appFilter)
		if err != nil {
			if isAuthError(err) {
				return err
			}
			errorList = append(errorList, err)
			continue
		}

		apps := appResults.Values()
		if len(apps) != 1 {
			msg := fmt.Sprintf("should have recieved only a single result matching AppID, received %d instead", len(apps))
			errorList = append(errorList, errors.New(msg))
			continue
		}

		_, err = appClient.Delete(ctx, *apps[0].ObjectID)
		if err != nil {
			if isAuthError(err) {
				return err
			}
			errorList = append(errorList, err)
			continue
		}
		logger.Info("deleted")
	}

	return utilerrors.NewAggregate(errorList)
}

func getServicePrincipalsByTag(ctx context.Context, spClient graphrbac.ServicePrincipalsClient, matchTag, infraID string) ([]graphrbac.ServicePrincipal, error) {
	matchedSPs := []graphrbac.ServicePrincipal{}

	infraFilter := fmt.Sprintf("startswith(displayName,'%s')", infraID)

	for spResults, err := spClient.List(ctx, infraFilter); spResults.NotDone(); err = spResults.NextWithContext(ctx) {
		if err != nil {
			return matchedSPs, err
		}

		for _, sp := range spResults.Values() {
			for _, tag := range *sp.Tags {
				if tag == matchTag {
					matchedSPs = append(matchedSPs, sp)
					break
				}
			}
		}
	}

	return matchedSPs, nil
}
