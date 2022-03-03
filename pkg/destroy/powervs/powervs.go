package powervs

import (
	"context"
	"fmt"
	gohttp "net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/resource/resourcev2/controllerv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/networking-go-sdk/dnsrecordsv1"
	"github.com/IBM/networking-go-sdk/zonesv1"
	"github.com/IBM/platform-services-go-sdk/resourcecontrollerv2"
	"github.com/IBM/platform-services-go-sdk/resourcemanagerv2"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/openshift/installer/pkg/destroy/providers"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

var (
	defaultTimeout = 15 * time.Minute
	stageTimeout   = 5 * time.Minute
)

const (
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

func fetchUserDetails(bxSession *bxsession.Session, generation int) (*User, error) {
	config := bxSession.Config
	user := User{}
	var bluemixToken string

	if strings.HasPrefix(config.IAMAccessToken, "Bearer") {
		bluemixToken = config.IAMAccessToken[7:len(config.IAMAccessToken)]
	} else {
		bluemixToken = config.IAMAccessToken
	}

	token, err := jwt.Parse(bluemixToken, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})
	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return &user, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if email, ok := claims["email"]; ok {
		user.Email = email.(string)
	}
	user.ID = claims["id"].(string)
	user.Account = claims["account"].(map[string]interface{})["bss"].(string)
	iss := claims["iss"].(string)
	if strings.Contains(iss, "https://iam.cloud.ibm.com") {
		user.cloudName = "bluemix"
	} else {
		user.cloudName = "staging"
	}
	user.cloudType = "public"

	user.generation = generation
	return &user, nil
}

// GetRegion converts from a zone into a region.
func GetRegion(zone string) (region string, err error) {
	err = nil
	switch {
	case strings.HasPrefix(zone, "dal"), strings.HasPrefix(zone, "us-south"):
		region = "us-south"
	case strings.HasPrefix(zone, "sao"):
		region = "sao"
	case strings.HasPrefix(zone, "us-east"):
		region = "us-east"
	case strings.HasPrefix(zone, "tor"):
		region = "tor"
	case strings.HasPrefix(zone, "eu-de-"):
		region = "eu-de"
	case strings.HasPrefix(zone, "lon"):
		region = "lon"
	case strings.HasPrefix(zone, "syd"):
		region = "syd"
	case strings.HasPrefix(zone, "tok"):
		region = "tok"
	case strings.HasPrefix(zone, "osa"):
		region = "osa"
	case strings.HasPrefix(zone, "mon"):
		region = "mon"
	default:
		return "", fmt.Errorf("region not found for the zone: %s", zone)
	}
	return
}

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	APIKey         string
	BaseDomain     string
	CISInstanceCRN string
	ClusterName    string
	Context        context.Context
	InfraID        string
	Logger         logrus.FieldLogger
	Region         string
	ServiceGUID    string
	VPCRegion      string
	Zone           string

	managementSvc  *resourcemanagerv2.ResourceManagerV2
	controllerSvc  *resourcecontrollerv2.ResourceControllerV2
	vpcSvc         *vpcv1.VpcV1
	zonesSvc       *zonesv1.ZonesV1
	dnsRecordsSvc  *dnsrecordsv1.DnsRecordsV1
	piSession      *ibmpisession.IBMPISession
	instanceClient *instance.IBMPIInstanceClient
	imageClient    *instance.IBMPIImageClient
	jobClient      *instance.IBMPIJobClient
	keyClient      *instance.IBMPIKeyClient

	resourceGroupID string
	cosInstanceID   string

	errorTracker
	pendingItemTracker
}

// New returns an IBMCloud destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (providers.Destroyer, error) {
	var apiKey string

	apiKey = os.Getenv("IBMCLOUD_API_KEY")
	if apiKey == "" {
		return nil, errors.New("environment variable IBMCLOUD_API_KEY must be set")
	}

	logger.Debugf("powervs.New: metadata.InfraID = %v", metadata.InfraID)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata = %v", metadata.ClusterPlatformMetadata)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS = %v", metadata.ClusterPlatformMetadata.PowerVS)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN = %v", metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup = %v", metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.Region = %v", metadata.ClusterPlatformMetadata.PowerVS.Region)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.VPCRegion = %v", metadata.ClusterPlatformMetadata.PowerVS.VPCRegion)
	logger.Debugf("powervs.New: metadata.ClusterPlatformMetadata.PowerVS.Zone = %v", metadata.ClusterPlatformMetadata.PowerVS.Zone)

	return &ClusterUninstaller{
		APIKey:             apiKey,
		BaseDomain:         metadata.ClusterPlatformMetadata.PowerVS.BaseDomain,
		ClusterName:        metadata.ClusterName,
		Context:            context.Background(),
		Logger:             logger,
		InfraID:            metadata.InfraID,
		CISInstanceCRN:     metadata.ClusterPlatformMetadata.PowerVS.CISInstanceCRN,
		Region:             metadata.ClusterPlatformMetadata.PowerVS.Region,
		VPCRegion:          metadata.ClusterPlatformMetadata.PowerVS.VPCRegion,
		Zone:               metadata.ClusterPlatformMetadata.PowerVS.Zone,
		pendingItemTracker: newPendingItemTracker(),
		resourceGroupID:    metadata.ClusterPlatformMetadata.PowerVS.PowerVSResourceGroup,
	}, nil
}

// Run is the entrypoint to start the uninstall process.
func (o *ClusterUninstaller) Run() (*types.ClusterQuota, error) {
	o.Logger.Debugf("powervs.Run")

	var ctx context.Context
	var deadline time.Time
	var ok bool
	var err error

	ctx, _ = o.contextWithTimeout()
	if ctx == nil {
		return nil, errors.Wrap(err, "powervs.Run: contextWithTimeout returns nil")
	}

	deadline, ok = ctx.Deadline()
	if !ok {
		return nil, errors.Wrap(err, "powervs.Run: failed to call ctx.Deadline")
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
		return false, errors.Wrap(err, "failed to destroy cluster")
	}

	return true, nil
}

func (o *ClusterUninstaller) destroyCluster() error {
	stagedFuncs := [][]struct {
		name    string
		execute func() error
	}{{
		{name: "Instances", execute: o.destroyInstances},
	}, {
		{name: "Load Balancers", execute: o.destroyLoadBalancers},
	}, {
		{name: "Images", execute: o.destroyImages},
		{name: "Security Groups", execute: o.destroySecurityGroups},
	}, {
		{name: "Cloud Object Storage Instances", execute: o.destroyCOSInstances},
		{name: "DNS Records", execute: o.destroyDNSRecords},
		{name: "SSH Keys", execute: o.destroySSHKeys},
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

	ctx, _ = o.contextWithTimeout()
	if ctx == nil {
		return errors.Wrap(err, "executeStageFunction contextWithTimeout returns nil")
	}

	deadline, ok = ctx.Deadline()
	if !ok {
		return errors.Wrap(err, "executeStageFunction failed to call ctx.Deadline")
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

// GetCISInstanceCRN gets the CRN name for the specified base domain.
func GetCISInstanceCRN(BaseDomain string) (string, error) {
	var CISInstanceCRN string = ""
	var APIKey string
	var bxSession *bxsession.Session
	var err error
	var tokenProviderEndpoint string = "https://iam.cloud.ibm.com"
	var tokenRefresher *authentication.IAMAuthRepository
	var authenticator *core.IamAuthenticator
	var controllerSvc *resourcecontrollerv2.ResourceControllerV2
	var listInstanceOptions *resourcecontrollerv2.ListResourceInstancesOptions
	var listResourceInstancesResponse *resourcecontrollerv2.ResourceInstancesList
	var instance resourcecontrollerv2.ResourceInstance
	var zonesService *zonesv1.ZonesV1
	var listZonesOptions *zonesv1.ListZonesOptions
	var listZonesResponse *zonesv1.ListZonesResp

	if APIKey = os.Getenv("IBMCLOUD_API_KEY"); len(APIKey) == 0 {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: environment variable IBMCLOUD_API_KEY not set")
	}
	bxSession, err = bxsession.New(&bluemix.Config{
		BluemixAPIKey:         APIKey,
		TokenProviderEndpoint: &tokenProviderEndpoint,
		Debug:                 false,
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: bxsession.New: %v", err)
	}
	tokenRefresher, err = authentication.NewIAMAuthRepository(bxSession.Config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authentication.NewIAMAuthRepository: %v", err)
	}
	err = tokenRefresher.AuthenticateAPIKey(bxSession.Config.BluemixAPIKey)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: tokenRefresher.AuthenticateAPIKey: %v", err)
	}
	authenticator = &core.IamAuthenticator{
		ApiKey: APIKey,
	}
	err = authenticator.Validate()
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: authenticator.Validate: %v", err)
	}
	// Instantiate the service with an API key based IAM authenticator
	controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
		ServiceName:   "cloud-object-storage",
		URL:           "https://resource-controller.cloud.ibm.com",
	})
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: creating ControllerV2 Service: %v", err)
	}
	listInstanceOptions = controllerSvc.NewListResourceInstancesOptions()
	listInstanceOptions.SetResourceID(cisServiceID)
	listResourceInstancesResponse, _, err = controllerSvc.ListResourceInstances(listInstanceOptions)
	if err != nil {
		return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListResourceInstances: %v", err)
	}
	for _, instance = range listResourceInstancesResponse.Resources {
		authenticator = &core.IamAuthenticator{
			ApiKey: APIKey,
		}

		err = authenticator.Validate()
		if err != nil {
		}

		zonesService, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
			Authenticator: authenticator,
			Crn:           instance.CRN,
		})
		if err != nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: NewZonesV1: %v", err)
		}
		listZonesOptions = zonesService.NewListZonesOptions()
		listZonesResponse, _, err = zonesService.ListZones(listZonesOptions)
		if listZonesResponse == nil {
			return CISInstanceCRN, fmt.Errorf("getCISInstanceCRN: ListZones: %v", err)
		}
		for _, zone := range listZonesResponse.Result {
			if *zone.Status == "active" {
				if *zone.Name == BaseDomain {
					CISInstanceCRN = *instance.CRN
				}
			}
		}
	}

	return CISInstanceCRN, nil
}

func (o *ClusterUninstaller) loadSDKServices() error {

	if o.APIKey == "" {
		return fmt.Errorf("loadSDKServices: missing APIKey in metadata.json")
	}

	var bxSession *bxsession.Session
	var tokenProviderEndpoint string = "https://iam.cloud.ibm.com"
	var err error

	bxSession, err = bxsession.New(&bluemix.Config{
		BluemixAPIKey:         o.APIKey,
		TokenProviderEndpoint: &tokenProviderEndpoint,
		Debug:                 false,
	})
	if err != nil {
		return fmt.Errorf("loadSDKServices: bxsession.New: %v", err)
	}

	tokenRefresher, err := authentication.NewIAMAuthRepository(bxSession.Config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		return fmt.Errorf("loadSDKServices: authentication.NewIAMAuthRepository: %v", err)
	}
	err = tokenRefresher.AuthenticateAPIKey(bxSession.Config.BluemixAPIKey)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		return fmt.Errorf("loadSDKServices: tokenRefresher.AuthenticateAPIKey: %v", err)
	}

	user, err := fetchUserDetails(bxSession, 2)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		return fmt.Errorf("loadSDKServices: fetchUserDetails: %v", err)
	}

	ctrlv2, err := controllerv2.New(bxSession)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		return fmt.Errorf("loadSDKServices: controllerv2.New: %v", err)
	}

	resourceClientV2 := ctrlv2.ResourceServiceInstanceV2()
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		return fmt.Errorf("loadSDKServices: ctrlv2.ResourceServiceInstanceV2: %v", err)
	}

	svcs, err := resourceClientV2.ListInstances(controllerv2.ServiceInstanceQuery{
		Type: "service_instance",
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		return fmt.Errorf("loadSDKServices: resourceClientV2.ListInstances: %v", err)
	}

	var serviceGUID string = ""

	for _, svc := range svcs {
		if svc.RegionID == o.Zone {
			serviceGUID = svc.Guid
			break
		}
	}
	if serviceGUID == "" {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		for _, svc := range svcs {
			o.Logger.Debugf("loadSDKServices: Guid = %v", svc.Guid)
			o.Logger.Debugf("loadSDKServices: RegionID = %v", svc.RegionID)
			o.Logger.Debugf("loadSDKServices: Name = %v", svc.Name)
			o.Logger.Debugf("loadSDKServices: Crn = %v", svc.Crn)
		}
		return fmt.Errorf("%s not found in list of service instances", o.Zone)
	}
	o.ServiceGUID = serviceGUID

	serviceInstance, err := resourceClientV2.GetInstance(o.ServiceGUID)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		return fmt.Errorf("loadSDKServices: resourceClientV2.GetInstance: %v", err)
	}

	region, err := GetRegion(serviceInstance.RegionID)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		return fmt.Errorf("loadSDKServices: GetRegion: %v", err)
	}

	var authenticator core.Authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		return fmt.Errorf("loadSDKServices: loadSDKServices: authenticator.Validate: %v", err)
	}

	var options *ibmpisession.IBMPIOptions = &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		Debug:         false,
		Region:        region,
		UserAccount:   user.Account,
		Zone:          serviceInstance.RegionID,
	}

	o.piSession, err = ibmpisession.NewIBMPISession(options)
	if (err != nil) || (o.piSession == nil) {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		if err != nil {
			return fmt.Errorf("loadSDKServices: ibmpisession.New: %v", err)
		}
		return fmt.Errorf("loadSDKServices: loadSDKServices: o.piSession is nil")
	}

	ctx, _ := o.contextWithTimeout()

	o.instanceClient = instance.NewIBMPIInstanceClient(ctx, o.piSession, o.ServiceGUID)
	if o.instanceClient == nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		return fmt.Errorf("loadSDKServices: loadSDKServices: o.instanceClient is nil")
	}

	o.imageClient = instance.NewIBMPIImageClient(ctx, o.piSession, o.ServiceGUID)
	if o.imageClient == nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		return fmt.Errorf("loadSDKServices: loadSDKServices: o.imageClient is nil")
	}

	o.jobClient = instance.NewIBMPIJobClient(ctx, o.piSession, o.ServiceGUID)
	if o.jobClient == nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		return fmt.Errorf("loadSDKServices: loadSDKServices: o.jobClient is nil")
	}

	o.keyClient = instance.NewIBMPIKeyClient(ctx, o.piSession, o.ServiceGUID)
	if o.keyClient == nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		return fmt.Errorf("loadSDKServices: loadSDKServices: o.keyClient is nil")
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		return fmt.Errorf("loadSDKServices: loadSDKServices: authenticator.Validate: %v", err)
	}

	// VpcV1
	o.vpcSvc, err = vpcv1.NewVpcV1(&vpcv1.VpcV1Options{
		Authenticator: authenticator,
		URL:           "https://" + o.VPCRegion + ".iaas.cloud.ibm.com/v1",
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		return fmt.Errorf("loadSDKServices: loadSDKServices: vpcv1.NewVpcV1: %v", err)
	}

	userAgentString := fmt.Sprintf("OpenShift/4.x Destroyer/%s", version.Raw)
	o.vpcSvc.Service.SetUserAgent(userAgentString)

	authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
	}

	// Instantiate the service with an API key based IAM authenticator
	o.managementSvc, err = resourcemanagerv2.NewResourceManagerV2(&resourcemanagerv2.ResourceManagerV2Options{
		Authenticator: authenticator,
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		return fmt.Errorf("loadSDKServices: loadSDKServices: creating ResourceManagerV2 Service: %v", err)
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
	}

	// Instantiate the service with an API key based IAM authenticator
	o.controllerSvc, err = resourcecontrollerv2.NewResourceControllerV2(&resourcecontrollerv2.ResourceControllerV2Options{
		Authenticator: authenticator,
		ServiceName:   "cloud-object-storage",
		URL:           "https://resource-controller.cloud.ibm.com",
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		o.Logger.Debugf("loadSDKServices: o.managementSvc = %v", o.managementSvc)
		return fmt.Errorf("loadSDKServices: loadSDKServices: creating ControllerV2 Service: %v", err)
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
	}

	o.zonesSvc, err = zonesv1.NewZonesV1(&zonesv1.ZonesV1Options{
		Authenticator: authenticator,
		Crn:           &o.CISInstanceCRN,
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		o.Logger.Debugf("loadSDKServices: o.managementSvc = %v", o.managementSvc)
		o.Logger.Debugf("loadSDKServices: o.controllerSvc = %v", o.controllerSvc)
		return fmt.Errorf("loadSDKServices: loadSDKServices: creating zonesSvc: %v", err)
	}

	// Get the Zone ID
	zoneOptions := o.zonesSvc.NewListZonesOptions()
	zoneResources, detailedResponse, err := o.zonesSvc.ListZonesWithContext(ctx, zoneOptions)
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		o.Logger.Debugf("loadSDKServices: o.managementSvc = %v", o.managementSvc)
		o.Logger.Debugf("loadSDKServices: o.controllerSvc = %v", o.controllerSvc)
		return fmt.Errorf("loadSDKServices: loadSDKServices: Failed to list Zones: %v and the response is: %s", err, detailedResponse)
	}

	zoneID := ""
	for _, zone := range zoneResources.Result {
		o.Logger.Debugf("loadSDKServices: Zone: %v", *zone.Name)
		if strings.Contains(o.BaseDomain, *zone.Name) {
			zoneID = *zone.ID
		}
	}

	authenticator = &core.IamAuthenticator{
		ApiKey: o.APIKey,
	}

	err = authenticator.Validate()
	if err != nil {
	}

	o.dnsRecordsSvc, err = dnsrecordsv1.NewDnsRecordsV1(&dnsrecordsv1.DnsRecordsV1Options{
		Authenticator:  authenticator,
		Crn:            &o.CISInstanceCRN,
		ZoneIdentifier: &zoneID,
	})
	if err != nil {
		o.Logger.Debugf("loadSDKServices: bxSession = %v", bxSession)
		o.Logger.Debugf("loadSDKServices: tokenRefresher = %v", tokenRefresher)
		o.Logger.Debugf("loadSDKServices: ctrlv2 = %v", ctrlv2)
		o.Logger.Debugf("loadSDKServices: resourceClientV2 = %v", resourceClientV2)
		o.Logger.Debugf("loadSDKServices: o.ServiceGUID = %v", o.ServiceGUID)
		o.Logger.Debugf("loadSDKServices: serviceInstance = %v", serviceInstance)
		o.Logger.Debugf("loadSDKServices: o.piSession = %v", o.piSession)
		o.Logger.Debugf("loadSDKServices: o.instanceClient = %v", o.instanceClient)
		o.Logger.Debugf("loadSDKServices: o.imageClient = %v", o.imageClient)
		o.Logger.Debugf("loadSDKServices: o.jobClient = %v", o.jobClient)
		o.Logger.Debugf("loadSDKServices: o.vpcSvc = %v", o.vpcSvc)
		o.Logger.Debugf("loadSDKServices: o.managementSvc = %v", o.managementSvc)
		o.Logger.Debugf("loadSDKServices: o.controllerSvc = %v", o.controllerSvc)
		return fmt.Errorf("loadSDKServices: loadSDKServices: Failed to instantiate dnsRecordsSvc: %v", err)
	}

	return nil
}

func (o *ClusterUninstaller) contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(o.Context, defaultTimeout)
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
		return errors.Errorf("%d items pending", pending[0])
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
