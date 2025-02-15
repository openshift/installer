package azure

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	azurestackdns "github.com/Azure/azure-sdk-for-go/profiles/2018-03-01/dns/mgmt/dns"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azcoreto "github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcegraph/armresourcegraph"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/services/preview/dns/mgmt/2018-03-01-preview/dns"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/go-autorest/autorest"
	azureenv "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
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
	CloudName azure.CloudEnvironment
	Session   *azuresession.Session

	InfraID                     string
	ResourceGroupName           string
	BaseDomainResourceGroupName string
	NetworkResourceGroupName    string
	ZoneName                    string
	ClusterName                 string

	Logger logrus.FieldLogger

	resourceGroupsClient    resources.GroupsClient
	zonesClient             dns.ZonesClient
	recordsClient           dns.RecordSetsClient
	privateRecordSetsClient privatedns.RecordSetsClient
	privateZonesClient      privatedns.PrivateZonesClient
	msgraphClient           *msgraphsdk.GraphServiceClient
	resourceGraphClient     *armresourcegraph.Client
	tagsClient              *armresources.TagsClient
}

func (o *ClusterUninstaller) configureClients() error {
	subscriptionID := o.Session.Credentials.SubscriptionID
	endpoint := o.Session.Environment.ResourceManagerEndpoint

	o.resourceGroupsClient = resources.NewGroupsClientWithBaseURI(endpoint, subscriptionID)
	o.resourceGroupsClient.Authorizer = o.Session.Authorizer

	o.zonesClient = dns.NewZonesClientWithBaseURI(endpoint, subscriptionID)
	o.zonesClient.Authorizer = o.Session.Authorizer

	o.recordsClient = dns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	o.recordsClient.Authorizer = o.Session.Authorizer

	o.privateZonesClient = privatedns.NewPrivateZonesClientWithBaseURI(endpoint, subscriptionID)
	o.privateZonesClient.Authorizer = o.Session.Authorizer

	o.privateRecordSetsClient = privatedns.NewRecordSetsClientWithBaseURI(endpoint, subscriptionID)
	o.privateRecordSetsClient.Authorizer = o.Session.Authorizer

	adapter, err := msgraphsdk.NewGraphRequestAdapter(o.Session.AuthProvider)
	if err != nil {
		return err
	}
	// This can be empty for StackCloud
	if o.Session.Environment.MicrosoftGraphEndpoint != "" {
		// Set the service root to the Microsoft Graph for the appropriate
		// cloud endpoint (e.g, GovCloud). Failing to do so results in an
		// unhelpful `context deadline exceeded` error.
		// NOTE: The API version must be included in the URL
		// See https://issues.redhat.com/browse/OCPBUGS-4549
		// See https://learn.microsoft.com/en-us/graph/sdks/national-clouds?tabs=go
		adapter.SetBaseUrl(fmt.Sprintf("%s/v1.0", o.Session.Environment.MicrosoftGraphEndpoint))
	}
	o.msgraphClient = msgraphsdk.NewGraphServiceClient(adapter)

	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: o.Session.CloudConfig,
		},
	}

	rgClient, err := armresourcegraph.NewClient(o.Session.TokenCreds, clientOpts)
	if err != nil {
		return err
	}
	o.resourceGraphClient = rgClient

	tagsClient, err := armresources.NewTagsClient(o.Session.Credentials.SubscriptionID, o.Session.TokenCreds, clientOpts)
	if err != nil {
		return err
	}
	o.tagsClient = tagsClient

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

	return &ClusterUninstaller{
		Session:                     session,
		InfraID:                     metadata.InfraID,
		ResourceGroupName:           metadata.Azure.ResourceGroupName,
		Logger:                      logger,
		BaseDomainResourceGroupName: metadata.Azure.BaseDomainResourceGroupName,
		CloudName:                   cloudName,
		ZoneName:                    metadata.Azure.BaseDomainName,
		ClusterName:                 metadata.ClusterName,
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

	// Retrieve metadata from resource group tags, if available
	filter := fmt.Sprintf("tagName eq 'kubernetes.io_cluster.%s' and tagValue eq 'owned'", o.InfraID)
	groupPager, err := o.resourceGroupsClient.ListComplete(waitCtx, filter, to.Int32Ptr(1))
	if err != nil {
		return nil, fmt.Errorf("could not list resource groups: %w", err)
	}

	for ; groupPager.NotDone(); err = groupPager.NextWithContext(waitCtx) {
		if err != nil {
			o.Logger.Debugf("failed to advance to next resource group list page: %v", err)
			continue
		}
		group := groupPager.Value()
		if len(o.ResourceGroupName) == 0 {
			o.ResourceGroupName = to.String(group.Name)
			o.Logger.Debugf("found resource group name=%s from tags", o.ResourceGroupName)
		}
		if len(o.BaseDomainResourceGroupName) == 0 {
			o.BaseDomainResourceGroupName = to.String(group.Tags[azure.TagMetadataBaseDomainRG])
			o.Logger.Debugf("found base domain resource group name=%s from tags", o.BaseDomainResourceGroupName)
		}
		if len(o.NetworkResourceGroupName) == 0 {
			o.NetworkResourceGroupName = to.String(group.Tags[azure.TagMetadataNetworkRG])
			o.Logger.Debugf("found network resource group name=%s from tags", o.NetworkResourceGroupName)
		}
	}

	if len(o.ResourceGroupName) == 0 {
		o.ResourceGroupName = o.InfraID + "-rg"
		o.Logger.Debugf("using default resource group name=%s", o.ResourceGroupName)
	}

	err = wait.PollUntilContextCancel(
		waitCtx,
		1*time.Second,
		false,
		func(ctx context.Context) (bool, error) {
			o.Logger.Debugf("deleting public records")
			if o.CloudName == azure.StackCloud {
				err = deleteAzureStackPublicRecords(ctx, o)
			} else {
				err = deletePublicRecords(ctx, o)
			}
			if err != nil {
				o.Logger.Debug(err)
				if isAuthError(err) {
					errs = append(errs, fmt.Errorf("unable to authenticate when deleting public DNS records: %w", err))
					return true, err
				}
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to delete public DNS records: %w", err))
		o.Logger.Debug(err)
	}

	err = wait.PollUntilContextCancel(
		waitCtx,
		1*time.Second,
		false,
		func(ctx context.Context) (bool, error) {
			o.Logger.Debugf("deleting resource group")
			err = deleteResourceGroup(ctx, o.resourceGroupsClient, o.Logger, o.ResourceGroupName)
			if err != nil {
				o.Logger.Debug(err)
				if isAuthError(err) {
					errs = append(errs, fmt.Errorf("unable to authenticate when deleting resource group: %w", err))
					return true, err
				} else if isResourceGroupBlockedError(err) {
					errs = append(errs, fmt.Errorf("unable to delete resource group, resources in the group are in use by others: %w", err))
					return true, err
				}
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to delete resource group: %w", err))
		o.Logger.Debug(err)
	}

	err = wait.PollUntilContextCancel(
		waitCtx,
		1*time.Second,
		false,
		func(ctx context.Context) (bool, error) {
			o.Logger.Debugf("deleting application registrations")
			err = deleteApplicationRegistrations(ctx, o.msgraphClient, o.Logger, o.InfraID)
			if err != nil {
				oDataErr := extractODataError(err)
				o.Logger.Debug(oDataErr)
				if isAuthError(err) {
					errs = append(errs, fmt.Errorf("unable to authenticate when deleting application registrations and their service principals: %w", oDataErr))
					return true, err
				}
				return false, nil
			}
			return true, nil
		},
	)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to delete application registrations and their service principals: %w", err))
		o.Logger.Debug(err)
	}

	// do not attempt to remove shared tags on azure stack hub,
	// as the resource graph api is not supported there.
	if o.CloudName != azure.StackCloud {
		if err := removeSharedTags(
			waitCtx, o.resourceGraphClient, o.tagsClient, o.InfraID, o.Session.Credentials.SubscriptionID, o.Logger,
		); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove shared tags: %w", err))
			o.Logger.Debug(err)
		}
	}

	return nil, utilerrors.NewAggregate(errs)
}

func removeSharedTags(
	ctx context.Context,
	graphClient *armresourcegraph.Client,
	tagsClient *armresources.TagsClient,
	infraID, subscriptionID string,
	logger logrus.FieldLogger,
) error {
	tagKey := fmt.Sprintf("kubernetes.io_cluster.%s", infraID)
	query := fmt.Sprintf(
		"resources | where tags.['%s'] == 'shared' | project id, name, type",
		tagKey,
	)
	results, err := graphClient.Resources(ctx,
		armresourcegraph.QueryRequest{
			Query: &query,
			Subscriptions: []*string{
				&subscriptionID,
			},
			Options: &armresourcegraph.QueryRequestOptions{
				ResultFormat: azcoreto.Ptr(armresourcegraph.ResultFormatObjectArray),
			},
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to query resources with shared tag: %w", err)
	}

	tagsParam := armresources.TagsPatchResource{
		Operation: azcoreto.Ptr(armresources.TagsPatchOperationDelete),
		Properties: &armresources.Tags{
			Tags: map[string]*string{
				tagKey: to.StringPtr("shared"),
			},
		},
	}

	m, ok := results.Data.([]any)
	if !ok {
		logger.Debugf("could not cast results data (of type %T) to []any, skipping", results.Data)
		return nil
	}

	var errs []error
	for _, r := range m {
		items, ok := r.(map[string]any)
		if !ok {
			logger.Debugf("could not cast items (of type %T) to map[strin]any, skipping", items)
			continue
		}
		resourceName, ok := items["name"].(string)
		if !ok {
			logger.Debugf("could not cast resource name (of type %T) to string, skipping", items["name"])
			continue
		}
		resourceType, ok := items["type"].(string)
		if !ok {
			logger.Debugf("could not cast resource type (of type %T) to string, skipping", items["type"])
			continue
		}
		resourceID, ok := items["id"].(string)
		if !ok {
			logger.Debugf("could not cast resource id (of type %T) to string, skipping", items["id"])
			continue
		}
		logger := logger.WithFields(logrus.Fields{
			"resource": resourceName,
			"type":     resourceType,
		})
		logger.Debugf("removing shared tag from resource %q", resourceName)
		if _, err := tagsClient.UpdateAtScope(ctx, resourceID, tagsParam, nil); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove shared tag from %s: %w", resourceName, err))
		}
		logger.Infoln("removed shared tag")
	}
	return utilerrors.NewAggregate(errs)
}

func deleteAzureStackPublicRecords(ctx context.Context, o *ClusterUninstaller) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	logger := o.Logger
	rgName := o.BaseDomainResourceGroupName

	dnsClient := azurestackdns.NewZonesClientWithBaseURI(o.Session.Environment.ResourceManagerEndpoint, o.Session.Credentials.SubscriptionID)
	dnsClient.Authorizer = o.Session.Authorizer

	recordsClient := azurestackdns.NewRecordSetsClientWithBaseURI(o.Session.Environment.ResourceManagerEndpoint, o.Session.Credentials.SubscriptionID)
	recordsClient.Authorizer = o.Session.Authorizer

	var errs []error

	zonesPage, err := dnsClient.ListByResourceGroup(ctx, rgName, to.Int32Ptr(100))
	logger.Debug(err)
	if err != nil {
		if zonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			logger.Debug("already deleted the AzureStack zones")
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, fmt.Errorf("failed to list dns zone: %w", err))
		if isAuthError(err) {
			return err
		}
	}

	allZones := sets.NewString()
	for ; zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to advance to next dns zone: %w", err))
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
						return fmt.Errorf("failed to delete record %s in zone %s: %w", to.String(record.Name), zone, err)
					}
					logger.WithField("record", to.String(record.Name)).Info("deleted")
				}
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func deleteRecordsFromBaseDomain(ctx context.Context, o *ClusterUninstaller) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if o.BaseDomainResourceGroupName == "" || o.ZoneName == "" || o.ClusterName == "" {
		o.Logger.Debugf("could not find values in the metadata to get record set")
		return nil
	}

	var errs []error
	apiURL := fmt.Sprintf("api.%s", o.ClusterName)
	appsURL := fmt.Sprintf("*.apps.%s", o.ClusterName)
	errs = append(errs, deleteRecordsets(ctx, o, apiURL, dns.CNAME))
	errs = append(errs, deleteRecordsets(ctx, o, appsURL, dns.A))
	return utilerrors.NewAggregate(errs)
}

func deleteRecordsets(ctx context.Context, o *ClusterUninstaller, url string, recordType dns.RecordType) error {
	var errs []error
	tag := fmt.Sprintf("kubernetes.io_cluster.%s", o.InfraID)
	result, err := o.recordsClient.Get(ctx, o.BaseDomainResourceGroupName, o.ZoneName, url, recordType)
	if err != nil {
		logrus.Debugf("unable to find %s: already deleted or insufficient permissions", url)
		if isAuthError(err) {
			return err
		}
		return nil
	}

	if value, ok := result.Metadata[tag]; ok && *value == "owned" {
		deleteResult, err := o.recordsClient.Delete(ctx, o.BaseDomainResourceGroupName, o.ZoneName, url, recordType, "")
		if err != nil {
			if deleteResult.IsHTTPStatus(http.StatusNotFound) {
				o.Logger.Debug("already deleted")
				return utilerrors.NewAggregate(errs)
			}
			errs = append(errs, fmt.Errorf("failed to delete base domain dns zone: %w", err))
			if isAuthError(err) {
				return err
			}
		} else {
			o.Logger.WithField("record", url).Info("deleted")
		}
	} else {
		o.Logger.WithField("record", url).Debugf("metadata mismatch")
	}
	return utilerrors.NewAggregate(errs)
}

func deletePublicRecords(ctx context.Context, o *ClusterUninstaller) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// collect records from private zones in rgName
	var errs []error

	zonesPage, err := o.zonesClient.ListByResourceGroup(ctx, o.ResourceGroupName, to.Int32Ptr(100))
	if err != nil {
		if zonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			o.Logger.Debug("private zone not found, checking public records")
			err2 := deleteRecordsFromBaseDomain(ctx, o)
			if err2 != nil {
				o.Logger.Debugf("failed to delete record sets from the base domain: %w", err)
			}
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, fmt.Errorf("failed to list dns zone: %w", err))
		if isAuthError(err) {
			return err
		}
	}

	pageCount := 0
	for ; zonesPage.NotDone(); err = zonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to advance to next dns zone: %w", err))
			continue
		}
		pageCount++

		for _, zone := range zonesPage.Values() {
			if zone.ZoneType == dns.Private {
				if err := deletePublicRecordsForZone(ctx, o.zonesClient, o.recordsClient, o.Logger, o.ResourceGroupName, to.String(zone.Name)); err != nil {
					errs = append(errs, fmt.Errorf("failed to delete public records for %s: %w", to.String(zone.Name), err))
					if isAuthError(err) {
						return err
					}
					continue
				}
			}
		}
	}

	privateZonesPage, err := o.privateZonesClient.ListByResourceGroup(ctx, o.ResourceGroupName, to.Int32Ptr(100))
	if err != nil {
		if privateZonesPage.Response().IsHTTPStatus(http.StatusNotFound) {
			o.Logger.Debug("already deleted")
			return utilerrors.NewAggregate(errs)
		}
		errs = append(errs, fmt.Errorf("failed to list private dns zone: %w", err))
		if isAuthError(err) {
			return err
		}
	}

	for ; privateZonesPage.NotDone(); err = privateZonesPage.NextWithContext(ctx) {
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to advance to next dns zone: %w", err))
			continue
		}
		pageCount++

		for _, zone := range privateZonesPage.Values() {
			if err := deletePublicRecordsForPrivateZone(ctx, o.privateRecordSetsClient, o.zonesClient, o.recordsClient, o.Logger, o.ResourceGroupName, to.String(zone.Name)); err != nil {
				errs = append(errs, fmt.Errorf("failed to delete public records for %s: %w", to.String(zone.Name), err))
				if isAuthError(err) {
					return err
				}
				continue
			}
		}
	}

	if pageCount == 0 {
		o.Logger.Warn("no DNS records found: either they were already deleted or the service principal lacks permissions to list them")
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
		return fmt.Errorf("failed to find shared zone for %s: %w", zoneName, err)
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
						return fmt.Errorf("failed to delete record %s in zone %s: %w", to.String(record.Name), sharedZone.Name, err)
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
		return fmt.Errorf("failed to delete %s: %w", name, err)
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
	if errors.As(err, &dErr) {
		if dErr.StatusCode == http.StatusNotFound {
			return true
		}

		if dErr.StatusCode == 0 {
			var serviceErr *azureenv.ServiceError
			if errors.As(dErr.Original, &serviceErr) && strings.HasSuffix(serviceErr.Code, "NotFound") {
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
		return fmt.Errorf("failed to gather list of Service Principals by tag: %w", err)
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
