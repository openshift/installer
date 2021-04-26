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
    defRegion = "us_south"
)

// Session is an object representing a session for the IBM Power VS API.
type Session struct {
	session *ibmpisession.IBMPISession
}

// GetSession returns an IBM Cloud session by using credentials found in default locations in order:
// env IBMID,
// env IBMID_PASSWORD,
// and, if no creds are found, asks for them
func GetSession() (*Session, error) {
	s, err := getPISession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load credentials")
	}

    return &Session{session:s}, nil
}


func getPISession() (*ibmpisession.IBMPISession, error) {

    var (
        id, passwd, region, zone  string
        //passwd string
        //region string
        //zone   string
    )

	if id = os.Getenv("IBMID"); len(id) == 0 {
        return nil, errors.New("empty IBMID environment variable")
    }
    if passwd = os.Getenv("IBMID_PASSWORD"); len(passwd) == 0 {
        return nil, errors.New("empty IBMID_PASSWORD variable")
    }

    region = os.Getenv("IBMCLOUD_REGION")
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
	return ibmpisession.New(passwd,region, false, defSessionTimeout, id, zone)
}
