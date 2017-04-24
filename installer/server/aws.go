package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// getAWSSession returns an AWS client session which should be re-used across
// service clients and cached. It is safe for concurrent use, but not
// modification.
func getAWSSession(accessKeyID, secretAccessKey, sessionToken, region string) (*session.Session, error) {
	// create an AWS client Config
	creds := credentials.NewStaticCredentials(accessKeyID, secretAccessKey, sessionToken)
	awsConfig := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(region).
		WithCredentialsChainVerboseErrors(true)
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}
	return sess, nil
}
