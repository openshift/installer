package powervs

import (
	"context"
	"errors"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"sync"
	"time"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/dnszonesv1"
	"github.com/IBM/networking-go-sdk/resourcerecordsv1"
	"github.com/IBM/networking-go-sdk/transitgatewayapisv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/version"
)

var (
	defaultTimeout = 30 * time.Minute
	stageTimeout   = 15 * time.Minute
)

func leftInContext(ctx context.Context) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		return math.MaxInt64
	}

	duration := time.Until(deadline)

	return duration
}

const (
	// resource Id for Power Systems Virtual Server in the Global catalog.
	powerIAASResourceID = "abd259f0-9990-11e8-acc8-b9f54a8f1661"

	// cisServiceID is the Cloud Internet Services' catalog service ID.
	cisServiceID = "75874a60-cb12-11e7-948e-37ac098eb1b9"
)

// User information.
type User struct {
	ID         string
	Email      string
	Account    string
	cloudName  string `default:"bluemix"`
	cloudType  string `default:"public"`
	generation int    `default:"2"`
}

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	APIKey             string
	BaseDomain         string
	CISInstanceCRN     string
	ClusterName        string
	DNSInstanceCRN     string
	DNSZone            string
	InfraID            string
	Logger             logrus.FieldLogger
	Region             string
	ServiceGUID        string
	VPCRegion          string
	Zone               string
	TransitGatewayName string

	managementSvc      *resourcemanagerv2.ResourceManagerV2
	controllerSvc      *resourcecontrollerv2.ResourceControllerV2
	vpcSvc             *vpcv1.VpcV1
	zonesSvc           *zonesv1.ZonesV1
	dnsRecordsSvc      *dnsrecordsv1.DnsRecordsV1
	dnsZonesSvc        *dnszonesv1.DnsZonesV1
	resourceRecordsSvc *resourcerecordsv1.ResourceRecordsV1
	piSession          *ibmpisession.IBMPISession
	instanceClient     *instance.IBMPIInstanceClient
	imageClient        *instance.IBMPIImageClient
	jobClient          *instance.IBMPIJobClient
	keyClient          *instance.IBMPIKeyClient
	dhcpClient         *instance.IBMPIDhcpClient
	networkClient      *instance.IBMPINetworkClient
	tgClient           *transitgatewayapisv1.TransitGatewayApisV1

	resourceGroupID string
	cosInstanceID   string
	dnsZoneID       string

	errorTracker
	pendingItemTracker
}

// New returns an IBMCloud destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	var (
		bxClient *powervs.BxClient
		APIKey   string
		err      error
	)

	// We need to prompt for missing variables because NewPISession requires them!
	bxClient, err = powervs.NewBxClient(true)
	if err != nil {
		return nil, err
	}
	APIKey = bxClient.GetBxClientAPIKey()
	if APIKey == "" {
		return nil, errors.New("powervs.GetSession did not return an API key")
	}

	logger.Debugf("powervs.New: metadata.InfraID = %v", metadata.InfraID)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.BaseDomain = %v", metadata.ClusterPlatformMetadata.PowerVS.BaseDomain)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN = %v", metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.DNSInstanceCRN = %v", metadata.ClusterPlatformMetadata.PowerVS.DNSInstanceCRN)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup = %v", metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.Region = %v", metadata.ClusterPlatformMetadata.PowerVS.Region)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.VPCRegion = %v", metadata.ClusterPlatformMetadata.PowerVS.VPCRegion)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.Zone = %v", metadata.ClusterPlatformMetadata.PowerVS.Zone)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.ServiceInstanceGUID = %v", metadata.ClusterPlatformMetadata.PowerVS.ServiceInstanceGUID)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.TransitGatewayName = %v", metadata.ClusterPlatformMetadata.PowerVS.TransitGatewayName)

	// Handle an optional setting in install-config.yaml
	if metadata.ClusterPlatformMetadata.PowerVS.VPCRegion == "" {
		var derivedVPCRegion string
		if derivedVPCRegion, err = powervstypes.VPCRegionForPowerVSRegion(metadata.ClusterPlatformMetadata.PowerVS.Region); err != nil {
			return nil, fmt.Errorf("powervs.New failed to derive VPCRegion: %w", err)
		}
		logger.Debugf("powervs.New: PowerVS.VPCRegion is missing, derived VPCRegion = %v", derivedVPCRegion)
		metadata.ClusterPlatformMetadata.PowerVS.VPCRegion = derivedVPCRegion
	}

	return &ClusterUninstaller{
		APIKey:             APIKey,
		BaseDomain:         metadata.ClusterPlatformMetadata.PowerVS.BaseDomain,
		ClusterName:        metadata.ClusterName,
		Logger:             logger,
		InfraID:            metadata.InfraID,
		CISInstanceCRN:     metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN,
		DNSInstanceCRN:     metadata.ClusterPlatformMetadata.PowerVS.DNSInstanceCRN,
		Region:             metadata.ClusterPlatformMetadata.PowerVS.Region,
		VPCRegion:          metadata.ClusterPlatformMetadata.PowerVS.VPCRegion,
		Zone:               metadata.ClusterPlatformMetadata.PowerVS.Zone,
		pendingItemTracker: newPendingItemTracker(),
		resourceGroupID:    metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup,
		ServiceGUID:        metadata.ClusterPlatformMetadata.PowerVS.ServiceInstanceGUID,
		TransitGatewayName: metadata.ClusterPlatformMetadata.PowerVS.TransitGatewayName,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	o.Logger.Debugf("powervs.Run")

	var ctx context.Context
	var deadline time.Time
	var ok bool
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	if ctx == nil {
		return nil, fmt.Errorf("powervs.Run: contextWithTimeout returns nil: %w", err)
	}

	deadline, ok = ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("powervs.Run: failed to call ctx.Deadline: %w", err)
	}

	var duration time.Duration = deadline.Sub(time.Now())

	o.Logger.Debugf("powervs.Run: duration = %v", duration)

	if duration <= 0 {
		return nil, fmt.Errorf("powervs.Run: duration is <= 0 (%v)", duration)
	}

	err = wait.PollImmediateInfinite(
		duration,
		o.PolledRun,
	)

	o.Logger.Debugf("powervs.Run: after wait.PollImmediateInfinite, err = %v", err)

	return nil, err
}

// PolledRun is the Run function which will be called with Polling.
func (o *ClusterUninstaller) PolledRun() (bool, error) {
	o.Logger.Debugf("powervs.PolledRun")

	var err error

	err = o.loadSDKServices()
	if err != nil {
		o.Logger.Debugf("powervs.PolledRun: Failed loadSDKServices")
		return false, err
	}

	err = o.destroyCluster()
	if err != nil {
		o.Logger.Debugf("powervs.PolledRun: Failed destroyCluster")
		return false, fmt.Errorf("failed to destroy cluster: %w", err)
	}

	return true, nil
}

func (o *ClusterUninstaller) destroyCluster() error {
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Transit Gateways", execute: o.destroyTransitGateways},
	}, {
		{name: "Cloud Instances", execute: o.destroyCloudInstances},
	}, {
		{name: "Power Instances", execute: o.destroyPowerInstances},
	}, {
		{name: "Load Balancers", execute: o.destroyLoadBalancers},
	}, {
		{name: "Cloud Subnets", execute: o.destroyCloudSubnets},
	}, {
		{name: "Public Gateways", execute: o.destroyPublicGateways},
	}, {
		{name: "DHCPs", execute: o.destroyDHCPNetworks},
	}, {
		{name: "Images", execute: o.destroyImages},
		{name: "VPCs", execute: o.destroyVPCs},
	}, {
		{name: "Security Groups", execute: o.destroySecurityGroups},
	}, {
		{name: "Cloud Object Storage Instances", execute: o.destroyCOSInstances},
		{name: "Cloud SSH Keys", execute: o.destroyCloudSSHKeys},
		{name: "Power SSH Keys", execute: o.destroyPowerSSHKeys},
	}, {
		{name: "DNS Records", execute: o.destroyDNSRecords},
		{name: "DNS Resource Records", execute: o.destroyResourceRecords},
	}, {
		{name: "Power Subnets", execute: o.destroyPowerSubnets},
	}, {
		{name: "Service Instances", execute: o.destroyServiceInstances},
	}}

	for _, stage := range stagedFuncs {
		var wg sync.WaitGroup
		errCh := make(chan error)
		wgDone := make(chan bool)

		for _, f := range stage {
			wg.Add(1)
			// Start a parallel goroutine
			go o.executeStageFunction(f, errCh, &wg)
		}

		// Start a parallel goroutine
		go func() {
			wg.Wait()
			close(wgDone)
		}()

		select {
		// Did the wait group goroutine finish?
		case <-wgDone:
			// On to the next stage
			o.Logger.Debugf("destroyCluster: <-wgDone")
			continue
		// Have we taken too long?
		case <-time.After(stageTimeout):
			return fmt.Errorf("destroyCluster: timed out")
		// Has an error been sent via the channel?
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
	o.Logger.Debugf("executeStageFunction: Adding: %s", f.name)

	defer wg.Done()

	var ctx context.Context
	var deadline time.Time
	var ok bool
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	if ctx == nil {
		return fmt.Errorf("executeStageFunction contextWithTimeout returns nil: %w", err)
	}

	deadline, ok = ctx.Deadline()
	if !ok {
		return fmt.Errorf("executeStageFunction failed to call ctx.Deadline: %w", err)
	}

	var duration time.Duration = deadline.Sub(time.Now())

	o.Logger.Debugf("executeStageFunction: duration = %v", duration)
	if duration <= 0 {
		return fmt.Errorf("executeStageFunction: duration is <= 0 (%v)", duration)
	}

	err = wait.PollImmediateInfinite(
		duration,
		func() (bool, error) {
			var err error

			o.Logger.Debugf("executeStageFunction: Executing: %s", f.name)

			err = f.execute()
			if err != nil {
				o.Logger.Debugf("ERROR: executeStageFunction: %s: %v", f.name, err)

				return false, err
			}

			return true, nil
		},
	)

	if err != nil {
		errCh <- err
	}
	return nil
}

func (o *ClusterUninstaller) newAuthenticator(apikey string) (core.Authenticator, error) {
	var (
		authenticator core.Authenticator
		err           error
	)

	if apikey == "" {
		return nil, errors.New("newAuthenticator: apikey is empty")
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: apikey,
	}

	err = authenticator.Validate()
	if err != nil {
		return nil, fmt.Errorf("newAuthenticator: authenticator.Validate: %w", err)
	}

	return authenticator, nil
}

func (o *ClusterUninstaller) loadSDKServices() error {
	var (
		err           error
		authenticator core.Authenticator
		versionDate   = "2023-07-04"
		tgOptions     *transitgatewayapisv1.TransitGatewayApisV1Options
		serviceName   string
	)

	defer func() {
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.keyClient = %v", o.keyClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		o.Logger.Debugf("loadSDKServices: o.managementSvc = %v", o.managementSvc)
		o.Logger.Debugf("loadSDKServices: o.controllerSvc = %v", o.controllerSvc)
	}()

	if o.APIKey == "" {
		return fmt.Errorf("loadSDKServices: missing APIKey in metadata.json")
	}

	user, err := powervs.FetchUserDetails(o.APIKey)
	if err != nil {
		return fmt.Errorf("loadSDKServices: fetchUserDetails: %w", err)
	}

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return err
	}

	var options *ibmpisession.IBMPIOptions = &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		Debug:         false,
		UserAccount:   user.Account,
		Zone:          o.Zone,
	}

	o.piSession, err = ibmpisession.NewIBMPISession(options)
	if (err != nil) || (o.piSession == nil) {
		if err != nil {
			return fmt.Errorf("loadSDKServices: ibmpisession.New: %w", err)
		}
		return fmt.Errorf("loadSDKServices: o.piSession is nil")
	}

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return err
	}

	// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
	o.vpcSvc, err = vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: authenticator,
		URL:           "https://" + o.VPCRegion + ".iaas.cloud.ibm.com/v1",
	})
	if err != nil {
		return fmt.Errorf("loadSDKServices: vpcv1.NewVpcV1: %w", err)
	}

	userAgentString := fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)
	o.vpcSvc.Service.SetUserAgent(userAgentString)

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return err
	}

	// Instantiate the service with an API key based IAM authenticator
	o.managementSvc, err = resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	})
	if err != nil {
		return fmt.Errorf("loadSDKServices: creating ResourceManagerV2 Service: %w", err)
	}

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return err
	}

	// Instantiate the service with an API key based IAM authenticator
	o.controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
		ServiceName:   "cloud-object-storage",
		URL:           "https://resource-controller.cloud.ibm.com",
	})
	if err != nil {
		return fmt.Errorf("loadSDKServices: creating ControllerV2 Service: %w", err)
	}

	authenticator, err = o.newAuthenticator(o.APIKey)
	if err != nil {
		return err
	}

	tgOptions = &transitgatewayapisv1.TransitGatewayApisV1Options{
		Authenticator: authenticator,
		Version:       &versionDate,
	}

	o.tgClient, err = transitgatewayapisv1.NewTransitGatewayApisV1(tgOptions)
	if err != nil {
		return fmt.Errorf("loadSDKServices: NewTransitGatewayApisV1: %w", err)
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	// Either CISInstanceCRN is set or DNSInstanceCRN is set. Both should not be set at the same time,
	// but check both just to be safe.
	if len(o.CISInstanceCRN) > 0 {
		authenticator, err = o.newAuthenticator(o.APIKey)
		if err != nil {
			return err
		}

		o.zonesSvc, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           &o.CISInstanceCRN,
		})
		if err != nil {
			return fmt.Errorf("loadSDKServices: creating zonesSvc: %w", err)
		}

		// Get the Zone ID
		zoneOptions := o.zonesSvc.NewListZonesOptions()
		zoneResources, detailedResponse, err := o.zonesSvc.ListZonesWithContext(ctx, zoneOptions)
		if err != nil {
			return fmt.Errorf("loadSDKServices: Failed to list Zones: %w and the response is: %s", err, detailedResponse)
		}

		for _, zone := range zoneResources.Result {
			o.Logger.Debugf("loadSDKServices: Zone: %v", *zone.Name)
			if strings.Contains(o.BaseDomain, *zone.Name) {
				o.dnsZoneID = *zone.ID
			}
		}

		authenticator, err = o.newAuthenticator(o.APIKey)
		if err != nil {
			return err
		}

		o.dnsRecordsSvc, err = dnsrecordsv1.NewDnsRecordsV1(&dnsrecordsv1.DnsRecordsV1Options{
			Authenticator:  authenticator,
			Crn:            &o.CISInstanceCRN,
			ZoneIdentifier: &o.dnsZoneID,
		})
		if err != nil {
			return fmt.Errorf("loadSDKServices: Failed to instantiate dnsRecordsSvc: %w", err)
		}
	}

	if len(o.DNSInstanceCRN) > 0 {
		authenticator, err = o.newAuthenticator(o.APIKey)
		if err != nil {
			return err
		}

		o.dnsZonesSvc, err = dnszonesv1.NewDnsZonesV1(&dnszonesv1.DnsZonesV1Options{
			Authenticator: authenticator,
		})
		if err != nil {
			return fmt.Errorf("loadSDKServices: creating zonesSvc: %w", err)
		}

		// Get the Zone ID
		dnsCRN, err := crn.Parse(o.DNSInstanceCRN)
		if err != nil {
			return fmt.Errorf("failed to parse DNSInstanceCRN: %w", err)
		}
		options := o.dnsZonesSvc.NewListDnszonesOptions(dnsCRN.ServiceInstance)
		listZonesResponse, detailedResponse, err := o.dnsZonesSvc.ListDnszones(options)
		if err != nil {
			return fmt.Errorf("loadSDKServices: Failed to list Zones: %w and the response is: %s", err, detailedResponse)
		}

		for _, zone := range listZonesResponse.Dnszones {
			o.Logger.Debugf("loadSDKServices: Zone: %v", *zone.Name)
			if strings.Contains(o.BaseDomain, *zone.Name) {
				o.dnsZoneID = *zone.ID
			}
		}

		authenticator, err = o.newAuthenticator(o.APIKey)
		if err != nil {
			return err
		}

		o.resourceRecordsSvc, err = resourcerecordsv1.NewResourceRecordsV1(&resourcerecordsv1.ResourceRecordsV1Options{
			Authenticator: authenticator,
		})
		if err != nil {
			return fmt.Errorf("loadSDKServices: Failed to instantiate resourceRecordsSvc: %w", err)
		}
	}

	o.Logger.Debugf("loadSDKServices: o.resourceGroupID = %s", o.resourceGroupID)
	// If the user passes in a human readable resource group id, then we need to convert it to a UUID
	listGroupOptions := o.managementSvc.NewListResourceGroupsOptions()
	groups, _, err := o.managementSvc.ListResourceGroupsWithContext(ctx, listGroupOptions)
	if err != nil {
		return fmt.Errorf("loadSDKServices: Failed to list resource groups: %w", err)
	}
	for _, group := range groups.Resources {
		if *group.Name == o.resourceGroupID {
			o.Logger.Debugf("loadSDKServices: resource FOUND: %s %s", *group.Name, *group.ID)
			o.resourceGroupID = *group.ID
		} else {
			o.Logger.Debugf("loadSDKServices: resource SKIP:  %s %s", *group.Name, *group.ID)
		}
	}
	o.Logger.Debugf("loadSDKServices: o.resourceGroupID = %s", o.resourceGroupID)

	// If we should have created a service instance dynamically
	if o.ServiceGUID == "" {
		serviceName = fmt.Sprintf("%s-power-iaas", o.InfraID)
		o.Logger.Debugf("loadSDKServices: serviceName = %v", serviceName)

		o.ServiceGUID, err = o.ServiceInstanceNameToGUID(context.Background(), serviceName)
		if err != nil {
			return fmt.Errorf("loadSDKServices: ServiceInstanceNameToGUID: %w", err)
		}
	}
	if o.ServiceGUID == "" {
		// The rest of this function relies on o.ServiceGUID, so finish now!
		return nil
	}

	o.instanceClient = instance.NewIBMPIInstanceClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.instanceClient == nil {
		return fmt.Errorf("loadSDKServices: o.instanceClient is nil")
	}

	o.imageClient = instance.NewIBMPIImageClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.imageClient == nil {
		return fmt.Errorf("loadSDKServices: o.imageClient is nil")
	}

	o.jobClient = instance.NewIBMPIJobClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.jobClient == nil {
		return fmt.Errorf("loadSDKServices: o.jobClient is nil")
	}

	o.keyClient = instance.NewIBMPIKeyClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.keyClient == nil {
		return fmt.Errorf("loadSDKServices: o.keyClient is nil")
	}

	o.dhcpClient = instance.NewIBMPIDhcpClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.dhcpClient == nil {
		return fmt.Errorf("loadSDKServices: o.dhcpClient is nil")
	}

	o.networkClient = instance.NewIBMPINetworkClient(context.Background(), o.piSession, o.ServiceGUID)
	if o.networkClient == nil {
		return fmt.Errorf("loadSDKServices: o.networkClient is nil")
	}

	return nil
}

// ServiceInstanceNameToGUID returns the GUID of the matching service instance name which was passed in.
func (o *ClusterUninstaller) ServiceInstanceNameToGUID(ctx context.Context, name string) (string, error) {
	var (
		options   *resourcecontrollerv2.ListResourceInstancesOptions
		resources *resourcecontrollerv2.ResourceInstancesList
		err       error
		perPage   int64 = 10
		moreData        = true
		nextURL   *string
	)

	options = o.controllerSvc.NewListResourceInstancesOptions()
	options.SetResourceGroupID(o.resourceGroupID)
	// resource ID for Power Systems Virtual Server in the Global catalog
	options.SetResourceID(powerIAASResourceID)
	options.SetLimit(perPage)

	for moreData {
		resources, _, err = o.controllerSvc.ListResourceInstancesWithContext(ctx, options)
		if err != nil {
			return "", fmt.Errorf("failed to list resource instances: %w", err)
		}

		for _, resource := range resources.Resources {
			var (
				getResourceOptions *resourcecontrollerv2.GetResourceInstanceOptions
				resourceInstance   *resourcecontrollerv2.ResourceInstance
				response           *core.DetailedResponse
			)

			o.Logger.Debugf("ServiceInstanceNameToGUID: resource.Name = %s", *resource.Name)

			getResourceOptions = o.controllerSvc.NewGetResourceInstanceOptions(*resource.ID)

			resourceInstance, response, err = o.controllerSvc.GetResourceInstance(getResourceOptions)
			if err != nil {
				return "", fmt.Errorf("failed to get instance: %w", err)
			}
			if response != nil && response.StatusCode == gohttp.StatusNotFound || response.StatusCode == gohttp.StatusInternalServerError {
				return "", fmt.Errorf("failed to get instance, response is: %v", response)
			}

			if resourceInstance.Type == nil {
				o.Logger.Debugf("ServiceInstanceNameToGUID: type: nil")
				continue
			}
			o.Logger.Debugf("ServiceInstanceNameToGUID: type: %v", *resourceInstance.Type)
			if resourceInstance.GUID == nil {
				o.Logger.Debugf("ServiceInstanceNameToGUID: GUID: nil")
				continue
			}
			if *resourceInstance.Type != "service_instance" && *resourceInstance.Type != "composite_instance" {
				continue
			}
			if *resourceInstance.Name != name {
				continue
			}

			o.Logger.Debugf("ServiceInstanceNameToGUID: Found match!")

			return *resourceInstance.GUID, nil
		}

		// Based on: https://cloud.ibm.com/apidocs/resource-controller/resource-controller?code=go#list-resource-instances
		nextURL, err = core.GetQueryParam(resources.NextURL, "start")
		if err != nil {
			return "", fmt.Errorf("failed to GetQueryParam on start: %w", err)
		}
		if nextURL == nil {
			options.SetStart("")
		} else {
			options.SetStart(*nextURL)
		}

		moreData = *resources.RowsCount == perPage
	}

	return "", nil
}

func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

func (o *ClusterUninstaller) timeout(ctx context.Context) bool {
	var deadline time.Time
	var ok bool

	deadline, ok = ctx.Deadline()
	if !ok {
		o.Logger.Debugf("timeout: deadline, ok = %v, %v", deadline, ok)
		return true
	}

	var after bool = time.Now().After(deadline)

	if after {
		// 01/02 03:04:05PM ‘06 -0700
		o.Logger.Debugf("timeout: after deadline! (%v)", deadline.Format("2006-01-02 03:04:05PM"))
	}

	return after
}

type ibmError struct {
	Status  int
	Message string
}

func isNoOp(err *ibmError) bool {
	if err == nil {
		return false
	}

	return err.Status == gohttp.StatusNotFound
}

// aggregateError is a utility function that takes a slice of errors and an
// optional pending argument, and returns an error or nil.
func aggregateError(errs []error, pending ...int) error {
	err := utilerrors.NewAggregate(errs)
	if err != nil {
		return err
	}
	if len(pending) > 0 && pending[0] > 0 {
		return fmt.Errorf("%d items pending", pending[0])
	}
	return nil
}

// pendingItemTracker tracks a set of pending item names for a given type of resource.
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
