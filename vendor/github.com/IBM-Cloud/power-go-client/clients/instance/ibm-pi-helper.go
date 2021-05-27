package instance

import (
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
)

/*
Helper methods that will be used by the client classes
*/

// IBMPIHelperClient ...
type IBMPIHelperClient struct {
	session         *ibmpisession.IBMPISession
	powerinstanceid string
}

// NewIBMPIHelperClient ...
func NewIBMPIHelperClient(sess *ibmpisession.IBMPISession, powerinstanceid string) *IBMPIHelperClient {
	return &IBMPIHelperClient{
		session:         sess,
		powerinstanceid: powerinstanceid,
	}
}
