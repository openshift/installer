package powervs

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"

	"github.com/sirupsen/logrus"
)

var (
	//reqAuthEnvs = []string{"IBMID", "IBMID_PASSWORD"}
	//optAuthEnvs = []string{"IBMCLOUD_REGION", "IBMCLOUD_ZONE"}
	//debug = false
	defSessionTimeout time.Duration = 9000000000000000000.0
	defRegion                       = "us_south"
)

// Session is an object representing a session for the IBM Power VS API.
// A bluemix session object may be a better fit here
type Session struct {
	Session *ibmpisession.IBMPISession
}

// PISessionVars is an object that holds the variables required to create an ibmpisession object
type PISessionVars struct {
	ID     string `json:"id,omitempty"`
	APIKey string `json:"apikey,omitempty"`
	Region string `json:"region,omitempty"`
	Zone   string `json:"zone,omitempty"`
}

// GetSession returns an ibmpisession object
func GetSession() (*Session, error) {
	s, err := getPISession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}
	return &Session{Session: s}, nil
}

var (
	defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".powervs", "config.json")
)

func getPISession() (*ibmpisession.IBMPISession, error) {

	var err error
	var pisv PISessionVars

	// Grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err = getPISessionVarsFromAuthFile(&pisv)
	if err != nil {
		return nil, err
	}

	// Frab variables from the users environment
	logrus.Debug("Gathering variables from user environment")
	err = getPISessionVarsFromEnv(&pisv)
	if err != nil {
		return nil, err
	}

	// Prompt the user for the remaining variables
	logrus.Debug("Gathering variables from user")
	err = getPISessionVarsFromUser(&pisv)
	if err != nil {
		return nil, err
	}

	// Save variables to disk
	err = savePISessionVars(&pisv)
	if err != nil {
		return nil, err
	}

	// This is needed by ibmcloud code to gather dns information later
	os.Setenv("IC_API_KEY", pisv.APIKey)

	// We are using the iamtoken field to hold the api key
	s, err := ibmpisession.New(pisv.APIKey, pisv.Region, false, defSessionTimeout, pisv.ID, pisv.Zone)
	if err != nil {
		return nil, err
	}

	return s, err

}

func getPISessionVarsFromAuthFile(pisv *PISessionVars) error {

	if pisv == nil {
		return errors.New("PISession Variable Object pointer cannot be nil")
	}

	authFilePath := defaultAuthFilePath
	if f := os.Getenv("POWERVS_AUTH_FILEPATH"); len(f) > 0 {
		authFilePath = f
	}

	// Check if AuthFile exists, return if it does not
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
		return errors.New("PISession Variable Object pointer cannot be nil")
	}

	if len(pisv.ID) == 0 {
		pisv.ID = os.Getenv("IBMID")
	}

	if len(pisv.APIKey) == 0 {
		// APIKeyEnvVars is a list of environment variable names containing an IBM Cloud API key
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
		return errors.New("PISession Variable Object pointer cannot be nil")
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
			return errors.New("Error saving the IBM Cloud User ID")
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
			return errors.New("Error saving the API Key")
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
	logrus.Debug("Saving variables to ", authFilePath)
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
