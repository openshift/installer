// Package aws collects AWS-specific configuration.
package aws

//go:generate ../../../../hack/generate-aws-permissions.sh

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	ccaws "github.com/openshift/cloud-credential-operator/pkg/aws"
	credvalidator "github.com/openshift/cloud-credential-operator/pkg/controller/utils"
)

// ValidateCreds will try to create an AWS session, and also verify that the current credentials
// are sufficient to perform an installation, and that they can be used for cluster runtime
// as either capable of creating new credentials for components that interact with the cloud or
// being able to be passed through as-is to the components that need cloud credentials
func ValidateCreds(ssn *session.Session) error {
	creds, err := ssn.Config.Credentials.Get()
	if err != nil {
		return errors.Wrap(err, "getting creds from session")
	}

	client, err := ccaws.NewClient([]byte(creds.AccessKeyID), []byte(creds.SecretAccessKey))
	if err != nil {
		return errors.Wrap(err, "initialize cloud-credentials client")
	}

	exists := struct{}{}
	installerPermissionsMap := make(map[string]struct{}, len(createPermissions)+len(destroyPermissions))
	for _, permission := range createPermissions {
		installerPermissionsMap[permission] = exists
	}
	for _, permission := range destroyPermissions {
		installerPermissionsMap[permission] = exists
	}
	installerPermissions := make([]string, 0, len(installerPermissionsMap))
	for permission := range installerPermissionsMap {
		installerPermissions = append(installerPermissions, permission)
	}

	// Check whether we can do an installation
	logger := logrus.StandardLogger()
	canInstall, err := credvalidator.CheckPermissionsAgainstActions(client, installerPermissions, logger)
	if err != nil {
		return errors.Wrap(err, "checking install permissions")
	}
	if !canInstall {
		return errors.New("current credentials insufficient for performing cluster installation")
	}

	// Check whether we can mint new creds for cluster services needing to interact with the cloud
	canMint, err := credvalidator.CheckCloudCredCreation(client, logger)
	if err != nil {
		return errors.Wrap(err, "mint credentials check")
	}
	if canMint {
		return nil
	}

	// Check whether we can use the current credentials in passthrough mode to satisfy
	// cluster services needing to interact with the cloud
	canPassthrough, err := credvalidator.CheckCloudCredPassthrough(client, logger)
	if err != nil {
		return errors.Wrap(err, "passthrough credentials check")
	}
	if canPassthrough {
		return nil
	}

	return errors.New("AWS credentials cannot be used to either create new creds or use as-is")
}
