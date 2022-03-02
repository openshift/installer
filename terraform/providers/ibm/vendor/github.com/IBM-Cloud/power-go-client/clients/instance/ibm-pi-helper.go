package instance

import (
	"context"

	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	"github.com/go-openapi/runtime"
)

/*
Helper methods that will be used by the client classes
*/

// IBMPIHelperClient ...
type IBMPIClient struct {
	session         *ibmpisession.IBMPISession
	cloudInstanceID string
	authInfo        runtime.ClientAuthInfoWriter
	ctx             context.Context
}

// NewIBMPIClient ...
func NewIBMPIClient(ctx context.Context, sess *ibmpisession.IBMPISession, cloudInstanceID string) *IBMPIClient {
	authInfo := ibmpisession.NewAuth(sess, cloudInstanceID)
	return &IBMPIClient{
		session:         sess,
		cloudInstanceID: cloudInstanceID,
		authInfo:        authInfo,
		ctx:             ctx,
	}
}
