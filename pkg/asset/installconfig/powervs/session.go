package powervs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	gohttp "net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/bluemix-go"
	"github.com/IBM-Cloud/bluemix-go/api/account/accountv2"
	"github.com/IBM-Cloud/bluemix-go/authentication"
	"github.com/IBM-Cloud/bluemix-go/http"
	"github.com/IBM-Cloud/bluemix-go/rest"
	bxsession "github.com/IBM-Cloud/bluemix-go/session"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/form3tech-oss/jwt-go"
	"github.com/sirupsen/logrus"

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
	AccountAPIV2         accountv2.Accounts
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
	PowerVSResourceGroup string `json:"resourcegroup,omitempty"`
}

// SessionVars is an object that holds the variables required to create an ibmpisession object.
type SessionVars struct {
	ID                   string
	APIKey               string
	Region               string
	Zone                 string
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

	accClient, err := accountv2.New(bxSess)
	if err != nil {
		return nil, err
	}

	c.AccountAPIV2 = accClient.Accounts()
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
	ss.PowerVSResourceGroup = sv.PowerVSResourceGroup

	// Save the session store to the disk.
	err = saveSessionStoreToAuthFile(&ss)
	if err != nil {
		return sv, err
	}

	return sv, nil
}

// GetAccountType func return the type of account TRAIL/PAID
func (c *BxClient) GetAccountType() (string, error) {
	myAccount, err := c.AccountAPIV2.Get((*c.User).Account)
	if err != nil {
		return "", err
	}

	return myAccount.Type, nil
}

// ValidateAccountPermissions Checks permission for provisioning Power VS resources
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
