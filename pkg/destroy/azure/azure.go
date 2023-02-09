package azure

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	azurestackdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	azurekiota "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
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
	SubscriptionID string
	TenantID       string
	Authorizer     autorest.Authorizer
	AuthProvider   *azurekiota.AzureIdentityAuthenticationProvider
	Environment    azureenv.Environment
	CloudName      azure.CloudEnvironment

	InfraID                     string
	ResourceGroupName           string
	BaseDomainResourceGroupName string

	Logger logrus.FieldLogger

	resourceGroupsClient    resources.GroupsClient
	zonesClient             dns.ZonesClient
	recordsClient           dns.RecordSetsClient
	privateRecordSetsClient privatedns.RecordSetsClient
	privateZonesClient      privatedns.PrivateZonesClient
	msgraphClient           *msgraphsdk.GraphServiceClient
}

func (o *ClusterUninstaller) configureClients() error {
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

	adapter, err := msgraphsdk.NewGraphRequestAdapter(o.AuthProvider)
	if err != nil {
		return err
	}
	// This can be empty for StackCloud
	if o.Environment.MicrosoftGraphEndpoint != "" {
		// Set the service root to the Microsoft Graph for the appropriate
		// cloud endpoint (e.g, GovCloud). Failing to do so results in an
		// unhelpful `context deadline exceeded` error.
		// NOTE: The API version must be included in the URL
		// See https://issues.redhat.com/browse/OCPBUGS-4549
		// See https://learn.microsoft.com/en-us/graph/sdks/national-clouds?tabs=go
		adapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", o.Environment.MicrosoftGraphEndpoint))
	}
	o.msgraphClient = msgraphsdk.NewGraphServiceClient(adapter)

	return nil
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
		SubscriptionID:              session.Credentials.SubscriptionID,
		TenantID:                    session.Credentials.TenantID,
		Authorizer:                  session.Authorizer,
		Environment:                 session.Environment,
		AuthProvider:                session.AuthProvider,
		InfraID:                     metadata.InfraID,
		ResourceGroupName:           group,
		Logger:                      logger,
		BaseDomainResourceGroupName: metadata.Azure.BaseDomainResourceGroupName,
		CloudName:                   cloudName,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	var errs []error
	var err error

	err = o.configureClients()
	if err != nil {
		return nil, err
	}

	// 2 hours
	timeout := 120 * time.Minute
	waitCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	wait.UntilWithContext(
		waitCtx,
		func(ctx context.Context) {
			o.Logger.Debugf("deleting public records")
			if o.CloudName == azure.StackCloud {
				err = deleteAzureStackPublicRecords(ctx, o)
			} else {
				err = deletePublicRecords(ctx, o.zonesClient, o.recordsClient, o.privateZonesClient, o.privateRecordSetsClient, o.Logger, o.ResourceGroupName)
			}
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
				} else if isResourceGroupBlockedError(err) {
					cancel()
					errs = append(errs, errors.Wrap(err, "unable to delete resource group, resources in the group are in use by others"))
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
			err = deleteApplicationRegistrations(ctx, o.msgraphClient, o.Logger, o.InfraID)
			if err != nil {
				oDataErr := extractODataError(err)
				o.Logger.Debug(oDataErr)
				if isAuthError(err) {
					cancel()
					errs = append(errs, errors.Wrap(oDataErr, "unable to authenticate when deleting application registrations and their service principals"))
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

	return nil, utilerrors.NewAggregate(errs)
}

func deleteAzureStackPublicRecords(ctx context.Context, o *ClusterUninstaller) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	logger := o.Logger
	rgName := o.BaseDomainResourceGroupName

	dnsClient := azurestackdns.NewZonesClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	dnsClient.Authorizer = o.Authorizer

	recordsClient := azurestackdns.NewRecordSetsClientWithBaseURI(o.Environment.ResourceManagerEndpoint, o.SubscriptionID)
	recordsClient.Authorizer = o.Authorizer

	var errs []error

	zonesPage, err := dnsClient.ListByResourceGroup(ctx, rgName, to.Int32Ptr(100))
	logger.Debug(err)
	if err != nil {
		if zonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			logger.Debug("already deleted the AzureStack zones")
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, errors.Wrap(err, "failed to list dns zone"))
		if isAuthError(err) {
			return err
		}
	}

	allZones := sets.NewString()
	for ; zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, errors.Wrap(err, "failed to advance to next dns zone"))
			continue
		}
		for _, zone := range zonesPage.Values() {
			allZones.Insert(to.String(zone.Name))
		}
	}

	clusterTag := fmt.Sprintf("kubernetes.io_cluster.%s", o.InfraID)
	for _, zone := range allZones.List() {
		for recordPages, err := recordsClient.ListByDNSZone(ctx, rgName, zone, to.Int32Ptr(100), ""); recordPages.NotDone(); err = recordPages.NextWithContext(ctx) {
			if err != nil {
				return err
			}
			for _, record := range recordPages.Values() {
				metadata := to.StringMap(record.Metadata)
				_, found := metadata[clusterTag]
				if found {
					resp, err := recordsClient.Delete(ctx, rgName, zone, to.String(record.Name), toAzureStackRecordType(to.String(record.Type)), "")
					if err != nil {
						if wasNotFound(resp.Response) {
							logger.WithField("record", to.String(record.Name)).Debug("already deleted")
							continue
						}
						return errors.Wrapf(err, "failed to delete record %s in zone %s", to.String(record.Name), zone)
					}
					logger.WithField("record", to.String(record.Name)).Info("deleted")
				}
			}
		}
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

func toAzureStackRecordType(t string) azurestackdns.RecordType {
	return azurestackdns.RecordType(strings.TrimPrefix(t, "Microsoft.Network/dnszones/"))
}

func deleteResourceGroup(ctx context.Context, client resources.GroupsClient, logger logrus.FieldLogger, name string) error {
	logger = logger.WithField("resource group", name)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	delFuture, err := client.Delete(ctx, name)
	if err == nil {
		err = delFuture.WaitForCompletionRef(ctx, client.Client)
	}
	if err != nil {
		if isNotFoundError(err) {
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

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var dErr autorest.DetailedError
	errors.As(err, &dErr)

	if dErr.StatusCode == http.StatusNotFound {
		return true
	}

	if dErr.StatusCode == 0 {
		serviceErr, ok := dErr.Original.(*azureenv.ServiceError)
		if ok {
			if strings.HasSuffix(serviceErr.Code, "NotFound") {
				return true
			}
		}
	}

	return false
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

	// https://github.com/Azure/azure-sdk-for-go/issues/16736
	// https://github.com/Azure/azure-sdk-for-go/blob/sdk/azidentity/v1.1.0/sdk/azidentity/errors.go#L36
	var authErr *azidentity.AuthenticationFailedError
	if errors.As(err, &authErr) {
		if authErr.RawResponse.StatusCode >= 400 && authErr.RawResponse.StatusCode <= 403 {
			return true
		}
	}

	return false
}

func isResourceGroupBlockedError(err error) bool {
	if err == nil {
		return false
	}

	var dErr autorest.DetailedError
	if errors.As(err, &dErr) {
		switch statusCode := dErr.StatusCode.(type) {
		case int:
			if statusCode == 409 {
				return true
			}
		}
	}

	return false
}

// Errors returned by the new Azure SDK are not very helpful. They just say
// "error status code received from the API". This function unwraps the
// ODataErr, if possible, and returns a new error with a more friendly "code:
// message" format.
func extractODataError(err error) error {
	var oDataErr *odataerrors.ODataError
	if errors.As(err, &oDataErr) {
		if typed := oDataErr.GetError(); typed != nil {
			return fmt.Errorf("%s: %s", *typed.GetCode(), *typed.GetMessage())
		}
	}
	return err
}

func deleteApplicationRegistrations(ctx context.Context, graphClient *msgraphsdk.GraphServiceClient, logger logrus.FieldLogger, infraID string) error {
	tag := fmt.Sprintf("kubernetes.io_cluster.%s=owned", infraID)
	servicePrincipals, err := getServicePrincipalsByTag(ctx, graphClient, tag, infraID)
	if err != nil {
		return errors.Wrap(err, "failed to gather list of Service Principals by tag")
	}
	// msgraphsdk can return a `nil` response even if no errors occurred
	if servicePrincipals == nil {
		logger.Debug("Empty response from API when listing Service Principals by tag")
		return nil
	}

	var errorList []error
	for _, sp := range servicePrincipals {
		appID := *sp.GetAppId()
		logger := logger.WithField("appID", appID)

		filter := fmt.Sprintf("appId eq '%s'", appID)
		listQuery := applications.ApplicationsRequestBuilderGetRequestConfiguration{
			QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}

		resp, err := graphClient.Applications().Get(ctx, &listQuery)
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		// msgraphsdk can return a `nil` response even if no errors occurred
		if resp == nil {
			logger.Debugf("Empty response getting Application from Service Principal %s", *sp.GetDisplayName())
			continue
		}
		apps := resp.GetValue()
		if len(apps) != 1 {
			err = fmt.Errorf("should have received only a single matching AppID, received %d instead", len(apps))
			errorList = append(errorList, err)
		}

		err = graphClient.ApplicationsById(*apps[0].GetId()).Delete(ctx, nil)
		if err != nil {
			errorList = append(errorList, err)
		}
		logger.Info("Deleted")
	}

	return utilerrors.NewAggregate(errorList)
}

func getServicePrincipalsByTag(ctx context.Context, graphClient *msgraphsdk.GraphServiceClient, matchTag, infraID string) ([]models.ServicePrincipalable, error) {
	filter := fmt.Sprintf("startswith(displayName, '%s') and tags/any(s:s eq '%s')", infraID, matchTag)
	listQuery := serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}
	resp, err := graphClient.ServicePrincipals().Get(ctx, &listQuery)
	if err != nil || resp == nil {
		return nil, err
	}
	return resp.GetValue(), nil
}
