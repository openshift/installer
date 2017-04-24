package cloudforms

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func maybeAwsErr(e error) error {
	if e == nil {
		return nil
	}
	if awsErr, ok := e.(awserr.Error); ok {
		return fmt.Errorf("%s", awsErr.Message())
	}
	return e
}
