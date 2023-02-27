package ibmcloud

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnssvcsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	icibmcloud "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

var (
	defaultTimeout = 2 * time.Minute
)

// ClusterUninstaller holds the various options for the cluster we want to delete
type ClusterUninstaller struct {
	ClusterName         string
	Context             context.Context
	Logger              logrus.FieldLogger
	InfraID             string
	AccountID           string
	BaseDomain          string
	CISInstanceCRN      string
	DNSInstanceID       string
	Region              string
	ResourceGroupName   string
	UserProvidedSubnets []string
	UserProvidedVPC     string

	managementSvc          *resourcemanagerv2.ResourceManagerV2
	controllerSvc          *resourcecontrollerv2.ResourceControllerV2
	vpcSvc                 *vpcv1.VpcV1
	iamPolicyManagementSvc *iampolicymanagementv1.IamPolicyManagementV1
	zonesSvc               *zonesv1.ZonesV1
	dnsZonesSvc            *dnszonesv1.DnsZonesV1
	dnsServicesSvc         *dnssvcsv1.DnsSvcsV1
	dnsRecordsSvc          *dnsrecordsv1.DnsRecordsV1
	maxRetryAttempt        int

	resourceGroupID string
	cosInstanceID   string
	zoneID          string

	errorTracker
	pendingItemTracker
}

// New returns an IBMCloud destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	return &ClusterUninstaller{
		ClusterName:         metadata.ClusterName,
		Context:             context.Background(),
		Logger:              logger,
		InfraID:             metadata.InfraID,
		AccountID:           metadata.ClusterPlatformMetadata.IBMCloud.AccountID,
		BaseDomain:          metadata.ClusterPlatformMetadata.IBMCloud.BaseDomain,
		CISInstanceCRN:      metadata.ClusterPlatformMetadata.IBMCloud.CISInstanceCRN,
		DNSInstanceID:       metadata.ClusterPlatformMetadata.IBMCloud.DNSInstanceID,
		Region:              metadata.ClusterPlatformMetadata.IBMCloud.Region,
		ResourceGroupName:   metadata.ClusterPlatformMetadata.IBMCloud.ResourceGroupName,
		UserProvidedSubnets: metadata.ClusterPlatformMetadata.IBMCloud.Subnets,
		UserProvidedVPC:     metadata.ClusterPlatformMetadata.IBMCloud.VPC,
		pendingItemTracker:  newPendingItemTracker(),
		maxRetryAttempt:     30,
	}, nil
}

// Retry ...
func (o *ClusterUninstaller) Retry(funcToRetry func() (error, bool)) error {
	var err error
	var stopRetry bool
	retryGap := 10
	for i := 0; i < o.maxRetryAttempt; i++ {
		if i > 0 {
			time.Sleep(time.Duration(retryGap) * time.Second)
		}
		// Call function which required retry, retry is decided by function itself
		err, stopRetry = funcToRetry()
		if stopRetry {
			break
		}

		if (i + 1) < o.maxRetryAttempt {
			o.Logger.Infof("UNEXPECTED RESULT, Re-attempting execution .., attempt=%d, retry-gap=%d, max-retry-Attempts=%d, stopRetry=%t, error=%v", i+1,
				retryGap, o.maxRetryAttempt, stopRetry, err)
		}
	}
	return err
}

// Run is the entrypoint to start the uninstall process
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	err := o.loadSDKServices()
	if err != nil {
		return nil, err
	}

	err = o.destroyCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failed to destroy cluster")
	}

	return nil, nil
}

func (o *ClusterUninstaller) destroyCluster() error {
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Stop instances", execute: o.stopInstances},
	}, {
		// Instances must occur before LB cleanup
		{name: "Instances", execute: o.destroyInstances},
		{name: "Disks", execute: o.destroyDisks},
	}, {
		// LB's must occur before Subnet cleanup
		{name: "Load Balancers", execute: o.destroyLoadBalancers},
	}, {
		{name: "Subnets", execute: o.destroySubnets},
	}, {
		// Public Gateways must occur before FIP's cleanup
		// Security Groups must occur before VPC cleanup
		{name: "Images", execute: o.destroyImages},
		{name: "Public Gateways", execute: o.destroyPublicGateways},
		{name: "Security Groups", execute: o.destroySecurityGroups},
	}, {
		{name: "Floating IPs", execute: o.destroyFloatingIPs},
	}, {
		{name: "Dedicated Hosts", execute: o.destroyDedicatedHosts},
		{name: "VPCs", execute: o.destroyVPCs},
	}, {
		// IAM must occur before COS cleanup
		{name: "IAM Authorizations", execute: o.destroyIAMAuthorizations},
	}, {
		// COS must occur before RG cleanup
		{name: "Cloud Object Storage Instances", execute: o.destroyCOSInstances},
		{name: "Dedicated Host Groups", execute: o.destroyDedicatedHostGroups},
	}, {
		{name: "DNS Records", execute: o.destroyDNSRecords},
		{name: "Resource Groups", execute: o.destroyResourceGroups},
	}}

	for _, stage := range stagedFuncs {
		var wg sync.WaitGroup
		errCh := make(chan error)
		wgDone := make(chan bool)

		for _, f := range stage {
			wg.Add(1)
			go o.executeStageFunction(f, errCh, &wg)
		}

		go func() {
			wg.Wait()
			close(wgDone)
		}()

		select {
		case <-wgDone:
			// On to the next stage
			continue
		case err := <-errCh:
			return err
		}
	}

	return nil
}

func (o *ClusterUninstaller) executeStageFunction(f struct {
	name    string
	execute func() error
}, errCh chan error, wg *sync.WaitGroup) error {
	defer wg.Done()

	err := wait.PollImmediateInfinite(
		time.Second*10,
		func() (bool, error) {
			ferr := f.execute()
			if ferr != nil {
				o.Logger.Debugf("%s: %v", f.name, ferr)
				return false, nil
			}
			return true, nil
		},
	)

	if err != nil {
		errCh <- err
	}
	return nil
}

func (o *ClusterUninstaller) loadSDKServices() error {
	apiKey := os.Getenv("IC_API_KEY")

	userAgentString := fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)

	// ResourceManagerV2
	rmAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
	if err != nil {
		return err
	}
	o.managementSvc, err = resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: rmAuthenticator,
	})
	if err != nil {
		return err
	}
	o.managementSvc.Service.SetUserAgent(userAgentString)

	// Attempt to retrieve the ResourceGroupID as soon as possible to validate ResourceGroupName
	_, err = o.ResourceGroupID()
	if err != nil {
		return err
	}

	// ResourceControllerV2
	rcAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
	if err != nil {
		return err
	}
	o.controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: rcAuthenticator,
	})
	if err != nil {
		return err
	}
	o.controllerSvc.Service.SetUserAgent(userAgentString)

	// IamPolicyManagementV1
	ipmAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
	if err != nil {
		return err
	}
	o.iamPolicyManagementSvc, err = iampolicymanagementv1.NewIamPolicyManagementV1(&iampolicymanagementv1.IamPolicyManagementV1Options{
		Authenticator: ipmAuthenticator,
	})
	if err != nil {
		return err
	}
	o.iamPolicyManagementSvc.Service.SetUserAgent(userAgentString)

	if len(o.CISInstanceCRN) > 0 {
		// ZonesV1
		zAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
		if err != nil {
			return err
		}
		o.zonesSvc, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: zAuthenticator,
			Crn:           core.StringPtr(o.CISInstanceCRN),
		})
		if err != nil {
			return err
		}
		o.zonesSvc.Service.SetUserAgent(userAgentString)

		// Get the Zone ID
		options := o.zonesSvc.NewListZonesOptions()
		resources, _, err := o.zonesSvc.ListZonesWithContext(o.Context, options)
		if err != nil {
			return err
		}

		zoneID := ""
		for _, zone := range resources.Result {
			if strings.Contains(o.BaseDomain, *zone.Name) {
				zoneID = *zone.ID
			}
		}
		if zoneID == "" {
			return errors.Errorf("Could not determine CIS DNS zone ID from base domain %q", o.BaseDomain)
		}

		// DnsRecordsV1
		dnsAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
		if err != nil {
			return err
		}
		o.dnsRecordsSvc, err = dnsrecordsv1.NewDnsRecordsV1(&dnsrecordsv1.DnsRecordsV1Options{
			Authenticator:  dnsAuthenticator,
			Crn:            core.StringPtr(o.CISInstanceCRN),
			ZoneIdentifier: core.StringPtr(zoneID),
		})
		if err != nil {
			return err
		}
		o.dnsRecordsSvc.Service.SetUserAgent(userAgentString)
	} else if len(o.DNSInstanceID) > 0 {
		// DnsSvcsV1
		dnsAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
		if err != nil {
			return err
		}
		o.dnsServicesSvc, err = dnssvcsv1.NewDnsSvcsV1(&dnssvcsv1.DnsSvcsV1Options{
			Authenticator: dnsAuthenticator,
		})
		if err != nil {
			return err
		}
		o.dnsServicesSvc.Service.SetUserAgent(userAgentString)

		// Get the Zone ID
		dzOptions := o.dnsServicesSvc.NewListDnszonesOptions(o.DNSInstanceID)
		dzResult, _, err := o.dnsServicesSvc.ListDnszonesWithContext(o.Context, dzOptions)
		if err != nil {
			return err
		}

		zoneID := ""
		for _, zone := range dzResult.Dnszones {
			if o.BaseDomain == *zone.Name {
				zoneID = *zone.ID
				break
			}
		}
		if zoneID == "" {
			return errors.Errorf("Could not determine DNS Services DNS zone ID from base domain %q", o.BaseDomain)
		}
		o.Logger.Debugf("Found DNS Services DNS zone ID for base domain %q: %s", o.BaseDomain, zoneID)
		o.zoneID = zoneID
	}

	// VpcV1
	vpcAuthenticator, err := icibmcloud.NewIamAuthenticator(apiKey)
	if err != nil {
		return err
	}
	o.vpcSvc, err = vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: vpcAuthenticator,
	})
	if err != nil {
		return err
	}
	o.vpcSvc.Service.SetUserAgent(userAgentString)

	region, _, err := o.vpcSvc.GetRegion(o.vpcSvc.NewGetRegionOptions(o.Region))
	if err != nil {
		return err
	}

	err = o.vpcSvc.SetServiceURL(fmt.Sprintf("%s/v1", *region.Endpoint))
	if err != nil {
		return err
	}

	return nil
}

func (o *ClusterUninstaller) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(o.Context, defaultTimeout)
}

// ResourceGroupID returns the ID of the resource group using its name
func (o *ClusterUninstaller) ResourceGroupID() (string, error) {
	if o.resourceGroupID != "" {
		return o.resourceGroupID, nil
	}

	// If no ResourceGroupName is available, raise an error
	if o.ResourceGroupName == "" {
		return "", errors.Errorf("No ResourceGroupName provided")
	}

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.managementSvc.NewListResourceGroupsOptions()
	options.SetAccountID(o.AccountID)
	options.SetName(o.ResourceGroupName)
	resources, _, err := o.managementSvc.ListResourceGroupsWithContext(ctx, options)
	if err != nil {
		return "", err
	}
	if len(resources.Resources) == 0 {
		return "", errors.Errorf("ResourceGroup '%q' not found", o.ResourceGroupName)
	} else if len(resources.Resources) > 1 {
		return "", errors.Errorf("Too many resource groups matched name %q", o.ResourceGroupName)
	}

	o.SetResourceGroupID(*resources.Resources[0].ID)
	return o.resourceGroupID, nil
}

// SetResourceGroupID sets the resource group ID
func (o *ClusterUninstaller) SetResourceGroupID(id string) {
	o.resourceGroupID = id
}

type ibmError struct {
	Status  int
	Message string
}

func isNoOp(err *ibmError) bool {
	if err == nil {
		return false
	}

	return err.Status == http.StatusNotFound
}

// aggregateError is a utility function that takes a slice of errors and an
// optional pending argument, and returns an error or nil
func aggregateError(errs []error, pending ...int) error {
	err := utilerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	if len(pending) > 0 && pending[0] > 0 {
		return errors.Errorf("%d items pending", pending[0])
	}
	return nil
}

// pendingItemTracker tracks a set of pending item names for a given type of resource
type pendingItemTracker struct {
	pendingItems map[string]cloudResources
}

func newPendingItemTracker() pendingItemTracker {
	return pendingItemTracker{
		pendingItems: map[string]cloudResources{},
	}
}

// GetAllPendintItems returns a slice of all of the pending items across all types.
func (t pendingItemTracker) GetAllPendingItems() []cloudResource {
	var items []cloudResource
	for _, is := range t.pendingItems {
		for _, i := range is {
			items = append(items, i)
		}
	}
	return items
}

// getPendingItems returns the list of resources to be deleted.
func (t pendingItemTracker) getPendingItems(itemType string) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	return lastFound.list()
}

// insertPendingItems adds to the list of resources to be deleted.
func (t pendingItemTracker) insertPendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	lastFound = lastFound.insert(items...)
	t.pendingItems[itemType] = lastFound
	return lastFound.list()
}

// deletePendingItems removes from the list of resources to be deleted.
func (t pendingItemTracker) deletePendingItems(itemType string, items []cloudResource) []cloudResource {
	lastFound, exists := t.pendingItems[itemType]
	if !exists {
		lastFound = cloudResources{}
	}
	lastFound = lastFound.delete(items...)
	t.pendingItems[itemType] = lastFound
	return lastFound.list()
}

func isErrorStatus(code int64) bool {
	return code != 0 && (code < 200 || code >= 300)
}

func (o *ClusterUninstaller) clusterLabelFilter() string {
	return fmt.Sprintf("kubernetes-io-cluster-%s:owned", o.InfraID)
}
