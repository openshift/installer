package powervs

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
)

var (
	//reqAuthEnvs = []string{"IBMID", "IBMID_PASSWORD"}
	//optAuthEnvs = []string{"IBMCLOUD_REGION", "IBMCLOUD_ZONE"}
	//debug = false
	defSessionTimeout time.Duration = 9000000000000000000.0
	defRegion                       = "us_south"
)

// Session is an object representing a session for the IBM Power VS API.
type Session struct {
	Session *ibmpisession.IBMPISession
	Creds   *UserCredentials
}

// UserCredentials is an object representing the credentials used for IBM Power VS during
// the creation of the install_config.yaml
type UserCredentials struct {
	APIKey string
	UserID string
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
	s, uc, err := getPISession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}

	return &Session{Session: s, Creds: uc}, nil
}

/*
//  https://github.com/IBM-Cloud/power-go-client/blob/master/ibmpisession/ibmpowersession.go
*/
func getPISession() (*ibmpisession.IBMPISession, *UserCredentials, error) {

	var (
		id, passwd, apikey, region, zone string
	)

	if id = os.Getenv("IBMID"); len(id) == 0 {
		err := survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Input{
					Message: "IBM Cloud User ID",
					Help:    "The login for \nhttps://cloud.ibm.com/",
				},
			},
		}, &id)
		if err != nil {
			return nil, nil, errors.New("Error saving the IBMID variable")
		}
	}

	if apikey = os.Getenv("API_KEY"); len(apikey) == 0 {
		err := survey.Ask([]*survey.Question{
			{
				Prompt: &survey.Password{
					Message: "IBM Cloud API Key",
					Help:    "The api key installation.\nhttps://cloud.ibm.com/iam/apikeys",
				},
			},
		}, &apikey)
		if err != nil {
			return nil, nil, errors.New("Error saving the API_KEY variable")
		}
	}

	region = os.Getenv("IBMCLOUD_REGION")
	// this can also be pulled from  ~/bluemix/config.json
	if r2 := os.Getenv("IC_REGION"); len(r2) > 0 {
		if len(region) > 0 && region != r2 {
			return nil, nil, errors.New(fmt.Sprintf("conflicting values for IBM Cloud Region: IBMCLOUD_REGION: %s and IC_REGION: %s", region, r2))
		}
		if len(region) == 0 {
			region = r2
		}
	}

	if zone = os.Getenv("IBMCLOUD_ZONE"); len(zone) == 0 {
		zone = region
	}

	// @TOOD: query if region is multi-zone? or just pass through err...
	// @TODO: pass through debug?
	s, err := ibmpisession.New(passwd, region, false, defSessionTimeout, id, zone)
	uc := &UserCredentials{UserID: id, APIKey: apikey}
	return s, uc, err
}
