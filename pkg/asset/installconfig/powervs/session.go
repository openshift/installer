package powervs

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

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

/*
//  https://github.com/IBM-Cloud/power-go-client/blob/master/ibmpisession/ibmpowersession.go
*/
func getPISession() (*ibmpisession.IBMPISession, error) {

	var (
		id, passwd, region, zone string
	)

	if id = os.Getenv("IBMID"); len(id) == 0 {
		return nil, errors.New("empty IBMID environment variable")
	}
	if passwd = os.Getenv("IBMID_PASSWORD"); len(passwd) == 0 {
		return nil, errors.New("empty IBMID_PASSWORD variable")
	}

	region = os.Getenv("IBMCLOUD_REGION")
	// this can also be pulled from  ~/bluemix/config.json
	if r2 := os.Getenv("IC_REGION"); len(r2) > 0 {
		if len(region) > 0 && region != r2 {
			return nil, errors.New(fmt.Sprintf("conflicting values for IBM Cloud Region: IBMCLOUD_REGION: %s and IC_REGION: %s", region, r2))
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
	return ibmpisession.New(passwd, region, false, defSessionTimeout, id, zone)
}
