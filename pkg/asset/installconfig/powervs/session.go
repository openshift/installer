package powervs

import (
        "encoding/json"
        "io/ioutil"
	"path/filepath"
	"os"
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

// piSessionVars is an object that holds the variables required to create an ibmpisession object
type PISessionVars struct {
	ID	string `json:"id,omitempty"`
	ApiKey 	string `json:"apikey,omitempty"`
	Region 	string `json:"region,omitempty"`
	Zone 	string `json:"zone,omitempty"`
}

// GetSession returns an IBM Cloud session by using credentials found in default locations in order:
// env IBMID & env IBMID_PASSWORD,
// ~/.bluemix/config.json ? (see TODO below)
// and, if no creds are found, asks for them
/* @TODO: if you do an `ibmcloud login` (or in my case ibmcloud login --sso), you get
//  a very nice creds file at ~/.bluemix/config.json, with an IAMToken. There's no username,
//  though (just the account's owner id, but that's not the same). It may be necessary
//  to use the IAMToken vs the password env var mentioned here:
//  https://github.com/IBM-Cloud/power-go-client#ibm-cloud-sdk-for-power-cloud
//  Yes, I think we'll need to use the IAMToken. There's a two-factor auth built into the ibmcloud login,
//  so the password alone isn't enough. The IAMToken is generated as a result. So either:
     1) require the user has done this already and pull from the file
     2) ask the user to paste in their IAMToken.
     3) let the password env var be the IAMToken? (Going with this atm since it's how I started)
     4) put it into Platform {userid: , iamtoken: , ...}
*/
func GetSession() (*Session, error) {
	s, err := getPISession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}
	return &Session{Session: s}, nil
}

var (
        defaultAuthFilePath = filepath.Join(os.Getenv("HOME"), ".powervs", "config.json")
	defaultBluemixFilePath = filepath.Join(os.Getenv("HOME"), ".bluemix", "config.json")
)

func getPISession() (*ibmpisession.IBMPISession, error) {

	var err error
	var pisv PISessionVars

	// First grab variables from the installer written authFilePath
	logrus.Debug("Gathering variables from AuthFile")
	err = getPISessionVarsFromAuthFile( &pisv )
	if err != nil {
                return nil, err
        }

	// Second grab variables from files

	// .bluemix/config.json doesn't seem to hold anything useful but other files might
	// region can be found in this file, but may not be the region the user would like to use
	// zone cannot be found in this file
	// username cannot be found in this file
	// apikey cannot be found in this file, iamtoken is not the apikey we need
	
	// Third grab variables from the users enviornment
	logrus.Debug("Gathering variables from user enviornment")
	err = getPISessionVarsFromEnv( &pisv )
	if err != nil {
                return nil, err
        }

	// Fourth prompt the user for the remaining variables
	logrus.Debug("Gathering variables from user")
	err = getPISessionVarsFromUser( &pisv )
	if err != nil {
                return nil, err
        }
	
	// Save variables to disk
	err = savePISessionVars( &pisv )
	if err != nil {
                return nil, err
        }
	
	// This is needed by ibmcloud code to gather dns information later
	os.Setenv("IC_API_KEY", pisv.ApiKey)
	// We are using the iamtoken field to hold the api key
	iamtoken := pisv.ApiKey
	
	s, err := ibmpisession.New( iamtoken, pisv.Region, false, defSessionTimeout, pisv.ID, pisv.Zone )
	if err != nil {
		return nil, err
	}
	
	return s, err
	
}

func getPISessionVarsFromAuthFile( pisv *PISessionVars ) error {

	authFilePath := defaultAuthFilePath
        if f := os.Getenv( "POWERVS_AUTH_FILEPATH" ); len(f) > 0 {
                authFilePath = f
        }
	
	// Check if AuthFile exists, return if it does not
	if _, err := os.Stat( authFilePath ); os.IsNotExist(err) {
		return nil
	}
	
	content, err := ioutil.ReadFile( authFilePath )
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(content, pisv)
	if err != nil {
		return err
	}
	
	return nil
}

func getPISessionVarsFromEnv( pisv *PISessionVars ) error {
	
	if len( pisv.ID ) == 0 {
		pisv.ID = os.Getenv("IBMID")
	}
	
	if len( pisv.ApiKey ) == 0 {
	        // APIKeyEnvVars is a list of environment variable names containing an IBM Cloud API key
        	var APIKeyEnvVars = []string{"IC_API_KEY", "IBMCLOUD_API_KEY", "BM_API_KEY", "BLUEMIX_API_KEY"}
        	pisv.ApiKey = getEnv(APIKeyEnvVars)
	}

	if len( pisv.Region ) == 0 {
		var regionEnvVars = []string{"IBMCLOUD_REGION", "IC_REGION"}
        	pisv.Region = getEnv(regionEnvVars)
	}	
	
	if len( pisv.Zone ) == 0 {
		var zoneEnvVars = []string{"IBMCLOUD_ZONE"}
	        pisv.Zone = getEnv(zoneEnvVars)
	}
	
	return nil
}

func getPISessionVarsFromUser( pisv *PISessionVars ) error {
	var err error

	if len( pisv.ID ) == 0 {
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
	
	if len( pisv.ApiKey ) == 0 {
                err = survey.Ask([]*survey.Question{
                        {
                                Prompt: &survey.Password{
                                        Message: "IBM Cloud API Key",
                                        Help:    "The api key installation.\nhttps://cloud.ibm.com/iam/apikeys",
                                },
                        },
                }, &pisv.ApiKey)
                if err != nil {
                        return errors.New("Error saving the API Key")
                }
	
	}
	
	if len( pisv.Region ) == 0 {
		pisv.Region, err = GetRegion()
                if err != nil {
                        return err
                }
		
	}

        if len( pisv.Zone ) == 0 {
                pisv.Zone, err = GetZone(pisv.Region)
                if err != nil {
                        return err
                }
        }	
	
	return nil
}

func savePISessionVars( pisv *PISessionVars) error {
	
	authFilePath := defaultAuthFilePath
        if f := os.Getenv( "POWERVS_AUTH_FILEPATH" ); len(f) > 0 {
                authFilePath = f
        }

        jsonVars, err := json.Marshal( *pisv )
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
