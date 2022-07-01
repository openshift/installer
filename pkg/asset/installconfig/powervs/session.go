package powervs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	gohttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/pkg/errors"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/form3tech-oss/jwt-go"

	"github.com/sirupsen/logrus"
)

var (
	defSessionTimeout   time.Duration = 9000000000000000000.0
	defRegion                         = "us_south"
	defaultAuthFilePath               = filepath.Join(os.Getenv("HOME"), ".powervs", "config.json")
)

//BxClient is struct which provides bluemix session details
type BxClient struct {
	*bxsession.Session
	APIKey       string
	PISession    *ibmpisession.IBMPISession
	User         *User
	AccountAPIV2 accountv2.Accounts
}

//User is struct with user details
type User struct {
	ID      string
	Email   string
	Account string
}

// PISessionVars is an object that holds the variables required to create an ibmpisession object.
type PISessionVars struct {
	ID     string `json:"id,omitempty"`
	APIKey string `json:"apikey,omitempty"`
	Region string `json:"region,omitempty"`
	Zone   string `json:"zone,omitempty"`
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

//NewBxClient func returns bluemix client
func NewBxClient() (*BxClient, error) {
	c := &BxClient{}

	var pisv PISessionVars
	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err := getPISessionVarsFromAuthFile(&pisv)
	if err != nil {
		return nil, err
	}

	// Grab variables from the users environment
	logrus.Debug("Gathering variables from user environment")
	err = getPISessionVarsFromEnv(&pisv)
	if err != nil {
		return nil, err
	}

	// Prompt the user for the remaining variables.
	err = getPISessionVarsFromUser(&pisv)
	if err != nil {
		return nil, err
	}

	// Save variables to disk.
	err = savePISessionVars(&pisv)
	if err != nil {
		return nil, err
	}

	c.APIKey = pisv.APIKey

	bxSess, err := bxsession.New(&bluemix.Config{
		BluemixAPIKey: pisv.APIKey,
	})
	if err != nil {
		return nil, err
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

	accClient, err := accountv2.New(bxSess)
	if err != nil {
		return nil, err
	}

	c.AccountAPIV2 = accClient.Accounts()
	c.Session.Config.Region = powervs.Regions[pisv.Region].VPCRegion
	return c, nil
}

//GetAccountType func return the type of account TRAIL/PAID
func (c *BxClient) GetAccountType() (string, error) {
	myAccount, err := c.AccountAPIV2.Get((*c.User).Account)
	if err != nil {
		return "", err
	}

	return myAccount.Type, nil
}

//ValidateAccountPermissions Checks permission for provisioning Power VS resources
func (c *BxClient) ValidateAccountPermissions() error {
	accType, err := c.GetAccountType()
	if err != nil {
		return err
	}
	if accType == "TRIAL" {
		return fmt.Errorf("account type must be of Pay-As-You-Go/Subscription type for provision Power VS resources")
	}
	return nil
}

//ValidateDhcpService checks for existing Dhcp service for the provided PowerVS cloud instance
func (c *BxClient) ValidateDhcpService(ctx context.Context, svcInsID string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	//create Power VS DHCP Client
	dhcpClient := instance.NewIBMPIDhcpClient(ctx, c.PISession, svcInsID)
	//Get all DHCP Services
	dhcpServices, err := dhcpClient.GetAll()
	if err != nil {
		return errors.Wrap(err, "failed to get DHCP service details")
	}
	if len(dhcpServices) > 0 {
		return fmt.Errorf("DHCP service already exists for provided cloud instance")
	}
	return nil
}

//ValidateCloudConnectionInPowerVSRegion counts cloud connection in PowerVS Region
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

// NewPISession updates pisession details, return error on fail
func (c *BxClient) NewPISession() error {
	var pisv PISessionVars

	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err := getPISessionVarsFromAuthFile(&pisv)
	if err != nil {
		return err
	}

	var authenticator core.Authenticator = &core.IamAuthenticator{
		ApiKey: c.APIKey,
	}

	// Create the session
	options := &ibmpisession.IBMPIOptions{
		Authenticator: authenticator,
		UserAccount:   c.User.Account,
		Region:        pisv.Region,
		Zone:          pisv.Zone,
		Debug:         false,
	}

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

func getPISessionVarsFromAuthFile(pisv *PISessionVars) error {

	if pisv == nil {
		return errors.New("nil var: PISessionVars")
	}

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	if _, err := os.Stat(authFilePath); os.IsNotExist(err) {
		return nil
	}

	content, err := ioutil.ReadFile(authFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, pisv)
	if err != nil {
		return err
	}

	return nil
}

func getPISessionVarsFromEnv(pisv *PISessionVars) error {

	if pisv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(pisv.ID) == 0 {
		pisv.ID = os.Getenv("IBMID")
	}

	if len(pisv.APIKey) == 0 {
		// APIKeyEnvVars is a list of environment variable names containing an IBM Cloud API key.
		var APIKeyEnvVars = []string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}
		pisv.APIKey = getEnv(APIKeyEnvVars)
	}

	if len(pisv.Region) == 0 {
		var regionEnvVars = []string{"IBMCLOUD_REGION", "IC_REGION"}
		pisv.Region = getEnv(regionEnvVars)
	}

	if len(pisv.Zone) == 0 {
		var zoneEnvVars = []string{"IBMCLOUD_ZONE"}
		pisv.Zone = getEnv(zoneEnvVars)
	}

	return nil
}

func getPISessionVarsFromUser(pisv *PISessionVars) error {
	var err error

	if pisv == nil {
		return errors.New("nil var: PiSessionVars")
	}

	if len(pisv.ID) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "IBM Cloud User ID",
					Help:    "The login for \nhttps://cloud.ibm.com/",
				},
			},
		}, &pisv.ID)
		if err != nil {
			return errors.New("error saving the IBM Cloud User ID")
		}

	}

	if len(pisv.APIKey) == 0 {
		err = survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Password{
					Message: "IBM Cloud API Key",
					Help:    "The API key installation.\nhttps://cloud.ibm.com/iam/apikeys",
				},
			},
		}, &pisv.APIKey)
		if err != nil {
			return errors.New("error saving the API Key")
		}

	}

	if len(pisv.Region) == 0 {
		pisv.Region, err = GetRegion()
		if err != nil {
			return err
		}

	}

	if len(pisv.Zone) == 0 {
		pisv.Zone, err = GetZone(pisv.Region)
		if err != nil {
			return err
		}
	}

	return nil
}

func savePISessionVars(pisv *PISessionVars) error {

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	jsonVars, err := json.Marshal(*pisv)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(authFilePath), 0700)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(authFilePath, jsonVars, 0600)
}

func getEnv(envs []string) string {
	for _, k := range envs {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
