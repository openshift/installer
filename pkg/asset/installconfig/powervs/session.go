package powervs

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	gohttp "net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	surveycore "github.com/AlecAivazis/survey/v2/core"
	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/form3tech-oss/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/intstr"

	machinev1 "github.com/openshift/api/machine/v1"
	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

var (
	defSessionTimeout   time.Duration = 9000000000000000000.0
	defRegion                         = "us_south"
	defaultAuthFilePath               = filepath.Join(os.Getenv("HOME"), ".powervs", "config.json")
)

// BxClient is struct which provides bluemix session details
type BxClient struct {
	*bxsession.Session
	APIKey               string
	Region               string
	Zone                 string
	PISession            *ibmpisession.IBMPISession
	User                 *User
	ServiceInstanceID    string
	PowerVSResourceGroup string
}

// User is struct with user details
type User struct {
	ID      string
	Email   string
	Account string
}

// SessionStore is an object and store that holds credentials and variables required to create a SessionVars object.
type SessionStore struct {
	ID                   string `json:"id,omitempty"`
	APIKey               string `json:"apikey,omitempty"`
	DefaultRegion        string `json:"region,omitempty"`
	DefaultZone          string `json:"zone,omitempty"`
	ServiceInstanceID    string `json:"serviceinstance,omitempty"`
	PowerVSResourceGroup string `json:"resourcegroup,omitempty"`
}

// SessionVars is an object that holds the variables required to create an ibmpisession object.
type SessionVars struct {
	ID                   string
	APIKey               string
	Region               string
	Zone                 string
	ServiceInstanceID    string
	PowerVSResourceGroup string
}

func authenticateAPIKey(sess *bxsession.Session) error {
	config := sess.Config
	tokenRefresher, err := authentication.NewIAMAuthRepository(config, &rest.Client{
		DefaultHeader: gohttp.Header{
			"User-Agent": []string{http.UserAgent()},
		},
	})
	if err != nil {
		return err
	}
	return tokenRefresher.AuthenticateAPIKey(config.BluemixAPIKey)
}

func fetchUserDetails(sess *bxsession.Session) (*User, error) {
	config := sess.Config
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

	return &user, nil
}

// NewBxClient func returns bluemix client
func NewBxClient(survey bool) (*BxClient, error) {
	c := &BxClient{}
	sv, err := getSessionVars(survey)
	if err != nil {
		return nil, err
	}

	c.APIKey = sv.APIKey
	c.Region = sv.Region
	c.Zone = sv.Zone
	c.ServiceInstanceID = sv.ServiceInstanceID
	c.PowerVSResourceGroup = sv.PowerVSResourceGroup

	bxSess, err := bxsession.New(&bluemix.Config{
		BluemixAPIKey: sv.APIKey,
	})
	if err != nil {
		return nil, err
	}
	if bxSess == nil {
		return nil, errors.New("failed to create bxsession.New in NewBxClient")
	}

	c.Session = bxSess

	err = authenticateAPIKey(bxSess)
	if err != nil {
		return nil, err
	}

	c.User, err = fetchUserDetails(bxSess)
	if err != nil {
		return nil, err
	}

	c.Session.Config.Region = powervs.Regions[sv.Region].VPCRegion

	return c, nil
}

func getSessionVars(survey bool) (SessionVars, error) {
	var sv SessionVars
	var ss SessionStore

	// Grab the session store from the installer written authFilePath
	logrus.Debug("Gathering credentials from AuthFile")
	err := getSessionStoreFromAuthFile(&ss)
	if err != nil {
		return sv, err
	}

	// Transfer the store to vars if they were found in the AuthFile
	sv.ID = ss.ID
	sv.APIKey = ss.APIKey
	sv.Region = ss.DefaultRegion
	sv.Zone = ss.DefaultZone
	sv.ServiceInstanceID = ss.ServiceInstanceID
	sv.PowerVSResourceGroup = ss.PowerVSResourceGroup

	// Grab variables from the users environment
	logrus.Debug("Gathering variables from user environment")
	err = getSessionVarsFromEnv(&sv)
	if err != nil {
		return sv, err
	}

	// Grab variable from the user themselves
	if survey {
		// Prompt the user for the first set of remaining variables.
		err = getFirstSessionVarsFromUser(&sv, &ss)
		if err != nil {
			return sv, err
		}

		// Transfer vars to the store to write out to the AuthFile
		ss.ID = sv.ID
		ss.APIKey = sv.APIKey
		ss.DefaultRegion = sv.Region
		ss.DefaultZone = sv.Zone
		ss.ServiceInstanceID = sv.ServiceInstanceID
		ss.PowerVSResourceGroup = sv.PowerVSResourceGroup

		// Save the session store to the disk.
		err = saveSessionStoreToAuthFile(&ss)
		if err != nil {
			return sv, err
		}

		// Since there is a minimal store at this point, it is safe
		// to call the function.
		// Prompt the user for the second set of remaining variables.
		err = getSecondSessionVarsFromUser(&sv, &ss)
		if err != nil {
			return sv, err
		}
	}

	// Transfer vars to the store to write out to the AuthFile
	ss.ID = sv.ID
	ss.APIKey = sv.APIKey
	ss.DefaultRegion = sv.Region
	ss.DefaultZone = sv.Zone
	ss.ServiceInstanceID = sv.ServiceInstanceID
	ss.PowerVSResourceGroup = sv.PowerVSResourceGroup

	// Save the session store to the disk.
	err = saveSessionStoreToAuthFile(&ss)
	if err != nil {
		return sv, err
	}

	return sv, nil
}

// ValidateDhcpService checks for existing Dhcp service for the provided PowerVS cloud instance
func (c *BxClient) ValidateDhcpService(ctx context.Context, svcInsID string, machineNetworks []types.MachineNetworkEntry) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Create PowerVS network client
	networkClient := instance.NewIBMPINetworkClient(ctx, c.PISession, svcInsID)
	if networkClient == nil {
		return errors.New("Failed to create a networkClient in ValidateDhcpService")
	}

	// Create PowerVS CloudConnection client
	cloudConnectionClient := instance.NewIBMPICloudConnectionClient(ctx, c.PISession, svcInsID)
	if cloudConnectionClient == nil {
		return errors.New("Failed to create a cloudConnectionClient in ValidateDhcpService")
	}

	allCloudConnecitons, err := cloudConnectionClient.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to get all existing Cloud Connections")
	}

	for _, singleCloudConnection := range allCloudConnecitons.CloudConnections {
		// Unfortunately, the Networks array is not filled in for a GetAll call :(
		cloudConnection, err := cloudConnectionClient.Get(*singleCloudConnection.CloudConnectionID)
		if err != nil {
			return errors.Wrap(err, "failed to get existing Cloud Connection details")
		}
		for _, ccNetwork := range cloudConnection.Networks {
			// The NetworkReference object does not provide subnet CIDRs.
			// So you have to get the network object based on the ID to find the CIDR.
			network, err := networkClient.Get(*ccNetwork.NetworkID)
			if err != nil {
				return errors.Wrap(err, "failed to get CC's network")
			}

			_, n1, err := net.ParseCIDR(*network.Cidr)
			if err != nil {
				return errors.Wrap(err, "failed to parse network.Cidr")
			}

			// Check each machineNetwork, typically one
			for _, machineNetwork := range machineNetworks {
				_, n2, err := net.ParseCIDR(machineNetwork.CIDR.String())
				if err != nil {
					return errors.Wrap(err, "failed to parse machineNetwork.CIDR")
				}
				if n2.Contains(n1.IP) || n1.Contains(n2.IP) {
					return fmt.Errorf("cidr conflicts with existing network")
				}
			}
		}
	}

	return nil
}

// ValidateCloudConnectionInPowerVSRegion counts cloud connection in PowerVS Region
func (c *BxClient) ValidateCloudConnectionInPowerVSRegion(ctx context.Context, svcInsID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	var cloudConnectionsIDs []string

	cloudConnectionClient := instance.NewIBMPICloudConnectionClient(ctx, c.PISession, svcInsID)

	//check number of cloudconnections
	getAllResp, err := cloudConnectionClient.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to get existing Cloud connection details")
	}

	if len(getAllResp.CloudConnections) >= 2 {
		return fmt.Errorf("cannot create new Cloud connection in Power VS. Only two Cloud connections are allowed per zone")
	}

	for _, cc := range getAllResp.CloudConnections {
		cloudConnectionsIDs = append(cloudConnectionsIDs, *cc.CloudConnectionID)
	}

	//check for Cloud connection attached to DHCP Service
	for _, cc := range cloudConnectionsIDs {
		cloudConn, err := cloudConnectionClient.Get(cc)
		if err != nil {
			return errors.Wrap(err, "failed to get Cloud connection details")
		}
		if cloudConn != nil {
			for _, nw := range cloudConn.Networks {
				if nw.DhcpManaged {
					return fmt.Errorf("only one Cloud connection can be attached to any DHCP network per account per zone")
				}
			}
		}
	}
	return nil
}

// GetSystemPools returns the system pools that are in the cloud.
func (c *BxClient) GetSystemPools(ctx context.Context, serviceInstanceID string) (models.SystemPools, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	systemPoolClient := instance.NewIBMPISystemPoolClient(ctx, c.PISession, serviceInstanceID)

	systemPools, err := systemPoolClient.GetSystemPools()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get system pools")
	}

	return systemPools, nil
}

// ValidateCapacityWithPools validates that the VMs created for both the controlPlanes and the
// computes will fit inside the given systemPools.
func ValidateCapacityWithPools(controlPlanes []machinev1beta1.Machine, computes []machinev1beta1.MachineSet, systemPools models.SystemPools) error {
	var (
		numCompute           int
		computeSystemType    string
		computeProcessorType string
		computeProcessors    float64
		computeMemoryGiB     int64
		numWorker            int64
		workerSystemType     string
		workerProcessorType  string
		workerProcessors     float64
		workerMemoryGiB      int64
		ok                   bool
	)

	// Find out the control plane master information
	numCompute = len(controlPlanes)
	ctrplConfigs := make([]*machinev1.PowerVSMachineProviderConfig, numCompute)
	for i, m := range controlPlanes {
		ctrplConfigs[i], ok = m.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig)
		if !ok {
			return errors.New("m.Spec.ProviderSpec.Value.Object failed")
		}
	}
	computeSystemType = ctrplConfigs[0].SystemType
	computeProcessorType = string(ctrplConfigs[0].ProcessorType)
	if ctrplConfigs[0].Processors.Type == intstr.Int {
		computeProcessors = float64(numCompute) * float64(ctrplConfigs[0].Processors.IntVal)
	} else {
		cores, err := strconv.ParseFloat(ctrplConfigs[0].Processors.StrVal, 64)
		if err != nil {
			return errors.Wrap(err, "failed to convert compute cores to a float")
		}
		computeProcessors = float64(numCompute) * cores
	}
	computeMemoryGiB = int64(numCompute) * int64(ctrplConfigs[0].MemoryGiB)

	// Find out the worker information
	computeReplicas := make([]int64, len(computes))
	computeConfigs := make([]*machinev1.PowerVSMachineProviderConfig, len(computes))
	for i, w := range computes {
		computeReplicas[i] = int64(*w.Spec.Replicas)
		numWorker = computeReplicas[i]
		computeConfigs[i], ok = w.Spec.Template.Spec.ProviderSpec.Value.Object.(*machinev1.PowerVSMachineProviderConfig)
		if !ok {
			return errors.New("w.Spec.Template.Spec.ProviderSpec.Value.Object")
		}

		workerSystemType = computeConfigs[i].SystemType
		workerProcessorType = string(computeConfigs[i].ProcessorType)
		if computeConfigs[i].Processors.Type == intstr.Int {
			workerProcessors = float64(computeReplicas[i]) * float64(computeConfigs[0].Processors.IntVal)
		} else {
			cores, err := strconv.ParseFloat(computeConfigs[0].Processors.StrVal, 64)
			if err != nil {
				return errors.Wrap(err, "failed to convert worker cores to a float")
			}
			workerProcessors = float64(computeReplicas[i]) * cores
		}
		workerMemoryGiB += numWorker * int64(computeConfigs[i].MemoryGiB)
	}

	// Helpful debug statement to save typing
	// fmt.Printf("ValidateCapacityWithPools: compute(%v) = {%v, %v, %v, %v}, worker(%v) = {%v, %v, %v, %v}\n", numCompute, computeSystemType, computeProcessorType, computeProcessors, computeMemoryGiB, numWorker, workerSystemType, workerProcessorType, workerProcessors, workerMemoryGiB)

	switch computeProcessorType {
	case "Capped":
	case "Dedicated":
	case "Shared":
		// @TODO I would think we should reduce the number of cores by some factor.
		// However, I cannot currently find documentation which describes what
		// PowerVS uses internally.
		computeProcessors = 0
	default:
		return errors.Errorf("Unknown compute processor type (%v)", computeProcessorType)
	}

	switch workerProcessorType {
	case "Capped":
	case "Dedicated":
	case "Shared":
		// @TODO I would think we should reduce the number of cores by some factor.
		// However, I cannot currently find documentation which describes what
		// PowerVS uses internally.
		workerProcessors = 0
	default:
		return errors.Errorf("Unknown worker processor type (%v)", workerProcessorType)
	}

	for _, systemPool := range systemPools {
		// Helpful debug statement to save typing
		// fmt.Printf("ValidateCapacityWithPools: pool %v, cores %v, memory %v\n", systemPool.Type, *systemPool.MaxCoresAvailable.Cores, *systemPool.MaxCoresAvailable.Memory)

		if computeSystemType == systemPool.Type {
			if computeProcessors > *systemPool.MaxCoresAvailable.Cores {
				return errors.Errorf("Not enough cores available (%v) for the compute nodes (need %v)", *systemPool.MaxCoresAvailable.Cores, computeProcessors)
			}
			*systemPool.MaxCoresAvailable.Cores -= computeProcessors

			if computeMemoryGiB > *systemPool.MaxCoresAvailable.Memory {
				return errors.Errorf("Not enough memory available (%v) for the compute nodes (need %v)", *systemPool.MaxCoresAvailable.Memory, computeMemoryGiB)
			}
			*systemPool.MaxCoresAvailable.Memory -= computeMemoryGiB
		}
		if workerSystemType == systemPool.Type {
			if workerProcessors > *systemPool.MaxCoresAvailable.Cores {
				return errors.Errorf("Not enough cores available (%v) for the worker nodes (need %v)", *systemPool.MaxCoresAvailable.Cores, workerProcessors)
			}
			*systemPool.MaxCoresAvailable.Cores -= workerProcessors

			if workerMemoryGiB > *systemPool.MaxCoresAvailable.Memory {
				return errors.Errorf("Not enough memory available (%v) for the worker nodes (need %v)", *systemPool.MaxCoresAvailable.Memory, workerMemoryGiB)
			}
			*systemPool.MaxCoresAvailable.Memory -= workerMemoryGiB
		}
	}

	return nil
}

// ValidateCapacity validates space for processors and storage in the cloud.
func (c *BxClient) ValidateCapacity(ctx context.Context, controlPlanes []machinev1beta1.Machine, computes []machinev1beta1.MachineSet, serviceInstanceID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	systemPools, err := c.GetSystemPools(ctx, serviceInstanceID)
	if err != nil {
		return errors.Wrap(err, "failed to get system pools")
	}

	// Call another function which we can also test with mock
	return ValidateCapacityWithPools(controlPlanes, computes, systemPools)
}

// NewPISession updates pisession details, return error on fail.
func (c *BxClient) NewPISession() error {
	var authenticator core.Authenticator = &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}

	// Create the session
	options := &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		UserAccount:   c.User.Account,
		Region:        c.Region,
		Zone:          c.Zone,
		Debug:         false,
	}

	// Avoid by defining err as a variable: non-name c.PISession on left side of :=
	var err error
	c.PISession, err = ibmpisession.NewIBMPISession(options)
	if err != nil {
		return err
	}

	return nil
}

// GetBxClientAPIKey returns the API key used by the Blue Mix Client.
func (c *BxClient) GetBxClientAPIKey() string {
	return c.APIKey
}

// getSessionStoreFromAuthFile gets the session creds from the auth file.
func getSessionStoreFromAuthFile(pss *SessionStore) error {
	if pss == nil {
		return errors.New("nil var: SessionStore")
	}

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	if _, err := os.Stat(authFilePath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(authFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, pss)
	if err != nil {
		return err
	}

	return nil
}

func getSessionVarsFromEnv(psv *SessionVars) error {
	if psv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(psv.ID) == 0 {
		psv.ID = os.Getenv("IBMID")
	}

	if len(psv.APIKey) == 0 {
		// APIKeyEnvVars is a list of environment variable names containing an IBM Cloud API key.
		var APIKeyEnvVars = []string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}
		psv.APIKey = getEnv(APIKeyEnvVars)
	}

	if len(psv.Region) == 0 {
		var regionEnvVars = []string{"IBMCLOUD_REGION", "IC_REGION"}
		psv.Region = getEnv(regionEnvVars)
	}

	if len(psv.Zone) == 0 {
		var zoneEnvVars = []string{"IBMCLOUD_ZONE"}
		psv.Zone = getEnv(zoneEnvVars)
	}

	if len(psv.ServiceInstanceID) == 0 {
		var serviceEnvVars = []string{"IBMCLOUD_SERVICE_INSTANCE"}
		psv.ServiceInstanceID = getEnv(serviceEnvVars)
	}

	if len(psv.PowerVSResourceGroup) == 0 {
		var resourceEnvVars = []string{"IBMCLOUD_RESOURCE_GROUP"}
		psv.PowerVSResourceGroup = getEnv(resourceEnvVars)
	}

	return nil
}

// Prompt the user for the first set of remaining variables.
// This is a chicken and egg problem.  We cannot call NewBxClient() or NewClient()
// yet for complicated questions to the user since those calls load the session
// variables from the store.  There is the possibility that the are empty at the
// moment.
func getFirstSessionVarsFromUser(psv *SessionVars, pss *SessionStore) error {
	var err error

	if psv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(psv.ID) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "IBM Cloud User ID",
					Help:    "The login for \nhttps://cloud.ibm.com/",
				},
			},
		}, &psv.ID)
		if err != nil {
			return errors.New("error saving the IBM Cloud User ID")
		}
	}

	if len(psv.APIKey) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Password{
					Message: "IBM Cloud API Key",
					Help:    "The API key installation.\nhttps://cloud.ibm.com/iam/apikeys",
				},
			},
		}, &psv.APIKey)
		if err != nil {
			return errors.New("error saving the API Key")
		}
	}

	return nil
}

// Prompt the user for the second set of remaining variables.
// This is a chicken and egg problem.  Now we can call NewBxClient() or NewClient()
// because the session store should at least have some minimal settings like the
// APIKey.
func getSecondSessionVarsFromUser(psv *SessionVars, pss *SessionStore) error {
	var (
		client *Client
		err    error
	)

	if psv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(psv.Region) == 0 {
		psv.Region, err = GetRegion(pss.DefaultRegion)
		if err != nil {
			return err
		}
	}

	if len(psv.Zone) == 0 {
		psv.Zone, err = GetZone(psv.Region, pss.DefaultZone)
		if err != nil {
			return err
		}
	}

	if len(psv.ServiceInstanceID) == 0 {
		if client == nil {
			client, err = NewClient()
			if err != nil {
				return fmt.Errorf("failed to powervs.NewClient: %w", err)
			}
		}

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
		defer cancel()

		serviceInstances, err := client.ListServiceInstances(ctx)
		if err != nil {
			return fmt.Errorf("failed to list serviceInstances: %w", err)
		}

		serviceInstancesSurvey := make([]string, len(serviceInstances))
		for i, serviceInstance := range serviceInstances {
			serviceInstancesSurvey[i] = strings.SplitN(serviceInstance, " ", 2)[0]
		}

		var serviceTransform survey.Transformer = func(ans interface{}) interface{} {
			switch v := ans.(type) {
			case surveycore.OptionAnswer:
				return surveycore.OptionAnswer{Value: strings.SplitN(v.Value, " ", 2)[1], Index: v.Index}
			case string:
				return strings.SplitN(ans.(string), " ", 2)[0]
			default:
				return ""
			}
		}

		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message: "Service Instance",
					Help:    "The Power VS service instance to be used for installation.",
					Default: "",
					Options: serviceInstances,
				},
				Transform: serviceTransform,
			},
		}, &psv.ServiceInstanceID)
		if err != nil {
			return fmt.Errorf("survey.ask failed with: %w", err)
		}
	}

	if len(psv.PowerVSResourceGroup) == 0 {
		if client == nil {
			client, err = NewClient()
			if err != nil {
				return fmt.Errorf("failed to powervs.NewClient: %w", err)
			}
		}

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
		defer cancel()

		resourceGroups, err := client.ListResourceGroups(ctx)
		if err != nil {
			return fmt.Errorf("failed to list resourceGroups: %w", err)
		}

		resourceGroupsSurvey := make([]string, len(resourceGroups.Resources))
		for i, resourceGroup := range resourceGroups.Resources {
			resourceGroupsSurvey[i] = *resourceGroup.Name
		}

		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Select{
					Message: "Resource Group",
					Help:    "The Power VS resource group to be used for installation.",
					Default: "",
					Options: resourceGroupsSurvey,
				},
			},
		}, &psv.PowerVSResourceGroup)
		if err != nil {
			return fmt.Errorf("survey.ask failed with: %w", err)
		}
	}

	return nil
}

func saveSessionStoreToAuthFile(pss *SessionStore) error {
	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	jsonVars, err := json.Marshal(*pss)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(authFilePath), 0700)
	if err != nil {
		return err
	}

	return os.WriteFile(authFilePath, jsonVars, 0o600)
}

func getEnv(envs []string) string {
	for _, k := range envs {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
